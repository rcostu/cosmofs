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
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/base64"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

const (
	hostAlgoRSA = "ssh-rsa"
)

var (
	MyPrivatePeer *localPeer
	MyPublicPeer *Peer
	PeerList map[string]*Peer = make(map[string]*Peer)

	knownPeersFileName = filepath.Join(os.Getenv("HOME"), ".ssh", "cosmofs_known_peers")
)

type localPeer struct {
	id string
	key *rsa.PrivateKey
	rawKey []byte
}

type Peer struct {
	ID string
	PubKey *rsa.PublicKey
	RawKey []byte
}

func init() {
	buffer := parseKeyFile(*pubkeyFileName)

	key, _, id, ok := parsePubKey(buffer)

	if !ok {
		log.Fatal("Cannot parse Public Key File")
	}

	MyPublicPeer = &Peer{
		ID: string(id),
		PubKey: key.(*rsa.PublicKey),
		RawKey: buffer,
	}

	buffer = parseKeyFile(*privkeyFileName)

	key, err := parsePrivateKey(buffer)

	if err != nil{
		log.Fatal("Cannot parse Private Key File:", err)
	}

	MyPrivatePeer = &localPeer{
		id: string(id),
		key: key.(*rsa.PrivateKey),
		rawKey: buffer,
	}

	createKnownPeersFile()

	_, err = os.Lstat(knownPeersFileName)

	if err != nil {
		err := createKnownPeersFile()

		if err != nil {
			log.Printf("Error creating known peers file: %s", err)
		}
	}

	err = decodeKnownPeersFile()

	if err != nil {
		log.Printf("Error decoding known peers file: %s", err)
	}
}

func SearchPeer(id string) (*Peer, bool){
	if peer, ok := PeerList[id]; ok {
		return peer, ok
	}
	return nil, false
}

func StorePeer(peer *Peer) {
	PeerList[peer.ID] = peer
}

func createKnownPeersFile() (err error) {
	knownPeersFile, err := os.Create(knownPeersFileName)

	if err != nil {
		return err
	}

	configEnc := gob.NewEncoder(knownPeersFile)

	err = configEnc.Encode(PeerList)

	if err != nil {
		log.Fatal("Error encoding peer list in known peers file: ", err)
	}

	return err
}

func decodeKnownPeersFile() (err error) {
	knownPeersFile, err := os.Open(knownPeersFileName)

	if err != nil {
		log.Printf("Error opening config file: %s", err)
		return err
	}

	configDec := gob.NewDecoder(knownPeersFile)

	err = configDec.Decode(&PeerList)

	if err != nil {
		log.Fatal("Error decoding list of files config file: ", err)
	}

	return err
}

func encodeKnownPeersFile() (err error) {
	err = os.Remove(knownPeersFileName)

	if err != nil {
		log.Printf("Error removing known peers file: %s", err)
		return err
	}

	createKnownPeersFile()

	return err
}

func parseKeyFile(keyFileName string) ([]byte) {
	fi, err := os.Lstat(keyFileName)

	if err != nil {
		log.Fatal("Error: Cannot find SSH Key file.")
	}

	keyFile, err := os.Open(keyFileName)

	if err != nil {
		log.Fatal("Error: Cannot open SSH Key file.")
	}

	defer keyFile.Close()

	buffer := make([]byte, fi.Size())

	keyFile.Read(buffer)

	return buffer
}

// parsePrivateKey parses a private RSA PKCS1 Key
func parsePrivateKey(in []byte) (out interface{}, err error) {
	block, _ := pem.Decode(in)

	if block == nil {
		return nil, errors.New("SSH: no private key found")
	}

	out, err = x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return
}

// parsePubKey parses a Public SSH-RSA Key encoded in Base64 format
func parsePubKey(in []byte) (out interface{}, rest, id []byte, ok bool) {
	algo, key, id, ok := parseString(in)

	if !ok {
		return
	}

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(key)))

	_, err := base64.StdEncoding.Decode(dst, key)

	if err != nil {
		log.Println("Error decoding rsa key:", err)
		return
	}

	algo, key, ok = parseKey(dst)

	if !ok {
		return
	}


	switch string(algo) {
	case hostAlgoRSA:
		pubkey, rest, ok := parseRSA(key)
		return pubkey, rest, id, ok
	}
	panic("ssh: unknown public key type")
}

// parseRSA parses an RSA key according to RFC 4253, section 6.6.
func parseRSA(in []byte) (out *rsa.PublicKey, rest []byte, ok bool) {
	key := new(rsa.PublicKey)

	bigE, in, ok := parseInt(in)
	if !ok || bigE.BitLen() > 24 {
		return
	}
	e := bigE.Int64()
	if e < 3 || e&1 == 0 {
		ok = false
		return
	}
	key.E = int(e)

	if key.N, in, ok = parseInt(in); !ok {
		return
	}

	ok = true
	return key, in, ok
}

func parseString(in []byte) (kind, key, id []byte, ok bool) {
	parts := bytes.Split(in, []byte(" "))

	kind = parts[0]
	key = parts[1]
	id = parts[2]
	id = id[:len(id)-1]
	ok = true
	return
}

func parseKey(in []byte) (out, rest []byte, ok bool) {
	if len(in) < 4 {
		return
	}

	length := binary.BigEndian.Uint32(in)
	if uint32(len(in)) < 4+length {
		return
	}
	out = in[4 : 4+length]
	rest = in[4+length:]
	ok = true
	return
}

func parseInt(in []byte) (out *big.Int, rest []byte, ok bool) {
	contents, rest, ok := parseKey(in)
	if !ok {
		return
	}
	out = new(big.Int)

	if len(contents) > 0 && contents[0]&0x80 == 0x80 {
		// This is a negative number
		notBytes := make([]byte, len(contents))
		for i := range notBytes {
			notBytes[i] = ^contents[i]
		}
		out.SetBytes(notBytes)
		out.Add(out, big.NewInt(1))
		out.Neg(out)
	} else {
		// Positive number
		out.SetBytes(contents)
	}
	ok = true
	return
}
