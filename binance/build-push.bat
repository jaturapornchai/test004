@echo off
REM Binance Trading Bot - Build and Push Only

echo 🚀 Building and pushing Binance Trading Bot...

REM Build and push image
echo 📦 Building Docker image...
docker buildx build --platform linux/amd64 --no-cache -t jaturapornchai/getspot:latest --push .

if %errorlevel% equ 0 (
    echo ✅ Docker image built and pushed successfully!
    echo.
    echo 🚢 To deploy on server, run this command:
    echo ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/binance && sudo docker pull jaturapornchai/getspot:latest && sudo docker-compose down && sudo docker-compose up -d && sudo docker logs -f binance"
    echo.
    echo Or simply run: deploy-server.bat
) else (
    echo ❌ Failed to build/push image
    exit /b 1
)

pause
