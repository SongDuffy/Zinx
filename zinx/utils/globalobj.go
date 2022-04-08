package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/zinx/ziface"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数可以通过zinx.json用户进行配置
*/

type GlobalObj struct {
	TcpServer ziface.IServer //当前Zinx全局server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器名称

	Version        string //当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值

}

//全局对外Globalobj
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("myDemo/ZinxV0.6/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json加载
	json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//初始化当前globalobj
func init() {
	//配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.6",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//从conf/zinx.json加载用户自定义参数
	GlobalObject.Reload()
}
