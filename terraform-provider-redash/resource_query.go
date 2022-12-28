package main

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
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
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	tags, ok := d.Get("tags").([]interface{})
	if ok {
		updatePayload := redash.QueryUpdatePayload{
			Name:         createPayload.Name,
			Query:        createPayload.Query,
			Description:  createPayload.Description,
			DataSourceID: createPayload.DataSourceID,
			Tags:         convertToStringList(tags),
		}

		_, err = c.UpdateQuery(query.ID, &updatePayload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(query.ID))
	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func convertToStringList(slice []interface{}) []string {
	var arr []string
	for _, v := range slice {
		if v == nil {
			continue
		}
		arr = append(arr, v.(string))
	}

	return arr
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
	_ = d.Set("tags", query.Tags)

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

	tags, ok := d.Get("tags").([]interface{})
	if ok {
		updatePayload.Tags = convertToStringList(tags)
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
