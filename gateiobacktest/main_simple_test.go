package main

import (
	"fmt"
	"strings"
)

func mainSimpleTest() {
	fmt.Println("🚀 SIMPLE PROFITABLE STRATEGIES - REAL MARKET TEST")
	fmt.Println("💡 Finding simple strategies that work with real data")
	fmt.Println(strings.Repeat("━", 80))

	fmt.Println("\n📋 Testing 4 Simple Strategies:")
	fmt.Println("1. 📊 RSI Bounce - ตีกลับจากแนวรับ/แนวต้าน")
	fmt.Println("2. 📈 Volume Breakout - ปริมาณการซื้อขายสูง")
	fmt.Println("3. 🕯️ Price Action - รูปแบบเทียน")
	fmt.Println("4. 📉 Mean Reversion - กลับสู่ค่าเฉลี่ย")

	fmt.Println("\n💰 Risk Management:")
	fmt.Println("• Risk per trade: 0.5% (Conservative)")
	fmt.Println("• Stop Loss: 1% from entry")
	fmt.Println("• Take Profit: 2% (1:2 R/R)")
	fmt.Println("• Time Limit: 12 hours max")

	fmt.Println(strings.Repeat("━", 80))

	// ทดสอบกลยุทธ์ง่ายๆ ทั้งหมด
	TestAllSimpleStrategies()

	fmt.Println(strings.Repeat("━", 80))
	fmt.Printf("🎉 SIMPLE STRATEGIES TEST COMPLETE!\n")
	fmt.Printf("💡 Focus on the best performing strategy\n")
	fmt.Printf("🚀 Simple strategies often work better in real markets\n")
}
