# Authgrid Production Deployment Guide

Complete guide for deploying Authgrid to **authgrid.org** with SSL, monitoring, and production best practices.

---

## 🚀 Quick Deploy

**Prerequisites:**
- Ubuntu 20.04+ server with root access
- Docker & Docker Compose installed
- Domain pointing to your server (A records configured)

**Deploy in 5 minutes:**
```bash
# 1. Clone repository
git clone https://github.com/Kelsidavis/authgrid.git
cd authgrid

# 2. Configure environment
cp .env.production.example .env.production
nano .env.production
# Change POSTGRES_PASSWORD to a secure password

# 3. Run deployment
cd deploy
chmod +x deploy.sh
./deploy.sh
```

That's it! Your site will be live at https://authgrid.org

---

## 📋 Detailed Deployment Steps

### Step 1: Server Setup

#### 1.1 Server Requirements

**Minimum:**
- 1 CPU core
- 1GB RAM
- 20GB storage
- Ubuntu 20.04+

**Recommended:**
- 2 CPU cores
- 2GB RAM
- 40GB storage
- Ubuntu 22.04 LTS

#### 1.2 Install Docker

```bash
# Update system
sudo apt update
sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
docker compose version
```

#### 1.3 Firewall Configuration

```bash
# Allow SSH, HTTP, HTTPS
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

---

### Step 2: DNS Configuration

Configure your DNS records for **authgrid.org**:

| Type | Name | Value | TTL |
|------|------|-------|-----|
| A | @ | YOUR_SERVER_IP | 3600 |
| A | www | YOUR_SERVER_IP | 3600 |
| A | api | YOUR_SERVER_IP | 3600 |

**Wait for DNS propagation** (usually 5-30 minutes):

```bash
# Check DNS propagation
dig authgrid.org
dig www.authgrid.org
dig api.authgrid.org
```

---

### Step 3: Clone Repository

```bash
# Clone to server
git clone https://github.com/Kelsidavis/authgrid.git
cd authgrid

# Or if you already have it locally, use rsync:
rsync -avz --exclude 'node_modules' \
  /local/path/to/authgrid/ \
  user@your-server:/opt/authgrid/
```

---

### Step 4: Configuration

#### 4.1 Create Production Environment File

```bash
cp .env.production.example .env.production
nano .env.production
```

**Configure:**
```bash
# Database - CHANGE THIS PASSWORD!
POSTGRES_USER=authgrid
POSTGRES_PASSWORD=YOUR_SECURE_PASSWORD_HERE

# Domain
DOMAIN=authgrid.org

# Email for Let's Encrypt
LETSENCRYPT_EMAIL=admin@authgrid.org

# API Configuration
AUTHGRID_DOMAIN=authgrid.org
```

**Generate secure password:**
```bash
openssl rand -base64 32
```

#### 4.2 Update Demo for Production

```bash
# Copy production demo
cp examples/demo/index.production.html examples/demo/index.html
```

This updates the API URL to `https://api.authgrid.org`.

---

### Step 5: Deploy

#### 5.1 Run Deployment Script

```bash
cd deploy
./deploy.sh
```

