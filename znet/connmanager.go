package znet

import (
	"errors"
	"fmt"
	"gorik/ziface"
	"sync"
)

/*
Connection manager module
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection // Map to hold connection information
	connLock    sync.RWMutex                  // Read-write lock for concurrent access to the map
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// Protect shared resource (map) with a write lock
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// Add the connection to ConnManager
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("Connection added to ConnManager successfully: conn num =", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// Protect shared resource (map) with a write lock
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// Remove the connection
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("Connection removed: ConnID =", conn.GetConnID(), "successfully: conn num =", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// Protect shared resource (map) with a read lock
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// Get the current number of connections
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	// Protect shared resource (map) with a write lock
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// Stop and remove all connections
	for connID, conn := range connMgr.connections {
		// Stop the connection
		conn.Stop()
		// Remove the connection
		delete(connMgr.connections, connID)
	}

	fmt.Println("All connections cleared successfully: conn num =", connMgr.Len())
}
