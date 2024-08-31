package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	store "go-eth/abi" // for demo
)

func main() {
	client, err := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x0D7D359dbe3dB61af0bC45D35a9ff06Db0Bad197")
	instance, err := store.NewAbi(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")
	_ = instance
}
