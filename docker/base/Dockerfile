FROM nginx:stable

ADD docker/base/nginx.conf /etc/nginx/nginx.conf
ADD docker/base/loadcat.conf /etc/loadcat.conf
ADD loadcatd /usr/bin/loadcat
ADD ui/templates /var/lib/loadcat/ui/templates

WORKDIR /var/lib/loadcat

ENTRYPOINT ["loadcat", "-config", "/etc/loadcat.conf"]
