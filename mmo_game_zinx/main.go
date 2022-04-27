package main

import "zinx/zinx/znet"

func main() {
	//创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	//链接创建和销毁的HOOK函数

	//注册一些路由业务

	s.Serve()
}
