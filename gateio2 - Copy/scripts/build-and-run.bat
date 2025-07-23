@echo off
echo Building Gate.io Trading Bot...
go build -o gateio-trading-bot.exe
if %ERRORLEVEL% EQU 0 (
    echo Build successful!
    echo Running bot...
    .\gateio-trading-bot.exe
) else (
    echo Build failed!
    pause
)
