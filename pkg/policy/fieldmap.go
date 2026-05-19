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
	"fmt"

	butlerv1alpha1 "github.com/butlerdotdev/butler-api/api/v1alpha1"
)

// ReadField extracts the option value from a TenantCluster spec for a
// given (provider type, option type) pair. The returned path string is
// the dotted field path used in admission rejection messages, e.g.
// `spec.infrastructureOverride.nutanix.imageUUID`.
//
// Supported is false when the provider does not expose the option type
// (for example a policy targeting `cluster` against a Harvester
// provider). Callers should treat unsupported pairs as no-ops; the
// policy webhook accepts the policy because the same policy may target
// multiple providers where the option type is meaningful.
//
// ReadField is the only per-provider code path the resolver touches.
// New providers add a switch case here; the resolver does not branch
// on provider type anywhere else.
func ReadField(tc *butlerv1alpha1.TenantCluster, providerType butlerv1alpha1.ProviderType, ot butlerv1alpha1.OptionType) (value string, path string, supported bool) {
	if tc == nil || tc.Spec.InfrastructureOverride == nil {
		return "", "", false
	}
	override := tc.Spec.InfrastructureOverride

	switch providerType {
	case butlerv1alpha1.ProviderTypeNutanix:
		if override.Nutanix == nil {
			return "", "", false
		}
		switch ot {
		case butlerv1alpha1.OptionTypeImage:
			return override.Nutanix.ImageUUID, "spec.infrastructureOverride.nutanix.imageUUID", true
		case butlerv1alpha1.OptionTypeNetwork:
			return override.Nutanix.SubnetUUID, "spec.infrastructureOverride.nutanix.subnetUUID", true
		case butlerv1alpha1.OptionTypeCluster:
			return override.Nutanix.ClusterUUID, "spec.infrastructureOverride.nutanix.clusterUUID", true
		case butlerv1alpha1.OptionTypeStorageContainer:
			return override.Nutanix.StorageContainerUUID, "spec.infrastructureOverride.nutanix.storageContainerUUID", true
		}
	case butlerv1alpha1.ProviderTypeHarvester:
		if override.Harvester == nil {
			return "", "", false
		}
		switch ot {
		case butlerv1alpha1.OptionTypeImage:
			return override.Harvester.ImageName, "spec.infrastructureOverride.harvester.imageName", true
		case butlerv1alpha1.OptionTypeNetwork:
			return override.Harvester.NetworkName, "spec.infrastructureOverride.harvester.networkName", true
		case butlerv1alpha1.OptionTypeCluster:
			return "", "", false
		case butlerv1alpha1.OptionTypeStorageContainer:
			return "", "", false
		}
	}
	return "", "", false
}

// ValidateAgainstRule checks a single TenantCluster spec field against a
// resolved policy rule. Returns (ok, message). When ok is false, message
// is the human-readable rejection text suitable for admission errors.
// For modes that do not enforce (default, recommended), ValidateAgainstRule
// always returns ok=true.
func ValidateAgainstRule(tc *butlerv1alpha1.TenantCluster, providerType butlerv1alpha1.ProviderType, ot butlerv1alpha1.OptionType, rule butlerv1alpha1.OptionRule, policyName string) (ok bool, fieldPath string, message string) {
	switch rule.Mode {
	case butlerv1alpha1.OptionModeDefault, butlerv1alpha1.OptionModeRecommended:
		return true, "", ""
	case butlerv1alpha1.OptionModePin, butlerv1alpha1.OptionModeAllowList:
		value, path, supported := ReadField(tc, providerType, ot)
		if !supported {
			return true, "", ""
		}
		if value == "" {
			return true, "", ""
		}
		for _, allowed := range rule.Values {
			if allowed == value {
				return true, "", ""
			}
		}
		verb := "pins"
		if rule.Mode == butlerv1alpha1.OptionModeAllowList {
			verb = "allow-lists"
		}
		return false, path, fmt.Sprintf("ClusterCreationPolicy %q %s %s to %v; got %q", policyName, verb, ot, rule.Values, value)
	default:
		return true, "", ""
	}
}
