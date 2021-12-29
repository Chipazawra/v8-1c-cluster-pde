package messages

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/khorevaa/ras-client/protocol/codec"
	"io"
)

const magic = 475223888

type ConnectMessageAck struct {
	data []byte
}

func (r *ConnectMessageAck) Type() byte {
	return CONNECT_ACK
}

func (r *ConnectMessageAck) Parse(codec.Decoder, io.Reader) {}

type ConnectMessage struct {
	Params map[string]interface{}
}

func (m *ConnectMessage) String() string {
	return ""
}

func (m *ConnectMessage) Type() byte {
	return CONNECT
}

func (m ConnectMessage) Format(c codec.Encoder, w io.Writer) {

	size := len(m.Params)
	if size == 0 {
		c.Null(w)
		return
	}

	c.NullableSize(size, w)

	for key, value := range m.Params {

		c.String(key, w)
		c.TypedValue(value, w)

	}

}

type NegotiateMessage struct {
	magic           int
	ProtocolVersion int16
	CodecVersion    int16
}

func (n NegotiateMessage) Type() byte {
	return NEGOTIATE
}

func (n NegotiateMessage) Format(c codec.Encoder, w io.Writer) {

	c.Int(n.magic, w)
	c.Short(n.ProtocolVersion, w)
	c.Short(n.CodecVersion, w)

}

func NewNegotiateMessage(protocol, codec int16) NegotiateMessage {
	return NegotiateMessage{
		magic:           magic,
		ProtocolVersion: protocol,
		CodecVersion:    codec,
	}
}

const endpointPrefix = "v8.service.Admin.Cluster"

type OpenEndpointMessage struct {
	Encoding string
	Version  string
	params   map[string]interface{}
}

func (m *OpenEndpointMessage) String() string {
	return pp.Sprintln(m)
}

func (m *OpenEndpointMessage) Type() byte {
	return ENDPOINT_OPEN
}

func (m *OpenEndpointMessage) Format(c codec.Encoder, w io.Writer) {

	c.String(endpointPrefix, w)
	c.String(m.Version, w)
	size := len(m.params)
	if size == 0 {
		c.Null(w)
		return
	}

	c.NullableSize(size, w)

	for key, value := range m.params {

		c.String(key, w)
		c.TypedValue(value, w)

	}

}

type OpenEndpointMessageAck struct {
	ServiceID  string
	Version    string
	EndpointID int

	params map[string]interface{}
}

func (m *OpenEndpointMessageAck) Parse(c codec.Decoder, r io.Reader) {

	c.StringPtr(&m.ServiceID, r)
	c.StringPtr(&m.Version, r)

	m.EndpointID = c.EndpointId(r)

	// TODO Params

}

func (m *OpenEndpointMessageAck) String() string {
	return pp.Sprintln(m)
}

func (m *OpenEndpointMessageAck) Type() byte {
	return ENDPOINT_OPEN_ACK
}

type EndpointFailure struct {
	ServiceID  string      `json:"service_id"`
	Version    string      `json:"version"`
	EndpointID int         `json:"endpoint_id,omitempty"`
	ClassCause string      `json:"class_cause,omitempty"`
	Message    string      `json:"message"`
	Trace      []string    `json:"trace,omitempty"`
	Cause      *CauseError `json:"cause,omitempty"`
}

type CloseEndpointMessage struct {
	EndpointID int
}

func (m *CloseEndpointMessage) Type() byte {
	return ENDPOINT_CLOSE
}

func (m *CloseEndpointMessage) Format(c codec.Encoder, w io.Writer) {

	c.EndpointId(m.EndpointID, w)

}

type CauseError struct {
	Service string      `json:"service"`
	Message string      `json:"message"`
	Cause   *CauseError `json:"cause,omitempty"`
}

func (e *CauseError) Error() string {

	if e.Cause != nil {
		return fmt.Sprintf("service-err: %s msg-err: %s %s", e.Service, e.Message, e.Cause.Error())
	}

	return fmt.Sprintf("service-err: %s msg-err: %s", e.Service, e.Message)

}

func (m *CauseError) Parse(c codec.Decoder, r io.Reader) {

	m.Service = c.String(r)
	m.Message = c.String(r)
	errSize := c.Size(r)

	if errSize > 0 {

		panic("TODO ")

	}

	m.Cause = tryParseCauseError(c, r)

	if m.Cause != nil && len(m.Cause.Message) == 0 && m.Cause.Cause == nil {
		m.Cause = nil
	}

}

func tryParseCauseError(c codec.Decoder, r io.Reader) (err *CauseError) {
	defer func() {
		if e := recover(); e != nil {
			err = nil
		}
	}()

	err = &CauseError{}
	err.Parse(c, r)
	return
}

func (m *EndpointFailure) Parse(c codec.Decoder, r io.Reader) {

	c.StringPtr(&m.ServiceID, r)
	c.StringPtr(&m.Version, r)

	m.EndpointID = c.EndpointId(r)
	m.ClassCause = c.String(r)
	m.Message = c.String(r)
	errSize := c.Size(r)

	if errSize > 0 {

		panic("TODO ")

	}

	m.Cause = tryParseCauseError(c, r)
}

func (m *EndpointFailure) String() string {
	return m.Cause.Error()
}

func (m *EndpointFailure) Type() byte {
	return ENDPOINT_FAILURE
}

func (m *EndpointFailure) Error() string {

	return fmt.Sprintf("service-id: %s class:%s message: %s %s",
		m.ServiceID, m.ClassCause, m.Message, m.Cause.Error())
}
