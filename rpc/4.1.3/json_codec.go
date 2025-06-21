package main

import (
	"encoding/json"
	"io"
	"net/rpc"
)

// JsonCodec implements rpc.Codec interface
type JsonCodec struct {
	conn io.ReadWriteCloser // Underlying connection
	enc  *json.Encoder      // JSON encoder
	dec  *json.Decoder      // JSON decoder
}

// NewJsonCodec creates a new JsonCodec
func NewJsonCodec(conn io.ReadWriteCloser) rpc.Codec {
	return &JsonCodec{
		conn: conn,
		enc:  json.NewEncoder(conn),
		dec:  json.NewDecoder(conn),
	}
}

// ReadRequestHeader reads the RPC request header
func (c *JsonCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.dec.Decode(r)
}

// ReadRequestBody reads the RPC request body
func (c *JsonCodec) ReadRequestBody(body interface{}) error {
	if body == nil {
		return nil // No request body
	}
	return c.dec.Decode(body)
}

// WriteResponse writes the RPC response (header and body)
func (c *JsonCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	// Encode response header
	if err := c.enc.Encode(r); err != nil {
		return err
	}
	// Encode response body
	if err := c.enc.Encode(body); err != nil {
		return err
	}
	return nil
}

// ReadResponseHeader reads the RPC response header
func (c *JsonCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

// ReadResponseBody reads the RPC response body
func (c *JsonCodec) ReadResponseBody(body interface{}) error {
	if body == nil {
		return nil // No response body
	}
	return c.dec.Decode(body)
}

// Close closes the underlying connection
func (c *JsonCodec) Close() error {
	return c.conn.Close()
}
