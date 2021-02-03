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