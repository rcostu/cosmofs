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
	"path/filepath"
	"regexp"
	"strings"
)

type FileList []*File
type DirTable map[string]FileList
type IDTable map[string]DirTable

type NameServerError struct {
	e error
}

func (e *NameServerError) Error() string {
	return "Error in the NameServer"
}

var (
	Table IDTable = make(IDTable)
)

// TODO: Do not add duplicate IDs
// TODO: Check ID correctness
func (t IDTable) AddID (id string) {
	if _, ok := t[id]; !ok {
		t[id] = make(DirTable)
	}
}

// TODO: Do not add duplicate dirs
func (t IDTable) AddDir (id, dir, baseDir string, recursive bool) (err error) {
	// Read the directory and include the files on it.
	fi, err := os.Lstat(dir)

	if err != nil {
		log.Printf("Error reading dir: %s - %s", dir, err)
		return err
	}

	if fi.IsDir() {
		file, err := os.Open(dir)

		if err != nil {
			log.Printf("Error reading dir: %s - %s", dir, err)
			return err
		}

		fi, err := file.Readdir(0)

		if err != nil {
			log.Printf("Error reading dir contents: %s - %s", dir, err)
			return err
		}

		files := make(FileList, 0)
		//globalBaseDir := filepath.Join(id, baseDir)

		for _, ent := range fi {
			if strings.HasPrefix(ent.Name(), ".") {
				continue
			}
			files = append(files, &File{
				localPath: filepath.Clean(dir),
				GlobalPath: filepath.Join(id,baseDir,ent.Name()),
				Filename: ent.Name(),
				Size: ent.Size(),
				IsDir: ent.IsDir(),
			})

			if recursive && ent.IsDir() {
				t.AddDir(id, filepath.Join(dir, ent.Name()),
				filepath.Join(baseDir, ent.Name()), recursive)
			}
		}

		t.AddID(id)
		t[id][baseDir] = files

		return err
	}
	return &NameServerError{}
}

func (t IDTable) ListIDs() (ids []string, err error) {
	if len(t) > 0 {
		for k := range t {
			ids = append(ids, k)
		}
		return ids, err
	}
	return nil, &NameServerError{}
}

func (t IDTable) ListDirs(id string) (dirs []string, err error) {
	if v, ok := t[id]; ok {
		for k := range v {
			dirs = append(dirs, filepath.Join(id, k))
		}
		return dirs, err
	}
	return nil, &NameServerError{}
}

func (t IDTable) ListDir (id, dir string) (content []string, err error) {
	if _, ok := t[id]; ok {
		if _, ok := t[id][dir]; ok {
			for _, file := range t[id][dir] {
				content = append(content, filepath.Join(id, dir, file.Filename))
			}
			return content, err
		}
	}
	return content, &NameServerError{}
}

func (t IDTable) IDExists (id string) (i string, err error) {
	if _, ok := t[id]; ok {
		return id, err
	}
	return "", &NameServerError{}
}

func (t IDTable) SearchDir (dir string) (result []string, err error) {
	if len(t) > 0 {
		found := false
		for k, v := range t {
			for d := range v {
				if strings.Contains(d, dir) {
					result = append(result, filepath.Join(k,d))
					found = true
				}
			}
		}
		if found {
			return result, err
		}
	}
	return result, &NameServerError{}
}

func (t IDTable) SearchFile (name string) (result []string, err error) {
	if len(t) > 0 {
		found := false
		for k, v := range t {
			for d, files := range v {
				for _, file := range files {
					if strings.Contains(file.Filename, name) && !file.IsDir {
						result = append(result, filepath.Join(k,d,file.Filename))
						found = true
					}
				}
			}
		}
		if found {
			return result, err
		}
	}
	return result, &NameServerError{}
}

func (t IDTable) Search (s string) (result []string, err error) {
	res1, err := t.SearchDir(s)

	if err != nil {
		return result, err
	}

	res2, err := t.SearchFile(s)

	if err != nil {
		return result, err
	}

	result = make([]string, len(res1) + len(res2))
	copy(result, res1)
	copy(result[len(res1):], res2)

	return result, err
}

func (t IDTable) DeleteID (id string) {
	delete (t, id)
}

func (t IDTable) DeleteDir (id, dir string) {
	if _, ok := t[id]; ok {
		if _, ok := t[id][dir]; ok {
			delete(t[id], dir)
			if len(t[id]) == 0 {
				t.DeleteID(id)
			}
		}
	}
}

func checkID (id string) (err error) {
	mailregexp, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`)

	if err != nil {
		return err
	}

	if mailregexp.MatchString(id) {
		return err
	}

	return &NameServerError{}
}

func SplitPath (path string) (id, dir string, err error) {
	res := strings.SplitN(path, "/", 2)

	if len(res) != 2 {
		return id, dir, &NameServerError{}
	}

	if err := checkID(res[0]); err != nil {
		return id, dir, &NameServerError{}
	}

	return res[0], res[1], err
}