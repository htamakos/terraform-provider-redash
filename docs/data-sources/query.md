# Query Data Source

Data source representation of a Redash Query

## Example Usage

```hcl
data "redash_query" "my_query" {
  id = 1
}

output "example" {
  value = jsonencode(data.redash_query.my_query)
}
```

## Argument Reference

* `id` - (Required) Query ID to load

## Attribute Reference

* `id` - Redash ID of this query
* `name` - Name of Redash query
* `query` - Query using the query language native to the data source
* `data_source_id` - ID of the data source
* `description` - Description of the Redash query
