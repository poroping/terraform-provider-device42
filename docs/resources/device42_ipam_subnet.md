---
page_title: "device42_ipam_subnet Resource - terraform-provider-device42"
subcategory: ""
description: |-
  Manage ipam_subnet in the Terraform provider device42.
---

# Resource device42_ipam_subnet

Manage ipam_subnet in the Terraform provider device42.

## Example Usage

```terraform
resource "device42_ipam_subnet" "parent" {
    name = "SUPERNET"
    tags = "TEST"
    mask_bits = "21"
    network = "10.25.0.0"
}

resource "device42_ipam_subnet" "example" {
    create_from_parent = true

    name = "TEST-SUBNET"
    tags = "TEST,FART"
    mask_bits = "29"
    parent_subnet_id = device42_ipam_subnet.parent.subnet_id
}

output "example" {
    value = device42_ipam_subnet.example
}
```

## Argument Reference

* `mask_bits` - (Required) Netmask bits.
* `customer_id` - Customer ID.
* `name` - Name.
* `network` - Netmask address.
* `parent_mask_bits` - Parent netmask bits.
* `parent_subnet_id` - ID of the parent subnet.
* `parent_vlan_id` - Parent vlan ID.
* `mask_bits` - Netmask bits.
* `subnet_id` - ID of the subnet.
* `tags` - Tags.
* `create_from_parent` - Use to create subnet from parent.
* `check_if_exists` - Use to check if subnet exists already.

If addition to above the resource exports the following attributes:

## Attribute Reference

* `id` - Resource ID.
* `parent_vlan_name` - Parent vlan name.
* `parent_vlan_number` - Parent vlan number.


