package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get configuration from environment
	host := getEnv("CLICKHOUSE_HOST", "localhost")
	port := getEnv("CLICKHOUSE_PORT", "9000")
	user := getEnv("CLICKHOUSE_USER", "default")
	password := getEnv("CLICKHOUSE_PASSWORD", "")
	database := getEnv("CLICKHOUSE_DATABASE", "default")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", user, password, host, port, database)

	fmt.Printf("🔗 Connecting to ClickHouse...\n")
	fmt.Printf("📍 Host: %s:%s\n", host, port)
	fmt.Printf("👤 User: %s\n", user)
	fmt.Printf("🗄️  Database: %s\n", database)
	fmt.Println()

	// Connect to ClickHouse
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open connection: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping ClickHouse: %v", err)
	}

	fmt.Println("✅ Successfully connected to ClickHouse!")

	// Test query - get ClickHouse version
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("❌ Failed to get version: %v", err)
	}
	fmt.Printf("📊 ClickHouse Version: %s\n", version)

	// List all tables
	fmt.Println("\n📋 Available tables:")
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("❌ Failed to list tables: %v", err)
	}
	defer rows.Close()

	tableCount := 0
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("⚠️  Error scanning table: %v", err)
			continue
		}
		tableCount++
		fmt.Printf("  %d. %s\n", tableCount, tableName)
	}

	if tableCount == 0 {
		fmt.Println("  (No tables found)")
	} else {
		fmt.Printf("\n📈 Total tables: %d\n", tableCount)
	}

	// Test simple query on a table (if exists)
	if tableCount > 0 {
		fmt.Println("\n🔍 Testing query on first table...")

		// Get first table name
		var firstTable string
		err = db.QueryRow("SELECT name FROM system.tables WHERE database = ? LIMIT 1", database).Scan(&firstTable)
		if err == nil {
			fmt.Printf("📊 Querying table: %s\n", firstTable)

			// Count rows
			var count int64
			err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", firstTable)).Scan(&count)
			if err != nil {
				fmt.Printf("⚠️  Error counting rows: %v\n", err)
			} else {
				fmt.Printf("📊 Total rows in %s: %d\n", firstTable, count)
			}

			// Show first 3 rows
			fmt.Printf("\n📝 Sample data from %s:\n", firstTable)
			rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 3", firstTable))
			if err != nil {
				fmt.Printf("⚠️  Error querying sample data: %v\n", err)
			} else {
				defer rows.Close()

				// Get column names
				columns, err := rows.Columns()
				if err != nil {
					fmt.Printf("⚠️  Error getting columns: %v\n", err)
				} else {
					fmt.Printf("📋 Columns: %v\n", columns)

					rowNum := 0
					for rows.Next() && rowNum < 3 {
						// Create a slice of interface{} to scan into
						values := make([]interface{}, len(columns))
						valuePtrs := make([]interface{}, len(columns))
						for i := range columns {
							valuePtrs[i] = &values[i]
						}

						if err := rows.Scan(valuePtrs...); err != nil {
							fmt.Printf("⚠️  Error scanning row: %v\n", err)
							continue
						}

						rowNum++
						fmt.Printf("  Row %d: ", rowNum)
						for i, val := range values {
							if i > 0 {
								fmt.Print(", ")
							}
							fmt.Printf("%v", val)
						}
						fmt.Println()
					}
				}
			}
		}
	}

	fmt.Println("\n🎉 Connection test completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
