// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package init_shutdown_script

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

var _ list.ListResource = &InitShutdownScriptListResource{}
var _ list.ListResourceWithConfigure = &InitShutdownScriptListResource{}

// InitShutdownScriptListResource implements the truenas_init_shutdown_script list resource (terraform query).
type InitShutdownScriptListResource struct{ client *client.Client }

// NewListResource returns a new InitShutdownScriptListResource.
func NewListResource() list.ListResource { return &InitShutdownScriptListResource{} }

func (l *InitShutdownScriptListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_init_shutdown_script"
}

func (l *InitShutdownScriptListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *InitShutdownScriptListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *InitShutdownScriptListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "initshutdownscript.query", req, stream, l.mapRow)
}

func (l *InitShutdownScriptListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api initShutdownScriptAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse init/shutdown script row", err.Error())
		return res
	}
	// This resource has no user-supplied name field; combine type and
	// trigger point for a readable label.
	res.DisplayName = fmt.Sprintf("%s @ %s", api.Type, api.When)
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m InitShutdownScriptModel
		responseToModel(&api, &m)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
