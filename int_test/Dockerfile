FROM nginx
LABEL description="Nginx server for serving images"
LABEL service="image-previewer-test-server"
ENV CODE_DIR /go/src/
WORKDIR ${CODE_DIR}
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./temp /usr/share/nginx/html/images

EXPOSE 80
