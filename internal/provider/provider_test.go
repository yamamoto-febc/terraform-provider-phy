package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"phy": func() (tfprotov6.ProviderServer, error) {
		return tfsdk.NewProtocol6Server(New()), nil
	},
}

func testAccPreCheck(t *testing.T) {
	testEnvVars := []string{"SAKURACLOUD_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN_SECRET"}
	for _, env := range testEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}