package websocket

import "time"

type ServerOption func(opt *serverOption)

type serverOption struct {
	Authentication

	ack        AckType
	ackTimeout time.Duration

	patten            string
	maxConnectionIdle time.Duration
}

func newOption(opts ...ServerOption) serverOption {
	o := serverOption{
		Authentication:    new(authentication),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		patten:            "/ws",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithAuthentication(authentication Authentication) ServerOption {
	return func(opt *serverOption) {
		opt.Authentication = authentication
	}
}

func WithServerAck(ack AckType) ServerOption {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithHandlerPatten(pattern string) ServerOption {
	return func(opt *serverOption) {
		opt.patten = pattern
	}
}

func WithServerMaxConnectionIdle(maxConnection time.Duration) ServerOption {
	return func(opt *serverOption) {
		if maxConnection > 0 {
			opt.maxConnectionIdle = maxConnection
		}
	}
}
