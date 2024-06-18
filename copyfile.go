package copyfile

import "github.com/rrgmc/debefix"

func New(options ...Option) *CopyFile {
	ret := &CopyFile{}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func NewOptions(options ...Option) (*CopyFile, []debefix.LoadOption, []debefix.ResolveOption) {
	c := New(options...)
	return c,
		[]debefix.LoadOption{debefix.WithLoadValueParser(c)},
		[]debefix.ResolveOption{debefix.WithRowResolvedCallback(c)}
}

