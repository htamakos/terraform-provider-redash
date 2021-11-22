package main

import (
	"context"
	"strconv"

	"github.com/digitalpoetry/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRedashQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedashQueryCreate,
		ReadContext:   resourceRedashQueryRead,
		UpdateContext: resourceRedashQueryUpdate,
		DeleteContext: resourceRedashQueryArchive,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRedashQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	createPayload := redash.QueryCreatePayload{
		Name:         d.Get("name").(string),
		Query:        d.Get("query").(string),
		DataSourceID: d.Get("data_source_id").(int),
		Description:  d.Get("description").(string),
	}

	query, err := c.CreateQuery(&createPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(query.ID))
	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func resourceRedashQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	query, err := c.GetQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", query.Name)
	_ = d.Set("query", query.Query)
	_ = d.Set("data_source_id", query.DataSourceID)
	_ = d.Set("description", query.Description)

	return diags
}

func resourceRedashQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updatePayload := redash.QueryUpdatePayload{
		Name:         d.Get("name").(string),
		Query:        d.Get("query").(string),
		DataSourceID: d.Get("data_source_id").(int),
		Description:  d.Get("description").(string),
	}

	_, err = c.UpdateQuery(id, &updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func resourceRedashQueryArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.ArchiveQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
