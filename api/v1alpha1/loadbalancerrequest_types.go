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

// LoadBalancerPhase represents the lifecycle phase of a LoadBalancerRequest.
// +kubebuilder:validation:Enum=Pending;Creating;Ready;Failed;Deleting
type LoadBalancerPhase string

const (
	// LoadBalancerPhasePending indicates the request has been received but not yet processed.
	LoadBalancerPhasePending LoadBalancerPhase = "Pending"

	// LoadBalancerPhaseCreating indicates the load balancer is being provisioned.
	LoadBalancerPhaseCreating LoadBalancerPhase = "Creating"

	// LoadBalancerPhaseReady indicates the load balancer is provisioned and has an endpoint.
	LoadBalancerPhaseReady LoadBalancerPhase = "Ready"

	// LoadBalancerPhaseFailed indicates the load balancer provisioning failed.
	LoadBalancerPhaseFailed LoadBalancerPhase = "Failed"

	// LoadBalancerPhaseDeleting indicates the load balancer is being torn down.
	LoadBalancerPhaseDeleting LoadBalancerPhase = "Deleting"
)

// LoadBalancerRequestSpec defines the desired state of a cloud load balancer
// for a management cluster control plane endpoint.
type LoadBalancerRequestSpec struct {
	// ClusterName is the name of the cluster this load balancer serves.
	// Used as a prefix for cloud resource naming.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	ClusterName string `json:"clusterName"`

	// ProviderConfigRef references the ProviderConfig with cloud credentials.
	// +kubebuilder:validation:Required
	ProviderConfigRef ProviderReference `json:"providerConfigRef"`

	// Port is the target port on backend instances.
	// +optional
	// +kubebuilder:default=6443
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port,omitempty"`

	// HealthCheckPort is the port used for backend health checks.
	// Defaults to the same value as Port.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	HealthCheckPort int32 `json:"healthCheckPort,omitempty"`

	// Targets are the backend instances to register with the load balancer.
	// Updated incrementally as machines come online.
	// +optional
	Targets []LoadBalancerTarget `json:"targets,omitempty"`
}

// LoadBalancerTarget identifies a backend instance for load balancer registration.
type LoadBalancerTarget struct {
	// IP is the target IP address.
	// +optional
	IP string `json:"ip,omitempty"`

	// InstanceID is the cloud-provider instance identifier.
	// Used for instance-based target groups (e.g., AWS NLB).
	// +optional
	InstanceID string `json:"instanceID,omitempty"`

	// InstanceName is the cloud-provider instance name.
	// Used for name-based registration (e.g., GCP target pools).
	// +optional
	InstanceName string `json:"instanceName,omitempty"`
}

// LoadBalancerRequestStatus defines the observed state of LoadBalancerRequest.
type LoadBalancerRequestStatus struct {
	// Phase represents the current lifecycle phase.
	// +optional
	Phase LoadBalancerPhase `json:"phase,omitempty"`

	// Endpoint is the load balancer's IP address or DNS name.
	// Populated when Phase reaches Ready.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// ResourceID is the cloud-provider resource identifier for the load balancer.
	// Used during cleanup to locate the resource for deletion.
	// +optional
	ResourceID string `json:"resourceID,omitempty"`

	// FailureReason provides a machine-readable failure reason.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage provides a human-readable failure message.
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// RegisteredTargets is the number of targets currently registered with the load balancer.
	// +optional
	RegisteredTargets int32 `json:"registeredTargets,omitempty"`

	// Conditions represent the latest available observations of the request's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastUpdated is the timestamp of the last status update.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`

	// ObservedGeneration is the generation most recently observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=lbr
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.clusterName",description="Target cluster"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current phase"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.endpoint",description="LB endpoint"
// +kubebuilder:printcolumn:name="Targets",type="integer",JSONPath=".status.registeredTargets",description="Registered backends"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// LoadBalancerRequest is the Schema for the loadbalancerrequests API.
// It represents a request to create a cloud load balancer for a management
// cluster's control plane endpoint. Provider controllers watch this resource
// and provision the appropriate cloud-native load balancer (GCP forwarding rule,
// AWS NLB, or Azure Standard LB).
type LoadBalancerRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadBalancerRequestSpec   `json:"spec,omitempty"`
	Status LoadBalancerRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoadBalancerRequestList contains a list of LoadBalancerRequest.
type LoadBalancerRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadBalancerRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoadBalancerRequest{}, &LoadBalancerRequestList{})
}

// IsReady returns true if the load balancer is provisioned and has an endpoint.
func (lbr *LoadBalancerRequest) IsReady() bool {
	return lbr.Status.Phase == LoadBalancerPhaseReady && lbr.Status.Endpoint != ""
}

// IsFailed returns true if the load balancer provisioning failed.
func (lbr *LoadBalancerRequest) IsFailed() bool {
	return lbr.Status.Phase == LoadBalancerPhaseFailed
}

// IsTerminating returns true if the load balancer is being deleted.
func (lbr *LoadBalancerRequest) IsTerminating() bool {
	return lbr.Status.Phase == LoadBalancerPhaseDeleting
}

// SetPhase updates the phase and last updated timestamp.
func (lbr *LoadBalancerRequest) SetPhase(phase LoadBalancerPhase) {
	lbr.Status.Phase = phase
	now := metav1.Now()
	lbr.Status.LastUpdated = &now
}

// SetFailure sets the failure reason and message.
func (lbr *LoadBalancerRequest) SetFailure(reason, message string) {
	lbr.Status.FailureReason = reason
	lbr.Status.FailureMessage = message
	lbr.SetPhase(LoadBalancerPhaseFailed)
}

// GetHealthCheckPort returns the health check port, defaulting to the target port.
func (lbr *LoadBalancerRequest) GetHealthCheckPort() int32 {
	if lbr.Spec.HealthCheckPort > 0 {
		return lbr.Spec.HealthCheckPort
	}
	if lbr.Spec.Port > 0 {
		return lbr.Spec.Port
	}
	return 6443
}
