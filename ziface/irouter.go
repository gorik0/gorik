package ziface

type IRouter interface {
	PreHandle(request IRequest)  // Hook method executed before handling the conn business
	Handle(request IRequest)     // Method to handle the conn business
	PostHandle(request IRequest) // Hook method executed after handling the conn business
}
