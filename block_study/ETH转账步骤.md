# Go语言实现以太坊ETH转账步骤

## 概述

本文档介绍如何使用Go语言和`go-ethereum`客户端库在以太坊网络上进行ETH转账。

## 前置准备

### 1. 安装依赖

```bash
go get github.com/ethereum/go-ethereum
```

### 2. 准备材料

- 发送方的私钥（Private Key）
- 接收方的公钥地址（Public Address）
- 以太坊节点RPC地址（如Infura）

## 转账步骤详解

### 步骤1: 连接以太坊网络

```go
client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_INFURA_KEY")
defer client.Close()
```

- 使用`ethclient.Dial()`连接到以太坊节点
- 可以连接到主网、测试网（如Sepolia）或本地节点
- 记得在使用完毕后关闭连接

### 步骤2: 获取链ID

```go
chainID, _ := client.ChainID(context.Background())
```

- 获取当前连接的区块链网络ID
- 用于后续交易签名（EIP-155标准）

### 步骤3: 加载发送方私钥

```go
privateKey, err := crypto.HexToECDSA(formPrivateKey)
```

- 将十六进制格式的私钥转换为ECDSA私钥对象
- **注意**: 私钥需妥善保管，不要泄露

### 步骤4: 获取发送方公钥和地址

```go
publicKey := privateKey.Public()
publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
formAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
```

- 从私钥推导出公钥
- 将公钥转换为以太坊地址格式

### 步骤5: 获取交易Nonce值

```go
nonce, err := client.PendingNonceAt(context.Background(), formAddress)
```

- Nonce是账户发送交易的序号
- 用于防止重放攻击
- 每次交易后自动递增

### 步骤6: 设置转账金额

```go
value := big.NewInt(10000000000000000) // 0.01 ETH (单位: wei)
```

- 金额单位为**wei**（1 ETH = 10^18 wei）
- 使用`big.Int`处理大数值

### 步骤7: 设置Gas参数

```go
gasLimit := uint64(21000)
gasPrice, err := client.SuggestGasPrice(context.Background())
```

- **Gas Limit**: 转账ETH的标准限制为21000
- **Gas Price**: 从网络获取建议的Gas价格

### 步骤8: 设置接收方地址

```go
toAddress := common.HexToAddress(toPublicKey)
```

- 将十六进制字符串转换为以太坊地址对象

### 步骤9: 创建未签名交易

```go
tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
```

- 组装交易对象，包含所有必要参数
- 最后一个参数为交易数据（普通转账为nil）

### 步骤10: 对交易进行签名

```go
signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
```

- 使用私钥对交易进行签名
- 采用EIP-155签名标准（包含链ID）

### 步骤11: 广播交易到网络

```go
err = client.SendTransaction(context.Background(), signTx)
```

- 将已签名的交易发送到以太坊网络
- 交易进入待处理池（mempool）

### 步骤12: 获取交易哈希

```go
log.Println("交易hash:", signTx.Hash().Hex())
```

- 获取交易的唯一标识哈希值
- 可用于在区块链浏览器上查询交易状态

## 关键参数说明

| 参数             | 说明            | 示例                                         |
|----------------|---------------|--------------------------------------------|
| Private Key    | 发送方私钥（不含0x前缀） | 64位十六进制字符串                                 |
| Public Address | 接收方地址（含0x前缀）  | 0x7997135454469971a9953a7b4a39dd757a8e3fdb |
| Chain ID       | 网络标识          | Sepolia测试网: 11155111                       |
| Nonce          | 交易序号          | 自动获取                                       |
| Value          | 转账金额（wei）     | 10000000000000000 = 0.01 ETH               |
| Gas Limit      | Gas上限         | 21000（ETH转账标准值）                            |
| Gas Price      | Gas价格（wei）    | 动态获取建议值                                    |

## 注意事项

1. **安全性**:
    - 永远不要在代码中硬编码私钥
    - 使用环境变量或密钥管理服务存储私钥

2. **Gas费用**:
    - 实际花费 = Gas Limit × Gas Price
    - 确保账户有足够余额支付转账金额 + Gas费用

3. **网络选择**:
    - 测试时使用测试网（Sepolia、Goerli等）
    - 测试网ETH可从水龙头免费获取

4. **交易确认**:
    - 交易发送后不代表已完成
    - 需要等待矿工打包确认（通常需要12-30个区块）

5. **错误处理**:
    - 生产环境中必须对每个错误进行妥善处理
    - 建议添加交易状态查询和重试机制

## 扩展功能

- 查询账户余额: `client.BalanceAt()`
- 等待交易确认: `bind.WaitMined()`
- 查询交易收据: `client.TransactionReceipt()`
- 估算Gas: `client.EstimateGas()`