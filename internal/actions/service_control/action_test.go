// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package service_control

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMetadata(t *testing.T) {
	a := &Action{}
	resp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{ProviderTypeName: "truenas"}, resp)
	if resp.TypeName != "truenas_service_control" {
		t.Errorf("TypeName = %q", resp.TypeName)
	}
}

func TestBuildParams_VerbThenService(t *testing.T) {
	got := buildParams(model{Service: types.StringValue("cifs"), Verb: types.StringValue("RESTART")})
	if len(got) != 2 || got[0] != "RESTART" || got[1] != "cifs" {
		t.Errorf("buildParams = %#v, want [RESTART cifs]", got)
	}
}

func TestSchema_VerbValidated(t *testing.T) {
	a := &Action{}
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	attr, ok := resp.Schema.Attributes["verb"].(interface {
		StringValidators() []validator.String
	})
	if !ok {
		t.Fatal("verb attribute does not expose string validators")
	}
	if len(attr.StringValidators()) == 0 {
		t.Error("verb should carry a OneOf validator")
	}
}

func TestVersionGate(t *testing.T) {
	for _, v := range []string{"25.10.3.1", "24.10", ""} {
		if !versionGateDiagnostics(v).HasError() {
			t.Errorf("version %q should be rejected", v)
		}
	}
	for _, v := range []string{"26.0.0", "26.1.0", "27.0.0"} {
		if versionGateDiagnostics(v).HasError() {
			t.Errorf("version %q should be accepted", v)
		}
	}
}

func TestSchema_HasWait(t *testing.T) {
	a := &Action{}
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	if _, ok := resp.Schema.Attributes["wait"]; !ok {
		t.Fatal("schema missing 'wait' attribute")
	}
}
