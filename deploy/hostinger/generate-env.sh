#!/usr/bin/env sh
set -eu

if [ -f .env ] && [ "${1:-}" != "--force" ]; then
  printf '%s\n' ".env already exists. Run: sh deploy/hostinger/generate-env.sh --force"
  exit 1
fi

require_value() {
  name="$1"
  prompt="$2"
  value=""
  while [ -z "$value" ]; do
    printf '%s: ' "$prompt"
    IFS= read -r value
  done
  eval "$name=\$value"
}

read_with_default() {
  name="$1"
  prompt="$2"
  default_value="$3"
  printf '%s [%s]: ' "$prompt" "$default_value"
  IFS= read -r value
  if [ -z "$value" ]; then
    value="$default_value"
  fi
  eval "$name=\$value"
}

random_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 32
    return
  fi

  date +%s%N | sha256sum | cut -d ' ' -f 1
}

umask 077

require_value DOMAIN "API domain, example api.example.com"
require_value FRONTEND_URL "Frontend URL, example https://example.com"
read_with_default APP_HOST_PORT "Host port for API on 127.0.0.1" "8083"
read_with_default ALLOWED_ORIGINS "Allowed frontend origin(s), comma-separated" "$FRONTEND_URL"
read_with_default SMTP_HOST "SMTP host" "smtp.gmail.com"
read_with_default SMTP_PORT "SMTP port" "587"
require_value SMTP_USER "SMTP user/email"
require_value SMTP_PASS "SMTP password"

DB_PASSWORD="$(random_secret)"
ACCESS_TOKEN_SECRET="$(random_secret)"
REFRESH_TOKEN_SECRET="$(random_secret)"
VERIFICATION_TOKEN_SECRET="$(random_secret)"

cat > .env <<EOF
APP_ENV=production
APP_PORT=8083
APP_HOST_PORT=$APP_HOST_PORT

DOMAIN=$DOMAIN
FRONTEND_URL=$FRONTEND_URL
ALLOWED_ORIGINS=$ALLOWED_ORIGINS

DOCKER_IMAGE=zidanindratama/backend-brevet:latest

DB_HOST=db
DB_USER=postgres
DB_PASSWORD=$DB_PASSWORD
DB_NAME=brevetdb
DB_PORT=5432
DB_SSLMODE=disable

REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

ACCESS_TOKEN_SECRET=$ACCESS_TOKEN_SECRET
REFRESH_TOKEN_SECRET=$REFRESH_TOKEN_SECRET
VERIFICATION_TOKEN_SECRET=$VERIFICATION_TOKEN_SECRET
ACCESS_TOKEN_EXPIRY_HOURS=24
REFRESH_TOKEN_EXPIRY_HOURS=24
VERIFICATION_TOKEN_EXPIRY_MINUTES=15
TOKEN_BLACKLIST_TTL=86400
CLEANUP_INTERVAL_HOURS=1

UPLOAD_DIR=/root/public/uploads

SMTP_HOST=$SMTP_HOST
SMTP_PORT=$SMTP_PORT
SMTP_USER=$SMTP_USER
SMTP_PASS=$SMTP_PASS
EOF

printf '%s\n' ".env production created. Keep this file only on the VPS and do not commit it."
