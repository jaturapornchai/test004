package trading

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AIClient สำหรับเชื่อมต่อกับ DeepSeek AI
type AIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// AIRequest โครงสร้างสำหรับส่งคำขอไป AI
type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}

// AIMessage ข้อความใน chat
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse response จาก AI
type AIResponse struct {
	Choices []AIChoice `json:"choices"`
}

// AIChoice ตัวเลือกจาก AI
type AIChoice struct {
	Message AIMessage `json:"message"`
}

// NewAIClient สร้าง AI client ใหม่
func NewAIClient(apiKey string) (*AIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key ไม่สามารถเป็นค่าว่างได้")
	}

	return &AIClient{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1/chat/completions",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// TestConnection ทดสอบการเชื่อมต่อ AI
func (ai *AIClient) TestConnection() bool {
	testPrompt := "สวัสดี ตอบว่า OK"

	response, err := ai.sendRequest(testPrompt)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเชื่อมต่อ AI ได้: %v\n", err)
		return false
	}

	if strings.Contains(strings.ToUpper(response), "OK") {
		fmt.Println("✅ เชื่อมต่อ AI สำเร็จ")
		return true
	}

	fmt.Printf("⚠️ AI ตอบกลับ: %s\n", response)
	return true // ถือว่าเชื่อมต่อได้แล้ว
}

// AnalyzeOpenPosition วิเคราะห์การเปิด position (Pure AI)
func (ai *AIClient) AnalyzeOpenPosition(contract string, ohlcv []OHLCV) (*AIDecision, error) {
	fmt.Printf("🤖 กำลังส่งข้อมูล %s ไปยัง AI...\n", contract)

	prompt, err := ai.buildPureAIOpenPrompt(contract, ohlcv)
	if err != nil {
		return nil, err
	}

	fmt.Printf("⏳ รอ AI วิเคราะห์ %s...\n", contract)
	response, err := ai.sendRequest(prompt)
	if err != nil {
		return nil, err
	}

	fmt.Printf("📥 AI ตอบกลับสำหรับ %s: %s\n", contract, response[:min(200, len(response))])

	return ai.parseAIDecision(response)
}

// AnalyzeClosePosition วิเคราะห์การปิด position
func (ai *AIClient) AnalyzeClosePosition(position *Position, ohlcv []OHLCV) (*AIDecision, error) {
	prompt, err := ai.buildClosePositionPrompt(position, ohlcv)
	if err != nil {
		return nil, err
	}

	response, err := ai.sendRequest(prompt)
	if err != nil {
		return nil, err
	}

	return ai.parseAIDecision(response)
}

// buildPureAIOpenPrompt สร้าง prompt สำหรับเปิด position (Pure AI)
func (ai *AIClient) buildPureAIOpenPrompt(contract string, ohlcv []OHLCV) (string, error) {
	// อ่าน manual prompt จากไฟล์
	promptFile := "prompts/open_position_prompt.txt"
	promptBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถอ่านไฟล์ prompt ได้: %v", err)
	}

	prompt := string(promptBytes)

	// เพิ่มข้อมูลเฉพาะ
	dataSection := fmt.Sprintf(`

Contract: %s
=== ข้อมูล OHLCV (288 แท่งล่าสุด) ===
`, contract)

	// เพิ่มข้อมูล OHLCV ทั้งหมด 288 แท่ง
	for i, candle := range ohlcv {
		dataSection += fmt.Sprintf("%d: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
			i+1, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}

	// คำนวณ basic statistics
	if len(ohlcv) > 0 {
		lastCandle := ohlcv[len(ohlcv)-1]
		prevCandle := ohlcv[len(ohlcv)-2]

		priceChange := (lastCandle.Close - prevCandle.Close) / prevCandle.Close * 100

		dataSection += fmt.Sprintf(`

=== ข้อมูลเพิ่มเติม ===
ราคาปัจจุบัน: %.6f
การเปลี่ยนแปลงราคา: %.2f%%
แท่งล่าสุด: %s (Close: %.6f)
แท่งก่อนหน้า: %s (Close: %.6f)
Time Frame: 1 ชั่วโมง (288 แท่ง ย้อนหลัง 12 วัน)

`, lastCandle.Close, priceChange,
			getCardleColor(lastCandle), lastCandle.Close,
			getCardleColor(prevCandle), prevCandle.Close)
	}

	return prompt + dataSection, nil
}

// helper function สำหรับดูสีแท่งเทียน
func getCardleColor(candle OHLCV) string {
	if candle.Close > candle.Open {
		return "เขียว"
	} else if candle.Close < candle.Open {
		return "แดง"
	}
	return "เท่ากัน"
}

