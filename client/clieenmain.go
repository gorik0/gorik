package main

import (
	"fmt"
	"gorik/znet"
	"io"
	"net"
	"time"
)

func main() {
	Client()
}

func Client() {
	fmt.Println("Client Test... start")
	// Wait for 3 seconds to give the server a chance to start the service
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// Send packet message
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx V0.5 Client Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err", err)
			return
		}

		// Read the head part of the stream first
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) // ReadFull will fill msg until it's full
		if err != nil {
			fmt.Println("read head error")
			break
		}
		// Unpack the headData byte stream into msg
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// msg has data, need to read data again
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// Read byte stream from the io based on dataLen
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
