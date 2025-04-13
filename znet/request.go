package znet

import "gorik/ziface"

type Request struct {
	conn ziface.IConnection // The connection already established with the client
	msg  ziface.IMessage    // Data requested by the client
}

// Get the connection information
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// Get the data of the request message
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// Get the ID of the request message
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
