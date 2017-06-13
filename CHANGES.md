# terraform-provider-matchbox

Notable changes between releases.

## Latest

## v0.2.0 (2017-06-13)

* Add Profile `generic_config` field to write generic/experimental config templates to Matchbox
* Add Profile `raw_ignition` field to write raw Ignition to Matchbox. Note that providing a `container_linux_config` is preferred.

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
