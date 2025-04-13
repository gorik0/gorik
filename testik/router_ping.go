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

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")

	_, err := request.GetConenction().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("callback ping ping ping error")
	}
}
