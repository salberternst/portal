FROM quay.io/oauth2-proxy/oauth2-proxy:v7.5.1

# this image is just for development
# enable root access in the container again
USER 0:0

RUN apk add --no-cache gettext curl
COPY ./entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod a+x /usr/local/bin/entrypoint.sh

COPY ./config.template.yaml /app/config.template.yaml

HEALTHCHECK --interval=120s --timeout=20s CMD curl -f http://localhost:4180 || exit 1

# UID/GID 65532 is also known as nonroot user in distroless image
USER 65532:65532

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
