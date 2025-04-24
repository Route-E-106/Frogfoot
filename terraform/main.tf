resource "google_cloud_run_service" "frogfoot_server" {
  name     = "frogfoot"
  location = "europe-north1"

}

import {
  id = "europe-north1/stable-terminus-457813-p5/frogfoot"
  to = google_cloud_run_service.frogfoot_server
}
