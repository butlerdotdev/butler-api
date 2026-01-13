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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IdentityProviderType defines the type of identity provider.
// +kubebuilder:validation:Enum=oidc
type IdentityProviderType string

const (
	// IdentityProviderTypeOIDC is an OpenID Connect provider.
	// Supports Google Workspace, Microsoft Entra ID, Okta, Auth0, Keycloak, etc.
	IdentityProviderTypeOIDC IdentityProviderType = "oidc"
)

// IdentityProviderPhase represents the current phase of an IdentityProvider.
// +kubebuilder:validation:Enum=Pending;Ready;Failed
type IdentityProviderPhase string

const (
	// IdentityProviderPhasePending indicates the provider is being validated.
	IdentityProviderPhasePending IdentityProviderPhase = "Pending"

	// IdentityProviderPhaseReady indicates the provider is validated and ready.
	IdentityProviderPhaseReady IdentityProviderPhase = "Ready"

	// IdentityProviderPhaseFailed indicates validation failed.
	IdentityProviderPhaseFailed IdentityProviderPhase = "Failed"
)

// IdentityProviderSpec defines the desired state of IdentityProvider.
type IdentityProviderSpec struct {
	// Type specifies the identity provider type.
	// Currently only "oidc" is supported.
	// +kubebuilder:validation:Required
	Type IdentityProviderType `json:"type"`

	// DisplayName is a human-readable name for this provider.
	// Shown in the login UI when multiple providers are configured.
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// OIDC contains OpenID Connect configuration.
	// Required when type is "oidc".
	// +optional
	OIDC *OIDCConfig `json:"oidc,omitempty"`
}

// OIDCConfig contains OpenID Connect provider configuration.
// Butler uses OIDC Discovery to automatically configure endpoints.
type OIDCConfig struct {
	// IssuerURL is the OIDC provider's issuer URL.
	// Butler appends /.well-known/openid-configuration for discovery.
	// Examples:
	//   - Google: https://accounts.google.com
	//   - Microsoft: https://login.microsoftonline.com/{tenant}/v2.0
	//   - Okta: https://{domain}.okta.com
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https://`
	IssuerURL string `json:"issuerURL"`

	// ClientID is the OAuth2 client ID from the identity provider.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	ClientID string `json:"clientID"`

	// ClientSecretRef references a Secret containing the OAuth2 client secret.
	// The Secret must contain a key named "client-secret".
	// +kubebuilder:validation:Required
	ClientSecretRef SecretReference `json:"clientSecretRef"`

	// RedirectURL is the OAuth2 callback URL.
	// Must match the redirect URI configured in the identity provider.
	// Example: https://butler.example.com/api/auth/callback
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https?://`
	RedirectURL string `json:"redirectURL"`

	// Scopes are the OAuth2 scopes to request.
	// Defaults to ["openid", "email", "profile"] if not specified.
	// Some providers require additional scopes for group information:
	//   - Microsoft: add "groups" or use Graph API
	//   - Okta: add "groups"
	//   - Google: groups require separate Cloud Identity API call
	// +optional
	Scopes []string `json:"scopes,omitempty"`

	// GroupsClaim is the JWT claim containing group memberships.
	// Defaults to "groups". Set to empty string to disable group extraction.
	// Note: Google Workspace doesn't include groups in the ID token by default.
	// +kubebuilder:default="groups"
	// +optional
	GroupsClaim string `json:"groupsClaim,omitempty"`

	// EmailClaim is the JWT claim containing the user's email.
	// Defaults to "email".
	// +kubebuilder:default="email"
	// +optional
	EmailClaim string `json:"emailClaim,omitempty"`

	// HostedDomain restricts authentication to a specific domain.
	// Only supported by Google Workspace. Users outside this domain
	// will see an error during Google authentication.
	// Example: "butlerlabs.dev"
	// +optional
	HostedDomain string `json:"hostedDomain,omitempty"`

	// InsecureSkipVerify disables TLS certificate verification.
	// WARNING: Only use for development with self-signed certificates.
	// +kubebuilder:default=false
	// +optional
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}

