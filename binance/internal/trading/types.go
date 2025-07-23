package trading

// OHLCV ข้อมูลแท่งเทียน
type OHLCV struct {
	Timestamp int64   `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

// Position ข้อมูล position
type Position struct {
	Contract      string  `json:"contract"`
	Size          int64   `json:"size"`
	EntryPrice    float64 `json:"entry_price"`
	MarkPrice     float64 `json:"mark_price"`
	UnrealizedPnl float64 `json:"unrealized_pnl"`
	Margin        float64 `json:"margin"`
	Leverage      float64 `json:"leverage"`
	Mode          string  `json:"mode"`
}

// AIDecision การตัดสินใจของ AI
type AIDecision struct {
	Action          string  `json:"action"`      // "LONG", "SHORT", "CLOSE", "HOLD"
	RiskRewardRatio float64 `json:"risk_reward"` // อัตราส่วน Risk:Reward
	Reason          string  `json:"reason"`      // เหตุผล
	StopLoss        float64 `json:"stop_loss"`   // แนะนำ SL
	TakeProfit      float64 `json:"take_profit"` // แนะนำ TP
}

// TechnicalData ข้อมูลทางเทคนิค (Pure AI Mode)
type TechnicalData struct {
	ATR    []float64 `json:"atr"`
	EMA100 []float64 `json:"ema100"`
}
