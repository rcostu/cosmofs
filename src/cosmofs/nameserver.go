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

package cosmofs

import (
	"log"
	"os"
	"strings"
)

type FileList []*File
type DirTable map[string]FileList
type IDTable map[string]DirTable

var (
	Table IDTable = make(IDTable)
)

func (t IDTable) AddID (id string) {
	t[id] = make(DirTable)
}

func (t IDTable) AddDir (id, dir string) {
	// Read the directory and include the files on it.
	fi, err := os.Lstat(dir)

	if err != nil {
		log.Printf("Error reading dir: %s - %s", dir, err)
	}

	if fi.IsDir() {
		file, err := os.Open(dir)

		if err != nil {
			log.Printf("Error reading dir: %s - %s", dir, err)
		}

		fi, err := file.Readdir(0)

		if err != nil {
			log.Printf("Error reading dir contents: %s - %s", dir, err)
		}

		files := make(FileList, 0)

		for _, ent := range fi {
			if strings.HasPrefix(ent.Name(), ".") {
				continue
			}
			log.Printf("%s",ent.Name())
			files = append(files, &File{
				Path: dir,
				Filename: ent.Name(),
				Size: ent.Size(),
			})
		}

		t[id][dir] = files
	}
}
