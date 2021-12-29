package rclient

import (
	"bytes"
	"context"
	"github.com/k0kubun/pp"
	"github.com/khorevaa/ras-client/internal/pool"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize/esig"
	"net"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const protocolVersion = 256

var serviceVersions = []string{"3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0"}

var _ Api = (*Client)(nil)

type Client struct {
	addr  string
	laddr net.Addr

	ctx          context.Context
	stopRoutines context.CancelFunc // остановить ping, read, и подобные горутины

	agentUser     string
	agentPassword string

	base pool.EndpointPool

	codec codec.Codec

	serviceVersion string
}

func (c *Client) Version() string {
	return c.serviceVersion
}

func (c *Client) Close() error {
	return c.base.Close()
}

func NewClient(addr string, opts ...Option) *Client {

	opt := &Options{
		serviceVersion: "9.0",
		codec:          codec.NewCodec1_0(),
	}

	for _, fn := range opts {
		fn(opt)
	}

	m := new(Client)
	m.addr = addr
	m.codec = opt.codec
	if opt.poolOptions != nil {
		m.base = pool.NewEndpointPool(opt.poolOptions)
	} else {
		m.base = pool.NewEndpointPool(m.poolOptions())
	}

	m.serviceVersion = opt.serviceVersion

	return m
}

func (c *Client) poolOptions() *pool.Options {
	return &pool.Options{
		Dialer:             c.dialFn,
		OpenEndpoint:       c.openEndpoint,
		CloseEndpoint:      c.closeEndpoint,
		InitConnection:     c.initConnection,
		PoolSize:           5,
		MinIdleConns:       1,
		MaxConnAge:         30 * time.Minute,
		IdleTimeout:        10 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
		PoolTimeout:        10 * time.Minute,
	}
}

func (c *Client) initConnection(ctx context.Context, conn *pool.Conn) error {

	negotiateMessage := messages.NewNegotiateMessage(protocolVersion, c.codec.Version())

	err := c.sendRequestMessage(conn, negotiateMessage)

	if err != nil {
		return err
	}

	err = c.sendRequestMessage(conn, &messages.ConnectMessage{Params: map[string]interface{}{
		"connect.timeout": int64(2000),
	}})

	packet, err := conn.GetPacket(ctx)

	if err != nil {
		return err
	}

	answer, err := c.tryParseMessage(packet)

	if err != nil {
		return err
	}

	if _, ok := answer.(*messages.ConnectMessageAck); !ok {
		return errors.New("unknown ack")
	}

	return nil
}

func (c *Client) openEndpoint(ctx context.Context, conn *pool.Conn) (info pool.EndpointInfo, err error) {

	var ack *messages.OpenEndpointMessageAck

	ack, err = c.tryOpenEndpoint(ctx, conn)
	if err != nil {

		message, ok := err.(*messages.EndpointFailure)

		if !ok {
			return nil, err
		}
		supportedVersion := detectSupportedVersion(message)
		if len(supportedVersion) == 0 {
			return nil, err
		}

		c.serviceVersion = supportedVersion
		ack, err = c.tryOpenEndpoint(ctx, conn)
	}

	if err != nil {
		return nil, err
	}

	endpointVersion, err := strconv.ParseFloat(ack.Version, 10)
	if err != nil {
		return nil, err
	}

	return endpointInfo{
		id:        ack.EndpointID,
		version:   int(endpointVersion),
		format:    0, // defaultFormat,
		serviceID: ack.ServiceID,
		codec:     c.codec,
	}, nil
}

type endpointInfo struct {
	id        int
	version   int
	format    int16
	serviceID string
	codec     codec.Codec
}

func (e endpointInfo) ID() int {
	return e.id
}

func (e endpointInfo) Version() int {
	return e.version
}

func (e endpointInfo) Format() int16 {
	return e.format
}

func (e endpointInfo) ServiceID() string {
	return e.serviceID
}

func (e endpointInfo) Codec() codec.Codec {
	return e.codec
}

func (c *Client) tryOpenEndpoint(ctx context.Context, conn *pool.Conn) (*messages.OpenEndpointMessageAck, error) {

	err := c.sendRequestMessage(conn, &messages.OpenEndpointMessage{Version: c.serviceVersion})

	packet, err := conn.GetPacket(ctx)

	if err != nil {
		return nil, err
	}

	answer, err := c.tryParseMessage(packet)

	if err != nil {
		return nil, err
	}

	switch t := answer.(type) {

	case *messages.EndpointFailure:

		return nil, t

	case *messages.OpenEndpointMessageAck:

		return t, nil

	default:

		pp.Println(answer)
		panic("unknown answer type")
	}

}

func (c *Client) closeEndpoint(_ context.Context, conn *pool.Conn, endpoint *pool.Endpoint) error {

	err := c.sendRequestMessage(conn, &messages.CloseEndpointMessage{EndpointID: endpoint.ID()})

	if err != nil {
		return err
	}

	return nil
}
func (c *Client) sendRequestMessage(conn *pool.Conn, message messages.RequestMessage) error {

	body := bytes.NewBuffer([]byte{})
	message.Format(c.codec.Encoder(), body)
	packet := pool.NewPacket(message.Type(), body.Bytes())

	err := conn.SendPacket(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) tryParseMessage(packet *pool.Packet) (message messages.ResponseMessage, err error) {
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

	case messages.CONNECT_ACK:

		decoder := c.codec.Decoder()

		message = &messages.ConnectMessageAck{}
		message.Parse(decoder, packet)

	case messages.KEEP_ALIVE:
		// nothing
	case messages.ENDPOINT_OPEN_ACK:

		decoder := c.codec.Decoder()

		message = &messages.OpenEndpointMessageAck{}
		message.Parse(decoder, packet)

	case messages.ENDPOINT_FAILURE:

		decoder := c.codec.Decoder()

		message = &messages.EndpointFailure{}
		message.Parse(decoder, packet)

	case messages.NULL_TYPE:

		panic(pp.Sprintln(int(packet.Type), "packet", packet))

	default:

		panic(pp.Sprintln(int(packet.Type), "packet", packet))
	}

	return
}

func (c *Client) dialFn(ctx context.Context) (net.Conn, error) {

	_, err := net.ResolveTCPAddr("tcp", c.addr)
	if err != nil {
		return nil, errors.Wrap(err, "resolving tcp")
	}

	var dialer net.Dialer

	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return nil, errors.Wrap(err, "dialing tcp")
	}

	return conn, nil

}

