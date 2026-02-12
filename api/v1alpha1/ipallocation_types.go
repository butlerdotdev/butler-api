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

// IPAllocationType defines the purpose of an IP allocation.
// +kubebuilder:validation:Enum=nodes;loadbalancer
type IPAllocationType string

const (
	// IPAllocationTypeNodes is for worker node IPs.
	IPAllocationTypeNodes IPAllocationType = "nodes"

	// IPAllocationTypeLoadBalancer is for load balancer IPs.
	IPAllocationTypeLoadBalancer IPAllocationType = "loadbalancer"
)

// IPAllocationPhase represents the current phase of an IPAllocation.
// +kubebuilder:validation:Enum=Pending;Allocated;Released;Failed
type IPAllocationPhase string

const (
	// IPAllocationPhasePending indicates the allocation is waiting to be fulfilled.
	IPAllocationPhasePending IPAllocationPhase = "Pending"

	// IPAllocationPhaseAllocated indicates IPs have been assigned.
	IPAllocationPhaseAllocated IPAllocationPhase = "Allocated"

	// IPAllocationPhaseReleased indicates IPs have been released.
	IPAllocationPhaseReleased IPAllocationPhase = "Released"

	// IPAllocationPhaseFailed indicates the allocation failed.
	IPAllocationPhaseFailed IPAllocationPhase = "Failed"
)

// IPAllocationSpec defines the desired state of IPAllocation.
type IPAllocationSpec struct {
	// PoolRef references the NetworkPool to allocate from.
	// +kubebuilder:validation:Required
	PoolRef LocalObjectReference `json:"poolRef"`

	// TenantClusterRef references the TenantCluster this allocation is for.
	// +kubebuilder:validation:Required
	TenantClusterRef NamespacedObjectReference `json:"tenantClusterRef"`

	// Type specifies the purpose of the allocation.
	// +kubebuilder:validation:Required
	Type IPAllocationType `json:"type"`

	// Count is the number of IPs to allocate.
	// If not specified, defaults from the NetworkPool are used.
	// Ignored when PinnedRange is set.
	// +optional
	// +kubebuilder:validation:Minimum=1
	Count *int32 `json:"count,omitempty"`

	// PinnedRange requests a specific IP range instead of automatic allocation.
	// Used for migrating existing clusters to IPAM or reserving well-known addresses.
	// The allocator validates the range is within the pool and not already allocated.
	// +optional
	PinnedRange *PinnedIPRange `json:"pinnedRange,omitempty"`
}

// PinnedIPRange specifies an exact IP range to allocate.
type PinnedIPRange struct {
	// StartAddress is the first IP of the pinned range.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(\d{1,3}\.){3}\d{1,3}$`
	StartAddress string `json:"startAddress"`

	// EndAddress is the last IP of the pinned range.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(\d{1,3}\.){3}\d{1,3}$`
	EndAddress string `json:"endAddress"`
}

// IPAllocationStatus defines the observed state of IPAllocation.
type IPAllocationStatus struct {
	// Phase represents the current phase of the allocation.
	// +optional
	Phase IPAllocationPhase `json:"phase,omitempty"`

	// Conditions represent the latest available observations.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// CIDR is the allocated range in CIDR notation if power-of-2 aligned.
	// +optional
	CIDR string `json:"cidr,omitempty"`

	// StartAddress is the first IP in the allocated range.
	// +optional
	StartAddress string `json:"startAddress,omitempty"`

	// EndAddress is the last IP in the allocated range.
	// +optional
	EndAddress string `json:"endAddress,omitempty"`

	// Addresses lists all individual IPs in the allocated range.
	// +optional
	Addresses []string `json:"addresses,omitempty"`

	// AllocatedCount is the number of IPs allocated.
	// +optional
	AllocatedCount int32 `json:"allocatedCount,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// AllocatedAt is the timestamp when IPs were assigned.
	// +optional
	AllocatedAt *metav1.Time `json:"allocatedAt,omitempty"`

	// AllocatedBy identifies the controller that fulfilled the allocation.
	// +optional
	AllocatedBy string `json:"allocatedBy,omitempty"`

	// ReleasedAt is the timestamp when IPs were released.
	// +optional
	ReleasedAt *metav1.Time `json:"releasedAt,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ipa
// +kubebuilder:printcolumn:name="Pool",type="string",JSONPath=".spec.poolRef.name",description="Network pool"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.tenantClusterRef.name",description="Tenant cluster"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type",description="Allocation type"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Allocation phase"
// +kubebuilder:printcolumn:name="Start",type="string",JSONPath=".status.startAddress",description="Start IP"
// +kubebuilder:printcolumn:name="End",type="string",JSONPath=".status.endAddress",description="End IP"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// IPAllocation represents an individual IP allocation from a NetworkPool.
type IPAllocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPAllocationSpec   `json:"spec,omitempty"`
	Status IPAllocationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPAllocationList contains a list of IPAllocation.
type IPAllocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPAllocation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPAllocation{}, &IPAllocationList{})
}
