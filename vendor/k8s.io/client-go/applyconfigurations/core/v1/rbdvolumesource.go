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

// RBDVolumeSourceApplyConfiguration represents an declarative configuration of the RBDVolumeSource type for use
// with apply.
type RBDVolumeSourceApplyConfiguration struct {
	CephMonitors []string                                `json:"monitors,omitempty"`
	RBDImage     *string                                 `json:"image,omitempty"`
	FSType       *string                                 `json:"fsType,omitempty"`
	RBDPool      *string                                 `json:"pool,omitempty"`
	RadosUser    *string                                 `json:"user,omitempty"`
	Keyring      *string                                 `json:"keyring,omitempty"`
	SecretRef    *LocalObjectReferenceApplyConfiguration `json:"secretRef,omitempty"`
	ReadOnly     *bool                                   `json:"readOnly,omitempty"`
}

// RBDVolumeSourceApplyConfiguration constructs an declarative configuration of the RBDVolumeSource type for use with
// apply.
func RBDVolumeSource() *RBDVolumeSourceApplyConfiguration {
	return &RBDVolumeSourceApplyConfiguration{}
}

// WithCephMonitors adds the given value to the CephMonitors field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the CephMonitors field.
func (b *RBDVolumeSourceApplyConfiguration) WithCephMonitors(values ...string) *RBDVolumeSourceApplyConfiguration {
	for i := range values {
		b.CephMonitors = append(b.CephMonitors, values[i])
	}
	return b
}

// WithRBDImage sets the RBDImage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RBDImage field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithRBDImage(value string) *RBDVolumeSourceApplyConfiguration {
	b.RBDImage = &value
	return b
}

// WithFSType sets the FSType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FSType field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithFSType(value string) *RBDVolumeSourceApplyConfiguration {
	b.FSType = &value
	return b
}

// WithRBDPool sets the RBDPool field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RBDPool field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithRBDPool(value string) *RBDVolumeSourceApplyConfiguration {
	b.RBDPool = &value
	return b
}

// WithRadosUser sets the RadosUser field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RadosUser field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithRadosUser(value string) *RBDVolumeSourceApplyConfiguration {
	b.RadosUser = &value
	return b
}

// WithKeyring sets the Keyring field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Keyring field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithKeyring(value string) *RBDVolumeSourceApplyConfiguration {
	b.Keyring = &value
	return b
}

// WithSecretRef sets the SecretRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretRef field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithSecretRef(value *LocalObjectReferenceApplyConfiguration) *RBDVolumeSourceApplyConfiguration {
	b.SecretRef = value
	return b
}

// WithReadOnly sets the ReadOnly field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadOnly field is set to the value of the last call.
func (b *RBDVolumeSourceApplyConfiguration) WithReadOnly(value bool) *RBDVolumeSourceApplyConfiguration {
	b.ReadOnly = &value
	return b
}
