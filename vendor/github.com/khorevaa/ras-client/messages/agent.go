package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	"io"
)

// GetAgentVersionRequest получение версии агента
//
//  type GET_AGENT_VERSION_REQUEST
//  respond GetAgentAdminsResponse
type GetAgentVersionRequest struct{}

func (r *GetAgentVersionRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *GetAgentVersionRequest) Format(_ codec.Encoder, _ int, _ io.Writer) {}

func (_ *GetAgentVersionRequest) Type() EndpointMessageType {
	return GET_AGENT_VERSION_REQUEST
}

// GetAgentVersionResponse ответ с версией агента кластера
//
//  type GET_AGENT_VERSION_RESPONSE
//  Users serialize.UsersList
type GetAgentVersionResponse struct {
	Version string
}

func (res *GetAgentVersionResponse) Parse(decoder codec.Decoder, _ int, r io.Reader) {

	decoder.StringPtr(&res.Version, r)
}

func (_ *GetAgentVersionResponse) Type() EndpointMessageType {
	return GET_AGENT_VERSION_RESPONSE
}

// GetAgentAdminsRequest получение списка админов агента
//
//  type GET_AGENT_ADMINS_REQUEST
//  respond GetAgentAdminsResponse
type GetAgentAdminsRequest struct{}

func (r *GetAgentAdminsRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *GetAgentAdminsRequest) Format(_ codec.Encoder, _ int, _ io.Writer) {}

func (_ *GetAgentAdminsRequest) Type() EndpointMessageType {
	return GET_AGENT_ADMINS_REQUEST
}

// GetAgentAdminsResponse ответ со списком админов агента кластера
//
//  type REG_CLUSTER_RESPONSE
//  Users serialize.UsersList
type GetAgentAdminsResponse struct {
	Users serialize.UsersList
}

func (res *GetAgentAdminsResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.UsersList{}
	list.Parse(decoder, version, r)

	res.Users = list

}

func (_ *GetAgentAdminsResponse) Type() EndpointMessageType {
	return GET_AGENT_ADMINS_RESPONSE
}

// RegAgentAdminRequest регистрация админа агента
//
//  type REG_AGENT_ADMIN_REQUEST
type RegAgentAdminRequest struct {
	User serialize.UserInfo
}

func (r *RegAgentAdminRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *RegAgentAdminRequest) Format(e codec.Encoder, v int, w io.Writer) {

	r.User.Format(e, v, w)

}

func (_ *RegAgentAdminRequest) Type() EndpointMessageType {
	return REG_AGENT_ADMIN_REQUEST
}

// UnregAgentAdminRequest удаление админа агента
//
//  type REG_AGENT_ADMIN_REQUEST
type UnregAgentAdminRequest struct {
	User string
}

func (r *UnregAgentAdminRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *UnregAgentAdminRequest) Format(e codec.Encoder, v int, w io.Writer) {

	e.String(r.User, w)

}

func (_ *UnregAgentAdminRequest) Type() EndpointMessageType {
	return UNREG_AGENT_ADMIN_REQUEST
}
