// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package certificate

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestListMetadata(t *testing.T) {
	l := &CertificateListResource{}
	req := resource.MetadataRequest{ProviderTypeName: "truenas"}
	resp := &resource.MetadataResponse{}
	l.Metadata(context.Background(), req, resp)
	if resp.TypeName != "truenas_certificate" {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, "truenas_certificate")
	}
}

func TestCertificateIdentitySchema(t *testing.T) {
	r := &CertificateResource{}
	resp := &resource.IdentitySchemaResponse{}
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, resp)
	if _, ok := resp.IdentitySchema.Attributes["id"]; !ok {
		t.Error("identity schema missing 'id' attribute")
	}
}
