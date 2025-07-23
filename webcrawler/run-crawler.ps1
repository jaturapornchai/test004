# Thai Web Crawler PowerShell Script
Write-Host "===============================================" -ForegroundColor Green
Write-Host "    Thai Web Crawler for Email Extraction" -ForegroundColor Green  
Write-Host "===============================================" -ForegroundColor Green
Write-Host ""

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Go detected: $goVersion" -ForegroundColor Yellow
}
catch {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Starting Thai web crawler..." -ForegroundColor Cyan
Write-Host "Press Ctrl+C to stop the crawler" -ForegroundColor Cyan
Write-Host ""

# Check if URLs provided as arguments
if ($args.Count -eq 0) {
    Write-Host "Using default Thai websites..." -ForegroundColor Yellow
    & go run main.go
}
else {
    Write-Host "Using custom URLs: $($args -join ', ')" -ForegroundColor Yellow
    & go run main.go @args
}

Write-Host ""
Write-Host "Crawler stopped." -ForegroundColor Green
Write-Host "Check found_emails.json for results." -ForegroundColor Green

# Show results if file exists
if (Test-Path "found_emails.json") {
    $content = Get-Content "found_emails.json" | ConvertFrom-Json
    Write-Host ""
    Write-Host "Found $($content.Count) emails total:" -ForegroundColor Green
    $content | Select-Object email, url | Format-Table -AutoSize
}

Read-Host "Press Enter to exit"
