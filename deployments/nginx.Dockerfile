FROM nginx:alpine

COPY deployments/nginx.conf /etc/nginx/nginx.conf
