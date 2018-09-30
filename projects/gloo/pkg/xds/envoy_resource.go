// Copyright 2018 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package xds

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/cache"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/util"
)

type EnvoyResource struct {
	resourceProto cache.ResourceProto
}

var _ cache.Resource = &EnvoyResource{}

func NewEnvoyResource(r cache.ResourceProto) *EnvoyResource {
	return &EnvoyResource{resourceProto: r}
}

// Resource types in xDS v2.
const (
	typePrefix   = cache.TypePrefix + "/envoy.api.v2."
	EndpointType = typePrefix + "ClusterLoadAssignment"
	ClusterType  = typePrefix + "Cluster"
	RouteType    = typePrefix + "RouteConfiguration"
	ListenerType = typePrefix + "Listener"
)

var (
	// ResponseTypes are supported response types.
	ResponseTypes = []string{
		EndpointType,
		ClusterType,
		RouteType,
		ListenerType,
	}
)

func (e *EnvoyResource) Self() cache.ResourceReference {
	return cache.ResourceReference{
		Name: e.Name(),
		Type: e.Type(),
	}
}

// GetResourceName returns the resource name for a valid xDS response type.
func (e *EnvoyResource) Name() string {
	switch v := e.resourceProto.(type) {
	case *v2.ClusterLoadAssignment:
		return v.GetClusterName()
	case *v2.Cluster:
		return v.GetName()
	case *v2.RouteConfiguration:
		return v.GetName()
	case *v2.Listener:
		return v.GetName()
	default:
		return ""
	}
}

func (e *EnvoyResource) ResourceProto() cache.ResourceProto {
	return e.resourceProto
}

func (e *EnvoyResource) Type() string {
	switch e.resourceProto.(type) {
	case *v2.ClusterLoadAssignment:
		return EndpointType
	case *v2.Cluster:
		return ClusterType
	case *v2.RouteConfiguration:
		return RouteType
	case *v2.Listener:
		return ListenerType
	default:
		return ""
	}
}

func (e *EnvoyResource) References() []cache.ResourceReference {
	out := make(map[cache.ResourceReference]bool)
	res := e.resourceProto
	if res == nil {
		return nil
	}
	switch v := res.(type) {
	case *v2.ClusterLoadAssignment:
		// no dependencies
	case *v2.Cluster:
		// for EDS type, use cluster name or ServiceName override
		if v.Type == v2.Cluster_EDS {
			rr := cache.ResourceReference{
				Type: EndpointType,
			}
			if v.EdsClusterConfig != nil && v.EdsClusterConfig.ServiceName != "" {
				rr.Name = v.EdsClusterConfig.ServiceName
			} else {
				rr.Name = v.Name
			}
			out[rr] = true
		}
	case *v2.RouteConfiguration:
		// References to clusters in both routes (and listeners) are not included
		// in the result, because the clusters are retrieved in bulk currently,
		// and not by name.
	case *v2.Listener:
		// extract route configuration names from HTTP connection manager
		for _, chain := range v.FilterChains {
			for _, filter := range chain.Filters {
				if filter.Name != util.HTTPConnectionManager {
					continue
				}

				config := &hcm.HttpConnectionManager{}
				if util.StructToMessage(filter.Config, config) == nil && config != nil {
					if rds, ok := config.RouteSpecifier.(*hcm.HttpConnectionManager_Rds); ok && rds != nil && rds.Rds != nil {
						rr := cache.ResourceReference{
							Type: RouteType,
						}
						rr.Name = rds.Rds.RouteConfigName
						out[rr] = true
					}
				}
			}
		}
	}

	var references []cache.ResourceReference
	for k, _ := range out {
		references = append(references, k)
	}
	return references
}

// GetResourceReferences returns the names for dependent resources (EDS cluster
// names for CDS, RDS routes names for LDS).
func GetResourceReferences(resources map[string]cache.Resource) map[string]bool {
	out := make(map[string]bool)
	for _, res := range resources {
		if res == nil {
			continue
		}
		switch v := res.ResourceProto().(type) {
		case *v2.ClusterLoadAssignment:
			// no dependencies
		case *v2.Cluster:
			// for EDS type, use cluster name or ServiceName override
			if v.Type == v2.Cluster_EDS {
				if v.EdsClusterConfig != nil && v.EdsClusterConfig.ServiceName != "" {
					out[v.EdsClusterConfig.ServiceName] = true
				} else {
					out[v.Name] = true
				}
			}
		case *v2.RouteConfiguration:
			// References to clusters in both routes (and listeners) are not included
			// in the result, because the clusters are retrieved in bulk currently,
			// and not by name.
		case *v2.Listener:
			// extract route configuration names from HTTP connection manager
			for _, chain := range v.FilterChains {
				for _, filter := range chain.Filters {
					if filter.Name != util.HTTPConnectionManager {
						continue
					}

					config := &hcm.HttpConnectionManager{}
					if util.StructToMessage(filter.Config, config) == nil && config != nil {
						if rds, ok := config.RouteSpecifier.(*hcm.HttpConnectionManager_Rds); ok && rds != nil && rds.Rds != nil {
							out[rds.Rds.RouteConfigName] = true
						}
					}
				}
			}
		}
	}
	return out
}

// GetResourceName returns the resource name for a valid xDS response type.
func GetResourceName(res cache.ResourceProto) string {
	switch v := res.(type) {
	case *v2.ClusterLoadAssignment:
		return v.GetClusterName()
	case *v2.Cluster:
		return v.GetName()
	case *v2.RouteConfiguration:
		return v.GetName()
	case *v2.Listener:
		return v.GetName()
	default:
		return ""
	}
}