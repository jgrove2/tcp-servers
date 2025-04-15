package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/jgrove2/tcp-servers/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Method string

const (
	POST    Method = "POST"
	GET     Method = "GET"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	HEAD    Method = "HEAD"
	OPTIONS Method = "OPTIONS"
	CONNECT Method = "CONNECT"
	TRACE   Method = "TRACE"
	PATCH   Method = "PATCH"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 0, 8) // initial buffer capacity
	tmp := make([]byte, 8)
	req := &Request{RequestLine: RequestLine{}, state: 0, Headers: headers.NewHeaders()}

	for req.state != 3 {
		consumed, err := req.parse(buf)
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			// Shift unparsed data to the front
			copy(buf, buf[consumed:])
			buf = buf[:len(buf)-consumed]
			continue
		}
		if err == io.EOF {
			break
		}

		n, err := reader.Read(tmp)
		if err != nil {
			if err == io.EOF {
				if len(buf) == 0 && req.state != 3 {
					return nil, errors.New("unexpected EOF while parsing request")
				}
				// We'll keep processing whatever we've got
			} else {
				return nil, err
			}
		}

		buf = append(buf, tmp[:n]...)

		consumed, err = req.parse(buf)
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			// Shift unparsed data to the front
			copy(buf, buf[consumed:])
			buf = buf[:len(buf)-consumed]
		}

		if err == io.EOF {
			break
		}
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	// Exit if the request is already fully parsed
	if r.state == 3 {
		return 0, nil
	}
	in := strings.Index(string(data), "\r\n")
	if len(data) == 2 && in == 0 {
		r.state = 3
		return 2, nil
	}
	// Parse the request line
	if r.state == 0 {
		index := strings.Index(string(data), "\r\n")
		if index == -1 {
			return 0, nil
		}
		err := r.parseRequestLine(string(data[:index]))
		if err != nil {
			return 0, err
		}
		r.state = 1
		return index + 2, nil
		// Parse the headers
	} else if r.state == 1 {
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			if r.Headers.Get("content-length") != "" {
				r.state = 2
				return n + 2, nil
			}
			r.state = 3
		}
		return n, nil
		// Parse the body
	} else if r.state == 2 {
		// Retrieve content-length from headers
		contentLengthStr := r.Headers.Get("content-length")
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, err
		}

		// Calculate how much data is still needed for the body
		bodyLength := len(r.Body)
		remaining := contentLength - bodyLength

		// throw an error if there is no more data and contentLength is still there
		if bodyLength == 0 && len(data) == 0 {
			r.state = 3
			return 0, errors.New("no data available for body")
		}

		// If body length exceeds content-length, return an error
		if bodyLength+len(data) > contentLength {
			return 0, errors.New("body length is greater than content length")
		}

		if remaining > 0 {
			// Append only the remaining required bytes to the body
			if remaining > len(data) {
				r.Body = append(r.Body, data...)
				return len(data), nil // Return the number of bytes appended
			}

			r.Body = append(r.Body, data[:remaining]...)
			if len(r.Body) == contentLength {
				r.state = 3 // Body fully parsed, move to done state
			}
			return remaining, nil // Return the number of bytes that were consumed
		}

		// If we've fully parsed the body, transition to done state
		r.state = 3
		return 0, nil
	}

	return 0, nil
}

func (r *Request) parseRequestLine(line string) error {
	parsed := strings.Split(line, " ")
	if len(parsed) != 3 {
		return fmt.Errorf("invalid request line: %s", line)
	}
	httpVersion := strings.Split(parsed[2], "/")
	method := parsed[0]
	switch method {
	case string(POST):
	case string(GET):
	case string(PUT):
	case string(DELETE):
	case string(HEAD):
	case string(OPTIONS):
	case string(CONNECT):
	case string(TRACE):
	case string(PATCH):
		break
	default:
		return fmt.Errorf("invalid method %s", method)
	}
	r.RequestLine.HttpVersion = httpVersion[1]
	r.RequestLine.RequestTarget = parsed[1]
	r.RequestLine.Method = method
	log.Println(r.RequestLine)
	return nil
}
