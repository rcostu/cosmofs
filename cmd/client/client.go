package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	verbose *bool = flag.Bool("v", false, "Verbose mode")
	list *bool = flag.Bool("l", false, "List directories")
)

func main () {
	conn, err := net.Dial("tcp", "localhost:5453")

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	if *list {
		fmt.Fprintf(conn, "List\n")
	}
	conn.Close()
}
