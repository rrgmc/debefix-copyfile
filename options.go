package copyfile

import "github.com/rrgmc/debefix"

// GetPathsCallback gets the source and destination file names from [FileData].
type GetPathsCallback func(ctx debefix.ValueResolveContext, fieldname string,
	fileData FileData) (source string, destination string, err error)

// GetValueCallback allows getting a value for the file field. If addField is false, no field will be added to the data,
// only the file will be copied.
type GetValueCallback func(ctx debefix.ValueCallbackResolveContext,
	fileData FileData) (value any, addField bool, err error)

// CopyFileCallback copy a file from a source to a destination.
type CopyFileCallback func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error

type Option func(*CopyFile)

// WithSourcePath sets the source path, the root of all source filenames.
func WithSourcePath(sourcePath string) Option {
	return func(c *CopyFile) {
		c.sourcePath = sourcePath
	}
}

// WithDestinationPath sets the destination path, the root of all destination filenames.
func WithDestinationPath(destinationPath string) Option {
	return func(c *CopyFile) {
		c.destinationPath = destinationPath
	}
}

// WithGetPathsCallback sets the callback to get the source and destination filenames from the FileData parameters.
// The default implementation DefaultGetPathsCallback replaces field filters in both source and destination using
// ReplaceFieldsWithFilter.
func WithGetPathsCallback(callback GetPathsCallback) Option {
	return func(c *CopyFile) {
		c.getPathsCallback = callback
	}
}

// WithGetValueCallback sets the callback to have the option to set the file field value.
// The default implementation DefaultGetValueCallback don't add the field if [FileData.Value] was not set,
// otherwise it replaces field filters in it using ReplaceFieldsWithFilter.
func WithGetValueCallback(setValueCallback GetValueCallback) Option {
	return func(c *CopyFile) {
		c.getValueCallback = setValueCallback
	}
}

// WithCopyFileCallback sets the callback that copy files from a source to a destination.
func WithCopyFileCallback(callback CopyFileCallback) Option {
	return func(c *CopyFile) {
		c.copyFileCallback = callback
	}
}
