package rclient

import (
	"github.com/khorevaa/ras-client/internal/pool"
	"github.com/khorevaa/ras-client/protocol/codec"
)

type Options struct {
	serviceVersion string
	poolOptions    *pool.Options
	codec          codec.Codec
}

type Option func(opt *Options)

func WithVersion(version string) Option {
	return func(opt *Options) {
		if len(version) == 0 {
			return
		}
		opt.serviceVersion = version
	}
}

func WithPool(o *pool.Options) Option {
	return func(opt *Options) {
		if o == nil {
			return
		}
		opt.poolOptions = o
	}
}
