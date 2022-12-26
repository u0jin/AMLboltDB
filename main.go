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

type Blockchain struct { // 블록체인은 다수의 블록을 가진다. 애초에 블록체인 자체가 블록들의 연결이기때문임
	blocks []*Block
}

type Block struct { // 블록체인은 블록들을 체인형태로 연결한것 >> 블록의 역할을 하는것이 Block 타입임 > 이건 블록헤더임
	PrevBlockHash []byte
	Hash          []byte
	Timestamp     int64
	Data          []byte // 머크트리루트 대신 그냥 데이터 가지도록 함
}

// 블록바디는 일반적으로 거래 >> 트랜잭션을 가지고 있음,지금은 트잭이 없으므로 생략
func NewBlock(data string, prevBlockHash []byte) *Block { // 외부에서 머클트리 루트를 대신할 데이터와 이전블록해시를 파라매터로 받음
	// 새로운 블록을 만듬
	// Hash는 별도로 설정함

	block := Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data)}

	block.SetHash()
	return &block
}

// 블록해시를 설정함
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// 부호 있는 정수를 문자열로 변환

	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	// 블록 헤더들을 모아서 해시함수를 통해 블록의 해시를 만들어 낸다.
	// 해시함수 == sha256 알고리즘 :: 32바이트의 해시값이 반환됨

	b.Hash = hash[:]
}
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) AddBlock(data string) { // 블록에 새로운 블록포함.
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock) // 블록배열에 추가시켜주면 됨
}
func NewGenesisBlock() *Block { // 블록을 새로만듬 >> 제네시스 블록을 미리 포함시킴  :: 제네시스 블록은 이전 블록이 없기 때문에 비워둬야한다.
	return NewBlock("Genesis Block", []byte{})
}
func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to ujin")
	bc.AddBlock("Send 2 more BTC to ujin")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
