package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/mmo_game_zinx/core"
	pb "zinx/mmo_game_zinx/pb"
	"zinx/zinx/ziface"
	"zinx/zinx/znet"
)

//玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	//解析客户端传送过来的proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move : Position unmarshal error ", err)
		return
	}

	//得到当前发送位置的是哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error,", err)
		return
	}

	fmt.Printf("Pllayer pid = %d, move(%f,%f,%f,%f)\n", pid,
		proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	//给其他玩家进行当前玩家的位置广播
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePosition(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
