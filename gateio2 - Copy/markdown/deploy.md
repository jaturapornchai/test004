# GateIO Trading Bot Deployment

## Build and push Docker image

**⚠️ Important: Edit .env with real API keys before building!**

```bash
# 1. Edit .env with your real API keys first:
# GATE_API_KEY=your_real_gate_api_key
# GATE_SECRET_KEY=your_real_gate_secret_key  
# AI_API_KEY=your_real_deepseek_api_key

# 2. Then build and push
docker buildx build --platform linux/amd64 --no-cache -t jaturapornchai/gateio:latest --push .
```

## One-command deploy to server

```bash
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/gateio && sudo docker pull jaturapornchai/gateio:latest && sudo docker-compose down && sudo docker-compose up -d && sudo docker logs -f gateio"
```

## Alternative: Deploy with environment setup

```bash
ssh root@178.128.55.234 "cd /mnt/volume_sgp1_02/gateio && \
sudo docker pull jaturapornchai/gateio:latest && \
sudo docker stop gateio 2>/dev/null || true && \
sudo docker rm gateio 2>/dev/null || true && \
sudo docker run -d --name gateio \
  --env-file .env \
  -v /mnt/volume_sgp1_02/gateio/logs:/root/logs \
  --restart unless-stopped \
  jaturapornchai/gateio:latest && \
sudo docker logs -f gateio"
```

## Manual step-by-step (if needed)

```bash
ssh root@178.128.55.234
cd /mnt/volume_sgp1_02/gateio
sudo docker pull jaturapornchai/gateio:latest
sudo docker-compose down
sudo docker-compose up -d
sudo docker ps
sudo docker logs -f gateio
```
