// Create a bootkube-worker profile
resource "matchbox_profile" "bootkube-worker" {
  name = "bootkube-worker"
  config = "${matchbox_config.bootkube-worker.name}"
}

// Create bootkube-worker CoreOS config resource
resource "matchbox_config" "bootkube-worker" {
  name = "bootkube-worker.yaml.tmpl"
  contents = "${file("./cl/bootkube-worker.yaml.tmpl")}"
}

// Create matcher groups for worker nodes

resource "matchbox_group" "node2" {
  name = "node2"
  profile = "${matchbox_profile.bootkube-worker.name}"
  selector {
    mac = "52:54:00:b2:2f:86"
    os = "installed"
  }
  metadata {
    domain_name = "node2.example.com"
    etcd_endpoints = "node1.example.com:2380"
    k8s_dns_service_ip = "${var.k8s_dns_service_ip}"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

resource "matchbox_group" "node3" {
  name = "node3"
  profile = "${matchbox_profile.bootkube-worker.name}"
  selector {
    mac = "52:54:00:c3:61:77"
    os = "installed"
  }
  metadata {
    domain_name = "node3.example.com"
    etcd_endpoints = "node1.example.com:2380"
    k8s_dns_service_ip = "${var.k8s_dns_service_ip}"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

