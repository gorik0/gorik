package ziface

import "net"

type IConnection interface {
	// Start the connection, making the current connection work
	Start()
	// Stop the connection, ending the current connection state
	Stop()
	// Get the raw socket TCPConn from the current connection
	GetTCPConnection() *net.TCPConn
	// Get the current connection ID
	GetConnID() uint32
	// Get the remote client's address information
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error
	SendBuffMsg(msgID uint32, data []byte) error
	SetProperty(key string, value interface{})
	// Get connection attributes
	GetProperty(key string) (interface{}, error)
	// Remove connection attributes
	RemoveProperty(key string)
}
type HandFunc func(*net.TCPConn, []byte, int) error
