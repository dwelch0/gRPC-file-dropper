package main

import (
	"flag"
	"sync"
	"time"

	"github.com/dwelch0/gRPC-practice/file_dropper/client"
	"github.com/dwelch0/gRPC-practice/file_dropper/server"
)

func main() {
	port := flag.Int("port", 50051, "The server port")
	//workdir := flag.String("workdir", "/file-dropper", "The server port")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run(*port)
	}()
	wg.Wait()

	time.Sleep(time.Duration(time.Second * 5))

	client.Run()
}
