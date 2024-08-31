package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-eth/abi"
)

func main() {
	//client, err := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/3e6d1dd20efe4d13bfdbc95865417e20")
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
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := abi.DeployAbi(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   //合约地址： 0x0D7D359dbe3dB61af0bC45D35a9ff06Db0Bad197
	fmt.Println(tx.Hash().Hex()) // 0x1f639b92aba60b8c6fcb5ec6d31a71151446d0f4da5b010ef81d925f831b5e94

	_ = instance
}
