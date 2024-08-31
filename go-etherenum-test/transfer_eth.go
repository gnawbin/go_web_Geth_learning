package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	if err != nil {
		log.Fatal(err)
	}

	// 加载私钥-需要去除0x
	privateKey, _ := crypto.HexToECDSA("746ef5b96fa933740daf2aa3d2740de8f72b7689a7b2f3c086bc8bf3967c32d3")
	// 获取公共地址
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	PublicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("PublicAddress:: ", PublicAddress.String())
	// 获取账户交易随机数
	Nonce, _ := client.PendingNonceAt(context.Background(), PublicAddress)
	value := big.NewInt(30000000000000000) // in wei (0.01 eth)
	gasLimit := uint64(3000000)            // in units.21000失败，但是3000000会成功
	//gasPricer := big.NewInt(100000000000000) // in wei (30 gwei)
	// 通过市场获取燃气价格
	gasPricer, _ := client.SuggestGasPrice(context.Background())
	// 接收地址 0xE280029a7867BA5C9154434886c241775ea87e53
	ToAddress := common.HexToAddress("0xE280029a7867BA5C9154434886c241775ea87e53")
	// 智能合约通信
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    Nonce,
		GasPrice: gasPricer,
		Gas:      gasLimit,
		To:       &ToAddress,
		Value:    value,
		Data:     make([]byte, 0),
	})
	// 利用私钥进行签名
	ChanID, _ := client.NetworkID(context.Background())
	SignedTx, _ := types.SignTx(tx, types.LatestSignerForChainID(ChanID), privateKey)
	// 广播节点
	if err := client.SendTransaction(context.Background(), SignedTx); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Transaction Success Address:: ", SignedTx.Hash().Hex())
}
