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