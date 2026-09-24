// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package vm_device

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/truenas/terraform-provider-truenas/internal/client"
	"github.com/truenas/terraform-provider-truenas/internal/listing"
)

var _ list.ListResource = &VMDeviceListResource{}
var _ list.ListResourceWithConfigure = &VMDeviceListResource{}

// VMDeviceListResource implements the truenas_vm_device list resource (terraform query).
type VMDeviceListResource struct{ client *client.Client }

// NewListResource returns a new VMDeviceListResource.
func NewListResource() list.ListResource { return &VMDeviceListResource{} }

func (l *VMDeviceListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_device"
}

func (l *VMDeviceListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *VMDeviceListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *VMDeviceListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "vm.device.query", req, stream, l.mapRow)
}

func (l *VMDeviceListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api deviceAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse VM device row", err.Error())
		return res
	}
	res.DisplayName = fmt.Sprintf("vm %d / device %d", api.VM, api.ID)
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m VMDeviceModel
		responseToModel(&api, &m)
		// responseToModel deliberately leaves Attributes untouched; populate
		// it from the API response here, mirroring ImportState.
		attrJSON, diags := apiAttributesJSON(&api)
		res.Diagnostics.Append(diags...)
		if res.Diagnostics.HasError() {
			return res
		}
		m.Attributes = types.StringValue(attrJSON)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
