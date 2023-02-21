FROM linux-docker-local.repo.pnet.ch/pf/distroless/static:latest

COPY discoveryd /usr/bin/discoveryd

USER 65534

ENTRYPOINT ["/usr/bin/discoveryd"]
