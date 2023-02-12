# Profile Resource

A Profile defines network boot and declarative provisioning configurations.

```tf
variable "os_stream" {
  type        = string
  description = "Fedora CoreOS release stream (e.g. testing, stable)"
  default     = "stable"
}

variable "os_version" {
  type        = string
  description = "Fedora CoreOS version to PXE and install (e.g. 32.20200715.3.0)"
}

locals {
  kernel = "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-kernel-x86_64"
  initrd = "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-initramfs.x86_64.img"
}
```

```tf
resource "matchbox_profile" "worker" {
  name = "worker"
  kernel = local.kernel
  initrd = [
    local.initrd
  ]
  args = [
    "ip=dhcp",
    "rd.neednet=1",
    "initrd=fedora-coreos-${var.os_version}-live-initramfs.x86_64.img",
    "coreos.inst.image_url=https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-metal.x86_64.raw.xz",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "coreos.inst.install_dev=sda",
    "console=tty0",
    "console=ttyS0",
  ]

  raw_ignition = data.ct_config.worker.rendered
}

// Transpile Fedora CoreOS config to Ignition
data "ct_config" "worker" {
  content      = file("worker.yaml")
  strict       = true
}
```

## Argument Reference

* `name` - Unqiue name for the machine matcher
* `kernel` - URL of the kernel image to boot
* `initrd` - List of URLs to init RAM filesystems
* `args` - List of kernel arguments
* `raw_ignition` - Fedora CoreOS or Flatcar Linux Ignition content (see [terraform-provider-ct](https://github.com/poseidon/terraform-provider-ct))
* `generic_config` - Generic configuration
* `container_linux_config` -  CoreOS Container Linux Config (CLC) (for backwards compatibility)
