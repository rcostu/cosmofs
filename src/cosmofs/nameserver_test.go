package cosmofs

import (
	"testing"
)

func TestTable(t *testing.T) {
	Table.AddID("roberto@costumero.es")

	t.Log(len(Table))
	for k := range Table {
		t.Log(k)
	}

	Table.AddDir("roberto@costumero.es", "/Users")
	v := Table["roberto@costumero.es"]
	vv := v["/Users"]
	for _, vvv := range vv {
			t.Log(vvv.Filename)
		}

	t.Fatal("FIN")
}
