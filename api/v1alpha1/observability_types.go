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

// ObservabilityConfig configures platform-level observability.
// This is stored in ButlerConfig and used as the default for all observability operations.
type ObservabilityConfig struct {
	// Pipeline configures the centralized observability pipeline cluster.
	// +optional
	Pipeline *ObservabilityPipelineConfig `json:"pipeline,omitempty"`

	// Collection configures default collection settings for tenant clusters.
	// +optional
	Collection *ObservabilityCollectionConfig `json:"collection,omitempty"`
}

// ObservabilityPipelineConfig configures the centralized pipeline.
type ObservabilityPipelineConfig struct {
	// ClusterRef references the TenantCluster serving as the observability pipeline.
	// +optional
	ClusterRef *NamespacedObjectReference `json:"clusterRef,omitempty"`

	// LogEndpoint is the Vector aggregator ingestion URL.
	// Example: "http://vector-aggregator.vector.svc:9000"
	// +optional
	LogEndpoint string `json:"logEndpoint,omitempty"`

	// MetricEndpoint is the optional remote-write endpoint for metrics.
	// +optional
	MetricEndpoint string `json:"metricEndpoint,omitempty"`

	// TraceEndpoint is the optional OTLP endpoint for traces.
	// Example: "tempo.tracing.svc:4317"
	// +optional
	TraceEndpoint string `json:"traceEndpoint,omitempty"`
}

// ObservabilityCollectionConfig configures default collection settings.
type ObservabilityCollectionConfig struct {
	// AutoEnroll controls whether new tenant clusters automatically get
	// observability agents installed. Stores intent only — not yet implemented
	// by a controller.
	// +optional
	AutoEnroll bool `json:"autoEnroll,omitempty"`

	// Logs configures default log collection settings.
	// +optional
	Logs *LogCollectionDefaults `json:"logs,omitempty"`

	// Metrics configures default metric collection settings.
	// +optional
	Metrics *MetricCollectionDefaults `json:"metrics,omitempty"`
}

// LogCollectionDefaults configures which log sources are collected by default.
type LogCollectionDefaults struct {
	// PodLogs enables collection of container stdout/stderr logs.
	// +optional
	PodLogs bool `json:"podLogs,omitempty"`

	// Journald enables collection of systemd journal logs.
	// +optional
	Journald bool `json:"journald,omitempty"`

	// KubernetesEvents enables collection of Kubernetes events.
	// +optional
	KubernetesEvents bool `json:"kubernetesEvents,omitempty"`
}

// MetricCollectionDefaults configures default metric collection settings.
type MetricCollectionDefaults struct {
	// Enabled controls whether metrics collection is enabled by default.
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// Retention is the Prometheus data retention period (e.g., "15d").
	// +optional
	Retention string `json:"retention,omitempty"`
}

// ObservabilityStatus shows the observed state of platform observability.
type ObservabilityStatus struct {
	// PipelineReady indicates whether the observability pipeline cluster is healthy.
	// +optional
	PipelineReady bool `json:"pipelineReady,omitempty"`

	// EnrolledCount is the number of tenant clusters with observability agents installed.
	// +optional
	EnrolledCount int32 `json:"enrolledCount,omitempty"`

	// TotalCount is the total number of tenant clusters.
	// +optional
	TotalCount int32 `json:"totalCount,omitempty"`
}
