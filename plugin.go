package copyfile

import (
	"errors"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/rrgmc/debefix"
)

type CopyFile struct {
	debefix.ValueImpl
	callback         Callback
	setValueCallback SetValueCallback
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
		if c.callback != nil {
			err := c.callback(ctx, fieldname, file)
			if err != nil {
				return err
			}
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
	if !c.fileData.SetValue {
		// don't add a data field
		return nil, false, nil
	}
	if c.cf.setValueCallback == nil {
		return nil, false, errors.New("setValueCallback not set")
	}
	return c.cf.setValueCallback(ctx, c.fileData)
}

const (
	metadataName = "__copyfile__"
)

func getMetadata(metadata map[string]any) *FileDataList {
	if md, ok := metadata[metadataName]; ok {
		if mdfl, ok := md.(*FileDataList); ok {
			return mdfl
		}
	}
	return &FileDataList{
		Fields: map[string]FileData{},
	}
}

func setMetadata(ctx debefix.ValueCallbackResolveContext, fileData FileData) {
	md := getMetadata(ctx.Metadata())
	md.Fields[ctx.FieldName()] = fileData
	ctx.SetMetadata(metadataName, md)
}
