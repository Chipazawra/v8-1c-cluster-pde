package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetLocksRequest)(nil)

// GetLocksRequest получение списка блокировок кластера
//
//  type GET_LOCKS_REQUEST = 66
//  kind MESSAGE_KIND = 1
//  respond GetSessionsResponse
type GetLocksRequest struct {
	ClusterID uuid.UUID
}

func (r *GetLocksRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetLocksRequest) Type() EndpointMessageType {
	return GET_LOCKS_REQUEST
}

func (r *GetLocksRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

// GetLocksResponse ответ со списком блокировок кластера
//
//  type GET_LOCKS_RESPONSE = 67
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetLocksResponse struct {
	List serialize.LocksList
}

func (_ *GetLocksResponse) Type() EndpointMessageType {
	return GET_LOCKS_RESPONSE
}

func (res *GetLocksResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)

	res.List = list

}

var _ EndpointRequestMessage = (*GetInfobaseLockRequest)(nil)

// GetInfobaseLockRequest получение списка блокировок информационной базы кластера
//
//  type GET_INFOBASE_LOCKS_REQUEST = 68
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseSessionsResponse
type GetInfobaseLockRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
}

func (r *GetInfobaseLockRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *GetInfobaseLockRequest) Type() EndpointMessageType {
	return GET_INFOBASE_LOCKS_REQUEST
}

func (r *GetInfobaseLockRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

// GetInfobaseLockResponse ответ со списком сблокировок иб
//
//  type GET_INFOBASE_LOCKS_RESPONSE = 69
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetInfobaseLockResponse struct {
	List serialize.LocksList
}

func (_ *GetInfobaseLockResponse) Type() EndpointMessageType {
	return GET_INFOBASE_LOCKS_RESPONSE
}

func (res *GetInfobaseLockResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)

	res.List = list

}

var _ EndpointRequestMessage = (*GetSessionLockRequest)(nil)

// GetSessionLockRequest получение списка блокировок сессии информационной базы кластера
//
//  type GET_SESSION_LOCKS_REQUEST = 72
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseSessionsResponse
type GetSessionLockRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	SessionID  uuid.UUID
}

func (r *GetSessionLockRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *GetSessionLockRequest) Type() EndpointMessageType {
	return GET_SESSION_LOCKS_REQUEST
}

func (r *GetSessionLockRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
	encoder.Uuid(r.SessionID, w)
}

// GetSessionLockResponse ответ со списком блокировок сессии иб
//
//  type GET_SESSION_LOCKS_RESPONSE = 73
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetSessionLockResponse struct {
	List serialize.LocksList
}

func (_ *GetSessionLockResponse) Type() EndpointMessageType {
	return GET_SESSION_LOCKS_RESPONSE
}

func (res *GetSessionLockResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)

	res.List = list

}

var _ EndpointRequestMessage = (*GetConnectionLockRequest)(nil)

// GetSessionLockRequest получение списка блокировок сессии информационной базы кластера
//
//  type GET_CONNECTION_LOCKS_REQUEST = 70
//  kind MESSAGE_KIND = 1
//  respond GetConnectionLockResponse
type GetConnectionLockRequest struct {
	ClusterID    uuid.UUID
	ConnectionID uuid.UUID
}

func (r *GetConnectionLockRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetConnectionLockRequest) Type() EndpointMessageType {
	return GET_CONNECTION_LOCKS_REQUEST
}

func (r *GetConnectionLockRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ConnectionID, w)
}

// GetSessionLockResponse ответ со списком блокировок сессии иб
//
//  type GET_CONNECTION_LOCKS_RESPONSE = 71
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetConnectionLockResponse struct {
	List serialize.LocksList
}

func (_ *GetConnectionLockResponse) Type() EndpointMessageType {
	return GET_CONNECTION_LOCKS_RESPONSE
}

func (res *GetConnectionLockResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)

	res.List = list

}
