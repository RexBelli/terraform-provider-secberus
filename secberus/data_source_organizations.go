package secberus

import (
	"context"
	"time"

	"github.com/RexBelli/go-secberus/secberus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"organizations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	//orderID := d.Get("id").(string)

	orgs, err := c.GetOrganizations()
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenOrgsData(orgs)
	if err := d.Set("organizations", orderItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))

	return diags
}

func flattenOrgsData(orgs *[]secberus.Organization) []interface{} {
	if orgs != nil {
		ois := make([]interface{}, len(*orgs), len(*orgs))

		for i, org := range *orgs {
			oi := make(map[string]interface{})

			oi["id"] = org.Id
			oi["description"] = org.Description
			oi["name"] = org.Name

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
