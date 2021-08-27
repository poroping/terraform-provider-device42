resource "device42_ipam_vlan" "example" {
  name   = "VLAN-CUST1-EXAMPLE"
  tags   = "CUST1,L2-WAN-01,DC-01"
  number = "666"
}

output "debug" {
  value = device42_ipam_vlan.example
}

## Will check if VLAN exists based on tags and number and import that into state if so. Else create a new VLAN.
## Limitation on filtering here is on API side. Unable to filter by name/vlan_id as of 27/8/21.

resource "device42_ipam_vlan" "example2" {
  check_if_exists = true

  name       = "VLAN-CUST2-EXAMPLE"
  tags_exist = "TEST,TEST2,CUST2"           # used for matching
  tags       = "TEST,TEST2,CUST2,TERRAFORM" # tags to update/create vlan with
}

output "debug" {
  value = device42_ipam_vlan.example2
}

## Will create a VLAN using the next sequentially available VLAN within provider range.
## VLANs that match the tags_range will be considered 'used'.
## Limitation on filtering here is on API side. Unable to filter by name/vlan_id as of 27/8/21.

resource "device42_ipam_vlan" "example3" {
  create_within_range = "666-766"

  name       = "VLAN-CUST3-EXAMPLE"
  tags_range = "L2-WAN-01,DC-01"
  tags       = "L2-WAN-01,DC-01,CUST3"
}

output "debug" {
  value = device42_ipam_vlan.example3
}