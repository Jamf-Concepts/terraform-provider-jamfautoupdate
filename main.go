// Copyright 2025 Jamf Software LLC.

package main

import (
	"context"
	"log"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/Jamf-Concepts/jamfautoupdate",
	}

	err := providerserver.Serve(context.Background(), provider.New, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
