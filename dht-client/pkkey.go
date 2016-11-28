package main

import (
    "bytes"
    "golang.org/x/crypto/openpgp"
    "golang.org/x/crypto/openpgp/packet"
    "fmt"
)


func main () {
    entity,err := getEntityFromFile("priv.key")

    b := new(bytes.Buffer)

    err = entity.Serialize(b)
    if err != nil {
        panic(err)
    }

    r := packet.NewReader(b)

    entity2, err := openpgp.ReadEntity(r)


    fmt.Println(entity)
    fmt.Println(entity2)
}
