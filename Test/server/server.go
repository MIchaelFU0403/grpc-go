package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	// pb "github.com/EDDYCJY/go-grpc-example/proto"
	pb "github.com/MinH-09/grpc-go/proto"
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
}

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	return nil
}

func (s *StreamService) mustEmbedUnimplementedStreamServiceServer() {
	print()
}
