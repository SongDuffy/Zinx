package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/zinx/znet"
)

/*
	模拟一个客户端
*/
func main() {
	fmt.Println("client0 start...")

	time.Sleep(1 * time.Second)
	//1 直接连接远程服务器，得到一个connect连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	//2 链接调用Write 写数据
	for {
		//发送封包的message消息 MsgId:0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.9 client test ping")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器回复message数据， MsgId：1
		//先读取流中的head部分 得到ID和datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error:", err)
			break
		}
		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			//再根据流中的head部分， 将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("---> Receive Server Msg : ID = ", msg.Id, "len = ", msg.DataLen,
				"data = ", string(msg.Data))
		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
