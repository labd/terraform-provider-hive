resource "hive_schema_publish" "example" {
  service = "example-service"
  commit  = "57ee05c"
  url     = "https://checkout.example.com/graphql"
  schema  = file("schema.graphql")
}
