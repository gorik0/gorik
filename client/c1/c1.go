package main

import (
	"fmt"
	"gorik/znet"
	"io"
	"net"
	"time"
)

/*
Simulate client
*/
func main() {

	fmt.Println("Client Test ... start")
	// Wait for 3 seconds before sending the test request to give the server a chance to start
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// Pack the message
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx V0.6 Client0 Test Message")))
		println("writing msg .. ")
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		// Read the head part from the stream
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) // ReadFull fills the buffer until it's full
		println("ASGTER ")
		if err != nil {
			fmt.Println("read head error")
			break
		}

		// Unpack the headData into a message
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// The message has data, so we need to read the data part
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// Read the data bytes from the stream based on the dataLen
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
