docker network create shoplite-net
docker run -d --name frontend --restart unless-stopped --network shoplite-net \
  -p 127.0.0.1:3000:80 YOURUSER/shoplite-frontend:v1.0.0

apt-get update && apt-get install -y nginx
cat > /etc/nginx/sites-available/shoplite <<'EOF'
server {
    listen 80 default_server;
    location /api/ {
        proxy_pass http://10.0.2.10:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
    }
}
EOF
ln -sf /etc/nginx/sites-available/shoplite /etc/nginx/sites-enabled/shoplite
rm -f /etc/nginx/sites-enabled/default
nginx -t && systemctl reload nginx
