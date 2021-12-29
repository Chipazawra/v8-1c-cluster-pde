package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetConnectionsShortRequest)(nil)

// GetConnectionsShortRequest получение списка соединений кластера
//
//  type GET_CONNECTIONS_SHORT_REQUEST = 51
//  kind MESSAGE_KIND = 1
//  respond GetConnectionsShortResponse
type GetConnectionsShortRequest struct {
	ClusterID uuid.UUID
}

func (r *GetConnectionsShortRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetConnectionsShortRequest) Type() EndpointMessageType {
	return GET_CONNECTIONS_SHORT_REQUEST
}

func (r *GetConnectionsShortRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

// GetConnectionsShortResponse ответ со списком соединений кластера
//
//  type GET_CONNECTIONS_SHORT_RESPONSE = 52
//  kind MESSAGE_KIND = 1
//  respond serialize.ConnectionShortInfoList
type GetConnectionsShortResponse struct {
	Connections serialize.ConnectionShortInfoList
}

func (_ *GetConnectionsShortResponse) Type() EndpointMessageType {
	return GET_CONNECTIONS_SHORT_RESPONSE
}

func (res *GetConnectionsShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.ConnectionShortInfoList{}
	list.Parse(decoder, version, r)

	res.Connections = list

}

var _ EndpointRequestMessage = (*DisconnectConnectionRequest)(nil)

// DisconnectConnectionRequest отключение соединения
//
//  type DISCONNECT_REQUEST = 59
//  respond nothing
type DisconnectConnectionRequest struct {
	ClusterID    uuid.UUID
	ProcessID    uuid.UUID
	InfobaseID   uuid.UUID
	ConnectionID uuid.UUID
}

func (r *DisconnectConnectionRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *DisconnectConnectionRequest) Type() EndpointMessageType {
	return DISCONNECT_REQUEST
}

func (r *DisconnectConnectionRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ProcessID, w)
	encoder.Uuid(r.ConnectionID, w)
}

var _ EndpointRequestMessage = (*GetInfobaseConnectionsShortRequest)(nil)

// GetInfobaseConnectionsShortRequest получение списка соединений кластера
//
//  type GET_INFOBASE_CONNECTIONS_SHORT_REQUEST = 52
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseConnectionsShortResponse
type GetInfobaseConnectionsShortRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
}

func (r *GetInfobaseConnectionsShortRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *GetInfobaseConnectionsShortRequest) Type() EndpointMessageType {
	return GET_INFOBASE_CONNECTIONS_SHORT_REQUEST
}

func (r *GetInfobaseConnectionsShortRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

// GetConnectionsShortResponse ответ со списком соединений кластера
//
//  type GET_INFOBASE_CONNECTIONS_SHORT_RESPONSE = 53
//  kind MESSAGE_KIND = 1
//  respond Connections serialize.ConnectionShortInfoList
type GetInfobaseConnectionsShortResponse struct {
	Connections serialize.ConnectionShortInfoList
}

func (_ *GetInfobaseConnectionsShortResponse) Type() EndpointMessageType {
	return GET_INFOBASE_CONNECTIONS_SHORT_RESPONSE
}

func (res *GetInfobaseConnectionsShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.ConnectionShortInfoList{}
	list.Parse(decoder, version, r)

	res.Connections = list

}
