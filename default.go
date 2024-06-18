package copyfile

import (
	"fmt"
	"io"
	"os"

	"github.com/rrgmc/debefix"
)

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

func DefaultCopyFileCallback(sourceFilename, destinationFilename string) error {
	sourceFileStat, err := os.Stat(sourceFilename)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFilename)
	}

	source, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationFilename)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}
