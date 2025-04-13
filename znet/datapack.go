package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gorik/utils"
	"gorik/ziface"
)

// DataPack class for packaging and unpacking, no need for members currently
type DataPack struct{}

// Initialization method for the DataPack class
func NewDataPack() *DataPack {
	return &DataPack{}
}
func (dp *DataPack) GetHeadLen() uint32 {
	// Id uint32(4 bytes) + DataLen uint32(4 bytes)
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// Create a buffer to store the bytes
	dataBuff := bytes.NewBuffer([]byte{})

	// Write DataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// Write MsgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// Write data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// Create an io.Reader from the input binary data
	dataBuff := bytes.NewReader(binaryData)

	// Only extract the information from the header, obtaining dataLen and msgID
	msg := &Message{}

	// Read dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// Read msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// Check if the dataLen exceeds the maximum allowed packet size
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data received")
	}

	// We only need to unpack the header data, and then read the data once more from the connection based on the length of the header
	return msg, nil
}
