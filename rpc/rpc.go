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

func DecodeMessage(msg []byte) (string, int, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return "", 0, errors.New("Did not found separator")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", 0, err
	}
	// TODO: WE'll get to this soon
	_ = content

	var baseMessage BaseMessage
	fmt.Println(content[:contentLength], contentLength, &baseMessage)
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", 0, err
	}
	return baseMessage.Method, contentLength, nil
}
