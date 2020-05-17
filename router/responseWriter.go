package router

import (
	"bytes"
	"net/http"
)

type responseWriter struct {
	resp *http.Response
}

func NewResponseWriter() responseWriter {
	resp := http.Response{
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{},
		Close:         false,
		ContentLength: 0,
	}
	return responseWriter{resp: &resp}
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb ClosingBuffer) Close() (err error) {
	//we don't actually have to do anything here, since the buffer is just some data in memory
	//and the error is initialized to no-error
	return
}

func (rw responseWriter) Header() http.Header {
	return rw.resp.Header
}

func (rw responseWriter) WriteHeader(statusCode int) {
	rw.resp.StatusCode = statusCode
}

func (rw responseWriter) Write(data []byte) (int, error) {
	if rw.resp.StatusCode == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	buf := new(bytes.Buffer)
	if rw.resp.Body != nil {
		buf.ReadFrom(rw.resp.Body)
	}

	n, err := buf.Write(data)
	if err != nil {
		return 0, err
	}

	cr := ClosingBuffer{
		Buffer: buf,
	}

	rw.resp.Body = cr

	rw.resp.ContentLength += int64(len(data))

	return n, nil
}
