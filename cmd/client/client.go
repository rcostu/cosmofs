/**

Copyright (C) 2012  Roberto Costumero Moreno <roberto@costumero.es>

This file is part of Cosmofs.

Cosmofs is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cosmofs is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cosmofs.  If not, see <http://www.gnu.org/licenses/>.

**/

package main

import (
	//ssh "code.google.com/p/go.crypto/ssh"
	"cosmofs"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
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

func parseKey() {
	keyFileName := filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")

	fi, err := os.Lstat(keyFileName)

	if err != nil {
		log.Fatal("Error: SSH Key file doesn't found.")
	}

	keyFile, err := os.Open(keyFileName)

	if err != nil {
		log.Fatal("Error: Cannot open SSH Key file.")
	}

	defer keyFile.Close()

	buffer := make([]byte, fi.Size())

	keyFile.Read(buffer)

	fmt.Printf("%s\n", buffer)

	//var agKey *ssh.AgentKey = &ssh.AgentKey{blob: buffer, Comment: ""}

	//fmt.Println(agKey.Key())
}

func main () {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("A server must be specified")
	}

	//parseKey()

	//conn, err := net.Dial("tcp", flag.Arg(0) + ":" + PORT)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:		net.IPv4(10,0,0,255),
		Port:	5453,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	fmt.Println(conn.RemoteAddr())

	defer conn.Close()

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
