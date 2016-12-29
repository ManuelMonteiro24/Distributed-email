package kademlia

import (
    "bytes"
    "encoding/gob"
    "golang.org/x/crypto/openpgp"
    "math/rand"
)

type Onion struct {
    Next NetworkNode
    Data []byte
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
        buc = rand.Intn(160)
        l = len(ht.RoutingTable[buc])
        e = extracted[buc]
        if l > e {
            onion_nodes = append(onion_nodes, ht.RoutingTable[buc][l-1-e].NetworkNode)
            extracted[buc] += 1
        }
    }

    return onion_nodes
}

func buildOnion(onion_nodes []*NetworkNode, data []byte) ([]byte, error) {
    // build onion from data, given nodes and their keys
    enc_buf := new(bytes.Buffer)
    dummy_list := []*openpgp.Entity{onion_nodes[0].PublicEntity}
    plaintext, err := openpgp.Encrypt(enc_buf, dummy_list, nil, nil, nil)

    plaintext.Write(data)
    plaintext.Close()

    onion := Onion{
        *onion_nodes[0],
        enc_buf.Bytes(),
    }

    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err = enc.Encode(onion)
    if err != nil {
        return nil, nil
    }

    if len(onion_nodes) > 1 {
        return buildOnion(onion_nodes[1:], buf.Bytes())
    } else {
        return buf.Bytes(), nil
    }
}

func removeOnionLayer(onion Onion, ent *openpgp.Entity) ([]byte, error) {
    // the function name is self-explaining
    buf := bytes.NewBuffer(onion.Data)
    var dummy_list openpgp.EntityList
    dummy_list = append(dummy_list, ent)
    md, err := openpgp.ReadMessage(buf, dummy_list, nil, nil)
    if err != nil {
        return nil, err
    }

    data_buf := new(bytes.Buffer)
    data_buf.ReadFrom(md.UnverifiedBody)

    return data_buf.Bytes(), nil
}
