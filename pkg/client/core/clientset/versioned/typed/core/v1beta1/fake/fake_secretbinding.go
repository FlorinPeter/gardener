/*
Copyright (c) SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSecretBindings implements SecretBindingInterface
type FakeSecretBindings struct {
	Fake *FakeCoreV1beta1
	ns   string
}

var secretbindingsResource = schema.GroupVersionResource{Group: "core.gardener.cloud", Version: "v1beta1", Resource: "secretbindings"}

var secretbindingsKind = schema.GroupVersionKind{Group: "core.gardener.cloud", Version: "v1beta1", Kind: "SecretBinding"}

// Get takes name of the secretBinding, and returns the corresponding secretBinding object, and an error if there is any.
func (c *FakeSecretBindings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(secretbindingsResource, c.ns, name), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}

// List takes label and field selectors, and returns the list of SecretBindings that match those selectors.
func (c *FakeSecretBindings) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.SecretBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(secretbindingsResource, secretbindingsKind, c.ns, opts), &v1beta1.SecretBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.SecretBindingList{ListMeta: obj.(*v1beta1.SecretBindingList).ListMeta}
	for _, item := range obj.(*v1beta1.SecretBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested secretBindings.
func (c *FakeSecretBindings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(secretbindingsResource, c.ns, opts))

}

// Create takes the representation of a secretBinding and creates it.  Returns the server's representation of the secretBinding, and an error, if there is any.
func (c *FakeSecretBindings) Create(ctx context.Context, secretBinding *v1beta1.SecretBinding, opts v1.CreateOptions) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(secretbindingsResource, c.ns, secretBinding), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}

// Update takes the representation of a secretBinding and updates it. Returns the server's representation of the secretBinding, and an error, if there is any.
func (c *FakeSecretBindings) Update(ctx context.Context, secretBinding *v1beta1.SecretBinding, opts v1.UpdateOptions) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(secretbindingsResource, c.ns, secretBinding), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}

// Delete takes name of the secretBinding and deletes it. Returns an error if one occurs.
func (c *FakeSecretBindings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(secretbindingsResource, c.ns, name, opts), &v1beta1.SecretBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSecretBindings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(secretbindingsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.SecretBindingList{})
	return err
}

// Patch applies the patch and returns the patched secretBinding.
func (c *FakeSecretBindings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(secretbindingsResource, c.ns, name, pt, data, subresources...), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}
