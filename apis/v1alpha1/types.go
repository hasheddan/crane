/*
Copyright 2020 The Crossplane Authors.

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

// A JobSpec defines the desired state of a Job.
type JobSpec struct {
	Provider string `json:"provider"`
	Repo     string `json:"repo"`
}

// +kubebuilder:object:root=true

// A Job configures a crane Job.
// +kubebuilder:printcolumn:name="PROVIDER",type="string",JSONPath=".spec.provider"
// +kubebuilder:printcolumn:name="REPO",type="string",JSONPath=".spec.repo"
// +kubebuilder:resource:scope=Cluster
type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec JobSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// JobList contains a list of Job
type JobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Job `json:"items"`
}
