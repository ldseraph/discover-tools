package main

import (
	"fmt"
	"os"
)

func main() {
	eth, err := NewEth()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	server := NewServer(eth, ":1323")
	server.Run()
}
