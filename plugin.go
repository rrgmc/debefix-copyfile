package copyfile

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/rrgmc/debefix"
)

// CopyFile is a debefix plugin to configure files to be copied during the data generation process.
type CopyFile struct {
	debefix.ValueImpl
	sourcePath       string
	destinationPath  string
	getPathsCallback GetPathsCallback
	getValueCallback GetValueCallback
	copyFileCallback CopyFileCallback
}

var (
	_ debefix.ValueParser         = (*CopyFile)(nil)
	_ debefix.RowResolvedCallback = (*CopyFile)(nil)
)

func (c *CopyFile) ParseValue(tag *ast.TagNode) (bool, any, error) {
	// parse "!copyfile" tag
	if tag.Start.Value != "!copyfile" {
		return false, nil, nil
	}

	var fileData FileData
	err := yaml.NodeToValue(tag.Value, &fileData, yaml.Strict())
	if err != nil {
		return false, nil, err
	}

	// return a [debefix.Value] to be processed later.
	return true, &copyFileValue{cf: c, fileData: fileData}, nil
}

func (c *CopyFile) RowResolved(ctx debefix.ValueResolveContext) error {
	// after row was resolved, call the callback to copy the file
	md := getMetadata(ctx.Row().Metadata)
	for fieldname, file := range md.Fields {
		getPathsCallback := c.getPathsCallback
		if getPathsCallback == nil {
			getPathsCallback = DefaultGetPathsCallback
		}
		source, destination, err := getPathsCallback(ctx, fieldname, file)
		if err != nil {
			return err
		}

		copyFileCallback := c.copyFileCallback
		if copyFileCallback == nil {
			copyFileCallback = DefaultCopyFileCallback
		}
		err = copyFileCallback(c.sourcePath, source, c.destinationPath, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

type copyFileValue struct {
	debefix.ValueImpl
	cf       *CopyFile
	fileData FileData
}

var (
	_ debefix.ValueCallback = (*copyFileValue)(nil)
)

func (c *copyFileValue) GetValueCallback(ctx debefix.ValueCallbackResolveContext) (resolvedValue any, addField bool, err error) {
	// copy the metadata to the row being processed, so it is available to [CopyFile.RowResolved].
	setMetadata(ctx, c.fileData)

	getValueCallback := c.cf.getValueCallback
	if getValueCallback == nil {
		getValueCallback = DefaultGetValueCallback
	}
	return getValueCallback(ctx, c.fileData)
}

const (
	metadataName = "__copyfile__"
)

type fileDataList struct {
	Fields map[string]FileData `json:"fields"`
}

func getMetadata(metadata map[string]any) *fileDataList {
	if md, ok := metadata[metadataName]; ok {
		if mdfl, ok := md.(*fileDataList); ok {
			return mdfl
		}
	}
	return &fileDataList{
		Fields: map[string]FileData{},
	}
}

func setMetadata(ctx debefix.ValueCallbackResolveContext, fileData FileData) {
	md := getMetadata(ctx.Metadata())
	md.Fields[ctx.FieldName()] = fileData
	ctx.SetMetadata(metadataName, md)
}
