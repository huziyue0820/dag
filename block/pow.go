package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const targetBits = 12 //目标位数 难度值  x*4  x就是得到的哈希前面0的个数

type ProofOfWork struct {
	Block  *Block   //区块指针
	Target *big.Int //目标值
}

/**
功能：新建一个pow任务
参数：
	b pow解决的区块
返回：
	一个新的pow任务
*/
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)                  //初始化一个Big整型数据 值为1
	target.Lsh(target, uint(256-targetBits)) //左移 256-targetBits 位 得到用来比较的目标

	pow := &ProofOfWork{b, target}

	return pow
}

/**
功能：准备pow需要的数据，将区块信息（区块前哈希、区块数据、时间戳）和难度值targetBits以及临时值nonce合并
参数：
	nonce 参与计算的随机数
返回：
	用于计算的数据
*/
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int       //哈希的整数代表
	var hash [32]byte         //对pow准备的数据求hash
	nonce := 0                //计数器 pow所求的解
	maxNonce := math.MaxInt64 //求解的循环次数
	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce) //准备数据
		hash = sha256.Sum256(data)     //对数据求哈希
		fmt.Printf("\r%x", hash)       //清空当前并回到行首输出哈希
		hashInt.SetBytes(hash[:])      //使用hash字符串初始化哈希大整数

		if hashInt.Cmp(pow.Target) == -1 { //大整数比较
			break //如果小于 则成功找到解
		} else {
			nonce++ //失败则 ++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:] //返回找到的随机数 和 哈希结果
}

//验证pow结果是否有效
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.Target) == -1

	return isValid
}

//十进制数转成16进制字符串
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}