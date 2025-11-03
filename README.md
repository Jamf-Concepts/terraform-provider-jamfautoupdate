# Jamf Auto Update Terraform Provider

Provides a data source for Jamf Auto Update titles using the Definitions API

> This provider is intended for use with internal Jamf projects. It is made public to aid in installation/project initialization and serve as an educational aid to those wishing to develop their own Terraform providers.

## Requirements

* Requires network access to the Definitions API (e.g. via Jamf Trust when run on macOS or operating from a trusted IP range when used as part of a CI/CD pipeline)

## Using the Provider in projects

In your `main.tf`:

```hcl
terraform {
  required_providers {
    jamfautoupdate = {
      source  = "Jamf-Concepts/jamfautoupdate"
      version = "~> 1.0"
    }
  }
}

provider "jamfautoupdate" {
  # your configuration here
}
```

---

### Initialize Terraform

```bash
terraform init
```

Terraform will install the provider and you're good to go!

## Provider Configuration Reference and Example Usage

Refer to [the documentation](https://registry.terraform.io/providers/Jamf-Concepts/jamfautoupdate/latest/docs).

## Included components

The following third party acknowledgements and licenses are incorporated by reference:

* [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) ([MPL](https://github.com/hashicorp/terraform-plugin-framework?tab=MPL-2.0-1-ov-file))
* [Terraform Plugin Log](https://github.com/hashicorp/terraform-plugin-log) ([MPL](https://github.com/hashicorp/terraform-plugin-log?tab=MPL-2.0-1-ov-file))
