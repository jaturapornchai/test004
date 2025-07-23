package main

import (
	"fmt"
	"log"
	"os"

	"binance-trading-bot/internal/trading"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 เริ่มต้น Binance Trading Bot - Pure AI Trading System V-004")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ ไม่สามารถโหลดไฟล์ .env ได้:", err)
	}

	// ตรวจสอบ API credentials
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	deepseekKey := os.Getenv("DEEPSEEK_API_KEY")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("❌ BINANCE_API_KEY หรือ BINANCE_API_SECRET ไม่ได้ตั้งค่าใน .env")
	}

	if deepseekKey == "" {
		log.Fatal("❌ DEEPSEEK_API_KEY ไม่ได้ตั้งค่าใน .env")
	}

	fmt.Println("✅ โหลด configuration สำเร็จ")

	// สร้าง trading bot instance
	bot, err := trading.NewTradingBot(apiKey, apiSecret, deepseekKey)
	if err != nil {
		log.Fatal("❌ ไม่สามารถสร้าง trading bot ได้:", err)
	}

	// ทดสอบการเชื่อมต่อ
	fmt.Println("🔍 ทดสอบการเชื่อมต่อ Binance และ AI...")
	if !bot.TestConnections() {
		log.Fatal("❌ การเชื่อมต่อไม่สำเร็จ")
	}

	fmt.Println("✅ การเชื่อมต่อทั้งหมดสำเร็จ")
	fmt.Println("🔄 เริ่มต้นระบบ Trading Loop...")

	// เริ่ม main trading loop
	bot.Start()
}
