package messages

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/khorevaa/ras-client/protocol/codec"
	"io"
	"strings"
)

//goland:noinspection ALL
const (
	VOID_MESSAGE_KIND EndpointMessageKind = 0
	MESSAGE_KIND      EndpointMessageKind = 1
	EXCEPTION_KIND    EndpointMessageKind = 0xff
)

type EndpointMessageKind int

func (e EndpointMessageKind) Type() byte {
	return byte(e)
}

type EndpointMessageFailure struct {
	ServiceID  string `json:"service_id"`
	ErrorType  string `json:"type"`
	EndpointID int    `json:"endpoint_id,omitempty"`
	Message    string `json:"message"`
}

func (m *EndpointMessageFailure) Parse(c codec.Decoder, r io.Reader) {

	c.StringPtr(&m.ServiceID, r)
	c.StringPtr(&m.Message, r)

	msg := strings.Split(m.ServiceID, "#")

	if len(msg) == 2 {
		m.ServiceID = msg[0]
		m.ErrorType = msg[1]
	}

}

func (m *EndpointMessageFailure) String() string {
	return pp.Sprintln(m)
}

func (m *EndpointMessageFailure) Type() EndpointMessageKind {
	return EXCEPTION_KIND
}

func (m *EndpointMessageFailure) Error() string {

	return fmt.Sprintf("endpoint: %d service: %s msg: %s", m.EndpointID, m.ServiceID, m.Message)
}

type EndpointMessage struct {
	EndpointID     int
	EndpointFormat int16
	Message        interface{}
	Type           EndpointMessageType
	Kind           EndpointMessageKind
}

func (m *EndpointMessage) Parse(decoder codec.Decoder, version int, reader io.Reader) {

	m.Kind = EndpointMessageKind(decoder.Byte(reader))

	switch m.Kind {

	case VOID_MESSAGE_KIND:
		return
	case EXCEPTION_KIND:

		fail := &EndpointMessageFailure{EndpointID: m.EndpointID}
		fail.Parse(decoder, reader)
		m.Message = fail

	case MESSAGE_KIND:

		respondType := decoder.Byte(reader)

		m.Type = EndpointMessageType(respondType)

		respond := m.Type.Parser()

		parser := respond.(codec.BinaryParser)

		// TODO Сделать получение ответа по типу
		parser.Parse(decoder, version, reader)

		m.Message = parser
	default:
		panic("unknown message kind")
	}

}

func (m *EndpointMessage) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.EndpointId(m.EndpointID, w)
	encoder.Short(m.EndpointFormat, w)

	encoder.Byte(byte(m.Kind), w)
	encoder.Byte(byte(m.Type), w) // МАГИЯ без этого байта требует авторизации на центральном кластере

	formatter := m.Message.(codec.BinaryWriter)
	formatter.Format(encoder, version, w) // запись тебя сообщения

}
