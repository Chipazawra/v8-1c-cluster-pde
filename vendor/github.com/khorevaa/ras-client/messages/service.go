package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetClusterServicesRequest)(nil)

// GetClusterServicesRequest получение списка сервисов кластера
//
//  type GET_CLUSTER_SERVICES_REQUEST = 38
//  kind MESSAGE_KIND = 1
//  respond GetClusterServicesResponse
type GetClusterServicesRequest struct {
	ClusterID uuid.UUID
	response  *GetClusterServicesResponse
}

func (r *GetClusterServicesRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *GetClusterServicesRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ *GetClusterServicesRequest) Type() EndpointMessageType {
	return GET_CLUSTER_SERVICES_REQUEST
}

// GetClusterServicesResponse содержит список сервисов кластера
//  type GET_CLUSTER_SERVICES_RESPONSE = 37
//  Services serialize.ManagerInfo
type GetClusterServicesResponse struct {
	Services []*serialize.ServiceInfo
}

func (res *GetClusterServicesResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ServiceInfo{}
		info.Parse(decoder, version, r)
		res.Services = append(res.Services, info)
	}

}

func (_ *GetClusterServicesResponse) Type() EndpointMessageType {
	return GET_CLUSTER_SERVICES_RESPONSE
}
