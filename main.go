package main

import (
	"flag"
	"time"

	"github.com/dwelch0/gRPC-practice/file_dropper/server"
)

func main() {
	port := flag.Int("port", 50051, "The server port")
	workdir := flag.String("workdir", "/file-dropper", "The server port")
	flag.Parse()

	go server.Run(*port, *workdir)
	time.Sleep(time.Duration(time.Second * 5))
}
