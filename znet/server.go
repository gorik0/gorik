package znet

import (
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
		go func() {
			// Continuously loop to get data from the client.
			for {
				buf := make([]byte, 512)
				cnt, err := conn.Read(buf)
				if err != nil {
					fmt.Println("recv buf err ", err)
					continue

				}
				// Echo back the received data.
				if _, err := conn.Write(buf[:cnt]); err != nil {
					fmt.Println("write back buf err ", err)
					continue

				}

			}

		}()

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
