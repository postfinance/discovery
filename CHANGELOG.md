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
