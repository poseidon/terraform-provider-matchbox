// Configure the matchbox provider
provider "matchbox" {
  endpoint = "matchbox.example.com:8081"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key = "${file("~/.matchbox/client.key")}"
  ca         = "${file("~/.matchbox/ca.crt")}"
}


// Create an etcd3 profile
resource "matchbox_profile" "etcd3" {
  name = "etcd3"
  /*ignition = "${matchbox_ct.etcd3}"
  boot {
    kernel = "/assets/coreos/1235.9.0/coreos_production_pxe.vmlinuz"
    initrd = "/assets/coreos/1235.9.0/coreos_production_pxe_image.cpio.gz"
    args = [
      "coreos.config.url=http://matchbox.foo:8080/ignition?uuid=${uuid}&mac=${mac:hexhyp}",
      "coreos.first_boot=yes",
      "console=tty0",
      "console=ttyS0",
      "coreos.autologin"
    ]
  }*/
}


// Create etcd3 CoreOS config resource
resource "matchbox_config" "etcd3" {
  name = "etcd3.yaml.tmpl"
  //  contents = "${file("./etcd3.yaml.tmpl")}"
}


// Create a matcher group
resource "matchbox_group" "etcd3" {
  name = "node1"
  /*profile = "${matchbox_profile.etcd3.name}"
  selector {
    mac = "52:54:00:a1:9c:ae"
    os = "installed"
  }
  metadata {
    domain_name = "node1.example.com"
    etcd_name = "node1"
    etcd_initial_cluster = "node1=http://node1.example.com:2380"
  }*/
}

