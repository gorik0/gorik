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
	TcpServer     ziface.Iserver // Current global Server object of Zinx
	Host          string         // Current server host IP
	TcpPort       int            // Current server host listening port
	Name          string         // Current server name
	Version       string         // Current Zinx version
	MaxPacketSize uint32         // Maximum size of data packet to be read
	MaxConn       int            // Maximum number of allowed connections on the current server host
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
		Name:          "ZinxServerApp",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	// Load some user-configurable parameters from the configuration file
	GlobalObject.Reload()
}
