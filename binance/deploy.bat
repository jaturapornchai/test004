@echo off
REM GateIO Trading Bot One-Command Deploy Script for Windows

echo 🚀 Starting GateIO Trading Bot deployment...

REM Build and push image
echo 📦 Building and pushing Docker image...
docker buildx build --platform linux/amd64 --no-cache -t jaturapornchai/gateio:latest --push .

if %errorlevel% equ 0 (
    echo ✅ Docker image built and pushed successfully!
) else (
    echo ❌ Failed to build/push image
    exit /b 1
)

REM Deploy to server
echo 🚢 Deploying to server...
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/binance && echo '📥 Pulling latest image...' && sudo docker pull jaturapornchai/getspot:latest && echo '🔄 Restarting bot...' && sudo docker-compose down && sudo docker-compose up -d && echo '✅ Bot deployed successfully!' && echo '📊 Showing logs (Ctrl+C to exit):' && sudo docker logs -f binance"

echo 🎉 Deployment completed!
