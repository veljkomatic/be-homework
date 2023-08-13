package jsonrpc

type ID int64

func NewID(id int64) ID {
	return ID(id)
}
