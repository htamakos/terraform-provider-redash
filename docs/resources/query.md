# Query Resource

Allows creation/management of a Redash Query.

## Example Usage

```hcl
resource "redash_query" "my_query" {
  name           = "my_query"
  data_source_id = redash_data_source.gps.id
  query          = "SELECT 1+1"
  description    = "hello world"
}


output "example" {
  value = "${jsonencode(redash_query.my_query)}"
}
```

## Argument Reference

* `name` - (Required) Name of Redash query
* `query` - (Required) Query using the query language native to the data source
* `data_source_id` - (Required) ID of the data source
* `description` - (Optional) Description of the Redash query

## Attribute Reference

* `id` - Redash query ID
* `name` - Name of Redash query
* `query` - Query using the query language native to the data source
* `data_source_id` - ID of the data source
* `description` - Description of the Redash query
