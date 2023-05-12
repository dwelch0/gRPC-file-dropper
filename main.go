package main

import (
	"fmt"

	ds "github.com/dwelch0/gRPC-practice/dropper_service"
)

func main() {
	x := ds.NewDropperServiceClient
	fmt.Println(x)
}
