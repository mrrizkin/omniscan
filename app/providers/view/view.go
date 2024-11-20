package view

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mrrizkin/omniscan/app/providers/cache"
	"github.com/mrrizkin/omniscan/config"
	"github.com/mrrizkin/omniscan/resources"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/exec"
)

type ViewProvider interface {
	Render(w io.Writer, template string, data ...map[string]interface{}) error
}

type View struct {
	provider  ViewProvider
	fs        http.FileSystem
	directory string
	extension string
	config    *config.App
	cache     *cache.Cache
	env       *exec.Environment
}

func (*View) Construct() interface{} {
	return func(
		cfg *config.App,
		cache *cache.Cache,
	) (*View, error) {
		fs := http.FS(resources.Views)
		directory := cfg.VIEW_DIRECTORY
		extension := cfg.VIEW_EXTENSION
		env := &exec.Environment{
			Context:           exec.EmptyContext().Update(builtins.GlobalFunctions).Update(builtins.GlobalVariables),
			Filters:           builtins.Filters,
			Tests:             builtins.Tests,
			ControlStructures: builtins.ControlStructures,
			Methods:           builtins.Methods,
		}

		return &View{
			config: cfg,
			cache:  cache,
			env:    env,

			fs:        fs,
			directory: directory,
			extension: extension,
		}, nil
	}
}

func (v *View) AddContext(ctx *exec.Context) {
	v.env.Context.Update(ctx)
}

func (v *View) AddFilter(filter *exec.FilterSet) {
	v.env.Filters.Update(filter)
}

func (v *View) AddTest(test *exec.TestSet) {
	v.env.Tests.Update(test)
}

func (v *View) AddControlStructure(controlStructure *exec.ControlStructureSet) {
	v.env.ControlStructures.Update(controlStructure)
}

func (v *View) Compile() error {
	provider, err := newJinja2(v.fs, v.directory, v.extension, v.env)
	if err != nil {
		return err
	}

	v.provider = provider

	return nil
}

func (v *View) Render(template string, data map[string]interface{}) ([]byte, error) {
	var cacheKey string
	if v.config.VIEW_CACHE && v.config.IsProduction() {
		var hash [16]byte
		if data != nil {
			encodedData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}

			hash = md5.Sum(append([]byte(template), encodedData...))
		} else {
			hash = md5.Sum([]byte(template))
		}

		cacheKey = hex.EncodeToString(hash[:])
		if value, ok := v.cache.Get(cacheKey); ok {
			return value.([]byte), nil
		}
	}

	var buf bytes.Buffer
	err := v.provider.Render(&buf, template, data)
	if err != nil {
		return nil, err
	}

	if v.config.VIEW_CACHE && v.config.IsProduction() {
		v.cache.Set(cacheKey, buf.Bytes())
	}

	return buf.Bytes(), nil
}
