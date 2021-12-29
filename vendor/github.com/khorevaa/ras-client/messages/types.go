package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize/esig"
	"io"
)

type RequestMessage interface {
	Type() byte
	Format(codec codec.Encoder, w io.Writer)
}

type ResponseMessage interface {
	Type() byte
	Parse(codec codec.Decoder, r io.Reader)
}

type EndpointRequestMessage interface {
	Type() EndpointMessageType
	Format(encoder codec.Encoder, version int, w io.Writer)
	Sig() esig.ESIG
}

type EndpointResponseMessage interface {
	Type() EndpointMessageType
	Parse(decoder codec.Decoder, version int, r io.Reader)
}
