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

// server implementation of server-side streaming rpc
func (s *dropperServiceServer) Watch(_ *fd.WatchRequest, stream fd.DropperService_WatchServer) error {
	// initial response after connection is always `OpSync`
	wr := &fd.WatchResponse{}
	wr.Op = &fd.WatchResponse_Sync{
		Sync: &fd.WatchResponse_OpSync{
			Items: s.fileItems,
		},
	}
	if err := stream.Send(wr); err != nil {
		return nil
	}
	log.Println("Connection initiated. OpSync type WatchResponse sent")
	time.Sleep(time.Duration(time.Second * 10))

	// begin loop of synthetic file events
	var wrType string
	for len(s.fileItems) > 0 {
		opInfo := s.generateChange()

		wr := &fd.WatchResponse{}

		switch opInfo.op {
		case sync:
			wr.Op = &fd.WatchResponse_Sync{
				Sync: &fd.WatchResponse_OpSync{
					Items: s.fileItems,
				},
			}
			wrType = "OpSync"
		case updated:
			wr.Op = &fd.WatchResponse_Update{
				Update: &fd.WatchResponse_OpUpdate{
					Item: &opInfo.file,
				},
			}
			wrType = "OpUpdate"
		case removed:
			wr.Op = &fd.WatchResponse_Remove{
				Remove: &fd.WatchResponse_OpRemove{
					Path: opInfo.path,
				},
			}
			wrType = "OpRemove"
			// no op
		default:
			continue
		}

		if err := stream.Send(wr); err != nil {
			return nil
		}
		log.Printf("%s type WatchResponse sent/n", wrType)

		time.Sleep(time.Duration(time.Second * 10))
	}
	return nil
}

// catch all for generateChange() return info
type fileOp struct {
	op   operation
	file fd.FileItem
	path string
}

type operation int

const (
	sync    operation = iota // 0
	updated operation = iota // 1
	removed operation = iota // 2
	noop    operation = iota // 3
)

// introduces continuous change on sample data
func (s *dropperServiceServer) generateChange() fileOp {
	rand.Seed(time.Now().UnixNano())

	// select an item tracked by server
	randItemIndex := rand.Intn(len(s.fileItems))

	switch rand.Intn(5) { // allow room for no-ops
	case 0:
		return fileOp{op: sync}
	case 1:
		// randomly updates a file's owner
		users := []string{"ravioli", "linguine", "penne"}
		randUser := users[rand.Intn(len(users))]
		if s.fileItems[randItemIndex].User == randUser {
			return fileOp{op: noop}
		} else {
			s.fileItems[randItemIndex].User = randUser
			return fileOp{op: updated, file: *s.fileItems[randItemIndex]}
		}
	case 2:
		tmpPath := s.fileItems[randItemIndex].Path
		// removes the randomly selected file item
		s.fileItems = append(s.fileItems[:randItemIndex], s.fileItems[randItemIndex+1:]...)
		return fileOp{op: removed, path: tmpPath}
	default:
		return fileOp{op: noop}
	}
}

func newServer() *dropperServiceServer {
	// real implementation would search a specified local directory
	// and populate FileItems based upon content discovered
	// this is a lazy approach to supply sample data
	s := &dropperServiceServer{fileItems: []*fd.FileItem{
		{
			Path:        "a.txt",
			Contents:    []byte("contents of the a.txt file"),
			Permissions: 777,
			User:        "foo",
			Group:       "bar",
			Command:     "head",
		},
		{
			Path:        "b.txt",
			Contents:    []byte("contents of the b.txt file"),
			Permissions: 777,
			User:        "foo",
			Group:       "bar",
			Command:     "cat",
		},
		{
			Path:        "c.txt",
			Contents:    []byte("contents of the c.txt file"),
			Permissions: 777,
			User:        "foo",
			Group:       "bar",
			Command:     "less",
		},
	}}
	return s
}

func Run(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fd.RegisterDropperServiceServer(grpcServer, newServer())
	log.Printf("file dropper server listening on localhost:%d", port)
	grpcServer.Serve(lis)
}
