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
	t.Log(Table.ListDir("roberto@costumero.es", "/Users"))

	t.Log(Table.IDExists("nonexistent"))
	t.Log(Table.IDExists("roberto@costumero.es"))

	t.Log(Table.SearchDir("empty"))
	t.Log(Table.SearchDir("out"))

	t.Fatal("FIN")
}
