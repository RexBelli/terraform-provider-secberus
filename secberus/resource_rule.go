package secberus

import (
	"context"
	"time"

	"github.com/RexBelli/go-secberus/secberus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleCreate,
		ReadContext:   resourceRuleRead,
		UpdateContext: resourceRuleUpdate,
		DeleteContext: resourceRuleDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"summary": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"logic": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"remediation_steps": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"alert_summary_tmpl": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: false,
				Required: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"resources": &schema.Schema{
				Type:     schema.TypeList,
				Computed: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Required: true,
						},
						"resource_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Optional: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Required: true,
						},
						"data_provider": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Required: true,
						},
						"required": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: false,
							Optional: true,
						},
						"score": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: false,
							Optional: true,
						},
					},
				},
			},
			"compliances": &schema.Schema{
				Type:     schema.TypeList,
				Computed: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Required: true,
						},
					},
				},
			},
			"alert_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"score": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: false,
				Required: true,
			},
			"subscribed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: false,
				Optional: true,
			},
		},
	}
}

func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	rule := ruleFromTF(ctx, d, m)

	r, err := c.CreateRule(rule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))
	d.SetId(r.ID)

	resourceRuleRead(ctx, d, m)

	return diags
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	id := d.Get("id").(string)

	_, err := c.GetRule(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	rule := ruleFromTF(ctx, d, m)
	id := d.Id()

	_, err := c.SetRule(id, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceRuleRead(ctx, d, m)
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)
	var diags diag.Diagnostics

	id := d.Id()

	err := c.DeleteRule(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func ruleFromTF(ctx context.Context, d *schema.ResourceData, m interface{}) *secberus.Rule {
	rule := &secberus.Rule{}

	// required fields
	desc := d.Get("description").(string)
	policyID := d.Get("policy_id").(string)
	priority := d.Get("priority").(float64)
	orgID := d.Get("org_id").(string)
	score := d.Get("score").(float64)

	rule.Description = desc
	rule.PolicyID = policyID
	rule.Priority = priority
	rule.OrgID = orgID
	rule.Score = score

	// optional fields
	summary, ok := d.Get("summary").(string)
	if !ok {
		summary = ""
	}
	logic, ok := d.Get("logic").(string)
	if !ok {
		logic = ""
	}
	remSteps, ok := d.Get("remediations_steps").(string)
	if !ok {
		remSteps = ""
	}
	alertSummaryTmpl, ok := d.Get("alert_summary_tmpl").(string)
	if !ok {
		alertSummaryTmpl = ""
	}
	subscribed, ok := d.Get("subscribed").(bool)
	if !ok {
		subscribed = false
	}

	rule.Summary = summary
	rule.Logic = logic
	rule.RemediationSteps = remSteps
	rule.AlertSummaryTmpl = alertSummaryTmpl
	rule.Subscribed = subscribed

	// complex fields
	resources := d.Get("resources").([]interface{})
	if !ok {
		resources = nil
	}
	compliances, ok := d.Get("compliances").([]interface{})
	if !ok {
		compliances = nil
	}

	for _, res := range resources {
		r := res.(map[string]interface{})

		required, ok := r["required"].(bool)
		if !ok {
			required = true
		}

		rule.Resources = append(rule.Resources, secberus.Resource{
			ID:           r["id"].(string),
			ResourceID:   r["id"].(string),
			Name:         r["name"].(string),
			DataProvider: r["data_provider"].(string),
			Score:        r["score"].(int),
			Required:     required,
		})
	}

	for _, comp := range compliances {
		rule.Compliances = append(rule.Compliances, comp.(secberus.Compliance))
	}

	return rule
}
