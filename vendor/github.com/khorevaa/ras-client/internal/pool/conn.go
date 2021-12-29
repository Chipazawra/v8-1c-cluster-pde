package pool

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// IOError is the data type for errors occurring in case of failure.
type IOError struct {
	Connection net.Conn
	Error      error
}

type Conn struct {
	connMU *sync.Mutex

	_locked uint32
	_closed uint32
	netConn net.Conn
	onError func(err IOError)

	endpoints []*Endpoint
	closer    func(ctx context.Context, conn *Conn, endpoint *Endpoint) error

	createdAt time.Time
	usedAt    uint32 // atomic
	pooled    bool
	Inited    bool
}

func NewConn(netConn net.Conn) *Conn {

	cn := &Conn{
		createdAt: time.Now(),
		connMU:    &sync.Mutex{},
	}
	cn.SetNetConn(netConn)
	cn.SetUsedAt(time.Now())

	return cn
}

func (c *Conn) SendPacket(packet *Packet) error {

	c.SetUsedAt(time.Now())
	err := packet.Write(c.netConn)
	return err
}

func (c *Conn) GetPacket(ctx context.Context) (packet *Packet, err error) {

	c.SetUsedAt(time.Now())
	return c.readContext(ctx)
}

func (c *Conn) UsedAt() time.Time {
	unix := atomic.LoadUint32(&c.usedAt)
	return time.Unix(int64(unix), 0)
}

func (c *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&c.usedAt, uint32(tm.Unix()))
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.netConn.RemoteAddr()
}

func (c *Conn) SetNetConn(netConn net.Conn) {
	c.netConn = netConn
}

func (c *Conn) closed() bool {

	if atomic.LoadUint32(&c._closed) == 1 {
		return true
	}
	_ = c.netConn.SetReadDeadline(time.Now())
	_, err := c.netConn.Read(make([]byte, 0))
	var zero time.Time
	_ = c.netConn.SetReadDeadline(zero)

	if err == nil {
		return false
	}

	netErr, _ := err.(net.Error)
	if err != io.EOF && !netErr.Timeout() {
		atomic.StoreUint32(&c._closed, 1)
		return true
	}
	return false
}

func (c *Conn) Close() error {

	if !atomic.CompareAndSwapUint32(&c._closed, 0, 1) {
		return nil
	}

	if c.closer != nil {

		for _, endpoint := range c.endpoints {
			_ = c.closer(context.Background(), c, endpoint)
		}
	}

	return c.netConn.Close()
}

//func (conn *Conn) lock() {
//
//	conn.connMU.Lock()
//	atomic.StoreUint32(&conn._locked, 1)
//}
//
//func (conn *Conn) unlock() {
//
//	atomic.StoreUint32(&conn._locked, 0)
//	conn.connMU.Unlock()
//
//}

func (c *Conn) Locked() bool {
	return atomic.LoadUint32(&c._locked) == 1
}

func (c *Conn) readContext(ctx context.Context) (*Packet, error) {

	recvDone := make(chan *Packet)
	errChan := make(chan error)

	go c.readPacket(recvDone, errChan)

	// setup the cancellation to abort reads in process
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
			// Close() can be used if this isn't necessarily a TCP connection
		case err := <-errChan:
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				go c.readPacket(recvDone, errChan)
				continue
			}
			return nil, err
		case packet := <-recvDone:
			return packet, nil
		}
	}

}

func (c *Conn) readPacket(recvDone chan *Packet, errChan chan error) {

	//c.lock()
	//defer c.unlock()

	err := c.netConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		errChan <- err
		return
	}

	typeBuffer := make([]byte, 1)

	_, err = c.netConn.Read(typeBuffer)

	if err != nil {

		if c.onError != nil {
			c.onError(IOError{c.netConn, err})
		}
		errChan <- err
		return
	}

	size, err := decodeSize(c.netConn)

	if err != nil {
		if c.onError != nil {
			c.onError(IOError{c.netConn, err})
		}
		errChan <- err
		return
	}

	data := make([]byte, size)
	readLength := 0
	n := 0

	for readLength < len(data) {
		n, err = c.netConn.Read(data[readLength:])
		readLength += n

		if err != nil {
			if c.onError != nil {
				c.onError(IOError{c.netConn, err})
			}
			errChan <- err
			return
		}
	}

	recvDone <- NewPacket(typeBuffer[0], data)

}

func decodeSize(r io.Reader) (int, error) {
	ff := 0xFFFFFF80
	b1, err := readByte(r)

	if err != nil {
		return 0, err
	}
	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := 7; (cur & ff) != 0x0; {

		b1, err = readByte(r)

		if err != nil {
			return 0, err
		}

		cur = int(b1 & 0xFF)
		size += (cur & 0x7F) << shift
		shift += 7

	}

	return size, nil
}

func readByte(r io.Reader) (byte, error) {

	byteBuffer := make([]byte, 1)
	_, err := r.Read(byteBuffer)

	return byteBuffer[0], err
}
