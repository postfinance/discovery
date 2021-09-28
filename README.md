<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [discovery](#discovery)
  - [Architecture](#architecture)
  - [Example Workflow:](#example-workflow)
  - [Authentication](#authentication)
  - [Configuration](#configuration)
  - [API](#api)
    - [GRPC](#grpc)
    - [REST](#rest)
  - [Prometheus Scrape Configuration](#prometheus-scrape-configuration)
    - [http_sd](#http_sd)
    - [file_sd](#file_sd)
  - [Systemd](#systemd)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# discovery
Service discovery for prometheus with etcd backend. This service can be useful in environments
where no prometheus service discovery other than [file-sd](https://prometheus.io/docs/guides/file-sd/) is possible.

## Architecture

![Architecture](./architecture.svg)

The service discovery consists of three components:

* A GPRC service to register services and store them to [etcd](https://etcd.io) backend.
* A rest endpoint for prometheus [http_sd](https://prometheus.io/docs/prometheus/latest/http_sd/).
* An exporter service to export stored services to filesystem for prometheus [file-sd](https://prometheus.io/docs/guides/file-sd/).
* CLI to register or unregister services and perform admin tasks.

## Example Workflow:

First we have to register a (prometheus) server:

```console
$ discovery server register prometheus1.example.com --labels=environment=test
```
The labels above can be used on service registration to select a server via kubernetes style label selectors (see below).

To list all registered servers:
```console
$ discovery server list
NAME                    MODIFIED             STATE  LABELS
prometheus1.example.com 2021-02-04T14:12:50Z active environment=test
```

Now you can register a service:

```console
$ discovery service register -e http://example.com/metrics example --labels=label1=value1,label2=value2 --selector=environment=test
2021-02-04T15:13:50.085+0100    INFO    client/service.go:82    service registered      {"id": "93c156b1-f218-5d79-88a5-219307e59d29"}
```

The selector is a kubernetes style [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) to select a server from the registered servers.
When you start the discovery service with `--replicas=n` and n>1, a service is distributed to n servers with the corresponding labels. The service discovery uses a [consistent hashing algorithm](https://arxiv.org/pdf/1406.2294v1.pdf)
to distribute services among servers.

You can see the regsitered services:

```console
$ discovery service list
NAME    NAMESPACE ID                                   ENDPOINT                   SERVERS                 LABELS                      SELECTOR         MODIFIED             DESCRIPTION
example default   93c156b1-f218-5d79-88a5-219307e59d29 http://example.com/metrics prometheus1.example.com label1=value1,label2=value2 environment=test 2021-02-04T14:13:50Z
```

You can specify a namespace with `-n`. The namespace defines how its services are exported (standard or blackbox). You can also create namespaces to group services. The service endpoint is unique per namepace.

To view all configured namespaces:

```console
$ discovery namespace list
NAME    EXPORTCONFIG MODIFIED
default standard     2021-02-04T14:07:03Z
```

Register a blackbox namespace:

```console
$ discovery namespace register -e blackbox default-blackbox
$ discovery namespace list
NAME             EXPORTCONFIG MODIFIED
default          standard     2021-02-05T07:58:35Z
default-blackbox blackbox     2021-02-05T08:00:18Z
```

Now we can register a blackbox service:

```console
$ discovery service register -e http://blackbox.example.com -n default-blackbox blackbox -s environment=test
$ discovery service list
NAME     NAMESPACE        ID                                   ENDPOINT                    SERVERS                 LABELS                      SELECTOR         MODIFIED             DESCRIPTION
blackbox default-blackbox e988791c-2c4e-5eeb-b3b8-db3c0cf82719 http://blackbox.example.com prometheus1.example.com                             environment=test 2021-02-05T08:05:44Z
example  default          93c156b1-f218-5d79-88a5-219307e59d29 http://example.com/metrics  prometheus1.example.com label1=value1,label2=value2 environment=test 2021-02-05T07:59:17Z
```

Now you can start the exporter service on the corresponding prometheus server:

```console
$ discoveryd exporter --directory=/tmp/discovery --server=prometheus1.example.com
2021-02-05T09:08:29.671+0100    INFO    server/exporter.go:28   starting exporter       {"buildinfo-date": "2021-02-04T06:48:33Z", "buildinfo-go": "go1.15.5", "buildinfo-revision": "1c8d0652", "buildinfo-version": "v0.1.0-SNAPSHOT-1c8d065", "debug": false, "directory": "/tmp/discovery", "etcd-auto-sync-interval": "10s", "etcd-ca": "", "etcd-ca-file": "", "etcd-cert": "", "etcd-cert-file": "", "etcd-dial-timeout": "5s", "etcd-endpoints": ["localhost:2379"], "etcd-key": "", "etcd-key-file": "", "etcd-password": "", "etcd-prefix": "/discovery", "etcd-request-timeout": "5s", "etcd-user": "", "profiler-enabled": false, "profiler-listen": ":6666", "profiler-timeout": "5m0s", "resync-interval": "1h0m0s", "server": "prometheus1.example.com", "show-config": false}
2021-02-05T09:08:29.675+0100    INFO    exporter/exporter.go:76 sync services
2021-02-05T09:08:29.677+0100    INFO    exporter/file.go:202    updating discovery file {"path": "/tmp/discovery/prometheus1.example.com/default-blackbox/blackbox.json"}
2021-02-05T09:08:29.677+0100    INFO    exporter/file.go:202    updating discovery file {"path": "/tmp/discovery/prometheus1.example.com/default/example.json"}
```

As you can see, the exporter created two files which can be used by prometheus for file_sd. The created files have different content depending on the namespace export configuration (`standard` or `blackbox`):

```console
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

```console
$ cat /tmp/discovery/prometheus1.example.com/blackbox/default-blackbox/blackbox.json |jq
[
  {
    "targets": [
      "http://blackbox.example.com"
    ]
  }
]
```

## Authentication
Discovery is meant to work with an openid connect server (Password Grant Flow). The following options exist for configuration:

```
--oidc-endpoint=STRING                       OIDC endpoint URL ($DISCOVERY_OIDC_ENDPOINT).
--oidc-client-id=STRING                      OIDC client ID ($DISCOVERY_OIDC_CLIENT_ID).
--oidc-roles=OIDC-ROLES,...                  The the roles that are allowed to change servers and namespaces and to issue machine tokens ($DISCOVERY_OIDC_ROLES).
--ca-cert=STRING                             Path to a custom tls ca pem file. Certificates in this file are added to system cert pool ($DISCOVERY_CA_CERT).
```

With `--oidc-roles` you can specify a comma separated list of roles, that can register servers, namespaces and services. Those roles can also issue machine access tokens.

To login run:

```console
$ discovery login
```

On successful login the token is saved to `~/.config/discovery/.token` for all subsequent requests.

You can create machine tokens with:

```console
$ token create -n default,default-blackbox ansible-user
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhbnNpYmxlLXVzZXIiLCJpYXQiOjE2MTI1MTM1ODkuMDA3ODY2OSwiaXNzIjoidGhlc2VjcmV0IiwibmJmIjoxNjEyNTEzNTg5LjAwNzg2NjksIm5hbWVzcGFjZXMiOlsiZGVmYXVsdCIsImRlZmF1bHQtYmxhY2tib3giXX0.IUKFuyKMAU5aRZJPLp67Uei9o2G5neJz_Ha86JZnd8o
```

The above command creates a token with access to `default` and `default-blackbox` namespaces. `ansible-user` is an ID to identify the token.

To view a token run:

```console
$ discovery token info eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhbnNpYmxlLXVzZXIiLCJpYXQiOjE2MTI1MTM1ODkuMDA3ODY2OSwiaXNzIjoidGhlc2VjcmV0IiwibmJmIjoxNjEyNTEzNTg5LjAwNzg2NjksIm5hbWVzcGFjZXMiOlsiZGVmYXVsdCIsImRlZmF1bHQtYmxhY2tib3giXX0.IUKFuyKMAU5aRZJPLp67Uei9o2G5neJz_Ha86JZnd8o
id: ansible-user
namespaces: [default default-blackbox]
expiry: never
```

You can use the created token for automating tasks. You can set the token as follows in `~/.config/discovery/.token`:

```yaml
machine_token: <generated jwt token>
```

With the above token you can regster services only (no namespaces or servers can be registered)

## Configuration

Every flag can be set with environment variables. Run `discovery --help` to check which variables are available. It is also possible to use yaml configuration files. You can check which config files are used with:

```console
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

## API

### GRPC
The service discovery has a GRPC API. The proto files can be found [here](./proto/postfinance/discovery/v1). The generated GRPC go code is also in that directory. To send the authorization token with the go client you
can use an [UnaryClientInterceptor](https://github.com/grpc/grpc-go/blob/master/interceptor.go) like below:

```go
func buildClientInterceptor(token string) func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, grpc.UnaryInvoker, ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
```

### REST
It is also possible to access the a rest api generated with [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). The swagger api documentation is available under http://localhost:3002/swagger/
(when running the server with default settings).

To access the API you have to set an authorization header:

```console
curl -X GET "http://localhost:3002/v1/namespaces" -H  "accept: application/json" -H " -H "authorization: bearer <token>""
```

To unregister a service by endpoint URL run:

```console
curl -G -v -X DELETE --data-urlencode "id=http://example.com:9182/metrics" -H "accept: application/json" -H "authorization: bearer $TOKEN" http://localhost:3002/v1/services/windows
```

which is the same as:

```console
curl -v -X DELETE -H "accept: application/json" -H "authorization: bearer $TOKEN" 'http://localhost:3002/v1/services/windows?id=http%3A%2F%2Fexample.com%3A9182%2Fmetrics'
```

## Prometheus Scrape Configuration

### http_sd
The following is an example on how to configure prometheus for http_sd:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'http'
    http_sd_configs:
      - url: https://<service-discovery-url>:<port>/v1/sd/prometheus1.example.com/default?export-config=standard
        authorization:
          type: Bearer
          credentials: <jwt token>
```

### file_sd
The following is an example on how to configure prometheus for file_sd:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'file'
    file_sd_configs:     
      - files: [/tmp/discovery/prometheus1.example.com/default/*.json]
```

## Systemd
It is possible to register and unregister services on start/stop with systemd. An example for auto registering [node_exporter](https://github.com/prometheus/node_exporter):

* `/usr/lib/systemd/system/node_exporter.service`:

```
[Unit]
Description=prometheus node exporter
Wants=network-online.target
After=network-online.target

[Service]
EnvironmentFile=/etc/sysconfig/node_exporter
ExecStart=/usr/bin/node_exporter --web.listen-address=${LISTEN} --log.level=${LOGLEVEL} --collector.netstat.fields=(.*) --collector.processes
ExecStartPost=-/bin/discovery service register ${DISCOVERY_NAME}
ExecStopPost=-/bin/discovery service unregister

[Install]
WantedBy=multi-user.target
```


* `/etc/sysconfig/node_exporter`:

```yaml
LISTEN=:9549
LOGLEVEL=info
DISCOVERY_NAMESPACE=systemd
DISCOVERY_ENDPOINTS="<endpoint_url>"
DISCOVERY_NAME=node
```

* `etc/discovery/config.yaml`:

```yaml
address: <discovery-service-url>
oidc-endpoint: <oidc-endpoint-url>
oidc-client-id: <oidc-endpoint-client-id>
selector: environment=test
```

Create a token and add it to `/root/.config/discovery/.token`:

```yaml
machine_token: <jwt-token>
```

After you have run `systemctl daemon-reload` the service registers and unregisters automatically on `systemctl start|stop node_exporter`.
