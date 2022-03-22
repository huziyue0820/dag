package block

import (
	"bytes"
	"crypto/sha256"
	"dag/transaction"
	"time"
)

//区块结构
type Block struct {
	Timestamp     int64  //时间戳
	Transactions  []*transaction.Transaction //数据
	PrevBlockHash []byte //父区块哈希 链式一个 dag数组
	Hash          []byte //区块哈希
	Nonce         int //随机数 解pow用的
}

/**
功能：求新区块的哈希值
参数：
	无
返回：
	无
*/
//func (b *Block) SetHash() {
//	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))                       //时间戳
//	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{}) //区块头
//	hash := sha256.Sum256(headers)
//
//	b.Hash = hash[:] //区块哈希
//}

/**
功能：创建一个新的区块
参数：
	data 区块中存储的数据
	prevBlockHash 父区块的哈希
返回：
	*Block 新生成的区块
*/
func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	//block.SetHash()
	//return block

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

/**
功能：对区块中的交易进行哈希
参数：
	无
返回：
	交易哈希后的结果
*/
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

