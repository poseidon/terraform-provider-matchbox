// CoreOS Install Profile
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
  container_linux_config = "${file("${path.module}/cl/coreos-install.yaml.tmpl")}"
}

// etcd3 profile
resource "matchbox_profile" "etcd3" {
  name = "etcd3"
  container_linux_config = "${file("${path.module}/cl/etcd3.yaml.tmpl")}"
}

// etcd3 Gateway profile
resource "matchbox_profile" "etcd3-gateway" {
  name = "etcd3-gateway"
  container_linux_config = "${file("${path.module}/cl/etcd3-gateway.yaml.tmpl")}"
}

// Self-hosted Kubernetes (bootkube) Controller profile
resource "matchbox_profile" "bootkube-controller" {
  name = "bootkube-controller"
  container_linux_config = "${file("${path.module}/cl/bootkube-controller.yaml.tmpl")}"
}

// Self-hosted Kubernetes (bootkube) Worker profile
resource "matchbox_profile" "bootkube-worker" {
  name = "bootkube-worker"
  container_linux_config = "${file("${path.module}/cl/bootkube-worker.yaml.tmpl")}"
}
