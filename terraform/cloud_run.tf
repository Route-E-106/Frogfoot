resource "google_cloud_run_service" "frogfoot_server" {
  autogenerate_revision_name = false
  location                   = "europe-north1"
  name                       = "frogfoot"

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"     = "1"
        "run.googleapis.com/client-name"       = "cloud-console"
        "run.googleapis.com/startup-cpu-boost" = "false"
      }
      labels = {
        "run.googleapis.com/startupProbeType" = "Default"
      }
    }
    spec {
      container_concurrency = 80
      node_selector         = {}
      timeout_seconds       = 300

      containers {
        args    = []
        command = []
        image   = "registry-1.docker.io/fulcrum29/frogfoot:latest"
        name    = "frogfoot-1"

        ports {
          container_port = 8080
          name           = "http1"
        }

        resources {
          limits = {
            "cpu"    = "1000m"
            "memory" = "256Mi"
          }
          requests = {}
        }

        startup_probe {
          failure_threshold     = 1
          initial_delay_seconds = 0
          period_seconds        = 240
          timeout_seconds       = 240

          tcp_socket {
            port = 8080
          }
        }

        volume_mounts {
          mount_path = "/volume"
          name       = "memory"
        }
      }

      volumes {
        name = "memory"

        empty_dir {
          medium     = "Memory"
          size_limit = "100"
        }
      }
    }
  }

  traffic {
    latest_revision = true
    percent         = 100
  }
}
