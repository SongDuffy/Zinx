package znet

import (
	"fmt"
	"strconv"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

/*
	消息处理模块实现
*/

type MsgHandle struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
	}
}

//调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! Need Register!")
		return
	}

	//根据msgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//ID存在
		panic("repeat api , msgID = " + strconv.Itoa(int(msgID)))
	}
	//添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " successful!")
}

//启动一个worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerpoolsize 分别开启worker,每个worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//启动一个worker
		//当前的worker对应channel开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTasklen)
		//启动当前worker 阻塞等待消息
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQuere chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is start...")

	//不断阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQuere:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue, 由work进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//平均分配消息给不同worker
	//根据客户端建立的ConnID进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request msgID = ", request.GetMsgID(), "workerID = ", workerID)
	//将消息发送给对应worker的task queue
	mh.TaskQueue[workerID] <- request
}
