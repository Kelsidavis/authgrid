# 🚀 Deploy Authgrid to Production

Your production deployment is **ready to go**! All configuration files, scripts, and documentation have been created.

## ✅ What's Ready

### Production Files Created
- ✅ `docker-compose.production.yml` - Production orchestration
- ✅ `nginx/nginx.conf` - Nginx main config with security headers
- ✅ `nginx/conf.d/authgrid.conf` - Site-specific routing for authgrid.org
- ✅ `.env.production` - Environment configuration template
- ✅ `deploy/deploy.sh` - Automated deployment script
- ✅ `deploy/init-letsencrypt.sh` - SSL certificate setup
- ✅ `examples/demo/index.production.html` - Production demo (uses api.authgrid.org)
- ✅ `DEPLOYMENT.md` - Comprehensive 500+ line deployment guide

### Services Configured
- ✅ PostgreSQL database with secure password from env
- ✅ API server at api.authgrid.org
- ✅ Main website at authgrid.org and www.authgrid.org
- ✅ Nginx reverse proxy with SSL termination
- ✅ Certbot for automated Let's Encrypt SSL certificates
- ✅ CORS properly configured for API access
- ✅ HTTP to HTTPS redirect

---

## 📋 Pre-Deployment Checklist

Before deploying, ensure you have:

1. **Server Access**
   - [ ] Ubuntu 20.04+ server with root/sudo access
   - [ ] Server IP address noted down

2. **DNS Configuration** (CRITICAL - Do this first!)
   - [ ] Point `authgrid.org` A record to your server IP
   - [ ] Point `www.authgrid.org` A record to your server IP
   - [ ] Point `api.authgrid.org` A record to your server IP
   - [ ] Wait 5-30 minutes for DNS propagation
   - [ ] Verify with: `dig authgrid.org`

3. **Server Requirements**
   - [ ] Docker installed
   - [ ] Docker Compose installed
   - [ ] Ports 80 and 443 open in firewall
   - [ ] At least 1GB RAM available

---

## 🎯 Deploy in 3 Steps

### Step 1: Transfer Files to Server

From your local machine:

```bash
# Option A: Using rsync (recommended)
rsync -avz --exclude 'node_modules' \
  /home/k/Desktop/authgrid/ \
  root@YOUR_SERVER_IP:/opt/authgrid/

# Option B: Using git (if repo is pushed)
ssh root@YOUR_SERVER_IP
git clone https://github.com/Kelsidavis/authgrid.git /opt/authgrid
cd /opt/authgrid
```

### Step 2: Configure Environment

SSH into your server:

```bash
ssh root@YOUR_SERVER_IP
cd /opt/authgrid
```

Create production environment file:

```bash
cp .env.production .env.production
nano .env.production
```

**CRITICAL: Change the database password!**

```bash
# Generate a secure password
openssl rand -base64 32

# Update .env.production with the generated password:
POSTGRES_PASSWORD=<YOUR_GENERATED_PASSWORD>
```

Your `.env.production` should look like:
```
POSTGRES_USER=authgrid
POSTGRES_PASSWORD=xK9mP2nQ8vL5wE7yT3zR6jH4bN1sA0fD==  # Your generated password
DOMAIN=authgrid.org
LETSENCRYPT_EMAIL=admin@authgrid.org
AUTHGRID_DOMAIN=authgrid.org
```

### Step 3: Run Deployment

Make deployment script executable and run:

```bash
cd /opt/authgrid/deploy
chmod +x deploy.sh init-letsencrypt.sh
./deploy.sh
```

The script will:
1. Build Docker images (~2 minutes)
2. Start PostgreSQL database
3. Start API server
4. Prompt you to generate SSL certificates (press 'y')
5. Start Nginx with SSL
6. Start Certbot for auto-renewal

**When prompted for SSL generation, press `y`**

---

## ✅ Verify Deployment

After deployment completes:

### 1. Check Services Status

```bash
cd /opt/authgrid
docker compose -f docker-compose.production.yml ps
```

All services should show "Up" or "Up (healthy)":
```
authgrid-postgres    Up (healthy)
authgrid-api         Up
authgrid-nginx       Up
authgrid-certbot     Up
```

### 2. Test API Health

```bash
curl https://api.authgrid.org/health
```

Should return:
```json
{"status":"healthy","time":"2025-10-28T..."}
```

### 3. Test in Browser

Open these URLs:
- https://authgrid.org (main demo)
- https://www.authgrid.org (should redirect to main)
- https://api.authgrid.org/health (API health check)

### 4. Test Registration & Login

1. Go to https://authgrid.org
2. Click "Register Now"
3. You should get a handle like `a1b2c3d4e5@authgrid.org`
4. Click "Login" with your handle
5. Should successfully authenticate

---

## 📊 View Logs

If anything doesn't work:

```bash
# View all logs
docker compose -f docker-compose.production.yml logs -f

# View specific service logs
docker compose -f docker-compose.production.yml logs -f api
docker compose -f docker-compose.production.yml logs -f nginx
docker compose -f docker-compose.production.yml logs -f postgres
```

---

## 🔧 Common Issues

### Issue: SSL Certificate Generation Fails

**Symptoms**: Certbot errors during deployment

**Solution**: Ensure DNS is fully propagated first
```bash
# Check DNS propagation
dig authgrid.org
dig api.authgrid.org

# If not propagated, wait and try again
cd /opt/authgrid/deploy
./init-letsencrypt.sh
```

### Issue: 502 Bad Gateway

**Symptoms**: Nginx shows 502 error

**Solution**: API server may not be ready yet
```bash
# Check API logs
docker compose -f docker-compose.production.yml logs api

# Restart API if needed
docker compose -f docker-compose.production.yml restart api
```

### Issue: Database Connection Errors

**Symptoms**: API logs show "connection refused"

**Solution**: Database may need more time to initialize
```bash
# Check database status
docker compose -f docker-compose.production.yml ps postgres

# Wait for health check to pass (10-20 seconds)
# Then restart API
docker compose -f docker-compose.production.yml restart api
```

---

## 🎉 Post-Deployment

Once everything is working:

1. **Set Up Automated Backups**
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
   find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
   EOF

   chmod +x /opt/authgrid/backup.sh

   # Add to crontab (daily at 2 AM)
   (crontab -l 2>/dev/null; echo "0 2 * * * /opt/authgrid/backup.sh") | crontab -
   ```

2. **Monitor Service Health**
   ```bash
   # Check service status
   docker compose -f docker-compose.production.yml ps

   # Check disk usage
   df -h
   docker system df
   ```

3. **Share the Demo**
   - Demo URL: https://authgrid.org
   - API docs: https://api.authgrid.org/health

---

## 📚 Full Documentation

For comprehensive information:
- **DEPLOYMENT.md** - Complete deployment guide (500+ lines)
- **README.md** - Project overview
- **TECHNICAL.md** - Architecture details
- **FAQ.md** - Common questions

---

## 🆘 Need Help?

If you encounter any issues:

1. Check logs: `docker compose -f docker-compose.production.yml logs -f`
2. Review DEPLOYMENT.md for detailed troubleshooting
3. Verify DNS propagation: `dig authgrid.org`
4. Ensure firewall allows ports 80 and 443
5. Verify .env.production password was changed

---

## 🚀 You're Almost Live!

Your Authgrid deployment will be available at:

- **Main Site**: https://authgrid.org
- **API**: https://api.authgrid.org
- **Health Check**: https://api.authgrid.org/health

**Next:** SSH to your server and run the 3 deployment steps above!

*Passwords die here. Welcome to production! 🔐*
