terraform {
  required_providers {
    jamfautoupdate = {
      source  = "Jamf-Concepts/jamfautoupdate"
      version = "~> 1.0"
    }
  }
}

# Configure the provider with the definitions URL
provider "jamfautoupdate" {
  definitions_url = var.definitions_url
}

# Data source to fetch all titles
data "jamfautoupdate_titles" "all" {}

# Output example to show all title names
output "all_title_names" {
  value = [
    for title in data.jamfautoupdate_titles.all.titles :
    title.title_name
  ]
}

# Data source to fetch title details for a specific list
data "jamfautoupdate_titles" "specific" {
  title_names = [
    "GoogleChrome",
    "MicrosoftOutlook"
  ]
}

# Save icons for all titles
resource "local_file" "title_icons" {
  for_each = {
    for title in data.jamfautoupdate_titles.specific.titles :
    title.title_name => title
    if title.icon_base64 != null
  }
  content_base64 = each.value.icon_base64
  filename       = "${path.module}/icons/${each.value.title_name}.png"
}

# Save multiple profiles with dynamic names for the first title
locals {
  profile_types = {
    notifications = data.jamfautoupdate_titles.specific.titles[0].notifications_profile
    pppcp         = data.jamfautoupdate_titles.specific.titles[0].pppcp_profile
    screen        = data.jamfautoupdate_titles.specific.titles[0].screen_recording_profile
  }
}

resource "local_file" "first_title_profiles" {
  for_each = {
    for type, content in local.profile_types :
    type => content
    if content != null
  }

  content_base64 = each.value
  filename       = "${path.module}/profiles/first_title_${each.key}.mobileconfig"
}

# Output which profiles were available and saved
output "available_profiles" {
  value = {
    for type, content in local.profile_types :
    type => content != null
  }
}
