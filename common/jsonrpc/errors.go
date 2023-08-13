package jsonrpc

// Error indicates any exceptional situation during operation execution,
type Error struct {
	// Code is the value indicating the certain error type.
	Code int `json:"code"`

	// Message is the description of this error.
	Message string `json:"message"`

	// Data any kind of data which provides additional
	// information about the error e.g. stack trace, error time.
	Data any `json:"data,omitempty"`
}
