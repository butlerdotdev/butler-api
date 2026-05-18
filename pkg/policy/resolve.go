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

// Package policy implements the ADR-018 ClusterCreationPolicy
// resolution algorithm as a pure function set. Callers (butler-server,
// butler-controller) supply the candidate policies and the resolution
// context; the package returns the effective rule per option type.
//
// No Kubernetes client calls live in this package. Listing policies
// happens at callers. Per ADR-018 Decision section 6.
package policy

import (
	butlerv1alpha1 "github.com/butlerdotdev/butler-api/api/v1alpha1"
)

// ResolutionContext is the tuple under which a set of ClusterCreationPolicy
// resources is resolved. For list-time resolution (butler-server) the
// fields come from the operator's session. For admission-time resolution
// (butler-controller) they come from the TenantCluster being validated.
type ResolutionContext struct {
	// TeamName is the team owning the TenantCluster. Empty means the
	// request has no team context; only ClusterWide policies match.
	TeamName string

	// EnvironmentName is the team environment. Empty means no environment
	// is selected; only ClusterWide and Team scope policies match.
	EnvironmentName string

	// ProviderType is the underlying provider for the cluster. A policy
	// without TargetProviders matches every provider; a policy with a
	// non-empty TargetProviders matches only when the type appears in
	// the list.
	ProviderType butlerv1alpha1.ProviderType
}

// Resolve runs the ADR-018 section 6 algorithm and returns the effective
// rule per option type for the given context.
//
// Specificity wins: a teamAndEnvironment-scoped rule for an option type
// shadows any team or clusterWide rule for the same option type. Modes
// do not stack across tiers. Option types with no matching rule do not
// appear in the result.
//
// Resolve assumes that intra-tier conflicts (two policies in the same
// tier defining a rule for the same option type for the same context)
// have already been rejected at policy admission time. If the caller
// passes a conflicting set, Resolve picks the first rule it encounters
// per option type within that tier. Callers should treat such inputs
// as misconfigured admission.
func Resolve(ctx ResolutionContext, policies []butlerv1alpha1.ClusterCreationPolicy) map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule {
	tiers := binByTier(ctx, policies)
	out := map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{}

	for _, tier := range tiers {
		for _, pol := range tier {
			for optType, rule := range pol.Spec.Options {
				if _, already := out[optType]; already {
					continue
				}
				out[optType] = rule
			}
		}
	}
	return out
}

// ResolveWithSources is the same as Resolve but also returns the name of
// the policy that contributed each effective rule. Used at admission
// time so rejection messages can name the policy.
func ResolveWithSources(ctx ResolutionContext, policies []butlerv1alpha1.ClusterCreationPolicy) (map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule, map[butlerv1alpha1.OptionType]string) {
	tiers := binByTier(ctx, policies)
	rules := map[butlerv1alpha1.OptionType]butlerv1alpha1.OptionRule{}
	sources := map[butlerv1alpha1.OptionType]string{}

	for _, tier := range tiers {
		for _, pol := range tier {
			for optType, rule := range pol.Spec.Options {
				if _, already := rules[optType]; already {
					continue
				}
				rules[optType] = rule
				sources[optType] = pol.Name
			}
		}
	}
	return rules, sources
}

// binByTier returns the three tiers in most-specific-first order. Each
// tier holds the policies that match the resolution context.
func binByTier(ctx ResolutionContext, policies []butlerv1alpha1.ClusterCreationPolicy) [3][]butlerv1alpha1.ClusterCreationPolicy {
	var tiers [3][]butlerv1alpha1.ClusterCreationPolicy

	for _, pol := range policies {
		if !providerMatches(ctx.ProviderType, pol.Spec.TargetProviders) {
			continue
		}
		switch tierOf(ctx, pol.Spec.Scope) {
		case tierTeamAndEnv:
			tiers[0] = append(tiers[0], pol)
		case tierTeam:
			tiers[1] = append(tiers[1], pol)
		case tierClusterWide:
			tiers[2] = append(tiers[2], pol)
		}
	}
	return tiers
}

type tier int

const (
	tierNone tier = iota
	tierTeamAndEnv
	tierTeam
	tierClusterWide
)

// tierOf returns the specificity tier of a policy's scope under the
// given context, or tierNone if the scope does not match the context.
func tierOf(ctx ResolutionContext, scope butlerv1alpha1.PolicyScope) tier {
	switch {
	case scope.TeamAndEnvironment != nil:
		if scope.TeamAndEnvironment.TeamRef.Name == ctx.TeamName &&
			scope.TeamAndEnvironment.EnvironmentName == ctx.EnvironmentName &&
			ctx.TeamName != "" && ctx.EnvironmentName != "" {
			return tierTeamAndEnv
		}
		return tierNone
	case scope.Team != nil:
		if scope.Team.TeamRef.Name == ctx.TeamName && ctx.TeamName != "" {
			return tierTeam
		}
		return tierNone
	case scope.ClusterWide != nil:
		return tierClusterWide
	default:
		return tierNone
	}
}

// providerMatches returns true when the policy's TargetProviders is
// empty (matches all providers) or contains the given provider type.
func providerMatches(pt butlerv1alpha1.ProviderType, targets []butlerv1alpha1.ProviderType) bool {
	if len(targets) == 0 {
		return true
	}
	for _, t := range targets {
		if t == pt {
			return true
		}
	}
	return false
}
