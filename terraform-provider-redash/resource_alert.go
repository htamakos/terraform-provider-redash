package main

import (
	"context"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
)

func resourceRedashAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashAlertRead,
		CreateContext: resourceRedashAlertCreate,
		UpdateContext: resourceRedashAlertUpdate,
		DeleteContext: resourceRedashAlertDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"op": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"column": {
							Type:     schema.TypeString,
							Required: true,
						},
						"muted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"custom_body": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_subject": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_triggered_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rearm": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRedashAlertRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := c.GetAlert(id)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", alert.Name)
	_ = d.Set("query_id", alert.Query.ID)
	_ = d.Set("options", alert.Options)
	_ = d.Set("state", alert.State)
	_ = d.Set("last_triggered_at", alert.LastTriggeredAt)
	_ = d.Set("updated_at", alert.UpdatedAt)
	_ = d.Set("created_at", alert.CreatedAt)
	_ = d.Set("rearm", alert.Rearm)
	_ = d.Set("user_id", alert.User.ID)

	return diags
}

func ptrString(s string) *string {
	return &s
}

func expandAlertOptions(v interface{}, d TerraformResourceData) (*redash.AlertOption, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	options := redash.AlertOption{}
	raw := l[0]
	original := raw.(map[string]interface{})

	op := original["op"]
	if val := reflect.ValueOf(op); val.IsValid() && !IsEmptyValue(val) {
		options.Op = val.String()
	} else {
		return nil, nil
	}

	value := original["value"]
	if val := reflect.ValueOf(value); val.IsValid() && !IsEmptyValue(val) {
		options.Value = val.Interface()
	}

	muted := original["muted"]
	if val := reflect.ValueOf(muted); val.IsValid() && !IsEmptyValue(val) {
		options.Muted = val.Bool()
	}

	column := original["column"]
	if val := reflect.ValueOf(column); val.IsValid() && !IsEmptyValue(val) {
		options.Column = val.String()
	}

	customSubject := original["custom_subject"]
	if val := reflect.ValueOf(customSubject); val.IsValid() && !IsEmptyValue(val) {
		options.CustomSubject = ptrString(val.String())
	}

	customBody := original["custom_body"]
	if val := reflect.ValueOf(customBody); val.IsValid() && !IsEmptyValue(val) {
		options.CustomBody = ptrString(val.String())
	}

	return &options, nil

}

func resourceRedashAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	createPayload := redash.CreateAlertPayload{
		Name:    d.Get("name").(string),
		QueryId: d.Get("query_id").(int),
	}

	options, ok := d.GetOkExists("options")
	if ok {
		val, err := expandAlertOptions(options, d)
		if err != nil {
			return diag.FromErr(err)
		}

		if val != nil {
			createPayload.Options = *val
		}
	}

	alert, err := c.CreateAlert(createPayload)
	if err != nil {
		return nil
	}

	d.SetId(strconv.Itoa(alert.ID))
	resourceRedashAlertRead(ctx, d, meta)

	return diags
}

func resourceRedashAlertUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	updatePayload := redash.UpdateAlertPayload{
		Name:    d.Get("name").(string),
		QueryId: d.Get("query_id").(int),
	}
	options, ok := d.GetOkExists("options")
	if ok {
		val, err := expandAlertOptions(options, d)
		if err != nil {
			return diag.FromErr(err)
		}

		if val != nil {
			updatePayload.Options = *val
		}
	}

	alert, err := c.UpdateAlert(id, &updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", alert.Name)
	_ = d.Set("query_id", alert.Query.ID)
	_ = d.Set("options", alert.Options)
	_ = d.Set("state", alert.State)
	_ = d.Set("last_triggered_at", alert.LastTriggeredAt)
	_ = d.Set("updated_at", alert.UpdatedAt)
	_ = d.Set("created_at", alert.CreatedAt)
	_ = d.Set("rearm", alert.Rearm)
	_ = d.Set("user_id", alert.User.ID)

	return diags
}

func resourceRedashAlertDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteAlert(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
