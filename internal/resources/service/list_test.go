// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestListMetadata(t *testing.T) {
	l := &ServiceListResource{}
	req := resource.MetadataRequest{ProviderTypeName: "truenas"}
	resp := &resource.MetadataResponse{}
	l.Metadata(context.Background(), req, resp)
	if resp.TypeName != "truenas_service" {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, "truenas_service")
	}
}
