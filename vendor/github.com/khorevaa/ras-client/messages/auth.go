package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

// ClusterAuthenticateRequest установка авторизации на кластере
//
//  type AUTHENTICATE_REQUEST = 10
//  kind MESSAGE_KIND = 1
//  respond nothing
type ClusterAuthenticateRequest struct {
	ClusterID      uuid.UUID
	User, Password string
}

func (r ClusterAuthenticateRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r ClusterAuthenticateRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.String(r.User, w)
	encoder.String(r.Password, w)
}

func (_ ClusterAuthenticateRequest) Type() EndpointMessageType {
	return AUTHENTICATE_REQUEST
}

// AuthenticateAgentRequest установка авторизации на агенте
//
//  type AUTHENTICATE_AGENT_REQUEST = 9
//  kind MESSAGE_KIND = 1
//  respond nothing
type AuthenticateAgentRequest struct {
	User, Password string
}

func (_ AuthenticateAgentRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (_ AuthenticateAgentRequest) Type() EndpointMessageType {
	return AUTHENTICATE_AGENT_REQUEST
}

func (r AuthenticateAgentRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.String(r.User, w)
	encoder.String(r.Password, w)

}

// AuthenticateInfobaseRequest установка авторизации в информационной базе
//
//  type ADD_AUTHENTICATION_REQUEST = 11
//  kind MESSAGE_KIND = 1
//  respond nothing
type AuthenticateInfobaseRequest struct {
	ClusterID      uuid.UUID
	User, Password string
}

func (r AuthenticateInfobaseRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ AuthenticateInfobaseRequest) Type() EndpointMessageType {
	return ADD_AUTHENTICATION_REQUEST
}

func (r AuthenticateInfobaseRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	encoder.String(r.User, w)
	encoder.String(r.Password, w)

}
