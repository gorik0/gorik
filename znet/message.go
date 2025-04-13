package znet

type Message struct {
	Id      uint32 // ID of the message
	DataLen uint32 // Length of the message
	Data    []byte //
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// Get the length of the message data segment
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// Get the message ID
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

// Get the message content
func (msg *Message) GetData() []byte {
	return msg.Data
}

// Set the length of the message data segment
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

// Set the message ID
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

// Set the message content
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
