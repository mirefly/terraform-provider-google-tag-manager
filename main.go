package main

import (
	"context"
	"terraform-provider-google-tag-manager/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/mirefly/google-tag-manager",
	})
}
