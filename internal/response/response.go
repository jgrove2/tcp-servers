package response

import (
	"fmt"
	"io"

	"github.com/jgrove2/tcp-servers/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK          StatusCode = 200
	StatusCodeBadRequest  StatusCode = 400
	StatusCodeServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusCodeOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
		break
	case StatusCodeBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
		break
	case StatusCodeServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header["content-length"] = fmt.Sprintf("%d", contentLen)
	header["connection"] = "close"
	header["content-type"] = "text/plain"
	return header
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for k, v := range h {
		str := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(str))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte{'\r', '\n'})
	if err != nil {
		return err
	}
	return nil
}
