package client

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"sync"

	fd "github.com/dwelch0/gRPC-practice/file_dropper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// processWatch performs rpc call to Watch (server-side streaming) and
// performs actions to local filesystem objects
func processWatch(client fd.DropperServiceClient) error {
	log.Println("Processing server-side stream from Watch()")

	stream, err := client.Watch(context.Background(),
		&fd.WatchRequest{},
	)
	if err != nil {
		return err
	}

	fo := newFileOrchestrator()
	// process stream
	for {
		wr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// determine op and execute concurrently
		switch op := wr.Op.(type) {
		case *fd.WatchResponse_Sync:
			log.Printf("Received sync operation: %v", op.Sync)
			go fo.synchronize(op.Sync.Items)
		case *fd.WatchResponse_Update:
			log.Printf("Received update operation: %v", op.Update)
			go fo.update(op.Update.Item)
		case *fd.WatchResponse_Remove:
			log.Printf("Received remove operation: %v", op.Remove)
			go fo.remove(op.Remove.Path)
		default:
			log.Printf("Received unknown operation: %v", op)
		}
	}
	return nil
}

type fileOpOrchestrator struct {
	fileLocks     map[string]*sync.Mutex // keys are file paths
	fileLocksLock sync.Mutex
}

func newFileOrchestrator() *fileOpOrchestrator {
	return &fileOpOrchestrator{
		fileLocks: make(map[string]*sync.Mutex),
	}
}

// reconciles local files with items slice sent by server
func (fo *fileOpOrchestrator) synchronize(items []*fd.FileItem) error {
	for _, item := range items {
		// acquire or create file-specific lock
		fo.fileLocksLock.Lock()
		fileLock, exists := fo.fileLocks[item.Path]
		if !exists {
			fileLock = &sync.Mutex{}
			fo.fileLocks[item.Path] = fileLock
		}
		fo.fileLocksLock.Unlock()

		fileLock.Lock()
		defer fileLock.Unlock()

		// create file
	}
	return nil
}

// updates file at specific path
// note: made design choice that server implementation does not consider alteration of file
// path to be a valid update operation. A remove and sync should be sent if file path is changed
// this choice affects logic around fileLocks map
func (fo *fileOpOrchestrator) update(item *fd.FileItem) error {
	// acquire file-specific lock
	fo.fileLocksLock.Lock()
	fileLock, exists := fo.fileLocks[item.Path]
	if !exists {
		fo.fileLocksLock.Unlock()
		return errors.New("file path does not exist in file lock map")
	}
	fo.fileLocksLock.Unlock()

	// lock operations on specific file
	fileLock.Lock()
	defer fileLock.Unlock()

	// update file

	return nil
}

// removes file at path from filesystem
func (fo *fileOpOrchestrator) remove(path string) error {
	// acquire file-specific lock
	fo.fileLocksLock.Lock()
	fileLock, exists := fo.fileLocks[path]
	if !exists {
		fo.fileLocksLock.Unlock()
		return errors.New("file path does not exist in file lock map")
	}
	fo.fileLocksLock.Unlock()

	// lock operations on specific file
	fileLock.Lock()
	defer fileLock.Unlock()

	// remove file

	// remove specific file lock from orchestrator
	fo.fileLocksLock.Lock()
	defer fo.fileLocksLock.Unlock()
	delete(fo.fileLocks, path)

	return nil
}

func Run() error {
	serverAddr := flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	flag.Parse()

	// configure/init gRPC client
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := fd.NewDropperServiceClient(conn)

	if err := processWatch(client); err != nil {
		return err
	}

	return nil
}