// IdentityProviderStatus defines the observed state of IdentityProvider.
type IdentityProviderStatus struct {
	// Conditions represent the latest available observations.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current phase of the provider.
	// +optional
	Phase IdentityProviderPhase `json:"phase,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastValidatedTime is when the provider was last successfully validated.
	// +optional
	LastValidatedTime *metav1.Time `json:"lastValidatedTime,omitempty"`

	// DiscoveredEndpoints contains endpoints discovered via OIDC Discovery.
	// +optional
	DiscoveredEndpoints *OIDCDiscoveredEndpoints `json:"discoveredEndpoints,omitempty"`

	// Message provides additional status information.
	// +optional
	Message string `json:"message,omitempty"`
}

// OIDCDiscoveredEndpoints contains endpoints from OIDC Discovery.
type OIDCDiscoveredEndpoints struct {
	// AuthorizationEndpoint is the OAuth2 authorization URL.
	// +optional
	AuthorizationEndpoint string `json:"authorizationEndpoint,omitempty"`

	// TokenEndpoint is the OAuth2 token URL.
	// +optional
	TokenEndpoint string `json:"tokenEndpoint,omitempty"`

	// UserInfoEndpoint is the OIDC userinfo URL.
	// +optional
	UserInfoEndpoint string `json:"userInfoEndpoint,omitempty"`

	// JWKSURI is the JSON Web Key Set URL for token validation.
	// +optional
	JWKSURI string `json:"jwksURI,omitempty"`
}

// IdentityProvider condition types.
const (
	// IdentityProviderConditionDiscovered indicates OIDC discovery succeeded.
	IdentityProviderConditionDiscovered = "Discovered"

	// IdentityProviderConditionSecretValid indicates the client secret is valid.
	IdentityProviderConditionSecretValid = "SecretValid"

	// IdentityProviderConditionReady indicates the provider is ready for use.
	IdentityProviderConditionReady = "Ready"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=idp
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type",description="Provider type"
// +kubebuilder:printcolumn:name="Issuer",type="string",JSONPath=".spec.oidc.issuerURL",description="OIDC issuer URL"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// IdentityProvider defines an external identity provider for authentication.
// Butler uses IdentityProviders to authenticate users via OIDC (OpenID Connect).
// Users are matched to Teams based on their email or group memberships.
//
// Example for Google Workspace:
//
//	apiVersion: butler.butlerlabs.dev/v1alpha1
//	kind: IdentityProvider
//	metadata:
//	  name: google-workspace
//	spec:
//	  type: oidc
//	  displayName: "Google Workspace"
//	  oidc:
//	    issuerURL: "https://accounts.google.com"
//	    clientID: "xxx.apps.googleusercontent.com"
//	    clientSecretRef:
//	      name: google-oidc-secret
//	      namespace: butler-system
//	      key: client-secret
//	    redirectURL: "https://butler.example.com/api/auth/callback"
//	    hostedDomain: "example.com"
type IdentityProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IdentityProviderSpec   `json:"spec,omitempty"`
	Status IdentityProviderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IdentityProviderList contains a list of IdentityProvider.
type IdentityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IdentityProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IdentityProvider{}, &IdentityProviderList{})
}

// Helper methods

// GetScopes returns the configured scopes or defaults.
func (idp *IdentityProvider) GetScopes() []string {
	if idp.Spec.OIDC != nil && len(idp.Spec.OIDC.Scopes) > 0 {
		return idp.Spec.OIDC.Scopes
	}
	return []string{"openid", "email", "profile"}
}

// GetGroupsClaim returns the configured groups claim or default.
func (idp *IdentityProvider) GetGroupsClaim() string {
	if idp.Spec.OIDC != nil && idp.Spec.OIDC.GroupsClaim != "" {
		return idp.Spec.OIDC.GroupsClaim
	}
	return "groups"
}

// GetEmailClaim returns the configured email claim or default.
func (idp *IdentityProvider) GetEmailClaim() string {
	if idp.Spec.OIDC != nil && idp.Spec.OIDC.EmailClaim != "" {
		return idp.Spec.OIDC.EmailClaim
	}
	return "email"
}

// GetDisplayName returns the display name or a default based on issuer.
func (idp *IdentityProvider) GetDisplayName() string {
	if idp.Spec.DisplayName != "" {
		return idp.Spec.DisplayName
	}
	if idp.Spec.OIDC != nil {
		switch {
		case contains(idp.Spec.OIDC.IssuerURL, "accounts.google.com"):
			return "Google"
		case contains(idp.Spec.OIDC.IssuerURL, "login.microsoftonline.com"):
			return "Microsoft"
		case contains(idp.Spec.OIDC.IssuerURL, "okta.com"):
			return "Okta"
		}
	}
	return idp.Name
}

// IsReady returns true if the provider is in Ready phase.
func (idp *IdentityProvider) IsReady() bool {
	return idp.Status.Phase == IdentityProviderPhaseReady
}

// contains checks if s contains substr (simple helper to avoid importing strings).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
