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

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:		net.IPv4(127,0,0,1),
		Port:	5453,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
		return
	}

	defer conn.Close()

	if *list_dirs {
		fmt.Printf("List directories\n")
		//fmt.Fprintf(conn, "List Directories\n")
		_, err = conn.Write([]byte("List Directories\n"))

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		decod := gob.NewDecoder(conn)

		var dirs []string

		decod.Decode(&dirs)

		for _, v := range dirs {
			fmt.Println(v)
		}
	}
}
