package main

import (
	"context"
	"terraform-provider-google-tag-manager/gtm"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), gtm.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/mirefly/google-tag-manager",
	})
}
