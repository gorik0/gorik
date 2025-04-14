package znet

import (
	"errors"
	"fmt"
	"gorik/ziface"
	"io"
	"net"
)

type Connection struct {
	Conn *net.TCPConn
	// Current connection's ID, also known as SessionID, ID is   globally unique
	ConnID uint32
	// Current connection's close status
	isClosed bool

	// The handle function of this connection's api
	MsgHandler ziface.IMsgHandle
	// Channel to inform that the connection has exited/stopped
	ExitBuffChan chan bool
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// Package the data and send it
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id =", msgId)
		return errors.New("Pack error msg")
	}

	// Write back to the client
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id", msgId, "error")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}

	return nil
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
	return c.Conn.RemoteAddr()
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

func NewConntion(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc, handler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   handler,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	println("pr datad ")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		// Create a data packing/unpacking object
		dp := NewDataPack()
		// Read the client's message header
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			c.ExitBuffChan <- true
			continue
		}
		println("read header ...")
		// Unpack the message, obtain msgid and datalen, and store them in msg
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			c.ExitBuffChan <- true
			continue
		}
		println("unpack to msg,,,")
		// Read the data based on dataLen and store it in msg.Data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		println("reading data ...")
		msg.SetData(data)

		// Get the Request data of the current client request
		req := Request{
			conn: c,
			msg:  msg, // Replace buf with msg
		}

		// Find the corresponding Handle registered in Routers based on the bound Conn
		go c.MsgHandler.DoMsgHandler(&req)
	}
}
