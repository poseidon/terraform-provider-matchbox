# terraform-provider-matchbox [![Build Status](https://github.com/poseidon/terraform-provider-matchbox/workflows/test/badge.svg)](https://github.com/poseidon/terraform-provider-matchbox/actions?query=workflow%3Atest+branch%3Amaster)

`terraform-provider-matchbox` allows defining [Matchbox](https://github.com/poseidon/matchbox) Profiles and Groups in Terraform. Matchbox matches machines, by label (e.g. MAC address), to Profiles with iPXE configs, Ignition configs, or generic free-form configs to provision clusters. Resources are created via the client certificate authenticated Matchbox API.

## Usage

[Setup](https://matchbox.psdn.io/network-setup/) a PXE network boot environment and [deploy](https://matchbox.psdn.io/deployment/) a Matchbox instance. Be sure to enable the gRPC API and follow the instructions to generate TLS credentials.

Configure the Matchbox provider with the Matchbox API endpoint and client certificate (e.g. `providers.tf`).

```tf
provider "matchbox" {
  endpoint    = "matchbox.example.com:8081"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key  = "${file("~/.matchbox/client.key")}"
  ca          = "${file("~/.matchbox/ca.crt")}"
}

terraform {
  required_providers {
    matchbox = {
      source = "poseidon/matchbox"
      version = "0.5.1"
    }
  }
}
```

Define a Matchbox Profile or Group resource in Terraform.

```tf
// Fedora CoreOS profile
resource "matchbox_profile" "fedora-coreos-install" {
  name  = "worker"
  kernel = "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-kernel-x86_64"

  initrd = [
    "--name main https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-initramfs.x86_64.img"
  ]

  args = [
    "initrd=main",
    "coreos.live.rootfs_url=https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-rootfs.x86_64.img",
    "coreos.inst.install_dev=/dev/sda",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}"
  ]

  raw_ignition = data.ct_config.worker.rendered
}

data "ct_config" "worker" {
  content = templatefile("fcc/fedora-coreos.yaml", {
    ssh_authorized_key = var.ssh_authorized_key
  })
  strict = true
}

// Default matcher group for machines
resource "matchbox_group" "default" {
  name    = "default"
  profile = matchbox_profile.fedora-coreos-install.name
  selector = {}
  metadata = {}
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

See [examples](https://github.com/poseidon/matchbox/tree/master/examples/terraform) for Terraform configs which PXE boot, install CoreOS, and provision entire clusters.

## Requirements

* Terraform v0.13+ [installed](https://www.terraform.io/downloads.html)
* Matchbox v0.8+ [installed](https://matchbox.psdn.io/deployment/)
* Matchbox credentials `client.crt`, `client.key`, `ca.crt`

## Development

### Binary

To develop the provider plugin locally, build an executable with Go 1.17+.

```
make
```
