/*


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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EchoSpec defines the desired state of Echo
type EchoSpec struct {
	Message string `json:"message"`
}

// EchoStatus defines the observed state of Echo
type EchoStatus struct {
	Message string `json:"message"`
	// The generation observed by the echo controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Represents the latest available observations of an echo's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []EchoCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Echo is the Schema for the echoes API
type Echo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EchoSpec   `json:"spec,omitempty"`
	Status EchoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EchoList contains a list of Echo
type EchoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Echo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Echo{}, &EchoList{})
}

type EchoCondition struct {
	// Type of deployment condition.
	Type string `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
}
