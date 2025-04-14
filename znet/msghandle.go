package znet

import (
	"fmt"
	"gorik/utils"
	"gorik/ziface"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter // Map to store the handler methods for each MsgID
	WorkerPoolSize uint32                    // Number of workers in the business worker pool
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("👨‍🌾👨‍🌾👨‍🌾Worker ID =", workerID, "is started.")
	// Continuously wait for messages in the queue
	for {
		select {
		// If there is a message, take the Request from the queue and execute the bound business method
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	// Start the required number of workers
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// A worker is started
		// Allocate space for the current worker's task queue
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// Start the current worker, blocking and waiting for messages in the corresponding task queue
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// Assign the current connection to the worker responsible for processing this connection based on ConnID
	// Round-robin average allocation policy
	// Get the workerID responsible for processing this connection
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	println("MSG go to ------->>>>>> ")
	println("MSG go to ------->>>>>> ")
	println("MSG go to ------->>>>>> ")
	println("MSG go to ------->>>>>> ")
	println("MSG go to ------->>>>>> ", workerID)

	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)

	// Send the request message to the task queue
	mh.TaskQueue[workerID] <- request
}
