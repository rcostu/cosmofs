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
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileList []*File
type DirTable map[string]FileList
type IDTable map[string]DirTable

var (
	Table IDTable = make(IDTable)
	myID string
)

// TODO: Multiple different kinds of errors
type NameServerError struct {
	e error
}

func (e *NameServerError) Error() string {
	return "Error in the NameServer"
}

func init() {
	// Check if COSMOFSIN environment is set
	if *Cosmofsin == "" {
		log.Fatalf("COSMOFSIN not set correctly. Current content <%s>", *Cosmofsin)
	}

	// Check if COSMOFSIN is a correct directory
	if _, err := os.Lstat(*Cosmofsin); err != nil {
		log.Fatalf("COSMOFSIN not set correctly. Current content <%s>", *Cosmofsin)
	}

	sharedDirList := filepath.SplitList(*Cosmofsout)

	// There shall be at least one shared directory
	if len(sharedDirList) == 0 {
		log.Fatal("COSMOFSOUT should have at least one directory or file.")
	}

	// TODO: Fix this
	// HACK to get here myID correctly
	buffer := parseKeyFile(*pubkeyFileName)

	_, _, id, ok := parsePubKey(buffer)

	if !ok {
		log.Fatal("Cannot parse Public Key File")
	}

	myID = string(id)

	// Create a new user in the table
	// TODO: Decode and create correct ID
	err := Table.AddID(myID)

	if err != nil {
		log.Fatal("Could not create new ID")
	}

	// Shared directories are initialized
	for _, dir := range sharedDirList {
		dir = filepath.Clean(dir)

		// Check wether we can read the current directory
		fi, err := os.Lstat(dir);

		if err != nil {
			continue
		}

		// If it is a directory, look for the config file and decode it, or
		// generate it if it does not already exists.
		if fi.IsDir() {
			configFileName := filepath.Join(dir, COSMOFSCONFIGFILE)

			if *resetConfig {
				_, err := os.Lstat(configFileName)

				if err == nil {
					err := os.Remove(configFileName)
					if err != nil {
						log.Fatal("Error re-generating config files.")
					}
				}
			}

			_, err := os.Lstat(configFileName)

			if err != nil {
				err := createConfigFile(dir, configFileName)

				if err != nil {
					log.Printf("Error creating config file: %s", err)
					continue
				}
			}

			// Decode the config file and update data structures.
			err = decodeConfigFile(configFileName)
			if err != nil {
				log.Printf("Error decoding config file: %s", err)
				continue
			}
		}
	}
}

func (t IDTable) AddID (id string) (err error) {
	err = checkID(id)

	if err != nil {
		return err
	}

	if _, ok := t[id]; !ok {
		t[id] = make(DirTable)
	}

	return err
}

func (t IDTable) AddDir (id, dir, baseDir string, recursive bool) (err error) {
	// Check for existing dir
	err = t.ExistsDir(id,baseDir)

	if err == nil {
		log.Printf("Error: Dir %s already exists in the table", baseDir)
		return err
	}

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
				LocalPath: filepath.Clean(dir),
				GlobalPath: filepath.Join(id,baseDir,ent.Name()),
				Filename: ent.Name(),
				Size: ent.Size(),
				IsDir: ent.IsDir(),
				Owner: MyPublicPeer,
				KeepCopy: true,
				Online: false,
				NumChunks: 1,
				Chunks: nil,
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

func (t IDTable) ListAllDirs() (dirs []string, err error) {
	for id, v := range t {
		for k := range v {
			dirs = append(dirs, filepath.Join(id, k))
		}
	}
	return dirs, err
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

func (t IDTable) ExistsID (id string) (i string, err error) {
	if _, ok := t[id]; ok {
		return id, err
	}
	return "", &NameServerError{}
}

func (t IDTable) ExistsDir (id, dir string) (err error) {
	_, err = t.ExistsID(id)

	if err != nil {
		return &NameServerError{}
	}

	if _, ok := t[id][dir]; !ok {
		return &NameServerError{}
	}

	return err
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
	if _, ok := t[id][dir]; ok {
		delete(t[id], dir)
		if len(t[id]) == 0 {
			t.DeleteID(id)
		}
	}
}

func (t IDTable) ReceiveAndMergeTable (decod *gob.Decoder) {
	var recvTable IDTable = make(IDTable)

	err := decod.Decode(&recvTable)

	if err != nil {
		log.Fatal("Error decoding table: ", err)
	}

	log.Printf("LOCAL TABLE: %v\n", Table)
	log.Printf("REMOTE TABLE: %v\n", recvTable)

	for k, v := range recvTable {
		for d, files := range v {
			if _, ok := t[k][d]; !ok {
				t.AddID(k)
				t[k][d] = files
				log.Printf("Added dir %v from %v\n", d, k)
			}
		}
	}

	encodeConfigFiles()
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

	return res[0], filepath.Clean(res[1]), err
}

func createConfigFile(dir, configFileName string) (err error) {
	// Create the config file.
	configFile, err := os.Create(configFileName)

	if err != nil {
		return err
	}

	// Add directory and subdirectories
	// TODO: Should it be recursive?
	err = Table.AddDir(myID, dir, filepath.Base(dir), true)

	if err != nil {
		log.Printf("Error adding new directory: %s", err)
		return err
	}

	configEnc := gob.NewEncoder(configFile)

	err = configEnc.Encode(Table)
	if err != nil {
		log.Fatal("Error encoding table in config file: ", err)
	}

	return err
}

func decodeConfigFile(configFileName string) (err error){
	configFile, err := os.Open(configFileName)

	if err != nil {
		log.Printf("Error opening config file: %s", err)
		return err
	}

	configDec := gob.NewDecoder(configFile)

	err = configDec.Decode(&Table)

	if err != nil {
		log.Fatal("Error decoding list of files config file: ", err)
	}

	return err
}

func encodeConfigFiles() (err error){
	sharedDirList := filepath.SplitList(*Cosmofsout)

	// Shared directories are initialized
	for _, dir := range sharedDirList {
		dir = filepath.Clean(dir)

		// Check wether we can read the current directory
		fi, err := os.Lstat(dir);

		if err != nil {
			continue
		}

		// If it is a directory, look for the config file and decode it, or
		// generate it if it does not already exists.
		if fi.IsDir() {
			configFileName := filepath.Join(dir, COSMOFSCONFIGFILE)

			_, err := os.Lstat(configFileName)

			if err == nil {
				err := os.Remove(configFileName)
				if err != nil {
					log.Fatal("Error re-generating config files.")
				}
			}

			err = createConfigFile(dir, configFileName)

			if err != nil {
				log.Printf("Error creating config file: %s", err)
				continue
			}

			// Decode the config file and update data structures.
			err = decodeConfigFile(configFileName)
			if err != nil {
				log.Printf("Error decoding config file: %s", err)
				continue
			}
		}
	}

	return err
}

func PrintTable() {
	for k, v := range Table {
		log.Printf("- %v\n", k)
		for kk, vv := range v {
			log.Printf("-- %v\n", kk)
			for _, vvv := range vv {
				log.Printf("--- %v : %v : %v\n", vvv.Filename, vvv.GlobalPath,
				vvv.LocalPath)
			}
		}
	}
}

