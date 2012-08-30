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
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
)

var (
	// Flags
	verbose *bool = flag.Bool("v", false, "Verbose output ON")
	myIP net.Addr
)

const (
	PORT int = 5453
	LOCALPORT int = 5454 
)

func debug (format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v)
	}
}

func listDirectories(conn *net.TCPConn) {
	dirs, err := cosmofs.Table.ListAllDirs()

	if err != nil {
		log.Printf("Error reading dirs %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(dirs)
}

func listKnownIDs(conn *net.TCPConn) {
	ids, err := cosmofs.Table.ListIDs()

	if err != nil {
		log.Printf("Error reading ids %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(ids)
}

func listConnectedIDs(conn *net.TCPConn) {
	encod := gob.NewEncoder(conn)

	encod.Encode(cosmofs.ConnectedPeers)
}

func listDirectoriesID(conn *net.TCPConn, reader *bufio.Reader) {
	id, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	id = strings.TrimRight(id, "\n")

	log.Printf("List directories for id %s from %s\n", id, conn.RemoteAddr())

	dirs, err := cosmofs.Table.ListDirs(id)

	if err != nil {
		log.Printf("Error reading dirs %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(dirs)
}

func listDirectory(conn *net.TCPConn, reader *bufio.Reader) {
	dirRecv, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	dirRecv = strings.TrimRight(dirRecv, "\n")

	id, dir, _ := cosmofs.SplitPath(dirRecv)

	log.Printf("List directory %s for id %s from %s\n", dir, id, conn.RemoteAddr())

	dirs, err := cosmofs.Table.ListDir(id, dir)

	if err != nil {
		log.Printf("Error reading dirs %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(dirs)
}

func search(conn *net.TCPConn, reader *bufio.Reader) {
	search, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	search = strings.TrimRight(search, "\n")

	log.Printf("Searching for %s from %s\n", search, conn.RemoteAddr())

	result, err := cosmofs.Table.Search(search)

	if err != nil {
		log.Printf("Error searching %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(result)
}

func searchDir(conn *net.TCPConn, reader *bufio.Reader) {
	search, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	search = strings.TrimRight(search, "\n")

	log.Printf("Searching Directories for %s from %s\n", search, conn.RemoteAddr())

	result, err := cosmofs.Table.SearchDir(search)

	if err != nil {
		log.Printf("Error searching directories %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(result)
}

func searchFile(conn *net.TCPConn, reader *bufio.Reader) {
	search, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	search = strings.TrimRight(search, "\n")

	log.Printf("Searching files for %s from %s\n", search, conn.RemoteAddr())

	result, err := cosmofs.Table.SearchFile(search)

	if err != nil {
		log.Printf("Error searching files %s", err)
	}

	encod := gob.NewEncoder(conn)

	encod.Encode(result)
}

func openFile(conn *net.TCPConn, reader *bufio.Reader) {
	file, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	file = strings.TrimRight(file, "\n")

	id, dirC, _ := cosmofs.SplitPath(file)

	fileName := filepath.Base(dirC)

	dir := strings.SplitN(dirC, "/", 2)

	log.Printf("Opening File %s in dir %s from %s\n", fileName, dir[0], conn.RemoteAddr())

	// Local file
	if strings.EqualFold(id, cosmofs.MyPublicPeer.ID) {
		files := cosmofs.Table[id][dir[0]]

		for _, v := range files {
			if strings.EqualFold(fileName, v.Filename) {
				encod := gob.NewEncoder(conn)
				debug("Encoding %v\n", filepath.Join(v.LocalPath, v.Filename))

				file, err := ioutil.ReadFile(filepath.Join(v.LocalPath, v.Filename))

				if err != nil {
					log.Printf("Error reading file %s\n", err)
					return
				}

				encod.Encode(file)
				break
			}
		}
	} else {	//Remote file
		if ip, ok := cosmofs.ConnectedPeers[id]; ok {
			connTCPS, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:		net.ParseIP(ip),
				Port:	PORT,
			})

			if err != nil {
				log.Fatalf("Error: %s\n", err)
				return
			}

			_, err = connTCPS.Write([]byte("Open File\n"))

			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}

			_, err = connTCPS.Write([]byte(file+"\n"))

			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}

			var content []byte

			decod := gob.NewDecoder(connTCPS)
			decod.Decode(&content)

			connTCPS.Close()

			encod := gob.NewEncoder(conn)

			encod.Encode(content)
		} else {
			log.Printf("Peer %v doesn't seem to be online\n", id)
		}
	}
}

func handleLocalPetition (conn *net.TCPConn) {
	defer conn.Close()

	debug("LOCAL PETITION")

	reader := bufio.NewReader(conn)

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
			debug("Table is now: %v\n", cosmofs.Table)
			listDirectories(conn)
		case "List Directories ID":
			debug("List Directories ID")
			listDirectoriesID(conn, reader)
		case "List Directory":
			debug("List Directory")
			listDirectory(conn, reader)
		case "List Known IDs":
			debug("List Known IDs from: %s\n", conn.RemoteAddr())
			listKnownIDs(conn)
		case "List Connected IDs":
			debug("List Connected IDs from: %s\n", conn.RemoteAddr())
			listConnectedIDs(conn)
		case "Search":
			debug("Search from %s\n", conn.RemoteAddr())
			search(conn, reader)
		case "Search Directory":
			debug("Search Directory from %s\n", conn.RemoteAddr())
			searchDir(conn, reader)
		case "Search File":
			debug("Search File from %s\n", conn.RemoteAddr())
			searchFile(conn, reader)
		case "Open File":
			debug("Open File from %s\n", conn.RemoteAddr())
			openFile(conn, reader)
	}
}

// Handles petitions from the peers.
func handleTCPPetition (lnTCP *net.TCPListener) {
	debug("WAITING FOR TCP CONN\n")

	conn, err := lnTCP.AcceptTCP()

	if err != nil {
		debug("Error: %s\n", err)
		go handleTCPPetition(lnTCP)
		conn.Close()
		return
	}

	remIP := strings.Split(conn.RemoteAddr().String(), ":")

	if strings.EqualFold(remIP[0], "127.0.0.1") {
		go handleTCPPetition(lnTCP)
		go handleLocalPetition(conn)
		return
	}

	defer conn.Close()

	debug("Connection made from: %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		debug("Error reading connection: %s", err)
		return
	}

	line = strings.TrimRight(line, "\n")

 	switch line {
		case "General TCP":
			debug("GENERAL TCP CONNECTION\n")
			connTCPS, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:		net.ParseIP(remIP[0]),
				Port:	PORT,
			})

			if err != nil {
				log.Fatalf("Error: %s\n", err)
				go handleTCPPetition(lnTCP)
				return
			}

			_, err = connTCPS.Write([]byte("General ANSWER\n"))

			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}

			encod := gob.NewEncoder(connTCPS)

			cosmofs.SendPeer(encod)

			// Send the number of shared directories
			err = encod.Encode(cosmofs.Table)

			if err != nil {
				log.Fatal("Error sending shared Table: ", err)
			}

			debug("List of Peers: %v\n", cosmofs.PeerList)

			decod := gob.NewDecoder(conn)

			id := cosmofs.ReceivePeer(decod)

			cosmofs.ConnectedPeer(id, remIP[0])

			log.Printf("CONNECTED: %v\n", cosmofs.ConnectedPeers)

			debug("List of Peers: %v\n", cosmofs.PeerList)

			cosmofs.Table.ReceiveAndMergeTable(decod)

			cosmofs.PrintTable()

			connTCPS.Close()

			go handleTCPPetition(lnTCP)

		case "General ANSWER":
			debug("GENERAL ANSWER\n")

			debug("List of Peers: %v\n", cosmofs.PeerList)

			decod := gob.NewDecoder(conn)

			id := cosmofs.ReceivePeer(decod)

			cosmofs.ConnectedPeer(id, remIP[0])

			log.Printf("CONNECTED: %v\n", cosmofs.ConnectedPeers)

			debug("List of Peers: %v\n", cosmofs.PeerList)

			cosmofs.Table.ReceiveAndMergeTable(decod)

			cosmofs.PrintTable()

			go handleTCPPetition(lnTCP)

		case "Open File":
			debug("OPEN FILE CONNECTION\n")
			file, err := reader.ReadString('\n')

			if err != nil && err != io.EOF {
				debug("Error reading connection: %s", err)
				return
			}

			file = strings.TrimRight(file, "\n")

			id, dirC, _ := cosmofs.SplitPath(file)

			fileName := filepath.Base(dirC)

			dir := strings.SplitN(dirC, "/", 2)

			log.Printf("Opening File %s in dir %s from %s\n", fileName, dir[0], conn.RemoteAddr())

			// Local file
			if strings.EqualFold(id, cosmofs.MyPublicPeer.ID) {
				files := cosmofs.Table[id][dir[0]]

				for _, v := range files {
					if strings.EqualFold(fileName, v.Filename) {
						encod := gob.NewEncoder(conn)
						debug("Encoding %v\n", filepath.Join(v.LocalPath, v.Filename))

						file, err := ioutil.ReadFile(filepath.Join(v.LocalPath, v.Filename))

						if err != nil {
							log.Printf("Error reading file %s\n", err)
							return
						}

						encod.Encode(file)
						break
					}
				}
			} else {
				log.Printf("Cannot find file %v\n", dirC)
			}
	}
}

func handleUDPPetition (lnUDP *net.UDPConn, ch chan int) {
	data := make([]byte, 4096)
	_, remoteIP, err := lnUDP.ReadFromUDP(data)

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	remIP := strings.Split(remoteIP.String(), ":")
	locIP := strings.Split(myIP.String(), ":")

	log.Printf("REM IP: %v, LOCAL IP: %v\n", remIP[0], locIP[0])

	cosmofs.ConnectedPeer(string(data), remIP[0])

	log.Printf("CONNECTED: %v\n", cosmofs.ConnectedPeers)

	if strings.EqualFold(remIP[0], locIP[0]) {
		ch <- 1
		return
	}

	_, remoteIP, err = lnUDP.ReadFromUDP(data)

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	log.Printf("FINAL IP: %v\n", net.ParseIP(remIP[0]))

	connTCPS, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:		net.ParseIP(remIP[0]),
		Port:	PORT,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
		return
	}

	_, err = connTCPS.Write([]byte("General TCP\n"))

	if err != nil {
		log.Fatalf("Error: %s\n", err)
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

	connTCPS.Close()

	debug("FINALIZING UDP CONN\n")

	ch <- 1
}

func main () {
	flag.Parse()

	// Leave the process listening for other peers
	lnUDP, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:		net.IPv4zero,
		Port:	PORT,
	})

	if err != nil {
		debug("Error: %s\n", err)
	}

	//Leave the process listening for other peers
	lnTCP, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:		net.IPv4zero,
		Port:	PORT,
	})

	if err != nil {
		debug("Error: %s\n", err)
		return
	}

	// Every server sends a broadcast message to anyone connected on the same
	// network.
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:		net.IPv4(255,255,255,255),
		Port:	PORT,
	})

	if err != nil {
		log.Fatalf("Error: %s\n", err)
		return
	}

	myIP = conn.LocalAddr()

	log.Printf("My IP: %v\n", myIP)

	_, err = conn.Write([]byte(cosmofs.MyPublicPeer.ID))

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	conn.Close()

	ch := make(chan int, 1)

	go handleTCPPetition(lnTCP)

	for {
		go handleUDPPetition(lnUDP, ch)
		<-ch
	}
}
