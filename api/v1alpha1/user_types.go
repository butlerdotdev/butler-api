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

// UserAuthType defines how a user authenticates.
// +kubebuilder:validation:Enum=sso;internal
type UserAuthType string

const (
	// UserAuthTypeSSO indicates the user authenticates via SSO/OIDC.
	UserAuthTypeSSO UserAuthType = "sso"

	// UserAuthTypeInternal indicates the user authenticates with email/password.
	UserAuthTypeInternal UserAuthType = "internal"
)

// UserSpec defines the desired state of a Butler user.
// Note: Passwords are NEVER stored in spec - users set their own via invite flow.
type UserSpec struct {
	// Email is the user's email address, used for Team membership matching.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=email
	Email string `json:"email"`

	// DisplayName is the user's display name shown in the UI.
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Disabled prevents the user from logging in.
	// +optional
	// +kubebuilder:default=false
	Disabled bool `json:"disabled,omitempty"`

	// Avatar is an optional URL to the user's avatar image.
	// +optional
	Avatar string `json:"avatar,omitempty"`

	// AuthType indicates how this user authenticates.
	// SSO users are created automatically on first login.
	// Internal users are created by admins and use email/password.
	// +kubebuilder:default="internal"
	// +optional
	AuthType UserAuthType `json:"authType,omitempty"`

	// SSOProvider is the name of the SSO provider (e.g., "Google", "Okta").
	// Only set for SSO users.
	// +optional
	SSOProvider string `json:"ssoProvider,omitempty"`

	// SSOSubject is the unique subject identifier from the SSO provider.
	// This is the "sub" claim from the OIDC token.
	// Only set for SSO users.
	// +optional
	SSOSubject string `json:"ssoSubject,omitempty"`
}

// UserStatus defines the observed state of User.
type UserStatus struct {
	// Phase represents the current state of the user.
	// +kubebuilder:validation:Enum=Pending;Active;Disabled;Locked
	Phase UserPhase `json:"phase,omitempty"`

	// PasswordSecretRef references the Secret containing the bcrypt password hash.
	// This is automatically created when the user sets their password.
	// Only used for internal users.
	// +optional
	PasswordSecretRef *SecretReference `json:"passwordSecretRef,omitempty"`

	// InviteTokenHash is the SHA256 hash of the invite token.
	// The raw token is only shown once when the user is created.
	// Only used for internal users.
	// +optional
	InviteTokenHash string `json:"inviteTokenHash,omitempty"`

	// InviteExpiresAt is when the invite token expires.
	// Only used for internal users.
	// +optional
	InviteExpiresAt *metav1.Time `json:"inviteExpiresAt,omitempty"`

	// InviteSentAt is when the invite was generated.
	// Only used for internal users.
	// +optional
	InviteSentAt *metav1.Time `json:"inviteSentAt,omitempty"`

	// PasswordChangedAt is when the password was last set/changed.
	// Only used for internal users.
	// +optional
	PasswordChangedAt *metav1.Time `json:"passwordChangedAt,omitempty"`

	// LastLoginTime is when the user last successfully logged in.
	// +optional
	LastLoginTime *metav1.Time `json:"lastLoginTime,omitempty"`

	// LoginCount is the total number of successful logins.
	// +optional
	LoginCount int64 `json:"loginCount,omitempty"`

	// FailedLoginAttempts is the number of consecutive failed login attempts.
	// Resets to 0 on successful login.
	// Only used for internal users.
	// +optional
	FailedLoginAttempts int32 `json:"failedLoginAttempts,omitempty"`

	// LockedUntil is set when the account is temporarily locked due to failed attempts.
	// Only used for internal users.
	// +optional
	LockedUntil *metav1.Time `json:"lockedUntil,omitempty"`

	// Teams lists the teams this user belongs to (resolved from Team CRDs).
	// This is informational and updated periodically.
	// +optional
	Teams []UserTeamMembership `json:"teams,omitempty"`

	// Conditions represent the latest available observations.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// UserPhase represents the current phase of a user.
// +kubebuilder:validation:Enum=Pending;Active;Disabled;Locked
type UserPhase string

const (
	// UserPhasePending indicates the user has been invited but hasn't set password.
	UserPhasePending UserPhase = "Pending"

	// UserPhaseActive indicates the user can log in.
	UserPhaseActive UserPhase = "Active"

	// UserPhaseDisabled indicates the user has been disabled by an admin.
	UserPhaseDisabled UserPhase = "Disabled"

	// UserPhaseLocked indicates the user is temporarily locked due to failed attempts.
	UserPhaseLocked UserPhase = "Locked"
)

// UserTeamMembership represents a user's membership in a team.
type UserTeamMembership struct {
	// Name is the team name.
	Name string `json:"name"`

	// Role is the user's role in the team.
	Role string `json:"role"`
}

// User condition types.
const (
	// UserConditionReady indicates the user account is ready for login.
	UserConditionReady = "Ready"

	// UserConditionInvitePending indicates waiting for user to accept invite.
	UserConditionInvitePending = "InvitePending"

	// UserConditionInviteExpired indicates the invite has expired.
	UserConditionInviteExpired = "InviteExpired"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=usr
// +kubebuilder:printcolumn:name="Email",type=string,JSONPath=`.spec.email`
// +kubebuilder:printcolumn:name="Display Name",type=string,JSONPath=`.spec.displayName`
// +kubebuilder:printcolumn:name="Auth",type=string,JSONPath=`.spec.authType`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Last Login",type=date,JSONPath=`.status.lastLoginTime`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// User represents a Butler user account.
// Users can authenticate via SSO (OIDC) or with email/password (internal).
//
// SSO User Flow:
// 1. User clicks "Sign in with Google/Okta/etc"
// 2. Butler creates User CRD automatically on first login
// 3. User is matched to Teams by email address
//
// Internal User Flow:
// 1. Admin creates User with email (no password)
// 2. Butler generates invite token, returns URL to admin
// 3. Admin shares invite URL with user (via Slack, email, etc.)
// 4. User clicks link, sets their own password
// 5. Password is hashed (bcrypt) and stored in a Secret
// 6. User status changes from Pending to Active
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserSpec   `json:"spec,omitempty"`
	Status UserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UserList contains a list of Users.
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}

// Helper methods

// IsSSO returns true if this is an SSO user.
func (u *User) IsSSO() bool {
	return u.Spec.AuthType == UserAuthTypeSSO
}

// IsInternal returns true if this is an internal user.
func (u *User) IsInternal() bool {
	return u.Spec.AuthType == UserAuthTypeInternal || u.Spec.AuthType == ""
}

// IsActive returns true if the user can log in.
func (u *User) IsActive() bool {
	return u.Status.Phase == UserPhaseActive && !u.Spec.Disabled
}

// IsDisabled returns true if the user is disabled.
func (u *User) IsDisabled() bool {
	return u.Spec.Disabled || u.Status.Phase == UserPhaseDisabled
}
