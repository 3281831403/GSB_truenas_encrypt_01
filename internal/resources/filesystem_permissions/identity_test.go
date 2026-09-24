// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package filesystem_permissions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestFilesystemPermissionsIdentitySchema verifies identity-based import
// support. This resource has NO list.go: filesystem.setperm/stat operate on
// an arbitrary, caller-supplied filesystem path, and there is no "list every
// managed path" query method to enumerate instances with.
func TestFilesystemPermissionsIdentitySchema(t *testing.T) {
	r := &FilesystemPermissionsResource{}
	resp := &resource.IdentitySchemaResponse{}
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, resp)
	if _, ok := resp.IdentitySchema.Attributes["id"]; !ok {
		t.Error("identity schema missing 'id' attribute")
	}
}
