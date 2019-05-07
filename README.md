# terraform-provider-matchbox

`terraform-provider-matchbox` allows defining CoreOS [Matchbox](https://github.com/poseidon/matchbox) Profiles and Groups in Terraform. Matchbox matches machines, by label (e.g. MAC address), to Profiles with iPXE configs, Container Linux configs, or generic free-form configs to provision clusters. Resources are created via the client certificate authenticated Matchbox API.

## Requirements

* Terraform v0.11+ [installed](https://www.terraform.io/downloads.html)
* Matchbox v0.6+ [installed](https://coreos.com/matchbox/docs/latest/deployment.html) (v0.7+ to use the `generic_config` field)
* Matchbox credentials `client.crt`, `client.key`, `ca.crt`

## Install

Add the `terraform-provider-matchbox` plugin binary for your system to the Terraform 3rd-party [plugin directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) `~/.terraform.d/plugins`.

```sh
VERSION=v0.2.3
wget https://github.com/poseidon/terraform-provider-matchbox/releases/download/$VERSION/terraform-provider-matchbox-$VERSION-linux-amd64.tar.gz
tar xzf terraform-provider-matchbox-$VERSION-linux-amd64.tar.gz
mv terraform-provider-matchbox-$VERSION-linux-amd64/terraform-provider-matchbox ~/.terraform.d/plugins/terraform-provider-matchbox_$VERSION
```

Terraform plugin binary names are versioned to allow for migrations of managed infrastructure.

```
$ tree ~/.terraform.d/
/home/user/.terraform.d/
└── plugins
    ├── terraform-provider-matchbox_v0.2.2
    └── terraform-provider-matchbox_v0.2.3
```

## Usage

[Setup](https://coreos.com/matchbox/docs/latest/network-setup.html) a PXE network boot environment and [deploy](https://coreos.com/matchbox/docs/latest/deployment.html) a Matchbox instance. Be sure to enable the gRPC API and follow the instructions to generate TLS credentials.

Configure the Matchbox provider to use your Matchbox API endpoint and client certificate in a `providers.tf` file.

```tf
provider "matchbox" {
  version = "0.2.3"
  endpoint    = "matchbox.example.com:8081"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key  = "${file("~/.matchbox/client.key")}"
  ca          = "${file("~/.matchbox/ca.crt")}"
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

Define a Matchbox Profile or Group resource in Terraform.

```tf
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
  selector = {
    mac = "52:54:00:a1:9c:ae"
  }
  metadata = {
    custom_variable = "machine_specific_value_here"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
```

See [examples](https://github.com/poseidon/matchbox/tree/master/examples/terraform) for Terraform configs which PXE boot, install CoreOS, and provision entire clusters.

## Development

### Binary

To develop the provider plugin locally, build an executable with Go 1.11+.

```
make
```

### Vendor

Add or update dependencies in `go.mod` and vendor.

```
make update
make vendor
```

