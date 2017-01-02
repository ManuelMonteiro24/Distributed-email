package kademlia

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"io"
)

type Onion struct {
	Next NetworkNode
	Data []byte
}

func randInt() int {
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randReader := bytes.NewBuffer(randBytes)
	rando, err := binary.ReadUvarint(randReader)
	if err != nil {
		panic(err)
	}
	return int(rando)
}

func getRandomNodesForOnion(ht *hashTable) (onion_nodes []*NetworkNode) {
	// get at most 3 random nodes from the client's bucket, to build the onion circuit
	var buc, l, e int

	var extracted [160]int
	n := ht.totalNodes()

	if n > 3 {
		n = 3
	}

	for len(onion_nodes) < n {
		buc = randInt() % 160
		l = len(ht.RoutingTable[buc])
		e = extracted[buc]
		if l > e {
			onion_nodes = append(onion_nodes, ht.RoutingTable[buc][l-1-e].NetworkNode)
			extracted[buc] += 1
		}
	}

	return onion_nodes
}

func (dht *DHT) getNodePubKey(node *NetworkNode) *rsa.PublicKey {
	pubkeyBytes, found, err := dht.Get(Hashit(dht.GetSelfID()))

	if !found {
		panic("failed to retrieve node public key")
	}
	if err != nil {
		panic(err)
	}

	return DeserializePublicKey(pubkeyBytes)
}

func BuildOnion(dht *DHT, nodes []*NetworkNode, data []byte) ([]byte, error) {
	// build onion from data, given nodes and their keys

	pubkey := dht.getNodePubKey(nodes[0])

	cipher := Encrypt(pubkey, data)

	onion := Onion{
		*nodes[0],
		cipher,
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(onion)
	if err != nil {
		panic(err)
	}

	if len(nodes) > 1 {
		return BuildOnion(dht, nodes[1:], buf.Bytes())
	} else {
		return buf.Bytes(), nil
	}
}

func RemoveOnionLayer(onion Onion, privkey *rsa.PrivateKey) ([]byte, error) {
	// the function name is self-explaining
	result := Decrypt(privkey, onion.Data)
	if len(result) > 0 {
		return result, nil
	} else {
		panic("error decrypting onion")
	}
}

func DecryptOnion(data []byte) (Onion, error) {
	oni := &Onion{}

	reader := bytes.NewBuffer(data)
	dec := gob.NewDecoder(reader)
	err := dec.Decode(oni)
	if err != nil {
		panic(err)
	}
	return *oni, nil
}

func Encrypt(pubkey *rsa.PublicKey, data []byte) []byte {
	key := make([]byte, 16)
	_, err := rand.Read(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	enc_sym_key, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, key, []byte(""))
	if err != nil {
		panic(err)
	}

	var keyLengthBytes [8]byte
	binary.PutUvarint(keyLengthBytes[:], uint64(len(enc_sym_key)))

	var result []byte

	result = append(result, keyLengthBytes[:]...)
	result = append(result, enc_sym_key[:]...)
	result = append(result, ciphertext[:]...)

	return result
}

func Decrypt(privkey *rsa.PrivateKey, data []byte) []byte {
	keyLengthBytes := data[:8]
	lengthReader := bytes.NewBuffer(keyLengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		panic(err)
	}

	enc_sym_key := data[8 : 8+length]
	ciphertext := data[8+length:]

	sym_key, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privkey, enc_sym_key, []byte(""))
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(sym_key)
	if err != nil {
		panic(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}

func SerializePublicKey(pubkey *rsa.PublicKey) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(pubkey)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DeserializePublicKey(pkeyBytes []byte) *rsa.PublicKey {
	reader := bytes.NewBuffer(pkeyBytes)
	pubkey := &rsa.PublicKey{}
	dec := gob.NewDecoder(reader)

	err := dec.Decode(pubkey)
	if err != nil {
		panic(err)
	}

	return pubkey
}
