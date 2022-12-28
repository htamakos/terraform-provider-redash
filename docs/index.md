# Redash Provider

Terraform provider for managing Redash configurations.

## Example Usage

```hcl
terraform {
  required_providers {
    redash = {
      source  = "htamakos/redash"
      version = "~> 0.0"
    }
  }
}

provider "redash" {
  redash_uri = "https://redash.exmaple.com"
}
```

## Argument Reference

* `redash_uri` - (Optional) The host URL to the Redash instance you will be managing, including protocol
  (for example `https://redash.exmaple.com` or `http://localhost:5000`). It must be provided, but it can also be sourced
  from the `REDASH_HOST` environment variable.
* `api_key` - (Optional) A Redash API token. It must be provided, but it can also be sourced from the `REDASH_API_KEY`
  environment variable.