func (c *Client) getEndpoint(ctx context.Context, sig esig.ESIG) (*pool.Endpoint, error) {

	return c.base.Get(ctx, sig)

}

func (c *Client) putEndpoint(ctx context.Context, endpoint *pool.Endpoint) {

	c.base.Put(ctx, endpoint)

}

func (c *Client) withEndpoint(ctx context.Context, sig esig.ESIG, fn func(context.Context, *pool.Endpoint) error) error {

	cn, err := c.getEndpoint(ctx, sig)
	if err != nil {
		return err
	}

	defer c.putEndpoint(ctx, cn)

	err = fn(ctx, cn)

	return err

}

func (c *Client) sendEndpointRequest(ctx context.Context, req messages.EndpointRequestMessage) (interface{}, error) {

	var value interface{}

	err := c.withEndpoint(ctx, req.Sig(), func(ctx context.Context, p *pool.Endpoint) error {

		message, err := p.SendRequest(ctx, req)

		if err != nil {
			return err
		}

		value = message.Message

		return err
	})

	return value, err

}

func (c *Client) Disconnect() error {
	// stop all routines
	c.stopRoutines()

	//err := c.conn.Close()
	//if err != nil {
	//	return errors.Wrap(err, "closing TCP connection")
	//}

	// TODO: закрыть каналы

	// возвращаем в false, потому что мы теряем конфигурацию
	// сессии, и можем ее потерять во время отключения.

	return nil
}
