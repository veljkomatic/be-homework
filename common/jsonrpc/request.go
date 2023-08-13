package jsonrpc

import (
	"encoding/json"
	"sync/atomic"
)

var id int64

// Request is the identified call of the method.
type Request struct {

	// ID is the unique identifier of this operation request.
	// If a client needs to identify the result of the operation execution,
	// the ID should be passed by the client, then it is guaranteed
	// that the client will receive the result frame with the same id.
	// If ID == 0 this is treated as a notification.
	ID ID `json:"id,omitempty"`

	// Version of a request, set to "2.0" as per specification.
	Version string `json:"jsonrpc"`

	// Method is the name which will be proceeded by this request.
	Method string `json:"method"`

	// Params are parameters which are needed for operation execution.
	Params json.RawMessage `json:"params,omitempty"`
}

// NewRequest creates a new request with the passed method and params.
// Every time NewRequest is called a counter is incremented and the request returned has an incremented id.
// Version field of request is always "2.0" as per jsonrpc specification.
func NewRequest(method string, params json.RawMessage) *Request {
	return &Request{
		ID:      NewID(atomic.AddInt64(&id, 1)),
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
}
