package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/sacloud/phy-go"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	client *phy.Client
}

// GetSchema .
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"token": {
				Type:     types.StringType,
				Optional: true,
			},
			"secret": {
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
			},
			"api_root_url": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type providerData struct {
	Token      types.String `tfsdk:"token"`
	Secret     types.String `tfsdk:"secret"`
	APIRootURL types.String `tfsdk:"api_root_url"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "token")
		return
	}
	if config.Secret.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "secret")
		return
	}

	var token string
	if config.Token.Null {
		token = os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	} else {
		token = config.Token.Value
	}
	if token == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find token",
			"Token cannot be an empty string",
		)
		return
	}

	var secret string
	if config.Secret.Null {
		secret = os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")
	} else {
		secret = config.Token.Value
	}
	if secret == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find secret",
			"Secret cannot be an empty string",
		)
		return
	}

	var apiRootURL string
	if config.APIRootURL.Null {
		apiRootURL = os.Getenv("SAKURACLOUD_PHY_API_ROOT_URL")
	} else {
		apiRootURL = config.Token.Value
	}
	if apiRootURL == "" {
		apiRootURL = phy.DefaultAPIRootURL
	}

	p.client = &phy.Client{
		Token:      token,
		Secret:     secret,
		APIRootURL: apiRootURL,
		Trace:      os.Getenv("SAKURACLOUD_TRACE") != "",
	}
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"phy_server": dataSourceServerType{},
	}, nil
}

func addCannotInterpolateInProviderBlockError(resp *tfsdk.ConfigureProviderResponse, attr string) {
	resp.Diagnostics.AddAttributeError(
		tftypes.NewAttributePath().WithAttributeName(attr),
		"Can't interpolate into provider block",
		"Interpolating that value into the provider block doesn't give the provider enough information to run. Try hard-coding the value, instead.",
	)
}

func errorConvertingProvider(typ interface{}) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Error converting provider", fmt.Sprintf("An unexpected error was encountered converting the provider. This is always a bug in the provider.\n\nType: %T", typ))
}
