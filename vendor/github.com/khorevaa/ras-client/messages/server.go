package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetWorkingServersRequest)(nil)

// GetWorkingServersRequest получение списка рабочих серверов кластера
//
//  type GET_WORKING_SERVERS_REQUEST = 38
type GetWorkingServersRequest struct {
	ClusterID uuid.UUID
}

func (r *GetWorkingServersRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *GetWorkingServersRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ *GetWorkingServersRequest) Type() EndpointMessageType {
	return GET_WORKING_SERVERS_REQUEST
}

// GetWorkingServersResponse содержит список рабочих серверов кластера
//  type GET_WORKING_SERVERS_RESPONSE = 37
//  Servers []*serialize.ServerInfo
type GetWorkingServersResponse struct {
	Servers []*serialize.ServerInfo
}

func (res *GetWorkingServersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ServerInfo{}
		info.Parse(decoder, version, r)
		res.Servers = append(res.Servers, info)
	}

}

func (_ *GetWorkingServersResponse) Type() EndpointMessageType {
	return GET_WORKING_SERVERS_RESPONSE
}

var _ EndpointRequestMessage = (*GetWorkingServerInfoRequest)(nil)

// GetWorkingServerInfoRequest получение информации о рабочем сервере кластера
//
//  type GET_WORKING_SERVER_INFO_REQUEST = 38
type GetWorkingServerInfoRequest struct {
	ClusterID uuid.UUID
	ServerID  uuid.UUID
}

func (r *GetWorkingServerInfoRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *GetWorkingServerInfoRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ServerID, w)
}

func (_ *GetWorkingServerInfoRequest) Type() EndpointMessageType {
	return GET_WORKING_SERVER_INFO_REQUEST
}

// GetWorkingServerInfoResponse содержит информацию о рабочем сервере кластера
//  type GET_WORKING_SERVER_INFO_RESPONSE = 37
//  Info *serialize.ServerInfo
type GetWorkingServerInfoResponse struct {
	Info *serialize.ServerInfo
}

func (res *GetWorkingServerInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := &serialize.ServerInfo{}
	info.Parse(decoder, version, r)
	res.Info = info

}

func (_ *GetWorkingServerInfoResponse) Type() EndpointMessageType {
	return GET_WORKING_SERVER_INFO_RESPONSE
}

var _ EndpointRequestMessage = (*RegWorkingServerRequest)(nil)

// GetWorkingServerInfoRequest регистрация информации о рабочем сервере на кластере
//
//  type REG_WORKING_SERVER_REQUEST = 38
type RegWorkingServerRequest struct {
	ClusterID uuid.UUID
	Info      *serialize.ServerInfo
}

func (r *RegWorkingServerRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *RegWorkingServerRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	r.Info.Format(encoder, version, w)
}

func (_ *RegWorkingServerRequest) Type() EndpointMessageType {
	return REG_WORKING_SERVER_REQUEST
}

// GetWorkingServerInfoResponse содержит информацию о рабочем сервере кластера
//  type GET_WORKING_SERVER_INFO_RESPONSE = 37
//  Info *serialize.ServerInfo
type RegWorkingServerResponse struct {
	ServerID uuid.UUID
}

func (res *RegWorkingServerResponse) Parse(decoder codec.Decoder, _ int, r io.Reader) {
	decoder.UuidPtr(&res.ServerID, r)
}

func (_ *RegWorkingServerResponse) Type() EndpointMessageType {
	return REG_WORKING_SERVER_RESPONSE
}

var _ EndpointRequestMessage = (*UnregWorkingServerRequest)(nil)

// GetWorkingServerInfoRequest отмена регистрации информации о рабочем сервере на кластере
//
//  type UNREG_WORKING_SERVER_REQUEST = 38
type UnregWorkingServerRequest struct {
	ClusterID uuid.UUID
	ServerID  uuid.UUID
}

func (r *UnregWorkingServerRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *UnregWorkingServerRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ServerID, w)
}

func (_ *UnregWorkingServerRequest) Type() EndpointMessageType {
	return UNREG_WORKING_SERVER_REQUEST
}
