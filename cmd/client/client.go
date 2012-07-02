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
		log.Fatal("Error: Cannot find SSH Key file.")
	}

	keyFile, err := os.Open(keyFileName)

	if err != nil {
		log.Fatal("Error: Cannot open SSH Key file.")
	}

	defer keyFile.Close()

	buffer := make([]byte, fi.Size())

	keyFile.Read(buffer)

	fmt.Printf("%s\n", buffer)

	/*out, rest, ok := ParseString(buffer)

	if !ok {
		fmt.Println("Error")
	}

	fmt.Println("OUT: ", out)
	fmt.Println("REST: ", rest)

	os.Exit(1)*/
}

func main () {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("A server must be specified")
	}

	//parseKey()

	//conn, err := net.Dial("tcp", flag.Arg(0) + ":" + PORT)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:		net.IPv4(127,0,0,1),
		Port:	5453,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	defer conn.Close()

	_, err = conn.Write([]byte("CosmoFS conn\n"))

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	if *list_dirs {
		fmt.Printf("List directories\n")
		//fmt.Fprintf(conn, "List Directories\n")
		_, err = conn.Write([]byte("List Directories\n"))

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		configDec := gob.NewDecoder(conn)
		err = configDec.Decode(&cosmofs.Table)

		if err != nil {
			log.Fatal("Error decoding table: ", err)
		}

		dirs, err := cosmofs.Table.ListAllDirs()

		if err != nil {
			log.Printf("Error reading dirs %s", err)
		}

		for _, v := range dirs {
			fmt.Println(v)
		}
	}
}
