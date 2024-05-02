package rpc_test

import (
	"fmt"
	"jims-lsp/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Expected: %s Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"Hi\"}"
	method, contentLength, err := rpc.DecodeMessage([]byte(incomingMessage))
	fmt.Println(method, "this is msg", contentLength)
	if err != nil {
		t.Fatal(err)
	}

	if contentLength != 15 {
		t.Fatalf("Expected: 16 got %d", contentLength)
	}

	if method != "Hi" {
		t.Fatalf("Expected: Hey, Got : %s", method)
	}
}
