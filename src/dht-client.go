package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"distmail/kademlia"
	"flag"
	"fmt"
	//"io/ioutil"
	"encoding/gob"
	b58 "github.com/jbenet/go-base58"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"gopkg.in/readline.v1"
)

func main() {
	var ip = flag.String("ip", "0.0.0.0", "IP Address to use")
	var port = flag.String("port", "0", "Port to use")
	var bIP = flag.String("bip", "", "IP Address to bootstrap against")
	var bPort = flag.String("bport", "", "Port to bootstrap against")
	var help = flag.Bool("help", false, "Display Help")
	var stun = flag.Bool("stun", true, "Use STUN")
	//var pkeyfile = flag.String("pkeyfile", "", "File containing PGP private key")
	//var username = flag.String("username", "", "Username")

	flag.Parse()

	//if *username == "" {
	//	displayFlagHelp()
	//	os.Exit(0)
	//}

	if *help {
		displayFlagHelp()
		os.Exit(0)
	}

	if *ip == "" {
		displayFlagHelp()
		os.Exit(0)
	}

	if *port == "" {
		displayFlagHelp()
		os.Exit(0)
	}

	//if *pkeyfile == "" {
	//	displayFlagHelp()
	//	os.Exit(0)
	//}

	var bootstrapNodes []*kademlia.NetworkNode
	if *bIP != "" || *bPort != "" {
		bootstrapNode := kademlia.NewNetworkNode(*bIP, *bPort)
		bootstrapNodes = append(bootstrapNodes, bootstrapNode)
	}

	options := kademlia.Options{
		BootstrapNodes: bootstrapNodes,
		IP:             *ip,
		Port:           *port,
		UseStun:        *stun,
	}

	dht, err := kademlia.NewDHT(&kademlia.MemoryStore{}, &options)
	if err != nil {
		panic(err)
	}

	fmt.Println("Opening socket..")

	if *stun {
		fmt.Println("Discovering public address using STUN..")
	}

	err = dht.CreateSocket()
	if err != nil {
		panic(err)
	}
	fmt.Println("..done")

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

	//TODO : read private key from file
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	public_key_id, err := dht.Store(kademlia.SerializePublicKey(&privkey.PublicKey), kademlia.Hashit(dht.GetSelfID()), true)
	fmt.Println("Stored public key with Id: " + public_key_id)

}
