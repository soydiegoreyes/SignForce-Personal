# Deployment — SignForce

## Server Details
- **Host:** Skynet
- **IP:** 187.77.6.110
- **OS:** Ubuntu 24
- **Access:** `ssh root@187.77.6.110`
- **Project path:** `/root/projects/signforce/` (symlinked from `/root/signforce/`)

## Current Deployment (Frontend Only)

### Architecture
```
Internet → Nginx (host, SSL) → signforce-web container (Nginx Alpine, :8090)
```

### Domain
- **Active:** `testsignforce.luxspace.org` (SSL via Let's Encrypt)
- **Pending DNS:** `app.sign-force.com` (A record → 187.77.6.110)

### Deploy / Update
```bash
cd ~/signforce && git pull && docker compose -f docker-compose.frontend.yml up -d --build
```

### Container Details
- Image: `nginx:alpine` with sffront copied to `/usr/share/nginx/html/`
- Port mapping: `127.0.0.1:8090:80` (NEVER change this port)
- Restart policy: `unless-stopped`

### Nginx Config (Host)
Location: `/etc/nginx/sites-available/signforce`
```
server {
    listen 443 ssl;
    server_name testsignforce.luxspace.org;
    ssl_certificate /etc/letsencrypt/live/testsignforce.luxspace.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/testsignforce.luxspace.org/privkey.pem;
    
    location / {
        proxy_pass http://127.0.0.1:8090;
    }
}
```

## Full Stack Deployment (MySQL + Go + Nginx)

### Prerequisites
- Docker and Docker Compose
- `.env` file configured (copy from `.env.example`)

### Steps
```bash
cp .env.example .env
# Edit .env with real values for DB_PASS, JWT_KEY, DOCS_KEY, etc.
nano .env

docker compose up -d
```

### Services Started
| Service | Port | Container |
|---------|------|-----------|
| MySQL 8.0 | 3306 | signforce-mysql |
| sfback | 5001 | signforce-back |
| sfmiddle | 5002 | signforce-middle |
| Nginx (frontend) | 8090 | signforce-web |

## SSL Certificate Renewal
```bash
certbot renew --nginx
```

## Backups
- **Automatic:** Daily at 3:00 AM via cron
- **Location:** `/root/backups/daily/`
- **Retention:** 7 days
- **Manual:** `bash /root/scripts/backup.sh`

## Troubleshooting

### Container won't start (port conflict)
```bash
sudo lsof -i :8090
# Kill the PID, then retry
docker compose -f docker-compose.frontend.yml up -d --build
```

### Changes not visible after deploy
- Hard refresh: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows)
- Check Nginx cache: assets cached for 7 days

### Check logs
```bash
docker logs signforce-web --tail 100
docker logs signforce-web -f  # Follow live
```
