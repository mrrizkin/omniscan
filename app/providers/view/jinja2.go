package view

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

type jinja2 struct {
	templates map[string]*exec.Template
}

func newJinja2(fs http.FileSystem, directory, extension string, env *exec.Environment) (ViewProvider, error) {
	templates := make(map[string]*exec.Template)

	loader, err := newHttpFileSystemLoader(fs, directory)
	if err != nil {
		return nil, err
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil || info.IsDir() {
			return nil
		}

		if len(extension) >= len(path) || path[len(path)-len(extension):] != extension {
			return nil
		}

		rel, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		name := filepath.ToSlash(rel)
		name = strings.TrimSuffix(name, extension)

		file, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		buf, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		tmpl, err := fromBytes(buf, loader, env)
		if err != nil {
			return err
		}

		templates[name] = tmpl

		return err
	}

	info, err := stat(fs, directory)
	if err != nil {
		err := walkFn(directory, nil, err)
		if err == nil {
			return nil, err
		}
	}
	err = walkInternal(fs, directory, info, walkFn)
	if err != nil {
		return nil, err
	}

	return &jinja2{
		templates: templates,
	}, nil
}

func (j *jinja2) Render(w io.Writer, template string, data ...map[string]interface{}) error {
	ctx := map[string]interface{}{}

	if len(data) > 0 {
		ctx = data[0]
	}

	tmpl, ok := j.templates[template]
	if !ok {
		return fmt.Errorf("template %s not found", template)
	}

	return tmpl.Execute(w, exec.NewContext(ctx))
}

type httpFilesystemLoader struct {
	fs      http.FileSystem
	baseDir string
}

func newHttpFileSystemLoader(
	httpfs http.FileSystem,
	baseDir string,
) (loaders.Loader, error) {
	hfs := &httpFilesystemLoader{
		fs:      httpfs,
		baseDir: baseDir,
	}
	if httpfs == nil {
		err := errors.New("httpfs cannot be nil")
		return nil, err
	}
	return hfs, nil
}

func (h *httpFilesystemLoader) Resolve(name string) (string, error) {
	return name, nil
}

// Get returns an io.Reader where the template's content can be read from.
func (h *httpFilesystemLoader) Read(path string) (io.Reader, error) {
	fullPath := path
	if h.baseDir != "" {
		fullPath = fmt.Sprintf(
			"%s/%s",
			h.baseDir,
			fullPath,
		)
	}

	return h.fs.Open(fullPath)
}

func (h *httpFilesystemLoader) Inherit(from string) (loaders.Loader, error) {
	hfs := &httpFilesystemLoader{
		fs:      h.fs,
		baseDir: h.baseDir,
	}

	return hfs, nil
}
