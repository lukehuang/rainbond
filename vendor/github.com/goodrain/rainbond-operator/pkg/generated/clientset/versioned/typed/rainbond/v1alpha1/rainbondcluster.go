// RAINBOND, Application Management Platform
// Copyright (C) 2014-2020 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/goodrain/rainbond-operator/pkg/apis/rainbond/v1alpha1"
	scheme "github.com/goodrain/rainbond-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RainbondClustersGetter has a method to return a RainbondClusterInterface.
// A group's client should implement this interface.
type RainbondClustersGetter interface {
	RainbondClusters(namespace string) RainbondClusterInterface
}

// RainbondClusterInterface has methods to work with RainbondCluster resources.
type RainbondClusterInterface interface {
	Create(*v1alpha1.RainbondCluster) (*v1alpha1.RainbondCluster, error)
	Update(*v1alpha1.RainbondCluster) (*v1alpha1.RainbondCluster, error)
	UpdateStatus(*v1alpha1.RainbondCluster) (*v1alpha1.RainbondCluster, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.RainbondCluster, error)
	List(opts v1.ListOptions) (*v1alpha1.RainbondClusterList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RainbondCluster, err error)
	RainbondClusterExpansion
}

// rainbondClusters implements RainbondClusterInterface
type rainbondClusters struct {
	client rest.Interface
	ns     string
}

// newRainbondClusters returns a RainbondClusters
func newRainbondClusters(c *RainbondV1alpha1Client, namespace string) *rainbondClusters {
	return &rainbondClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the rainbondCluster, and returns the corresponding rainbondCluster object, and an error if there is any.
func (c *rainbondClusters) Get(name string, options v1.GetOptions) (result *v1alpha1.RainbondCluster, err error) {
	result = &v1alpha1.RainbondCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rainbondclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RainbondClusters that match those selectors.
func (c *rainbondClusters) List(opts v1.ListOptions) (result *v1alpha1.RainbondClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.RainbondClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rainbondclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested rainbondClusters.
func (c *rainbondClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("rainbondclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a rainbondCluster and creates it.  Returns the server's representation of the rainbondCluster, and an error, if there is any.
func (c *rainbondClusters) Create(rainbondCluster *v1alpha1.RainbondCluster) (result *v1alpha1.RainbondCluster, err error) {
	result = &v1alpha1.RainbondCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("rainbondclusters").
		Body(rainbondCluster).
		Do().
		Into(result)
	return
}

// Update takes the representation of a rainbondCluster and updates it. Returns the server's representation of the rainbondCluster, and an error, if there is any.
func (c *rainbondClusters) Update(rainbondCluster *v1alpha1.RainbondCluster) (result *v1alpha1.RainbondCluster, err error) {
	result = &v1alpha1.RainbondCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rainbondclusters").
		Name(rainbondCluster.Name).
		Body(rainbondCluster).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *rainbondClusters) UpdateStatus(rainbondCluster *v1alpha1.RainbondCluster) (result *v1alpha1.RainbondCluster, err error) {
	result = &v1alpha1.RainbondCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rainbondclusters").
		Name(rainbondCluster.Name).
		SubResource("status").
		Body(rainbondCluster).
		Do().
		Into(result)
	return
}

// Delete takes name of the rainbondCluster and deletes it. Returns an error if one occurs.
func (c *rainbondClusters) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rainbondclusters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rainbondClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rainbondclusters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched rainbondCluster.
func (c *rainbondClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RainbondCluster, err error) {
	result = &v1alpha1.RainbondCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("rainbondclusters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
