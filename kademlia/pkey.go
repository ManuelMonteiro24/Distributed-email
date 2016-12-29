package kademlia

import (
    "bytes"
    "golang.org/x/crypto/openpgp"
    "io/ioutil"
)

func GetEntityFromFile(path string) (privkey *openpgp.Entity, err error) {
    b, err := ioutil.ReadFile(path) 
    if err != nil {
        return nil, err
    }

    privkeyList,err := openpgp.ReadArmoredKeyRing(bytes.NewReader(b))
    if err != nil {
        return nil, err
    }
    privkey = privkeyList[0]

    return privkey, nil
}
