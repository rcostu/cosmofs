package main

import (
	"fmt"
	"net"
	"os"
)

func main () {
	conn, err := net.Dial("tcp", "localhost:5453")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}

	conn.Close()
}
