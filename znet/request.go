package znet

import "gorik/ziface"

type Request struct {
	conn ziface.IConnection
	data []byte
}

// GetConenction implements IRequest.
func (r *Request) GetConenction() IConnection {
	return r.conn
}

// GetData implements IRequest.
func (r *Request) GetData() []byte {
	return r.data
}
