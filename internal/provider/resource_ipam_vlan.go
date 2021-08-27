package provider

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/libdevice42/client"
	ipam "github.com/poroping/libdevice42/client/ip_a_m"
	"github.com/poroping/libdevice42/models"

	funk "github.com/thoas/go-funk"
)

func resourceIpamVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage IPAM vlans.",

		CreateContext: resourceIpamVlanCreate,
		ReadContext:   resourceIpamVlanRead,
		UpdateContext: resourceIpamVlanUpdate,
		DeleteContext: resourceIpamVlanDelete,

		Importer: nil,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"number": {
				Description: "VLAN number.",
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
			// "tags_or": {
			// 	Description: "Tags (OR).",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			"tags_and": {
				Description: "Tags (AND).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vlan_id": {
				Description: "VLAN ID.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"create_within_range": {
				Description:   "Use to create vlan from a range of vlans.",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"tags_and", "name"},
				ConflictsWith: []string{"number"},
				ValidateFunc:  validation.StringMatch(validateRegExpVlanRange(), "Provide valid VLAN range"),
			},
			"check_if_exists": {
				Description: "Use to check if vlan exists already.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceIpamVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	if v, ok := d.GetOk("vlan_id"); ok {
		if s, ok := v.(string); ok {
			d.SetId(s)
			return resourceIpamVlanUpdate(ctx, d, meta)
		}
	}

	params := ipam.NewPostIPAMvlansParams()

	if d.Get("check_if_exists").(bool) {
		err, vlan_id := ipamVlansCheckExist(ctx, d, meta)

		if err != nil {
			return err
		}

		if vlan_id != nil {
			d.SetId(*vlan_id)
			return resourceIpamVlanUpdate(ctx, d, meta)
		}
	}

	if _, ok := d.GetOk("create_within_range"); ok {
		next_vlan, _ := ipamVlanFromRange(ctx, d, meta)
		params.Number = strconv.Itoa(*next_vlan)
	}

	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			params.Name = &s
		}
	}
	if v, ok := d.GetOk("number"); ok {
		if s, ok := v.(string); ok {
			params.Number = s
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		if s, ok := v.(string); ok {
			params.Tags = &s
		}
	}

	resp, err := client.IPam.PostIPAMvlans(params)

	if err != nil {
		return diag.Errorf("error creating vlan. %s", err)
	}

	j_code := resp.Payload.Code.(json.Number)
	code, _ := j_code.Int64()
	msg := intList(resp.Payload.Msg.([]interface{}))
	if code != 0 {
		return diag.Errorf("error creating vlan. %s", msg[0])
	}

	d.SetId(string(msg[1]))

	return resourceIpamVlanRead(ctx, d, meta)
}

func ipamVlanFromRange(ctx context.Context, d *schema.ResourceData, meta interface{}) (*int, diag.Diagnostics) {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMvlansParams()

	if v, ok := d.GetOk("tags"); ok {
		if s, ok := v.(string); ok {
			params.Tags = &s
		}
	}

	vlan_range := make([]int, 0)

	if v, ok := d.GetOk("create_within_range"); ok {
		if s, ok := v.(string); ok {
			mins := strings.Split(s, "-")[0]
			maxs := strings.Split(s, "-")[1]
			min, _ := strconv.Atoi(mins)
			max, _ := strconv.Atoi(maxs)
			for i := min; i < max; i++ {
				vlan_range = append(vlan_range, i)
			}
		}
	}

	if len(vlan_range) == 0 {
		return nil, diag.Errorf("no vlans in range.")
	}

	// find used vlans based on tag

	used_vlans := make([]int, 0)

	resp, err := client.IPam.GetIPAMvlans(params)

	if err != nil {
		return nil, diag.Errorf("error reading vlans. %s", err)
	}

	vlans := resp.Payload.Vlans

	for _, vlan := range vlans {
		num, _ := vlan.Number.(json.Number).Int64()
		used_vlans = append(used_vlans, int(num))
		sort.Ints(used_vlans)
	}

	next_vlan := funk.Subtract(vlan_range, used_vlans).([]int)[0]

	return &next_vlan, nil
}

func ipamVlansCheckExist(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, *string) {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMvlansParams()

	if v, ok := d.GetOk("number"); ok {
		if s, ok := v.(string); ok {
			params.Number = &s
		}
	}

	if v, ok := d.GetOk("tags_and"); ok {
		if s, ok := v.(string); ok {
			params.TagsAnd = &s
		}
	}

	resp, err := client.IPam.GetIPAMvlans(params)

	if err != nil {
		return diag.Errorf("error reading response. %s", err), nil
	}

	vlans := resp.Payload.Vlans

	if len(vlans) == 0 {
		return nil, nil
	}

	if len(vlans) > 1 {
		return diag.Errorf("error multiple vlans found, filter better."), nil
	}

	vlan_id := resp.Payload.Vlans[0].VlanID.(json.Number).String()

	return nil, &vlan_id
}

func resourceIpamVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewGetIPAMvlansIDParams()
	id := d.Id()
	i, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("error getting vlanid. %s", err)
	}
	params.SetID(int64(i))

	resp, err := client.IPam.GetIPAMvlansID(params)

	if err != nil {
		return diag.Errorf("error reading IPAM vlan. %s", err)
	}

	setIpamVlan(d, resp.Payload)

	return nil
}

func resourceIpamVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewPutIPAMvlansParams()
	id := d.Id()
	i, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("error getting vlanid. %s", err)
	}
	params.SetID(int64(i))

	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			params.Name = &s
		}
	}
	if v, ok := d.GetOk("number"); ok {
		if s, ok := v.(string); ok {
			params.Number = &s
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		if s, ok := v.(string); ok {
			params.Tags = &s
		}
	}

	resp, err := client.IPam.PutIPAMvlans(params)

	if err != nil {
		return diag.Errorf("error updating vlan. %s", err)
	}

	j_code := resp.Payload.Code.(json.Number)
	code, _ := j_code.Int64()
	msg := intList(resp.Payload.Msg.([]interface{}))
	if code != 0 {
		return diag.Errorf("error updating vlan. %s", msg[0])
	}

	d.SetId(string(msg[1]))

	return resourceIpamVlanRead(ctx, d, meta)
}

func resourceIpamVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Device42)

	params := ipam.NewDeleteIPAMvlansParams()
	id := d.Id()
	i, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("error getting vlanid. %s", err)
	}
	params.SetID(int64(i))

	resp, err := client.IPam.DeleteIPAMvlans(params)

	if err != nil {
		return diag.Errorf("error deleting IPAM vlan. %s", err)
	}

	del := resp.Payload.Deleted.(string)

	b, err := strconv.ParseBool(del)
	if err != nil {
		return diag.Errorf("error reading delete response. %s", err)
	}

	if !b {
		return diag.Errorf("error deleting IPAM vlan.")
	}

	d.SetId("")

	return nil
}

func setIpamVlan(d *schema.ResourceData, resp *models.IPAMvlans) {
	if v, ok := resp.Name.(string); ok {
		d.Set("name", v)
	}
	if v, ok := resp.Number.(json.Number); ok {
		d.Set("number", v.String())
	}
	if v := resp.Tags; v != nil {
		d.Set("tags", strings.Join(v, ","))
	}
	if v, ok := resp.VlanID.(json.Number); ok {
		d.Set("vlan_id", v.String())
	}
}
