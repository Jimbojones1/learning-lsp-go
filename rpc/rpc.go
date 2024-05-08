package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("did not found separator")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}
	return baseMessage.Method, content[:contentLength], nil
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	// Content-Length: ...\r\n
	if !found {
		// this means we are not ready yet, just keep going
		return 0, nil, nil
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))

	if len(content) < contentLength {
		// this means we haven't read enough bytes yet, keop waiting
		return 0, nil, nil
	}

	if err != nil {
		// there was no number for Content-Length
		return 0, nil, err
	}
	// +4 for the seperator \r\n\r\n
	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
