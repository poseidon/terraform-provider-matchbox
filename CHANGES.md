# terraform-provider-matchbox

Notable changes between releases.

## Latest

* Add Profile generic_config support
* Add Profile raw_ignition support

## v0.1.1 (2017-05-15)

Fix darwin release, which was compiled for Linux.

* Fix Makefile cross-compilation

## v0.1.0 (2017-04-07)

Initial release of the Matchbox Terraform Provider Plugin

* Configure a Provider with a matchbox TLS client cert/key
* Create matchbox machine Profile resources with Container Linux Configs
* Create matchbox matcher Group resources to match bare-metal machines to profiles
* Requires matchbox v0.6.0 installation
* See examples to provision etcd3 or self-hosted Kubernetes clusters
