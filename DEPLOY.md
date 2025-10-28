# Authgrid Deployment Guide

Complete guide to deploying Authgrid backend online (cheapest options).

---

## 🆓 Option 1: Fly.io (FREE - Recommended)

**Best for:** Production-ready free hosting with PostgreSQL

### Prerequisites
```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.fly/bin:$PATH"
source ~/.bashrc

# Login
fly auth login
```

### Deploy Steps

**1. Create PostgreSQL Database**
```bash
cd /home/k/Desktop/authgrid

# Create Postgres cluster (FREE tier: 3GB storage)
fly postgres create --name authgrid-db --region ord

# Note the connection string shown (or get it later with: fly postgres db show authgrid-db)
```

**2. Create and Deploy API**
```bash
# Create app (don't deploy yet)
fly apps create authgrid-api --region ord

# Attach database (sets DATABASE_URL automatically)
fly postgres attach authgrid-db --app authgrid-api

# Deploy!
fly deploy
```

**3. Check Status**
```bash
# View app
fly status

# Check logs
fly logs

# Test endpoint
curl https://authgrid-api.fly.dev/health
```

**4. Update Frontend Config**
Edit `examples/demo/config.js`:
```javascript
window.AUTHGRID_API_URL = 'https://authgrid-api.fly.dev';
```

**5. Update CORS (Important!)**

Your API needs to allow your frontend domain. Edit `src/api/main.go:57-61`:
```go
allowedOrigins := []string{
    "https://authgrid.org",
    "https://www.authgrid.org",
    "https://authgrid-api.fly.dev", // Fly.io domain
    "http://localhost:3000",
}
```

Then redeploy:
```bash
fly deploy
```

### Run Database Migrations

**Option A: Connect and run manually**
```bash
# Get connection string
fly postgres connect -a authgrid-db

# In Postgres shell:
\i /path/to/migrations/001_init.sql
\q
```

**Option B: Use psql from local machine**
```bash
# Get connection URL
fly postgres db show authgrid-db

# Connect and run migration
psql "postgres://user:pass@host:5432/authgrid?sslmode=require" < migrations/001_init.sql
```

### Useful Commands
```bash
fly status                    # App status
fly logs                      # View logs
fly logs -a authgrid-db      # Database logs
fly ssh console              # SSH into container
fly scale show               # Current scaling
fly scale memory 512         # Increase memory (if needed, costs $)
fly secrets list             # List env vars
fly secrets set KEY=value    # Add env var
fly apps destroy authgrid-api # Delete app
```

### Costs
- **Free tier:** 3 shared VMs (1GB RAM), 3GB persistent storage
- **Your app:** Uses ~256MB RAM, well within free tier
- **Cost: $0/month** ✅

---

## 🆓 Option 2: Render (FREE with limitations)

**Best for:** Quick testing (auto-sleeps after 15 min inactivity)

### Prerequisites
- GitHub account
- Push your code to GitHub

### Deploy Steps

**1. Push to GitHub**
```bash
cd /home/k/Desktop/authgrid
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/YOUR_USERNAME/authgrid.git
git push -u origin main
```

**2. Create PostgreSQL Database**
1. Go to https://dashboard.render.com
2. Click **New +** → **PostgreSQL**
3. Settings:
   - Name: `authgrid-db`
   - Region: Choose closest
   - PostgreSQL Version: 15
   - Instance Type: **Free**
4. Click **Create Database**
5. **Copy Internal Database URL** (starts with `postgres://`)

**3. Create Web Service**
1. Click **New +** → **Web Service**
2. Connect your GitHub repository
3. Settings:
   - Name: `authgrid-api`
   - Region: Same as database
   - Branch: `main`
   - Runtime: **Docker**
   - Dockerfile Path: `src/api/Dockerfile`
   - Instance Type: **Free**
4. **Environment Variables:**
   - Key: `DATABASE_URL`
   - Value: *paste Internal Database URL from step 2*
5. Click **Create Web Service**

