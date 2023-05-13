package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	fd "github.com/dwelch0/gRPC-practice/file_dropper"
	"google.golang.org/grpc"
)

type dropperServiceServer struct {
	fd.UnimplementedDropperServiceServer
	fileItems []*fd.FileItem
}

func (s *dropperServiceServer) Watch(req *fd.WatchRequest, stream fd.DropperService_WatchServer) error {
	for len(s.fileItems) > 0 {
		op := s.generateChange()

		wr := &fd.WatchResponse{}

		switch op {
		case updated:
		case removed:
		case noop:
			wr.Op = &fd.WatchResponse_Sync{
				Sync: &fd.WatchResponse_OpSync{
					Items: s.fileItems,
				},
			}
		}

		if err := stream.Send(wr); err != nil {
			return nil
		}

		time.Sleep(time.Duration(time.Second * 10))
	}
	return nil
}

type operation int

const (
	updated operation = iota // 0
	removed operation = iota // 1
	noop    operation = iota // 2
)

// lazy func to introduce continuous change on sample data
func (s *dropperServiceServer) generateChange() operation {
	rand.Seed(time.Now().UnixNano())

	// select an item tracked by server
	randItemIndex := rand.Intn(len(s.fileItems))

	switch rand.Intn(5) { // allow room for no-ops
	case 0:
		// randomly updates a file's owner
		users := []string{"ravioli", "linguine", "penne"}
		randUser := users[rand.Intn(len(users))]
		if s.fileItems[randItemIndex].User == randUser {
			return noop
		} else {
			s.fileItems[randItemIndex].User = randUser
			return updated
		}
	case 1:
		// removes the randomly selected file item
		s.fileItems = append(s.fileItems[:randItemIndex], s.fileItems[randItemIndex+1:]...)
		return removed
	default:
		return noop
	}
}

func newServer(workdir string) *dropperServiceServer {
	s := &dropperServiceServer{fileItems: []*fd.FileItem{
		{
			Path:        fmt.Sprintf("%s/a.txt", workdir),
			Contents:    []byte("contents of the a.txt file"),
			Permissions: 777,
			User:        "foo",
			Group:       "bar",
			Command:     fmt.Sprintf("echo %s/a.txt", workdir),
		},
		{
			Path:        fmt.Sprintf("%s/b.txt", workdir),
			Contents:    []byte("contents of the b.txt file"),
			Permissions: 777,
			User:        "foo",
			Group:       "bar",
			Command:     fmt.Sprintf("echo %s/b.txt", workdir),
		},
		{
			Path:        fmt.Sprintf("%s/c.txt", workdir),
			Contents:    []byte("contents of the c.txt file"),
			Permissions: 777,
			User:        "foo",
			Group:       "bar",
			Command:     fmt.Sprintf("echo %s/c.txt", workdir),
		},
	}}
	return s
}

func Run(port int, workdir string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fd.RegisterDropperServiceServer(grpcServer, newServer(workdir))
	grpcServer.Serve(lis)
}
