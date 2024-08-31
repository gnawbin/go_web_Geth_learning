package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	store "go-eth/abi" // for demo
)

func main() {
	client, err := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("746ef5b96fa933740daf2aa3d2740de8f72b7689a7b2f3c086bc8bf3967c32d3")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress("0x0D7D359dbe3dB61af0bC45D35a9ff06Db0Bad197")
	instance, err := store.NewAbi(address, client)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := instance.Inc(auth)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	result, err := instance.Get(&bind.CallOpts{From: fromAddress, Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Big int: %v\n", result)

}
