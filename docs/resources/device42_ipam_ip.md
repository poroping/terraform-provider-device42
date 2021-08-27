---
page_title: "device42_ipam_ip Resource - terraform-provider-device42"
subcategory: ""
description: |-
  Manage ipam_ip in the Terraform provider device42.
---

# Resource device42_ipam_ip

Manage ipam_ip in the Terraform provider device42.

## Example Usage

```terraform
resource "device42_ipam_subnet" "example" {
  name      = "EXAMPLE"
  tags      = "EXAMPLE"
  mask_bits = "29"
  network   = "10.25.0.0"
}

resource "device42_ipam_ip" "example" {
  subnet_id = "1312"
  ipaddress = "10.25.0.1"
  notes     = "server1.example.com"
}

output "example" {
  value = device42_ipam_ip.example
}

resource "device42_ipam_ip" "example2" {
  subnet_id  = "1312"
  suggest_ip = true
  notes      = "server2.example2.com"
}

output "example2" {
  value = device42_ipam_ip.example2
}
```

## Argument Reference

* `subnet_id` - (Required) Subnet ID.
* `ipaddress` - IP address.
* `notes` - Notes.
* `suggest_ip` - Get next free IP in subnet.

In addition to above the resource exports the following attributes:

## Attribute Reference

* `id` - Resource ID.


