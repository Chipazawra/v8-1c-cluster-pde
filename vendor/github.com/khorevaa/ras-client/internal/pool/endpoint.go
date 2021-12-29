package pool

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize/esig"
	"io"
	"strings"
	"sync/atomic"
	"time"
)

func NewEndpoint(endpoint EndpointInfo) *Endpoint {

	return &Endpoint{
		id:        endpoint.ID(),
		version:   endpoint.Version(),
		format:    endpoint.Format(),
		serviceID: endpoint.ServiceID(),
		codec:     endpoint.Codec(),
	}
}

type EndpointInfo interface {
	ID() int
	Version() int
	Format() int16
	ServiceID() string
	Codec() codec.Codec
}

type Endpoint struct {
	id        int
	version   int
	format    int16
	serviceID string
	codec     codec.Codec

	conn      *Conn
	createdAt time.Time
	usedAt    uint32 // atomic
	pooled    bool
	Inited    bool

	sig          esig.ESIG
	clusterHash  string
	infobaseHash string
	onRequest    func(ctx context.Context, endpoint *Endpoint, req messages.EndpointRequestMessage) error
}

func calcHash(in string) string {

	str := base64.StdEncoding.EncodeToString([]byte(in))
	return str
}

func checkHash(val1, val2 string) bool {

	return strings.EqualFold(val1, val2)

}

func (e *Endpoint) Sig() esig.ESIG {
	return e.sig
}

func (e *Endpoint) SetSig(sig esig.ESIG) {
	e.sig = sig
}

func (e *Endpoint) UsedAt() time.Time {
	unix := atomic.LoadUint32(&e.usedAt)
	return time.Unix(int64(unix), 0)
}

func (e *Endpoint) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&e.usedAt, uint32(tm.Unix()))
}

func (e *Endpoint) ID() int {
	return e.id
}

func (e *Endpoint) Version() int {
	return e.version
}

func (e *Endpoint) Format() int16 {
	return e.format
}

func (e *Endpoint) ServiceID() string {
	return e.serviceID
}

func (e *Endpoint) Codec() codec.Codec {
	return e.codec
}

func (e *Endpoint) CheckClusterAuth(user, pwd string) bool {
	return checkHash(e.clusterHash, calcHash(fmt.Sprintf("%s:%s", user, pwd)))
}

func (e *Endpoint) SetClusterAuth(user, pwd string) {
	e.clusterHash = calcHash(fmt.Sprintf("%s:%s", user, pwd))
}

func (e *Endpoint) CheckInfobaseAuth(user, pwd string) bool {
	return checkHash(e.infobaseHash, calcHash(fmt.Sprintf("%s:%s", user, pwd)))
}

func (e *Endpoint) SetInfobaseAuth(user, pwd string) {
	e.infobaseHash = calcHash(fmt.Sprintf("%s:%s", user, pwd))
}

func (e *Endpoint) sendRequest(ctx context.Context, message *messages.EndpointMessage) (*messages.EndpointMessage, error) {

	e.SetUsedAt(time.Now())

	body := bytes.NewBuffer([]byte{})

	message.Format(e.codec.Encoder(), e.version, body)

	packet := NewPacket(messages.ENDPOINT_MESSAGE, body.Bytes())

	err := e.conn.SendPacket(packet)
	if err != nil {
		return nil, err
	}

	answer, err := e.conn.GetPacket(ctx)

	if err != nil {
		return nil, err
	}

	return e.tryParseMessage(answer)

}

func (e *Endpoint) sendVoidRequest(_ context.Context, conn *Conn, m messages.EndpointMessage) error {

	body := bytes.NewBuffer([]byte{})

	m.Format(e.codec.Encoder(), e.version, body)

	packet := NewPacket(byte(m.Type), body.Bytes())

	err := conn.SendPacket(packet)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) tryParseMessage(packet *Packet) (message *messages.EndpointMessage, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch val := e.(type) {

			case string:

				err = errors.New(val)

			case error:
				err = val
			default:
				panic(e)
			}
		}
	}()

	switch packet.Type {

	case messages.ENDPOINT_MESSAGE:

		decoder := e.codec.Decoder()

		endpointID := decoder.EndpointId(packet)
		format := decoder.Short(packet)

		message = &messages.EndpointMessage{
			EndpointID:     endpointID,
			EndpointFormat: format,
		}

		message.Parse(decoder, e.version, packet)

	case messages.ENDPOINT_FAILURE:

		decoder := e.codec.Decoder()

		err := &messages.EndpointFailure{}
		err.Parse(decoder, packet)

		return nil, err

	default:

		return nil, &messages.UnknownMessageError{
			Type:       packet.Type,
			Data:       packet.Data,
			EndpointID: e.id,
			ServiceID:  e.serviceID,
			Err:        ErrUnknownMessage}
	}

	return
}

func (e *Endpoint) tryFormatMessage(message *messages.EndpointMessage, writer io.Writer) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch val := e.(type) {

			case string:

				err = errors.New(val)

			case error:
				err = val
			default:
				panic(e)
			}
		}
	}()

	encoder := e.codec.Encoder()
	message.Format(encoder, e.version, writer)

	return
}

func (e *Endpoint) SendRequest(ctx context.Context, req messages.EndpointRequestMessage) (*messages.EndpointMessage, error) {

	if e.onRequest != nil {

		err := e.onRequest(ctx, e, req)

		if err != nil {
			return nil, err
		}

	}

	message := e.newEndpointMessage(req)
	answer, err := e.sendRequest(ctx, message)

	if err != nil {
		return nil, err
	}

	switch err := answer.Message.(type) {

	case *messages.EndpointMessageFailure:

		return nil, err

	case *messages.EndpointFailure:

		return nil, err

	}

	return answer, err

}

func (e *Endpoint) newEndpointMessage(req messages.EndpointRequestMessage) *messages.EndpointMessage {

	message := &messages.EndpointMessage{
		EndpointID:     e.id,
		EndpointFormat: e.format,
		Message:        req,
		Type:           req.Type(),
		Kind:           messages.MESSAGE_KIND,
	}

	return message

}
