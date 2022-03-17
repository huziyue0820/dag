package blockchain

import (
	"dag/block"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//区块链结构
type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// 区块链迭代器 用于遍历区块链
type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

const dbFile = "blockchain.db" //存储区块链的文件
const blocksBucket = "blocks"  //用来存储区块的bucket

/**
功能：向区块链中添加区块
参数：
	data 区块中的数据
返回：
	无
*/
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte //添加区块的时候 最后一个区块的哈希

	err := bc.DB.View(func(tx *bolt.Tx) error { //bolt数据库的只读模式
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l")) //读取最后一个区块的哈希

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := block.NewBlock(data, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize()) //存入新的区块
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash) //更新叶子区块
		if err != nil {
			log.Panic(err)
		}

		bc.Tip = newBlock.Hash //更新叶子区块哈希

		return nil
	})

}

/**
功能：创建创世区块，初始化区块链
参数：
	无
返回：
	无
*/
func NewGenesisBlock() *block.Block {
	return block.NewBlock("Genesis Block", []byte{})
}

/**
功能：初始化区块链
参数：
	无
返回：
	无
*/
func NewBlockchain() *Blockchain {
	var tip []byte                          //叶子区块
	db, err := bolt.Open(dbFile, 0600, nil) //连接数据库
	if err != nil {
		log.Fatal("encode error:", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil { //如果不存在，即还没有创建区块链，也没有存储在文件中
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize()) //存储键值对 键为区块的哈希 值为序列化后的区块
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash) //存储键值对 键表示叶子区块 值为叶子区块的哈希
			tip = genesis.Hash                     //叶子区块的哈希
		} else {
			tip = b.Get([]byte("l")) //如果当前已经存在区块链，则读取当前区块链的叶子区块
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

/**
功能：创建一个区块链迭代器
参数：
	无
返回：
	创建的区块链迭代器
*/
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.DB}

	return bci
}

/**
功能：从区块链中返回下一个区块
参数：
	无
返回：
	下一个区块
*/
func (i *BlockchainIterator) Next() *block.Block {
	var pre_block *block.Block

	err := i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.CurrentHash)
		pre_block = block.DeserializeBlock(encodedBlock)  //得到上一个区块

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.CurrentHash = pre_block.PrevBlockHash  //向前遍历

	return pre_block
}

