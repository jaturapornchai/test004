# Start Product API Server
Write-Host "Starting Product API Server..." -ForegroundColor Green
Write-Host ""
Write-Host "Server will be available at: http://localhost:8080" -ForegroundColor Yellow
Write-Host "API Documentation: http://localhost:8080" -ForegroundColor Yellow
Write-Host ""
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Cyan
Write-Host ""

Set-Location $PSScriptRoot
go run main.go
