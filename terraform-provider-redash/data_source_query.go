package main

import (
	"context"
	"fmt"

	"github.com/digitalpoetry/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRedashQuery() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: dataSourceRedashQueryRead,
	}
}

func dataSourceRedashQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id := d.Get("id").(int)
	query, err := c.GetQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(query.ID))
	_ = d.Set("name", query.Name)
	_ = d.Set("query", query.Query)
	_ = d.Set("data_source_id", query.DataSourceID)
	_ = d.Set("description", query.Description)

	return diags
}
