package headers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Headers map[string]string

type IHeaders interface {
	Parse(data []byte) (n int, done bool, err error)
	Get(key string) string
}

func NewHeaders() Headers {
	return Headers{}
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
	clrf := []byte{'\r', '\n'}
	lineEnd := bytes.Index(data, clrf)
	log.Println(lineEnd)

	if lineEnd == -1 {
		return 0, false, nil
	}

	if lineEnd == 0 {
		return 0, true, nil
	}

	line := string(data[:lineEnd])

	// Check for malformed spacing (no space before colon allowed)
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return 0, false, fmt.Errorf("malformed header: missing colon")
	}

	if data[colonIndex-1] == ' ' {
		return 0, false, fmt.Errorf("malformed header: space before colon")
	}

	// Trim whitespace from key and value
	key := strings.TrimSpace(line[:colonIndex])
	value := strings.TrimSpace(line[colonIndex+1:])

	if !validateHeaderFieldName(key) {
		return 0, false, errors.New("Field must only contain A-Z a-z 0-9 !, @ #, $, %, ^, &, *, (, ), -, _, +, ., /, :, ;, <, >, ?, |, ~")
	}

	key = strings.ToLower(key)

	if (*h)[key] != "" {
		(*h)[key] = fmt.Sprintf("%s, %s", (*h)[key], value)
	} else {
		(*h)[strings.ToLower(key)] = value
	}

	return lineEnd + 2, false, nil

}

func (h *Headers) Get(key string) string {
	return (*h)[strings.ToLower(key)]
}

func validateHeaderFieldName(field string) bool {
	pattern := `^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(field)
}
