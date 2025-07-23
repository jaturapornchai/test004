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

	if len(ohlcv) > 0 {
		ema20 := ai.calculateEMA(ohlcv, 20)
		ema50 := ai.calculateEMA(ohlcv, 50)
		ema100 := ai.calculateEMA(ohlcv, 100)
		ema200 := ai.calculateEMA(ohlcv, 200)

		prompt += fmt.Sprintf("\n\nEMA Values:\nEMA20: %.6f\nEMA50: %.6f\nEMA100: %.6f\nEMA200: %.6f\n",
			ema20, ema50, ema100, ema200)

		alignment := ai.checkEMAAlignment(ohlcv)
		touch := ai.checkEMA200Touch(ohlcv)
		prompt += fmt.Sprintf("EMA Alignment: %s\nEMA200 Touch (48 candles): %s\nTime Frame: 1 hour (244 candles, ~10 days)\n", alignment, touch)
	}

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

	prompt := string(promptBytes)
	prompt += fmt.Sprintf("\n\nContract: %s\nPosition Side: %s\nPosition Size: %d\nCandles: %d\n\n",
		position.Contract, positionSide, position.Size, len(ohlcv))

	for i, candle := range ohlcv {
		prompt += fmt.Sprintf("%d: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
			i+1, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}

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
	k := 2.0 / 201.0
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
