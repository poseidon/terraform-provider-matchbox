// Create a bootkube-controller profile
resource "matchbox_profile" "bootkube-controller" {
  name = "bootkube-controller"
  container_linux_config = "${matchbox_config.bootkube-controller.name}"
}

// Create bootkube-controller Container Linux Config resource
resource "matchbox_config" "bootkube-controller" {
  name = "bootkube-controller.yaml.tmpl"
  contents = "${file("./cl/bootkube-controller.yaml.tmpl")}"
}

// Create a matcher group

resource "matchbox_group" "node1" {
  name = "node1"
  profile = "${matchbox_profile.bootkube-controller.name}"
  selector {
    mac = "52:54:00:a1:9c:ae"
    os = "installed"
  }
  metadata {
    domain_name = "node1.example.com"
    etcd_name = "node1"
    etcd_initial_cluster = "node1=http://node1.example.com:2380"
    k8s_dns_service_ip = "${var.k8s_dns_service_ip}"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
