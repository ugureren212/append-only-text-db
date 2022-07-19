package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/icza/backscanner"
	pb "github.com/net-reply-future-networks/k8s-golang-append-only-store/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var port = flag.Int("port", 3000, "The server port")

// server struct/type is used to implement the actual grpc service
// to register the server struct/type with grpc you need to embed pb.UnimplementedDatastoreServer within the struct/type
type server struct {
	pb.UnimplementedDatastoreServer
}

// handler to set key value pair into a single txt file
// path of txt file can be found in appendToFile()
func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	var err error

	err = appendToFile(in.Key, in.Value, false)
	if err != nil {
		msg := fmt.Sprintf("cannot set key value pair. error description %e", err)
		err = status.Error(codes.InvalidArgument, msg)
		return nil, err
	}

	log.Printf("Appended: %v %v to file /tmp/text.log", in.Key, in.Value)
	return &pb.SetReply{Message: "Status ok. Key Value set."}, nil
}

// gets key value pair from txt file specified in appendToFile()
// txt file is read bottom to top, instead of top to bottom, as the db is append only
// if a key is deleted (marked with a "#" tombstone) all previous key value pairs with the same key are ignored
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	var err error

	// Prevents reading from top but setting key/value pair at bottom
	key, value, err := readFileBackwards(in.Key)
	if err != nil {
		msg := fmt.Sprintf("cannot get key value pair. error description %e", err)
		err = status.Error(codes.InvalidArgument, msg)
		return nil, err
	}
	log.Printf("Received: %v %v", key, value)
	return &pb.GetReply{Key: key, Value: value}, nil
}

// deletes key value pairs by appending a # "tombstone" as a 3rd item
// This means that any key value pairs with the same key will not be retrieved
func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	var err error

	err = appendToFile(in.Key, "", true)
	if err != nil {
		msg := fmt.Sprintf("cannot delete key. error description %e", err)
		err = status.Error(codes.InvalidArgument, msg)
		return nil, err
	}

	log.Printf("Delete: %v from file /tmp/text.log", in.Key)
	return &pb.DelReply{Message: "Key deleted"}, nil
}

// reads /tmp/text.log bottom to top and returns key value pair from text file if it exists
// as this is a append only db, data needs to be read from the bottom first
// key value pairs starting with a # "tombstone" will not be retrieved as that key value pair and previous identical keys are considered to be deleted
func readFileBackwards(queryKey string) (string, string, error) {
	f, err := os.Open("/tmp/text.log")
	if err != nil {
		err2 := fmt.Errorf("could not access db file. error description : %v", err)
		return "", "", err2
	}
	fi, err := f.Stat()
	if err != nil {
		err2 := fmt.Errorf("could not access file info structure describing file. error description : %v", err)
		return "", "", err2
	}
	defer f.Close()

	scanner := backscanner.New(f, int(fi.Size()))
	targetKey := []byte(queryKey)
	for {
		line, pos, err := scanner.LineBytes()
		if err != nil {
			if err == io.EOF {
				err2 := fmt.Errorf("key: %v is not found in file", targetKey)
				return "", "", err2
			} else {
				err2 := fmt.Errorf("error searching for key in db file. error description : %v", err)
				return "", "", err2
			}
		}

		if bytes.Contains(line, targetKey) {

			log.Printf("Found key: %q at line position: %d, line: %q\n", targetKey, pos, line)

			items := bytes.Split([]byte(line), []byte(","))

			if string(items[0]) == "#" {
				strTargetKey := string(targetKey)
				err := fmt.Errorf("key '%v' does not exist. key has already been deleted", strTargetKey)
				return "", "", err
			}

			return string(items[0]), string(items[1]), nil
		}
	}
}

// appends key value pair to /tmp/text.log
// this function is also used to delete files by appending a # "tombstone" to mark the key value pair as deleted
func appendToFile(key string, value string, tombstone bool) error {
	var data string
	var err error

	if tombstone {
		data = "#" + "," + key + "," + value + "\n"
	} else {
		data = key + "," + value + "\n"
	}

	f, err := os.OpenFile("/tmp/text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("could not access db file. error description : %w", err)
	}

	// makes sure os writes data
	defer f.Sync()
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return fmt.Errorf("could not write to db file. error description : %w", err)
	}
	return err
}

// TODO: Unit test (readfilebackwards & append)
// TODO: Deploys k8s
// TODO: Improve Docker file
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDatastoreServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
