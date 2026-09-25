// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

// Unit tests for configurable dataset encryption and lock state.
package dataset

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestDatasetCreateAPIPayloadEncryptionOn verifies that configuring
// encrypted = true produces a create payload that really enables
// encryption (and disables inherit_encryption, which pool.dataset.create
// rejects in combination with encryption=true).
func TestDatasetCreateAPIPayloadEncryptionOn(t *testing.T) {
	m := &DatasetModel{
		Name:      types.StringValue("tank/mydata"),
		Encrypted: types.BoolValue(true),
	}

	payload := m.apiPayload()

	if payload["encryption"] != true {
		t.Errorf("expected encryption=true in create payload, got %v", payload["encryption"])
	}
	if payload["inherit_encryption"] != false {
		t.Errorf("expected inherit_encryption=false when enabling encryption, got %v", payload["inherit_encryption"])
	}
}

// TestDatasetCreateAPIPayloadEncryptionOff verifies that configuring
// encrypted = false is sent explicitly (not silently dropped), so the
// created dataset's state matches the configuration.
func TestDatasetCreateAPIPayloadEncryptionOff(t *testing.T) {
	m := &DatasetModel{
		Name:      types.StringValue("tank/mydata"),
		Encrypted: types.BoolValue(false),
	}

	payload := m.apiPayload()

	if payload["encryption"] != false {
		t.Errorf("expected encryption=false in create payload, got %v", payload["encryption"])
	}
	if _, ok := payload["inherit_encryption"]; ok {
		t.Errorf("inherit_encryption must only be sent when enabling encryption, got %v", payload)
	}
}

// TestDatasetCreateAPIPayloadEncryptionUnset verifies that an unmanaged
// (null/unknown) encrypted attribute leaves the create payload free of
// encryption keys, preserving TrueNAS defaults.
func TestDatasetCreateAPIPayloadEncryptionUnset(t *testing.T) {
	for name, enc := range map[string]types.Bool{
		"null":    types.BoolNull(),
		"unknown": types.BoolUnknown(),
	} {
		m := &DatasetModel{
			Name:      types.StringValue("tank/mydata"),
			Encrypted: enc,
		}
		payload := m.apiPayload()
		if _, ok := payload["encryption"]; ok {
			t.Errorf("%s encrypted: create payload must not include \"encryption\", got %v", name, payload)
		}
		if _, ok := payload["inherit_encryption"]; ok {
			t.Errorf("%s encrypted: create payload must not include \"inherit_encryption\", got %v", name, payload)
		}
	}
}

// TestDatasetUpdateAPIPayloadOmitsEncryption verifies that encryption keys
// never reach pool.dataset.update, which rejects them as create-only.
// Toggling encryption on an existing dataset converges via replacement.
func TestDatasetUpdateAPIPayloadOmitsEncryption(t *testing.T) {
	for _, enc := range []bool{true, false} {
		m := &DatasetModel{
			Name:      types.StringValue("tank/mydata"),
			Encrypted: types.BoolValue(enc),
		}
		payload := m.updateAPIPayload()
		if _, ok := payload["encryption"]; ok {
			t.Errorf("encrypted=%v: update payload must not include \"encryption\", got %v", enc, payload)
		}
		if _, ok := payload["inherit_encryption"]; ok {
			t.Errorf("encrypted=%v: update payload must not include \"inherit_encryption\", got %v", enc, payload)
		}
	}
}

// TestDatasetAPIResponseParsesLocked verifies that the locked field
// returned by pool.dataset.get_instance is decoded.
func TestDatasetAPIResponseParsesLocked(t *testing.T) {
	raw := []byte(`{"name":"tank/mydata","type":"FILESYSTEM","mountpoint":"/mnt/tank/mydata","encrypted":true,"locked":true,"pool":"tank"}`)

	var apiResp apiResponse
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		t.Fatalf("unmarshal apiResponse: %v", err)
	}
	if !apiResp.Locked {
		t.Error("expected Locked=true from API response")
	}
	if !apiResp.Encrypted {
		t.Error("expected Encrypted=true from API response")
	}
}

