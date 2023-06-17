package main

import (
	"context"
	"terraform-provider-google-tag-manager/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/mirefly/google-tag-manager",
	})
}
