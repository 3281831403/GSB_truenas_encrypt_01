// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package filesystem_acl

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestFilesystemAclIdentitySchema verifies identity-based import support.
// This resource has NO list.go: filesystem.getacl/setacl operate on an
// arbitrary, caller-supplied filesystem path, and there is no "list every
// ACL-managed path" query method to enumerate instances with.
func TestFilesystemAclIdentitySchema(t *testing.T) {
	r := &FilesystemAclResource{}
	resp := &resource.IdentitySchemaResponse{}
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, resp)
	if _, ok := resp.IdentitySchema.Attributes["id"]; !ok {
		t.Error("identity schema missing 'id' attribute")
	}
}
