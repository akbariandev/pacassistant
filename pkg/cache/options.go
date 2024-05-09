package cache

var defaultServerOptions = options{}

type options struct{}

type Option interface {
	apply(*options)
}

type EmptyServerOption struct{}

func (EmptyServerOption) apply(*options) {}

// joinOption provides a way to combine arbitrary number of
// options into one.
type joinOption struct {
	opts []Option
}

func (mdo *joinOption) apply(do *options) {
	for _, opt := range mdo.opts {
		opt.apply(do)
	}
}

func newJoinOption(opts ...Option) Option {
	return &joinOption{opts: opts}
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncServerOption(f func(options2 *options)) *funcOption {
	return &funcOption{
		f: f,
	}
}
