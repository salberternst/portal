FROM quay.io/oauth2-proxy/oauth2-proxy

# this image is just for development
# enable root access in the container again
USER 0:0

RUN apk add --no-cache gettext
COPY ./entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]