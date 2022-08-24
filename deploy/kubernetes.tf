resource "kubernetes_namespace" "collectable" {
  metadata {
    name = "collectable"
  }
}