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
import "net/http/httptest"
import "testing"

import "encoding/base64"

import "github.com/milochristiansen/axis2"
import "github.com/milochristiansen/axis2/sources/zip"

func TestStaticOcclusion(t *testing.T) {
	server := getTestServer(t)

	req, err := http.NewRequest("GET", "/dont-serve.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler, pattern := server.Handlers.Handler(req)
	if pattern == "" {
		// This should be impossible, because the / handler will catch it.
		t.Fatal("No handler found.")
	}
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Wrong response. Expected %v, got %v", http.StatusNotFound, rr.Code)
	}
}

func TestTemplates(t *testing.T) {
	server := getTestServer(t)

	req, err := http.NewRequest("GET", "/template", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler, pattern := server.Handlers.Handler(req)
	if pattern == "" {
		// This should be impossible, because the / handler will catch it.
		t.Fatal("No handler found.")
	}
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Wrong response. Expected %v, got %v", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "test" {
		t.Errorf("Wrong body. Expected %q, got %q", "test", rr.Body.String())
	}

	// Double check resource occlusion while we are here:
	req, err = http.NewRequest("GET", "/template.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler, pattern = server.Handlers.Handler(req)
	if pattern == "" {
		// This should be impossible, because the / handler will catch it.
		t.Fatal("No handler found.")
	}
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Wrong response. Expected %v, got %v", http.StatusNotFound, rr.Code)
	}
}

func TestStatic(t *testing.T) {
	server := getTestServer(t)

	req, err := http.NewRequest("GET", "/static.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler, pattern := server.Handlers.Handler(req)
	if pattern == "" {
		// This should be impossible, because the / handler will catch it.
		t.Fatal("No handler found.")
	}
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Wrong response. Expected %v, got %v", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "This is a static file" {
		t.Errorf("Wrong body. Expected %q, got %q", "This is a static file", rr.Body.String())
	}
}

// Helpers
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}

var test_server *Server

func getTestServer(t *testing.T) *Server {
	if test_server != nil {
		return test_server
	}

	// Load the test data into AXIS
	fs := new(axis2.FileSystem)
	dir, err := zip.NewRawDir(TestData)
	if err != nil {
		t.Fatal(err)
	}
	fs.Mount("resources", dir, false)

	// Set up a simple test server.
	err, test_server = Initialize(fs, "resources", []Handler{
		&TemplateHandler{
			Resources: []string{"template.html"},
			Template:  "template.html",
			Path:      "/template",
			Data: func(w http.ResponseWriter, r *http.Request) interface{} {
				return "test"
			},
		},
		&SimpleHandler{
			Resources: []string{"dont-serve.txt"},
			Path:      "/simple",
			Logic: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		},
	}, errorHandler)
	if err != nil {
		t.Fatal(err)
	}

	return test_server
}

var TestData []byte

func init() {
	var err error
	TestData, err = base64.StdEncoding.DecodeString("UEsDBBQAAAAAAJRIglGShN+lLwAAAC8AAAAOAAAAZG9udC1zZXJ2ZS50eHRUaGlzIGZpbGUgaXMgYSByZXNvdXJjZSBhbmQgd2lsbCBub3QgYmUgc2VydmVkLlBLAwQUAAAAAACVSIJRk/66FBUAAAAVAAAACgAAAHN0YXRpYy5jc3NUaGlzIGlzIGEgc3RhdGljIGZpbGVQSwMEFAAAAAAAOkeCUX4doh8HAAAABwAAAA0AAAB0ZW1wbGF0ZS5odG1se3sgLiB9fVBLAQI/AxQAAAAAAJRIglGShN+lLwAAAC8AAAAOAAAAAAAAAAAAAACkgQAAAABkb250LXNlcnZlLnR4dFBLAQI/AxQAAAAAAJVIglGT/roUFQAAABUAAAAKAAAAAAAAAAAAAACkgVsAAABzdGF0aWMuY3NzUEsBAj8DFAAAAAAAOkeCUX4doh8HAAAABwAAAA0AAAAAAAAAAAAAAKSBmAAAAHRlbXBsYXRlLmh0bWxQSwUGAAAAAAMAAwCvAAAAygAAAAAA")
	if err != nil {
		panic("Failed to decode test date. This should be impossible.")
	}
}
