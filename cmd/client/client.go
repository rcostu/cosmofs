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
	list_dirs *bool = flag.Bool("l", false, "List directories")
	list_dir *string = flag.String("L", "", "List directory")
)

const (
	PORT string = "5453"
)

func debug (format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v)
	}
}

func main () {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("A server must be specified")
	}

	conn, err := net.Dial("tcp", flag.Arg(0) + ":" + PORT)
	defer conn.Close()

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	var decodedFiles map[string] []*cosmofs.File = make(map[string] []*cosmofs.File)
	configDec := gob.NewDecoder(conn)

	if *list_dirs {
		fmt.Printf("List directories\n")
		fmt.Fprintf(conn, "List Directories\n")

		var numDirs, numFiles int

		err = configDec.Decode(&numDirs)

		if err != nil {
			log.Fatal("Error decoding length config file: ", err)
		}

		debug("DECODED LENGTH VALUE: %v", numDirs)

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
}
