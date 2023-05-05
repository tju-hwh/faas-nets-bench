/*
Copyright 2019-2021 OpenFaaS Authors

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	openfaasv1 "github.com/openfaas/faas-netes/pkg/apis/openfaas/v1"
	versioned "github.com/openfaas/faas-netes/pkg/client/clientset/versioned"
	internalinterfaces "github.com/openfaas/faas-netes/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/openfaas/faas-netes/pkg/client/listers/openfaas/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ProfileInformer provides access to a shared informer and lister for
// Profiles.
type ProfileInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ProfileLister
}

type profileInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewProfileInformer constructs a new informer for Profile type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewProfileInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredProfileInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredProfileInformer constructs a new informer for Profile type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredProfileInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OpenfaasV1().Profiles(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OpenfaasV1().Profiles(namespace).Watch(context.TODO(), options)
			},
		},
		&openfaasv1.Profile{},
		resyncPeriod,
		indexers,
	)
}

func (f *profileInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredProfileInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *profileInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&openfaasv1.Profile{}, f.defaultInformer)
}

func (f *profileInformer) Lister() v1.ProfileLister {
	return v1.NewProfileLister(f.Informer().GetIndexer())
}
