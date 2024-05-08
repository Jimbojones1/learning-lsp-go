package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// will specify the type of the params in all the requests later
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"` // id could be null

	// Result
	// Error
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
