package encoding

import "fmt"

// Size constants for better readability and configuration
const (
	// DefaultMaxSize is the default maximum input size (10MB)
	DefaultMaxSize = 1024 * 1024 * 10

	// MinimumMaxSize is the minimum allowed size (1KB)
	MinimumMaxSize = 1024

	// MaximumMaxSize is the absolute maximum allowed size (1GB)
	MaximumMaxSize = 1024 * 1024 * 1024
)

// Options holds all configuration options for the JSON parser
type Options struct {
	// MaxSize defines the maximum size of the input that the parser will accept
	MaxSize int

	// DisableSizeLimit allows bypassing size limit checks when set to true
	DisableSizeLimit bool

	// StrictMode enables additional validation during parsing
	StrictMode bool
}

// Validate checks if the options are valid
func (o *Options) Validate() error {
	if o.DisableSizeLimit {
		return nil
	}

	if o.MaxSize < MinimumMaxSize {
		return fmt.Errorf("max size %d is below minimum allowed size %d", o.MaxSize, MinimumMaxSize)
	}

	if o.MaxSize > MaximumMaxSize {
		return fmt.Errorf("max size %d exceeds maximum allowed size %d", o.MaxSize, MaximumMaxSize)
	}

	return nil
}

// Option defines a function type that can modify Options
type Option func(*Options) error

// defaultOptions returns the default options
func defaultOptions() *Options {
	return &Options{
		MaxSize:          DefaultMaxSize,
		DisableSizeLimit: false,
		StrictMode:       false,
	}
}

// WithMaxSize sets the maximum allowed input size
func WithMaxSize(size int) Option {
	return func(o *Options) error {
		if size <= 0 {
			return fmt.Errorf("max size must be positive, got %d", size)
		}

		o.MaxSize = size

		return nil
	}
}

// WithDisableSizeLimit disables size limit checking
func WithDisableSizeLimit() Option {
	return func(o *Options) error {
		o.DisableSizeLimit = true

		return nil
	}
}

// WithStrictMode enables strict parsing mode
func WithStrictMode() Option {
	return func(o *Options) error {
		o.StrictMode = true

		return nil
	}
}

// applyOptions applies the given options to the default options
func applyOptions(opts ...Option) (*Options, error) {
	options := defaultOptions()

	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, fmt.Errorf("invalid option: %w", err)
		}
	}

	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	return options, nil
}
