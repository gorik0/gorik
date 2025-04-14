package ziface

type IServer interface {
	// Start the server
	Start()
	// Stop the server
	Stop()
	// Start the business service
	Serve()
	// Register a routing business method for the current server to be used for client connection processing
	AddRouter(msgID uint32, router IRouter)
	// Get the connection manager
	GetConnMgr() IConnManager
	// Set the hook function to be called when a connection is created for this server
	SetOnConnStart(func(IConnection))
	// Set the hook function to be called when a connection is about to be disconnected for this server
	SetOnConnStop(func(IConnection))
	// Invoke the OnConnStart hook function for the connection
	CallOnConnStart(conn IConnection)
	// Invoke the OnConnStop hook function for the connection
	CallOnConnStop(conn IConnection)
}
