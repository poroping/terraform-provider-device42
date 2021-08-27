---
page_title: "device42 Provider"
subcategory: ""
description: |-
  
---

# device42 Provider

Butchered some old SDK into this provider to replace some python scripts.

Open to pull requests / requests for other resources at the github repo.

## Example Usage

```terraform
provider "device42" {
  host        = "d42.example.com"
  username    = "terraform"
  password    = "superpassword"
  insecure    = true
}
```

## Schema
