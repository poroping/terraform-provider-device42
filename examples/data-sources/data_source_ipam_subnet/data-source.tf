data "device42_ipam_subnet" "example" {
  subnet_id = "1321"
}

output "example" {
  value = data.device42_ipam_subnet.example
}