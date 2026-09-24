// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package docker_config

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

var _ list.ListResource = &DockerConfigListResource{}
var _ list.ListResourceWithConfigure = &DockerConfigListResource{}

// DockerConfigListResource implements the truenas_docker_config list
// resource (terraform query). It is a singleton: List always emits exactly
// one result.
type DockerConfigListResource struct{ client *client.Client }

// NewListResource returns a new DockerConfigListResource.
func NewListResource() list.ListResource { return &DockerConfigListResource{} }

func (l *DockerConfigListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_docker_config"
}

func (l *DockerConfigListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *DockerConfigListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *DockerConfigListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "docker.config", req, stream, l.mapOne)
}

func (l *DockerConfigListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api dockerConfigAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse docker_config", err.Error())
		return res
	}
	res.DisplayName = "docker_config"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, dockerConfigResourceID)...)
	if req.IncludeResource {
		var m DockerConfigModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
