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

//玩家世界聊天消息
func (p *Player) Talk(content string) {
	//组件msgID为200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//向所有玩家发送msgID为200的消息
	for _, players := range players {
		players.SendMsg(200, proto_msg)
	}
}

//同步玩家
func (p *Player) SyncSurrounding() {
	//获取当前玩家周围玩家信息
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//将当前玩家位置信息通过msgID为200发送给周围玩家
	//组建msgID为200的proto数据
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

	//全部周围玩家都向各自的客户端发送200消息，proto_msg
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	//将周围玩家的位置信息发送给当前玩家
	//组建msgID为202的proto数据
	//制作pb.Player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个message Player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}

		players_proto_msg = append(players_proto_msg, p)
	}
	//封装SyncPlayer protobuf数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	//将组建好的数据发送给当前玩家的客户端
	p.SendMsg(202, SyncPlayers_proto_msg)
}

//更新玩家位置
func (p *Player) UpdatePosition(x, y, z, v float32) {
	//更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组建广播proto协议 MsgID:200 Tp-4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取当前玩家的周边玩家
	players := p.GetSurroundingPlayers()

	//一次给每个玩家对应的客户端发送当前玩家位置更新信息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI九宫格内所有玩家ID
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	//将所有pid对应的player防盗players切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

//玩家下线
func (p *Player) Offline() {
	//得到当前玩家周边九宫格内都有哪些玩家
	players := p.GetSurroundingPlayers()

	//给周围玩家广播MsgId：201消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	//将当前玩家从AOI删除
	WorldMgrObj.AoiMgr.RemoveFromGridbyPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
