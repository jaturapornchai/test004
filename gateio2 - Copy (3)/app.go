package main

import (
	"fmt"
	"log"
	"os"

	"gateio-trading-bot/internal/trading"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 เริ่มต้น Gate.io Trading Bot - Pure AI Trading System V-003")

	// ตรวจสอบ command line arguments
	if len(os.Args) > 1 && os.Args[1] == "prompts" {
		fmt.Println("📋 แสดง AI Prompts แบบละเอียด...")
		displayPromptsOnly()
		return
	}

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ ไม่สามารถโหลดไฟล์ .env ได้:", err)
	}

	// ตรวจสอบ API credentials
	apiKey := os.Getenv("GATE_API_KEY")
	apiSecret := os.Getenv("GATE_API_SECRET")
	deepseekKey := os.Getenv("DEEPSEEK_API_KEY")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("❌ GATE_API_KEY หรือ GATE_API_SECRET ไม่ได้ตั้งค่าใน .env")
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
	fmt.Println("🔍 ทดสอบการเชื่อมต่อ Gate.io และ AI...")
	if !bot.TestConnections() {
		log.Fatal("❌ การเชื่อมต่อไม่สำเร็จ")
	}

	fmt.Println("✅ การเชื่อมต่อทั้งหมดสำเร็จ")
	fmt.Println("🔄 เริ่มต้นระบบ Trading Loop...")

	// แสดง prompt เสมอก่อนเริ่มเทรด
	fmt.Println("📋 แสดง AI Prompts ที่ใช้ในระบบ...")
	bot.DisplayDetailedPrompts()
	fmt.Println("⏰ รอ 3 วินาทีก่อนเริ่มเทรด...")
	fmt.Println("📢 กด Ctrl+C เพื่อหยุดระบบ")
	fmt.Println()

	// เริ่ม main trading loop
	bot.StartTradingLoop()
}

// displayPromptsOnly แสดง AI Prompts เท่านั้น (ไม่รันระบบเทรด)
func displayPromptsOnly() {
	// สร้าง dummy AI client เพื่อแสดง prompts
	aiClient, err := trading.NewAIClient("dummy")
	if err != nil {
		log.Fatal("❌ ไม่สามารถสร้าง AI client ได้:", err)
	}

	aiClient.DisplayDetailedPrompts()

	fmt.Println("💡 วิธีการใช้งาน:")
	fmt.Println("   go run app.go          - รันระบบเทรดปกติ")
	fmt.Println("   go run app.go prompts  - แสดง AI prompts เท่านั้น")
}
