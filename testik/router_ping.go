package testik

import (
	"fmt"
	"gorik/ziface"
	"gorik/znet"
)

// Ping test custom router
type PingRouter struct {
	znet.BaseRouter // Must embed BaseRouter first
}

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")

	_, err := request.GetConenction().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("callback ping ping ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")

	_, err := request.GetConenction().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("callback ping ping ping error")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")

	_, err := request.GetConenction().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("callback ping ping ping error")
	}
}
