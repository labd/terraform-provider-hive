# GraphQL Hive Terraform Provider


## Usage

```hcl

terraform {
  required_providers {
    hive = {
      source  = "labd/hive"
      version = "0.0.1"
    }
  }
}

resource "hive_schema_check" "my-service" {
  service = "my-service"
  commit  = "57ee05c"
  schema  = file("schema.graphql")
}

resource "hive_schema_publish" "my-service" {
  service = "my-service"
  commit  = "57ee05c"
  url     = "https://checkout.example.com/graphql"
  schema  = file("schema.graphql")
}

resource "hive_app_create" "persisted_documents" {
  name      = "site"
  version   =  "1.0.0"
  documents = file("persisted_documents.json")
}

resource "hive_app_publish" "persisted_documents" {
  name      = "site"
  version   =  "1.0.0"
}
```
