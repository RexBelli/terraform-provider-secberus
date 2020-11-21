package secberus

import (
	"context"
	"time"

	"github.com/RexBelli/go-secberus/secberus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesRead,
		Schema: map[string]*schema.Schema{
			"resources": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": &schema.Schema{
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
						"data_provider": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"score": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceResourcesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*secberus.Client)

	var diags diag.Diagnostics

	resources, err := c.GetResources()
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenResourcesData(resources)
	if err := d.Set("resources", orderItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))

	return diags
}

func flattenResourcesData(resources *[]secberus.Resource) []interface{} {
	if resources != nil {
		ois := make([]interface{}, len(*resources), len(*resources))

		for i, res := range *resources {
			oi := make(map[string]interface{})

			oi["id"] = res.ID
			oi["resource_id"] = res.ID
			oi["description"] = res.Description
			oi["name"] = res.Name
			oi["data_provider"] = res.DataProvider
			oi["score"] = res.Score

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
