package mocker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

type Mocker struct {
	Router http.Handler
}

type Header map[string]string

func New(router http.Handler) *Mocker {
	return &Mocker{
		Router: router,
	}
}

func (c *Mocker) Request(method string, path string, body []byte, header *Header) *httptest.ResponseRecorder {
	reader := bytes.NewReader(body)
	req, _ := http.NewRequest(method, path, reader)

	if header != nil {
		for key, value := range *header {
			req.Header.Set(key, value)
		}
	}

	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)
	return w
}

func (c *Mocker) Head(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodHead, path, body, header)
}

func (c *Mocker) Options(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodOptions, path, body, header)
}

func (c *Mocker) Get(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodGet, path, body, header)
}

func (c *Mocker) Put(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodPut, path, body, header)
}

func (c *Mocker) Post(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodPost, path, body, header)
}

func (c *Mocker) Delete(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodDelete, path, body, header)
}

func (c *Mocker) Patch(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodPatch, path, body, header)
}

func (c *Mocker) Trace(path string, body []byte, header *Header) *httptest.ResponseRecorder {
	return c.Request(http.MethodTrace, path, body, header)
}
