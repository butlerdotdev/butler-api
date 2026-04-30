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

func TestGetEffectivePlatformRole(t *testing.T) {
	tests := []struct {
		name            string
		platformRole    string
		isPlatformAdmin bool
		want            string
	}{
		{
			name:         "PlatformRole admin",
			platformRole: "admin",
			want:         "admin",
		},
		{
			name:            "IsPlatformAdmin fallback",
			isPlatformAdmin: true,
			want:            "admin",
		},
		{
			name:         "PlatformRole viewer",
			platformRole: "viewer",
			want:         "viewer",
		},
		{
			name:            "IsPlatformAdmin overrides viewer",
			platformRole:    "viewer",
			isPlatformAdmin: true,
			want:            "admin",
		},
		{
			name: "both unset",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Spec: UserSpec{
					PlatformRole:    tt.platformRole,
					IsPlatformAdmin: tt.isPlatformAdmin,
				},
			}
			got := u.GetEffectivePlatformRole()
			if got != tt.want {
				t.Errorf("GetEffectivePlatformRole() = %q, want %q", got, tt.want)
			}
		})
	}
}
