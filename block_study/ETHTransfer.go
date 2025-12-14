package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const formPrivateKey = "xxxxxx"
const toPublicKey = "0x7997135454469971a9953a7b4a39dd757a8e3fdb"

func main() {

	// ethclient连接到以太坊测试网
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/a55cbed156ef4d638210967282f99266")
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 获取当前链ID
	chainID, _ := client.ChainID(context.Background())
	// 加载私钥，获取私钥实例（发送方的私钥）
	privateKey, err := crypto.HexToECDSA(formPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 根据私钥获取公钥
	publicKey := privateKey.Public()

	// 根据公钥获取发送方地址
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	formAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("发送方address:", formAddress)

	// 获取交易随机数nonce
	nonce, err := client.PendingNonceAt(context.Background(), formAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("nonce:", nonce)

	// 设置要转账的ETH，单位：wei
	value := big.NewInt(10000000000000000)

	// 设置转账的gas上限（21000）
	gasLimit := uint64(21000)
	// 获取建议的gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 设置接收方的地址，入参为接收方的公钥，返回值为接收方的地址
	toAddress := common.HexToAddress(toPublicKey)
	log.Println("接收方address:", toAddress)

	// 生成未签名事务
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 发送方对交易进行签名
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 将已签名的交易广播到区块链网络
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("交易hash:", signTx.Hash().Hex())
}
