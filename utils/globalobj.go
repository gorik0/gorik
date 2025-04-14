package utils

import (
	"encoding/json"
	"gorik/ziface"
	"io/ioutil"
)

/*
Stores all global parameters related to the Zinx framework for use by other modules.
Some parameters can also be configured by users in zinx.json.
*/
type GlobalObj struct {
	TcpServer ziface.IServer // Current global Server object of Zinx
	Host      string         // IP of the current server host
	TcpPort   int            // Listening port number of the current server host
	Name      string         // Name of the current server

	/*
	   Zinx
	*/
	Version          string // Current Zinx version number
	MaxPacketSize    uint32 // Maximum size of data packet
	MaxConn          int    // Maximum number of connections allowed on the current server host
	WorkerPoolSize   uint32 // Number of workers in the business worker pool
	MaxWorkerTaskLen uint32 // Maximum number of tasks stored in the task queue corresponding to each business worker
	MaxMsgChanLen    int
	/*
	   config file path
	*/
	ConfFilePath string
}

/*
Define a global object
*/
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		println("PANIC!!!!!! no cfg")
	}
	// Parse the JSON data into the struct
	// fmt.Printf("json: %s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
Provide the init() method, which is automatically loaded.
*/
func init() {
	// Initialize the GlobalObject variable and set some default values
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    10,
	}

	// Load some user-configured parameters from the configuration file
	GlobalObject.Reload()
}
