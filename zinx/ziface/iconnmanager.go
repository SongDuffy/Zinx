package ziface

/*
	链接管理模块
*/

type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//根据connID获取连接
	Get(connID uint32) (IConnection, error)
	//得到当前连接总数
	Len() int
	//清除终止所有连接
	ClearConn()
}
