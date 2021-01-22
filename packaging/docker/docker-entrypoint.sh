#!/bin/sh
set -e

if [ "$1" = 'discovery' ]; then
    echo "starting discovery"
    exec /discovery server
fi

exec "$@"
