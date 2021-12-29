package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*TerminateSessionRequest)(nil)

// TerminateSessionRequest отключение сеанса
//
//  type DISCONNECT_REQUEST = 71
//  kind MESSAGE_KIND = 1
//  respond nothing
type TerminateSessionRequest struct {
	ClusterID uuid.UUID
	SessionID uuid.UUID
	Message   string
}

func (r *TerminateSessionRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *TerminateSessionRequest) Type() EndpointMessageType {
	return TERMINATE_SESSION_REQUEST
}

func (r *TerminateSessionRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.SessionID, w)
	encoder.String(r.Message, w)
}

var _ EndpointRequestMessage = (*GetInfobaseSessionsRequest)(nil)

// GetInfobaseSessionsRequest получение списка сессий информационной базы кластера
//
//  type GET_INFOBASE_SESSIONS_REQUEST = 61
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseSessionsResponse
type GetInfobaseSessionsRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
}

func (r *GetInfobaseSessionsRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *GetInfobaseSessionsRequest) Type() EndpointMessageType {
	return GET_INFOBASE_SESSIONS_REQUEST
}

func (r *GetInfobaseSessionsRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

// GetInfobaseSessionsResponse ответ со списком сессий кластера
//
//  type GET_INFOBASE_SESSIONS_RESPONSE = 62
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetInfobaseSessionsResponse struct {
	Sessions serialize.SessionInfoList
}

func (_ *GetInfobaseSessionsResponse) Type() EndpointMessageType {
	return GET_INFOBASE_SESSIONS_RESPONSE
}

func (res *GetInfobaseSessionsResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.SessionInfoList{}
	list.Parse(decoder, version, r)

	res.Sessions = list

}

var _ EndpointRequestMessage = (*GetSessionsRequest)(nil)

// GetSessionsRequest получение списка сессий кластера
//
//  type GET_SESSIONS_REQUEST = 59
//  kind MESSAGE_KIND = 1
//  respond GetSessionsResponse
type GetSessionsRequest struct {
	ClusterID uuid.UUID
}

func (r *GetSessionsRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetSessionsRequest) Type() EndpointMessageType {
	return GET_SESSIONS_REQUEST
}

func (r *GetSessionsRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

// GetInfobaseSessionsResponse ответ со списком сессий кластера
//
//  type GET_SESSIONS_RESPONSE = 60
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetSessionsResponse struct {
	Sessions serialize.SessionInfoList
}

func (_ *GetSessionsResponse) Type() EndpointMessageType {
	return GET_SESSIONS_RESPONSE
}

func (res *GetSessionsResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.SessionInfoList{}
	list.Parse(decoder, version, r)

	res.Sessions = list

}
