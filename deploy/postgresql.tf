resource "helm_release" "postgresql" {
  name      = "postgresql"
  namespace = "collectable"

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "postgresql"

  set {
    name  = "auth.username"
    value = "collectable_admin"
  }

  set {
    name  = "auth.password"
    value = "collectable_admin_password"
  }

  set {
    name  = "primary.service.type"
    value = "LoadBalancer"
  }
}

resource "postgresql_database" "iam_database" {
  name       = "iam"
  depends_on = [helm_release.postgresql]
}