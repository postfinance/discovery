FROM docker.repo.pnet.ch/pf/alpine:latest

ENV TZ=Europe/Zurich

RUN http_proxy=http://outapplpf.pnet.ch:3128 apk add --no-cache tzdata

COPY discovery /discovery
COPY packaging/docker/docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["discovery"]

EXPOSE 3001
