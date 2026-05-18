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

package policy

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	butlerv1alpha1 "github.com/butlerdotdev/butler-api/api/v1alpha1"
)

func makePolicy(name string, scope butlerv1alpha1.PolicyScope, providers []butlerv1alpha1.ProviderType, opts map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule) butlerv1alpha1.ClusterCreationPolicy {
	return butlerv1alpha1.ClusterCreationPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: butlerv1alpha1.ClusterCreationPolicySpec{
			Scope:           scope,
			TargetProviders: providers,
			Options:         opts,
		},
	}
}

func clusterWide() butlerv1alpha1.PolicyScope {
	return butlerv1alpha1.PolicyScope{ClusterWide: &butlerv1alpha1.ClusterWideScope{}}
}

func teamScope(team string) butlerv1alpha1.PolicyScope {
	return butlerv1alpha1.PolicyScope{Team: &butlerv1alpha1.TeamScope{TeamRef: butlerv1alpha1.LocalObjectReference{Name: team}}}
}

func teamEnvScope(team, env string) butlerv1alpha1.PolicyScope {
	return butlerv1alpha1.PolicyScope{TeamAndEnvironment: &butlerv1alpha1.TeamEnvironmentScope{
		TeamRef:         butlerv1alpha1.LocalObjectReference{Name: team},
		EnvironmentName: env,
	}}
}

func TestResolve_NoPoliciesReturnsEmpty(t *testing.T) {
	out := Resolve(ResolutionContext{TeamName: "acme", ProviderType: butlerv1alpha1.ProviderTypeNutanix}, nil)
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(out))
	}
}

func TestResolve_ClusterWideOnly(t *testing.T) {
	policies := []butlerv1alpha1.ClusterCreationPolicy{
		makePolicy("global-image-allow", clusterWide(), []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix}, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModeAllowList, Values: []string{"img-a", "img-b"}},
		}),
	}
	out := Resolve(ResolutionContext{TeamName: "acme", EnvironmentName: "prod", ProviderType: butlerv1alpha1.ProviderTypeNutanix}, policies)
	if got, ok := out[butlerv1alpha1.OptionTypeImage]; !ok || got.Mode != butlerv1alpha1.OptionModeAllowList {
		t.Fatalf("expected image allowList from cluster-wide, got %+v", out)
	}
}

// TestResolve_SpecificityWins reproduces the ADR-018 §6 worked example:
// a clusterWide image allowList plus a teamAndEnvironment network pin.
// The image resolution falls through to clusterWide; the network
// resolution stops at teamAndEnvironment.
func TestResolve_SpecificityWins(t *testing.T) {
	policies := []butlerv1alpha1.ClusterCreationPolicy{
		makePolicy("global-no-deprecated-images", clusterWide(), []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix}, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModeAllowList, Values: []string{"img-rocky-9", "img-talos-1.7", "img-talos-1.8"}},
		}),
		makePolicy("acme-prod-network-pin", teamEnvScope("acme", "prod"), []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix}, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeNetwork: {Mode: butlerv1alpha1.OptionModePin, Values: []string{"net-prod-vlan-200"}},
		}),
	}
	ctx := ResolutionContext{TeamName: "acme", EnvironmentName: "prod", ProviderType: butlerv1alpha1.ProviderTypeNutanix}
	out := Resolve(ctx, policies)

	if got, ok := out[butlerv1alpha1.OptionTypeImage]; !ok || got.Mode != butlerv1alpha1.OptionModeAllowList || len(got.Values) != 3 {
		t.Errorf("expected image allowList[3] from clusterWide, got %+v", got)
	}
	if got, ok := out[butlerv1alpha1.OptionTypeNetwork]; !ok || got.Mode != butlerv1alpha1.OptionModePin || len(got.Values) != 1 {
		t.Errorf("expected network pin[1] from teamAndEnv, got %+v", got)
	}
	if _, ok := out[butlerv1alpha1.OptionTypeCluster]; ok {
		t.Errorf("expected no cluster rule, got one")
	}
	if _, ok := out[butlerv1alpha1.OptionTypeStorageContainer]; ok {
		t.Errorf("expected no storageContainer rule, got one")
	}
}

