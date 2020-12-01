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
import "strings"
import "errors"

import "github.com/milochristiansen/axis2"

// This uses a tag system much like Rubble.

// Server is a convenient holder for the HTTP handlers and the loaded files generated by Initialize.
type Server struct {
	//handlers []Handler
	Files    map[string]*File
	Handlers *http.ServeMux

	log        *logger
	hasHandler map[string]bool
	errhandler HTTPErrorHandler
}

type Handler interface {
	initalize(fs *axis2.FileSystem, s *Server) error
}

type File struct {
	Name    string // File name.
	Source  string // File path (AXIS syntax, including loc ids).
	Content []byte
	Tags    map[string]bool
}

func (f *File) FullPath() string {
	if f.Source == "" {
		return f.Name
	}
	return f.Source + "/" + f.Name
}

// HTTPErrorHandler is a superset of an HTTP handler that also takes a status code. Called whenever the server
// detects an error. Currently this is only called with 404 errors.
type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request, status int)

// Initialize creates a new Server based on the given data directory and handlers.
//
// If there is no handler for "/" one will automatically be created that simply calls the error handler with a 404.
//
// The Loggers are optional. If you provide one logger it will be used by everything. Two will be used for info and
// errors. Only the first two will be used. You may pass nil for any Logger, in which case that kind of message will
// not be logged.
func Initialize(fs *axis2.FileSystem, path string, handlers []Handler, errhandler HTTPErrorHandler, log ...Logger) (error, *Server) {
	s := &Server{}

	s.log = &logger{}
	switch len(log) {
	case 0:
		// Do nothing. RIP logging.
	case 1:
		s.log.i = log[0]
		s.log.e = log[0]
	default:
		s.log.i = log[0]
		s.log.e = log[1]
	}
	if s.log.i == nil {
		s.log.i = dummyLogger
	}
	if s.log.e == nil {
		s.log.e = dummyLogger
	}

	s.errhandler = errhandler

	// First build a tree of resources
	s.log.i.Println("Building data tree.")
	err := loadDir(fs, path, s)
	if err != nil {
		s.log.e.Println("Error: ", err, " while building data tree.")
		return err, nil
	}

	// Then mark off anything with an handler and set up the handlers.
	s.log.i.Println("Initializing handlers.")
	s.hasHandler = map[string]bool{}
	s.Handlers = http.NewServeMux()
	for _, h := range handlers {
		err := h.initalize(fs, s)
		if err != nil {
			s.log.e.Println("Error: ", err, " while initializing handlers.")
			return err, nil
		}
	}

	if !s.hasHandler["/"] {
		s.Handlers.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			errhandler(w, r, http.StatusNotFound)
		})
	}

	// Finally create handlers for the remaining stuff
	for _, f := range s.Files {
		p := strings.TrimPrefix(f.FullPath(), path)

		s.log.e.Println("Building handler for ", p)
		if s.hasHandler[f.FullPath()] {
			s.log.e.Println("A handler for ", p, " already exists.")
			return errors.New("A handler for " + p + " already exists."), nil
		}

		s.Handlers.HandleFunc(f.FullPath(), staticPageHandler(f, s))
	}

	return nil, s
}

func loadDir(fs *axis2.FileSystem, path string, s *Server) error {
	dirpath := path
	if path != "" {
		path += "/"
	}

	for _, filepath := range fs.ListFiles(dirpath) {
		if strings.HasPrefix(filepath, ".") {
			continue
		}

		content, err := fs.ReadAll(path + filepath)
		if err != nil {
			return err
		}

		file := &File{filepath, dirpath, content, map[string]bool{}}
		tags := GetFileTags(filepath)
		for _, tag := range tags {
			file.Tags[tag] = true
		}
		s.Files[filepath] = file
	}

	for _, dir := range fs.ListDirs(dirpath) {
		if strings.HasPrefix(dirpath, ".") {
			continue
		}

		err := loadDir(fs, path+dir, s)
		if err != nil {
			return err
		}
	}
	return nil
}
