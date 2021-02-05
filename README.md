<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [discovery](#discovery)
  - [Architecture](#architecture)
    - [Example Workflow:](#example-workflow)
    - [Authentication](#authentication)
    - [Configuration](#configuration)
    - [API](#api)
    - [Systemd](#systemd)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# discovery
Service discovery for prometheus with etcd backend. This service can be useful in environments
where no prometheus service discovery other than [file-sd](https://prometheus.io/docs/guides/file-sd/) is possible.

## Architecture

![Architecture](./architecture.svg)

The service discovery consists of three components:

* A GPRC service to register services and store them to [etcd](https://etcd.io) backend.
* An exporter service to export stored services to filesystem for promehteus [file-sd](https://prometheus.io/docs/guides/file-sd/).
* CLI to register or unregister services and perform admin tasks.

### Example Workflow:

First we have to create a (prometheus) server:

```
$ discovery server register prometheus1.example.com --labels=environment=test
```

To list all registered servers:
```
$ discovery server list
NAME                    MODIFIED             STATE  LABELS
prometheus1.example.com 2021-02-04T14:12:50Z active environment=test
```

Now you can register a service:

```
$ discovery service register -e http://example.com/metrics example --labels=label1=value1,label2=value2 --selector=environment=test
2021-02-04T15:13:50.085+0100    INFO    client/service.go:82    service registered      {"id": "93c156b1-f218-5d79-88a5-219307e59d29"}
```

The selector is a kubernetes style [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) to select a server from the registered servers.
When you start the discovery service with a number of replicas n>1, a service is distributed to n servers with the corresponding labels. The service discovery uses a [consistent hashing algorithm](https://arxiv.org/pdf/1406.2294v1.pdf)
to distribute services among servers.

You can see the regsitered services:

```
$ discovery service list
NAME    NAMESPACE ID                                   ENDPOINT                   SERVERS                 LABELS                      SELECTOR         MODIFIED             DESCRIPTION
example default   93c156b1-f218-5d79-88a5-219307e59d29 http://example.com/metrics prometheus1.example.com label1=value1,label2=value2 environment=test 2021-02-04T14:13:50Z
```

You can specify a namespace with `-n`. The namespace configures how the service is exported (standard or blackbox).

To view all configured namespaces:

```
$ discovery namespace list
NAME    EXPORTCONFIG MODIFIED
default standard     2021-02-04T14:07:03Z
```

Register a blackbox namespace:

```
$ discovery namespace register -e blackbox default-blackbox
$ discovery namespace list
NAME             EXPORTCONFIG MODIFIED
default          standard     2021-02-05T07:58:35Z
default-blackbox blackbox     2021-02-05T08:00:18Z
```

Now we can register a blackbox service:

```
$ discovery service register -e http://blackbox.example.com -n default-blackbox blackbox -s environment=test
$ discovery service list
NAME     NAMESPACE        ID                                   ENDPOINT                    SERVERS                 LABELS                      SELECTOR         MODIFIED             DESCRIPTION
blackbox default-blackbox e988791c-2c4e-5eeb-b3b8-db3c0cf82719 http://blackbox.example.com prometheus1.example.com                             environment=test 2021-02-05T08:05:44Z
example  default          93c156b1-f218-5d79-88a5-219307e59d29 http://example.com/metrics  prometheus1.example.com label1=value1,label2=value2 environment=test 2021-02-05T07:59:17Z
```

Now you can start the exporter service on the corresponding promehteus server:

```
$ discoveryd exporter --directory=/tmp/discovery --server=prometheus1.example.com
2021-02-05T09:08:29.671+0100    INFO    server/exporter.go:28   starting exporter       {"buildinfo-date": "2021-02-04T06:48:33Z", "buildinfo-go": "go1.15.5", "buildinfo-revision": "1c8d0652", "buildinfo-version": "v0.1.0-SNAPSHOT-1c8d065", "debug": false, "directory": "/tmp/discovery", "etcd-auto-sync-interval": "10s", "etcd-ca": "", "etcd-ca-file": "", "etcd-cert": "", "etcd-cert-file": "", "etcd-dial-timeout": "5s", "etcd-endpoints": ["localhost:2379"], "etcd-key": "", "etcd-key-file": "", "etcd-password": "", "etcd-prefix": "/discovery", "etcd-request-timeout": "5s", "etcd-user": "", "profiler-enabled": false, "profiler-listen": ":6666", "profiler-timeout": "5m0s", "resync-interval": "1h0m0s", "server": "prometheus1.example.com", "show-config": false}
2021-02-05T09:08:29.675+0100    INFO    exporter/exporter.go:76 sync services
2021-02-05T09:08:29.677+0100    INFO    exporter/file.go:202    updating discovery file {"path": "/tmp/discovery/prometheus1.example.com/blackbox/default-blackbox/blackbox.json"}
2021-02-05T09:08:29.677+0100    INFO    exporter/file.go:202    updating discovery file {"path": "/tmp/discovery/prometheus1.example.com/standard/default/example.json"}
```

As you can see, the exporter created two files which can be used by prometheus for file_sd. The created files have different content depending on the namespace export configuration (`standard` or `blackbox`):

```
$ cat /tmp/discovery/prometheus1.example.com/standard/default/example.json |jq
[
  {
    "targets": [
      "example.com"
    ],
    "labels": {
      "__metrics_path__": "/metrics",
      "__scheme__": "http",
      "instance": "example.com",
      "job": "example",
      "label1": "value1",
      "label2": "value2"
    }
  }
]
```

```
$ cat /tmp/discovery/prometheus1.example.com/blackbox/default-blackbox/blackbox.json |jq
[
  {
    "targets": [
      "http://blackbox.example.com"
    ]
  }
]
```

### Authentication
Discovery is meant to work with an openid connect server (Password Grant Flow). The following options exist for configuration:

```
--oidc-endpoint=STRING                       OIDC endpoint URL ($DISCOVERY_OIDC_ENDPOINT).
--oidc-client-id=STRING                      OIDC client ID ($DISCOVERY_OIDC_CLIENT_ID).
--oidc-roles=OIDC-ROLES,...                  The the roles that are allowed to change servers and namespaces and to issue machine tokens ($DISCOVERY_OIDC_ROLES).
--ca-cert=STRING                             Path to a custom tls ca pem file. Certificates in this file are added to system cert pool ($DISCOVERY_CA_CERT).
```

With `--oidc-roles` you can specify a comma separated list of roles, that can register servers, namespaces and services.

To login run:

```
$ discovery login
```

On successful login the token is saved to `~/.config/discovery/.token` for all subsequent requests.

You can create machine tokens with:

```
$ token create -n default,default-blackbox ansible-user
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhbnNpYmxlLXVzZXIiLCJpYXQiOjE2MTI1MTM1ODkuMDA3ODY2OSwiaXNzIjoidGhlc2VjcmV0IiwibmJmIjoxNjEyNTEzNTg5LjAwNzg2NjksIm5hbWVzcGFjZXMiOlsiZGVmYXVsdCIsImRlZmF1bHQtYmxhY2tib3giXX0.IUKFuyKMAU5aRZJPLp67Uei9o2G5neJz_Ha86JZnd8o
```

The above command allows that token to register services in `default` and `default-blackbox` namespaces. `ansible-user` is an ID to identify the token.

To view a token run:

```
$ discovery token info eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhbnNpYmxlLXVzZXIiLCJpYXQiOjE2MTI1MTM1ODkuMDA3ODY2OSwiaXNzIjoidGhlc2VjcmV0IiwibmJmIjoxNjEyNTEzNTg5LjAwNzg2NjksIm5hbWVzcGFjZXMiOlsiZGVmYXVsdCIsImRlZmF1bHQtYmxhY2tib3giXX0.IUKFuyKMAU5aRZJPLp67Uei9o2G5neJz_Ha86JZnd8o
id: ansible-user
namespaces: [default default-blackbox]
expiry: never
```

### Configuration

Every flag can be set with environment variables. Run `discovery --help` to check which variables are available. It is also possible to use yaml configuration files. You can check which config files are used with:

```
$ discovery --show-config
Configuration files:
  /home/user/.config/discovery/config.yaml                       parsed
  /etc/discovery/config.yaml                                     not found
```

Example config:

```yaml
selector: zone=default
oidc-client-id: client-id
oidc-roles:
  - role1
  - role2
oidc-endpoint: https://auth.example.com/auth/realms/discovery
```

### API
The service discovery has a grpc API. The proto files are [here](./pkg/discoverypb). The generated go grpc code is also in that directory.

### Systemd
It is possible to register and unregister services on start/stop with systemd. An example:

```
ExecStartPost=/usr/bin/discovery register --endpoint=https://${ETCD_NAME}.pnet.ch:7002/metrics etcdv3
ExecStopPost=/usr/bin/discovery unregister --endpoint=https://${ETCD_NAME}.pnet.ch:7002/metrics
```
