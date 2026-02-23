// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories returns the provider factories for acceptance tests.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jamfautoupdate": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccPreCheck validates that required environment variables are set.
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("JAMF_AUTO_UPDATE_DEFINITIONS_URL") == "" && os.Getenv("JAMF_AUTO_UPDATE_DEFINITIONS_FILE") == "" {
		t.Skip("JAMF_AUTO_UPDATE_DEFINITIONS_URL or JAMF_AUTO_UPDATE_DEFINITIONS_FILE must be set for acceptance tests")
	}
}

func TestAccTitlesDataSource_FetchSpecific(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "jamfautoupdate_titles" "test" {
  title_names = ["GoogleChrome", "1Password"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jamfautoupdate_titles.test", "titles.#", "2"),
				),
			},
		},
	})
}

func TestAccTitlesDataSource_FetchAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "jamfautoupdate_titles" "all" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jamfautoupdate_titles.all", "titles.#"),
				),
			},
		},
	})
}

func TestAccTitlesDataSource_VerifyAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "jamfautoupdate_titles" "test" {
  title_names = ["GoogleChrome"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jamfautoupdate_titles.test", "titles.#", "1"),
					resource.TestCheckResourceAttr("data.jamfautoupdate_titles.test", "titles.0.title_name", "GoogleChrome"),
					resource.TestCheckResourceAttrSet("data.jamfautoupdate_titles.test", "titles.0.title_display_name"),
					resource.TestCheckResourceAttrSet("data.jamfautoupdate_titles.test", "titles.0.title_version"),
					resource.TestCheckResourceAttrSet("data.jamfautoupdate_titles.test", "titles.0.icon_base64"),
					resource.TestCheckResourceAttrSet("data.jamfautoupdate_titles.test", "titles.0.uninstall_icon_base64"),
					resource.TestCheckResourceAttrSet("data.jamfautoupdate_titles.test", "titles.0.app_bundle_id"),
				),
			},
		},
	})
}

func TestAccTitlesDataSource_EmptyList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "jamfautoupdate_titles" "empty" {
  title_names = []
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jamfautoupdate_titles.empty", "titles.#", "0"),
				),
			},
		},
	})
}

func TestAccTitlesDataSource_InvalidTitle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "jamfautoupdate_titles" "bad" {
  title_names = ["NonExistentTitle12345"]
}`,
				ExpectError: regexp.MustCompile(`(titles do not exist|status code: 404)`),
			},
		},
	})
}

func TestAccProviderConfigure_BothSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "jamfautoupdate" {
  definitions_url  = "https://example.com"
  definitions_file = "/tmp/test.json"
}

data "jamfautoupdate_titles" "test" {
  title_names = []
}`,
				ExpectError: regexp.MustCompile(`Exactly one of definitions_url or definitions_file must be set`),
			},
		},
	})
}
