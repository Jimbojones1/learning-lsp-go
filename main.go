package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	writer := os.Stdout
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
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
		writeResponse(writer, msg)
		logger.Print("Sent the reply")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error in textDocument/didOpen %s", err)
			return
		}

		logger.Printf("Connected to: %s  ", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error in textDocument/didChange %s", err)
			return
		}

		logger.Printf("Changed to: %s  ", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			fmt.Println(change, " <= this is change!")
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error in textDocument/hover %s", err)
			return
		}
		// create Response
		fmt.Println(request.Params)
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		// write it back
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error in textDocument/definition %s", err)
			return
		}
		// create Response
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		// write it back
		writeResponse(writer, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s", err)
			return
		}

		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/completion: %s", err)
			return
		}

		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)

	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
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
