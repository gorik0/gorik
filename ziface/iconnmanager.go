package ziface

/*
Connection management abstract layer
*/
type IConnManager interface {
	Add(conn IConnection)                   // Add a connection
	Remove(conn IConnection)                // Remove a connection
	Get(connID uint32) (IConnection, error) // Get a connection using the connection ID
	Len() int                               // Get the current number of connections
	ClearConn()                             // Remove and stop all connections
}
