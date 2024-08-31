package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"

	"math/big"
)

func main() {
	client, _ := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	privateKey, _ := crypto.HexToECDSA("746ef5b96fa933740daf2aa3d2740de8f72b7689a7b2f3c086bc8bf3967c32d3")

	// 获取公共地址
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	PublicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("PublicAddress:: ", PublicAddress.String())

	EthValue := big.NewInt(0)
	// 使用文章地址，防止出现其他问题
	toAddress := common.HexToAddress("0xE280029a7867BA5C9154434886c241775ea87e53")
	// 代币合约地址
	tokenAddress := common.HexToAddress("0xb1D4538B4571d411F07960EF2838Ce337FE1E80E")
	// 函数方法切片传递
	transferFnSignature := []byte("transfer(address,uint256)")
	// 生成签名hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	// 切片
	methID := hash.Sum(nil)[:4]
	// 填充地址
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	// 发送代币数目
	Number := new(big.Int)
	Number.SetString("10000000000000000000", 10)
	// 代币量填充
	var data []byte
	data = append(data, methID...)
	data = append(data, paddedAddress...)
	data = append(data, common.LeftPadBytes(Number.Bytes(), 32)...)

	// 使用方法估算燃气费
	gasLime, _ := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	fmt.Println("Ethereum GasLime:", gasLime)

	// 获取账户交易随机数
	Nonce, _ := client.PendingNonceAt(context.Background(), PublicAddress)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 进行交易事务
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    Nonce,
		GasPrice: gasPrice,
		Gas:      gasLime,
		To:       &tokenAddress,
		Value:    EthValue,
		Data:     data,
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
