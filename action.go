package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Header struct {
	Key   string
	Value string
}

type errorResponse struct {
	Error string `json:"error"`
}

type Action func(r *http.Request) *Response

type Headers map[string]string

type Response struct {
	Status      int
	ContentType string
	Content     io.Reader
	Headers     Headers
}

func Error(status int, err error, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: bytes.NewBufferString(err.Error()),
		Headers: headers,
	}
}

func ErrorJSON(status int, err error, headers Headers) *Response {
	errResp := errorResponse{
		Error: err.Error(),
	}

	b, err := json.Marshal(errResp)

	if err != nil {
		return Error(http.StatusInternalServerError, err, headers)
	}
	return &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
		Headers:     headers,
	}
}

func Data(status int, content []byte, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: bytes.NewBuffer(content),
		Headers: headers,
	}
}

func DataJSON(status int, v interface{}, headers Headers) *Response {
	b, err := json.Marshal(v)

	if err != nil {
		return ErrorJSON(http.StatusInternalServerError, err, headers)
	}

	return &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
		Headers:     headers,
	}
}

func DataWithReader(status int, r io.Reader, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: r,
		Headers: headers,
	}
}

func (a Action) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if response := a(r); response != nil {
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
