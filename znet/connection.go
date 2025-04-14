package znet

import (
	"errors"
	"fmt"
	"gorik/utils"
	"gorik/ziface"
	"io"
	"net"
)

type Connection struct {
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	// Current connection's ID, also known as SessionID, ID is   globally unique
	ConnID uint32
	// Current connection's close status
	isClosed bool

	// The handle function of this connection's api
	MsgHandler ziface.IMsgHandle
	// Channel to inform that the connection has exited/stopped
	ExitBuffChan chan bool
	msgChan      chan []byte
	msgBuffChan  chan []byte
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when sending message")
	}

	// Package the data and send it
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error, msgID =", msgID)
		return errors.New("Pack error message")
	}

	println("to chan")
	// Write back to the client
	// Change the previous direct write using conn.Write to sending the message to the Channel for the Writer to read
	c.msgChan <- msg
	println("afte sendong ")

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

func (c *Connection) Start() {
	// 1. Start the Goroutine for reading data from the client
	go c.StartReader()
	// 2. Start the Goroutine for writing data back to the client
	go c.StartWriter()

	// Call the registered hook method for connection creation according to the user's requirements
	c.TcpServer.CallOnConnStart(c)
}

// Stop implements ziface.IConnection.
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	// If the current connection is already closed
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// ==================
	// If the user registered a callback function for this connection's closure, it should be called explicitly at this moment
	c.TcpServer.CallOnConnStop(c)
	// ==================

	// Close the socket connection
	c.Conn.Close()
	// Close the writer
	c.ExitBuffChan <- true

	// Remove the connection from the connection manager
	c.TcpServer.GetConnMgr().Remove(c)

	// Close all channels of this connection
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server, // Set the server object
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}

	// Add the newly created connection to the connection manager
	c.TcpServer.GetConnMgr().Add(c)
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// Worker pool mechanism has been started, send the message to the Worker for processing
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// Execute the corresponding Handle method from the bound message and its corresponding processing method
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			// Data to be written to the client
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:", err, "Conn Writer exit")
				return
			}

		case data, ok := <-c.msgBuffChan:
			// Handling data for buffered channel
			if ok {
				// Data to be written to the client
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buffered Data error:", err, "Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}

		case <-c.ExitBuffChan:
			return
		}
	}
}
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when sending buffered message")
	}
	// Pack the data and send it
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID =", msgID)
		return errors.New("Pack error message")
	}

	// Write to the client
	c.msgBuffChan <- msg

	return nil
}
