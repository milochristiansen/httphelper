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
import "errors"
import "html/template"
import filepath "path"

import "github.com/milochristiansen/axis2"

type SimpleHandler struct {
	// AXIS paths for resources assigned to this Handler. You may use other resources as well,
	// but anything listed here will be marked off the list of files to serve statically.
	Resources []string

	// The handler logic. See also http.HandlerFunc.
	Logic http.Handler

	Path  string // The path this handler is responsible for.
	Loose bool   // If true do no automatically insert a check for path supersets.
}

func (h *SimpleHandler) initalize(fs *axis2.FileSystem, s *Server) error {
	err := handlerBoilerplate(h.Path, h.Resources, s)
	if err != nil {
		return err
	}

	if h.Loose {
		s.Handlers.Handle(h.Path, h.Logic)
	} else {
		s.Handlers.HandleFunc(h.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != h.Path {
				s.errhandler(w, r, http.StatusNotFound)
				return
			}

			h.Logic.ServeHTTP(w, r)
		})
	}
	return nil
}

func handlerBoilerplate(path string, resources []string, s *Server) error {
	s.log.e.Println("Building handler for ", path)

	if s.hasHandler[path] {
		s.log.e.Println("A handler for ", path, " already exists.")
		return errors.New("A handler for " + path + " already exists.")
	}
	s.hasHandler[path] = true

	for _, p := range resources {
		f, ok := s.Files[p]
		if !ok {
			s.log.e.Println("Resource ", p, " does not exist.")
			return errors.New("Resource " + p + " does not exist.")
		}
		f.Tags["Resource"] = true
	}
	return nil
}

func staticPageHandler(f *File, s *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != f.FullPath() {
			s.errhandler(w, r, http.StatusNotFound)
			return
		}

		typ := mime.TypeByExtension(getExt(f.Name))
		if typ != "" {
			w.Header().Set("Content-Type", typ)
		}
		n, err := w.Write(f.Content) // This may error out, but we can't do anything about it (except maybe log it) at this point.
		if err != nil {
			s.log.e.Println("Error in static page handler: ", err, " bytes written: ", n)
		}
	}
}

type TemplateHandler struct {
	// AXIS paths for resources assigned to this Handler. You may use other resources as well,
	// but anything listed here will be marked off the list of files to serve statically.
	Resources []string
	Template  string // The AXIS path to the template file (also list in Resources)

	// Return the data object the template needs to operate.
	Data func(w http.ResponseWriter, r *http.Request) interface{}

	Path string // The path this handler is responsible for.

	page *template.Template
	name string
}

// TemplateHandler returns a http.Handler backed by a template. If 'data' returns nil, then nothing is
// done. It is assumed that you have already handled the error.
// TODO: This is unworkable. Returning an error here is super ugly, and it is simply impossible to have a Server when this is needed.
func (h *TemplateHandler) initalize(fs *axis2.FileSystem, s *Server) error {
	err := handlerBoilerplate(h.Path, h.Resources, s)
	if err != nil {
		return err
	}

	h.name = stripExt(filepath.Base(h.Path))

	content, err := fs.ReadAll(h.Path)
	if err != nil {
		s.log.e.Println("Error in TemplateHandler ", h.name, ": ", err)
		return err
	}

	h.page, err = template.New(h.name).Parse(string(content))
	if err != nil {
		s.log.e.Println("Error in TemplateHandler ", h.name, ": ", err)
		return err
	}

	s.Handlers.HandleFunc(h.Path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != h.Path {
			s.errhandler(w, r, http.StatusNotFound)
			return
		}

		d := h.Data(w, r)
		if d == nil {
			return
		}
		err := h.page.Execute(w, d)
		if err != nil {
			s.log.e.Println("Error in TemplateHandler ", h.name, ": ", err)
		}
	})

	return nil
}
