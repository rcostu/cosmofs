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
	//"bufio"
	"bytes"
	"cosmofs"
	"encoding/gob"
	"flag"
	//"io"
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
func handleTCPPetition (lnTCP *net.TCPListener, ch chan int) {
	debug("WAITING FOR TCP CONN\n")

	conn, err := lnTCP.AcceptTCP()

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	debug("Connection made from: %s\n", conn.RemoteAddr())

	defer conn.Close()

	debug("List of Peers: %v\n", cosmofs.PeerList)

	decod := gob.NewDecoder(conn)

	cosmofs.ReceivePeer(decod)

	debug("List of Peers: %v\n", cosmofs.PeerList)

	cosmofs.Table.ReceiveAndMergeTable(decod)

	debug("LISTA DE DIRECTORIOS: %v\n", cosmofs.Table)

	ch <- 1

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

func handleUDPPetition (lnUDP *net.UDPConn, ch chan int) {
	data := make([]byte, 4096)
	_, remoteIP, err := lnUDP.ReadFromUDP(data)

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	if !bytes.HasPrefix(data, []byte("CosmoFS conn")) {
		debug("Error in protocol")
		return
	}

	_, remoteIP, err = lnUDP.ReadFromUDP(data)

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	remIP := strings.Split(remoteIP.String(), ":")

	cosmofs.ConnectedPeer(string(data), remIP[0])

	log.Printf("CONNECTED: %v\n", cosmofs.ConnectedPeers)

	log.Printf("FINAL IP: %v\n", net.ParseIP(remIP[0]))

	connTCPS, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:		net.ParseIP(remIP[0]),
		Port:	5453,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
		return
	}

	debug("TCP DIAL DONE\n")

	encod := gob.NewEncoder(connTCPS)

	cosmofs.SendPeer(encod)

	debug("PEER SENT\n")

	// Send the number of shared directories
	err = encod.Encode(cosmofs.Table)

	if err != nil {
		log.Fatal("Error sending shared Table: ", err)
	}

	debug("FINALIZING UDP CONN\n")

	ch <- 1
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

	//Leave the process listening for other peers
	lnTCP, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:		net.IPv4zero,
		Port:	5453,
	})

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	// Every server sends a broadcast message to anyone connected on the same
	// network.
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

	ch := make(chan int, 1)

	for {
		go handleUDPPetition(lnUDP, ch)
		go handleTCPPetition(lnTCP, ch)
		<-ch
		<-ch
	}
}
