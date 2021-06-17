/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/errors"
)

// Constants block.
const (
	// MachinePoolFinalizer is the finalizer for the machine pool.
	MachinePoolFinalizer = "awsmachinepool.infrastructure.cluster.x-k8s.io"

	// LaunchTemplateLatestVersion defines the launching of the latest version of the template.
	LaunchTemplateLatestVersion = "$Latest"
)

// AWSMachinePoolSpec defines the desired state of AWSMachinePool
type AWSMachinePoolSpec struct {
	// ProviderID is the ARN of the associated ASG
	// +optional
	ProviderID string `json:"providerID,omitempty"`

	// MinSize defines the minimum size of the group.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	MinSize int32 `json:"minSize"`

	// MaxSize defines the maximum size of the group.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	MaxSize int32 `json:"maxSize"`

	// AvailabilityZones is an array of availability zones instances can run in
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// Subnets is an array of subnet configurations
	// +optional
	Subnets []infrav1.AWSResourceReference `json:"subnets,omitempty"`

	// AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
	// AWS provider.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// AWSLaunchTemplate specifies the launch template and version to use when an instance is launched.
	// +kubebuilder:validation:Required
	AWSLaunchTemplate AWSLaunchTemplate `json:"awsLaunchTemplate"`

	// MixedInstancesPolicy describes how multiple instance types will be used by the ASG.
	MixedInstancesPolicy *MixedInstancesPolicy `json:"mixedInstancesPolicy,omitempty"`

	// ProviderIDList are the identification IDs of machine instances provided by the provider.
	// This field must match the provider IDs as seen on the node objects corresponding to a machine pool's machine instances.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// The amount of time, in seconds, after a scaling activity completes before another scaling activity can start.
	// If no value is supplied by user a default value of 300 seconds is set
	// +optional
	DefaultCoolDown metav1.Duration `json:"defaultCoolDown,omitempty"`

	// RefreshPreferences describes set of preferences associated with the instance refresh request.
	// +optional
	RefreshPreferences *RefreshPreferences `json:"refreshPreferences,omitempty"`

	// Enable or disable the capacity rebalance autoscaling group feature
	// +optional
	CapacityRebalance bool `json:"capacityRebalance,omitempty"`
}

// RefreshPreferences defines the specs for instance refreshing.
type RefreshPreferences struct {
	// The strategy to use for the instance refresh. The only valid value is Rolling.
	// A rolling update is an update that is applied to all instances in an Auto
	// Scaling group until all instances have been updated.
	// +optional
	Strategy *string `json:"strategy,omitempty"`

	// The number of seconds until a newly launched instance is configured and ready
	// to use. During this time, the next replacement will not be initiated.
	// The default is to use the value for the health check grace period defined for the group.
	// +optional
	InstanceWarmup *int64 `json:"instanceWarmup,omitempty"`

	// The amount of capacity as a percentage in ASG that must remain healthy
	// during an instance refresh. The default is 90.
	// +optional
	MinHealthyPercentage *int64 `json:"minHealthyPercentage,omitempty"`
}

// AWSMachinePoolStatus defines the observed state of AWSMachinePool
type AWSMachinePoolStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas
	// +optional
	Replicas int32 `json:"replicas"`

	// Conditions defines current service state of the AWSMachinePool.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// Instances contains the status for each instance in the pool
	// +optional
	Instances []*AWSMachinePoolInstanceStatus `json:"instances"`

	// The ID of the launch template
	LaunchTemplateID string `json:"launchTemplateID,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	ASGStatus *ASGStatus `json:"asgStatus,omitempty"`
}

// AWSMachinePoolInstanceStatus defines the status of the AWSMachinePoolInstance.
type AWSMachinePoolInstanceStatus struct {
	// InstanceID is the identification of the Machine Instance within ASG
	// +optional
	InstanceID string `json:"instanceID,omitempty"`

	// Version defines the Kubernetes version for the Machine Instance
	// +optional
	Version *string `json:"version,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=awsmachinepools,scope=Namespaced,categories=cluster-api,shortName=awsmp
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="Machine ready status"
// +kubebuilder:printcolumn:name="MinSize",type="integer",JSONPath=".spec.minSize",description="Minimum instanes in ASG"
// +kubebuilder:printcolumn:name="MaxSize",type="integer",JSONPath=".spec.maxSize",description="Maximum instanes in ASG"
// +kubebuilder:printcolumn:name="LaunchTemplate ID",type="string",JSONPath=".status.launchTemplateID",description="Launch Template ID"

// AWSMachinePool is the Schema for the awsmachinepools API
type AWSMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSMachinePoolSpec   `json:"spec,omitempty"`
	Status AWSMachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSMachinePoolList contains a list of AWSMachinePool.
type AWSMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSMachinePool{}, &AWSMachinePoolList{})
}

// GetConditions returns the observations of the operational state of the AWSMachinePool resource.
func (r *AWSMachinePool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSMachinePool to the predescribed clusterv1.Conditions.
func (r *AWSMachinePool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// GetObjectKind will return the ObjectKind of an AWSMachinePool.
func (r *AWSMachinePool) GetObjectKind() schema.ObjectKind {
	return &r.TypeMeta
}

// GetObjectKind will return the ObjectKind of an AWSMachinePoolList.
func (r *AWSMachinePoolList) GetObjectKind() schema.ObjectKind {
	return &r.TypeMeta
}
