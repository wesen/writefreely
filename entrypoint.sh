#!/bin/sh

echo "Creating config.ini"
echo "$CONFIG_INI" > /go/config.ini
ls -la /go/config.ini

echo "Creating keys"
mkdir -p /go/keys/keys
echo "$COOKIES_AUTH_KEY" | base64 -d > /go/keys/keys/cookies_auth.aes256
echo "$COOKIES_ENC_KEY" | base64 -d > /go/keys/keys/cookies_enc.aes256
echo "$CSRF_KEY" | base64 -d > /go/keys/keys/csrf.aes256
echo "$EMAIL_KEY" | base64 -d > /go/keys/keys/email.aes256
echo "$OAUTH_KEY" | base64 -d > /go/keys/keys/oauth.aes256
ls -la /go/keys/keys

echo "Starting app"
cd /go
./cmd/writefreely/writefreely db init
./cmd/writefreely/writefreely db migrate
./cmd/writefreely/writefreely