**4. Run Migrations**

Once deployed, use Render shell:
```bash
# In Render dashboard, click your web service
# Go to Shell tab
psql $DATABASE_URL < /migrations/001_init.sql
```

Or connect from local:
```bash
# Use the External Database URL from Render
psql "YOUR_EXTERNAL_DB_URL" < migrations/001_init.sql
```

**5. Update Frontend**
Your API URL will be: `https://authgrid-api.onrender.com`

```javascript
// config.js
window.AUTHGRID_API_URL = 'https://authgrid-api.onrender.com';
```

### Limitations
- ⚠️ **Spins down after 15 min inactivity** (30s cold start)
- ⚠️ **Free PostgreSQL expires after 90 days**
- ✅ Great for testing/demos

---

## 🆓 Option 3: Oracle Cloud (FREE Forever)

**Best for:** Generous free tier forever (complex setup)

### Prerequisites
1. Sign up at https://www.oracle.com/cloud/free/
2. Verify account (credit card required, but not charged)

### Deploy Steps

**1. Create Compute Instance**
1. Go to **Compute** → **Instances**
2. Click **Create Instance**
3. Settings:
   - Name: `authgrid-server`
   - Image: Ubuntu 22.04
   - Shape: **VM.Standard.E2.1.Micro** (Always Free)
   - Add SSH key (generate or use existing)
4. Create and note the **Public IP**

**2. Configure Firewall**
```bash
# In Oracle Cloud Console
# Go to Instance → Subnet → Security List → Add Ingress Rule
# Port: 8080, Source: 0.0.0.0/0
```

**3. SSH Into Server**
```bash
ssh ubuntu@YOUR_PUBLIC_IP
```

**4. Install Dependencies**
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ubuntu
exit  # Logout and login again

# SSH back in
ssh ubuntu@YOUR_PUBLIC_IP
```

**5. Deploy Authgrid**
```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/authgrid.git
cd authgrid

# Start services
docker-compose up -d

# Check status
docker-compose ps
docker-compose logs -f
```

**6. Configure Domain (Optional)**

If you have a domain, point it to the public IP:
```
A Record: api.authgrid.org → YOUR_PUBLIC_IP
```

**7. Setup Nginx + SSL**
```bash
# Install nginx
sudo apt install -y nginx certbot python3-certbot-nginx

