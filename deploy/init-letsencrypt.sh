#!/bin/bash

# Initialize Let's Encrypt SSL certificates for authgrid.org

set -e

# Load environment variables
if [ ! -f ../.env.production ]; then
    echo "Error: .env.production file not found"
    echo "Please copy .env.production and configure it first"
    exit 1
fi

source ../.env.production

DOMAINS=($DOMAIN "www.$DOMAIN" "api.$DOMAIN")
RSA_KEY_SIZE=4096
DATA_PATH="../certbot"
EMAIL=${LETSENCRYPT_EMAIL}
STAGING=0 # Set to 1 for testing

echo "### Preparing directories ..."
mkdir -p "$DATA_PATH/conf"
mkdir -p "$DATA_PATH/www"

echo ""
echo "### Downloading recommended TLS parameters ..."
mkdir -p "$DATA_PATH/conf"
curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf > "$DATA_PATH/conf/options-ssl-nginx.conf"
curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem > "$DATA_PATH/conf/ssl-dhparams.pem"

echo ""
echo "### Creating dummy certificate for $DOMAIN ..."
PATH_TO_CERT="$DATA_PATH/conf/live/$DOMAIN"
mkdir -p "$PATH_TO_CERT"
docker compose -f ../docker-compose.production.yml run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$RSA_KEY_SIZE -days 1\
    -keyout '$PATH_TO_CERT/privkey.pem' \
    -out '$PATH_TO_CERT/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo ""

echo "### Starting nginx ..."
docker compose -f ../docker-compose.production.yml up --force-recreate -d nginx
echo ""

echo "### Deleting dummy certificate for $DOMAIN ..."
docker compose -f ../docker-compose.production.yml run --rm --entrypoint "\
  rm -Rf /etc/letsencrypt/live/$DOMAIN && \
  rm -Rf /etc/letsencrypt/archive/$DOMAIN && \
  rm -Rf /etc/letsencrypt/renewal/$DOMAIN.conf" certbot
echo ""

echo "### Requesting Let's Encrypt certificate for $DOMAIN ..."
# Join $DOMAINS to -d args
DOMAIN_ARGS=""
for domain in "${DOMAINS[@]}"; do
  DOMAIN_ARGS="$DOMAIN_ARGS -d $domain"
done

# Select appropriate email arg
case "$EMAIL" in
  "") EMAIL_ARG="--register-unsafely-without-email" ;;
  *) EMAIL_ARG="--email $EMAIL" ;;
esac

# Enable staging mode if needed
if [ $STAGING != "0" ]; then STAGING_ARG="--staging"; fi

docker compose -f ../docker-compose.production.yml run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    $STAGING_ARG \
    $EMAIL_ARG \
    $DOMAIN_ARGS \
    --rsa-key-size $RSA_KEY_SIZE \
    --agree-tos \
    --force-renewal" certbot
echo ""

echo "### Reloading nginx ..."
docker compose -f ../docker-compose.production.yml exec nginx nginx -s reload

echo ""
echo "### SSL certificates successfully generated!"
echo "### Your site should now be available at https://$DOMAIN"
