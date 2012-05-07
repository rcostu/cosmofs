package cosmofs

import (
	"path/filepath"
	"testing"
)

var dir = "/Users/"

func TestTable(t *testing.T) {
	err := Table.AddID("prueba@prueba.es")

	if err != nil {
		t.Error("Failure in AddID")
	}

	t.Log("Elements in the table: ", len(Table))

	err = Table.AddDir("prueba@prueba.es", dir, filepath.Base(dir), false)

	if err != nil {
		t.Error("Failure in AddDir")
	}

	ids, err := Table.ListIDs()

	if err != nil && len(ids) != 2 {
		t.Error("Failure in ListIDs")
	}

	dirs, err := Table.ListDirs("nanana")

	if err == nil {
		t.Error("Failure in ListDirs, nanana should not exist.")
	}

	dirs, err = Table.ListDirs("roberto@costumero.es")

	if err != nil {
		t.Error("Failure in ListDirs.")
	}

	dirs, err = Table.ListDir("roberto@costumero.es", "empty")

	if err == nil {
		t.Error("Failure in ListDir. empty should not exist.")
	}

	dirs, err = Table.ListDir("roberto@costumero.es", "pruebas_cosmofs")

	if err == nil {
		t.Error("Failure in ListDir. pruebas_cosmofs should not exist.")
	}

	dirs, err = Table.ListDir("roberto@costumero.es", "out1")

	if err != nil {
		t.Error("Failure in ListDir.")
	}

	_, err = Table.ExistsID("nonexistent")

	if err == nil {
		t.Error("Failure in ExistsID. nonexistent should not exist.")
	}

	id, err := Table.ExistsID("roberto@costumero.es")

	if err != nil && id != "roberto@costumero.es" {
		t.Error("Failure in ExistsID.")
	}

	dirs, err = Table.SearchDir("empty")

	if err == nil {
		t.Error("Failure in SearchDir. empty should not exist")
	}

	dirs, err = Table.SearchDir("out")

	if err != nil {
		t.Error("Failure in SearchDir.")
	}

	dirs, err = Table.SearchFile("empty")

	if err == nil {
		t.Error("Failure in SearchFile. empty should not exist")
	}

	dirs, err = Table.SearchFile("out")

	if err != nil {
		t.Error("Failure in SearchFile.")
	}

	dirs, err = Table.Search("empty")

	if err == nil {
		t.Error("Failure in Search. empty should not exist")
	}

	dirs, err = Table.Search("out")

	if err != nil {
		t.Error("Failure in Search.")
	}

	err = checkID("a@a.")

	if err == nil {
		t.Error("Failure in checkID. a@a. should not exist")
	}

	err = checkID("aaa2.com")

	if err == nil {
		t.Error("Failure in checkID. aaa2.com should not exist")
	}

	err = checkID("roberto@costumero.es")

	if err != nil {
		t.Error("Failure in checkID.")
	}

	_, _, err = SplitPath("yo@aa/lelele")

	if err == nil {
		t.Error("Failure in SplitPath. yo@aa/lelele should not pass.")
	}

	_, _, err = SplitPath("/Users/media/var")

	if err == nil {
		t.Error("Failure in SplitPath. /Users/media/var should not pass.")
	}

	_, _, err = SplitPath("roberto@costumero.es/out1")

	if err != nil {
		t.Error("Failure in SplitPath.")
	}

	dirs, err = Table.ListDirs("prueba@prueba.es")

	if err != nil || len(dirs) != 1 {
		t.Log(Table.ListDirs("prueba@prueba.es"))
		t.Error("Failure in ListDirs. prueba@prueba.es should have 3 dirs.")
	}

	Table.DeleteDir("prueba@prueba.es", filepath.Base(dir))

	dirs, err = Table.ListDirs("prueba@prueba.es")

	if err != nil && len(dirs) != 0 {
		t.Log(Table.ListDirs("prueba@prueba.es"))
		t.Error("Failure in DeleteDir. prueba@prueba.es should not have any dirs.")
	}

	Table.DeleteID("prueba@prueba.es")

	t.Log("Elements in the table: ", len(Table))

	if len(Table) != 1 {
		t.Error("Failure in DeleteID.")
	}
}
