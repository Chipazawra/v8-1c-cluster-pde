package pool

import (
	"context"
	"net"
	"time"
)

type Options struct {
	Dialer         func(ctx context.Context) (net.Conn, error)
	OnClose        func(conn *Conn) error
	OpenEndpoint   func(ctx context.Context, conn *Conn) (EndpointInfo, error)
	CloseEndpoint  func(ctx context.Context, conn *Conn, endpoint *Endpoint) error
	InitConnection func(ctx context.Context, conn *Conn) error

	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	MaxOpenEndpoints   int
}
