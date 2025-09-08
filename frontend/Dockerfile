FROM nginxinc/nginx-unprivileged:1.27-alpine

# Statik dosyalar
COPY index.html /usr/share/nginx/html/index.html
COPY app.js     /usr/share/nginx/html/app.js
COPY styles.css /usr/share/nginx/html/styles.css

# Basit server config
COPY default.conf /etc/nginx/conf.d/default.conf

# Bu imaj user=101 ve port=8080 ile gelir
EXPOSE 8080
