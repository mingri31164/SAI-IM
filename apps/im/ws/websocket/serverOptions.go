package websocket

type ServerOption func(opt *serverOption)

type serverOption struct {
	Authentication
	patten string
}

func newOption(opts ...ServerOption) serverOption {
	o := serverOption{
		Authentication: new(authentication),
		patten:         "/ws",
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

func WithHandlerPatten(pattern string) ServerOption {
	return func(opt *serverOption) {
		opt.patten = pattern
	}
}
