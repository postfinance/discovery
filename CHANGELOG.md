## 0.9.0 (2023-02-21)


### Build System

* **deps**: bump github.com/golang-jwt/jwt/v4 from 4.4.3 to 4.5.0 (#115) ([d393d692](https://github.com/postfinance/discovery/commit/d393d692))
* **deps**: bump golang.org/x/oauth2 from 0.4.0 to 0.5.0 (#112) ([2af7d82c](https://github.com/postfinance/discovery/commit/2af7d82c))
* **deps**: bump golang.org/x/term from 0.4.0 to 0.5.0 (#114) ([bc546bac](https://github.com/postfinance/discovery/commit/bc546bac))
* **deps**: bump google.golang.org/grpc from 1.52.1 to 1.52.3 (#111) ([1b943933](https://github.com/postfinance/discovery/commit/1b943933))
* **deps**: bump google.golang.org/grpc from 1.52.3 to 1.53.0 (#113) ([b693f0d1](https://github.com/postfinance/discovery/commit/b693f0d1))


### New Features

* **common**: add possibility to use an external command to create the oidc id_token ([a3654f1c](https://github.com/postfinance/discovery/commit/a3654f1c))
* **common**: building docker image from distroless ([8c729dd5](https://github.com/postfinance/discovery/commit/8c729dd5))
* **common**: retry register and unregister service exponentially ([2dfd5651](https://github.com/postfinance/discovery/commit/2dfd5651))



## 0.8.3 (2023-01-25)


### Build System

* **deps**: bump github.com/alecthomas/kong from 0.5.0 to 0.6.0 (#65) ([a1dfa4b3](https://github.com/postfinance/discovery/commit/a1dfa4b3))
* **deps**: bump github.com/alecthomas/kong from 0.6.0 to 0.6.1 (#68) ([de21b5f3](https://github.com/postfinance/discovery/commit/de21b5f3))
* **deps**: bump github.com/alecthomas/kong from 0.6.1 to 0.7.0 (#95) ([639f50ef](https://github.com/postfinance/discovery/commit/639f50ef))
  > Bumps [github.com/alecthomas/kong](https://github.com/alecthomas/kong) from 0.6.1 to 0.7.0.
  > - [Release notes](https://github.com/alecthomas/kong/releases)
  > - [Commits](https://github.com/alecthomas/kong/compare/v0.6.1...v0.7.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/alecthomas/kong
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/alecthomas/kong from 0.7.0 to 0.7.1 (#102) ([fe248ae8](https://github.com/postfinance/discovery/commit/fe248ae8))
* **deps**: bump github.com/coreos/go-oidc/v3 from 3.1.0 to 3.2.0 (#57) ([03f7fd11](https://github.com/postfinance/discovery/commit/03f7fd11))
* **deps**: bump github.com/coreos/go-oidc/v3 from 3.2.0 to 3.4.0 (#86) ([7cc5a460](https://github.com/postfinance/discovery/commit/7cc5a460))
* **deps**: bump github.com/golang-jwt/jwt/v4 from 4.4.1 to 4.4.2 (#70) ([8f010a90](https://github.com/postfinance/discovery/commit/8f010a90))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#60) ([5c49f0ff](https://github.com/postfinance/discovery/commit/5c49f0ff))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#61) ([28cf4c39](https://github.com/postfinance/discovery/commit/28cf4c39))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#64) ([bae94397](https://github.com/postfinance/discovery/commit/bae94397))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#74) ([4b55f76d](https://github.com/postfinance/discovery/commit/4b55f76d))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#78) ([f1b75d42](https://github.com/postfinance/discovery/commit/f1b75d42))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#93) ([abce63ea](https://github.com/postfinance/discovery/commit/abce63ea))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#96) ([23816705](https://github.com/postfinance/discovery/commit/23816705))
* **deps**: bump github.com/postfinance/flash from 0.3.0 to 0.4.0 (#92) ([c98c6fb5](https://github.com/postfinance/discovery/commit/c98c6fb5))
* **deps**: bump github.com/postfinance/single from 0.0.1 to 0.0.2 (#59) ([0a969030](https://github.com/postfinance/discovery/commit/0a969030))
* **deps**: bump github.com/prometheus/client_golang (#100) ([380c4e6b](https://github.com/postfinance/discovery/commit/380c4e6b))
* **deps**: bump github.com/prometheus/client_golang (#56) ([581e9191](https://github.com/postfinance/discovery/commit/581e9191))
* **deps**: bump github.com/prometheus/client_golang (#77) ([e494d773](https://github.com/postfinance/discovery/commit/e494d773))
* **deps**: bump github.com/prometheus/client_golang (#97) ([9f6ff9dd](https://github.com/postfinance/discovery/commit/9f6ff9dd))
* **deps**: bump github.com/stretchr/testify from 1.7.1 to 1.7.2 (#66) ([46238ff5](https://github.com/postfinance/discovery/commit/46238ff5))
* **deps**: bump github.com/stretchr/testify from 1.7.2 to 1.8.0 (#71) ([cfd238dd](https://github.com/postfinance/discovery/commit/cfd238dd))
* **deps**: bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#94) ([64f3330e](https://github.com/postfinance/discovery/commit/64f3330e))
  > Bumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.8.0 to 1.8.1.
  > - [Release notes](https://github.com/stretchr/testify/releases)
  > - [Commits](https://github.com/stretchr/testify/compare/v1.8.0...v1.8.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/stretchr/testify
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump go.uber.org/zap from 1.21.0 to 1.22.0 (#79) ([49cddd63](https://github.com/postfinance/discovery/commit/49cddd63))
* **deps**: bump golang.org/x/oauth2 from 0.1.0 to 0.2.0 (#98) ([70f29c82](https://github.com/postfinance/discovery/commit/70f29c82))
* **deps**: bump google.golang.org/grpc from 1.46.0 to 1.46.2 (#58) ([c9f26ba8](https://github.com/postfinance/discovery/commit/c9f26ba8))
* **deps**: bump google.golang.org/grpc from 1.46.2 to 1.47.0 (#63) ([d0b7f221](https://github.com/postfinance/discovery/commit/d0b7f221))
* **deps**: bump google.golang.org/grpc from 1.47.0 to 1.48.0 (#73) ([136060c4](https://github.com/postfinance/discovery/commit/136060c4))
  > Bumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.47.0 to 1.48.0.
  > - [Release notes](https://github.com/grpc/grpc-go/releases)
  > - [Commits](https://github.com/grpc/grpc-go/compare/v1.47.0...v1.48.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: google.golang.org/grpc
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump google.golang.org/grpc from 1.48.0 to 1.50.1 (#90) ([8c774d3f](https://github.com/postfinance/discovery/commit/8c774d3f))
* **deps**: bump google.golang.org/grpc from 1.50.1 to 1.51.0 (#103) ([31f5d6e3](https://github.com/postfinance/discovery/commit/31f5d6e3))
* **deps**: bump k8s.io/apimachinery from 0.24.0 to 0.24.1 (#62) ([a88497d6](https://github.com/postfinance/discovery/commit/a88497d6))
* **deps**: bump k8s.io/apimachinery from 0.24.1 to 0.24.2 (#67) ([c1c20ff9](https://github.com/postfinance/discovery/commit/c1c20ff9))
* **deps**: bump k8s.io/apimachinery from 0.24.2 to 0.24.3 (#72) ([856a1abe](https://github.com/postfinance/discovery/commit/856a1abe))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.24.2 to 0.24.3.
  > - [Release notes](https://github.com/kubernetes/apimachinery/releases)
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.24.2...v0.24.3)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/apimachinery from 0.24.3 to 0.25.3 (#91) ([34ab5b20](https://github.com/postfinance/discovery/commit/34ab5b20))
* **deps**: bump k8s.io/apimachinery from 0.25.3 to 0.25.4 (#101) ([55e5249d](https://github.com/postfinance/discovery/commit/55e5249d))
* **deps**: github.com/coreos/go-oidc/v3 3.4.0 -> 3.5.0 ([c78868c2](https://github.com/postfinance/discovery/commit/c78868c2))
* **deps**: github.com/golang-jwt/jwt/v4 4.4.2 -> 4.4.3 ([0848956a](https://github.com/postfinance/discovery/commit/0848956a))
* **deps**: github.com/grpc-ecosystem/grpc-gateway/v2 2.13.0 -> 2.15.0 ([f3a42303](https://github.com/postfinance/discovery/commit/f3a42303))
* **deps**: github.com/postfinance/flash 0.4.0 -> 0.5.0 ([38a824d2](https://github.com/postfinance/discovery/commit/38a824d2))
* **deps**: go.uber.org/zap 1.23.0 -> 1.24.0 ([bb6c1c61](https://github.com/postfinance/discovery/commit/bb6c1c61))
* **deps**: golang.org/x/oauth2 0.2.0 -> 0.4.0 ([e83a9c34](https://github.com/postfinance/discovery/commit/e83a9c34))
* **deps**: google.golang.org/genproto 0.0.0-20221027153422-115e99e71e1c -> 0.0.0-20230124163310-31e0e69b6fc2 ([5cb1f536](https://github.com/postfinance/discovery/commit/5cb1f536))
* **deps**: google.golang.org/grpc 1.51.0 -> 1.52.1 ([d926718d](https://github.com/postfinance/discovery/commit/d926718d))
* **deps**: k8s.io/apimachinery 0.25.4 -> 0.26.1 ([69692361](https://github.com/postfinance/discovery/commit/69692361))



## 0.8.3 (2023-01-25)


### Build System

* **deps**: bump github.com/alecthomas/kong from 0.5.0 to 0.6.0 (#65) ([a1dfa4b3](https://github.com/postfinance/discovery/commit/a1dfa4b3))
* **deps**: bump github.com/alecthomas/kong from 0.6.0 to 0.6.1 (#68) ([de21b5f3](https://github.com/postfinance/discovery/commit/de21b5f3))
* **deps**: bump github.com/alecthomas/kong from 0.6.1 to 0.7.0 (#95) ([639f50ef](https://github.com/postfinance/discovery/commit/639f50ef))
  > Bumps [github.com/alecthomas/kong](https://github.com/alecthomas/kong) from 0.6.1 to 0.7.0.
  > - [Release notes](https://github.com/alecthomas/kong/releases)
  > - [Commits](https://github.com/alecthomas/kong/compare/v0.6.1...v0.7.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/alecthomas/kong
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/alecthomas/kong from 0.7.0 to 0.7.1 (#102) ([fe248ae8](https://github.com/postfinance/discovery/commit/fe248ae8))
* **deps**: bump github.com/coreos/go-oidc/v3 from 3.1.0 to 3.2.0 (#57) ([03f7fd11](https://github.com/postfinance/discovery/commit/03f7fd11))
* **deps**: bump github.com/coreos/go-oidc/v3 from 3.2.0 to 3.4.0 (#86) ([7cc5a460](https://github.com/postfinance/discovery/commit/7cc5a460))
* **deps**: bump github.com/golang-jwt/jwt/v4 from 4.4.1 to 4.4.2 (#70) ([8f010a90](https://github.com/postfinance/discovery/commit/8f010a90))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#60) ([5c49f0ff](https://github.com/postfinance/discovery/commit/5c49f0ff))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#61) ([28cf4c39](https://github.com/postfinance/discovery/commit/28cf4c39))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#64) ([bae94397](https://github.com/postfinance/discovery/commit/bae94397))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#74) ([4b55f76d](https://github.com/postfinance/discovery/commit/4b55f76d))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#78) ([f1b75d42](https://github.com/postfinance/discovery/commit/f1b75d42))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#93) ([abce63ea](https://github.com/postfinance/discovery/commit/abce63ea))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#96) ([23816705](https://github.com/postfinance/discovery/commit/23816705))
* **deps**: bump github.com/postfinance/flash from 0.3.0 to 0.4.0 (#92) ([c98c6fb5](https://github.com/postfinance/discovery/commit/c98c6fb5))
* **deps**: bump github.com/postfinance/single from 0.0.1 to 0.0.2 (#59) ([0a969030](https://github.com/postfinance/discovery/commit/0a969030))
* **deps**: bump github.com/prometheus/client_golang (#100) ([380c4e6b](https://github.com/postfinance/discovery/commit/380c4e6b))
* **deps**: bump github.com/prometheus/client_golang (#56) ([581e9191](https://github.com/postfinance/discovery/commit/581e9191))
* **deps**: bump github.com/prometheus/client_golang (#77) ([e494d773](https://github.com/postfinance/discovery/commit/e494d773))
* **deps**: bump github.com/prometheus/client_golang (#97) ([9f6ff9dd](https://github.com/postfinance/discovery/commit/9f6ff9dd))
* **deps**: bump github.com/stretchr/testify from 1.7.1 to 1.7.2 (#66) ([46238ff5](https://github.com/postfinance/discovery/commit/46238ff5))
* **deps**: bump github.com/stretchr/testify from 1.7.2 to 1.8.0 (#71) ([cfd238dd](https://github.com/postfinance/discovery/commit/cfd238dd))
* **deps**: bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#94) ([64f3330e](https://github.com/postfinance/discovery/commit/64f3330e))
  > Bumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.8.0 to 1.8.1.
  > - [Release notes](https://github.com/stretchr/testify/releases)
  > - [Commits](https://github.com/stretchr/testify/compare/v1.8.0...v1.8.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/stretchr/testify
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump go.uber.org/zap from 1.21.0 to 1.22.0 (#79) ([49cddd63](https://github.com/postfinance/discovery/commit/49cddd63))
* **deps**: bump golang.org/x/oauth2 from 0.1.0 to 0.2.0 (#98) ([70f29c82](https://github.com/postfinance/discovery/commit/70f29c82))
* **deps**: bump google.golang.org/grpc from 1.46.0 to 1.46.2 (#58) ([c9f26ba8](https://github.com/postfinance/discovery/commit/c9f26ba8))
* **deps**: bump google.golang.org/grpc from 1.46.2 to 1.47.0 (#63) ([d0b7f221](https://github.com/postfinance/discovery/commit/d0b7f221))
* **deps**: bump google.golang.org/grpc from 1.47.0 to 1.48.0 (#73) ([136060c4](https://github.com/postfinance/discovery/commit/136060c4))
  > Bumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.47.0 to 1.48.0.
  > - [Release notes](https://github.com/grpc/grpc-go/releases)
  > - [Commits](https://github.com/grpc/grpc-go/compare/v1.47.0...v1.48.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: google.golang.org/grpc
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump google.golang.org/grpc from 1.48.0 to 1.50.1 (#90) ([8c774d3f](https://github.com/postfinance/discovery/commit/8c774d3f))
* **deps**: bump google.golang.org/grpc from 1.50.1 to 1.51.0 (#103) ([31f5d6e3](https://github.com/postfinance/discovery/commit/31f5d6e3))
* **deps**: bump k8s.io/apimachinery from 0.24.0 to 0.24.1 (#62) ([a88497d6](https://github.com/postfinance/discovery/commit/a88497d6))
* **deps**: bump k8s.io/apimachinery from 0.24.1 to 0.24.2 (#67) ([c1c20ff9](https://github.com/postfinance/discovery/commit/c1c20ff9))
* **deps**: bump k8s.io/apimachinery from 0.24.2 to 0.24.3 (#72) ([856a1abe](https://github.com/postfinance/discovery/commit/856a1abe))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.24.2 to 0.24.3.
  > - [Release notes](https://github.com/kubernetes/apimachinery/releases)
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.24.2...v0.24.3)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/apimachinery from 0.24.3 to 0.25.3 (#91) ([34ab5b20](https://github.com/postfinance/discovery/commit/34ab5b20))
* **deps**: bump k8s.io/apimachinery from 0.25.3 to 0.25.4 (#101) ([55e5249d](https://github.com/postfinance/discovery/commit/55e5249d))
* **deps**: github.com/coreos/go-oidc/v3 3.4.0 -> 3.5.0 ([c78868c2](https://github.com/postfinance/discovery/commit/c78868c2))
* **deps**: github.com/golang-jwt/jwt/v4 4.4.2 -> 4.4.3 ([0848956a](https://github.com/postfinance/discovery/commit/0848956a))
* **deps**: github.com/grpc-ecosystem/grpc-gateway/v2 2.13.0 -> 2.15.0 ([f3a42303](https://github.com/postfinance/discovery/commit/f3a42303))
* **deps**: github.com/postfinance/flash 0.4.0 -> 0.5.0 ([38a824d2](https://github.com/postfinance/discovery/commit/38a824d2))
* **deps**: go.uber.org/zap 1.23.0 -> 1.24.0 ([bb6c1c61](https://github.com/postfinance/discovery/commit/bb6c1c61))
* **deps**: golang.org/x/oauth2 0.2.0 -> 0.4.0 ([e83a9c34](https://github.com/postfinance/discovery/commit/e83a9c34))
* **deps**: google.golang.org/genproto 0.0.0-20221027153422-115e99e71e1c -> 0.0.0-20230124163310-31e0e69b6fc2 ([5cb1f536](https://github.com/postfinance/discovery/commit/5cb1f536))
* **deps**: google.golang.org/grpc 1.51.0 -> 1.52.1 ([d926718d](https://github.com/postfinance/discovery/commit/d926718d))
* **deps**: k8s.io/apimachinery 0.25.4 -> 0.26.1 ([69692361](https://github.com/postfinance/discovery/commit/69692361))



## 0.8.2 (2022-05-09)


### Bug Fixes

* **discoveryd**: return the correct services for a namespace ([f74115f6](https://github.com/postfinance/discovery/commit/f74115f6))
  > Prior ot this change, the server returned too many services for a namespace
  > when two or more namespaces started with the same name. For example with
  > the namespaces `default` and `default-blackbox` the service returned
  > services for both namespaces when selecting only `default` namespace.


### Build System

* **deps**: bump github.com/alecthomas/kong from 0.3.0 to 0.4.0 (#36) ([2b5b2041](https://github.com/postfinance/discovery/commit/2b5b2041))
  > Bumps [github.com/alecthomas/kong](https://github.com/alecthomas/kong) from 0.3.0 to 0.4.0.
  > - [Release notes](https://github.com/alecthomas/kong/releases)
  > - [Commits](https://github.com/alecthomas/kong/compare/v0.3.0...v0.4.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/alecthomas/kong
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/alecthomas/kong from 0.4.0 to 0.4.1 (#40) ([81bac235](https://github.com/postfinance/discovery/commit/81bac235))
* **deps**: bump github.com/alecthomas/kong from 0.4.1 to 0.5.0 (#43) ([36e5d4da](https://github.com/postfinance/discovery/commit/36e5d4da))
* **deps**: bump github.com/golang-jwt/jwt/v4 from 4.2.0 to 4.3.0 (#38) ([508d0d6a](https://github.com/postfinance/discovery/commit/508d0d6a))
* **deps**: bump github.com/golang-jwt/jwt/v4 from 4.3.0 to 4.4.0 (#46) ([2820e9c7](https://github.com/postfinance/discovery/commit/2820e9c7))
  > Bumps [github.com/golang-jwt/jwt/v4](https://github.com/golang-jwt/jwt) from 4.3.0 to 4.4.0.
  > - [Release notes](https://github.com/golang-jwt/jwt/releases)
  > - [Changelog](https://github.com/golang-jwt/jwt/blob/main/VERSION_HISTORY.md)
  > - [Commits](https://github.com/golang-jwt/jwt/compare/v4.3.0...v4.4.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/golang-jwt/jwt/v4
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/golang-jwt/jwt/v4 from 4.4.0 to 4.4.1 (#50) ([b6fc68d8](https://github.com/postfinance/discovery/commit/b6fc68d8))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 ([2d135b11](https://github.com/postfinance/discovery/commit/2d135b11))
  > Bumps [github.com/grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway) from 2.7.2 to 2.7.3.
  > - [Release notes](https://github.com/grpc-ecosystem/grpc-gateway/releases)
  > - [Changelog](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/.goreleaser.yml)
  > - [Commits](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.7.2...v2.7.3)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/grpc-ecosystem/grpc-gateway/v2
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#29) ([b06ff6d5](https://github.com/postfinance/discovery/commit/b06ff6d5))
  > Bumps [github.com/grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway) from 2.6.0 to 2.7.2.
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#42) ([d12ead46](https://github.com/postfinance/discovery/commit/d12ead46))
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#45) ([9b0a054d](https://github.com/postfinance/discovery/commit/9b0a054d))
  > Bumps [github.com/grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway) from 2.8.0 to 2.9.0.
  > - [Release notes](https://github.com/grpc-ecosystem/grpc-gateway/releases)
  > - [Changelog](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/.goreleaser.yml)
  > - [Commits](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.8.0...v2.9.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/grpc-ecosystem/grpc-gateway/v2
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/grpc-ecosystem/grpc-gateway/v2 (#52) ([7ff8c5c5](https://github.com/postfinance/discovery/commit/7ff8c5c5))
* **deps**: bump github.com/postfinance/flash from 0.2.0 to 0.3.0 (#47) ([112e85dc](https://github.com/postfinance/discovery/commit/112e85dc))
  > Bumps [github.com/postfinance/flash](https://github.com/postfinance/flash) from 0.2.0 to 0.3.0.
  > - [Release notes](https://github.com/postfinance/flash/releases)
  > - [Changelog](https://github.com/postfinance/flash/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/postfinance/flash/compare/v0.2.0...v0.3.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/postfinance/flash
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/prometheus/client_golang (#32) ([a67bde86](https://github.com/postfinance/discovery/commit/a67bde86))
  > Bumps [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang) from 1.11.0 to 1.12.0.
  > - [Release notes](https://github.com/prometheus/client_golang/releases)
  > - [Changelog](https://github.com/prometheus/client_golang/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/prometheus/client_golang/compare/v1.11.0...v1.12.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/prometheus/client_golang
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump github.com/prometheus/client_golang (#34) ([c0620988](https://github.com/postfinance/discovery/commit/c0620988))
  > Bumps [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang) from 1.12.0 to 1.12.1.
  > - [Release notes](https://github.com/prometheus/client_golang/releases)
  > - [Changelog](https://github.com/prometheus/client_golang/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/prometheus/client_golang/compare/v1.12.0...v1.12.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/prometheus/client_golang
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (#48) ([5f524e47](https://github.com/postfinance/discovery/commit/5f524e47))
  > Bumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.7.0 to 1.7.1.
  > - [Release notes](https://github.com/stretchr/testify/releases)
  > - [Commits](https://github.com/stretchr/testify/compare/v1.7.0...v1.7.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/stretchr/testify
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump go.uber.org/zap from 1.19.1 to 1.20.0 (#30) ([7cf9e380](https://github.com/postfinance/discovery/commit/7cf9e380))
* **deps**: bump go.uber.org/zap from 1.20.0 to 1.21.0 (#39) ([dce57ad7](https://github.com/postfinance/discovery/commit/dce57ad7))
* **deps**: bump google.golang.org/grpc from 1.43.0 to 1.44.0 (#37) ([ce89c171](https://github.com/postfinance/discovery/commit/ce89c171))
  > Bumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.43.0 to 1.44.0.
  > - [Release notes](https://github.com/grpc/grpc-go/releases)
  > - [Commits](https://github.com/grpc/grpc-go/compare/v1.43.0...v1.44.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: google.golang.org/grpc
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump google.golang.org/grpc from 1.44.0 to 1.45.0 (#44) ([38e28810](https://github.com/postfinance/discovery/commit/38e28810))
* **deps**: bump google.golang.org/grpc from 1.45.0 to 1.46.0 (#54) ([2a018ed2](https://github.com/postfinance/discovery/commit/2a018ed2))
* **deps**: bump google.golang.org/protobuf from 1.27.1 to 1.28.0 (#51) ([017da592](https://github.com/postfinance/discovery/commit/017da592))
* **deps**: bump k8s.io/apimachinery from 0.23.1 to 0.23.2 (#33) ([3caddcda](https://github.com/postfinance/discovery/commit/3caddcda))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.23.1 to 0.23.2.
  > - [Release notes](https://github.com/kubernetes/apimachinery/releases)
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.23.1...v0.23.2)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/apimachinery from 0.23.1 to 0.23.3 (#35) ([ff6425bb](https://github.com/postfinance/discovery/commit/ff6425bb))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.23.1 to 0.23.3.
  > - [Release notes](https://github.com/kubernetes/apimachinery/releases)
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.23.1...v0.23.3)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/apimachinery from 0.23.3 to 0.23.4 (#41) ([f247efcc](https://github.com/postfinance/discovery/commit/f247efcc))
* **deps**: bump k8s.io/apimachinery from 0.23.4 to 0.23.5 (#49) ([d24f049d](https://github.com/postfinance/discovery/commit/d24f049d))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.23.4 to 0.23.5.
  > - [Release notes](https://github.com/kubernetes/apimachinery/releases)
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.23.4...v0.23.5)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/apimachinery from 0.23.5 to 0.23.6 (#53) ([9829cf75](https://github.com/postfinance/discovery/commit/9829cf75))
* **deps**: bump k8s.io/apimachinery from 0.23.6 to 0.24.0 (#55) ([233fefe2](https://github.com/postfinance/discovery/commit/233fefe2))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.23.6 to 0.24.0.
  > - [Release notes](https://github.com/kubernetes/apimachinery/releases)
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.23.6...v0.24.0)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...



## 0.8.1 (2022-01-10)


### Bug Fixes

* **client**: use better help for filter by label selector flag ([b987d41e](https://github.com/postfinance/discovery/commit/b987d41e))
* **discovery**: help displays positional argument as required ([#13](https://github.com/postfinance/discovery/issues/13), [7653c937](https://github.com/postfinance/discovery/commit/7653c937))


### Dependencies

* **enproto**: 0.0.0-20210927142257-433400c27d05 -> 0.0.0-20211007155348-82e027067bd4 ([d67b5bc2](https://github.com/postfinance/discovery/commit/d67b5bc2))
* **oauth2**: 0.0.0-20210819190943-2bc19b11175f -> 0.0.0-20211005180243-6b3c2da341f1 ([1971b8e1](https://github.com/postfinance/discovery/commit/1971b8e1))



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
