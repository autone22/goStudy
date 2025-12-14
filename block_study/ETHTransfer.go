package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

const formPrivateKey = "xxxxx"
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
	value := big.NewInt(1000000000000000)

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
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	})

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

	receipt, err := waitForTransaction(client, signTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == 1 {
		log.Println("交易成功")
		queryBalance(client, formAddress, toAddress)
	} else {
		log.Println("交易失败")
	}
}

// waitForTransaction 等待交易被打包确认
func waitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()

	// 设置超时时间（5分钟）
	timeout := time.After(5 * time.Minute)
	// 每3秒检查一次
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("等待交易确认超时")
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil {
				// 找到交易收据，说明交易已被打包
				return receipt, nil
			}
			// 如果是"not found"错误，继续等待
			log.Println("交易还未被打包，继续等待...")
		}
	}
}

func queryBalance(client *ethclient.Client, formAddress common.Address, toAddress common.Address) {
	formBalance, _ := client.BalanceAt(context.Background(), formAddress, nil)
	formBalanceEth := new(big.Float).Quo(new(big.Float).SetInt(formBalance), big.NewFloat(params.Ether))
	log.Printf("发送方余额: %s ETH (%s wei)\n", formBalanceEth.Text('f', 6), formBalance.String())

	toBalance, _ := client.BalanceAt(context.Background(), toAddress, nil)
	toBalanceEth := new(big.Float).Quo(new(big.Float).SetInt(toBalance), big.NewFloat(params.Ether))
	log.Printf("接收方余额: %s ETH (%s wei)\n", toBalanceEth.Text('f', 6), toBalance.String())
}
