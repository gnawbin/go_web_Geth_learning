package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKey, err := crypto.HexToECDSA("746ef5b96fa933740daf2aa3d2740de8f72b7689a7b2f3c086bc8bf3967c32d3")
	if err != nil {
		log.Fatal(err)
	}

	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) //0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hexutil.Encode(signature)) //0xdcf2c0b51900bce2924d0694895a1e723a091e34a65b16aadb148cc8bc75b27b68b08ead051e01f78efab5bd98c6d5d59e2f44fcc1875fc8b1dd8293fcaf872201
}
