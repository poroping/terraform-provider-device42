package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poroping/libdevice42/client"
	ipam "github.com/poroping/libdevice42/client/ip_a_m"
	"github.com/poroping/libdevice42/models"
)

func resourceIpamIP() *schema.Resource {
	return &schema.Resource{
		Description: "Manage IPAM ips.",

		CreateContext: resourceIpamIPCreate,
		ReadContext:   resourceIpamIPRead,
		UpdateContext: resourceIpamIPUpdate,
		DeleteContext: resourceIpamIPDelete,

		Importer: nil,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "IP address ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ipaddress": {
				Description: "IP address.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"notes": {
				Description: "Notes.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"subnet_id": {
				Description: "Subnet ID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"suggest_ip": {
				Description:   "Get next free IP in subnet.",
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"ipaddress"},
			},
		},
	}
}

func resourceIpamIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewPostIPAMIpsParams()

	if v, ok := d.GetOk("ipaddress"); ok {
		if s, ok := v.(string); ok {
			params.Ipaddress = s
		}
	}

	if d.Get("suggest_ip").(bool) {
		err, ip := ipamSuggestIP(ctx, d, meta)
		if err != nil {
			return err
		}
		params.Ipaddress = *ip
	}

	if v, ok := d.GetOk("notes"); ok {
		if s, ok := v.(string); ok {
			params.Notes = &s
		}
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		if s, ok := v.(string); ok {
			params.SubnetID = &s
		}
	}

	resp, err := client.IPam.PostIPAMIps(params)

	if err != nil {
		return diag.Errorf("error creating IP. %s", err)
	}

	j_code := resp.Payload.Code.(json.Number)
	code, _ := j_code.Int64()
	msg := intList(resp.Payload.Msg.([]interface{}))
	if code != 0 {
		return diag.Errorf("error creating IP. %s", msg[0])
	}

	d.SetId(string(msg[1]))

	return resourceIpamIPRead(ctx, d, meta)
}

func ipamSuggestIP(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, *string) {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMSuggestIPParams()

	if v, ok := d.GetOk("subnet_id"); ok {
		if s, ok := v.(string); ok {
			params.SubnetID = &s
		}
	}

	resp, err := client.IPam.GetIPAMSuggestIP(params)

	if err != nil {
		return diag.Errorf("error reading suggest IP response. %s", err), nil
	}

	ip := resp.Payload.IP.(string)

	return nil, &ip
}

func resourceIpamIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMIpsParams()
	id := d.Id()
	params.SetIPID(&id)

	resp, err := client.IPam.GetIPAMIps(params)

	if err != nil {
		return diag.Errorf("error reading IPAM IP. %s", err)
	}

	ips := resp.Payload.Ips

	if len(ips) == 0 {
		return diag.Errorf("error IP ID not found.")
	}

	if len(ips) > 1 {
		return diag.Errorf("error more than one IP found.")
	}

	setIpamIP(d, resp.Payload.Ips[0])

	return nil
}

func resourceIpamIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewPostIPAMIpsParams()
	id := d.Id()
	params.SetIPID(&id)

	if v, ok := d.GetOk("notes"); ok {
		if s, ok := v.(string); ok {
			params.Notes = &s
		}
	}

	resp, err := client.IPam.PostIPAMIps(params)

	if err != nil {
		return diag.Errorf("error updating IP. %s", err)
	}

	j_code := resp.Payload.Code.(json.Number)
	code, _ := j_code.Int64()
	msg := intList(resp.Payload.Msg.([]interface{}))
	if code != 0 {
		return diag.Errorf("error updating IP. %s", msg[0])
	}

	d.SetId(string(msg[1]))

	return resourceIpamIPRead(ctx, d, meta)
}

func resourceIpamIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewDeleteIPAMIpsParams()
	id := d.Id()
	i, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("error getting IP ID. %s", err)
	}
	params.SetID(int64(i))

	resp, err := client.IPam.DeleteIPAMIps(params)

	if err != nil {
		return diag.Errorf("error deleting IPAM IP. %s", err)
	}

	del := resp.Payload.Deleted.(string)

	b, err := strconv.ParseBool(del)
	if err != nil {
		return diag.Errorf("error reading delete response. %s", err)
	}

	if !b {
		return diag.Errorf("error deleting IPAM IP.")
	}

	d.SetId("")

	return nil
}

func setIpamIP(d *schema.ResourceData, resp *models.IPAMips) {
	if v, ok := resp.IP.(string); ok {
		d.Set("ipaddress", v)
	}
	if v, ok := resp.Notes.(string); ok {
		d.Set("notes", v)
	}
}
