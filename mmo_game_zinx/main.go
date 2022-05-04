package main

import (
	"fmt"
	"zinx/mmo_game_zinx/apis"
	"zinx/mmo_game_zinx/core"
	"zinx/zinx/ziface"
	"zinx/zinx/znet"
)

//当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)

	//给客户端发送MsgID为1的消息  同步当前player的ID给客户端
	player.SyncPid()

	//给客户端发送MsgID为200的消息  同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()

	//将当前新上线的玩家添加到World manager中
	core.WorldMgrObj.AddPlayer(player)

	//将该链接绑定一个pid玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	//同步周边玩家，告知他们当前玩家已经上线，广播当前玩家信息位置
	player.SyncSurrounding()

	fmt.Println("=====> Player pid = ", player.Pid, " is arrived <=====")
}

//给当前链接断开之前触发的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	//通过链接属性得到当前链接所绑定的Pid
	pid, _ := conn.GetProperty("pid")

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//玩家下线
	player.Offline()
	fmt.Println("Player pid:", pid, "is offline")
}

func main() {
	//创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	//链接创建和销毁的HOOK函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动服务
	s.Serve()
}