// TestResolve_TeamCarveOutWinsOverClusterWide is the canonical
// specificity-wins illustration from ADR §Patterns: a cluster-wide
// deprecation list plus a team carve-out that allows the deprecated
// image. The team rule wins for the target team; everyone else falls
// through to the cluster-wide rule.
func TestResolve_TeamCarveOutWinsOverClusterWide(t *testing.T) {
	policies := []butlerv1alpha1.ClusterCreationPolicy{
		makePolicy("global-deny-rocky-8", clusterWide(), []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix}, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModeAllowList, Values: []string{"img-rocky-9", "img-talos-1.8"}},
		}),
		makePolicy("team-legacy-rocky-8", teamScope("team-legacy"), []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix}, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModeAllowList, Values: []string{"img-rocky-8", "img-rocky-9", "img-talos-1.8"}},
		}),
	}

	legacyOut := Resolve(ResolutionContext{TeamName: "team-legacy", ProviderType: butlerv1alpha1.ProviderTypeNutanix}, policies)
	if got := legacyOut[butlerv1alpha1.OptionTypeImage]; len(got.Values) != 3 || got.Values[0] != "img-rocky-8" {
		t.Errorf("team-legacy expected to see 3 images including rocky-8, got %+v", got.Values)
	}

	otherOut := Resolve(ResolutionContext{TeamName: "team-other", ProviderType: butlerv1alpha1.ProviderTypeNutanix}, policies)
	if got := otherOut[butlerv1alpha1.OptionTypeImage]; len(got.Values) != 2 {
		t.Errorf("team-other expected to see 2 images (cluster-wide fallback), got %+v", got.Values)
	}
}

func TestResolve_TargetProvidersFiltering(t *testing.T) {
	policies := []butlerv1alpha1.ClusterCreationPolicy{
		makePolicy("nutanix-only", clusterWide(), []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix}, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModePin, Values: []string{"img-only-nutanix"}},
		}),
	}
	nutanixOut := Resolve(ResolutionContext{TeamName: "acme", ProviderType: butlerv1alpha1.ProviderTypeNutanix}, policies)
	if _, ok := nutanixOut[butlerv1alpha1.OptionTypeImage]; !ok {
		t.Error("nutanix context expected to match nutanix-only policy")
	}
	harvOut := Resolve(ResolutionContext{TeamName: "acme", ProviderType: butlerv1alpha1.ProviderTypeHarvester}, policies)
	if _, ok := harvOut[butlerv1alpha1.OptionTypeImage]; ok {
		t.Error("harvester context expected no rule from nutanix-only policy")
	}
}

func TestResolve_EmptyTargetProvidersMatchesAll(t *testing.T) {
	policies := []butlerv1alpha1.ClusterCreationPolicy{
		makePolicy("any-provider", clusterWide(), nil, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModeRecommended, Values: []string{"img-recommended"}},
		}),
	}
	for _, pt := range []butlerv1alpha1.ProviderType{butlerv1alpha1.ProviderTypeNutanix, butlerv1alpha1.ProviderTypeHarvester} {
		out := Resolve(ResolutionContext{TeamName: "acme", ProviderType: pt}, policies)
		if _, ok := out[butlerv1alpha1.OptionTypeImage]; !ok {
			t.Errorf("provider %q expected to match empty-targets policy", pt)
		}
	}
}

func TestResolveWithSources_NamesContributingPolicy(t *testing.T) {
	policies := []butlerv1alpha1.ClusterCreationPolicy{
		makePolicy("cw-img", clusterWide(), nil, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModePin, Values: []string{"img-cw"}},
		}),
		makePolicy("team-img", teamScope("acme"), nil, map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{
			butlerv1alpha1.OptionTypeImage: {Mode: butlerv1alpha1.OptionModePin, Values: []string{"img-team"}},
		}),
	}
	rules, sources := ResolveWithSources(ResolutionContext{TeamName: "acme", ProviderType: butlerv1alpha1.ProviderTypeNutanix}, policies)
	if rules[butlerv1alpha1.OptionTypeImage].Values[0] != "img-team" {
		t.Errorf("expected team rule to win, got %+v", rules[butlerv1alpha1.OptionTypeImage])
	}
	if sources[butlerv1alpha1.OptionTypeImage] != "team-img" {
		t.Errorf("expected source name 'team-img', got %q", sources[butlerv1alpha1.OptionTypeImage])
	}
}

