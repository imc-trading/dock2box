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

package pongo2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Unknwon/macaron"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Render_HTML(t *testing.T) {
	Convey("Render HTML", t, func() {
		m := macaron.Classic()
		m.Use(Pongoers(Options{
			Directory: "fixtures/basic",
		}, "fixtures/basic2"))
		m.Get("/foobar", func(r macaron.Render) {
			r.HTML(200, "hello", map[string]interface{}{
				"Name": "jeremy",
			})
			r.SetTemplatePath("", "fixtures/basic2")
		})
		m.Get("/foobar2", func(r macaron.Render) {
			if r.HasTemplateSet("basic2") {
				r.HTMLSet(200, "basic2", "hello", map[string]interface{}{
					"Name": "jeremy",
				})
			}
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>")
		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>")

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/foobar2", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>What's up, jeremy</h1>")

		Convey("Change render templates path", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/foobar", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Header().Get(ContentType), ShouldEqual, ContentHTML+"; charset=UTF-8")
			So(resp.Body.String(), ShouldEqual, "<h1>What's up, jeremy</h1>")
		})
	})

	Convey("Render HTML and return string", t, func() {
		m := macaron.Classic()
		m.Use(Pongoers(Options{
			Directory: "fixtures/basic",
		}, "basic2:fixtures/basic2"))
		m.Get("/foobar", func(r macaron.Render) {
			result, err := r.HTMLString("hello", "jeremy")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "<h1>Hello jeremy</h1>")
		})
		m.Get("/foobar2", func(r macaron.Render) {
			result, err := r.HTMLSetString("basic2", "hello", "jeremy")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "<h1>What's up, jeremy</h1>")
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/foobar2", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
	})

	Convey("Render bad HTML", t, func() {
		m := macaron.Classic()
		m.Use(Pongoer(Options{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r macaron.Render) {
			r.HTML(200, "nope", nil)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		So(resp.Body.String(), ShouldEqual, "pongo2: template \"nope\" is undefined\n")
	})
}

func Test_Render_XHTML(t *testing.T) {
	Convey("Render XHTML", t, func() {
		m := macaron.Classic()
		m.Use(Pongoer(Options{
			Directory:       "fixtures/basic",
			HTMLContentType: ContentXHTML,
		}))
		m.Get("/foobar", func(r macaron.Render) {
			r.HTML(200, "hello", map[string]interface{}{
				"Name": "jeremy",
			})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get(ContentType), ShouldEqual, ContentXHTML+"; charset=UTF-8")
		So(resp.Body.String(), ShouldEqual, "<h1>Hello jeremy</h1>")
	})

	m := macaron.Classic()
	m.Use(Pongoer(Options{
		Directory:       "fixtures/basic",
		HTMLContentType: ContentXHTML,
	}))
}

func Test_Render_Extensions(t *testing.T) {
	Convey("Render with extensions", t, func() {
		m := macaron.Classic()
		m.Use(Pongoer(Options{
			Directory:  "fixtures/basic",
			Extensions: []string{".tmpl", ".html"},
		}))
		m.Get("/foobar", func(r macaron.Render) {
			r.HTML(200, "hypertext", map[string]interface{}{})
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/foobar", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)

		So(resp.Body.String(), ShouldEqual, "Hypertext!")
	})
}

func Test_Render_NoRace(t *testing.T) {
	Convey("Make sure render has no race", t, func() {
		m := macaron.Classic()
		m.Use(Pongoer(Options{
			Directory: "fixtures/basic",
		}))
		m.Get("/foobar", func(r macaron.Render) {
			r.HTML(200, "hello", map[string]interface{}{
				"Name": "world",
			})
		})

		done := make(chan bool)
		doreq := func() {
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/foobar", nil)
			m.ServeHTTP(resp, req)
			done <- true
		}
		// Run two requests to check there is no race condition
		go doreq()
		go doreq()
		<-done
		<-done
	})
}
