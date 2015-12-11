// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package pongo2 is a middleware that provides pongo2 template engine for Macaron.
package pongo2

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/macaron"
	"gopkg.in/flosch/pongo2.v3"
)

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentBinary  = "application/octet-stream"
	ContentJSON    = "application/json"
	ContentHTML    = "text/html"
	ContentXHTML   = "application/xhtml+xml"
	ContentXML     = "text/xml"
	defaultCharset = "UTF-8"
)

const (
	_DEFAULT_TPL_SET_NAME = "DEFAULT"
)

func compile(opt Options) map[string]*pongo2.Template {
	tplMap := make(map[string]*pongo2.Template)

	if opt.TemplateFileSystem == nil {
		opt.TemplateFileSystem = macaron.NewTemplateFileSystem(macaron.RenderOptions{
			Directory:  opt.Directory,
			Extensions: opt.Extensions,
		}, true)
	}

	for _, f := range opt.TemplateFileSystem.ListFiles() {
		t, err := pongo2.FromFile(path.Join(opt.Directory, f.Name()) + f.Ext())
		if err != nil {
			// Bomb out if parse fails. We don't want any silent server starts.
			log.Fatalf("\"%s\": %v", f.Name(), err)
		}
		tplMap[strings.Replace(f.Name(), "\\", "/", -1)] = t
	}

	return tplMap
}

// templateSet represents a template set of type *pongo2.Template.
type templateSet struct {
	lock sync.RWMutex
	sets map[string]map[string]*pongo2.Template
	dirs map[string]string
}

func newTemplateSet() *templateSet {
	return &templateSet{
		sets: make(map[string]map[string]*pongo2.Template),
		dirs: make(map[string]string),
	}
}

func (ts *templateSet) Set(name string, opt *Options) map[string]*pongo2.Template {
	tplMap := compile(*opt)

	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.sets[name] = tplMap
	ts.dirs[name] = opt.Directory
	return tplMap
}

func (ts *templateSet) Get(setName, tplName string) (*pongo2.Template, error) {
	ts.lock.RLock()
	defer ts.lock.RUnlock()

	set := ts.sets[setName]
	if set == nil {
		return nil, fmt.Errorf("pongo2: template set \"%s\" is undefined", setName)
	}
	t := set[tplName]
	if t == nil {
		return nil, fmt.Errorf("pongo2: template \"%s\" is undefined", tplName)
	}

	return t, nil
}

func (ts *templateSet) GetDir(name string) string {
	ts.lock.RLock()
	defer ts.lock.RUnlock()

	return ts.dirs[name]
}

// Options represents a struct for specifying configuration options for the Render middleware.
type Options struct {
	// Directory to load templates. Default is "templates"
	Directory string
	// Extensions to parse template files from. Defaults to [".tmpl", ".html"]
	Extensions []string
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON
	IndentJSON bool
	// Outputs human readable XML
	IndentXML bool
	// Prefixes the JSON output with the given bytes.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
	// TemplateFileSystem is the interface for supporting any implmentation of template file system.
	macaron.TemplateFileSystem
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".tmpl", ".html"}
	}
	if len(opt.HTMLContentType) == 0 {
		opt.HTMLContentType = ContentHTML
	}

	return opt
}

func renderHandler(opt Options, tplSets []string) macaron.Handler {
	cs := macaron.PrepareCharset(opt.Charset)
	ts := newTemplateSet()
	ts.Set(_DEFAULT_TPL_SET_NAME, &opt)

	var tmpOpt Options
	for _, tplSet := range tplSets {
		tplName, tplDir := macaron.ParseTplSet(tplSet)
		tmpOpt = opt
		tmpOpt.Directory = tplDir
		ts.Set(tplName, &tmpOpt)
	}

	return func(ctx *macaron.Context) {
		r := &render{
			TplRender: &macaron.TplRender{
				ResponseWriter: ctx.Resp,
				Opt: &macaron.RenderOptions{
					IndentJSON: opt.IndentJSON,
					IndentXML:  opt.IndentXML,
					PrefixJSON: opt.PrefixJSON,
					PrefixXML:  opt.PrefixXML,
				},
				CompiledCharset: cs,
			},
			templateSet:     ts,
			opt:             &opt,
			compiledCharset: cs,
		}
		ctx.Render = r
		ctx.MapTo(r, (*macaron.Render)(nil))
	}
}

