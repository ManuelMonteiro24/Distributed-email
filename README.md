# dISTmail
dISTmail is a distributed email system written in Go. It's implemented through a Kademlia DHT 
(based on [prettymuchbryce's implementation](https://github.com/prettymuchbryce/kademlia)).


## Node identification
A node registers by storing its RSA public key on the network and its network ID is the SHA1 hash of its public key.
A node's public key is stored on the hashtable key correspondent to SHA1(nodeID).

## Email sending
A sent e-mail is stored on the first available slot for the recipient. Nodes' "receiving slots" are the hashtable 
keys correspondent to SHA1(SHA1(SHA1( ... (nodeID) ... ) - the "composed SHA1" of the receiving node's ID. Note that
the first slot is SHA1(SHA1(nodeID)), because SHA1(nodeID) is reserved for storing the node's public key.

## Privacy
To hide the sender's identity, an e-mail is not stored in the network by the node who wants to send it, but rather by the
exit node of an Onion Routing Circuit, built by the node who wants to send the e-mail. The nodes used for this circuit are
random DHT fingers from the sending node.

## How to use it

* git clone this repo to your $GOROOT/src/ folder
* cd $GOROOT/src/distributed-email
* go run *.go

## Contacts
Upon registration a key pair associated to the user name is generated and stored in the Users directory. For sending and receiving the emails the program utilizes the public key contacts stored in the Contacts Directory, upon receiving a public key from another user this one must store it in the Contacts Directory in the format username_PublicKey (e.g. john_PublicKey).

---

This project was done in the context of a Cryptography and Communications Security course at IST, in Lisbon, 
in the Winter Semester of 2016/2017.

It's the first project we've done in Go and it's not linted, but we do plan to refactor it to make it more readable.
Feel free to take it and use it as you wish.
