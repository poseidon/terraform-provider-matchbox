# Matchbox Provider

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
