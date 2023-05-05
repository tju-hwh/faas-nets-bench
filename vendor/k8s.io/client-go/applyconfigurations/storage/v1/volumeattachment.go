/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	apistoragev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	managedfields "k8s.io/apimachinery/pkg/util/managedfields"
	internal "k8s.io/client-go/applyconfigurations/internal"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// VolumeAttachmentApplyConfiguration represents an declarative configuration of the VolumeAttachment type for use
// with apply.
type VolumeAttachmentApplyConfiguration struct {
	v1.TypeMetaApplyConfiguration    `json:",inline"`
	*v1.ObjectMetaApplyConfiguration `json:"metadata,omitempty"`
	Spec                             *VolumeAttachmentSpecApplyConfiguration   `json:"spec,omitempty"`
	Status                           *VolumeAttachmentStatusApplyConfiguration `json:"status,omitempty"`
}

// VolumeAttachment constructs an declarative configuration of the VolumeAttachment type for use with
// apply.
func VolumeAttachment(name string) *VolumeAttachmentApplyConfiguration {
	b := &VolumeAttachmentApplyConfiguration{}
	b.WithName(name)
	b.WithKind("VolumeAttachment")
	b.WithAPIVersion("storage.k8s.io/v1")
	return b
}

// ExtractVolumeAttachment extracts the applied configuration owned by fieldManager from
// volumeAttachment. If no managedFields are found in volumeAttachment for fieldManager, a
// VolumeAttachmentApplyConfiguration is returned with only the Name, Namespace (if applicable),
// APIVersion and Kind populated. It is possible that no managed fields were found for because other
// field managers have taken ownership of all the fields previously owned by fieldManager, or because
// the fieldManager never owned fields any fields.
// volumeAttachment must be a unmodified VolumeAttachment API object that was retrieved from the Kubernetes API.
// ExtractVolumeAttachment provides a way to perform a extract/modify-in-place/apply workflow.
// Note that an extracted apply configuration will contain fewer fields than what the fieldManager previously
// applied if another fieldManager has updated or force applied any of the previously applied fields.
// Experimental!
func ExtractVolumeAttachment(volumeAttachment *apistoragev1.VolumeAttachment, fieldManager string) (*VolumeAttachmentApplyConfiguration, error) {
	return extractVolumeAttachment(volumeAttachment, fieldManager, "")
}

// ExtractVolumeAttachmentStatus is the same as ExtractVolumeAttachment except
// that it extracts the status subresource applied configuration.
// Experimental!
func ExtractVolumeAttachmentStatus(volumeAttachment *apistoragev1.VolumeAttachment, fieldManager string) (*VolumeAttachmentApplyConfiguration, error) {
	return extractVolumeAttachment(volumeAttachment, fieldManager, "status")
}

func extractVolumeAttachment(volumeAttachment *apistoragev1.VolumeAttachment, fieldManager string, subresource string) (*VolumeAttachmentApplyConfiguration, error) {
	b := &VolumeAttachmentApplyConfiguration{}
	err := managedfields.ExtractInto(volumeAttachment, internal.Parser().Type("io.k8s.api.storage.v1.VolumeAttachment"), fieldManager, b, subresource)
	if err != nil {
		return nil, err
	}
	b.WithName(volumeAttachment.Name)

	b.WithKind("VolumeAttachment")
	b.WithAPIVersion("storage.k8s.io/v1")
	return b, nil
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithKind(value string) *VolumeAttachmentApplyConfiguration {
	b.Kind = &value
	return b
}

// WithAPIVersion sets the APIVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the APIVersion field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithAPIVersion(value string) *VolumeAttachmentApplyConfiguration {
	b.APIVersion = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithName(value string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithGenerateName(value string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithNamespace(value string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Namespace = &value
	return b
}

// WithSelfLink sets the SelfLink field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SelfLink field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithSelfLink(value string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.SelfLink = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithUID(value types.UID) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.UID = &value
	return b
}

// WithResourceVersion sets the ResourceVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceVersion field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithResourceVersion(value string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ResourceVersion = &value
	return b
}

// WithGeneration sets the Generation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Generation field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithGeneration(value int64) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Generation = &value
	return b
}

// WithCreationTimestamp sets the CreationTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CreationTimestamp field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithCreationTimestamp(value metav1.Time) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.CreationTimestamp = &value
	return b
}

// WithDeletionTimestamp sets the DeletionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionTimestamp field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithDeletionTimestamp(value metav1.Time) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionTimestamp = &value
	return b
}

// WithDeletionGracePeriodSeconds sets the DeletionGracePeriodSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionGracePeriodSeconds field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithDeletionGracePeriodSeconds(value int64) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionGracePeriodSeconds = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *VolumeAttachmentApplyConfiguration) WithLabels(entries map[string]string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *VolumeAttachmentApplyConfiguration) WithAnnotations(entries map[string]string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithOwnerReferences adds the given value to the OwnerReferences field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OwnerReferences field.
func (b *VolumeAttachmentApplyConfiguration) WithOwnerReferences(values ...*v1.OwnerReferenceApplyConfiguration) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOwnerReferences")
		}
		b.OwnerReferences = append(b.OwnerReferences, *values[i])
	}
	return b
}

// WithFinalizers adds the given value to the Finalizers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Finalizers field.
func (b *VolumeAttachmentApplyConfiguration) WithFinalizers(values ...string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		b.Finalizers = append(b.Finalizers, values[i])
	}
	return b
}

// WithClusterName sets the ClusterName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClusterName field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithClusterName(value string) *VolumeAttachmentApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ClusterName = &value
	return b
}

func (b *VolumeAttachmentApplyConfiguration) ensureObjectMetaApplyConfigurationExists() {
	if b.ObjectMetaApplyConfiguration == nil {
		b.ObjectMetaApplyConfiguration = &v1.ObjectMetaApplyConfiguration{}
	}
}

// WithSpec sets the Spec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Spec field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithSpec(value *VolumeAttachmentSpecApplyConfiguration) *VolumeAttachmentApplyConfiguration {
	b.Spec = value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *VolumeAttachmentApplyConfiguration) WithStatus(value *VolumeAttachmentStatusApplyConfiguration) *VolumeAttachmentApplyConfiguration {
	b.Status = value
	return b
}
