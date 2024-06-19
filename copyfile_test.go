package copyfile

import (
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/rrgmc/debefix"
	"gotest.tools/v3/assert"
)

func TestCopyFile(t *testing.T) {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tenants:
    rows:
      - tenant_id: 987
        name: "Joomla"
  tags:
    config:
      depends: ["tenants"]
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
        tagfilename:
          !copyfile
          value: "{value:tag_id}.png"
          source: "images/tags/javascript.png"
          destination: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
`),
		},
	})

	var source, destination string

	_, loadOptions, resolveOptions := NewOptions(
		WithSourcePath("/tmp/source"),
		WithDestinationPath("/tmp/destination"),
		WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error) {
			return DefaultGetPathsCallback(ctx, fieldname, fileData)
		}),
		WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error) {
			return DefaultGetValueCallback(ctx, fileData)
		}),
		WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
			assert.Equal(t, "", source)
			assert.Equal(t, "", destination)

			source = filepath.Join(sourcePath, sourceFilename)
			destination = filepath.Join(destinationPath, destinationFilename)
			return nil
		}),
	)

	data, err := debefix.Load(provider, loadOptions...)
	assert.NilError(t, err)

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	assert.NilError(t, err)

	assert.Equal(t, "559.png", resolvedData.Tables["tags"].Rows[0].Fields["tagfilename"])
	assert.Equal(t, "/tmp/source/images/tags/javascript.png", source)
	assert.Equal(t, "/tmp/destination/tenant/Joomla/images/tags/559.png", destination)
}

func TestCopyFileDefaultValue(t *testing.T) {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tenants:
    rows:
      - tenant_id: 987
        name: "Joomla"
  tags:
    config:
      depends: ["tenants"]
      default_values:
        tagfilename:
          !copyfile
          value: "{value:tag_id}.png"
          source: "images/tags/javascript.png"
          destination: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
`),
		},
	})

	var source, destination string

	_, loadOptions, resolveOptions := NewOptions(
		WithSourcePath("/tmp/source"),
		WithDestinationPath("/tmp/destination"),
		WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error) {
			return DefaultGetPathsCallback(ctx, fieldname, fileData)
		}),
		WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error) {
			return DefaultGetValueCallback(ctx, fileData)
		}),
		WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
			assert.Equal(t, "", source)
			assert.Equal(t, "", destination)

			source = filepath.Join(sourcePath, sourceFilename)
			destination = filepath.Join(destinationPath, destinationFilename)
			return nil
		}),
	)

	data, err := debefix.Load(provider, loadOptions...)
	assert.NilError(t, err)

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	assert.NilError(t, err)

	assert.Equal(t, "559.png", resolvedData.Tables["tags"].Rows[0].Fields["tagfilename"])
	assert.Equal(t, "/tmp/source/images/tags/javascript.png", source)
	assert.Equal(t, "/tmp/destination/tenant/Joomla/images/tags/559.png", destination)
}

func TestCopyFileNoValue(t *testing.T) {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tenants:
    rows:
      - tenant_id: 987
        name: "Joomla"
  tags:
    config:
      depends: ["tenants"]
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
        tagfilename:
          !copyfile
          source: "images/tags/javascript.png"
          destination: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
`),
		},
	})

	var source, destination string

	_, loadOptions, resolveOptions := NewOptions(
		WithSourcePath("/tmp/source"),
		WithDestinationPath("/tmp/destination"),
		WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error) {
			return DefaultGetPathsCallback(ctx, fieldname, fileData)
		}),
		WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error) {
			return DefaultGetValueCallback(ctx, fileData)
		}),
		WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
			assert.Equal(t, "", source)
			assert.Equal(t, "", destination)
			source = filepath.Join(sourcePath, sourceFilename)
			destination = filepath.Join(destinationPath, destinationFilename)
			return nil
		}),
	)

	data, err := debefix.Load(provider, loadOptions...)
	assert.NilError(t, err)

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	assert.NilError(t, err)

	_, isField := resolvedData.Tables["tags"].Rows[0].Fields["tagfilename"]
	assert.Assert(t, !isField)
	assert.Equal(t, "/tmp/source/images/tags/javascript.png", source)
	assert.Equal(t, "/tmp/destination/tenant/Joomla/images/tags/559.png", destination)
}

func TestCopyFileMetadata(t *testing.T) {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tenants:
    rows:
      - tenant_id: 987
        name: "Joomla"
  tags:
    config:
      depends: ["tenants"]
      default_values:
        tagfilename:
          !copyfile
          value: "{value:tag_id}.png"
          source: "images/tags/{metadata:sourceTag}.png"
          destination: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
        _metadata:
          !metadata
          sourceTag: "javascript_lang"
`),
		},
	})

	var source, destination string

	_, loadOptions, resolveOptions := NewOptions(
		WithSourcePath("/tmp/source"),
		WithDestinationPath("/tmp/destination"),
		WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error) {
			return DefaultGetPathsCallback(ctx, fieldname, fileData)
		}),
		WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error) {
			return DefaultGetValueCallback(ctx, fileData)
		}),
		WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
			assert.Equal(t, "", source)
			assert.Equal(t, "", destination)

			source = filepath.Join(sourcePath, sourceFilename)
			destination = filepath.Join(destinationPath, destinationFilename)
			return nil
		}),
	)

	data, err := debefix.Load(provider, loadOptions...)
	assert.NilError(t, err)

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	assert.NilError(t, err)

	assert.Equal(t, "559.png", resolvedData.Tables["tags"].Rows[0].Fields["tagfilename"])
	assert.Equal(t, "/tmp/source/images/tags/javascript_lang.png", source)
	assert.Equal(t, "/tmp/destination/tenant/Joomla/images/tags/559.png", destination)
}

func TestCopyFileMetadataDefault(t *testing.T) {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tenants:
    rows:
      - tenant_id: 987
        name: "Joomla"
  tags:
    config:
      depends: ["tenants"]
      default_values:
        tagfilename:
          !copyfile
          value: "{value:tag_id}.png"
          source: "images/tags/{metadata:sourceTag:default}.png"
          destination: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
`),
		},
	})

	var source, destination string

	_, loadOptions, resolveOptions := NewOptions(
		WithSourcePath("/tmp/source"),
		WithDestinationPath("/tmp/destination"),
		WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) (source string, destination string, err error) {
			return DefaultGetPathsCallback(ctx, fieldname, fileData)
		}),
		WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData FileData) (value any, addField bool, err error) {
			return DefaultGetValueCallback(ctx, fileData)
		}),
		WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
			assert.Equal(t, "", source)
			assert.Equal(t, "", destination)

			source = filepath.Join(sourcePath, sourceFilename)
			destination = filepath.Join(destinationPath, destinationFilename)
			return nil
		}),
	)

	data, err := debefix.Load(provider, loadOptions...)
	assert.NilError(t, err)

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	assert.NilError(t, err)

	assert.Equal(t, "559.png", resolvedData.Tables["tags"].Rows[0].Fields["tagfilename"])
	assert.Equal(t, "/tmp/source/images/tags/default.png", source)
	assert.Equal(t, "/tmp/destination/tenant/Joomla/images/tags/559.png", destination)
}
