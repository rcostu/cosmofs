package main

import (
	"bufio"
	"cosmofs"
	"flag"
	"encoding/gob"
	"log"
	"net"
	"os"
	"path/filepath"
)

var (
	// Flags
	verbose *bool = flag.Bool("v", false, "Verbose output ON")
	cosmofsin *string = flag.String("cosmofsin", os.Getenv("COSMOFSIN"), "Location of incoming packages")
	cosmofsout *string = flag.String("cosmofsout", os.Getenv("COSMOFSOUT"), "Location of shared directories")

	// Shared Directory List
	sharedDirList []string = make([]string, len(filepath.SplitList(*cosmofsout)))
)

const (
	COSMOFSDIR string = ".cosmofs"
	COSMOFSCONFIGFILE string = ".cosmofsconfig"
)

// Handles petitions from the peers.
func handlePetition (conn net.Conn) {
	if *verbose {
		log.Printf("Connection made from: %s\n", conn.RemoteAddr())
	}


	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Error reading connection: %s", err)
	}

	log.Println(line)
}

func main () {
	flag.Parse()

	// Check if COSMOFSIN environment is set
	if *cosmofsin == "" {
		log.Fatalf("COSMOFSIN not set correctly. Current content <%s>", *cosmofsin)
	}

	// Check if COSMOFSIN is a correct directory
	if _, err := os.Lstat(*cosmofsin); err != nil {
		log.Fatalf("COSMOFSIN not set correctly. Current content <%s>", *cosmofsin)
	}

	if *verbose {
		log.Printf("Inbound files arriving to %s", *cosmofsin)
	}

	// Initialize sharedDirList
	sharedDirList = filepath.SplitList(*cosmofsout)

	// There shall be at least one shared directory
	if len(sharedDirList) == 0 {
		log.Fatal("COSMOFSOUT should have at least one directory or file.")
	}

	for _, dir := range sharedDirList {
		dir = filepath.Clean(dir)
		log.Println(dir)

		fi, err := os.Lstat(dir);

		if err != nil {
			log.Printf("Error reading dir: %s - %s", dir, err)
			continue
		}

		if fi.IsDir() {
			configFileName := filepath.Join(dir, COSMOFSCONFIGFILE)
			_, err := os.Lstat(configFileName)

			if err != nil {
				log.Printf("Error config file does not exists: %s", err)

				// Create the config file.
				configFile, err := os.Create(configFileName)

				if err != nil {
					log.Printf("Error creating config file: %s", err)
					continue
				}

				// Read the directory and include the files on it.
				file, err := os.Open(dir)

				if err != nil {
					log.Printf("Error reading dir: %s - %s", dir, err)
					continue
				}

				fi, err := file.Readdir(0)

				if err != nil {
					log.Printf("Error reading dir contents: %s - %s", dir, err)
					continue
				}

				files := make([]*cosmofs.File, len(fi))

				for i, ent := range fi {
					log.Println(ent)
					files[i] = &cosmofs.File{
						Path: dir,
						Filename: ent.Name(),
					}
					log.Println(files[i])
				}

				configEnc := gob.NewEncoder(configFile)

				err = configEnc.Encode(len(files))

				if err != nil {
					log.Fatal("Error encoding length config file: ", err)
				}

				for i := range files {
					err = configEnc.Encode(files[i])
				}

				if err != nil {
					log.Fatal("Error encoding list of files config file: ", err)
				}
			}

			// Decode the config file and update data structures.
			configFile, err := os.Open(configFileName)

			if err != nil {
				log.Printf("Error opening config file: %s", err)
				continue
			}

			configDec := gob.NewDecoder(configFile)

			var numFiles int

			err = configDec.Decode(&numFiles)

			if err != nil {
				log.Fatal("Error decoding length config file: ", err)
			}

			log.Printf("DECODED LENGTH VALUE: %v", numFiles)

			var decodedFiles []*cosmofs.File = make([]*cosmofs.File, numFiles)

			for i := range decodedFiles {
				decodedFiles[i] = new(cosmofs.File)
				err = configDec.Decode(decodedFiles[i])
				if err != nil {
					log.Fatal("Error decoding list of files config file: ", err)
				}
				log.Printf("DECODED VALUES: %v", decodedFiles[i])
			}
		}
	}

	// Leave the process listening for other peers
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
