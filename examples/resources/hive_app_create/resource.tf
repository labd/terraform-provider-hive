resource "hive_app_create" "example" {
  name      = "example-service"
  version   = "1.0.0"
  documents = file("persisted-documents.json")
}
