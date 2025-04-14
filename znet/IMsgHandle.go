package znet

import (
	"fmt"
	"gorik/ziface"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter // Map to store the handler methods for each MsgID
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId =", request.GetMsgID(), "is not FOUND!")
		return
	}

	// Execute the corresponding handler methods
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// Add specific handling logic for a message
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1. Check if the current msg's API handler method already exists
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgId = " + strconv.Itoa(int(msgId)))
	}
	// 2. Add the binding relationship between msg and api
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId =", msgId)
}
