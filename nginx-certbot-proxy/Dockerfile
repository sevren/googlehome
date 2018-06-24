FROM yobasystems/alpine-nginx:armhf

ARG DNS_HOST=REPLACE_ME
ARG YOUR_EMAIL=REPLACE_ME

RUN apk update && apk add --no-cache certbot py2-future py2-pip
RUN certbot certonly --standalone -d ${DNS_HOST} --email ${YOUR_EMAIL} --agree-tos --quiet
RUN mv /etc/letsencrypt/archive/fullchain1.pem /etc/ssl/fullchain.pem
RUN mv /etc/letsencrypt/archive/privkey1.pem /etc/ssl/privkey.pem
COPY nginx.conf /etc/nginx/nginx.conf
COPY dhparams.pem /etc/ssl/dhparams.pem

## These 2 are only if you create the certificates outside the container and just want to copy them in
#COPY fullchain1.pem /etc/ssl/fullchain.pem
#COPY privkey1.pem /etc/ssl/privkey.pem

CMD ["nginx"]