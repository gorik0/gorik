package ziface

type Iserver interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32, router IRouter)
}
