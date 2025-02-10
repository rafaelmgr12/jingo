package encoding

// DefaultMaxSize is the default maximum input size (10MB)
const DefaultMaxSize = 1024 * 1024 * 10 // 10MB

// Options define a function type that can modify the behavior of the encoding.
type Option func(*Options)

// Options holds all configuration options for the JSON parser
type Options struct {
	// MaxSize defines the maximum size of the input that the parser will accept.
	MaxSize int
}

// defaultOptions returns the default options
func defaultOptions() *Options {
	return &Options{
		MaxSize: DefaultMaxSize,
	}
}

// WithMaxSize sets the maximum allowed input size
func WithMaxSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.MaxSize = size
		}
	}
}

// applyOptions applies the given options to the default options
func applyOptions(opts ...Option) *Options {
	options := defaultOptions()

	for _, opt := range opts {
		opt(options)
	}

	return options
}
