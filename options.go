package copyfile

import "github.com/rrgmc/debefix"

type GetPathsCallback func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error)

type GetValueCallback func(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error)

type CopyFileCallback func(sourceFilename, destinationFilename string) error

type Option func(*CopyFile)

func WithSourcePath(sourcePath string) Option {
	return func(c *CopyFile) {
		c.sourcePath = sourcePath
	}
}

func WithDestinationPath(destinationPath string) Option {
	return func(c *CopyFile) {
		c.destinationPath = destinationPath
	}
}

func WithCopyFileCallback(callback CopyFileCallback) Option {
	return func(c *CopyFile) {
		c.copyFileCallback = callback
	}
}

func WithGetPathsCallback(callback GetPathsCallback) Option {
	return func(c *CopyFile) {
		c.getPathsCallback = callback
	}
}

func WithGetValueCallback(setValueCallback GetValueCallback) Option {
	return func(c *CopyFile) {
		c.getValueCallback = setValueCallback
	}
}