// TestDatasetResponseToModelSetsLockedAndEncrypted verifies that the
// actual encryption and lock state read back from the API lands in state,
// so drift (manual changes) is detected and can be converged.
func TestDatasetResponseToModelSetsLockedAndEncrypted(t *testing.T) {
	r := &DatasetResource{}

	for _, tc := range []struct {
		name      string
		encrypted bool
		locked    bool
	}{
		{"encrypted and locked", true, true},
		{"encrypted and unlocked", true, false},
		{"unencrypted and unlocked", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			apiResp := apiResponse{
				Name:      "tank/mydata",
				Encrypted: tc.encrypted,
				Locked:    tc.locked,
			}
			var m DatasetModel
			if diags := r.responseToModel(&apiResp, &m); diags.HasError() {
				t.Fatalf("responseToModel returned errors: %v", diags)
			}
			if m.Encrypted.ValueBool() != tc.encrypted {
				t.Errorf("Encrypted = %v, want %v", m.Encrypted.ValueBool(), tc.encrypted)
			}
			if m.Locked.ValueBool() != tc.locked {
				t.Errorf("Locked = %v, want %v", m.Locked.ValueBool(), tc.locked)
			}
		})
	}
}

// TestLockAction verifies the (configured, actual) lock-state mapping:
// lock when configured locked but actually unlocked, unlock when configured
// unlocked but actually locked, and no API call when they already agree.
func TestLockAction(t *testing.T) {
	t.Run("lock when configured locked and actual unlocked", func(t *testing.T) {
		method, params, ok := lockAction("tank/mydata", true, false)
		if !ok {
			t.Fatal("expected a lock action")
		}
		if method != "pool.dataset.lock" {
			t.Errorf("method = %q, want pool.dataset.lock", method)
		}
		if len(params) != 2 || params[0] != "tank/mydata" {
			t.Fatalf("unexpected lock params: %v", params)
		}
		opts, ok := params[1].(map[string]any)
		if !ok || opts["force_umount"] != true {
			t.Errorf("expected force_umount=true lock options, got %v", params[1])
		}
	})

	t.Run("unlock when configured unlocked and actual locked", func(t *testing.T) {
		method, params, ok := lockAction("tank/mydata", false, true)
		if !ok {
			t.Fatal("expected an unlock action")
		}
		if method != "pool.dataset.unlock" {
			t.Errorf("method = %q, want pool.dataset.unlock", method)
		}
		if len(params) != 2 || params[0] != "tank/mydata" {
			t.Fatalf("unexpected unlock params: %v", params)
		}
		opts, ok := params[1].(map[string]any)
		if !ok {
			t.Fatalf("expected unlock options map, got %T", params[1])
		}
		datasets, ok := opts["datasets"].([]any)
		if !ok || len(datasets) != 1 {
			t.Fatalf("expected one unlock dataset spec, got %v", opts["datasets"])
		}
		spec, ok := datasets[0].(map[string]any)
		if !ok || spec["name"] != "tank/mydata" {
			t.Errorf("expected unlock spec for tank/mydata, got %v", datasets[0])
		}
	})

	t.Run("no action when states agree", func(t *testing.T) {
		for _, tc := range []struct {
			configured bool
			actual     bool
		}{
			{true, true},
			{false, false},
		} {
			if method, params, ok := lockAction("tank/mydata", tc.configured, tc.actual); ok {
				t.Errorf("configured=%v actual=%v: expected no action, got %s %v",
					tc.configured, tc.actual, method, params)
			}
		}
	})
}

// TestDatasetSchemaEncryptionLockAttributes verifies the schema exposes
// encrypted and locked as configurable (Optional) attributes while keeping
// them Computed for unmanaged configurations.
func TestDatasetSchemaEncryptionLockAttributes(t *testing.T) {
	s := resourceSchema()

	enc, ok := s.Attributes["encrypted"]
	if !ok {
		t.Fatal("resourceSchema missing attribute \"encrypted\"")
	}
	encAttr, ok := enc.(interface{ IsOptional() bool })
	if !ok || !encAttr.IsOptional() {
		t.Error("encrypted must be Optional so configuration can enable/disable encryption")
	}
	encComputed, ok := enc.(interface{ IsComputed() bool })
	if !ok || !encComputed.IsComputed() {
		t.Error("encrypted must stay Computed for unmanaged configurations")
	}

	locked, ok := s.Attributes["locked"]
	if !ok {
		t.Fatal("resourceSchema missing attribute \"locked\"")
	}
	lockedAttr, ok := locked.(interface{ IsOptional() bool })
	if !ok || !lockedAttr.IsOptional() {
		t.Error("locked must be Optional so configuration can lock/unlock the dataset")
	}
	lockedComputed, ok := locked.(interface{ IsComputed() bool })
	if !ok || !lockedComputed.IsComputed() {
		t.Error("locked must stay Computed for unmanaged configurations")
	}
}
