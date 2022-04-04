package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

//封包、拆包的具体模块
type DataPack struct {
}

//拆包封包示例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节)+ID uint32(4字节)
	return 8
}

//封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将dataLen 写进databuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包(将包的Head读出来) 之后再根据head信息里data长度 再进行一次读取
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和MsgId
	msg := &Message{}

	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data receive")
	}

	return msg, nil
}