func TestValidateAgainstRule_PinReject(t *testing.T) {
	tc := &butlerv1alpha1.TenantCluster{
		Spec: butlerv1alpha1.TenantClusterSpec{
			InfrastructureOverride: &butlerv1alpha1.InfrastructureOverride{
				Nutanix: &butlerv1alpha1.NutanixOverride{ImageUUID: "img-evil"},
			},
		},
	}
	rule := butlerv1alpha1.OptionRule{Mode: butlerv1alpha1.OptionModePin, Values: []string{"img-good"}}
	ok, path, msg := ValidateAgainstRule(tc, butlerv1alpha1.ProviderTypeNutanix, butlerv1alpha1.OptionTypeImage, rule, "test-pin")
	if ok {
		t.Fatalf("expected reject, got ok=true")
	}
	if path != "spec.infrastructureOverride.nutanix.imageUUID" {
		t.Errorf("unexpected path %q", path)
	}
	if msg == "" {
		t.Error("expected non-empty rejection message")
	}
}

func TestValidateAgainstRule_PinAccept(t *testing.T) {
	tc := &butlerv1alpha1.TenantCluster{
		Spec: butlerv1alpha1.TenantClusterSpec{
			InfrastructureOverride: &butlerv1alpha1.InfrastructureOverride{
				Nutanix: &butlerv1alpha1.NutanixOverride{ImageUUID: "img-good"},
			},
		},
	}
	rule := butlerv1alpha1.OptionRule{Mode: butlerv1alpha1.OptionModePin, Values: []string{"img-good"}}
	ok, _, _ := ValidateAgainstRule(tc, butlerv1alpha1.ProviderTypeNutanix, butlerv1alpha1.OptionTypeImage, rule, "test-pin")
	if !ok {
		t.Error("expected accept")
	}
}

func TestValidateAgainstRule_DefaultDoesNotEnforce(t *testing.T) {
	tc := &butlerv1alpha1.TenantCluster{
		Spec: butlerv1alpha1.TenantClusterSpec{
			InfrastructureOverride: &butlerv1alpha1.InfrastructureOverride{
				Nutanix: &butlerv1alpha1.NutanixOverride{ImageUUID: "img-anything"},
			},
		},
	}
	rule := butlerv1alpha1.OptionRule{Mode: butlerv1alpha1.OptionModeDefault, Default: "img-default"}
	ok, _, _ := ValidateAgainstRule(tc, butlerv1alpha1.ProviderTypeNutanix, butlerv1alpha1.OptionTypeImage, rule, "test-default")
	if !ok {
		t.Error("default mode should not reject")
	}
}

func TestValidateAgainstRule_UnsupportedProviderOptionIsNoOp(t *testing.T) {
	tc := &butlerv1alpha1.TenantCluster{
		Spec: butlerv1alpha1.TenantClusterSpec{
			InfrastructureOverride: &butlerv1alpha1.InfrastructureOverride{
				Harvester: &butlerv1alpha1.HarvesterOverride{ImageName: "img"},
			},
		},
	}
	// Harvester has no storage container concept.
	rule := butlerv1alpha1.OptionRule{Mode: butlerv1alpha1.OptionModePin, Values: []string{"sc-only"}}
	ok, _, _ := ValidateAgainstRule(tc, butlerv1alpha1.ProviderTypeHarvester, butlerv1alpha1.OptionTypeStorageContainer, rule, "p")
	if !ok {
		t.Error("storageContainer rule on Harvester should be a no-op (unsupported)")
	}
}

func TestReadField_HarvesterImageMapsToImageName(t *testing.T) {
	tc := &butlerv1alpha1.TenantCluster{
		Spec: butlerv1alpha1.TenantClusterSpec{
			InfrastructureOverride: &butlerv1alpha1.InfrastructureOverride{
				Harvester: &butlerv1alpha1.HarvesterOverride{ImageName: "default/talos-1.8"},
			},
		},
	}
	v, path, supported := ReadField(tc, butlerv1alpha1.ProviderTypeHarvester, butlerv1alpha1.OptionTypeImage)
	if !supported {
		t.Fatal("expected Harvester+image to be supported")
	}
	if v != "default/talos-1.8" {
		t.Errorf("expected ImageName value, got %q", v)
	}
	if path != "spec.infrastructureOverride.harvester.imageName" {
		t.Errorf("unexpected path %q", path)
	}
}
