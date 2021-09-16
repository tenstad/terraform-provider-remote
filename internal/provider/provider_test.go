package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"remote": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
	"remotehost": func() (*schema.Provider, error) {
		provider := New("dev")()
		configureProvider := provider.ConfigureContextFunc
		provider.ConfigureContextFunc = func(c context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			rd.Set("conn", []interface{}{
				map[string]interface{}{
					"host":     "remotehost",
					"user":     "root",
					"password": "password",
					"port":     22,
				},
			})
			return configureProvider(c, rd)
		}
		return provider, nil
	},
	"remotehost2": func() (*schema.Provider, error) {
		provider := New("dev")()
		configureProvider := provider.ConfigureContextFunc
		provider.ConfigureContextFunc = func(c context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			rd.Set("conn", []interface{}{
				map[string]interface{}{
					"host":     "remotehost2",
					"user":     "root",
					"password": "password",
					"port":     22,
				},
			})
			return configureProvider(c, rd)
		}
		return provider, nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
