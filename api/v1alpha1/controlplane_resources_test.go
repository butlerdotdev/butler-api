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

	"k8s.io/apimachinery/pkg/api/resource"
)

// TestControlPlaneResourcesKonnectivityOmitEmpty verifies the new konnectivity
// field is omitted when unset (backward compatible with existing objects) and
// round-trips when set.
func TestControlPlaneResourcesKonnectivityOmitEmpty(t *testing.T) {
	// Omitted: existing objects with only apiServer must not gain a konnectivity key.
	b, err := json.Marshal(ControlPlaneResourcesSpec{APIServer: &ComponentResources{}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "konnectivity") {
		t.Fatalf("konnectivity must be omitted when unset, got %s", b)
	}

	// Explicit: a konnectivity request round-trips.
	cpu := resource.MustParse("10m")
	mem := resource.MustParse("32Mi")
	in := ControlPlaneResourcesSpec{
		Konnectivity: &ComponentResources{
			Requests: &ResourceQuantities{CPU: &cpu, Memory: &mem},
		},
	}
	bb, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out ControlPlaneResourcesSpec
	if err := json.Unmarshal(bb, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Konnectivity == nil || out.Konnectivity.Requests == nil ||
		out.Konnectivity.Requests.Memory == nil || out.Konnectivity.Requests.Memory.String() != "32Mi" {
		t.Fatalf("konnectivity did not round-trip, got %+v", out.Konnectivity)
	}
}
