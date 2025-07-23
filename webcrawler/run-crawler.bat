@echo off
echo ===============================================
echo    Thai Web Crawler for Email Extraction
echo ===============================================
echo.
echo Starting Thai web crawler...
echo Press Ctrl+C to stop the crawler
echo.

if "%1"=="" (
    echo Using default Thai websites...
    go run main.go
) else (
    echo Using custom URLs: %*
    go run main.go %*
)

echo.
echo Crawler stopped. Check found_emails.json for results.
pause
