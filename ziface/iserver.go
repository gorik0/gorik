package ziface

type Iserver interface {
	Start()
	Stop()
	Serve()
	AddRouter(router IRouter)
}
