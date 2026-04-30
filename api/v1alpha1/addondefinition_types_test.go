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

func TestGetEffectiveTier(t *testing.T) {
	tests := []struct {
		name     string
		tier     AddonTier
		platform bool
		want     string
	}{
		{
			name: "explicit infrastructure tier",
			tier: AddonTierInfrastructure,
			want: "infrastructure",
		},
		{
			name: "explicit apps tier",
			tier: AddonTierApps,
			want: "apps",
		},
		{
			name:     "explicit tier overrides platform",
			tier:     AddonTierApps,
			platform: true,
			want:     "apps",
		},
		{
			name:     "platform fallback to infrastructure",
			platform: true,
			want:     "infrastructure",
		},
		{
			name: "non-platform fallback to apps",
			want: "apps",
		},
		{
			name:     "explicit infrastructure with platform true is redundant but valid",
			tier:     AddonTierInfrastructure,
			platform: true,
			want:     "infrastructure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad := &AddonDefinition{
				Spec: AddonDefinitionSpec{
					Tier:     tt.tier,
					Platform: tt.platform,
				},
			}
			got := ad.GetEffectiveTier()
			if got != tt.want {
				t.Errorf("GetEffectiveTier() = %q, want %q", got, tt.want)
			}
		})
	}
}
