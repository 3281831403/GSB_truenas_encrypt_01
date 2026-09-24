// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package truecommand_config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/truenas/terraform-provider-truenas/internal/client"
	"github.com/truenas/terraform-provider-truenas/internal/listing"
)

var _ list.ListResource = &TrueCommandConfigListResource{}
var _ list.ListResourceWithConfigure = &TrueCommandConfigListResource{}

// TrueCommandConfigListResource implements the truenas_truecommand_config list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type TrueCommandConfigListResource struct{ client *client.Client }

// NewListResource returns a new TrueCommandConfigListResource.
func NewListResource() list.ListResource { return &TrueCommandConfigListResource{} }

func (l *TrueCommandConfigListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_truecommand_config"
}

func (l *TrueCommandConfigListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *TrueCommandConfigListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *client.Client, got %T", req.ProviderData))
		return
	}
	l.client = c
}

func (l *TrueCommandConfigListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "truecommand.config", req, stream, l.mapOne)
}

func (l *TrueCommandConfigListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api truecommandConfigAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse truecommand_config", err.Error())
		return res
	}
	res.DisplayName = "truecommand_config"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, trueCommandConfigResourceID)...)
	if req.IncludeResource {
		var m TrueCommandConfigModel
		res.Diagnostics.Append(responseToModel(&api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
