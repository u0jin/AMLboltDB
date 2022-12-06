package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

// 단순하고 가벼움
// go작성
// 별도 서버 필요없음
// 데이터 구조 설계 자유로움

type Blockchain struct {
	db *bolt.DB
	l  []byte
}

type Block struct {
	Transactions []*Transaction
}

const (
	Blockchain = "blocks"
	dbFile     = "chain.db"
)

func NewBlockchain() *Blockchain {
	db, err := bolt.Open(dbFile, 0600, nil)

}

// 데이터 타입 존재안함.. 키와 바이트 배열로 저장
// 구조체를 직렬화해서 저장 >> 구조체를 바이트배열로 변환, 바이트배열 -> 구조체 복원 예정

// Serialize before sending
func (b *Block) Serialize() []byte {
	var value bytes.Buffer

	encoder := gob.NewEncoder(&value)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatal("Encode Error:", err)
	}

	return value.Bytes()
}

// Deserialize block(not a method)
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal("Decode Error:", err)
	}

	return &block
}

func main() {
	chain := GetBlockchain()
	defer chain.db.Close()

	for i := 1; i < 10; i++ {
		chain.AddBlock(strconv.Itoa(i))
	}
	chain.ShowBlocks()
}
