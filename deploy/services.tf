resource "docker_image" "gateway_image" {
  name = "gateway_image"
  build {
    path = "../services/gateway"
    tag  = ["collectable/gateway:latest"]
  }
}

resource "helm_release" "gateway_service_release" {
  name      = "gateway"
  namespace = "collectable"

  chart = "./charts/gateway"

  timeout    = 120
  depends_on = [docker_image.gateway_image, helm_release.iam_service_release]
}

resource "docker_image" "iam_image" {
  name = "iam_image"
  build {
    path = "../services/iam"
    tag  = ["collectable/iam:latest"]
  }
}

resource "helm_release" "iam_service_release" {
  name      = "iam"
  namespace = "collectable"

  chart = "./charts/iam"

  timeout    = 120
  depends_on = [docker_image.iam_image]
}