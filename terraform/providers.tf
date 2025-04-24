terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.31.0"
    }
  }
  required_version = "1.9.8"
}

provider "google" {
  project = "stable-terminus-457813-p5"
  region  = "europe-north1"
  zone    = "europe-north1-a"
}
