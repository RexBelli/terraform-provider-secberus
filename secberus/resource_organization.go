package secberus

import (
	"context"
	"time"

	"github.com/RexBelli/go-secberus/secberus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		UpdateContext: resourceOrganizationUpdate,
		DeleteContext: resourceOrganizationDelete,
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
		},
	}
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	//	orgId := d.Get("id").(string)
	name := d.Get("name").(string)
	desc, ok := d.Get("description").(string)
	if !ok {
		desc = ""
	}

	orgIn := secberus.Organization{
		Name:        name,
		Description: desc,
	}

	org, err := c.CreateOrganization(orgIn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))
	d.SetId(org.Id)

	resourceOrganizationRead(ctx, d, m)

	return diags
}

func resourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	orgID := d.Get("id").(string)

	org, err := c.GetOrganization(orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", org.Name)
	d.Set("description", org.Description)

	return diags
}

func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	orgID := d.Id()

	cur, err := c.GetOrganization(orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		cur.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		cur.Description = d.Get("description").(string)
	}

	err = c.UpdateOrganization(cur)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceOrganizationRead(ctx, d, m)
}

func resourceOrganizationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)
	var diags diag.Diagnostics

	orgID := d.Id()

	err := c.DeleteOrganization(orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
