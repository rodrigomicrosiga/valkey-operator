package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValkeyInstanceSpec define a configuração desejada para a instância do Valkey
type ValkeyInstanceSpec struct {
	// Image define a imagem do Valkey a ser utilizada
	// +kubebuilder:default:="valkey/valkey:8.0"
	// +optional
	Image string `json:"image,omitempty"`

	// Replicas define o número de pods do Valkey (Alta Disponibilidade)
	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// Port define a porta de comunicação do Valkey
	// +kubebuilder:default:=6379
	// +optional
	Port int32 `json:"port,omitempty"`

	// StorageSize define o tamanho do disco persistente (PVC)
	// +kubebuilder:default:="2Gi"
	// +optional
	StorageSize string `json:"storageSize,omitempty"`
}

// ValkeyInstanceStatus define o estado observado pelo Operator (Status)
type ValkeyInstanceStatus struct {
	// Phase representa o estado atual da instância (ex: Provisioning, Ready)
	// +optional
	Phase string `json:"phase,omitempty"`

	// Endpoint é a URL interna de conexão para as aplicações no cluster
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Conditions armazenam o histórico de transições de estado para auditoria
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.endpoint"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ValkeyInstance is the Schema for the valkeyinstances API
type ValkeyInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ValkeyInstanceSpec   `json:"spec,omitempty"`
	Status ValkeyInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ValkeyInstanceList contains a list of ValkeyInstance
type ValkeyInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ValkeyInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ValkeyInstance{}, &ValkeyInstanceList{})
}
