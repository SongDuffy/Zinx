package znet

import "zinx/zinx/ziface"

//实现router时，先嵌入BaseRouter基类，然后重写
type BaseRouter struct {
}

//在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

//在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

//在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
