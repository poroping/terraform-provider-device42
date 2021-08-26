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

func resourceIpamSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage IPAM subnets.",

		CreateContext: resourceIpamSubnetCreate,
		ReadContext:   resourceIpamSubnetRead,
		UpdateContext: resourceIpamSubnetUpdate,
		DeleteContext: resourceIpamSubnetDelete,

		Importer: nil,

		Schema: map[string]*schema.Schema{
			"mask_bits": {
				Description: "Netmask bits.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"customer_id": {
				Description: "Customer ID.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "Name.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"network": {
				Description: "Network address.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"parent_mask_bits": {
				Description: "Parent netmask bits.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"parent_subnet_id": {
				Description: "ID of the parent subnet.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"parent_vlan_id": {
				Description: "Parent vlan ID.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
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
				Computed:    true,
				Optional:    true,
			},
			"tags": {
				Description: "Tags.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"create_from_parent": {
				Description:  "Use to create subnet from parent.",
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				RequiredWith: []string{"parent_subnet_id"},
			},
			"check_if_exists": {
				Description: "Use to check if subnet exists already.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceIpamSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	if v, ok := d.GetOk("subnet_id"); ok {
		if s, ok := v.(string); ok {
			d.SetId(s)
			return resourceIpamSubnetUpdate(ctx, d, meta)
		}
	}

	params := ipam.NewPostIPAMsubnetsParams()

	if d.Get("check_if_exists").(bool) {
		err, subnet_id := ipamSubnetsCheckExist(ctx, d, meta)

		if err != nil {
			return err
		}

		if subnet_id != nil {
			params.SubnetID = subnet_id

			subnet_id_int, err := strconv.Atoi(*subnet_id)

			if err != nil {
				return diag.Errorf("error reading subnet_id. %s", err)
			}

			read_params := ipam.NewGetIPAMSubnetIDParams()
			read_params.SetSubnetID(int64(subnet_id_int))

			resp2, err := client.IPam.GetIPAMSubnetID(read_params)

			if err != nil {
				return diag.Errorf("error reading IPAM subnet. %s", err)
			}

			if v, ok := resp2.Payload.MaskBits.(string); ok {
				d.Set("mask_bits", v)
			}

			if v, ok := resp2.Payload.Network.(string); ok {
				d.Set("network", v)
			}
		}
	}

	if d.Get("create_from_parent").(bool) {
		if params.SubnetID == nil {
			return ipamSubnetsCreateChildCreate(ctx, d, meta)
		}
	}

	if v, ok := d.GetOk("mask_bits"); ok {
		if s, ok := v.(string); ok {
			params.MaskBits = s
		}
	}
	if v, ok := d.GetOk("customer_id"); ok {
		if s, ok := v.(string); ok {
			params.CustomerID = &s
		}
	}
	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			params.Name = &s
		}
	}
	if v, ok := d.GetOk("network"); ok {
		if s, ok := v.(string); ok {
			params.Network = s
		}
	}
	if v, ok := d.GetOk("network"); ok {
		if s, ok := v.(string); ok {
			params.Network = s
		}
	}
	if v, ok := d.GetOk("parent_subnet_id"); ok {
		if s, ok := v.(string); ok {
			params.ParentSubnetID = &s
		}
	}
	if v, ok := d.GetOk("parent_vlan_id"); ok {
		if s, ok := v.(string); ok {
			params.ParentVlanID = &s
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		if s, ok := v.(string); ok {
			params.Tags = &s
		}
	}

	resp, err := client.IPam.PostIPAMsubnets(params)

	if err != nil {
		return diag.Errorf("error creating subnet. %s", err)
	}

	j_code := resp.Payload.Code.(json.Number)
	code, _ := j_code.Int64()
	msg := intList(resp.Payload.Msg.([]interface{}))
	if code != 0 {
		return diag.Errorf("error creating subnet. %s", msg[0])
	}

	d.SetId(string(msg[1]))

	return resourceIpamSubnetRead(ctx, d, meta)
}

func ipamSubnetsCreateChildCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewPostIPAMSubnetsCreateChildParams()

	if v, ok := d.GetOk("parent_subnet_id"); ok {
		if s, ok := v.(string); ok {
			params.ParentSubnetID = &s
		}
	}

	if v, ok := d.GetOk("mask_bits"); ok {
		if s, ok := v.(string); ok {
			params.MaskBits = s
		}
	}

	resp, err := client.IPam.PostIPAMSubnetsCreateChild(params)

	if err != nil {
		return diag.Errorf("error creating child subnet. %s", err)
	}

	subnet_id, err := resp.Payload.SubnetID.(json.Number).Int64()

	if err != nil {
		return diag.Errorf("error read child subnet_id. %s", err)
	}

	subnet_id_s := resp.Payload.SubnetID.(json.Number).String()
	d.SetId(subnet_id_s)

	read_params := ipam.NewGetIPAMSubnetIDParams()
	read_params.SetSubnetID(subnet_id)

	resp2, err := client.IPam.GetIPAMSubnetID(read_params)

	if err != nil {
		return diag.Errorf("error reading IPAM subnet. %s", err)
	}

	if v, ok := resp2.Payload.Network.(string); ok {
		d.Set("network", v)
	}

	return resourceIpamSubnetUpdate(ctx, d, meta)
}

func ipamSubnetsCheckExist(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, *string) {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMsubnetsParams()

	if v, ok := d.GetOk("mask_bits"); ok {
		if s, ok := v.(string); ok {
			params.MaskBits = &s
		}
	}

	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			params.Name = &s
		}
	}

	if v, ok := d.GetOk("network"); ok {
		if s, ok := v.(string); ok {
			params.Network = &s
		}
	}

	if v, ok := d.GetOk("parent_subnet_id"); ok {
		if s, ok := v.(string); ok {
			params.ParentSubnetID = &s
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		if s, ok := v.(string); ok {
			params.Tags = &s
		}
	}

	resp, err := client.IPam.GetIPAMsubnets(params)

	if err != nil {
		return diag.Errorf("error reading response. %s", err), nil
	}

	subnets := resp.Payload.Subnets

	if len(subnets) == 0 {
		return nil, nil
	}

	if len(subnets) > 1 {
		return diag.Errorf("error multiple subnets found, filter better."), nil
	}

	subnet_id := resp.Payload.Subnets[0].SubnetID.(json.Number).String()

	return nil, &subnet_id
}

func resourceIpamSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMSubnetIDParams()
	id := d.Id()
	i, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("error getting subnetid. %s", err)
	}
	params.SetSubnetID(int64(i))

	resp, err := client.IPam.GetIPAMSubnetID(params)

	if err != nil {
		return diag.Errorf("error reading IPAM subnet. %s", err)
	}

	setIpamSubnet(d, resp.Payload)

	return nil
}

func resourceIpamSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewPostIPAMsubnetsParams()
	id := d.Id()

	params.SubnetID = &id

	if v, ok := d.GetOk("mask_bits"); ok {
		if s, ok := v.(string); ok {
			params.MaskBits = s
		}
	}
	if v, ok := d.GetOk("customer_id"); ok {
		if s, ok := v.(string); ok {
			params.CustomerID = &s
		}
	}
	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			params.Name = &s
		}
	}
	if v, ok := d.GetOk("network"); ok {
		if s, ok := v.(string); ok {
			params.Network = s
		}
	}
	if v, ok := d.GetOk("parent_subnet_id"); ok {
		if s, ok := v.(string); ok {
			params.ParentSubnetID = &s
		}
	}
	if v, ok := d.GetOk("parent_vlan_id"); ok {
		if s, ok := v.(string); ok {
			params.ParentVlanID = &s
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		if s, ok := v.(string); ok {
			params.Tags = &s
		}
	}

	resp, err := client.IPam.PostIPAMsubnets(params)

	if err != nil {
		return diag.Errorf("error updating subnet. %s", err)
	}

	j_code := resp.Payload.Code.(json.Number)
	code, _ := j_code.Int64()
	msg := intList(resp.Payload.Msg.([]interface{}))
	if code != 0 {
		return diag.Errorf("error updating subnet. %s", msg[0])
	}

	d.SetId(string(msg[1]))

	return resourceIpamSubnetRead(ctx, d, meta)
}

func resourceIpamSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewDeleteIPAMsubnetsParams()
	id := d.Id()
	i, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("error getting subnetid. %s", err)
	}
	params.SetSubnetID(int64(i))

	resp, err := client.IPam.DeleteIPAMsubnets(params)

	if err != nil {
		return diag.Errorf("error deleting IPAM subnet. %s", err)
	}

	del := resp.Payload.Deleted.(string)

	b, err := strconv.ParseBool(del)
	if err != nil {
		return diag.Errorf("error reading delete response. %s", err)
	}

	if !b {
		return diag.Errorf("error deleting IPAM subnet.")
	}

	d.SetId("")

	return nil
}

func setIpamSubnet(d *schema.ResourceData, resp *models.IPAMsubnets) {
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
