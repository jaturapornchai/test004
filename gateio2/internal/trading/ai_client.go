package trading

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type AIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewAIClient(apiKey string) (*AIClient, error) {
	return &AIClient{
		apiKey:     apiKey,
		baseURL:    "https://api.deepseek.com/v1/chat/completions",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}, nil
}

func (ai *AIClient) TestConnection() bool {
	return true
}

func (ai *AIClient) AnalyzeOpenPosition(contract string, ohlcv []OHLCV) (*AIDecision, error) {
	prompt, err := ai.buildOpenPrompt(contract, ohlcv)
	if err != nil {
		return nil, err
	}

	// แสดง real prompt ที่ส่งไป AI
	fmt.Println("🤖 ======= REAL PROMPT ส่งไป AI (OPEN POSITION) =======")
	fmt.Printf("📄 Contract: %s\n", contract)
	fmt.Printf("📊 Candles: %d แท่ง\n", len(ohlcv))
	fmt.Println("───────────────────────────────────────────────────────")
	// แสดง prompt แต่จำกัดความยาว (แสดงแค่ 10 แท่งแรกและ 5 แท่งสุดท้าย)
	ai.displayTruncatedPrompt(prompt, len(ohlcv))
	fmt.Println("════════════════════════════════════════════════════════")
	fmt.Println()

	response, err := ai.makeRequest(prompt)
	if err != nil {
		return nil, err
	}
	return ai.parseAIDecision(response)
}

func (ai *AIClient) AnalyzeClosePosition(position *Position, ohlcv []OHLCV) (*AIDecision, error) {
	prompt, err := ai.buildClosePrompt(position, ohlcv)
	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่าเป็น EMA force close หรือไม่
	if prompt == "FORCE_CLOSE_EMA50_BELOW" {
		fmt.Printf("⚠️  EMA50 FORCE CLOSE: %s (LONG) - ราคาต่ำกว่า EMA50\n", position.Contract)
		return &AIDecision{
			Action:     "CLOSE",
			Confidence: 100.0,
			Reason:     "ราคาปัจจุบันต่ำกว่า EMA50 - บังคับปิด LONG position",
		}, nil
	}

	if prompt == "FORCE_CLOSE_EMA50_ABOVE" {
		fmt.Printf("⚠️  EMA50 FORCE CLOSE: %s (SHORT) - ราคาสูงกว่า EMA50\n", position.Contract)
		return &AIDecision{
			Action:     "CLOSE",
			Confidence: 100.0,
			Reason:     "ราคาปัจจุบันสูงกว่า EMA50 - บังคับปิด SHORT position",
		}, nil
	}

	if prompt == "FORCE_CLOSE_EMA_MISALIGNMENT" {
		positionType := "LONG"
		if position.Size < 0 {
			positionType = "SHORT"
		}
		fmt.Printf("⚠️  EMA ALIGNMENT FORCE CLOSE: %s (%s) - EMA alignment หลุด\n", position.Contract, positionType)

		return &AIDecision{
			Action:     "CLOSE",
			Confidence: 100.0,
			Reason:     "EMA alignment เปลี่ยน - บังคับปิด position",
		}, nil
	}

	// แสดง real prompt ที่ส่งไป AI
	fmt.Println("🤖 ======= REAL PROMPT ส่งไป AI (CLOSE POSITION) =======")
	fmt.Printf("📄 Contract: %s\n", position.Contract)
	positionType := "LONG"
	if position.Size < 0 {
		positionType = "SHORT"
	}
	fmt.Printf("📊 Position: %s (Size: %d)\n", positionType, position.Size)
	fmt.Printf("📊 Candles: %d แท่ง\n", len(ohlcv))
	fmt.Println("───────────────────────────────────────────────────────")
	ai.displayTruncatedPrompt(prompt, len(ohlcv))
	fmt.Println("════════════════════════════════════════════════════════")
	fmt.Println()

	response, err := ai.makeRequest(prompt)
	if err != nil {
		return nil, err
	}
	return ai.parseAIDecision(response)
}

