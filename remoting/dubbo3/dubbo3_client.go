package dubbo3

import (
	"fmt"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/protocol"
	"github.com/apache/dubbo-go/remoting"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"net"
	"time"
)

type TripleClient struct {
	conn net.Conn
	h2Controller H2Controller
	exchangeClient *remoting.ExchangeClient
	addr string
}
func NewTripleClient()*TripleClient {
	return &TripleClient{}
}

func (t*TripleClient) SetExchangeClient(eclient *remoting.ExchangeClient){
	t.exchangeClient = eclient
}

func (t*TripleClient)Connect(url *common.URL) error {
	logger.Warn("want to connect to url = ", url.Location)
	conn, err := net.Dial("tcp", url.Location)
	if err != nil {
		return err
	}
	t.addr = url.Location
	t.conn = conn
	h2Controller := NewH2Controller(t.conn, false, grpc.MethodDesc{})
	h2Controller.H2ShakeHand()
	return nil
}

func (t*TripleClient)Request(request *remoting.Request, timeout time.Duration, response *remoting.PendingResponse) error{
	//invocation, ok := request.Data.(invocation.RPCInvocation)
	//if !ok{
	//	return perrors.New("request is not invocation")
	//}
	invocationPtr := request.Data.(*protocol.Invocation)
	argValues := (*invocationPtr).Arguments()
	reqData, err := proto.Marshal(argValues[0].(proto.Message))

	if err != nil{
		panic("client request marshal not ok ")
	}
	fmt.Printf("getInvocation first arg = %+v\n", reqData)
	t.h2Controller.UnaryInvoke((*invocationPtr).MethodName(), t.conn.RemoteAddr().String(), reqData)
	// todo call remote and recv rsp
	return nil
}

func (t*TripleClient) Close() {

}

func (t*TripleClient) IsAvailable() bool {
	return true
}
