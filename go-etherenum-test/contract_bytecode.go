package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x0D7D359dbe3dB61af0bC45D35a9ff06Db0Bad197")
	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(bytecode))
	// 跟https://sepolia.arbiscan.io/address/0x0D7D359dbe3dB61af0bC45D35a9ff06Db0Bad197#code 中的Deployed Bytecode一致
}
