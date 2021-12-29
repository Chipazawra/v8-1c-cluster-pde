package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetClusterManagersRequest)(nil)

// GetClusterManagersRequest получение списка менеджеров кластера
//
//  type GET_CLUSTER_MANAGERS_REQUEST = 19
//  kind MESSAGE_KIND = 1
//  respond GetClusterManagersResponse
type GetClusterManagersRequest struct {
	ClusterID uuid.UUID
}

func (r *GetClusterManagersRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetClusterManagersRequest) Type() EndpointMessageType {
	return GET_CLUSTER_MANAGERS_REQUEST
}

func (r *GetClusterManagersRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)

}

// GetClusterManagersResponse содержит список менеджеров кластера
//  type GET_CLUSTER_MANAGERS_RESPONSE = 20
//  Managers serialize.ManagerInfo
type GetClusterManagersResponse struct {
	Managers []*serialize.ManagerInfo
}

func (res *GetClusterManagersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ManagerInfo{}
		info.Parse(decoder, version, r)

		res.Managers = append(res.Managers, info)
	}
}

func (_ *GetClusterManagersResponse) Type() EndpointMessageType {
	return GET_CLUSTER_MANAGERS_RESPONSE
}
