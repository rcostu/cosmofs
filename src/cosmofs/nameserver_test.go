package cosmofs

import (
	"path/filepath"
	"testing"
)

var dir = "/Users/roberto/Documents/Facultad/5-Quinto/Sistemas Informáticos/pruebas_cosmofs"

func TestTable(t *testing.T) {
	Table.AddID("roberto@costumero.es")

	t.Log("Elements in the table: ", len(Table))

	Table.AddDir("roberto@costumero.es", dir, filepath.Base(dir), true)

	t.Log(Table.ListIDs())

	t.Log(Table.ListDirs("nanana"))
	t.Log(Table.ListDirs("roberto@costumero.es"))

	t.Log(Table.ListDir("roberto@costumero.es", "empty"))
	t.Log(Table.ListDir("roberto@costumero.es", "pruebas_cosmofs"))
	t.Log(Table.ListDir("roberto@costumero.es", "pruebas_cosmofs/out1"))

	t.Log(Table.IDExists("nonexistent"))
	t.Log(Table.IDExists("roberto@costumero.es"))

	t.Log(Table.SearchDir("empty"))
	t.Log(Table.SearchDir("out"))

	t.Log(Table.SearchFile("empty"))
	t.Log(Table.SearchFile("out"))

	t.Log(Table.Search("empty"))
	t.Log(Table.Search("out"))

	t.Log(checkID("a@a."))
	t.Log(checkID("aaa2.com"))
	t.Log(checkID("roberto@costumero.es"))

	t.Log(SplitPath("yo@aa/lelele"))
	t.Log(SplitPath("/Users/media/var"))
	t.Log(SplitPath("roberto@costumero.es/pruebas_cosmofs/out1"))

	t.Fatal("FIN")
}
