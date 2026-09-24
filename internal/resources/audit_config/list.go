// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package audit_config

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

var _ list.ListResource = &AuditConfigListResource{}
var _ list.ListResourceWithConfigure = &AuditConfigListResource{}

// AuditConfigListResource implements the truenas_audit_config list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type AuditConfigListResource struct{ client *client.Client }

// NewListResource returns a new AuditConfigListResource.
func NewListResource() list.ListResource { return &AuditConfigListResource{} }

func (l *AuditConfigListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_config"
}

func (l *AuditConfigListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *AuditConfigListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *AuditConfigListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "audit.config", req, stream, l.mapOne)
}

func (l *AuditConfigListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api auditConfigAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse audit_config", err.Error())
		return res
	}
	res.DisplayName = "audit_config"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, auditConfigResourceID)...)
	if req.IncludeResource {
		var m AuditConfigModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
