package pinger

import (
	"context"
	api "do/api/pinger"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
)

var _ api.PingerServer = (*service)(nil)

func RegisterService(s *grpc.Server) {
	api.RegisterPingerServer(s, &service{})
}

type service struct{}

func (s *service) Ping(_ context.Context, in *api.PingRequest) (*api.PingReply, error) {
	return &api.PingReply{
		Message: fmt.Sprintf("Pong %s", in.Message),
	}, nil
}

func (s *service) BidiStream(stream api.Pinger_BidiStreamServer) error {
	for {
		msg := fmt.Sprintf("server ping %d", time.Now().Unix())
		if err := stream.Send(&api.StreamReply{Message: msg}); err != nil {
			return fmt.Errorf("could not sent stream: %v", err)
		}
		rsp, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("client closed stream: %v", err)
			return nil
		}
		if err != nil {
			return fmt.Errorf("could not receive client stream: %v", err)
		}
		fmt.Printf("client stream: %v\n", rsp)
		time.Sleep(time.Second)
	}
}

func (s *service) ClientStream(stream api.Pinger_ClientStreamServer) error {
	for {
		rsp, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("client closed stream: %v", err)
			if err := stream.SendAndClose(&api.StreamReply{Message: "OK"}); err != nil {
				return fmt.Errorf("could not SendAndClose: %v", err)
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("could not receive client stream: %v", err)
		}
		fmt.Printf("client stream: %v\n", rsp)
	}
}

func (s *service) ServerStream(in *api.StreamRequest, stream api.Pinger_ServerStreamServer) error {
	fmt.Printf("got initial client request: %v\n", in)
	for {
		msg := fmt.Sprintf("server msg %d", time.Now().Unix())
		if err := stream.Send(&api.StreamReply{Message: msg}); err != nil {
			return fmt.Errorf("could not sent stream: %v", err)
		}
		time.Sleep(time.Second)
	}
}