**The script will:**
1. ✅ Build Docker images
2. ✅ Start PostgreSQL database
3. ✅ Start API server
4. ✅ Generate SSL certificates (Let's Encrypt)
5. ✅ Start Nginx reverse proxy
6. ✅ Start Certbot for auto-renewal

#### 5.2 Manual SSL Setup (Alternative)

If auto-SSL fails, run manually:

```bash
cd deploy
./init-letsencrypt.sh
```

---

### Step 6: Verify Deployment

#### 6.1 Check Services

```bash
docker compose -f docker-compose.production.yml ps
```

Expected output:
```
NAME                 STATUS
authgrid-postgres    Up (healthy)
authgrid-api         Up
authgrid-nginx       Up
authgrid-certbot     Up
```

#### 6.2 Test API

```bash
# Test health endpoint
curl https://api.authgrid.org/health

# Expected response:
# {"status":"healthy","time":"2025-10-28T..."}
```

#### 6.3 Test Website

Open in browser:
- https://authgrid.org
- https://www.authgrid.org
- https://api.authgrid.org/health

#### 6.4 Test Registration & Login

1. Visit https://authgrid.org
2. Click "Register Now"
3. Verify you get a handle
4. Click "Login"
5. Verify authentication works

---

## 🔧 Management Commands

### View Logs

```bash
# All services
docker compose -f docker-compose.production.yml logs -f

# Specific service
docker compose -f docker-compose.production.yml logs -f api
docker compose -f docker-compose.production.yml logs -f nginx
docker compose -f docker-compose.production.yml logs -f postgres
```

### Restart Services

```bash
# Restart all
docker compose -f docker-compose.production.yml restart

# Restart specific service
docker compose -f docker-compose.production.yml restart api
docker compose -f docker-compose.production.yml restart nginx
```

### Stop Services

```bash
docker compose -f docker-compose.production.yml down
```

### Update Deployment

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker compose -f docker-compose.production.yml build
docker compose -f docker-compose.production.yml up -d
```

---

## 🗄️ Database Management

### Backup Database

```bash
# Create backup
docker compose -f docker-compose.production.yml exec postgres \
  pg_dump -U authgrid authgrid > backup_$(date +%Y%m%d_%H%M%S).sql

# Backup to compressed file
docker compose -f docker-compose.production.yml exec postgres \
  pg_dump -U authgrid authgrid | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz
```

### Restore Database

```bash
# Restore from backup
cat backup.sql | docker compose -f docker-compose.production.yml exec -T postgres \
  psql -U authgrid authgrid
```

### Access Database Shell

```bash
docker compose -f docker-compose.production.yml exec postgres \
  psql -U authgrid -d authgrid
```

### Automated Backups

```bash
# Create backup script
cat > /opt/authgrid/backup.sh << 'EOF'
#!/bin/bash
cd /opt/authgrid
BACKUP_DIR="/opt/backups/authgrid"
mkdir -p $BACKUP_DIR
docker compose -f docker-compose.production.yml exec -T postgres \
  pg_dump -U authgrid authgrid | gzip > \
  $BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql.gz
# Keep only last 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
EOF

chmod +x /opt/authgrid/backup.sh

# Add to crontab (daily at 2 AM)
(crontab -l 2>/dev/null; echo "0 2 * * * /opt/authgrid/backup.sh") | crontab -
```

---

## 🔒 Security Hardening

### 1. Change Default Ports (Optional)

Edit `docker-compose.production.yml`:
```yaml
nginx:
  ports:
    - "8080:80"   # Custom HTTP port
    - "8443:443"  # Custom HTTPS port
```

### 2. Restrict API Access

In `nginx/conf.d/authgrid.conf`, add IP whitelist:
```nginx
# Only allow specific IPs to access API
location / {
    allow 1.2.3.4;      # Your IP
    allow 5.6.7.8/24;   # Your network
    deny all;

    proxy_pass http://api:8080;
    ...
}
```

### 3. Enable Fail2Ban

```bash
sudo apt install fail2ban -y

# Create nginx jail
sudo nano /etc/fail2ban/jail.local
```

Add:
```ini
[nginx-limit-req]
enabled = true
filter = nginx-limit-req
logpath = /var/log/nginx/error.log
maxretry = 5
bantime = 3600
```

### 4. Regular Updates

```bash
# Create update script
cat > /opt/authgrid/update.sh << 'EOF'
#!/bin/bash
cd /opt/authgrid
git pull
docker compose -f docker-compose.production.yml build
docker compose -f docker-compose.production.yml up -d
docker image prune -f
EOF

chmod +x /opt/authgrid/update.sh
```

---

## 📊 Monitoring

### Docker Stats

```bash
docker stats
```

### Check Disk Usage

```bash
df -h
docker system df
```

### Monitor Logs

```bash
# Watch API logs
docker compose -f docker-compose.production.yml logs -f api | grep "error"

# Watch nginx access logs
docker compose -f docker-compose.production.yml logs -f nginx
```

### Prometheus & Grafana (Advanced)

See [MONITORING.md](MONITORING.md) for full setup.

---

## 🔄 SSL Certificate Renewal

Certbot auto-renews certificates. To force renewal:

```bash
docker compose -f docker-compose.production.yml run --rm certbot renew
docker compose -f docker-compose.production.yml exec nginx nginx -s reload
```

Check certificate expiry:
```bash
echo | openssl s_client -servername authgrid.org -connect authgrid.org:443 2>/dev/null | \
  openssl x509 -noout -dates
```

---

## 🐛 Troubleshooting

### API Not Responding

```bash
# Check API logs
docker compose -f docker-compose.production.yml logs api

# Check if API container is running
docker compose -f docker-compose.production.yml ps api

# Restart API
docker compose -f docker-compose.production.yml restart api
```

### SSL Certificate Issues

```bash
# Check nginx logs
docker compose -f docker-compose.production.yml logs nginx

# Verify certificate exists
ls -la certbot/conf/live/authgrid.org/

# Regenerate certificate
cd deploy
./init-letsencrypt.sh
```

### Database Connection Errors

```bash
# Check database status
docker compose -f docker-compose.production.yml ps postgres

# Check database logs
docker compose -f docker-compose.production.yml logs postgres

# Restart database
docker compose -f docker-compose.production.yml restart postgres
```

### 502 Bad Gateway

Usually means API is not responding:

```bash
# Check API logs
docker compose -f docker-compose.production.yml logs api

# Check API health
curl http://localhost:8080/health

# Restart API
docker compose -f docker-compose.production.yml restart api
```

---

## 📈 Performance Tuning

### PostgreSQL Configuration

Edit `docker-compose.production.yml`:

```yaml
postgres:
  command:
    - "postgres"
    - "-c"
    - "max_connections=200"
    - "-c"
    - "shared_buffers=256MB"
    - "-c"
    - "effective_cache_size=1GB"
    - "-c"
    - "work_mem=16MB"
```

### Nginx Caching

Add to `nginx/conf.d/authgrid.conf`:

```nginx
# Cache zone
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=api_cache:10m max_size=100m;

# In location block
location /api/ {
    proxy_cache api_cache;
    proxy_cache_valid 200 10m;
    proxy_cache_valid 404 1m;
    ...
}
```

---

## 🔄 Scaling

### Horizontal Scaling

To handle more traffic, add multiple API servers:

```yaml
# docker-compose.production.yml
api:
  deploy:
    replicas: 3
```

Then configure nginx for load balancing:

```nginx
upstream api_backend {
    server api:8080;
    server api:8081;
    server api:8082;
}

location / {
    proxy_pass http://api_backend;
}
```

---

## ✅ Production Checklist

Before going live:

- [ ] DNS records configured and propagated
- [ ] SSL certificates generated and working
- [ ] Database password changed from default
- [ ] Backups configured and tested
- [ ] Firewall rules configured
- [ ] Monitoring set up
- [ ] Tested registration flow
- [ ] Tested login flow
- [ ] Tested API endpoints
- [ ] Checked all logs for errors
- [ ] Set up alerting
- [ ] Documented access credentials
- [ ] Tested recovery procedures

---

## 📞 Support

**Issues:**
- GitHub: https://github.com/Kelsidavis/authgrid/issues
- Email: support@authgrid.org

**Documentation:**
- [QUICKSTART.md](QUICKSTART.md)
- [TECHNICAL.md](TECHNICAL.md)
- [FAQ.md](FAQ.md)

---

## 🎉 You're Live!

Your Authgrid instance is now running at:

- **Website:** https://authgrid.org
- **API:** https://api.authgrid.org
- **Status:** https://api.authgrid.org/health

**Next steps:**
1. Share the demo URL
2. Monitor logs for the first few hours
3. Set up analytics (optional)
4. Announce on social media
5. Gather feedback

---

*Passwords die here. Welcome to production! 🚀*
