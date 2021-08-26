package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poroping/libdevice42/client"
	ipam "github.com/poroping/libdevice42/client/ip_a_m"
	"github.com/poroping/libdevice42/models"
)

func dataSourceIpamSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Read IPAM subnet.",

		ReadContext: dataSourceIpamSubnetRead,

		Schema: map[string]*schema.Schema{
			"mask_bits": {
				Description: "Netmask bits.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"customer_id": {
				Description: "Customer ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network": {
				Description: "Network address.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_mask_bits": {
				Description: "Parent netmask bits.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_subnet_id": {
				Description: "ID of the parent subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_vlan_id": {
				Description: "Parent vlan ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_vlan_name": {
				Description: "Parent vlan name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_vlan_number": {
				Description: "Parent vlan number.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subnet_id": {
				Description: "ID of the subnet.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tags": {
				Description: "Tags.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIpamSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMSubnetIDParams()

	id := d.Get("subnet_id").(string)

	if v, ok := d.GetOk("subnet_id"); ok {
		if s, ok := v.(string); ok {
			i, _ := strconv.Atoi(s)
			params.SubnetID = int64(i)
		}
	}

	resp, err := client.IPam.GetIPAMSubnetID(params)

	if err != nil {
		return diag.Errorf("error retrieving IPAM subnets. %s", err)
	}

	dataSetIpamSubnet(d, resp.Payload)

	d.SetId(id)

	return nil
}

func dataSetIpamSubnet(d *schema.ResourceData, resp *models.IPAMsubnets) {
	if v, ok := resp.CustomerID.(json.Number); ok {
		d.Set("customer_id", v.String())
	}
	if v, ok := resp.MaskBits.(json.Number); ok {
		d.Set("mask_bits", v.String())
	}
	if v, ok := resp.Name.(string); ok {
		d.Set("name", v)
	}
	if v, ok := resp.Network.(string); ok {
		d.Set("network", v)
	}
	if v, ok := resp.ParentVlanID.(json.Number); ok {
		d.Set("parent_vlan_id", v.String())
	}
	if v, ok := resp.ParentVlanName.(string); ok {
		d.Set("parent_vlan_name", v)
	}
	if v, ok := resp.ParentVlanNumber.(json.Number); ok {
		d.Set("parent_vlan_number", v)
	}
	if v, ok := resp.SubnetID.(json.Number); ok {
		d.Set("subnet_id", v.String())
	}
	if v := resp.Tags; v != nil {
		d.Set("tags", strings.Join(v, ","))
	}
}
