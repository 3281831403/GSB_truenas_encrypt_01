// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package container_device

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

var _ list.ListResource = &ContainerDeviceListResource{}
var _ list.ListResourceWithConfigure = &ContainerDeviceListResource{}

// ContainerDeviceListResource implements the truenas_container_device list resource (terraform query).
type ContainerDeviceListResource struct{ client *client.Client }

// NewListResource returns a new ContainerDeviceListResource.
func NewListResource() list.ListResource { return &ContainerDeviceListResource{} }

func (l *ContainerDeviceListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_device"
}

func (l *ContainerDeviceListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *ContainerDeviceListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *ContainerDeviceListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "container.device.query", req, stream, l.mapRow)
}

func (l *ContainerDeviceListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api containerDeviceAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse container device row", err.Error())
		return res
	}
	res.DisplayName = fmt.Sprintf("container %d / device %d", api.Container, api.ID)
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m ContainerDeviceModel
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
