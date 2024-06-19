package copyfile

import "github.com/rrgmc/debefix"

// New creates an instance of the CopyFile debefix plugin.
func New(options ...Option) *CopyFile {
	ret := &CopyFile{}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

// NewOptions returns debefix [debefix.LoadOption] and [debefix.ResolveOption] to use the plugin.
func NewOptions(options ...Option) (*CopyFile, []debefix.LoadOption, []debefix.ResolveOption) {
	c := New(options...)
	return c,
		[]debefix.LoadOption{debefix.WithLoadValueParser(c)},
		[]debefix.ResolveOption{debefix.WithRowResolvedCallback(c)}
}
