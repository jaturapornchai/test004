package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// EmailDatabase represents a simple file-based email database
type EmailDatabase struct {
	filePath string
	emails   map[string]EmailData
	mu       sync.RWMutex
	nextID   int
}

// NewEmailDatabase creates a new email database
func NewEmailDatabase(filePath string) *EmailDatabase {
	db := &EmailDatabase{
		filePath: filePath,
		emails:   make(map[string]EmailData),
		nextID:   1,
	}
	db.loadFromFile()
	return db
}

// loadFromFile loads emails from JSON file
func (db *EmailDatabase) loadFromFile() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, err := os.Stat(db.filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, that's okay
	}

	data, err := ioutil.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	var emailList []EmailData
	if err := json.Unmarshal(data, &emailList); err != nil {
		return err
	}

	for _, email := range emailList {
		db.emails[email.Email] = email
		if email.ID >= db.nextID {
			db.nextID = email.ID + 1
		}
	}

	return nil
}

// saveToFile saves emails to JSON file
func (db *EmailDatabase) saveToFile() error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var emailList []EmailData
	for _, email := range db.emails {
		emailList = append(emailList, email)
	}

	data, err := json.MarshalIndent(emailList, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(db.filePath, data, 0644)
}

// AddEmail adds an email to the database
func (db *EmailDatabase) AddEmail(email string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.emails[email]; exists {
		return false // Email already exists
	}

	emailData := EmailData{
		ID:    db.nextID,
		Email: email,
	}

	db.emails[email] = emailData
	db.nextID++

	return true // New email added
}

// GetCount returns the total number of emails
func (db *EmailDatabase) GetCount() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.emails)
}

// GetAllEmails returns all emails as a slice
func (db *EmailDatabase) GetAllEmails() []EmailData {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var emails []EmailData
	for _, email := range db.emails {
		emails = append(emails, email)
	}
	return emails
}

// EmailData represents an email found during crawling (simplified for personal emails only)
type EmailData struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// Crawler represents the web crawler
type Crawler struct {
	visitedURLs  map[string]bool
	emailPattern *regexp.Regexp
	mu           sync.Mutex
	maxDepth     int
	delay        time.Duration
	db           *EmailDatabase
}

// NewCrawler creates a new crawler instance with database
func NewCrawler(maxDepth int, delay time.Duration, dbPath string) *Crawler {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// Create email database
	db := NewEmailDatabase(dbPath)

	return &Crawler{
		visitedURLs:  make(map[string]bool),
		emailPattern: emailRegex,
		maxDepth:     maxDepth,
		delay:        delay,
		db:           db,
	}
}

// isValidURL checks if the URL is valid and should be crawled (Thai websites only)
func (c *Crawler) isValidURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Only crawl HTTP and HTTPS URLs
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Only allow Thai domains and Thai-related international domains
	host := strings.ToLower(parsedURL.Host)

	// Thai domains (.th, .co.th, .ac.th, etc.)
	if strings.HasSuffix(host, ".th") {
		return c.isValidPath(parsedURL.Path)
	}

	// Popular Thai websites on international domains
	thaiSites := []string{
		"sanook.com", "pantip.com", "kapook.com", "mthai.com",
		"thairath.co.th", "manager.co.th", "bangkokpost.com",
		"nationthailand.com", "thaipbs.or.th", "mcot.net",
		"komchadluek.net", "khaosod.co.th", "posttoday.com",
		"workpointtoday.com", "ch3thailand.com", "tnn24.com",
		"thaienquirer.com", "siamzone.com", "dek-d.com",
		"jeban.com", "bloggang.com", "thaivisa.com",
	}

	for _, site := range thaiSites {
		if strings.Contains(host, site) {
			return c.isValidPath(parsedURL.Path)
		}
	}

	return false
}

// isValidPath checks if the URL path should be crawled
func (c *Crawler) isValidPath(path string) bool {
	// Skip common file extensions that are not web pages
	skipExtensions := []string{".pdf", ".jpg", ".jpeg", ".png", ".gif", ".zip", ".rar", ".exe", ".doc", ".docx", ".mp4", ".avi", ".mp3", ".wav"}
	for _, ext := range skipExtensions {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return false
		}
	}
	return true
}

