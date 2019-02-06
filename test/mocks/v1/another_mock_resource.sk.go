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
func NewAnotherMockResource(namespace, name string) *AnotherMockResource {
	return &AnotherMockResource{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *AnotherMockResource) SetStatus(status core.Status) {
	r.Status = status
}

func (r *AnotherMockResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *AnotherMockResource) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	return hashutils.HashAll(
		metaCopy,
		r.BasicField,
	)
}

type AnotherMockResourceList []*AnotherMockResource
type AnothermockresourcesByNamespace map[string]AnotherMockResourceList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list AnotherMockResourceList) Find(namespace, name string) (*AnotherMockResource, error) {
	for _, anotherMockResource := range list {
		if anotherMockResource.Metadata.Name == name {
			if namespace == "" || anotherMockResource.Metadata.Namespace == namespace {
				return anotherMockResource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find anotherMockResource %v.%v", namespace, name)
}

func (list AnotherMockResourceList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, anotherMockResource := range list {
		ress = append(ress, anotherMockResource)
	}
	return ress
}

func (list AnotherMockResourceList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, anotherMockResource := range list {
		ress = append(ress, anotherMockResource)
	}
	return ress
}

func (list AnotherMockResourceList) Names() []string {
	var names []string
	for _, anotherMockResource := range list {
		names = append(names, anotherMockResource.Metadata.Name)
	}
	return names
}

func (list AnotherMockResourceList) NamespacesDotNames() []string {
	var names []string
	for _, anotherMockResource := range list {
		names = append(names, anotherMockResource.Metadata.Namespace+"."+anotherMockResource.Metadata.Name)
	}
	return names
}

func (list AnotherMockResourceList) Sort() AnotherMockResourceList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list AnotherMockResourceList) Clone() AnotherMockResourceList {
	var anotherMockResourceList AnotherMockResourceList
	for _, anotherMockResource := range list {
		anotherMockResourceList = append(anotherMockResourceList, proto.Clone(anotherMockResource).(*AnotherMockResource))
	}
	return anotherMockResourceList
}

func (list AnotherMockResourceList) Each(f func(element *AnotherMockResource)) {
	for _, anotherMockResource := range list {
		f(anotherMockResource)
	}
}

func (list AnotherMockResourceList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *AnotherMockResource) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

func (byNamespace AnothermockresourcesByNamespace) Add(anotherMockResource ...*AnotherMockResource) {
	for _, item := range anotherMockResource {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace AnothermockresourcesByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace AnothermockresourcesByNamespace) List() AnotherMockResourceList {
	var list AnotherMockResourceList
	for _, anotherMockResourceList := range byNamespace {
		list = append(list, anotherMockResourceList...)
	}
	return list.Sort()
}

func (byNamespace AnothermockresourcesByNamespace) Clone() AnothermockresourcesByNamespace {
	cloned := make(AnothermockresourcesByNamespace)
	for ns, list := range byNamespace {
		cloned[ns] = list.Clone()
	}
	return cloned
}

var _ resources.Resource = &AnotherMockResource{}

// Kubernetes Adapter for AnotherMockResource

func (o *AnotherMockResource) GetObjectKind() schema.ObjectKind {
	t := AnotherMockResourceCrd.TypeMeta()
	return &t
}

func (o *AnotherMockResource) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*AnotherMockResource)
}

var AnotherMockResourceCrd = crd.NewCrd("testing.solo.io",
	"anothermockresources",
	"testing.solo.io",
	"v1",
	"AnotherMockResource",
	"amr",
	false,
	&AnotherMockResource{})
