FROM nginx:1.21.0

ADD docker/base/nginx.conf /etc/nginx/nginx.conf
ADD docker/base/loadcat.conf /etc/loadcat.conf
ADD loadcat /usr/bin/loadcat
ADD ui/templates /var/lib/loadcat/ui/templates

WORKDIR /var/lib/loadcat

ENTRYPOINT ["loadcat", "-config", "/etc/loadcat.conf"]
