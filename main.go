package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"learning-lsp/analysis"
	"learning-lsp/lsp"
	"learning-lsp/rpc"
	"log"
	"os"
)

func main() {
	fmt.Println("Hi")
	logger := getLogger("/Users/jameshaff/go/src/learning-lsp/log.txt")
	logger.Println("Logger Started!")
	scanner := bufio.NewScanner(os.Stdin)
	// rpc.Split satisifies Split type
	scanner.Split(rpc.Split)
	// create a new State for the application
	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
	logger.Printf("Msg recieved with a method of %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("hey we couldn't parse this %s", err)
		}

		logger.Printf("Connected to: %s %s ", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		// lets replay
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Print("Sent the reply")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error in textDocument/didOpen %s", err)
		}

		logger.Printf("Connected to: %s %s ", request.Params.TextDocument.Text, request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error in textDocument/didChange %s", err)
		}

		logger.Printf("Helloooooo to: %s %s ", request.Params.ContentChanges, request.Params.TextDocument.URI)

	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		panic("hey this isn't a good file!")
	}
	// log.Lshortfile gives you filename:line
	return log.New(logfile, "[learning-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
