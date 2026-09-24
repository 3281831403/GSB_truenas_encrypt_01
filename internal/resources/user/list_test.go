// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestListMetadata(t *testing.T) {
	l := &UserListResource{}
	req := resource.MetadataRequest{ProviderTypeName: "truenas"}
	resp := &resource.MetadataResponse{}
	l.Metadata(context.Background(), req, resp)
	if resp.TypeName != "truenas_user" {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, "truenas_user")
	}
}

func TestUserIdentitySchema(t *testing.T) {
	r := &UserResource{}
	resp := &resource.IdentitySchemaResponse{}
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, resp)
	if _, ok := resp.IdentitySchema.Attributes["id"]; !ok {
		t.Error("identity schema missing 'id' attribute")
	}
}
