package main

import (
	"fmt"
	"zinx/zinx/ziface"
	"zinx/zinx/znet"
)

/*
	基于Zinx框架来开发的服务器端应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test Handle
func (PR *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端数据，再回写ping
	fmt.Println("receive from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

//hellozinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

//Test Handle
func (PR *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	//先读取客户端数据，再回写ping
	fmt.Println("receive from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.6]")
	//2 给当前zinx框架添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//3 启动server
	s.Serve()
}
