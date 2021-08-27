---
page_title: "device42_ipam_subnet Data Source - terraform-provider-device42"
subcategory: ""
description: |-
  Get subnet info with the Terraform provider device42.
---

# Data Source device42_ipam_subnet

Get subnet info with the Terraform provider device42.

## Example Usage

```terraform
data "device42_ipam_subnet" "example" {
  subnet_id = "1"
}

output "example" {
  value = data.device42_ipam_subnet.example
}
```

## Argument Reference

- **subnet_id** (Required) Subnet ID.

In addition to above the resource exports the following attributes:

## Attribute Reference

- **id** The ID of this resource.
- **mask_bits** Netmask bits.
- **customer_id** Customer ID. 
- **name** Name.
- **network** Network address.
- **parent_mask_bits** Parent netmask bits.
- **parent_subnet_id** ID of the parent subnet.
- **parent_vlan_id** Parent vlan ID.
- **parent_vlan_name** Parent vlan name.
- **parent_vlan_number** Parent vlan number.
- **tags** Tags.


