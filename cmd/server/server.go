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
	"strings"
)

var (
	// Flags
	verbose *bool = flag.Bool("v", false, "Verbose output ON")
	cosmofsin *string = flag.String("cosmofsin", os.Getenv("COSMOFSIN"), "Location of incoming packages")
	cosmofsout *string = flag.String("cosmofsout", os.Getenv("COSMOFSOUT"), "Location of shared directories")

	// Shared Directory List
	sharedDirList []string = make([]string, len(filepath.SplitList(*cosmofsout)))
	// Shared Files List
	sharedFileList map[string] []*cosmofs.File = make(map[string] []*cosmofs.File)
)

const (
	COSMOFSDIR string = ".cosmofs"
	COSMOFSCONFIGFILE string = ".cosmofsconfig"
)

func debug (format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v)
	}
}

// Handles petitions from the peers.
func handlePetition (conn net.Conn) {
	debug("Connection made from: %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	if err != nil {
		debug("Error reading connection: %s", err)
		return
	}

	line = strings.TrimRight(line, "\n")

	// Listing directories
	if line == "List" {
		debug("List directories from: %s\n", conn.RemoteAddr())

		configEnc := gob.NewEncoder(conn)

		// Send the number of shared directories
		err = configEnc.Encode(len(sharedFileList))

		if err != nil {
			log.Fatal("Error sending length of files: ", err)
		}

		debug("%d directories shared", len(sharedFileList))

		// For each directory some data is sent to the client
		for dir, files := range sharedFileList {
			// Send the number of files in the current directory
			err = configEnc.Encode(len(files))
			if err != nil {
				log.Fatal("%d files found on directory: ", len(files))
			}

			// Send directory name
			err = configEnc.Encode(dir)

			if err != nil {
				log.Fatal("Error sending dir of files: ", err)
			}

			debug("Sent directory %s", dir)

			// Send each one of the file names
			for _, file := range files {
				err = configEnc.Encode(file)
				if err != nil {
					log.Fatal("Error sending file: ", err)
				}
				debug("Sent file: %s", file)
			}
		}
	}
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

	debug("Inbound files arriving to %s", *cosmofsin)

	// Initialize sharedDirList
	sharedDirList = filepath.SplitList(*cosmofsout)

	// There shall be at least one shared directory
	if len(sharedDirList) == 0 {
		log.Fatal("COSMOFSOUT should have at least one directory or file.")
	}

	// Shared directories are initialized
	for _, dir := range sharedDirList {
		dir = filepath.Clean(dir)
		debug("%s", dir)

		// Check wether we can read the current directory
		fi, err := os.Lstat(dir);

		if err != nil {
			debug("Error reading dir: %s - %s", dir, err)
			continue
		}

		// If it is a directory, look for the config file and decode it, or
		// generate it if it does not already exists.
		if fi.IsDir() {
			configFileName := filepath.Join(dir, COSMOFSCONFIGFILE)
			_, err := os.Lstat(configFileName)

			if err != nil {
				debug("Error config file does not exists: %s", err)

				// Create the config file.
				configFile, err := os.Create(configFileName)

				if err != nil {
					debug("Error creating config file: %s", err)
					continue
				}

				// Read the directory and include the files on it.
				file, err := os.Open(dir)

				if err != nil {
					debug("Error reading dir: %s - %s", dir, err)
					continue
				}

				fi, err := file.Readdir(0)

				if err != nil {
					debug("Error reading dir contents: %s - %s", dir, err)
					continue
				}

				files := make([]*cosmofs.File, len(fi))

				for i, ent := range fi {
					debug("%s",ent)
					files[i] = &cosmofs.File{
						Path: dir,
						Filename: ent.Name(),
					}
					debug("%s", files[i])
				}

				configEnc := gob.NewEncoder(configFile)

				err = configEnc.Encode(len(files))

				if err != nil {
					log.Fatal("Error encoding length config file: ", err)
				}

				for i := range files {
					err = configEnc.Encode(files[i])
					if err != nil {
						log.Fatal("Error encoding list of files config file: ", err)
					}
				}
			}

			// Decode the config file and update data structures.
			configFile, err := os.Open(configFileName)

			if err != nil {
				debug("Error opening config file: %s", err)
				continue
			}

			configDec := gob.NewDecoder(configFile)

			var numFiles int

			err = configDec.Decode(&numFiles)

			if err != nil {
				log.Fatal("Error decoding length config file: ", err)
			}

			debug("DECODED LENGTH VALUE: %v", numFiles)

			var decodedFile *cosmofs.File
			sharedFileList[dir] = make([]*cosmofs.File, numFiles)

			for i := 0; i < numFiles; i++ {
				decodedFile = new(cosmofs.File)
				err = configDec.Decode(decodedFile)
				if err != nil {
					log.Fatal("Error decoding list of files config file: ", err)
				}
				sharedFileList[dir][i] = decodedFile
				debug("DECODED VALUES: %v", sharedFileList[dir][i])
			}
		} // IsDir()

		// TODO: What if it is a file?? 
	}

	// Leave the process listening for other peers
	ln, err := net.Listen("tcp", ":5453")

	if err != nil {
		debug("Error: %s\n", err)
	}

	debug("Listening on address ", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			debug("Error: %s\n", err)
			continue
		}
		go handlePetition(conn)
	}
}
