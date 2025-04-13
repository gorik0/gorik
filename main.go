package main

import (
	"fmt"
	"gorik/testik"
	"gorik/znet"
	"net"
	"time"
)

func main() {
	go ClientTest()

	s := znet.NewServer("Goriko server ;;; ")
	r := testik.PingRouter{}
	s.AddRouter(&r)
	s.Serve()
}
func ClientTest() {
	fmt.Println("Starting clietn !")
	time.Sleep(time.Second * 3)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		println("Client start err :::; ", err)
	}

	for {
		_, err := conn.Write([]byte("HELO server !!! "))
		if err != nil {
			println("error while writing ", err)
		}
		buffer := make([]byte, 512)
		cnt, err := conn.Read(buffer)
		if err != nil {
			println("error while reading .... ", err)
		}
		fmt.Printf("Succesfully reading .... msg:%s , cnt = %d \n", buffer, cnt)
		time.Sleep(time.Second * 1)
	}

}
