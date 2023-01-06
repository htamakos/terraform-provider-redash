package main

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
)

func resourceRedashAlertDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashAlertDestinationRead,
		CreateContext: resourceRedashAlertDestinationCreate,
		UpdateContext: resourceRedashAlertDestinationUpdate,
		DeleteContext: resourceRedashAlertDestinationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subject_template": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "({state}) {alert_name}",
						},
					}},
			},
			"slack_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icon_emoji": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"channel": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"webhook_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					}},
			},
			"hipchat_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					}},
			},
			"mattermost_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"channel": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"chatwork_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"room_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"message_template": {
							Type:     schema.TypeString,
							Default:  "{alert_name} changed state to {new_state}.\\n{alert_url}\\n{query_url}",
							Optional: true,
						},
					}},
			},
			"pagerduty_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"google_hangouts_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},

			"icon": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRedashAlertDestinationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	alertDestination, err := c.GetDestination(id)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", alertDestination.Name)
	_ = d.Set("type", alertDestination.Type)
	_ = d.Set("icon", alertDestination.Icon)

	switch alertDestination.Type {
	case "email":
		_ = d.Set("email_options", alertDestination.Options)
	case "slack":
		_ = d.Set("slack_options", alertDestination.Options)
	case "webhook":
		_ = d.Set("webhook_options", alertDestination.Options)
	case "hipchat":

		_ = d.Set("hipchat_options", alertDestination.Options)
	case "mattermost":

		_ = d.Set("mattermost_options", alertDestination.Options)
	case "chatwork":

		_ = d.Set("chatwork_options", alertDestination.Options)
	case "pagerduty":

		_ = d.Set("pagerduty_options", alertDestination.Options)
	case "hangouts_chat":
		_ = d.Set("google_hangouts_options", alertDestination.Options)
	}

	return diags
}

func resourceRedashAlertDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	createPayload := redash.CreateOrUpdateDestinationPayload{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	email_options, ok := d.GetOkExists("email_options")
	if ok {
		if val := expandOptions(email_options); val != nil {
			createPayload.Options = val
		}
	}

	slack_options, ok := d.GetOkExists("slack_options")
	if ok {
		if val := expandOptions(slack_options); val != nil {
			createPayload.Options = val
		}
	}

	webhook_options, ok := d.GetOkExists("webhook_options")
	if ok {
		if val := expandOptions(webhook_options); val != nil {
			createPayload.Options = val
		}
	}

	hipchat_options, ok := d.GetOkExists("hipchat_options")
	if ok {
		if val := expandOptions(hipchat_options); val != nil {
			createPayload.Options = val
		}
	}

	mattermost_options, ok := d.GetOkExists("mattermost_options")
	if ok {
		if val := expandOptions(mattermost_options); val != nil {
			createPayload.Options = val
		}
	}

	chatwork_options, ok := d.GetOkExists("chatwork_options")
	if ok {
		if val := expandOptions(chatwork_options); val != nil {
			createPayload.Options = val
		}
	}

	pagerduty_options, ok := d.GetOkExists("pagerduty_options")
	if ok {
		if val := expandOptions(pagerduty_options); val != nil {
			createPayload.Options = val
		}
	}

	google_hangouts_options, ok := d.GetOkExists("google_hangouts_options")
	if ok {
		if val := expandOptions(google_hangouts_options); val != nil {
			createPayload.Options = val
		}
	}

	alertDestination, err := c.CreateDestination(&createPayload)
	if err != nil {
		return nil
	}

	d.SetId(strconv.Itoa(alertDestination.Id))
	resourceRedashAlertDestinationRead(ctx, d, meta)

	return diags
}

func expandOptions(v interface{}) map[string]interface{} {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	raw := l[0]
	return raw.(map[string]interface{})
}

func resourceRedashAlertDestinationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}
	updatePayload := redash.CreateOrUpdateDestinationPayload{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	email_options, ok := d.GetOkExists("email_options")
	if ok {
		if val := expandOptions(email_options); val != nil {
			updatePayload.Options = val
		}
	}

	slack_options, ok := d.GetOkExists("slack_options")
	if ok {
		if val := expandOptions(slack_options); val != nil {
			updatePayload.Options = val
		}
	}

	webhook_options, ok := d.GetOkExists("webhook_options")
	if ok {
		if val := expandOptions(webhook_options); val != nil {
			updatePayload.Options = val
		}
	}

	hipchat_options, ok := d.GetOkExists("hipchat_options")
	if ok {
		if val := expandOptions(hipchat_options); val != nil {
			updatePayload.Options = val
		}
	}

	mattermost_options, ok := d.GetOkExists("mattermost_options")
	if ok {
		if val := expandOptions(mattermost_options); val != nil {
			updatePayload.Options = val
		}
	}

	chatwork_options, ok := d.GetOkExists("chatwork_options")
	if ok {
		if val := expandOptions(chatwork_options); val != nil {
			updatePayload.Options = val
		}
	}

	pagerduty_options, ok := d.GetOkExists("pagerduty_options")
	if ok {
		if val := expandOptions(pagerduty_options); val != nil {
			updatePayload.Options = val
		}
	}

	google_hangouts_options, ok := d.GetOkExists("google_hangouts_options")
	if ok {
		if val := expandOptions(google_hangouts_options); val != nil {
			updatePayload.Options = val
		}
	}

	alertDestination, err := c.UpdateDestination(id, &updatePayload)
	if err != nil {
		return nil
	}

	d.SetId(strconv.Itoa(alertDestination.Id))
	resourceRedashAlertDestinationRead(ctx, d, meta)

	return diags

}

func resourceRedashAlertDestinationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteDestination(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
