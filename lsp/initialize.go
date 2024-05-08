package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

// handling server lifecyle in spec
type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type (
	ServerCapabilities struct {
		TestDocumentSync   int  `json:"textDocumentSync"` // value 1 says send the whole document everytime
		HoverProvider      bool `json:"hoverProvider"`
		DefinitionProvider bool `json:"definitionProvider"`
		CodeActionProvider bool `json:"codeActionProvider"`
	}
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TestDocumentSync:   1, // 1 says send the whole document everytime
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "Learning Lsp",
				Version: "0.0.0.0.0-beta.final",
			},
		},
	}
}
