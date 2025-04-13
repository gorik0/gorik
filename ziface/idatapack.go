package ziface

/*
Packaging and unpacking data
Directly operates on the data stream of TCP connections, adding header information for transmitting data to handle TCP packet sticking problem.
*/
type IDataPack interface {
	GetHeadLen() uint32                // Method to get the length of the packet header
	Pack(msg IMessage) ([]byte, error) // Method to pack the message
	Unpack([]byte) (IMessage, error)   // Method to unpack the message
}
