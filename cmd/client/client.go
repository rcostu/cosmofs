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
	list_dirs *bool = flag.Bool("dirs", false, "List directories")
	list_dir_id *string = flag.String("dirID", "", "List directories for ID")
	list_dir *string = flag.String("dir", "", "List a dir")

	list_ids *bool = flag.Bool("ids", false, "List all IDs")

	search *string = flag.String("s", "", "Search")
	search_dir *string = flag.String("sDir", "", "Search directory")
	search_file *string = flag.String("sFile", "", "Search File")
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

	decod := gob.NewDecoder(conn)

	if *list_dirs {
		fmt.Printf("List directories\n")
		//fmt.Fprintf(conn, "List Directories\n")
		_, err = conn.Write([]byte("List Directories\n"))

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		var dirs []string

		decod.Decode(&dirs)

		for _, v := range dirs {
			fmt.Println(v)
		}
	}

	if *list_ids {
		fmt.Printf("List IDs\n")
		_, err = conn.Write([]byte("List IDs\n"))

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		var ids []string

		decod.Decode(&ids)

		for _, v := range ids {
			fmt.Println(v)
		}
	}

	if *list_dir_id != "" {
		fmt.Printf("Listing directories for ID %s\n", *list_dir_id)

		_, err = conn.Write([]byte("List Directories ID\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		_, err = conn.Write([]byte(*list_dir_id+"\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		var dirs []string

		decod.Decode(&dirs)

		for _, v := range dirs {
			fmt.Println(v)
		}

		if dirs == nil {
			fmt.Printf("There are no entries for ID %s\n", *list_dir_id)
		}
	}

	if *list_dir != "" {
		fmt.Printf("Listing directory %s\n", *list_dir)

		_, err = conn.Write([]byte("List Directory\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		_, err = conn.Write([]byte(*list_dir+"\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		var files []string

		decod.Decode(&files)

		for _, v := range files {
			fmt.Println(v)
		}

		if files == nil {
			fmt.Printf("There are no entries for Directory %s\n", *list_dir)
		}
	}

	if *search != "" {
		fmt.Printf("Searching for %s\n", *search)

		_, err = conn.Write([]byte("Search\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		_, err = conn.Write([]byte(*search+"\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		var result []string

		decod.Decode(&result)

		for _, v := range result {
			fmt.Println(v)
		}

		if result == nil {
			fmt.Printf("There are no entries for %s\n", *search)
		}
	}

	if *search_dir != "" {
		fmt.Printf("Searching for %s\n", *search_dir)

		_, err = conn.Write([]byte("Search Directory\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		_, err = conn.Write([]byte(*search_dir+"\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		var result []string

		decod.Decode(&result)

		for _, v := range result {
			fmt.Println(v)
		}

		if result == nil {
			fmt.Printf("There are no entries for %s\n", *search_dir)
		}
	}

	if *search_file != "" {
		fmt.Printf("Searching for %s\n", *search_file)

		_, err = conn.Write([]byte("Search File\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		_, err = conn.Write([]byte(*search_file+"\n"))

		if err !=nil {
			log.Fatalf("Error: %s\n", err)
		}

		var result []string

		decod.Decode(&result)

		for _, v := range result {
			fmt.Println(v)
		}

		if result == nil {
			fmt.Printf("There are no entries for %s\n", *search_file)
		}
	}

}
