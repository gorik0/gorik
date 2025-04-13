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

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	// Read the client's data first, then write back ping...ping...ping
	fmt.Println("recv from client: msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	// Write back data
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Create a server handle

	s := znet.NewServer()

	// Configure the router
	s.AddRouter(&PingRouter{})

	// Start the server
	s.Serve()
}
