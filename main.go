package main

import (
	"fmt"
	"gorik/znet"
	"io"
	"net"
	"time"
)

// func main() {
// 	go ClientTest()

// 	s := znet.NewServer("Goriko server ;;; ")
// 	r := testik.PingRouter{}
// 	s.AddRouter(&r)
// 	s.Serve()
// }
// func ClientTest() {
// 	fmt.Println("Starting clietn !")
// 	time.Sleep(time.Second * 3)
// 	conn, err := net.Dial("tcp", "127.0.0.1:7777")
// 	if err != nil {
// 		println("Client start err :::; ", err)
// 	}

// 	for {
// 		_, err := conn.Write([]byte("HELO server !!! "))
// 		if err != nil {
// 			println("error while writing ", err)
// 		}
// 		buffer := make([]byte, 512)
// 		cnt, err := conn.Read(buffer)
// 		if err != nil {
// 			println("error while reading .... ", err)
// 		}
// 		fmt.Printf("Succesfully reading .... msg:%s , cnt = %d \n", buffer, cnt)
// 		time.Sleep(time.Second * 1)
// 	}

// }

// -------------------CLIENT
// -------------------CLIENT
// -------------------CLIENT
// -------------------CLIENT
// Test the datapack packaging and unpacking functionality
func main() {

	go Client()
	// Create a TCP server socket
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Server listen error:", err)
		return
	}

	// Create a goroutine to handle reading and parsing the data from the client goroutine, which is responsible for handling sticky packets
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Server accept error:", err)
		}

		// Handle client requests
		go func(conn net.Conn) {
			// Create a datapack object dp for packaging and unpacking
			dp := znet.NewDataPack()
			for {
				// 1. Read the head part from the stream
				headData := make([]byte, dp.GetHeadLen())
				// ReadFull fills the buffer completely
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("Read head error")
					break
				}

				// 2. Unpack the headData bytes into msgHead
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("Server unpack error:", err)
					return
				}

				// 3. Read the data bytes from the IO based on dataLen
				if msgHead.GetDataLen() > 0 {
					// msg has data, need to read data again
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("Server unpack data error:", err)
						return
					}

					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}
			}
		}(conn)
	}

}

func Client() {
	// Create a client goroutine to simulate the data for sticky packets and send them
	time.Sleep(time.Second * 3)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client dial error:", err)
		return
	}

	// 1. Create a DataPack object dp
	dp := znet.NewDataPack()

	// 2. Pack msg1
	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Client pack msg1 error:", err)
		return
	}

	// 3. Pack msg2
	msg2 := &znet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Client pack msg2 error:", err)
		return
	}

	// 4. Concatenate sendData1 and sendData2 to create a sticky packet
	sendData1 = append(sendData1, sendData2...)

	// 5. Write data to the server
	conn.Write(sendData1)

	// Block the client
	select {}
}
