#!/bin/bash

# GateIO Trading Bot One-Command Deploy Script
echo "🚀 Starting GateIO Trading Bot deployment..."

# Build and push image
echo "📦 Building and pushing Docker image..."
docker buildx build --platform linux/amd64 --no-cache -t jaturapornchai/gateiobot:latest --push .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built and pushed successfully!"
else
    echo "❌ Failed to build/push image"
    exit 1
fi

# Deploy to server
echo "🚢 Deploying to server..."
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/gateio && \
echo '📥 Pulling latest image...' && \
sudo docker pull jaturapornchai/gateiobot:latest && \
echo '🔄 Restarting bot...' && \
sudo docker-compose down && \
sudo docker-compose up -d && \
echo '✅ Bot deployed successfully!' && \
echo '📊 Showing logs (Ctrl+C to exit):' && \
sudo docker logs -f gateiobot"

echo "🎉 Deployment completed!"
