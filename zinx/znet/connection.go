package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

/*
	连接模块
*/

type Connection struct {
	//当前Conn隶属于哪个server
	TcpServer ziface.IServer

	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前链接状态
	isClosed bool

	//告知当前链接已经退出的  停止channel
	ExitChan chan bool

	//无缓冲管道，用于读写goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理msgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

//初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandler,
		isClosed:  false,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID = ", c.Conn, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大的512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err:", err)
		//	return
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的msgHead 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		//拆包， 得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据dataLen   再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}

		msg.SetData(data)

		//得到当前Conn数据 的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到注册绑定的Conn对应的Router调用
			//根据绑定好的MsgID 找到对应处理api业务执行
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

//StartWriter 写消息goroutine,专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Start write...]")
	defer fmt.Println(c.RemoteAddr().String(), "conn writer exit!")

	//不断阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println(err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Connection Start, ConnID = ", c.ConnID)

	//启动从当前链接的读数据业务
	go c.StartReader()

	//启动从当前链接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop, ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//调用开发者注册的 销毁链接之前 需要执行的业务hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//关闭write
	c.ExitChan <- true

	//将当前链接从ConnMgr中摘除
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed ")
	}

	//将data进行封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个链接属性
	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("No property found ")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property, key)
}
