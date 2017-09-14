# Matchbox Terraform Provider

The Matchbox provider is used to interact with the [matchbox](https://github.com/coreos/matchbox) API. Matchbox matches bare-metal machines by labels (e.g. MAC address) to Profiles with iPXE configs, Container Linux configs, or generic free-form configs in order to provision clusters.

## Usage

A Matchbox v0.6+ [installation](https://coreos.com/matchbox/docs/latest/deployment.html) is required. Matchbox v0.7+ is required to use the `generic_config` field.

Install [Terraform](https://www.terraform.io/downloads.html) v0.9+ Add the `terraform-provider-matchbox` plugin binary somewhere on your filesystem.

```sh
# dev
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
  endpoint = "${var.matchbox_rpc_endpoint}"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key = "${file("~/.matchbox/client.key")}"
  ca         = "${file("~/.matchbox/ca.crt")}"
}

// Create a Container Linux install profile
resource "matchbox_profile" "container-linux-install" {
  name = "container-linux-install"
  kernel = "/assets/coreos/${var.container_linux_version}/coreos_production_pxe.vmlinuz"
  initrd = [
    "/assets/coreos/${var.container_linux_version}/coreos_production_pxe_image.cpio.gz"
  ]
  args = [
    "coreos.config.url=http://${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "coreos.first_boot=yes",
    "console=tty0",
    "console=ttyS0",
    "coreos.autologin"
  ]
  container_linux_config = "${file("./cl/coreos-install.yaml.tmpl")}"
  generic_config = "${file("./example.ks")}"
}

// Match a bare-metal machine
resource "matchbox_group" "node1" {
  name = "node1"
  profile = "${matchbox_profile.container-linux-install.name}"
  selector {
    mac = "52:54:00:a1:9c:ae"
  }
  metadata {
    custom_variable = "machine_specific_value_here"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
```

See [examples](https://github.com/coreos/matchbox/tree/master/examples/terraform) for Terraform configs which PXE boot, install CoreOS, and provision entire clusters.

## Development

### Binary

To develop the plugin locally, compile and install the executable with Go 1.8.

    make build
    make test

### Vendor

Add or update dependencies in `glide.yaml` and vendor. The [glide](https://github.com/Masterminds/glide) and [glide-vc](https://github.com/sgotti/glide-vc) tools vendor and prune dependencies.

    make vendor
