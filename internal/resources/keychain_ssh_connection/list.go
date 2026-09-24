// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package keychain_ssh_connection

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

var _ list.ListResource = &KeychainSSHConnectionListResource{}
var _ list.ListResourceWithConfigure = &KeychainSSHConnectionListResource{}

// KeychainSSHConnectionListResource implements the
// truenas_keychain_ssh_connection list resource (terraform query).
type KeychainSSHConnectionListResource struct{ client *client.Client }

// NewListResource returns a new KeychainSSHConnectionListResource.
func NewListResource() list.ListResource { return &KeychainSSHConnectionListResource{} }

func (l *KeychainSSHConnectionListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keychain_ssh_connection"
}

func (l *KeychainSSHConnectionListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *KeychainSSHConnectionListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// List queries keychaincredential.query filtered to type=SSH_CREDENTIALS
// (keychainCredentialType, defined in model.go): the same query method
// backs truenas_keychain_ssh_keypair (type=SSH_KEY_PAIR), and without this
// filter each list would also return the other kind's rows.
func (l *KeychainSSHConnectionListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	filter := [][]any{{"type", "=", keychainCredentialType}}
	listing.StreamCollectionFiltered(ctx, l.client, "keychaincredential.query", filter, req, stream, l.mapRow)
}

func (l *KeychainSSHConnectionListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api keychainSSHConnectionAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse keychain SSH connection row", err.Error())
		return res
	}
	res.DisplayName = api.Name
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m KeychainSSHConnectionModel
		res.Diagnostics.Append(responseToModel(&api, &m)...)
		if res.Diagnostics.HasError() {
			return res
		}
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
