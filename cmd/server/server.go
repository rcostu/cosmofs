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
	"bufio"
	"cosmofs"
	"encoding/gob"
	"flag"
	"io"
	"log"
	"net"
	"strings"
)

var (
	// Flags
	verbose *bool = flag.Bool("v", false, "Verbose output ON")
)

func debug (format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v)
	}
}

// Handles petitions from the peers.
func handlePetition (conn net.Conn) {
	debug("Connection made from: %s\n", conn.RemoteAddr())

	defer conn.Close()

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	if err != nil {
		debug("Error reading connection: %s", err)
		return
	}

	line = strings.TrimRight(line, "\n")

	if line != "CosmoFS conn" {
		debug("Error in protocol")
		return
	}

	line, err = reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	line = strings.TrimRight(line, "\n")

	// Listing directories
	switch line {
		case "List Directories":
			debug("List directories from: %s\n", conn.RemoteAddr())

			encod := gob.NewEncoder(conn)
			// Send the number of shared directories
			err = encod.Encode(cosmofs.Table)

			if err != nil {
				log.Fatal("Error sending shared Table: ", err)
			}
	}
}

func main () {
	flag.Parse()

	// Leave the process listening for other peers
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:		net.IPv4zero,
		Port:	5453,
	})

	if err != nil {
		debug("Error: %s\n", err)
	}

	debug("Listening on address ", ln.Addr())

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			debug("Error: %s\n", err)
			continue
		}
		go handlePetition(conn)
	}
}
