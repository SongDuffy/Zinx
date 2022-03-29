package znet

import (
	"fmt"
	"net"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

/*
	连接模块
*/

type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前链接状态
	isClosed bool

	//告知当前链接已经退出的  停止channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID = ", c.Conn, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大的512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			return
		}
		//得到当前Conn数据 的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		//从路由中，找到注册绑定的Conn对应的Router调用
	}
}

func (c *Connection) Start() {
	fmt.Println("Connection Start, ConnID = ", c.ConnID)

	//启动从当前链接的读数据业务
	go c.StartReader()

	//启动从当前链接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop, ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	c.Conn.Close()
	close(c.ExitChan)
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

func (c *Connection) Send(data []byte) error {
	return nil
}
