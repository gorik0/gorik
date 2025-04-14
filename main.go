package main

import (
	"fmt"
	"gorik/ziface"
	"gorik/znet"
)

// Ping test custom router
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	// Read the data from the client first, then write back ping...ping...ping
	fmt.Println("recv from client: msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	// Read the data from the client first, then write back Hello Zinx Router V0.8
	fmt.Println("recv from client: msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello Zinx Router V0.8"))
	if err != nil {
		fmt.Println(err)
	}
}

// Executed when a connection is created
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called...")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// Executed when a connection is lost
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called...")
}

func main() {
	// Create a server handler
	s := znet.NewServer()

	// Register connection hook callback functions
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// Configure routers
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// Start the server
	s.Serve()
}
