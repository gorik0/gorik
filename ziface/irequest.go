package ziface

type IRequest interface {
	GetConenction() IConnection
	GetData() []byte
}