# Create nginx config
sudo nano /etc/nginx/sites-available/authgrid
```

Add:
```nginx
server {
    listen 80;
    server_name api.authgrid.org;  # or use IP

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/authgrid /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Get SSL certificate (if using domain)
sudo certbot --nginx -d api.authgrid.org
```

**8. Update Frontend**
```javascript
// config.js
window.AUTHGRID_API_URL = 'https://api.authgrid.org';
// or
window.AUTHGRID_API_URL = 'http://YOUR_PUBLIC_IP:8080';
```

### Costs
- **Always Free:** 2 VMs (1GB RAM), 200GB storage, 10TB bandwidth
- **Cost: $0/month forever** ✅

---

## 💰 Option 4: Hetzner Cloud ($3.79/month - Best Value)

**Best for:** Cheap VPS with full control

### Deploy Steps

**1. Create Account**
1. Sign up at https://console.hetzner.cloud
2. Add payment method

**2. Create Server**
1. **New Project** → **Add Server**
2. Settings:
   - Location: Choose closest
   - Image: Ubuntu 22.04
   - Type: **CX11** (2GB RAM, $3.79/mo)
   - Add SSH key
3. Create and note **IPv4 address**

**3. Follow Oracle Cloud steps 3-8 above**

Same exact process: SSH in, install Docker, deploy, configure nginx.

### Costs
- **CX11:** $3.79/month (2GB RAM, 20GB SSD, 20TB traffic)
- **CX22:** $6.86/month (4GB RAM, 40GB SSD) - if you need more

---

## 💰 Option 5: DigitalOcean ($4-6/month)

**Best for:** Popular choice with good documentation

### Deploy Steps

**1. Create Account**
1. Sign up at https://www.digitalocean.com
2. Get $200 credit (60 days) with referral link

**2. Create Droplet**
1. **Create** → **Droplets**
2. Settings:
   - Image: Ubuntu 22.04 LTS
   - Plan: **Basic** ($4/mo or $6/mo)
   - Region: Choose closest
   - Add SSH key
3. Create and note **IP address**

**3. Follow Oracle Cloud steps 3-8 above**

### Costs
- **Basic $4/mo:** 512MB RAM, 10GB SSD, 500GB transfer
- **Basic $6/mo:** 1GB RAM, 25GB SSD, 1TB transfer

---

## 📊 Quick Comparison

| Provider | Monthly Cost | RAM | Storage | Setup Difficulty | Free Tier |
|----------|-------------|-----|---------|-----------------|-----------|
| **Fly.io** | $0 | 256MB | 3GB | ⭐ Easy | Yes |
| **Render** | $0* | 512MB | 1GB | ⭐ Easy | Yes (limited) |
| **Oracle Cloud** | $0 | 1GB | 50GB | ⭐⭐⭐ Hard | Forever |
| **Hetzner** | $3.79 | 2GB | 20GB | ⭐⭐ Medium | No |
| **DigitalOcean** | $4-6 | 512MB-1GB | 10-25GB | ⭐⭐ Medium | $200 credit |

*Render free tier: Sleeps after 15 min, DB expires in 90 days

---

## 🎯 My Recommendation

**For most users: Fly.io (FREE)**
- Easiest deployment
- True free tier
- PostgreSQL included
- No sleep/downtime
- Perfect for Authgrid

**Quick start:**
```bash
# Install CLI
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Deploy
cd /home/k/Desktop/authgrid
fly postgres create --name authgrid-db
fly apps create authgrid-api
fly postgres attach authgrid-db --app authgrid-api
fly deploy

# Done! Your API is at: https://authgrid-api.fly.dev
```

---

## 🔧 After Deployment

### 1. Update CORS
Edit `src/api/main.go` to allow your frontend domain:
```go
allowedOrigins := []string{
    "https://authgrid.org",
    "https://YOUR-API-DOMAIN.fly.dev",
    "http://localhost:3000",
}
```

### 2. Update Frontend Config
Edit `examples/demo/config.js`:
```javascript
window.AUTHGRID_API_URL = 'https://YOUR-API-DOMAIN.fly.dev';
```

### 3. Test
```bash
# Health check
curl https://YOUR-API-DOMAIN.fly.dev/health

# Try registering from frontend
# Check logs to see requests coming in
```

### 4. Monitor
```bash
# Fly.io
fly logs
fly status

# VPS (Hetzner, DO, Oracle)
ssh user@server
docker-compose logs -f
```

---

## 🆘 Troubleshooting

### Database Connection Failed
```bash
# Check DATABASE_URL is set
fly secrets list

# Verify database is running
fly postgres db list -a authgrid-db

# Test connection
fly ssh console -a authgrid-api
env | grep DATABASE_URL
```

### CORS Errors
Make sure your frontend domain is in the allowed origins list, then redeploy.

### 502 Bad Gateway
App might be crashing. Check logs:
```bash
fly logs
```

Common issues:
- DATABASE_URL not set
- Database migrations not run
- Port mismatch (should be 8080)

---

## 📚 Further Reading

- [Fly.io Docs](https://fly.io/docs/)
- [Render Docs](https://render.com/docs)
- [Oracle Cloud Free Tier](https://www.oracle.com/cloud/free/)
- [Hetzner Cloud Docs](https://docs.hetzner.com/cloud/)
- [DigitalOcean Tutorials](https://www.digitalocean.com/community/tutorials)

---

**Ready to deploy?** Start with Fly.io:
```bash
curl -L https://fly.io/install.sh | sh
fly auth login
cd /home/k/Desktop/authgrid
fly launch
```
