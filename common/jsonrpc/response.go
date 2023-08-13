package jsonrpc

import "encoding/json"

// Response is a jsonrpc response object
type Response struct {
	// ID is the unique identifier of the request.
	ID ID `json:"id,omitempty"`

	// Version of a request, set to "2.0" as per specification.
	Version string `json:"jsonrpc"`

	// Result is the result of the method call.
	// Result and Error are mutually exclusive.
	Result json.RawMessage `json:"result,omitempty"`

	// Error is an error returned by the server if the RPC failed.
	Error *Error `json:"error,omitempty"`
}
