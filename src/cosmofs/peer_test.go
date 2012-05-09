package cosmofs

import (
	"crypto/rsa"
	"os"
	"path/filepath"
	"testing"
)

func TestPeer(t *testing.T) {
	keyFileName := filepath.Join(os.Getenv("HOME"), ".ssh", "prueba.pub")

	fi, err := os.Lstat(keyFileName)

	if err != nil {
		t.Fatal("Error: Cannot find SSH Key file.")
	}

	keyFile, err := os.Open(keyFileName)

	if err != nil {
		t.Fatal("Error: Cannot open SSH Key file.")
	}

	defer keyFile.Close()

	buffer := make([]byte, fi.Size())

	t.Log(fi.Size())

	keyFile.Read(buffer)

	t.Logf("%s", buffer)

	key, _, id, ok := ParsePubKey(buffer)

	if !ok {
		t.Fatal("Error")
	}

	t.Logf("ID: %s", id)
	t.Logf("Exponente: %v", key.(*rsa.PublicKey).E)
	t.Logf("Modulo: %v", key.(*rsa.PublicKey).N)
	t.Fail()
}
