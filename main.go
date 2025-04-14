// Server.go

package main

import (
	"fmt"
	"gorik/ziface"
	"gorik/znet"
)

// ping test custom router
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	// Read client data first, then reply with ping...ping...ping
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
	// Read client data first, then reply with ping...ping...ping
	fmt.Println("recv from client: msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello Zinx Router V0.10"))
	if err != nil {
		fmt.Println(err)
	}
}

// Executed when a connection is created
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ... ")

	// Set two connection properties after the connection is created
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "🧗GOIRIKO🧗")
	conn.SetProperty("Home", "https://github.com/aceld/zinx")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// Executed when a connection is lost
func DoConnectionLost(conn ziface.IConnection) {
	// Before the connection is destroyed, query the "Name" and "Home" properties of the conn
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name =", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home =", home)
	}

	fmt.Println("DoConnectionLost is Called ... ")
}

func main() {
	// Create a server handle
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
