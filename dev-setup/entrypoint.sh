#!/bin/sh

envsubst < /app/config.template.yaml > /tmp/config.yaml

exec "$@"