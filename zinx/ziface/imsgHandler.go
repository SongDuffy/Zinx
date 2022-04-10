package ziface

/*
	消息管理抽象层
*/

type IMsgHandle interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体处理逻辑
	AddRouter(msgID uint32, router IRouter)
	//启动worker工作池
	StartWorkerPool()
	//发送消息给消息队列
	SendMsgToTaskQueue(request IRequest)
}
