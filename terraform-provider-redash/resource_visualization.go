package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
	"strconv"
)

func resourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashVisualizationRead,
		CreateContext: resourceRedashVisualizationCreate,
		UpdateContext: resourceRedashVisualizationUpdate,
		DeleteContext: resourceRedashVisualizationDelete,
		Schema: map[string]*schema.Schema{
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRedashVisualizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	queryId := d.Get("query_id").(int)
	visualizationId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	visualization, err := c.GetVisualization(queryId, visualizationId)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", visualization.Name)
	_ = d.Set("type", visualization.Type)
	_ = d.Set("description", visualization.Description)

	return diags
}

func resourceRedashVisualizationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	payload := redash.VisualizationCreatePayload{
		QueryId:     d.Get("query_id").(int),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Options:     redash.VisualizationOptions{},
	}
	visualization, err := c.CreateVisualization(&payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(visualization.ID))

	return diags
}

func resourceRedashVisualizationUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	payload := redash.VisualizationUpdatePayload{
		Name: d.Get("name").(string),
	}
	visualization, err := c.UpdateVisualization(id, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", visualization.Name)

	return diags
}

func resourceRedashVisualizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteVisualization(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
