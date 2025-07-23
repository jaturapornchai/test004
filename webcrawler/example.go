package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Example usage of the web crawler
func main() {
	fmt.Println("=== Web Crawler Example ===")

	// Example URLs to start crawling
	// You can replace these with actual websites you want to crawl
	testURLs := []string{
		"https://httpbin.org/html",
		"https://example.com",
		// Add more URLs here as needed
	}

	fmt.Println("This example shows how to use the web crawler.")
	fmt.Println("The main crawler is in main.go")
	fmt.Println()
	fmt.Println("To run the actual crawler, use:")
	fmt.Println("  go run main.go")
	fmt.Println()
	fmt.Println("Or with custom URLs:")
	for _, url := range testURLs {
		fmt.Printf("  go run main.go %s\n", url)
	}
	fmt.Println()
	fmt.Println("The crawler will:")
	fmt.Println("1. Visit the starting URLs")
	fmt.Println("2. Extract all links from each page")
	fmt.Println("3. Follow those links (up to max depth)")
	fmt.Println("4. Extract email addresses from page content")
	fmt.Println("5. Save found emails to found_emails.json")
	fmt.Println("6. Skip already visited URLs")
	fmt.Println()
	fmt.Println("Press Enter to run a quick test...")

	var input string
	fmt.Scanln(&input)

	// Quick demonstration
	fmt.Println("Running quick test...")

	// Create a test crawler with limited depth
	crawler := NewCrawlerExample(1, 3*time.Second, "test_emails.json")

	// Test with a simple URL
	testURL := "https://httpbin.org/html"
	fmt.Printf("Testing with: %s\n", testURL)

	go crawler.TestCrawl(testURL)

	// Wait a bit for the test
	time.Sleep(10 * time.Second)

	fmt.Println("\nTest completed! Check test_emails.json for results.")
	fmt.Println("For full crawling, run: go run main.go")
}

// Simplified crawler for demonstration
type CrawlerExample struct {
	maxDepth   int
	delay      time.Duration
	outputFile string
}

func NewCrawlerExample(maxDepth int, delay time.Duration, outputFile string) *CrawlerExample {
	return &CrawlerExample{
		maxDepth:   maxDepth,
		delay:      delay,
		outputFile: outputFile,
	}
}

func (c *CrawlerExample) TestCrawl(url string) {
	fmt.Printf("Test crawling: %s\n", url)
	fmt.Println("(This is just a demo - check main.go for the full implementation)")

	// Create a simple test result
	testData := `[
  {
    "email": "test@example.com",
    "url": "` + url + `",
    "date": "` + time.Now().Format("2006-01-02 15:04:05") + `"
  }
]`

	err := os.WriteFile(c.outputFile, []byte(testData), 0644)
	if err != nil {
		log.Printf("Error writing test file: %v", err)
		return
	}

	fmt.Printf("Test results saved to: %s\n", c.outputFile)
}
