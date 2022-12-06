package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/boltdb/bolt"
)

// 단순하고 가벼움
// go작성
// 별도 서버 필요없음
// 데이터 구조 설계 자유로움

type Blockchain struct {
	db *bolt.DB
	l  []byte // 이미 만들어져 있는 블록체인을 조회하거나 순회하기 위함 ->> l 키를 가진 마지막 블록의 해시
}

type Block struct { // 블록체인은 블록들을 체인형태로 연결한것 >> 블록의 역할을 하는것이 Block 타입임 > 이건 블록헤더임
	prevBlockHash []byte
	Hash          []byte
	Timestamp     int64
	Data          []byte // 머크트리루트 대신 그냥 데이터 가지도록 함
	Nonce         int64
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const targetBits = 24

const (
	BlocksBuket = "blocks" // 버킷 == RDBMS 테이블과 비슷한개념 -> 버킷에 블록을 담아 조회하거나, 내용갱신 가능
	dbFile      = "chain.db"
)

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.block.prevBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(nonce),
		IntToHex(targetBits),
	}, []byte{})

	return data
}

// 블록생성시 작업증명을 무조건 거쳐야함 그리고 결과로 나온것들에 대해 블록에 적용시킬 필요있음

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data), 0}
	pow := NewProofOfWork(block)
	block.Nonce, block.Hash = pow.Run()
	//block.SetHash()

	return block
}

func NewBlockchain() *Blockchain { // 버킷 == RDBMS 테이블과 비슷한개념 -> 버킷에 블록을 담아 조회하거나, 내용갱신 가능
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var l []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBuket))

		if b == nil {
			// 새로운 블록체인 만들어야하는 경우
			b, err := tx.CreateBucket([]byte(BlocksBuket))

			if err != nil {
				log.Panic(err)
			}

			genesis := NewBlock("Genesis Block", []byte{})

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// "l" 키는 마지막 블록해시를 저장한다
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			l = genesis.Hash
		} else {
			// 이미 블록체인이 있는 경우
			l = b.Get([]byte("l"))
		}
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	return &Blockchain{db, l}
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
	/*
		chain := GetBlockchain()
		defer chain.db.Close()

		for i := 1; i < 10; i++ {
			chain.AddBlock(strconv.Itoa(i))
		}
		chain.ShowBlocks()
	*/
	for _, b := range bc.blocks {

		pow := NewProofOfWork(b)
		fmt.Println("pow:", pow.Validate(b))

		fmt.Println()
	}
}
