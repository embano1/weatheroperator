package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WeatherServiceSpec defines the desired state of WeatherService
type WeatherServiceSpec struct {
	City string `json:"city"`
	Unit string `json:"unit"`
}

// WeatherServiceStatus defines the observed state of WeatherService
type WeatherServiceStatus struct {
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WeatherService is the Schema for the weatherservices API
// +k8s:openapi-gen=true
type WeatherService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WeatherServiceSpec   `json:"spec"`
	Status WeatherServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WeatherServiceList contains a list of WeatherService
type WeatherServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WeatherService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WeatherService{}, &WeatherServiceList{})
}
