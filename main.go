package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// 단순하고 가벼움
// go작성
// 별도 서버 필요없음
// 데이터 구조 설계 자유로움

type Blockchain struct {
	blocks []*Block
}

type Block struct { // 블록체인은 블록들을 체인형태로 연결한것 >> 블록의 역할을 하는것이 Block 타입임 > 이건 블록헤더임
	PrevBlockHash []byte
	Hash          []byte
	Timestamp     int64
	Data          []byte // 머크트리루트 대신 그냥 데이터 가지도록 함
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}
func NewBlock(data string, prevBlockHash []byte) *Block {
	//block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}

	block := Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data)}

	block.SetHash()
	return &block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
