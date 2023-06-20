/*
Copyright 2022 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// RecordParameters are the configurable fields of a Record.
type RecordParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// RecordObservation are the observable fields of a Record.
type RecordObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A RecordSpec defines the desired state of a Record.
type RecordSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RecordParameters `json:"forProvider"`
}

// A RecordStatus represents the observed state of a Record.
type RecordStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RecordObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Record is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,simplejsonapp}
type Record struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordSpec   `json:"spec"`
	Status RecordStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RecordList contains a list of Record
type RecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Record `json:"items"`
}

// Record type metadata.
var (
	RecordKind             = reflect.TypeOf(Record{}).Name()
	RecordGroupKind        = schema.GroupKind{Group: Group, Kind: RecordKind}.String()
	RecordKindAPIVersion   = RecordKind + "." + SchemeGroupVersion.String()
	RecordGroupVersionKind = SchemeGroupVersion.WithKind(RecordKind)
)

func init() {
	SchemeBuilder.Register(&Record{}, &RecordList{})
}
