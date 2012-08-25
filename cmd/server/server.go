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
	"bytes"
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

/*func listDirectories() {
	fmt.Printf("List directories\n")
	
	_, err = conn.Write([]byte("List Directories\n"))

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	decod := gob.NewDecoder(conn)
	err = decod.Decode(&cosmofs.Table)

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
}*/

// Handles petitions from the peers.
func handlePetition (conn net.Conn) {
	debug("Connection made from: %s\n", conn.RemoteAddr())

	defer conn.Close()

	debug("List of Peers: %v\n", cosmofs.PeerList)

	cosmofs.ReceivePeer(conn)

	/*reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

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
	}*/
}

func main () {
	flag.Parse()

	// Leave the process listening for other peers
	lnUDP, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:		net.IPv4zero,
		Port:	5453,
	})

	if err != nil {
		debug("Error: %s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:		net.IPv4(255,255,255,255),
		Port:	5453,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
		return
	}

	log.Printf("My IP: %v\n", conn.LocalAddr())

	_, err = conn.Write([]byte("CosmoFS conn\n"))

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	_, err = conn.Write([]byte(cosmofs.MyPublicPeer.ID))

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	conn.Close()

	//Leave the process listening for other peers
	lnTCP, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:		net.IPv4zero,
		Port:	5453,
	})

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	for {
		data := make([]byte, 4096)
		_, remoteIP, err := lnUDP.ReadFromUDP(data)

		if err != nil {
			debug("Error: %s\n", err)
			continue
		}

		if !bytes.HasPrefix(data, []byte("CosmoFS conn")) {
			debug("Error in protocol")
			continue
		}

		_, remoteIP, err = lnUDP.ReadFromUDP(data)

		if err != nil {
			debug("Error: %s\n", err)
			continue
		}

		remIP := strings.Split(remoteIP.String(), ":")

		cosmofs.ConnectedPeer(string(data), remIP[0])

		log.Printf("CONNECTED: %v\n", cosmofs.ConnectedPeers)

		log.Printf("FINAL IP: %v\n", net.ParseIP(remIP[0]))

		connTCPS, err = net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:		net.ParseIP(remIP[0]),
			Port:	5453,
		})

		if err != nil {
			log.Fatalf("Error: %s\n", err)
			return
		}

		cosmofs.SendPeer(connTCPS)

		connTCPR, err := lnTCP.AcceptTCP()

		if err != nil {
			debug("Error: %s\n", err)
			continue
		}

		go handlePetition(connTCPR)
	}
}
