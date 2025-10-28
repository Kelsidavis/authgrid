# Production Commands Quick Reference

Quick reference for managing your Authgrid production deployment.

## 🚀 Deployment Commands

```bash
# Initial deployment
cd /opt/authgrid/deploy
./deploy.sh

# Manual SSL setup (if needed)
cd /opt/authgrid/deploy
./init-letsencrypt.sh
```

## 📊 Service Management

```bash
# Check all services status
docker compose -f docker-compose.production.yml ps

# Start all services
docker compose -f docker-compose.production.yml up -d

# Stop all services
docker compose -f docker-compose.production.yml down

# Restart all services
docker compose -f docker-compose.production.yml restart

# Restart specific service
docker compose -f docker-compose.production.yml restart api
docker compose -f docker-compose.production.yml restart nginx
docker compose -f docker-compose.production.yml restart postgres
```

## 📝 Logs

```bash
# View all logs (live)
docker compose -f docker-compose.production.yml logs -f

# View specific service logs
docker compose -f docker-compose.production.yml logs -f api
docker compose -f docker-compose.production.yml logs -f nginx
docker compose -f docker-compose.production.yml logs -f postgres
docker compose -f docker-compose.production.yml logs -f certbot

# View last 100 lines
docker compose -f docker-compose.production.yml logs --tail=100 api

# Search logs for errors
docker compose -f docker-compose.production.yml logs api | grep -i error
```

## 🔄 Updates

```bash
# Pull latest code
cd /opt/authgrid
git pull origin main

# Rebuild and restart
docker compose -f docker-compose.production.yml build
docker compose -f docker-compose.production.yml up -d

# Clean up old images
docker image prune -f
```

## 🗄️ Database

```bash
# Create backup
docker compose -f docker-compose.production.yml exec postgres \
  pg_dump -U authgrid authgrid > backup_$(date +%Y%m%d_%H%M%S).sql

# Compressed backup
docker compose -f docker-compose.production.yml exec postgres \
  pg_dump -U authgrid authgrid | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Restore from backup
cat backup.sql | docker compose -f docker-compose.production.yml exec -T postgres \
  psql -U authgrid authgrid

# Access database shell
docker compose -f docker-compose.production.yml exec postgres \
  psql -U authgrid -d authgrid

# Check database size
docker compose -f docker-compose.production.yml exec postgres \
  psql -U authgrid -d authgrid -c "SELECT pg_size_pretty(pg_database_size('authgrid'));"
```

## 🔒 SSL Certificates

```bash
# Force certificate renewal
docker compose -f docker-compose.production.yml run --rm certbot renew

# Reload nginx after renewal
docker compose -f docker-compose.production.yml exec nginx nginx -s reload

# Check certificate expiry
echo | openssl s_client -servername authgrid.org -connect authgrid.org:443 2>/dev/null | \
  openssl x509 -noout -dates

# Test certificate
openssl s_client -connect authgrid.org:443 -servername authgrid.org
```

## 🧪 Testing

```bash
# Test API health
curl https://api.authgrid.org/health

# Test API locally (from server)
curl http://localhost:8080/health

# Test nginx configuration
docker compose -f docker-compose.production.yml exec nginx nginx -t

# Test database connection
docker compose -f docker-compose.production.yml exec postgres \
  pg_isready -U authgrid
```

## 📊 Monitoring

```bash
# Docker stats (CPU, memory usage)
docker stats

# Disk usage
df -h
docker system df

# Network ports
sudo netstat -tlnp | grep -E ':(80|443|5432|8080)'

# Check container health
docker compose -f docker-compose.production.yml ps

# View container resource usage
docker compose -f docker-compose.production.yml top
```

## 🔍 Troubleshooting

```bash
# Check if services are running
docker compose -f docker-compose.production.yml ps

# Inspect container
docker compose -f docker-compose.production.yml logs api --tail=50

# Check API connectivity from nginx
docker compose -f docker-compose.production.yml exec nginx ping api

# Check database connectivity from API
docker compose -f docker-compose.production.yml exec api sh -c \
  'apk add postgresql-client && psql $DATABASE_URL -c "SELECT 1;"'

# Restart problematic service
docker compose -f docker-compose.production.yml restart api

# Full restart (if needed)
docker compose -f docker-compose.production.yml down
docker compose -f docker-compose.production.yml up -d
```

## 🌐 DNS Verification

```bash
# Check DNS propagation
dig authgrid.org
dig www.authgrid.org
dig api.authgrid.org

# Check from different DNS servers
dig @8.8.8.8 authgrid.org
dig @1.1.1.1 authgrid.org

# Verify A records
nslookup authgrid.org
```

## 🔐 Security

```bash
# Check firewall status
sudo ufw status

# Allow required ports
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Check nginx security headers
curl -I https://authgrid.org

# View active connections
docker compose -f docker-compose.production.yml exec nginx \
  sh -c 'cat /var/log/nginx/access.log | tail -20'
```

## 📦 Container Management

```bash
# Remove stopped containers
docker container prune -f

# Remove unused images
docker image prune -f

# Remove unused volumes
docker volume prune -f

# Full cleanup (BE CAREFUL!)
docker system prune -af --volumes
```

## 🔄 Backup & Restore

```bash
# Create full backup (database + config)
tar -czf authgrid_backup_$(date +%Y%m%d).tar.gz \
  /opt/authgrid/.env.production \
  /opt/authgrid/certbot/conf
docker compose -f docker-compose.production.yml exec postgres \
  pg_dump -U authgrid authgrid > /tmp/db_backup.sql

# Restore configuration
tar -xzf authgrid_backup_20241028.tar.gz -C /

# Restore database
cat /tmp/db_backup.sql | docker compose -f docker-compose.production.yml exec -T postgres \
  psql -U authgrid authgrid
```

## 📈 Performance

```bash
# Check API response time
time curl https://api.authgrid.org/health

# Monitor logs in real-time
docker compose -f docker-compose.production.yml logs -f api | grep -E '(POST|GET)'

# Check PostgreSQL performance
docker compose -f docker-compose.production.yml exec postgres \
  psql -U authgrid -d authgrid -c \
  "SELECT * FROM pg_stat_activity WHERE state = 'active';"
```

## 🆘 Emergency Commands

```bash
# Emergency stop (immediate)
docker compose -f docker-compose.production.yml kill

# Force restart everything
docker compose -f docker-compose.production.yml down --remove-orphans
docker compose -f docker-compose.production.yml up -d --force-recreate

# Check what's consuming resources
docker stats --no-stream
top -bn1 | head -20

# View system logs
journalctl -xe
dmesg | tail -50
```

---

## 📍 Useful Paths

```
/opt/authgrid/                    # Main application directory
/opt/authgrid/.env.production     # Environment configuration
/opt/authgrid/certbot/conf        # SSL certificates
/opt/authgrid/deploy/             # Deployment scripts
/var/lib/docker/volumes/          # Docker volumes (database data)
~/.authgrid/                      # CLI keystore
```

## 🌐 URLs

```
https://authgrid.org              # Main website
https://www.authgrid.org          # WWW redirect
https://api.authgrid.org          # API endpoint
https://api.authgrid.org/health   # Health check
```

---

**Tip**: Bookmark this file for quick reference when managing your production deployment!
