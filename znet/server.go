package znet

import (
	"errors"
	"fmt"
	"gorik/utils"
	"gorik/ziface"
	"net"
)

type Server struct {
	Name        string
	IPversion   string
	IP          string
	Port        int
	msgHandler  ziface.IMsgHandle
	ConnMgr     ziface.IConnManager
	OnConnStart func(conn ziface.IConnection)
	// Hook function to be called when a connection is about to be disconnected for this server
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// AddRouter implements ziface.Iserver.
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router success!")
}

// Serve implements ziface.Iserver.
func (s *Server) Serve() {
	s.Start()
	select {}
}

// Start implements ziface.Iserver.
func (s *Server) Start() {

	fmt.Printf("[START] Server name: %s, listener at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	fmt.Printf("[sTart]sevrer listen on port %d in addr %s .... \n", s.Port, s.IP)

	s.msgHandler.StartWorkerPool()
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

		if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
			conn.Close()
			continue
		}
		dealConn := NewConnection(s, conn, id, s.msgHandler)
		id++
		go dealConn.Start()

	}
}

// Stop implements ziface.Iserver.
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name", s.Name)

	// Stop or clean up other necessary connection information or other information
	s.ConnMgr.ClearConn()

}
func NewServer() ziface.IServer {
	// Initialize the global configuration file first
	utils.GlobalObject.Reload()

	s := &Server{
		Name:       utils.GlobalObject.Name, // Get from global parameters
		IPversion:  "tcp4",
		IP:         utils.GlobalObject.Host,    // Get from global parameters
		Port:       utils.GlobalObject.TcpPort, // Get from global parameters
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	println("[CONN handle] callback to client  ... .")
	if _, err := conn.Write(data[:cnt]); err != nil {

		println("write back to client err   ::::  ", err)
		return errors.New("Calback to lcient err ")

	}
	return nil
}
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// Set the hook function to be called when a connection is about to be disconnected for the server
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// Invoke the OnConnStart hook function for the connection
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// Invoke the OnConnStop hook function for the connection
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
