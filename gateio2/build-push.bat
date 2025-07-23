@echo off
REM GateIO Trading Bot - Build and Push Only

echo 🚀 Building and pushing GateIO Trading Bot...

REM Build and push image
echo 📦 Building Docker image...
docker buildx build --platform linux/amd64 --no-cache -t jaturapornchai/gateio:latest --push .

if %errorlevel% equ 0 (
    echo ✅ Docker image built and pushed successfully!
    echo.
    echo 🚢 To deploy on server, run this command:
    echo ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/gateio && sudo docker pull jaturapornchai/gateio:latest && sudo docker-compose down && sudo docker-compose up -d && sudo docker logs -f gateio"
    echo.
    echo Or simply run: deploy-server.bat
) else (
    echo ❌ Failed to build/push image
    exit /b 1
)

pause
