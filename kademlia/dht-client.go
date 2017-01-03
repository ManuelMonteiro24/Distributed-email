package kademlia

import (
	"crypto/rsa"
	"fmt"
	"os"
	"os/signal"
)

func InitDHT(ID []byte, bIP string, bPort string, privkey *rsa.PrivateKey, extractor ExtractorFunc) (*DHT, string) {
	ip := ""
	port := ""

	var bootstrapNodes []*NetworkNode
	if bIP != "" || bPort != "" {
		bootstrapNode := NewNetworkNode(bIP, bPort)
		bootstrapNodes = append(bootstrapNodes, bootstrapNode)
	}

	options := Options{
		BootstrapNodes: bootstrapNodes,
		IP:             ip,
		Port:           port,
		UseStun:        true,
		ID:             ID,
		PrivKey:        privkey,
		mailExtractor:  extractor,
	}

	dht, err := NewDHT(&MemoryStore{}, &options)
	if err != nil {
		panic(err)
	}

	err = dht.CreateSocket()
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("Now listening on " + dht.GetNetworkAddr())
		err := dht.Listen()
		panic(err)
	}()

	if len(bootstrapNodes) > 0 {
		fmt.Println("Bootstrapping..")
		dht.Bootstrap()
		fmt.Println("..done")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			err := dht.Disconnect()
			if err != nil {
				panic(err)
			}
		}
	}()

	public_key_id, err := dht.Store(SerializePublicKey(&privkey.PublicKey), Hashit(dht.GetSelfID()), true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Stored public key with Id: " + public_key_id)

	return dht, public_key_id
}
