package secberus

import (
	"context"
	"strings"
	"time"

	"github.com/RexBelli/go-secberus/secberus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationUserBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationUserBindingCreate,
		ReadContext:   resourceOrganizationUserBindingRead,
		UpdateContext: resourceOrganizationUserBindingUpdate,
		DeleteContext: resourceOrganizationUserBindingDelete,
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
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: false,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceOrganizationUserBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	orgID := d.Get("org_id").(string)

	users, ok := d.Get("users").([]interface{})
	if !ok {
		users = nil
	}

	var userIDs []string
	for _, user := range users {
		userID, err := c.UserEmailToId(user.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		userIDs = append(userIDs, userID)
	}

	err := c.SetOrganizationUsers(orgID, userIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))
	d.Set("org_id", orgID)
	d.SetId(orgID)

	resourceOrganizationUserBindingRead(ctx, d, m)

	return diags
}

func resourceOrganizationUserBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	orgID := d.Get("id").(string)

	users, err := c.GetOrganizationUsers(orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	var userIDs []string
	for _, user := range *users {
		userIDs = append(userIDs, user.Id)
	}

	d.Set("users", strings.Join(userIDs, ","))

	return diags
}

func resourceOrganizationUserBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	orgID := d.Get("id").(string)

	users, ok := d.Get("users").([]interface{})
	if !ok {
		users = nil
	}

	var userIDs []string
	for _, user := range users {
		userID, err := c.UserEmailToId(user.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		userIDs = append(userIDs, userID)
	}

	err := c.SetOrganizationUsers(orgID, userIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceOrganizationUserBindingRead(ctx, d, m)
}

func resourceOrganizationUserBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)
	var diags diag.Diagnostics

	orgID := d.Id()

	emptyUsers := []string{}

	err := c.SetOrganizationUsers(orgID, emptyUsers)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
