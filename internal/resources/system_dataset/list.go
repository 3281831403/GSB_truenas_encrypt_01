// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package system_dataset

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

var _ list.ListResource = &SystemDatasetListResource{}
var _ list.ListResourceWithConfigure = &SystemDatasetListResource{}

// SystemDatasetListResource implements the truenas_system_dataset list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type SystemDatasetListResource struct{ client *client.Client }

// NewListResource returns a new SystemDatasetListResource.
func NewListResource() list.ListResource { return &SystemDatasetListResource{} }

func (l *SystemDatasetListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_dataset"
}

func (l *SystemDatasetListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *SystemDatasetListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *SystemDatasetListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "systemdataset.config", req, stream, l.mapOne)
}

func (l *SystemDatasetListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api systemDatasetAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse system_dataset", err.Error())
		return res
	}
	res.DisplayName = "system_dataset"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, systemDatasetResourceID)...)
	if req.IncludeResource {
		var m SystemDatasetModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
