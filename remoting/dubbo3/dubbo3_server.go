package dubbo3

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/protocol"
	"github.com/apache/dubbo-go/protocol/invocation"
	"net"
	"sync"
)


type tripleServer struct {
	lst net.Listener
	handler func(*invocation.RPCInvocation) protocol.RPCResult
	addr string
}

func NewTripleServer(url *common.URL, handlers func(*invocation.RPCInvocation) protocol.RPCResult) *tripleServer {
	return &tripleServer{
		addr: url.Location,
		handler: handlers,
	}
}

// todo stop server
func (t*tripleServer) Stop(){

}

func (t*tripleServer) Start(){
	logger.Warn("tripleServer Start at ", t.addr)
	lst, err := net.Listen("tcp", t.addr)
	if err != nil {
		panic(err)
	}
	t.lst = lst
	t.Run()
}

func (t *tripleServer) Run() {
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

func (t *tripleServer) handleRawConn(conn net.Conn) {
	h2Controller := NewH2Controller(conn, true)
	h2Controller.H2ShakeHand()
	h2Controller.run()
}