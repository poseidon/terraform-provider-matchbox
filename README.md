# Matchbox Terraform Provider

The Matchbox provider is used to interact with the [matchbox](https://github.com/coreos/matchbox) API. Matchbox matches bare-metal machines by labels (e.g. MAC address) to Profiles with iPXE configs and Container Linux configs to provision clusters.

## Status

Warning, this project is pre-alpha. Breaking changes are expected. Matchbox latest is required.

## Usage

Install [Terraform](https://www.terraform.io/downloads.html) v0.9.2. Add the `terraform-provider-matchbox` plugin binary somewhere on your filesystem.

```sh
go get -u github.com/coreos/terraform-provider-matchbox
```

Register the plugin in `~/.terraformrc`.

```hcl
providers {
  matchbox = "/path/to/terraform-provider-matchbox"
}
```

On-premise, [setup](https://coreos.com/matchbox/docs/latest/network-setup.html) a PXE network boot environment. [Install matchbox](https://coreos.com/matchbox/docs/latest/deployment.html) on a provisioner node or Kubernetes cluster. Be sure to enable the gRPC API and follow the instructions to generate TLS credentials.

### Examples

```tf
// Configure the matchbox provider
provider "matchbox" {
  endpoint = "matchbox.example.com:8081"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key = "${file("~/.matchbox/client.key")}"
  ca         = "${file("~/.matchbox/ca.crt")}"
}

// Create a CoreOS-install profile
resource "matchbox_profile" "coreos-install" {
  name = "coreos-install"
  kernel = "/assets/coreos/1235.9.0/coreos_production_pxe.vmlinuz"
  initrd = [
    "/assets/coreos/1235.9.0/coreos_production_pxe_image.cpio.gz"
  ]
  args = [
    "coreos.config.url=http://matchbox.example.com:8080/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "coreos.first_boot=yes",
    "console=tty0",
    "console=ttyS0",
    "coreos.autologin"
  ]
  container_linux_config = "${file("./cl/coreos-install.yaml.tmpl")}"
}

// Match a bare-metal machine
resource "matchbox_group" "node1" {
  name = "node1"
  profile = "${matchbox_profile.coreos-install.name}"
  selector {
    mac = "52:54:00:a1:9c:ae"
  }
  metadata {
    coreos_channel = "stable"
    coreos_version = "1235.9.0"
    ignition_endpoint = "http://matchbox.example.com:8080/ignition"
    baseurl = "http://matchbox.example.com:8080/assets/coreos"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
```

See [examples](https://github.com/coreos/matchbox/tree/master/examples/terraform) for Terraform configs which PXE boot, install CoreOS, and provision entire clusters.

## Development

### Binary

To develop the plguin locally, compile and install the executable with Go 1.8.

    make build
    make test

### Vendor

Add or update dependencies in glide.yaml and vendor. The glide and glide-vc tools vendor and prune dependencies.

    make vendor
