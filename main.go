package main

import (
	"flag"
	"time"

	"github.com/dwelch0/gRPC-practice/file_dropper/client"
	"github.com/dwelch0/gRPC-practice/file_dropper/server"
)

func main() {
	port := flag.Int("port", 50051, "The server port")
	workdir := flag.String("workdir", "/file-dropper", "The server port")
	flag.Parse()

	go server.Run(*port)

	time.Sleep(time.Duration(time.Second * 5))

	client.Run(*workdir)
}
