package main

import (
	"flag"
	"log"
	"net"
)

var (
	verbose *bool = flag.Bool("v", false, "Verbose output ON")
	cosmofsin *string = flag.String("cosmofsin", os.Getenv("COSMOFSIN"), "Location of incoming packages")
	cosmofsout *string = flag.String("cosmofsout", os.Getenv("COSMOFSOUT"), "Location of shared directories")
)

func handlePetition (conn net.Conn) {
	if *verbose {
		log.Printf("Connection made from: %s\n", conn.RemoteAddr())
	}
}

func main () {
	flag.Parse()

	if *cosmofsin == "" {
		log.Fatalf("COSMOFSIN not set correctly. Current content <%s>", *cosmofsin)
	}

	if _, err := os.Lstat(*cosmofsin); err != nil {
		log.Fatalf("COSMOFSIN not set correctly. Current content <%s>", *cosmofsin)
	}

	ln, err := net.Listen("tcp", ":5453")

	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	if *verbose {
		log.Println("Listening on address ", ln.Addr())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error: %s\n", err)
			continue
		}
		go handlePetition(conn)
	}
}
