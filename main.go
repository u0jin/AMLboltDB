package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

// 단순하고 가벼움
// go작성
// 별도 서버 필요없음
// 데이터 구조 설계 자유로움

// 작업 증명을 채굴이라 한다 >> 작업시간은 블록체인 내부에서 설정한 난이도에 따라 다름
// 작업 증명은 단순히 말하면, 끝 자리가 0으로 끝나는 비트의 수에 맞는 해시값을 찾는것 > 단순하지만, 해시값을 비교해보는 수밖에 없음

const targetBits = 24

type Blockchain struct { // 블록체인은 다수의 블록을 가진다. 애초에 블록체인 자체가 블록들의 연결이기때문임
	blocks []*Block
}

type Block struct { // 블록체인은 블록들을 체인형태로 연결한것 >> 블록의 역할을 하는것이 Block 타입임 > 이건 블록헤더임
	PrevBlockHash []byte
	Hash          []byte
	Timestamp     int64
	Data          []byte // 머크트리루트 대신 그냥 데이터 가지도록 함
	Nonce         int    // 블록이 작업증명으로 채굴했다는 증명
}
type ProofOfWork struct {
	block  *Block   // 해시값을 찾을 블록
	target *big.Int // 타겟값
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

// 블록바디는 일반적으로 거래 >> 트랜잭션을 가지고 있음,지금은 트잭이 없으므로 생략
func NewBlock(data string, prevBlockHash []byte) *Block { // 외부에서 머클트리 루트를 대신할 데이터와 이전블록해시를 파라매터로 받음
	block := &Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data), 0}
	pow := newProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func newProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1) // 숫자 크기 때문에 big.Int 사용
	// 1
	target.Lsh(target, uint(256-targetBits))
	// shift연산자 target = 1 << 256-targetBits
	// 256 == 해시함수 sha256알고리즘 결과값이 256bit에 해당하기 때문

	pow := &ProofOfWork{b, target}
	return pow
}
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	// byte 슬라이스
	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		fmt.Println(err)
	}

	return buff.Bytes()
}

// nonce 라는 것은 단순 counter이며, 블록헤더와 결합되어 target보다 더 작은 값을 찾기 위한 데이터를 준비시키기 위한 메서드
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		intToHex(pow.block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(nonce)),
	}, []byte{})
	return data
}
func (pow *ProofOfWork) Validate(b *Block) bool {
	var hashInt big.Int
	data := pow.prepareData(b.Nonce)
	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
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

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	max_number := ^uint(0)
	nonce := 0

	fmt.Printf("블록 마이닝 시작!!! %s\n", pow.block.Data)
	for uint(nonce) < max_number {
		data := pow.prepareData((nonce))
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println("마이닝 성공!!!")
	return nonce, hash[:]
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to ujin")
	bc.AddBlock("Send 2 more BTC to ujin")

	fmt.Println()

	for _, b := range bc.blocks {

		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Println()

		pow := newProofOfWork(b)
		fmt.Println("pow:", pow.Validate(b))
		fmt.Println()

	}
}
