package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sacloud/phy-go"
	v1 "github.com/sacloud/phy-go/apis/v1"
)

type dataSourceServerType struct{}

// GetSchema .
func (r dataSourceServerType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "ID of the resource.",
				Type:        types.StringType,
				Computed:    true,
			},
			"nickname": {
				Description: "NickName of the server.",
				Type:        types.StringType,
				Computed:    true,
			},
			"filter": {
				Description: "The free word filter of the server.",
				Type:        types.StringType,
				Required:    true,
			},
		},
	}, nil
}

// NewDataSource resource instance
func (r dataSourceServerType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, ok := p.(*provider)
	if !ok {
		return nil, diag.Diagnostics{errorConvertingProvider(p)}
	}
	return serverDataSource{client: provider.client}, nil
}

type serverDataSource struct {
	client *phy.Client
}

type serverData struct {
	Id       *string `tfsdk:"id"`
	NickName *string `tfsdk:"nickname"`

	Filter string `tfsdk:"filter"`
}

// Read resource information
func (r serverDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data serverData
	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	filter := v1.FreeWordFilter{data.Filter}
	found, err := phy.NewServerOp(r.client).List(ctx, &v1.ListServersParams{
		FreeWord: &filter,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Getting Servers", err.Error())
		return
	}
	if len(found.Servers) == 0 {
		resp.Diagnostics.AddError("Error Getting Server", fmt.Sprintf("No Server where %s", filter))
		return
	}

	data = serverData{
		Id:       &found.Servers[0].ServerId,
		NickName: &found.Servers[0].Service.Nickname,
		Filter:   data.Filter,
	}
	resp.Diagnostics = resp.State.Set(ctx, &data)
}
