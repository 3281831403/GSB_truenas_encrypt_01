// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package boot_environment

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

var _ list.ListResource = &BootEnvironmentListResource{}
var _ list.ListResourceWithConfigure = &BootEnvironmentListResource{}

// BootEnvironmentListResource implements the truenas_boot_environment list
// resource (terraform query).
type BootEnvironmentListResource struct{ client *client.Client }

// NewListResource returns a new BootEnvironmentListResource.
func NewListResource() list.ListResource { return &BootEnvironmentListResource{} }

func (l *BootEnvironmentListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_boot_environment"
}

func (l *BootEnvironmentListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *BootEnvironmentListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *BootEnvironmentListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "boot.environment.query", req, stream, l.mapRow)
}

// mapRow builds a ListResult from one boot.environment.query row. Note:
// "source" (the clone origin a boot environment was created from) is never
// present in boot.environment.query's response — see bootEnvAPI's doc
// comment in model.go — so a listed row's full resource representation
// leaves "source" null, exactly like this resource's own ImportState does
// today (it never learns the clone origin either).
func (l *BootEnvironmentListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api bootEnvAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse boot environment row", err.Error())
		return res
	}
	res.DisplayName = api.ID
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m BootEnvironmentModel
		responseToModel(&api, &m)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
