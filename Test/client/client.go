package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"

	// pb "github.com/EDDYCJY/go-grpc-example/proto"
	pb "github.com/MinH-09/grpc-go/proto"
)

const (
	PORT = "9002"
)

// multiReadCloser is MultiReadCloser instance.
type multiReadCloser struct {
	reader  io.Reader
	closers []io.Closer
}

// MultiReadCloser returns a ReadCloser that's the logical concatenation of
// the provided input readClosers.
func MultiReadCloser(readClosers ...io.ReadCloser) io.ReadCloser {
	readers := make([]io.Reader, len(readClosers))
	closers := make([]io.Closer, len(readClosers))
	for readCloser := range readClosers {
		closers[readCloser] = readClosers[readCloser]
		readers[readCloser] = readClosers[readCloser]
	}

	return &multiReadCloser{
		reader:  io.MultiReader(readers...),
		closers: closers,
	}
}

// Reader is the interface that wraps the basic Read method for multiReadCloser.
func (m *multiReadCloser) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

// Closer is the interface that wraps the basic Close method for multiReadCloser.
func (m *multiReadCloser) Close() error {
	var merr error
	for i := range m.closers {
		if err := m.closers[i].Close(); err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	return merr
}

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)
	/*
		err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 2018}})
		if err != nil {
			log.Fatalf("printLists.err: %v", err)
		}
	*/
	file, _ := os.Open("E:\\VirtualMachine\\ubuntu-22.10-live-server-amd64.iso")

	var readClosers []io.ReadCloser
	readClosers = append(readClosers, file)
	rc := MultiReadCloser(readClosers...)
	buffer := make([]byte, 4096)

	println(time.Now().Format("2006-01-02 15:04:05"))
	for {
		n, err := rc.Read(buffer)

		// s := string(buffer)
		// println(buffer)
		// println(s)
		if err != nil && err != io.EOF {
			log.Fatalf("printRecord.err: %v", err)
			return
		}
		if err == io.EOF || n == 0 {
			break
		}
		err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "n", Value: 2018, Buffer: buffer}})
		if err != nil {
			log.Fatalf("printRecord.err: %v", err)
		}
	}

	println(time.Now().Format("2006-01-02 15:04:05"))

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2018}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}
}

func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}

	/*
		for n := 0; n < 6; n++ {
			err := stream.Send(r)
			if err != nil {
				return err
			}
		}
	*/
	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	// log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)

	return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	return nil
}
