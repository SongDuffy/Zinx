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

//创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=====> DoConnectionBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	fmt.Println("Set connection ...")
	conn.SetProperty("Name", "syh")
}

//链接断开之前的需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("=====> DoConnection is Called...")
	fmt.Println("conn ID = ", conn.GetConnID(), "is lost...")

	//获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v1.0]")

	//2 注册链接hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3 给当前zinx框架添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//4 启动server
	s.Serve()
}
