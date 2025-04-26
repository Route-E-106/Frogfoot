resource "google_compute_instance" "frogfoot_server" {
  name         = "frogfoot-server"
  machine_type = "e2-small"
  zone         = "europe-north1-a"
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-minimal-2504-plucky-amd64-v20250415"
      type  = "pd-standard"
      size  = 10

    }
    auto_delete = true
    device_name = "frogfoot-server"

  }
  network_interface {
    network    = "default"
    stack_type = "IPV4_ONLY"
    access_config {
      network_tier = "STANDARD"
    }

  }
}

