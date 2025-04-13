package ziface

// IMessaage is an abstract interface that encapsulates a message
type IMessage interface {
	GetDataLen() uint32 // Get the length of the message data segment
	GetMsgId() uint32   // Get the message ID
	GetData() []byte    // Get the message content

	SetMsgId(uint32)   // Set the message ID
	SetData([]byte)    // Set the message content
	SetDataLen(uint32) // Set the length of the message data segment
}
