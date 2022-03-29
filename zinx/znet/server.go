package znet

import (
	"fmt"
	"net"
	"zinx/zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前的server添加一个router，server注册的链接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port:%d, is starting\n", s.IP, s.Port)

	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("Resolve tcp addr fail, %s", err)
			return
		}
		//2 Listen Tcp Addr
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("Listen fail, %s", err)
			return
		}

		fmt.Printf("Start zinx server successful, %s successful,Listening...", s.Name)
		var cid uint32
		cid = 0

		//3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err,", err)
				continue
			}

			//已经与客户端建立连接，做一些业务
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动当前业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止和回收

}

func (s *Server) Serve() {
	//启动服务器功能
	s.Start()

	//做一些启动服务器之后的额外业务

	//阻塞
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!!")
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
