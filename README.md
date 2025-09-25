# Jamf Auto Update Terraform Provider

Provides a data source for Jamf Auto Update titles using the Definitions API

## Requirements

* Requires network access to the Definitions API (e.g. via Jamf Trust when run on macOS)

## Manual Installation

### Step 1: Download the Release Zip

Pick your platform and architecture from the [latest releases](https://github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/releases/latest) page:

* If running on an Apple Silicon Mac, download `...darwin_arm64.zip`
* If running on an Intel Mac, download `...darwin_amd64.zip`

---

### Step 2: Extract to your local Terraform Plugin Directory

The plugin has to be extracted to the correct location in order for Terraform to find and use it. We also have to remove the Quarantine attribute, or it will not load.

#### Example (macOS arm64)

```bash
cd ~/Downloads
mkdir -p ~/.terraform.d/plugins/local/Jamf-Concepts/jamfautoupdate/1.4.4/darwin_arm64
unzip terraform-provider-jamfautoupdate_1.4.4_darwin_arm64.zip -d ~/.terraform.d/plugins/local/Jamf-Concepts/jamfautoupdate/1.4.4/darwin_arm64
xattr -r -d com.apple.quarantine ~/.terraform.d/plugins
```

This will result in:

```
~/.terraform.d/plugins/
└── local/
    └── jamf/
        └── jamfautoupdate/
            └── 1.4.4/
                └── darwin_arm64/
                    └── terraform-provider-jamfautoupdate_v1.4.4
```

### Step 3: Set up a local file system mirror

Terraform needs to know it will be using a locally installed plugin, so this is how we tell it where to find them.

Create the file: `~/.terraform.d/terraform.tfrc`

* `nano ~/.terraform.d/terraform.tfrc`

It must contain the following contents (copy/paste into the file):

```hcl
provider_installation {
  filesystem_mirror {
    path    = ~"/.terraform.d/plugins"
    include = ["local/Jamf-Concepts/*"]
  }
  
  direct {
    exclude = ["local/Jamf-Concepts/*"]
  }
}

cli_config {
  parallelism = 1
}
```

* In nano, use **ctrl+x** then enter **y** and **return** to save.

That's it! You're ready to use the provider!

---

## Using the Provider in your own Terraform Projects

In your `main.tf`:

```hcl
terraform {
  required_providers {
    jamfautoupdate = {
      source  = "local/Jamf-Concepts/jamfautoupdate"
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

Terraform will detect the manually-installed provider and you're good to go!

**Note: If you are updating this provider manually, you should also delete the lock file:**

```bash
rm .terraform.lock.hcl
terraform init
```

## Provider Configuration Reference and Example Usage

Refer to [this documentation](./docs).

## Included components

The following third party acknowledgements and licenses are incorporated by reference:

* [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) ([MPL](https://github.com/hashicorp/terraform-plugin-framework?tab=MPL-2.0-1-ov-file))
* [Terraform Plugin Log](https://github.com/hashicorp/terraform-plugin-log) ([MPL](https://github.com/hashicorp/terraform-plugin-log?tab=MPL-2.0-1-ov-file))
* [resize](https://github.com/nfnt/resize) ([ISC License](https://github.com/nfnt/resize/blob/master/LICENSE))
