package jsongrpctranscoder

import (
	"encoding/base64"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/api/v1/plugins/jsongrpctranscoder"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	glooplugins "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/util"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

//go:generate protoc -I$GOPATH/src/github.com/lyft/protoc-gen-validate -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/ filter.proto

const (
	filterName = "io.solo.filters.http.json_grpc_transcoder"

	// TODO(talnordan)
	pluginStage = plugins.FaultFilter
)

type plugin struct {
	descriptorSetPerUpstream map[string]*descriptor.FileDescriptorSet
}

func NewPlugin() plugins.Plugin {
	return &plugin{}
}

func (p *plugin) Init(params plugins.InitParams) error {
	return nil
}

// TODO(talnordan): Duplicated from `plugins/grpc/plugin.go`.
func convertProto(encodedBytes []byte) (*descriptor.FileDescriptorSet, error) {
	// base-64 encoded by function discovery
	rawDescriptors, err := base64.StdEncoding.DecodeString(string(encodedBytes))
	if err != nil {
		return nil, err
	}
	var fileDescriptor descriptor.FileDescriptorSet
	if err := proto.Unmarshal(rawDescriptors, &fileDescriptor); err != nil {
		return nil, err
	}
	return &fileDescriptor, nil
}

// TODO(talnordan): What if there is no gRPC upstream with this gRPC spec?
func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *envoyapi.Cluster) error {
	upstreamType, ok := in.UpstreamSpec.UpstreamType.(v1.ServiceSpecGetter)
	if !ok {
		return nil
	}

	if upstreamType.GetServiceSpec() == nil {
		return nil
	}

	grpcWrapper, ok := upstreamType.GetServiceSpec().PluginType.(*glooplugins.ServiceSpec_Grpc)
	if !ok {
		return nil
	}
	grpcSpec := grpcWrapper.Grpc

	descriptors, err := convertProto(grpcSpec.Descriptors)
	if err != nil {
		return errors.Wrapf(err, "parsing grpc spec as a proto descriptor set")
	}

	p.descriptorSetPerUpstream[in.Metadata.Name] = descriptors

	return nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if len(p.descriptorSetPerUpstream) == 0 {
		return nil, nil
	}

	var filters []plugins.StagedHttpFilter
	for _, descriptorSet := range p.descriptorSetPerUpstream {
		descriptorBytes, err := proto.Marshal(descriptorSet)
		if err != nil {
			return nil, errors.Wrapf(err, "marshaling proto descriptor")
		}

		filterConfig, err := util.MessageToStruct(&jsongrpctranscoder.JsonGrpcTranscoder{
			DescriptorSet: &envoytranscoder.JsonGrpcTranscoder_ProtoDescriptorBin{
				ProtoDescriptorBin: descriptorBytes,
			},
		})

		if err != nil {
			return nil, errors.Wrapf(err, "ERROR: marshaling JsonGrpcTranscoder config")
		}

		filters = append(filters, plugins.StagedHttpFilter{
			HttpFilter: &envoyhttp.HttpFilter{
				Name:   filterName,
				Config: filterConfig,
			},
			Stage: pluginStage,
		})
	}

	if len(filters) == 0 {
		return nil, errors.Errorf("ERROR: no valid JsonGrpcTranscoder available")
	}

	return filters, nil
}

// plugins.ListenerPlugin
func (p *plugin) ProcessListener(params Params, in *v1.Listener, out *envoyapi.Listener) error {
	// TODO(talnordan)
	httpListener, ok := in.ListenerType.(*v1.Listener_HttpListener)
	if !ok {
		// TODO(talnordan): err
		return nil
	}

	plugins := httpListener.ListenerPlugins
}
