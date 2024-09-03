package copyfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rrgmc/debefix"
)

// ReplaceFieldsWithFilter replaces curly-braces separated fields in str as debefix filter expressions.
func ReplaceFieldsWithFilter(str string, ctx debefix.ValueResolveContext) (string, error) {
	p := ParseFields(str)
	if len(p.Fields()) == 0 {
		return str, nil
	}

	rmap := map[string]string{}
	for _, fld := range p.Fields() {
		rmap[fld] = fld
	}

	replaceValues, err := ctx.ResolvedData().ExtractValues(ctx.Row(), rmap)
	if err != nil {
		return "", err
	}

	return p.Replace(replaceValues)
}

// DefaultGetPathsCallback is the default implementation of GetPathsCallback.
func DefaultGetPathsCallback(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error) {
	source, err = ReplaceFieldsWithFilter(fileData.Source, ctx)
	if err != nil {
		return "", "", err
	}
	destination, err = ReplaceFieldsWithFilter(fileData.Destination, ctx)
	if err != nil {
		return "", "", err
	}
	return source, destination, nil
}

// DefaultGetValueCallback is the default implementation of GetValueCallback.
func DefaultGetValueCallback(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error) {
	if fileData.Value == nil {
		return nil, false, nil
	}
	sv, err := ReplaceFieldsWithFilter(*fileData.Value, ctx)
	if err != nil {
		return nil, false, err
	}
	return sv, true, nil
}

// DefaultCopyFileCallback is the default implementation of CopyFileCallback.
func DefaultCopyFileCallback(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
	if sourcePath == "" || destinationPath == "" {
		return fmt.Errorf("source and destination paths are required")
	}
	if sourceFilename == "" || destinationFilename == "" {
		return fmt.Errorf("source and destination file names are required")
	}

	sourceFileStat, err := os.Stat(filepath.Join(sourcePath, sourceFilename))
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFilename)
	}

	source, err := os.Open(filepath.Join(sourcePath, sourceFilename))
	if err != nil {
		return err
	}
	defer source.Close()

	destinationFullFilename := filepath.Join(destinationPath, destinationFilename)

	err = os.MkdirAll(filepath.Dir(destinationFullFilename), os.ModePerm)
	if err != nil {
		return err
	}

	destination, err := os.Create(destinationFullFilename)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}
