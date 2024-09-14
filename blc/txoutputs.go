package blc

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutputs struct {
	//存储所有的output
	UTXOS []*UTXO
}

// Serialize
func (txOutputs *TxOutputs) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// unserialize
func DeserializeTxOutputs(data []byte) *TxOutputs {
	var outputs TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return &outputs
}
