/*
Copyright 2022 SAP SE.

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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VropsExporterSpec defines the desired state of VropsExporter
type VropsExporterSpec struct {
	Name              string                              `json:"name"`
	Prometheus        string                              `json:"prometheus,omitempty"`
	ScrapeInterval    int                                 `json:"scrapeInterval,omitempty"`
	ScrapeTimeout     int                                 `json:"scrapeTimeout,omitempty"`
	Namespace         string                              `json:"namespace,omitempty"`
	Image             string                              `json:"image"`
	Port              int32                               `json:"port"`
	User              string                              `json:"user"`
	Password          string                              `json:"password"`
	Debug             string                              `json:"debug,omitempty"`
	Target            string                              `json:"target"`
	Inventory         *VropsExporterInventorySpec         `json:"inventory"`
	InventoryExporter *VropsExporterInventoryExporterSpec `json:"inventory-exporter"`
	ExporterTypes     []ExporterType                      `json:"exporter-types"`
}

// VropsExporterStatus defines the observed state of VropsExporter
type VropsExporterStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VropsExporter is the Schema for the vropsexporters API
type VropsExporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VropsExporterSpec   `json:"spec,omitempty"`
	Status VropsExporterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VropsExporterList contains a list of VropsExporter
type VropsExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VropsExporter `json:"items"`
}

// VropsExporterInventorySpec defines the desired state of VropsExporterInventory
type VropsExporterInventorySpec struct {
	Tag       string `json:"tag,omitempty"`
	Port      int    `json:"port,omitempty"`
	Sleep     int    `json:"sleep,omitempty"`
	Timeout   string `json:"timeout,omitempty"`
	Image     string `json:"image,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	// Define resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

// VropsExporterInventoryExporterSpec defines the desired state of VropsExporterInventoryExporter
type VropsExporterInventoryExporterSpec struct {
	Port int `json:"port,omitempty"`
}

func init() {
	SchemeBuilder.Register(&VropsExporter{}, &VropsExporterList{})
}

// ExporterType defines the desired state of ExporterTypes
type ExporterType struct {
	Name       string                  `json:"name"`
	Collectors []string                `json:"collectors"`
	Resources  v1.ResourceRequirements `json:"resources,omitempty"`
}
