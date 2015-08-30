package main

import (
	"encoding/json"

	"github.com/gocraft/web"
)

var EMPTY_GIF = []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\xff\xff\xff\x00\x00\x00!\xf9\x04\x00\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")

// Used for convenient API error handling with panic()
type APIError interface {
	Handle(res web.ResponseWriter, req *web.Request)
}

type JSON map[string]interface{}

type JSONError struct {
	data JSON
	code int
}

// Panic with a JSONError, which should be handled by APIErrorMiddleware
func returnError(data JSON, code int) {
	// TODO: report statsd metric
	panic(JSONError{data, code})
}

func (e JSONError) Handle(res web.ResponseWriter, req *web.Request) {
	output, err := json.Marshal(e.data)
	if err != nil {
		// failed to marshal the error response
		panic(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(e.code)
	res.Write(output)
}

func APIErrorMiddleware(res web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	defer func() {
		if err := recover(); err != nil {
			// Will re-panic if not an APIError, resulting in HTTP 500
			e := err.(APIError)
			e.Handle(res, req)
		}
	}()

	next(res, req)
}