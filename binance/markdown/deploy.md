# Binance Trading Bot Deployment

## Build and push Docker image

**⚠️ Important: Edit .env with real API keys before building!**

```bash
# 1. Edit .env with your real API keys first:
# BINANCE_API_KEY=your_real_binance_api_key
# BINANCE_SECRET_KEY=your_real_binance_secret_key  
# AI_API_KEY=your_real_deepseek_api_key

# 2. Then build and push
docker buildx build --platform linux/amd64 --no-cache -t jaturapornchai/binance:latest --push .
```

## One-command deploy to server

```bash
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/binance && sudo docker pull jaturapornchai/binance:latest && sudo docker-compose down && sudo docker-compose up -d && sudo docker logs -f binance"
```

## Alternative: Deploy with environment setup

```bash
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/binance && \
sudo docker pull jaturapornchai/binance:latest && \
sudo docker stop binance 2>/dev/null || true && \
sudo docker rm binance 2>/dev/null || true && \
sudo docker run -d --name binance \
  -v /mnt/volume_sgp1_02/binance/logs:/root/logs \
  --restart unless-stopped \
  jaturapornchai/binance:latest && \
sudo docker logs -f binance"
```

## Manual step-by-step (if needed)

```bash
ssh root@178.128.55.234
cd /mnt/volume_sgp1_02/binance
sudo docker pull jaturapornchai/binance:latest
sudo docker-compose down
sudo docker-compose up -d
sudo docker ps
sudo docker logs -f binance
```
