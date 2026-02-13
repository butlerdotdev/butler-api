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

// ReservedRange defines a range of IPs excluded from allocation.
type ReservedRange struct {
	// CIDR is the reserved range in CIDR notation.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(\d{1,3}\.){3}\d{1,3}/\d{1,2}$`
	CIDR string `json:"cidr"`

	// Description explains why this range is reserved.
	// +optional
	Description string `json:"description,omitempty"`
}

// TenantAllocationDefaults defines default allocation sizes per tenant.
type TenantAllocationDefaults struct {
	// NodesPerTenant is the default number of node IPs per tenant.
	// +kubebuilder:default=5
	// +kubebuilder:validation:Minimum=1
	// +optional
	NodesPerTenant int32 `json:"nodesPerTenant,omitempty"`

	// LBPoolPerTenant is the default number of load balancer IPs per tenant.
	// +kubebuilder:default=8
	// +kubebuilder:validation:Minimum=1
	// +optional
	LBPoolPerTenant int32 `json:"lbPoolPerTenant,omitempty"`
}

// TenantAllocationConfig defines the allocatable sub-range and defaults.
type TenantAllocationConfig struct {
	// Start is the first allocatable IP address.
	// +kubebuilder:validation:Required
	Start string `json:"start"`

	// End is the last allocatable IP address.
	// +kubebuilder:validation:Required
	End string `json:"end"`

	// Defaults defines default allocation sizes per tenant.
	// +optional
	Defaults TenantAllocationDefaults `json:"defaults,omitempty"`
}

// NetworkPoolSpec defines the desired state of NetworkPool.
type NetworkPoolSpec struct {
	// CIDR is the network range in CIDR notation.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(\d{1,3}\.){3}\d{1,3}/\d{1,2}$`
	CIDR string `json:"cidr"`

	// Reserved defines ranges excluded from allocation.
	// +optional
	Reserved []ReservedRange `json:"reserved,omitempty"`

	// TenantAllocation configures the allocatable sub-range and defaults.
	// If not specified, the entire CIDR (minus reserved ranges) is allocatable.
	// +optional
	TenantAllocation *TenantAllocationConfig `json:"tenantAllocation,omitempty"`
}

// NetworkPoolStatus defines the observed state of NetworkPool.
type NetworkPoolStatus struct {
	// Conditions represent the latest available observations.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// TotalIPs is the total number of usable IPs (excluding reserved).
	// +optional
	TotalIPs int32 `json:"totalIPs,omitempty"`

	// AllocatedIPs is the number of currently allocated IPs.
	// +optional
	AllocatedIPs int32 `json:"allocatedIPs,omitempty"`

	// AvailableIPs is the number of available IPs.
	// +optional
	AvailableIPs int32 `json:"availableIPs,omitempty"`

	// AllocationCount is the total number of IPAllocations from this pool.
	// +optional
	AllocationCount int32 `json:"allocationCount,omitempty"`

	// FragmentationPercent indicates how fragmented the free space is (0-100).
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	FragmentationPercent *int32 `json:"fragmentationPercent,omitempty"`

	// LargestFreeBlock is the size of the largest contiguous free block.
	// +optional
	LargestFreeBlock int32 `json:"largestFreeBlock,omitempty"`

	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=np
// +kubebuilder:printcolumn:name="CIDR",type="string",JSONPath=".spec.cidr",description="Network CIDR"
// +kubebuilder:printcolumn:name="Available",type="integer",JSONPath=".status.availableIPs",description="Available IPs"
// +kubebuilder:printcolumn:name="Allocated",type="integer",JSONPath=".status.allocatedIPs",description="Allocated IPs"
// +kubebuilder:printcolumn:name="Total",type="integer",JSONPath=".status.totalIPs",description="Total usable IPs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// NetworkPool defines a platform-level IP pool for on-prem IPAM.
type NetworkPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkPoolSpec   `json:"spec,omitempty"`
	Status NetworkPoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkPoolList contains a list of NetworkPool.
type NetworkPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkPool{}, &NetworkPoolList{})
}
