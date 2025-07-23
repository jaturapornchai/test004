@echo off
echo Starting Product API Server...
echo.
echo Server will be available at: http://localhost:8080
echo API Documentation: http://localhost:8080
echo.
echo Press Ctrl+C to stop the server
echo.

cd /d "%~dp0"
go run main.go

pause
