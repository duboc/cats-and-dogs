FROM nginx:1.15.0
ADD hello.conf /etc/nginx/conf.d/
RUN rm /var/log/nginx/*
COPY index.html /usr/share/nginx/html/
COPY css /usr/share/nginx/html/css/
COPY js /usr/share/nginx/html/js/
COPY cat.jpg /usr/share/nginx/html/
COPY dog.jpg /usr/share/nginx/html/
COPY hedhog.jpg /usr/share/nginx/html/
COPY start.sh /start.sh

#CMD ["nginx", "-g", "daemon off;"]
CMD ["/start.sh"]

