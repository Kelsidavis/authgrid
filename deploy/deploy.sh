#!/bin/bash

# Authgrid Production Deployment Script

set -e

echo "╔══════════════════════════════════════════╗"
echo "║   Authgrid Production Deployment         ║"
echo "╚══════════════════════════════════════════╝"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "⚠️  Warning: This script may need sudo for Docker commands"
    echo ""
fi

# Check if .env.production exists
if [ ! -f ../.env.production ]; then
    echo "❌ Error: .env.production not found"
    echo ""
    echo "Please create .env.production from .env.production template:"
    echo "  cp .env.production.example .env.production"
    echo "  nano .env.production"
    echo ""
    exit 1
fi

# Load environment
source ../.env.production

# Check if password was changed
if [[ "$POSTGRES_PASSWORD" == "CHANGE_THIS_TO_SECURE_PASSWORD" ]]; then
    echo "❌ Error: Please change the database password in .env.production"
    echo ""
    echo "Generate a secure password:"
    echo "  openssl rand -base64 32"
    echo ""
    exit 1
fi

echo "✅ Environment configured"
echo ""

# Step 1: Build images
echo "📦 Building Docker images..."
cd ..
docker compose -f docker-compose.production.yml build
echo "✅ Images built"
echo ""

# Step 2: Start database
echo "🗄️  Starting database..."
docker compose -f docker-compose.production.yml up -d postgres
echo "⏳ Waiting for database to be ready..."
sleep 10
echo "✅ Database ready"
echo ""

# Step 3: Start API
echo "🚀 Starting API server..."
docker compose -f docker-compose.production.yml up -d api
echo "✅ API server started"
echo ""

# Step 4: Check if SSL certificates exist
if [ ! -d "certbot/conf/live/$DOMAIN" ]; then
    echo "🔒 SSL certificates not found"
    echo "❓ Do you want to generate Let's Encrypt certificates now? (y/n)"
    read -r GENERATE_SSL

    if [[ "$GENERATE_SSL" == "y" || "$GENERATE_SSL" == "Y" ]]; then
        echo "🔐 Generating SSL certificates..."
        cd deploy
        chmod +x init-letsencrypt.sh
        ./init-letsencrypt.sh
        cd ..
    else
        echo "⚠️  Skipping SSL generation"
        echo "   You can generate certificates later with:"
        echo "   cd deploy && ./init-letsencrypt.sh"
    fi
else
    echo "✅ SSL certificates found"
fi

echo ""

# Step 5: Start nginx
echo "🌐 Starting nginx..."
docker compose -f docker-compose.production.yml up -d nginx
echo "✅ Nginx started"
echo ""

# Step 6: Start certbot (for auto-renewal)
echo "🔄 Starting certbot for auto-renewal..."
docker compose -f docker-compose.production.yml up -d certbot
echo "✅ Certbot started"
echo ""

# Check status
echo "📊 Checking service status..."
docker compose -f docker-compose.production.yml ps
echo ""

# Test API
echo "🧪 Testing API..."
sleep 3
API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/health || echo "000")

if [ "$API_STATUS" == "200" ]; then
    echo "✅ API is responding"
else
    echo "⚠️  API returned status: $API_STATUS"
    echo "   Check logs: docker compose -f docker-compose.production.yml logs api"
fi

echo ""
echo "╔══════════════════════════════════════════╗"
echo "║   Deployment Complete!                   ║"
echo "╠══════════════════════════════════════════╣"
echo "║                                          ║"
echo "║   🌐 Website: https://$DOMAIN"
echo "║   🔌 API:     https://api.$DOMAIN"
echo "║                                          ║"
echo "║   View logs:                             ║"
echo "║   docker compose -f docker-compose.production.yml logs -f"
echo "║                                          ║"
echo "║   Stop services:                         ║"
echo "║   docker compose -f docker-compose.production.yml down"
echo "║                                          ║"
echo "╚══════════════════════════════════════════╝"
echo ""
