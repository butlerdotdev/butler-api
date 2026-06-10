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

import "testing"

func TestClusterBootstrapProviderPredicates(t *testing.T) {
	tests := []struct {
		provider  string
		wantLocal bool
		wantCloud bool
	}{
		{provider: "harvester", wantLocal: false, wantCloud: false},
		{provider: "nutanix", wantLocal: false, wantCloud: false},
		{provider: "proxmox", wantLocal: false, wantCloud: false},
		{provider: "aws", wantLocal: false, wantCloud: true},
		{provider: "azure", wantLocal: false, wantCloud: true},
		{provider: "gcp", wantLocal: false, wantCloud: true},
		{provider: "local", wantLocal: true, wantCloud: false},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			cb := &ClusterBootstrap{Spec: ClusterBootstrapSpec{Provider: tt.provider}}

			if got := cb.IsLocal(); got != tt.wantLocal {
				t.Errorf("IsLocal() for %q = %v, want %v", tt.provider, got, tt.wantLocal)
			}
			if got := cb.IsCloudProvider(); got != tt.wantCloud {
				t.Errorf("IsCloudProvider() for %q = %v, want %v", tt.provider, got, tt.wantCloud)
			}
			if cb.IsLocal() && cb.IsCloudProvider() {
				t.Errorf("provider %q is both local and cloud; predicates must be mutually exclusive", tt.provider)
			}
		})
	}
}
