FROM nginx
LABEL description="Nginx server for serving images"
LABEL service="image-previewer-test-server"

COPY ./int_test/nginx.conf /etc/nginx/nginx.conf
COPY ./int_test/images /usr/share/nginx/html/images

EXPOSE 80
