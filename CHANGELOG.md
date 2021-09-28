## 0.8.0 (2021-09-28)


### Bug Fixes

* **discoveryd**: return `codes.NotFound` intead of `codes.Internal` not existing entities ([d3e63e61](https://github.com/postfinance/discovery/commit/d3e63e61))


### Dependencies

* **apimachinery**: 0.22.0 -> 0.22.1 ([d466678c](https://github.com/postfinance/discovery/commit/d466678c))
* **apimachinery**: 0.22.1 -> 0.22.2 ([32997363](https://github.com/postfinance/discovery/commit/32997363))
* **genproto**: 0.0.0-20210804223703-f1db76f3300d -> 0.0.0-20210820002220-43fce44e7af1 ([6f114595](https://github.com/postfinance/discovery/commit/6f114595))
* **genproto**: 0.0.0-20210903162649-d08c68adba83 -> 0.0.0-20210927142257-433400c27d05 ([54ae20ad](https://github.com/postfinance/discovery/commit/54ae20ad))
* **go**: 1.16 -> 1.17 ([65a5c976](https://github.com/postfinance/discovery/commit/65a5c976))
* **go-oidc**: 3.0.0 -> 3.1.0 ([cb4efdc0](https://github.com/postfinance/discovery/commit/cb4efdc0))
* **grpc**: 1.39.0 -> 1.40.0 ([3a8938e1](https://github.com/postfinance/discovery/commit/3a8938e1))
* **grpc**: 1.40.0 -> 1.41.0 ([ddf370c9](https://github.com/postfinance/discovery/commit/ddf370c9))
* **grpc-gateway**: 2.5.0 -> 2.6.0 ([749eb613](https://github.com/postfinance/discovery/commit/749eb613))
* **jwt**: 4.0.0 -> 4.1.0 ([c1c62c2f](https://github.com/postfinance/discovery/commit/c1c62c2f))
* **king**: 0.2.0 -> 0.3.0 ([ecbbf305](https://github.com/postfinance/discovery/commit/ecbbf305))
* **term**: 0.0.0-20210615171337-6886f2dfbf5b -> 0.0.0-20210927222741-03fcf44c2211 ([998d10d8](https://github.com/postfinance/discovery/commit/998d10d8))
* **zap**: 1.18.1 -> 1.19.0 ([b34b8678](https://github.com/postfinance/discovery/commit/b34b8678))
* **zap**: 1.19.0 -> 1.19.1 ([2e1102e0](https://github.com/postfinance/discovery/commit/2e1102e0))


### New Features

* **cli**: add the possibility to list and unregister all unresolvable services ([#21](https://github.com/postfinance/discovery/issues/21), [87301176](https://github.com/postfinance/discovery/commit/87301176))
* **discoveryd**: support of prometheus http_sd ([#23](https://github.com/postfinance/discovery/issues/23), [4862d3d1](https://github.com/postfinance/discovery/commit/4862d3d1))
  > Now it is possible to use the rest endpoint
  > `/v1/sd/<prometheus-server>/<namespace>` as
  > http service discovery. See Readme for more
  > information.



## 0.7.8 (2021-08-05)


### Bug Fixes

* **common**: failed to parse machine token errors ([2cf360ba](https://github.com/postfinance/discovery/commit/2cf360ba))
  > After migrating to an updated jwt package, tokens created with the old
  > library could not be parsed anymore. Now parsing works for old and new
  > generation jwt tokens.



## 0.7.7 (2021-08-05)


### Bug Fixes

* **common**: update jwt lib to fix a security issue ([8417051b](https://github.com/postfinance/discovery/commit/8417051b))


### Dependencies

* **apimachinery**: 0.21.2 -> 0.22.0 ([ac6534c8](https://github.com/postfinance/discovery/commit/ac6534c8))
* **genproto**: 0.0.0-20210629200056-84d6f6074151 -> 0.0.0-20210804223703-f1db76f3300d ([4e2413c8](https://github.com/postfinance/discovery/commit/4e2413c8))
* **king**: 0.1.0 -> 0.2.0 ([baff26c5](https://github.com/postfinance/discovery/commit/baff26c5))



## 0.7.6 (2021-08-03)


### Bug Fixes

* **cli**: register service failed when no labels were used ([1a38a07f](https://github.com/postfinance/discovery/commit/1a38a07f))



## 0.7.5 (2021-06-30)


### Dependencies

* **apimachinery**: 0.21.0 -> v0.21.2 ([29da097f](https://github.com/postfinance/discovery/commit/29da097f))
* **genproto**: v0.0.0-20210617175327-b9e0b3197ced -> v0.0.0-20210629200056-84d6f6074151 ([4d78ae7c](https://github.com/postfinance/discovery/commit/4d78ae7c))
* **grpc**: 1.38.0 -> 1.39.0 ([9bb9e423](https://github.com/postfinance/discovery/commit/9bb9e423))
* **grpc-gateway**: v2.4.0 -> v2.5.0 ([4f7c0e53](https://github.com/postfinance/discovery/commit/4f7c0e53))
* **kong**: master -> 0.2.17 ([4a5575ba](https://github.com/postfinance/discovery/commit/4a5575ba))
* **oauth2**: 0.0.0-20210514164344-f6687ab2804c -> 0.0.0-20210628180205-a41e5a781914 ([b2ee8c69](https://github.com/postfinance/discovery/commit/b2ee8c69))
* **protobuf**: v1.26.0 -> v1.27.1 ([268646c2](https://github.com/postfinance/discovery/commit/268646c2))
* **store**: update from 0.2.0-pre to 0.2.0 ([0dc85d94](https://github.com/postfinance/discovery/commit/0dc85d94))
* **term**: 0.0.0-20210503060354-a79de5458b56 -> 0.0.0-20210615171337-6886f2dfbf5b ([812f0749](https://github.com/postfinance/discovery/commit/812f0749))
* **zap**: v1.17.0 -> v1.18.1 ([dcb32b0d](https://github.com/postfinance/discovery/commit/dcb32b0d))



## 0.7.4 (2021-05-19)


### Bug Fixes

* **discoveryd**: use correct grpc codes on validation and not found errors ([bcc67252](https://github.com/postfinance/discovery/commit/bcc67252))



## 0.7.3 (2021-05-18)


### Bug Fixes

* **discoveryd**: use grpc log middleware to log grpc methods ([e4524024](https://github.com/postfinance/discovery/commit/e4524024))



## 0.7.2 (2021-05-18)


### Bug Fixes

* **discoveryd**: return `NotFound` grpc error instead of `Internal` or `InvalidArgument` when namespace or servers not found ([c342497a](https://github.com/postfinance/discovery/commit/c342497a))



## 0.7.1 (2021-05-18)


### Bug Fixes

* **common**: increase server register timeout ([3457a2f3](https://github.com/postfinance/discovery/commit/3457a2f3))
* **common**: move import command to service subcommand ([4177c181](https://github.com/postfinance/discovery/commit/4177c181))
* **discoveryd**: `discovery_services_count` reports now the correct number ([136f77a9](https://github.com/postfinance/discovery/commit/136f77a9))
* **exporter**: no more false positive `failed to delete service` error messages ([#20](https://github.com/postfinance/discovery/issues/20), [ca215651](https://github.com/postfinance/discovery/commit/ca215651))



## 0.7.0 (2021-04-23)


### Bug Fixes

* **exporter**: show flags in metrics endpoint ([61745826](https://github.com/postfinance/discovery/commit/61745826))


### New Features

* **discoveryd**: add namespace to services count metric ([bd9a8186](https://github.com/postfinance/discovery/commit/bd9a8186))
* **exporter**: add go and process metrics ([18bdb0c6](https://github.com/postfinance/discovery/commit/18bdb0c6))
* **exporter**: add prometheus metrics endpoint ([ab758867](https://github.com/postfinance/discovery/commit/ab758867))



## 0.6.0 (2021-04-22)


### Bug Fixes

* **discoveryd**: remove `etcd-ca` and `etcd-cert` from metrics ([d4cfe276](https://github.com/postfinance/discovery/commit/d4cfe276))
  > This values can be pretty large. If you need to see the value you can
  > refer to the logs.


### New Features

* **discoveryd**: add logger prometheus metrics ([e0656157](https://github.com/postfinance/discovery/commit/e0656157))



## 0.5.3 (2021-04-21)


### Bug Fixes

* **common**: validate service label names and label values ([#18](https://github.com/postfinance/discovery/issues/18), [f8198afc](https://github.com/postfinance/discovery/commit/f8198afc))



## 0.5.2 (2021-04-14)


### Bug Fixes

* **common**: update kong dependency to master ([#17](https://github.com/postfinance/discovery/issues/17), [c24f227d](https://github.com/postfinance/discovery/commit/c24f227d))
* **common**: use correct command description for server and client ([978e507e](https://github.com/postfinance/discovery/commit/978e507e))



## 0.5.1 (2021-04-13)


### Bug Fixes

* **common**: update to newest gprc version ([2e97b74c](https://github.com/postfinance/discovery/commit/2e97b74c))



## 0.5.0 (2021-04-07)


### Bug Fixes

* **exporter**: remove service from export file on unregister ([#16](https://github.com/postfinance/discovery/issues/16), [bb56b1ed](https://github.com/postfinance/discovery/commit/bb56b1ed))


### New Features

* **cli**: add possibility to filter services ([#15](https://github.com/postfinance/discovery/issues/15), [b0859607](https://github.com/postfinance/discovery/commit/b0859607))
* **cli**: add possibility to sort services by endpoint or modification date ([08575b30](https://github.com/postfinance/discovery/commit/08575b30))



## 0.4.1 (2021-03-26)


### Bug Fixes

* **rest**: use URL query parameter to unregister service by endpoint ([f0a82b74](https://github.com/postfinance/discovery/commit/f0a82b74))
  > It is now possible to unregister a service by endpoint URL. Previously it was only possible to unregister
  > a service via REST by ID.



## 0.4.0 (2021-03-17)


### New Features

* **discovery**: allow to register multiple endpoints at once ([4dde6860](https://github.com/postfinance/discovery/commit/4dde6860))
  > In the `service unregister` subcommand, we also changed the environment variable from
  > `DISCOVERY_SERVICES` to `DISCOVERY_ENDPOINTS`. This makes the
  > register and unregister subcommands more consistent.



## 0.3.0 (2021-03-16)


### Bug Fixes

* **common**: cli returns error (instead of warning) when server is found ([#10](https://github.com/postfinance/discovery/issues/10), [e2300ce7](https://github.com/postfinance/discovery/commit/e2300ce7))
  > When registering a service with a selector that would result in a
  > service registration without corresponding server, the command
  > fails with `no server found for selector '<selector>'` error message.
* **common**: correctly use environment variable `DISOVERY_NAME` for service name ([#12](https://github.com/postfinance/discovery/issues/12), [af3c41b9](https://github.com/postfinance/discovery/commit/af3c41b9))
* **discovery**: do not start oidc client for machine tokens ([7a4f447b](https://github.com/postfinance/discovery/commit/7a4f447b))
* **discovery**: not printing usage when service name is empty ([671d3cee](https://github.com/postfinance/discovery/commit/671d3cee))
* **exporter**: create namespace export directories on startup ([#11](https://github.com/postfinance/discovery/issues/11), [aa7895df](https://github.com/postfinance/discovery/commit/aa7895df))


### New Features

* **discovery**: add command aliases ns, svc and svr ([32d1f51a](https://github.com/postfinance/discovery/commit/32d1f51a))



## 0.2.3 (2021-03-11)


### Bug Fixes

* **common**: remove namespace config type from export path ([#9](https://github.com/postfinance/discovery/issues/9), [0b17742f](https://github.com/postfinance/discovery/commit/0b17742f))



## 0.2.2 (2021-02-19)


### Bug Fixes

* **discoveryd**: namespace cache works with multiple discovery instances ([#8](https://github.com/postfinance/discovery/issues/8), [548c01b5](https://github.com/postfinance/discovery/commit/548c01b5))



## 0.2.1 (2021-02-12)


### Bug Fixes

* **common**: add exporter and discovery systemd files ([f15a4e81](https://github.com/postfinance/discovery/commit/f15a4e81))
* **common**: allow '-' character in etcd prefix. ([1d48aa02](https://github.com/postfinance/discovery/commit/1d48aa02))
* **exporter**: update to new etcd store package ([f5ed78ad](https://github.com/postfinance/discovery/commit/f5ed78ad))



## 0.2.0 (2021-02-12)


### New Features

* **common**: make etcd prefix configurable ([647b5fe3](https://github.com/postfinance/discovery/commit/647b5fe3))
* **discoveryd**: serve swagger-ui with grpc-gateway api doc ([2ca0e337](https://github.com/postfinance/discovery/commit/2ca0e337))



## 0.1.0 (2021-02-03)

### New Features

* **common**: add new option `--ca-cert` to specify additional ca certifiactes ([c9668c18](https://github.com/postfinance/discovery/commit/c9668c18))



## 0.0.1 initial release (2021-02-03)

### Bug Fixes

* **common**: actually use selector in import command ([d170490b](https://github.com/postfinance/discovery/commit/d170490b))
* **common**: add new stuff ([c73812e2](https://github.com/postfinance/discovery/commit/c73812e2))
* **common**: do not allow unregister server that has registered services ([a2b6308c](https://github.com/postfinance/discovery/commit/a2b6308c))
* **common**: exporter ([c3a8aaf3](https://github.com/postfinance/discovery/commit/c3a8aaf3))
* **common**: metrics count ([373298c0](https://github.com/postfinance/discovery/commit/373298c0))
* **discovery**: use correct token path ([07426526](https://github.com/postfinance/discovery/commit/07426526))
* **discoveryd**: fix authorization ([39290b97](https://github.com/postfinance/discovery/commit/39290b97))
* **discoveryd**: initialize prometheus metrics ([0a1302fb](https://github.com/postfinance/discovery/commit/0a1302fb))
* **exporter**: export service as blackbox or standard via namespace config ([5451a182](https://github.com/postfinance/discovery/commit/5451a182))
* **exporter**: ingore events for services not containing configured server ([eb8b2cd7](https://github.com/postfinance/discovery/commit/eb8b2cd7))
* **exporter**: only export services for configured server ([9b1e3d2b](https://github.com/postfinance/discovery/commit/9b1e3d2b))
* **exporter**: rewrite files when namespace exportconfig changes ([000a0c1a](https://github.com/postfinance/discovery/commit/000a0c1a))
* **server**: metric discovery_services_count metrics is now correctly calculated ([#5](https://github.com/postfinance/discovery/issues/5), [c8ec60c0](https://github.com/postfinance/discovery/commit/c8ec60c0))

### New Features

* **common**: add possibility to enable/disable a server ([fd0a9076](https://github.com/postfinance/discovery/commit/fd0a9076))
* **common**: add possibility to select server with selector in import command ([8fad2bce](https://github.com/postfinance/discovery/commit/8fad2bce))
* **common**: add the possibility to use a selector for server selection ([7458ea08](https://github.com/postfinance/discovery/commit/7458ea08))
* **common**: basic working implementation ([797a8197](https://github.com/postfinance/discovery/commit/797a8197))
* **common**: stricter parsing for endpoint url ([b2e98fe4](https://github.com/postfinance/discovery/commit/b2e98fe4))
* **discoveryd**: add grpc reflection support ([7bd5a2dc](https://github.com/postfinance/discovery/commit/7bd5a2dc))
* **discoveryd**: add oidc authentication ([95271cca](https://github.com/postfinance/discovery/commit/95271cca))
* **discoveryd**: add prometheus metric showing the configured replication factor ([#4](https://github.com/postfinance/discovery/issues/4), [3f179e6e](https://github.com/postfinance/discovery/commit/3f179e6e))
