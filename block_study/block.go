package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// ethclient连接到以太坊测试网
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/a55cbed156ef4d638210967282f99266")
	if err != nil {
		log.Fatal(err)
	}

	// 获取当前链ID
	chainID, _ := client.ChainID(context.Background())
	// 指定操作的区块高度为9816183（指定区块）
	newInt := big.NewInt(9816183)
	// 获取区块头信息
	header, err := client.HeaderByNumber(context.Background(), newInt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("当前区块:", header.Number.String())

	// 获取区块信息
	block, _ := client.BlockByNumber(context.Background(), newInt)

	// block.Transactions() 获取当前区块的交易列表
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
		fmt.Println(tx.Value().String())

		if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Println("sender", sender.Hex())
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		marshal, _ := json.Marshal(receipt)
		fmt.Println("receipt", string(marshal))
		fmt.Println("status", receipt.Status)
		break
	}
	fmt.Println("交易数量：", len(block.Transactions()))

}
