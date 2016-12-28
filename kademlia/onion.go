package kademlia

import (
    "bytes"
    "errors"
    "math/rand"
    "encoding/gob"
    "golang.org/crypto/x/openpgp"
)


type Onion struct {
    Next NetworkNode,
    Data []byte
}



func getRandomNodesForOnion(ht *hashTable) (onion_nodes []*NetworkNode) {
    var buc, l, e int
    var extracted [160]int
    n := ht.totalNodes()

    if n > 3 {
        n = 3
    }

    for len(onion_nodes) < n {
        buc = rand.Intn(160)
        l = len(RoutingTable[buc])
        e = extracted[buc]
        if l > e {
            onion_nodes = append(onion_nodes, RoutingTable[buc][l-1-e])
            extracted[buc] += 1
        }
    }

    return onion_nodes
}

func buildOnion(onion_nodes []*NetworkNode, data []byte) ([]byte, err){
    enc_buf := new(bytes.Buffer)
    dummy_list := []*Entity{onion_nodes[0]}
    plaintext, err := openpgp.Encrypt(enc_buf, dummy_list, nil, nil, nil) 

    plaintext.Write(data)
    plaintext.Close()

    onion := Onion{
            Next *onion_nodes[0],
            Data enc_buf.Bytes()
    }

    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err = enc.Encode(onion)
    if err != nil {
        return nil
    }

    if len(onion_nodes) > 1 {
        return buildOnion(onion_nodes[1:], buf.Bytes())
    } else {
        return buf.Bytes(), nil
    }
}

func removeOnionLayer(onion Onion, ent *Entity) ([]byte, err){
    buf := bytes.NewBuffer(onion.Data)
    dummy_list := []*Entity{ent}

    md, err := openpgp.ReadMessage(buf, dummy_list, nil, nil)
    if err != nil {
        return nil, err
    }

    data_buf := new(bytes.Buffer)
    data_buf.ReadFrom(md.UnverifiedBody)

    return data_buf.Bytes(), nil
}