func (ai *AIClient) buildOpenPrompt(contract string, ohlcv []OHLCV) (string, error) {
	promptBytes, err := os.ReadFile("prompts/open_position_prompt.txt")
	if err != nil {
		return "", err
	}

	prompt := string(promptBytes)
	prompt += fmt.Sprintf("\n\nContract: %s\nCandles: %d\n\n", contract, len(ohlcv))

	for i, candle := range ohlcv {
		prompt += fmt.Sprintf("%d: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
			i+1, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}

	// ไม่ส่ง EMA values ให้ AI - ให้ AI คำนวณและวิเคราะห์เอง
	prompt += "\nTime Frame: 1 hour (200 candles, ~8 days)\n"

	return prompt, nil
}

func (ai *AIClient) buildClosePrompt(position *Position, ohlcv []OHLCV) (string, error) {
	promptBytes, err := os.ReadFile("prompts/close_position_prompt.txt")
	if err != nil {
		return "", err
	}

	positionSide := "HOLD"
	if position.Size < 0 {
		positionSide = "SHORT"
	} else if position.Size > 0 {
		positionSide = "LONG"
	}

	// คำนวณ EMA values สำหรับการตรวจสอบ
	if len(ohlcv) >= 200 {
		ema20 := ai.calculateEMA(ohlcv, 20)
		ema50 := ai.calculateEMA(ohlcv, 50)
		ema100 := ai.calculateEMA(ohlcv, 100)
		ema200 := ai.calculateEMA(ohlcv, 200)
		currentPrice := ohlcv[len(ohlcv)-1].Close

		// เงื่อนไข EMA50: ตรวจสอบก่อนถาม AI
		if positionSide == "LONG" {
			// LONG position: ถ้าราคาปัจจุบันต่ำกว่า EMA50 ให้ CLOSE
			if currentPrice < ema50 {
				return "FORCE_CLOSE_EMA50_BELOW", nil
			}
		} else if positionSide == "SHORT" {
			// SHORT position: ถ้าราคาปัจจุบันสูงกว่า EMA50 ให้ CLOSE
			if currentPrice > ema50 {
				return "FORCE_CLOSE_EMA50_ABOVE", nil
			}
		}

		// ตรวจสอบ EMA alignment และบังคับ CLOSE ถ้าไม่เรียงตัว
		if positionSide == "LONG" {
			// LONG position: ต้องการ EMA20>EMA50>EMA100>EMA200
			if !(ema20 > ema50 && ema50 > ema100 && ema100 > ema200) {
				// EMA ไม่เรียงตัวแบบ bullish → บังคับ CLOSE
				return "FORCE_CLOSE_EMA_MISALIGNMENT", nil
			}
		} else if positionSide == "SHORT" {
			// SHORT position: ต้องการ EMA20<EMA50<EMA100<EMA200
			if !(ema20 < ema50 && ema50 < ema100 && ema100 < ema200) {
				// EMA ไม่เรียงตัวแบบ bearish → บังคับ CLOSE
				return "FORCE_CLOSE_EMA_MISALIGNMENT", nil
			}
		}
	}

	prompt := string(promptBytes)
	prompt += fmt.Sprintf("\n\nContract: %s\nPosition Side: %s\nPosition Size: %d\nCandles: %d\n\n",
		position.Contract, positionSide, position.Size, len(ohlcv))

	for i, candle := range ohlcv {
		prompt += fmt.Sprintf("%d: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
			i+1, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}

	// ไม่ส่ง EMA values ให้ AI - ให้ AI คำนวณและวิเคราะห์เอง
	prompt += "\nTime Frame: 1 hour (200 candles, ~8 days)\n"

	return prompt, nil
}

func (ai *AIClient) calculateEMA(ohlcv []OHLCV, period int) float64 {
	if len(ohlcv) < period {
		return 0
	}
	k := 2.0 / float64(period+1)
	ema := ohlcv[0].Close
	for i := 1; i < len(ohlcv); i++ {
		ema = ohlcv[i].Close*k + ema*(1-k)
	}
	return ema
}

func (ai *AIClient) checkEMA200Touch(ohlcv []OHLCV) string {
	if len(ohlcv) < 200 {
		return "Insufficient data"
	}
	k := 2.0 / float64(200+1) // แก้ไข: ใช้ EMA200 period ที่ถูกต้อง
	ema := ohlcv[0].Close
	startIdx := len(ohlcv) - 48
	if startIdx < 0 {
		startIdx = 0
	}
	for i := 1; i < len(ohlcv); i++ {
		ema = ohlcv[i].Close*k + ema*(1-k)
		if i >= startIdx {
			candle := ohlcv[i]
			if candle.Low <= ema && ema <= candle.High {
				return "Yes"
			}
		}
	}
	return "No"
}

func (ai *AIClient) checkEMAAlignment(ohlcv []OHLCV) string {
	if len(ohlcv) < 200 {
		return "Mixed"
	}
	ema20 := ai.calculateEMA(ohlcv, 20)
	ema50 := ai.calculateEMA(ohlcv, 50)
	ema100 := ai.calculateEMA(ohlcv, 100)
	ema200 := ai.calculateEMA(ohlcv, 200)

	if ema20 > ema50 && ema50 > ema100 && ema100 > ema200 {
		return "LONG"
	}
	if ema20 < ema50 && ema50 < ema100 && ema100 < ema200 {
		return "SHORT"
	}
	return "Mixed"
}

func (ai *AIClient) makeRequest(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ai.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.apiKey)

	resp, err := ai.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	return response.Choices[0].Message.Content, nil
}

func (ai *AIClient) parseAIDecision(response string) (*AIDecision, error) {
	action := "HOLD"
	upperResponse := strings.ToUpper(response)

	if strings.Contains(upperResponse, "CLOSE") {
		action = "CLOSE"
	} else if strings.Contains(upperResponse, "LONG") {
		action = "LONG"
	} else if strings.Contains(upperResponse, "SHORT") {
		action = "SHORT"
	}

	return &AIDecision{
		Action:     action,
		Confidence: 70.0,
		Reason:     response,
	}, nil
}

// DisplayDetailedPrompts แสดง prompt แบบละเอียดทั้งสองแบบ
func (ai *AIClient) DisplayDetailedPrompts() {
	fmt.Println("🤖 ==================== AI PROMPTS DETAILED ====================")
	fmt.Println()

	// แสดง Open Position Prompt
	fmt.Println("🟢 OPEN POSITION PROMPT:")
	fmt.Println("📄 ไฟล์: prompts/open_position_prompt.txt")
	fmt.Println("───────────────────────────────────────────────────────────────")
	openPromptBytes, err := os.ReadFile("prompts/open_position_prompt.txt")
	if err != nil {
		fmt.Printf("❌ ข้อผิดพลาดในการอ่านไฟล์ open prompt: %v\n", err)
	} else {
		fmt.Println(string(openPromptBytes))
	}

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println()

	// แสดง Close Position Prompt
	fmt.Println("🔴 CLOSE POSITION PROMPT:")
	fmt.Println("📄 ไฟล์: prompts/close_position_prompt.txt")
	fmt.Println("───────────────────────────────────────────────────────────────")
	closePromptBytes, err := os.ReadFile("prompts/close_position_prompt.txt")
	if err != nil {
		fmt.Printf("❌ ข้อผิดพลาดในการอ่านไฟล์ close prompt: %v\n", err)
	} else {
		fmt.Println(string(closePromptBytes))
	}

	fmt.Println()
	fmt.Println("🔍 สรุปข้อมูลสำคัญ:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 ข้อมูลที่ใช้: OHLCV 200 แท่งเทียน (8+ วันย้อนหลัง)")
	fmt.Println("💰 Position Size: 15 USDT")
	fmt.Println("⚡ Leverage: 5x")
	fmt.Println("📈 Time Frame: 1 ชั่วโมง")
	fmt.Println("🎯 กลยุทธ์: Pure Price Action")
	fmt.Println("🔄 EMA Filter: 20>50>100>200 (LONG) หรือ 20<50<100<200 (SHORT)")
	fmt.Println("� Coins Filter: เหรียญใหม่ (≤ 4 เดือน) + Volume > $100K")
	fmt.Println("�📝 Response Format: JSON {action, reason}")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

// displayTruncatedPrompt แสดง prompt แบบย่อให้เห็นตัวอย่าง
func (ai *AIClient) displayTruncatedPrompt(prompt string, candleCount int) {
	lines := strings.Split(prompt, "\n")

	// หาบรรทัดที่เริ่มข้อมูล OHLCV
	startIdx := -1
	for i, line := range lines {
		if strings.Contains(line, ": O=") {
			startIdx = i
			break
		}
	}

	if startIdx == -1 {
		// ถ้าไม่เจอข้อมูล OHLCV ให้แสดงทั้งหมด
		fmt.Println(prompt)
		return
	}

	// แสดงส่วน header (prompt template)
	for i := 0; i < startIdx; i++ {
		fmt.Println(lines[i])
	}

	// แสดงข้อมูล OHLCV แค่ 10 แท่งแรกและ 5 แท่งสุดท้าย
	fmt.Println("📊 OHLCV Data (แสดง 10 แท่งแรก ... 5 แท่งสุดท้าย):")

	// แสดง 10 แท่งแรก
	endFirstBatch := startIdx + 10
	if endFirstBatch > len(lines) {
		endFirstBatch = len(lines)
	}

	for i := startIdx; i < endFirstBatch && i < len(lines); i++ {
		if strings.Contains(lines[i], ": O=") {
			fmt.Println(lines[i])
		}
	}

	// ถ้ามีมากกว่า 15 แท่ง ให้แสดง ...
	if candleCount > 15 {
		fmt.Printf("... (ซ่อน %d แท่ง) ...\n", candleCount-15)

		// แสดง 5 แท่งสุดท้าย
		startLastBatch := len(lines) - 8 // 5 แท่งสุดท้าย + time frame line
		for i := startLastBatch; i < len(lines); i++ {
			if i >= 0 && i < len(lines) && strings.Contains(lines[i], ": O=") {
				fmt.Println(lines[i])
			}
		}
	}

	// แสดงส่วนท้าย (time frame)
	for i := len(lines) - 3; i < len(lines); i++ {
		if i >= 0 && i < len(lines) && !strings.Contains(lines[i], ": O=") {
			fmt.Println(lines[i])
		}
	}
}
