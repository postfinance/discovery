FROM alpine:latest

COPY discovery /usr/bin/discovery
COPY discoveryd /usr/bin/discoveryd
COPY packaging/docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod 750 /docker-entrypoint.sh

CMD ["server"]
ENTRYPOINT ["/docker-entrypoint.sh"]
