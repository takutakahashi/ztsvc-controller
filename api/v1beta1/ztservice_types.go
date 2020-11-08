package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZTServiceSpec defines the desired state of ZTService
type ZTServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ZTService. Edit ZTService_types.go to remove/update
	Service    *corev1.Service `json:"service,omitempty"`
	ServiceRef string          `json:"serviceRef,omitempty"`
}

// ZTServiceStatus defines the observed state of ZTService
type ZTServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	IP      string `json:"ip,omitempty"`
	Address string `json:"address,omitempty"`
}

// +kubebuilder:object:root=true

// ZTService is the Schema for the ztservices API
type ZTService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZTServiceSpec   `json:"spec,omitempty"`
	Status ZTServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ZTServiceList contains a list of ZTService
type ZTServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ZTService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ZTService{}, &ZTServiceList{})
}
