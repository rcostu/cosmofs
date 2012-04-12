package main

import (
	"cosmofs"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	verbose *bool = flag.Bool("v", false, "Verbose mode")
	list *bool = flag.Bool("l", false, "List directories")
)

func debug (format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v)
	}
}

func main () {
	flag.Parse()

	conn, err := net.Dial("tcp", "localhost:5453")

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	if *list {
		fmt.Printf("List directories\n")
		fmt.Fprintf(conn, "List\n")

		configDec := gob.NewDecoder(conn)

		var numDirs, numFiles int

		err = configDec.Decode(&numDirs)

		if err != nil {
			log.Fatal("Error decoding length config file: ", err)
		}

		debug("DECODED LENGTH VALUE: %v", numDirs)

		var decodedFiles map[string] []*cosmofs.File = make(map[string] []*cosmofs.File)
		var dir string

		for numDirs > 0 {
			err = configDec.Decode(&numFiles)

			if err != nil {
				log.Fatal("Error decoding length config file: ", err)
			}

			debug("DECODED NUM FILES VALUE: %v", numFiles)

			err = configDec.Decode(&dir)

			if err != nil {
				log.Fatal("Error decoding length config file: ", err)
			}

			debug("DECODED DIR NAME VALUE: %v", dir)
			fmt.Println("d: " + dir)

			var decodedFile *cosmofs.File
			decodedFiles[dir] = make([]*cosmofs.File, numFiles)

			for i := 0; i < numFiles; i++ {
				decodedFile = new(cosmofs.File)
				err = configDec.Decode(decodedFile)
				if err != nil {
					log.Fatal("Error decoding list of files config file: ", err)
				}
				decodedFiles[dir][i] = decodedFile
				debug("DECODED VALUES: %v", decodedFiles[dir][i])
				fmt.Println("--- " + decodedFile.Filename)
			}

			numDirs--
		}
	}
	conn.Close()
}
