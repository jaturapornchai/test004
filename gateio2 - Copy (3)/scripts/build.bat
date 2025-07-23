@echo off
echo Building Gate.io Trading Bot...
go build -o gateio-trading-bot.exe
if %ERRORLEVEL% EQU 0 (
    echo Build successful!
) else (
    echo Build failed!
    pause
)
