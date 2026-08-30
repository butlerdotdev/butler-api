/*
Copyright 2026 The Butler Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"encoding/json"
	"strings"
	"testing"
)

// ptrInt32 is a small helper for building explicit replica values in tests.
func ptrInt32(v int32) *int32 { return &v }

// TestControlPlaneReplicasOmitEmpty verifies that an omitted control plane
// replicas value serializes as absent (nil), so the hosted control plane
// provider owns the default. An explicit zero must serialize as 0 (a distinct,
// meaningful value), never be confused with omission.
func TestControlPlaneReplicasOmitEmpty(t *testing.T) {
	tests := []struct {
		name        string
		replicas    *int32
		wantPresent bool
		wantValue   string
	}{
		{name: "omitted is absent", replicas: nil, wantPresent: false},
		{name: "explicit zero is present", replicas: ptrInt32(0), wantPresent: true, wantValue: "0"},
		{name: "explicit one is present", replicas: ptrInt32(1), wantPresent: true, wantValue: "1"},
		{name: "explicit two is present", replicas: ptrInt32(2), wantPresent: true, wantValue: "2"},
		{name: "explicit three is present", replicas: ptrInt32(3), wantPresent: true, wantValue: "3"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(ControlPlaneSpec{Replicas: tc.replicas})
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			got := string(b)
			has := strings.Contains(got, "\"replicas\"")
			if has != tc.wantPresent {
				t.Fatalf("replicas presence = %v, want %v (json=%s)", has, tc.wantPresent, got)
			}
			if tc.wantPresent && !strings.Contains(got, "\"replicas\":"+tc.wantValue) {
				t.Fatalf("replicas value not %s in json=%s", tc.wantValue, got)
			}
		})
	}
}

// TestControlPlaneReplicasRoundTrip verifies existing serialized objects remain
// safe: a stored explicit replicas value deserializes back to the same explicit
// pointer, and an object without the field deserializes to nil (provider owns
// the default). This guarantees the int32 -> *int32 change does not migrate any
// existing TenantCluster.
func TestControlPlaneReplicasRoundTrip(t *testing.T) {
	// Existing object that stored an explicit replicas: 1 stays an explicit 1.
	var withOne ControlPlaneSpec
	if err := json.Unmarshal([]byte(`{"replicas":1}`), &withOne); err != nil {
		t.Fatalf("unmarshal replicas:1: %v", err)
	}
	if withOne.Replicas == nil || *withOne.Replicas != 1 {
		t.Fatalf("stored replicas:1 did not round-trip to explicit 1, got %v", withOne.Replicas)
	}

	// An object with no replicas field decodes to nil, not to a synthesized value.
	var withNone ControlPlaneSpec
	if err := json.Unmarshal([]byte(`{}`), &withNone); err != nil {
		t.Fatalf("unmarshal empty: %v", err)
	}
	if withNone.Replicas != nil {
		t.Fatalf("omitted replicas decoded to %v, want nil", *withNone.Replicas)
	}

	// An object with explicit two round-trips through marshal/unmarshal unchanged.
	orig := ControlPlaneSpec{Replicas: ptrInt32(2)}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back ControlPlaneSpec
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.Replicas == nil || *back.Replicas != 2 {
		t.Fatalf("replicas:2 did not round-trip, got %v", back.Replicas)
	}
}
