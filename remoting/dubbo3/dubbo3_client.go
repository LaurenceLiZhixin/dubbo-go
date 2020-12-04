package dubbo3

import (
	"fmt"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/protocol/invocation"
	"github.com/apache/dubbo-go/remoting"
	perrors "github.com/pkg/errors"
	"net"
	"time"
)

type TripleClient struct {
	conn net.Conn
	h2Controller H2Controller
	exchangeClient *remoting.ExchangeClient
	addr string
}

func(t*TripleClient) SetExchangeClient(client *remoting.ExchangeClient){
	t.exchangeClient= client
}

func (t*TripleClient)Connect(url *common.URL) error {
	logger.Warn("want to connect to url = ", url.Location)
	conn, err := net.Dial("tcp", url.Location)
	if err != nil {
		return err
	}
	t.addr = url.Location
	t.conn = conn
	h2Controller := NewH2Controller(t.conn, false)
	h2Controller.H2ShakeHand()
	return nil
}

func (t*TripleClient)Request(request *remoting.Request, timeout time.Duration, response *remoting.PendingResponse) error {
	invocation, ok := request.Data.(invocation.RPCInvocation)
	if !ok{
		return perrors.New("request is not invocation")
	}
	fmt.Printf("getInvocation = %+v\n", invocation)
	t.h2Controller.UnaryInvoke(invocation.MethodName(), t.conn.RemoteAddr().String(), request)
	// todo call remote and recv rsp
	return nil
}

func (t*TripleClient) Close() {

}

func (t*TripleClient) IsAvailable() bool {
	return true
}

func NewTripleClient() *TripleClient {
	return &TripleClient{}
}