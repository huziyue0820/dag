package block

import (
	"bytes"
	"encoding/gob"
	"log"
)

/**
功能：对区块进行编码以进行存储
参数：
	无
返回：
	无
*/
func (b *Block) Serialize() []byte{
	var result bytes.Buffer  //用来存储序列化后的数据
	encoder := gob.NewEncoder(&result)    //初始化编码器

	err:=encoder.Encode(b)   //对区块进行编码
	if err != nil {
		log.Fatal("encode error:", err)
	}

	return result.Bytes()
}

/**
功能：对区块进行反序列化
参数：
	d 序列化后的区块
返回：
	Block 反序列化后的区块
*/
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err!=nil{
		log.Fatal("encode error:", err)
	}

	return &block
}