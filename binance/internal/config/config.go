package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey    string
	APISecret string
	BaseURL   string

	// Trading parameters
	PositionSize float64
	Leverage     int

	// Technical parameters
	PivotPeriod int
	ATRFactor   float64
	ATRPeriod   int
	EMAPeriod   int

	// AI parameters
	AIModel   string
	AIBaseURL string

	// Risk management
	MinRiskReward float64
}

func Load() (*Config, error) {
	// โหลด .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("⚠️ ไม่สามารถโหลดไฟล์ .env ได้: %v\n", err)
		// แต่ยังคงดำเนินการต่อ เผื่อมี env vars อื่น
	} else {
		fmt.Println("✅ โหลดไฟล์ .env สำเร็จ")
	}

	cfg := &Config{
		APIKey:    os.Getenv("BINANCE_API_KEY"),
		APISecret: os.Getenv("BINANCE_API_SECRET"),
		BaseURL:   getEnvWithDefault("BINANCE_BASE_URL", "https://fapi.binance.com"),

		// Trading parameters ตามเอกสาร
		PositionSize: 15.0, // 15 USDT ต่อ position
		Leverage:     10,   // 10x leverage

		// Technical parameters ตาม Pine Script
		PivotPeriod: 2,   // Pivot Point Period
		ATRFactor:   3.0, // ATR Factor
		ATRPeriod:   10,  // ATR Period
		EMAPeriod:   100, // EMA Period

		// AI parameters
		AIModel:   getEnvWithDefault("AI_MODEL", "gpt-4"),
		AIBaseURL: getEnvWithDefault("AI_BASE_URL", "https://api.openai.com/v1"),

		// Risk management
		MinRiskReward: 3.0, // ต้อง >= 3.0
	}

	fmt.Printf("🔍 Debug - API Key length: %d\n", len(cfg.APIKey))
	fmt.Printf("🔍 Debug - API Secret length: %d\n", len(cfg.APISecret))

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("BINANCE_API_KEY ไม่ได้ถูกตั้งค่าใน .env")
	}

	if cfg.APISecret == "" {
		return nil, fmt.Errorf("BINANCE_API_SECRET ไม่ได้ถูกตั้งค่าใน .env")
	}

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
