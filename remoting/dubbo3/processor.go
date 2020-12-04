package dubbo3

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	pb "triple-protocol/protocol"

	"github.com/golang/protobuf/proto"
)

const defaultRWBufferMaxLen = 4096

type processor struct {
	stream                *stream
	codec                 CodeC
	readWriteMaxBufferLen uint32 // useless
}

// protoc config参数增加,对codec进行选择
func newProcessor(s *stream) *processor {
	return &processor{
		stream:                s,
		codec:                 NewProtobufCodeC(),
		readWriteMaxBufferLen: defaultRWBufferMaxLen,
	}
}

func (p *processor) processUnaryRPC(buf bytes.Buffer, stream *stream) (*bytes.Buffer, error) {
	readBuf := buf.Bytes()
	header := readBuf[:5]
	length := binary.BigEndian.Uint32(header[1:])

	desc := func(v interface{}) error {
		if err := p.codec.Unmarshal(readBuf[5:5+length], v.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	// 执行函数
	reply, err := pb.Greeter_serviceDesc.Methods[0].Handler(server{}, context.Background(), desc, nil)
	if err != nil {
		return nil, err
	}

	replyData, err := proto.Marshal(reply.(proto.Message))
	if err != nil {
		return nil, err
	}
	rsp := make([]byte, 5+len(replyData))
	rsp[0] = byte(0)
	binary.BigEndian.PutUint32(rsp[1:], uint32(len(replyData)))
	copy(rsp[5:], replyData[:])
	return bytes.NewBuffer(rsp), nil
}

/// 先放着
// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("######### get server request name :" + in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
