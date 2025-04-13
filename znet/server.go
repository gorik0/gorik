package znet

import (
	"errors"
	"fmt"
	"gorik/ziface"
	"net"
)

type Server struct {
	Name      string
	IPversion string
	IP        string
	Port      int
}

// Serve implements ziface.Iserver.
func (s *Server) Serve() {
	s.Start()
	select {}
}

// Start implements ziface.Iserver.
func (s *Server) Start() {

	fmt.Printf("[sTart]sevrer listen on port %d in addr %s .... \n", s.Port, s.IP)

	addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("errror while listentingn ::: ", err.Error())
	}

	fmt.Println("RESOLVER :::: ", addr)
	listener, err := net.ListenTCP(s.IPversion, addr)
	if err != nil {
		fmt.Println("listen", s.IPversion, "err", err)
		return

	}
	fmt.Println("start Zinx server  ", s.Name, " succ, now listening...")

	// 3. Start the server network connection business.
	for {
		// 3.1. Block and wait for client connection requests.
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err ", err)
			continue

		}

		var id uint32
		id = 0
		dealConn := NewConntion(conn, id, CallbackToClient)
		id++
		dealConn.Start()

	}
}

// Stop implements ziface.Iserver.
func (s *Server) Stop() {
	fmt.Printf("[STOOP server %s \n]", s.Name)
	// todo

}

func NewServer(name string) ziface.Iserver {
	return &Server{
		Name:      name,
		IPversion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
}

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	println("[CONN handle] callback to client  ... .")
	if _, err := conn.Write(data[:cnt]); err != nil {

		println("write back to client err   ::::  ", err)
		return errors.New("Calback to lcient err ")

	}
	return nil
}
