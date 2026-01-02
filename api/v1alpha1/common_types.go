/*
Copyright 2026 Butler Labs.
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

// ProviderReference references a ProviderConfig
type ProviderReference struct {
	// Name of the ProviderConfig
	Name string `json:"name"`

	// Namespace of the ProviderConfig
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SecretReference references a Secret
type SecretReference struct {
	// Name of the secret
	Name string `json:"name"`

	// Namespace of the secret
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Key in the secret
	// +optional
	Key string `json:"key,omitempty"`
}

// MachineRole defines the role of a machine
// +kubebuilder:validation:Enum=control-plane;worker
type MachineRole string

const (
	MachineRoleControlPlane MachineRole = "control-plane"
	MachineRoleWorker       MachineRole = "worker"
)

// MachinePhase defines the phase of a machine
type MachinePhase string

const (
	MachinePhasePending  MachinePhase = "Pending"
	MachinePhaseCreating MachinePhase = "Creating"
	MachinePhaseRunning  MachinePhase = "Running"
	MachinePhaseFailed   MachinePhase = "Failed"
	MachinePhaseDeleting MachinePhase = "Deleting"
)
