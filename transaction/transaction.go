package transaction

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
