package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
	pb "zinx/mmo_game_zinx/pb"
	"zinx/zinx/ziface"
)

//玩家对象
type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的链接（用于和客户端的链接）
	X    float32            //平面的X坐标
	Y    float32            //高度
	Z    float32            //平面y坐标
	V    float32            //旋转的角度（0-360度）
}

/*
	Player ID 生成器（可换数据库）
*/
var PidGen int32 = 1  //用来生产玩家ID的计数器
var IdLock sync.Mutex //保护pidgen的mutex

//创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), //随机在140坐标点 基于Y轴若干偏移
		V:    0,
	}

	return p
}

/*
	提供一个发送给客户端消息的方法
	将protobuf数据序列化后，再调用sendmsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将proto Message结构体序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}

	//将二进制文件通过zinx框架的sendmsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player Send Msg error!")
		return
	}

	return
}

//告知客户端玩家Pid， 同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组建MsgID：0 的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

//广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//将消息发送给客户端
	p.SendMsg(200, proto_msg)
}
