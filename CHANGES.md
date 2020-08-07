# terraform-provider-matchbox

Notable changes between releases.

## Latest

## v0.4.1

* Fix zip archive artifacts for Darwin and Windows ([#53](https://github.com/poseidon/terraform-provider-matchbox/pull/53))

## v0.4.0

* Migrate to the Terraform Plugin SDK ([#49](https://github.com/poseidon/terraform-provider-matchbox/pull/49))
* Add Linux ARM64 release artifacts
* Add zip archive format with signed checksum

## v0.3.0

* Add compatibility with Terraform v0.12. Retain v0.11 compatibility ([#42](https://github.com/poseidon/terraform-provider-matchbox/pull/42))

## v0.2.3

* Document usage with the Terraform [3rd-party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) directory ([#39](https://github.com/poseidon/terraform-provider-matchbox/pull/39))
* Use Go v1.11.5 for pre-compiled binaries

## v0.2.2

* Improve client endpoint validation ([#23](https://github.com/poseidon/terraform-provider-matchbox/pull/23))
  * Provide better errors if endpoint includes a scheme or is missing a port

## v0.2.1

* Statically link the `terraform-provider-matchbox` binaries

## v0.2.0

* Add Profile `generic_config` field to write generic/experimental config templates to Matchbox
* Add Profile `raw_ignition` field to write raw Ignition to Matchbox. Note that providing a `container_linux_config` is preferred.

## v0.1.1

Fix darwin release, which was compiled for Linux.

* Fix Makefile cross-compilation

## v0.1.0

Initial release of the Matchbox Terraform Provider Plugin

* Configure a Provider with a matchbox TLS client cert/key
* Create matchbox machine Profile resources with Container Linux Configs
* Create matchbox matcher Group resources to match bare-metal machines to profiles
* Requires matchbox v0.6.0 installation
* See examples to provision etcd3 or self-hosted Kubernetes clusters
