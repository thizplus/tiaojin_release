# Build stage
FROM node:20-alpine as build

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install

COPY . .

# รับ environment variables จาก build args
ARG VITE_APP_TITLE
ARG VITE_API_URL
ARG VITE_APP_VERSION
ARG VITE_WS_URL
ARG VITE_WS_AUTO_CONNECT
ARG VITE_WS_NOTIFICATIONS

# Export เป็น environment variables สำหรับ Vite
ENV VITE_APP_TITLE=$VITE_APP_TITLE
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_APP_VERSION=$VITE_APP_VERSION
ENV VITE_WS_URL=$VITE_WS_URL
ENV VITE_WS_AUTO_CONNECT=$VITE_WS_AUTO_CONNECT
ENV VITE_WS_NOTIFICATIONS=$VITE_WS_NOTIFICATIONS

# Build project
RUN npm run build

# Production stage  
FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html

RUN echo 'server { \
    listen 80; \
    server_name localhost; \
    \
    location / { \
        root /usr/share/nginx/html; \
        index index.html; \
        try_files $uri $uri/ /index.html; \
    } \
    \
    gzip on; \
    gzip_vary on; \
    gzip_min_length 10240; \
    gzip_proxied expired no-cache no-store private auth; \
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml application/javascript; \
    gzip_disable "MSIE [1-6]\."; \
}' > /etc/nginx/conf.d/default.conf

RUN echo "OK" > /usr/share/nginx/html/health

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]