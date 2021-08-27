resource "device42_ipam_subnet" "example" {
  name      = "EXAMPLE"
  tags      = "EXAMPLE"
  mask_bits = "29"
  network   = "10.25.0.0"
}

resource "device42_ipam_ip" "example" {
  subnet_id = device42_ipam_subnet.example.subnet_id
  ipaddress = "10.25.0.1"
  notes     = "server1.example.com"
}

output "example" {
  value = device42_ipam_ip.example
}

resource "device42_ipam_ip" "example2" {
  subnet_id  = device42_ipam_subnet.example.subnet_id
  suggest_ip = true
  notes      = "server2.example2.com"
}

output "example2" {
  value = device42_ipam_ip.example2
}