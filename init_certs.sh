#!/bin/bash

domains=(apis.greyboxpay.com wallet.greyboxpay.com)
email="fortuneosho@gmail.com" 

# Stopping Nginx before generating certificates
docker-compose stop nginx

for domain in "${domains[@]}"; do
  docker-compose run --rm --entrypoint "\
    certbot certonly --webroot -w /var/www/certbot \
    --email $email \
    --agree-tos \
    --no-eff-email \
    -d $domain" certbot
done

# Restart Nginx
docker-compose up -d nginx
