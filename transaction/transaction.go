package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10 // 币基交易的价格  挖矿奖励的数量

type Transaction struct {
	ID   []byte     //交易ID
	Vin  []TXInput  //交易记录输入
	Vout []TXOutput //交易记录输出
}

//交易记录输出
type TXOutput struct {
	Value        int    //当前输出包含的价值
	ScriptPubKey string //用于锁定该笔输出的脚本
}

//交易记录输入
type TXInput struct {
	Txid      []byte //交易的输入来源  即一个输入对应一个输出
	Vout      int    //输出交易中的输出索引
	ScriptSig string //用于解锁输出的脚本
}

/**
功能：为交易设置交易号
参数：
	无
返回：
*/
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}


/**
功能：创建一个币基交易
参数：
	to 目标地址
	data 父区块的哈希
返回：
	*Transaction 新生成的币基交易
*/
func NewCoinbaseTx(to, data string) *Transaction{
	if data == ""{
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil , []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx

}