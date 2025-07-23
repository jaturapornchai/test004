@echo off
REM GateIO Trading Bot - Server Deploy Only

echo 🚢 Deploying to server...
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/binance && echo '📥 Pulling latest image...' && sudo docker pull jaturapornchai/getspot:latest && echo '🔄 Restarting bot...' && sudo docker-compose down && sudo docker-compose up -d && echo '✅ Bot deployed successfully!' && echo '📊 Showing logs (Ctrl+C to exit):' && sudo docker logs -f binance"

echo 🎉 Deployment completed!
