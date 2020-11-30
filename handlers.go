/*
Copyright 2020 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package httphelper

import "net/http"
import "mime"
import "html/template"
import filepath "path"

import "github.com/milochristiansen/axis2"

func staticPageHandler(f *File, s *Server, errhandler HTTPErrorHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != f.FullPath() {
			errhandler(w, r, http.StatusNotFound)
			return
		}

		typ := mime.TypeByExtension(getExt(f.Name))
		if typ != "" {
			w.Header().Set("Content-Type", typ)
		}
		n, err := w.Write(f.Content) // This may error out, but we can't do anything about it (except maybe log it) at this point.
		if err != nil {
			s.log.ErrPrintln("Error in static page handler: ", err, " bytes written: ", n)
		}
	}
}

type templateHandler struct {
	page *template.Template
	data func(w http.ResponseWriter, r *http.Request) interface{}
	s    *Server
	name string
}

// TemplateHandler returns a http.Handler backed by a template. If 'data' returns nil, then nothing is
// done. It is assumed that you have already handled the error.
func TemplateHandler(fs *axis2.FileSystem, path string, s *Server, data func(w http.ResponseWriter, r *http.Request) interface{}) (error, http.Handler) {
	th := &templateHandler{}
	th.name = stripExt(filepath.Base(path))

	content, err := fs.ReadAll(path)
	if err != nil {
		s.log.ErrPrintln("Error in TemplateHandler ", th.name, ": ", err)
		return err, nil
	}

	th.page, err = template.New(th.name).Parse(string(content))
	if err != nil {
		s.log.ErrPrintln("Error in TemplateHandler ", th.name, ": ", err)
		return err, nil
	}

	th.data = data
	th.s = s
	return nil, th
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := th.data(w, r)
	if d == nil {
		return
	}
	err := th.page.Execute(w, d)
	if err != nil {
		th.s.log.ErrPrintln("Error in TemplateHandler ", th.name, ": ", err)
	}
}
