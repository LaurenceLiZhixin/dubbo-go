package dubbo3

import (
	"fmt"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol"
	"google.golang.org/grpc"
	"net"
	"sync"
)


type TripleServer struct {
	lst net.Listener
	addr string
	rpcService common.RPCService
	methodsDesc []grpc.MethodDesc
}

func NewTripleServer(url *common.URL) *TripleServer {
	return &TripleServer{
		addr: url.Location,
	}
}

// todo stop server
func (t*TripleServer) Stop(){

}

// DubboGrpcService is gRPC service
type DubboGrpcService interface {
	// SetProxyImpl sets proxy.
	SetProxyImpl(impl protocol.Invoker)
	// GetProxyImpl gets proxy.
	GetProxyImpl() protocol.Invoker
	// ServiceDesc gets an RPC service's specification.
	ServiceDesc() *grpc.ServiceDesc
}

func (t*TripleServer) Start(url *common.URL){
	key := url.GetParam(constant.BEAN_NAME_KEY, "")
	fmt.Println("key = ", key)
	service := config.GetProviderService(key)
	ds, ok := service.(DubboGrpcService)
	if !ok{
		panic("TripleServer Start failed, service not impl DubboGrpcService")
	}
	t.methodsDesc = ds.ServiceDesc().Methods
	logger.Warn("tripleServer Start at ", t.addr)
	lst, err := net.Listen("tcp", t.addr)
	if err != nil {
		panic(err)
	}
	t.lst = lst
	t.Run()
}

func (t *TripleServer) Run() {
	wg := sync.WaitGroup{}
	for {
		conn, err := t.lst.Accept()
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go func() {
			t.handleRawConn(conn)
			wg.Done()
		}()
	}
}

func (t *TripleServer) handleRawConn(conn net.Conn) {
	h2Controller := NewH2Controller(conn, true, t.methodsDesc[0])
	h2Controller.H2ShakeHand()
	h2Controller.run()
}