package znet

import (
	"fmt"
	"gorik/ziface"
	"net"
)

type Connection struct {
	Conn *net.TCPConn
	// Current connection's ID, also known as SessionID, ID is   globally unique
	ConnID uint32
	// Current connection's close status
	isClosed bool

	// The handle function of this connection's api
	Router ziface.IRouter

	// Channel to inform that the connection has exited/stopped
	ExitBuffChan chan bool
}

// GetConnID implements ziface.IConnection.
func (c Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetTCPConnection implements ziface.IConnection.
func (c Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// RemoteAddr implements ziface.IConnection.
func (c Connection) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}

// Start implements ziface.IConnection.
func (c Connection) Start() {
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

// Stop implements ziface.IConnection.
func (c Connection) Stop() {
	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

func NewConntion(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")

	defer fmt.Printf(c.Conn.RemoteAddr().String(), " conn reader exit!! \n")

	defer c.Stop()

	for {

		buffer := make([]byte, 512)

		_, err := c.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Error while reading data from connect connection ::: ", err)
			c.ExitBuffChan <- true
			continue
		}

		request := Request{conn: c, data: buffer}

		go func(r ziface.IRequest) {
			c.Router.PreHandle(r)
			c.Router.Handle(r)
			c.Router.PostHandle(r)
		}(&request)
	}
}

var _ ziface.IConnection = Connection{}