// Pongoer is a Middleware that maps a macaron.Render service into the Macaron handler chain.
// An single variadic pongo2.Options struct can be optionally provided to configure
// HTML rendering. The default directory for templates is "templates" and the default
// file extension is ".tmpl" and ".html".
//
// If MACARON_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MACARON_ENV environment variable to "production".
func Pongoer(options ...Options) macaron.Handler {
	return renderHandler(prepareOptions(options), []string{})
}

func Pongoers(options Options, tplSets ...string) macaron.Handler {
	return renderHandler(prepareOptions([]Options{options}), tplSets)
}

type render struct {
	*macaron.TplRender
	*templateSet
	opt             *Options
	compiledCharset string

	startTime time.Time
}

func data2Context(data interface{}) pongo2.Context {
	return pongo2.Context(data.(map[string]interface{}))
}

func (r *render) renderHTML(status int, setName, tplName string, data interface{}) {
	r.startTime = time.Now()

	t, err := r.templateSet.Get(setName, tplName)
	if macaron.Env == macaron.DEV {
		opt := *r.opt
		opt.Directory = r.templateSet.GetDir(setName)
		r.templateSet.Set(setName, &opt)
		t, err = r.templateSet.Get(setName, tplName)
	}
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Header().Set(ContentType, r.opt.HTMLContentType+r.compiledCharset)
	r.WriteHeader(status)
	if err := t.ExecuteWriter(data2Context(data), r); err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *render) HTML(status int, name string, data interface{}, _ ...macaron.HTMLOptions) {
	r.renderHTML(status, _DEFAULT_TPL_SET_NAME, name, data)
}

func (r *render) HTMLSet(status int, setName, tplName string, data interface{}, _ ...macaron.HTMLOptions) {
	r.renderHTML(status, setName, tplName, data)
}

func (r *render) HTMLSetBytes(setName, tplName string, data interface{}, _ ...macaron.HTMLOptions) ([]byte, error) {
	t, err := r.templateSet.Get(setName, tplName)
	if macaron.Env == macaron.DEV {
		opt := *r.opt
		opt.Directory = r.templateSet.GetDir(setName)
		r.templateSet.Set(setName, &opt)
		t, err = r.templateSet.Get(setName, tplName)
	}
	if err != nil {
		return []byte(""), err
	}

	return t.ExecuteBytes(data2Context(data))
}

func (r *render) HTMLBytes(name string, data interface{}, _ ...macaron.HTMLOptions) ([]byte, error) {
	return r.HTMLSetBytes(_DEFAULT_TPL_SET_NAME, name, data)
}

func (r *render) HTMLSetString(setName, tplName string, data interface{}, _ ...macaron.HTMLOptions) (string, error) {
	p, err := r.HTMLSetBytes(setName, tplName, data)
	return string(p), err
}

func (r *render) HTMLString(name string, data interface{}, _ ...macaron.HTMLOptions) (string, error) {
	p, err := r.HTMLBytes(name, data)
	return string(p), err
}

func (r *render) SetTemplatePath(setName, dir string) {
	if len(setName) == 0 {
		setName = _DEFAULT_TPL_SET_NAME
	}
	opt := *r.opt
	opt.Directory = dir
	r.templateSet.Set(setName, &opt)
}

func (r *render) HasTemplateSet(name string) bool {
	_, ok := r.templateSet.sets[name]
	return ok
}