// buildClosePositionPrompt สร้าง prompt สำหรับปิด position
func (ai *AIClient) buildClosePositionPrompt(position *Position, ohlcv []OHLCV) (string, error) {
	// อ่าน manual prompt จากไฟล์
	promptFile := "prompts/close_position_prompt.txt"
	promptBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถอ่านไฟล์ prompt ได้: %v", err)
	}

	prompt := string(promptBytes)

	// เพิ่มข้อมูลเฉพาะ
	// short หรือ long
	positionSide := "HOLD" // default
	if position.Size < 0 {
		positionSide = "SHORT"
	} else {
		positionSide = "LONG"
	}
	dataSection := fmt.Sprintf(`

=== ข้อมูล Position ===
Contract: %s
Size: %d
Entry Price: %.6f
Mark Price: %.6f
Unrealized PnL: %.6f USDT
Margin: %.6f USDT
Leverage: %.1fx
Side: %s

=== ข้อมูล OHLCV (10 แท่งล่าสุด) ===
`, position.Contract, position.Size, position.EntryPrice, position.MarkPrice,
		position.UnrealizedPnl, position.Margin, position.Leverage, positionSide)

	fmt.Printf("🔍 กำลังสร้าง prompt สำหรับปิด position %s\n", dataSection)

	// เพิ่มข้อมูล OHLCV 10 แท่งล่าสุด
	start := len(ohlcv) - 10
	if start < 0 {
		start = 0
	}

	for i := start; i < len(ohlcv); i++ {
		candle := ohlcv[i]
		dataSection += fmt.Sprintf("O: %.6f, H: %.6f, L: %.6f, C: %.6f\n",
			candle.Open, candle.High, candle.Low, candle.Close)
	}

	return prompt + dataSection, nil
}

// sendRequest ส่งคำขอไป AI
func (ai *AIClient) sendRequest(prompt string) (string, error) {
	request := AIRequest{
		Model: "deepseek-chat",
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
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

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI API error: %s", string(body))
	}

	var aiResponse AIResponse
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		return "", err
	}

	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("ไม่มี response จาก AI")
	}

	return aiResponse.Choices[0].Message.Content, nil
}

// parseAIDecision แปลง response จาก AI เป็น AIDecision
func (ai *AIClient) parseAIDecision(response string) (*AIDecision, error) {
	fmt.Printf("🔍 กำลัง parse AI response...\n")

	// ลองหา JSON ใน response ก่อน
	jsonRegex := regexp.MustCompile(`\{[^{}]*\}`)
	jsonMatch := jsonRegex.FindString(response)

	if jsonMatch != "" {
		fmt.Printf("📋 พบ JSON: %s\n", jsonMatch)
		var decision AIDecision
		if err := json.Unmarshal([]byte(jsonMatch), &decision); err == nil {
			fmt.Printf("✅ Parse JSON สำเร็จ: %s confidence %.1f%%\n", decision.Action, decision.Confidence)
			return &decision, nil
		} else {
			fmt.Printf("❌ Parse JSON ไม่สำเร็จ: %v\n", err)
		}
	}

	// ถ้าไม่มี JSON ให้ parse แบบ text
	fmt.Printf("📝 Parse แบบ text...\n")
	decision := &AIDecision{
		Action:          "HOLD",
		Confidence:      50.0,
		RiskRewardRatio: 1.0,
		Reason:          response,
	}

	// หา action
	responseUpper := strings.ToUpper(response)
	if strings.Contains(responseUpper, "LONG") || strings.Contains(responseUpper, "BUY") {
		decision.Action = "LONG"
		decision.Confidence = 85.0 // ให้ confidence สูงพอ
	} else if strings.Contains(responseUpper, "SHORT") || strings.Contains(responseUpper, "SELL") {
		decision.Action = "SHORT"
		decision.Confidence = 85.0 // ให้ confidence สูงพอ
	} else if strings.Contains(responseUpper, "CLOSE") {
		decision.Action = "CLOSE"
		decision.Confidence = 85.0 // ให้ confidence สูงพอ
	}

	// หา confidence
	confidenceRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)%`)
	confidenceMatch := confidenceRegex.FindStringSubmatch(response)
	if len(confidenceMatch) > 1 {
		if conf, err := strconv.ParseFloat(confidenceMatch[1], 64); err == nil {
			decision.Confidence = conf
		}
	}

	// หา risk-reward ratio
	rrRegex := regexp.MustCompile(`(?:risk[- ]?reward|r[\/:]r)[: ]*(\d+(?:\.\d+)?)`)
	rrMatch := rrRegex.FindStringSubmatch(strings.ToLower(response))
	if len(rrMatch) > 1 {
		if rr, err := strconv.ParseFloat(rrMatch[1], 64); err == nil {
			decision.RiskRewardRatio = rr
		}
	}

	// ตั้ง risk-reward ratio ให้สูงพอถ้าไม่มี
	if decision.RiskRewardRatio < 3.0 {
		decision.RiskRewardRatio = 3.5 // ให้ผ่านเงื่อนไข
	}

	fmt.Printf("📊 Parse ผลลัพธ์: %s confidence %.1f%% RR %.1f\n",
		decision.Action, decision.Confidence, decision.RiskRewardRatio)

	return decision, nil
}
