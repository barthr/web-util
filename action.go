package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

// HandlerFunc is a wrapper for the handler interface
type HandlerFunc func(r *http.Request) *Response

// Implement the http.Handle interface
// We handle the writing of the headers and the body inside the action method
// Which will be called by the router
func (hf HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if response := hf(r); response != nil {
		if response.ContentType != "" {
			rw.Header().Set("Content-Type", response.ContentType)
		}
		for k, v := range response.Headers {
			rw.Header().Set(k, v)
		}
		rw.WriteHeader(response.Status)
		_, err := io.Copy(rw, response.Content)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

// Headers are the headers used for sending to the callee
type Headers map[string]string

// Response is the actual response used to write to the response writer
type Response struct {
	Status      int
	ContentType string
	Content     io.Reader
	Headers     Headers
}

// Error returns create's an Response object with the content set to the error
func Error(status int, err error, headers ...Headers) *Response {
	r := &Response{
		Status:  status,
		Content: bytes.NewBufferString(err.Error()),
	}
	if len(headers) > 0 {
		r.Headers = headers[0]
	}
	return r
}

// ErrorJSON create's an Response object with the content set to the error
// and encoded in JSON
func ErrorJSON(status int, err error, headers ...Headers) *Response {
	errResp := errorResponse{
		Error: err.Error(),
	}

	b, err := json.Marshal(errResp)

	if err != nil {
		return Error(http.StatusInternalServerError, err, headers...)
	}
	r := &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
	}
	if len(headers) > 0 {
		r.Headers = headers[0]
	}
	return r
}

// Data create's an Response object with the content set to the passed byte array content
func Data(status int, content []byte, headers ...Headers) *Response {
	r := &Response{
		Status:  status,
		Content: bytes.NewBuffer(content),
	}
	if len(headers) > 0 {
		r.Headers = headers[0]
	}
	return r
}

// JSON create's an Response object with the content set to the encoded json value of v
func JSON(status int, v interface{}, headers ...Headers) *Response {
	b, err := json.Marshal(v)

	if err != nil {
		return ErrorJSON(http.StatusInternalServerError, err, headers...)
	}

	r := &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
	}
	if len(headers) > 0 {
		r.Headers = headers[0]
	}
	return r
}

// WithReader create's an Response object with the content set to the given reader
func WithReader(status int, reader io.Reader, headers ...Headers) *Response {
	r := &Response{
		Status:  status,
		Content: reader,
	}
	if len(headers) > 0 {
		r.Headers = headers[0]
	}
	return r
}
