data "hive_schema_check" "check" {
  schema  = <<EOF
      type Query {
    hello: String
  }

  schema {
    query: Query
  }
  EOF
  service = "subgraph"
}