// hasThaiContent checks if the text contains Thai content
func (c *Crawler) hasThaiContent(text string) bool {
	thaiCharCount := 0
	totalChars := 0

	for _, r := range text {
		if r >= 0x0E00 && r <= 0x0E7F { // Thai Unicode range
			thaiCharCount++
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= 0x0E00 && r <= 0x0E7F) {
			totalChars++
		}
	}

	// Must have at least 10 Thai characters and Thai content should be at least 20% of total text
	if thaiCharCount < 10 {
		return false
	}

	if totalChars == 0 {
		return false
	}

	thaiRatio := float64(thaiCharCount) / float64(totalChars)
	return thaiRatio >= 0.2 // At least 20% Thai content
}

func (c *Crawler) normalizeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Remove fragment and query parameters for URL comparison
	parsedURL.Fragment = ""
	parsedURL.RawQuery = ""

	return parsedURL.String()
}

// isPersonalEmail checks if an email is from a personal email provider
func (c *Crawler) isPersonalEmail(email string) bool {
	personalDomains := []string{
		"gmail.com", "hotmail.com", "yahoo.com", "outlook.com",
		"live.com", "msn.com", "icloud.com", "me.com",
		"protonmail.com", "tutanota.com", "yandex.com",
		"mail.com", "aol.com", "zoho.com", "fastmail.com",
		"gmx.com", "rambler.ru", "rambler.ua", "inbox.com",
		"yahoo.co.th", "yahoo.co.uk", "yahoo.ca", "yahoo.fr",
		"hotmail.co.th", "hotmail.co.uk", "hotmail.fr",
		"live.co.th", "live.co.uk", "live.fr", "live.ca",
		"rediffmail.com", "mailinator.com", "guerrillamail.com",
		"10minutemail.com", "temp-mail.org", "maildrop.cc",
	}

	emailLower := strings.ToLower(email)
	for _, domain := range personalDomains {
		if strings.HasSuffix(emailLower, "@"+domain) {
			return true
		}
	}
	return false
}

func (c *Crawler) extractEmails(text, pageURL string) []string {
	emails := c.emailPattern.FindAllString(text, -1)

	// Remove duplicates and filter personal emails only
	emailMap := make(map[string]bool)
	var personalEmails []string

	for _, email := range emails {
		email = strings.ToLower(email)
		// Skip if already found or not a personal email
		if emailMap[email] || !c.isPersonalEmail(email) {
			continue
		}
		emailMap[email] = true
		personalEmails = append(personalEmails, email)
	}

	return personalEmails
}

// saveEmailToDB saves an email to the database
func (c *Crawler) saveEmailToDB(email string) bool {
	return c.db.AddEmail(email)
}

// getEmailCount returns the total number of emails in the database
func (c *Crawler) getEmailCount() int {
	return c.db.GetCount()
}

// getAllEmails returns all emails from the database
func (c *Crawler) getAllEmails() []EmailData {
	return c.db.GetAllEmails()
}

