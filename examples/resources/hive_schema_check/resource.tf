resource "hive_schema_check" "example" {
  service = "example-service"
  commit  = "57ee05c"
  schema  = file("schema.graphql")
}
