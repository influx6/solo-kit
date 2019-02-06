// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewFakeResource(namespace, name string) *FakeResource {
	return &FakeResource{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *FakeResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *FakeResource) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	return hashutils.HashAll(
		metaCopy,
		r.Count,
	)
}

type FakeResourceList []*FakeResource
type FakesByNamespace map[string]FakeResourceList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list FakeResourceList) Find(namespace, name string) (*FakeResource, error) {
	for _, fakeResource := range list {
		if fakeResource.Metadata.Name == name {
			if namespace == "" || fakeResource.Metadata.Namespace == namespace {
				return fakeResource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find fakeResource %v.%v", namespace, name)
}

func (list FakeResourceList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, fakeResource := range list {
		ress = append(ress, fakeResource)
	}
	return ress
}

func (list FakeResourceList) Names() []string {
	var names []string
	for _, fakeResource := range list {
		names = append(names, fakeResource.Metadata.Name)
	}
	return names
}

func (list FakeResourceList) NamespacesDotNames() []string {
	var names []string
	for _, fakeResource := range list {
		names = append(names, fakeResource.Metadata.Namespace+"."+fakeResource.Metadata.Name)
	}
	return names
}

func (list FakeResourceList) Sort() FakeResourceList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list FakeResourceList) Clone() FakeResourceList {
	var fakeResourceList FakeResourceList
	for _, fakeResource := range list {
		fakeResourceList = append(fakeResourceList, proto.Clone(fakeResource).(*FakeResource))
	}
	return fakeResourceList
}

func (list FakeResourceList) Each(f func(element *FakeResource)) {
	for _, fakeResource := range list {
		f(fakeResource)
	}
}

func (list FakeResourceList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *FakeResource) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

func (byNamespace FakesByNamespace) Add(fakeResource ...*FakeResource) {
	for _, item := range fakeResource {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace FakesByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace FakesByNamespace) List() FakeResourceList {
	var list FakeResourceList
	for _, fakeResourceList := range byNamespace {
		list = append(list, fakeResourceList...)
	}
	return list.Sort()
}

func (byNamespace FakesByNamespace) Clone() FakesByNamespace {
	cloned := make(FakesByNamespace)
	for ns, list := range byNamespace {
		cloned[ns] = list.Clone()
	}
	return cloned
}

var _ resources.Resource = &FakeResource{}

// Kubernetes Adapter for FakeResource

func (o *FakeResource) GetObjectKind() schema.ObjectKind {
	t := FakeResourceCrd.TypeMeta()
	return &t
}

func (o *FakeResource) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*FakeResource)
}

var FakeResourceCrd = crd.NewCrd("testing.solo.io",
	"fakes",
	"testing.solo.io",
	"v1",
	"FakeResource",
	"fk",
	false,
	&FakeResource{})