// crawlPage crawls a single page and extracts emails and links
func (c *Crawler) crawlPage(pageURL string, depth int) {
	if depth > c.maxDepth {
		return
	}

	normalizedURL := c.normalizeURL(pageURL)

	c.mu.Lock()
	if c.visitedURLs[normalizedURL] {
		c.mu.Unlock()
		return
	}
	c.visitedURLs[normalizedURL] = true
	c.mu.Unlock()

	fmt.Printf("Crawling: %s (depth: %d)\n", pageURL, depth)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(pageURL)
	if err != nil {
		log.Printf("Error fetching %s: %v", pageURL, err)
		return
	}
	defer resp.Body.Close()

	// Only process HTML content
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML from %s: %v", pageURL, err)
		return
	}

	// Extract emails from page text
	pageText := doc.Text()

	// Skip pages that don't have Thai content
	if !c.hasThaiContent(pageText) {
		fmt.Printf("Skipping %s - No Thai content detected\n", pageURL)
		return
	}

	emails := c.extractEmails(pageText, pageURL)

	if len(emails) > 0 {
		fmt.Printf("Found %d personal emails on %s\n", len(emails), pageURL)

		c.mu.Lock()
		newEmailCount := 0
		for _, email := range emails {
			if c.saveEmailToDB(email) {
				newEmailCount++
			}
		}
		c.mu.Unlock()

		if newEmailCount > 0 {
			fmt.Printf("Added %d new emails to database\n", newEmailCount)
			// Save to file after adding new emails
			if err := c.db.saveToFile(); err != nil {
				log.Printf("Error saving database to file: %v", err)
			}
		}
	}

	// Extract links for further crawling
	var links []string
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Convert relative URLs to absolute
		absoluteURL, err := url.Parse(href)
		if err != nil {
			return
		}

		baseURL, err := url.Parse(pageURL)
		if err != nil {
			return
		}

		fullURL := baseURL.ResolveReference(absoluteURL).String()

		if c.isValidURL(fullURL) {
			links = append(links, fullURL)
		}
	})

	// Add delay to be respectful to the server
	time.Sleep(c.delay)

	// Crawl found links recursively
	for _, link := range links {
		normalizedLink := c.normalizeURL(link)
		c.mu.Lock()
		visited := c.visitedURLs[normalizedLink]
		c.mu.Unlock()

		if !visited {
			go c.crawlPage(link, depth+1)
		}
	}
}

// Start begins the crawling process
func (c *Crawler) Start(startURLs []string) {
	fmt.Printf("Starting web crawler with %d initial URLs\n", len(startURLs))
	fmt.Printf("Max depth: %d, Delay: %v\n", c.maxDepth, c.delay)
	fmt.Println("Database: personal_emails.json")

	var wg sync.WaitGroup

	for _, startURL := range startURLs {
		if c.isValidURL(startURL) {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				c.crawlPage(url, 0)
			}(startURL)
		}
	}

	// Wait for initial URLs to be processed
	wg.Wait()

	// Keep the program running to allow background goroutines to continue
	fmt.Println("Initial crawling completed. Continuing background crawling...")

	// Periodic status report
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			emailCount := c.getEmailCount()

			c.mu.Lock()
			visitedCount := len(c.visitedURLs)
			c.mu.Unlock()

			fmt.Printf("Status: Found %d unique personal emails, Visited %d URLs\n", emailCount, visitedCount)
		}
	}
}

// showDatabaseStats shows database statistics and recent emails
func (c *Crawler) showDatabaseStats() {
	count := c.getEmailCount()

	fmt.Printf("\n=== Database Statistics ===\n")
	fmt.Printf("Total personal emails found: %d\n", count)

	// Show recent 10 emails
	emails := c.getAllEmails()

	fmt.Println("\nPersonal emails in database:")
	maxShow := 10
	if len(emails) < maxShow {
		maxShow = len(emails)
	}

	for i := 0; i < maxShow; i++ {
		fmt.Printf("%d. %s\n", i+1, emails[i].Email)
	}

	if len(emails) > 10 {
		fmt.Printf("... and %d more emails\n", len(emails)-10)
	}
	fmt.Println("=============================\n")
}

func main() {
	// Configuration
	maxDepth := 3                    // Maximum crawling depth
	delay := 2 * time.Second         // Delay between requests
	dbPath := "personal_emails.json" // JSON database file

	// Starting URLs - Thai websites
	startURLs := []string{
		"https://pantip.com",
		"https://sanook.com",
		"https://kapook.com",
		"https://mthai.com",
		"https://thairath.co.th",
		"https://manager.co.th",
		"https://siamzone.com",
		"https://dek-d.com",
		// Add more Thai URLs here
	}

	// Check if custom URLs are provided via command line
	if len(os.Args) > 1 {
		startURLs = os.Args[1:]
		fmt.Printf("Using URLs from command line: %v\n", startURLs)
	}

	// Create and start crawler
	crawler := NewCrawler(maxDepth, delay, dbPath)

	fmt.Println("Web Crawler Started!")
	fmt.Printf("Personal emails will be saved to JSON database: %s\n", dbPath)
	fmt.Println("Press Ctrl+C to stop the crawler")

	crawler.Start(startURLs)
}
