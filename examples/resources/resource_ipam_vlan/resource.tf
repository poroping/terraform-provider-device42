resource "device42_ipam_vlan" "example" {
    check_if_exists = true
    create_within_range = "1-106"

    name = "API-TEST23"
    tags = "TEST,TEST2,TEST3"
    number = "666"
}

output "debug" {
    value = device42_ipam_vlan.example
}

resource "device42_ipam_vlan" "example2" {
    create_within_range = "100-199"

    name = "API-TEST23"
    tags = "TEST,TEST2,TEST3"
}

output "debug2" {
    value = device42_ipam_vlan.example2
}