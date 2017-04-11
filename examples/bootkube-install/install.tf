// Create a CoreOS install-reboot profile
resource "matchbox_profile" "install-reboot" {
  name = "install-reboot"
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
  container_linux_config = "${file("./cl/install-reboot.yaml.tmpl")}"
}

resource "matchbox_group" "default" {
  name = "default"
  profile = "${matchbox_profile.install-reboot.name}"
  metadata {
    coreos_channel = "stable"
    coreos_version = "1235.9.0"
    ignition_endpoint = "http://matchbox.example.com:8080/ignition"
    baseurl = "http://matchbox.example.com:8080/assets/coreos"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
