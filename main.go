package main

// 비트코인 ​​전체 노드가 실행되고 네트워크와 동기화 필요

// 동기화 단계 - 노션 (https://www.notion.so/0890171bdb654a8f82f3a3a60d637eba?pvs=4)

import (
	"bytes"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func createDatabase(dbPath string) (*bolt.DB, error) {
	// 성공하면 *bolt.DB 포인터(BoltDB 인스턴스)를 반환
	// 데이터베이스 생성에 문제가 있으면 오류를 반환

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func putData(db *bolt.DB, bucketName string, key string, value string) error {
	// 지정된 버킷에 키-값 쌍을 삽입
	// 버킷이 없으면 새 버킷을 생성
	// 프로세스 중에 문제가 발생하면 이 함수는 오류를 반환

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func getData(db *bolt.DB, bucketName string, key string) (string, error) {
	// 인수로 문자열. 지정된 버킷에서 지정된 '키'와 관련된 값을 가져옴
	// 지정된 버킷에서 지정된 키와 연결된 값을 검색하여 문자열로 반환
	// 성공하면 값을 문자열로 반환하고 작업에 문제가 있는 경우(예: 버킷 또는 키를 찾을 수 없음) 오류 반환

	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Bucket not found: %s", bucketName)
		}
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func iterateBucket(db *bolt.DB, bucketName string) error {
	// BoltDB 인스턴스에 대한 포인터와 버킷 이름을 입력으로 받음
	// 지정된 버킷의 모든 키-값 쌍을 반복하여 출력, 오류가 발생하면 오류를 반환

	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Bucket not found: %s", bucketName)
		}

		return b.ForEach(func(k, v []byte) error {
			fmt.Printf("Key: %s, Value: %s\n", string(k), string(v))
			return nil
		})
	})
}

func getKeysWithPrefix(db *bolt.DB, bucketName string, prefix string) ([]string, error) {
	// 제공된 접두사가 있는 지정된 버킷에서 모든 키를 검색하여 문자열 조각으로 반환
	// 오류가 발생하면 오류를 반환

	var keys []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Bucket not found: %s", bucketName)
		}

		c := b.Cursor()
		prefixBytes := []byte(prefix)

		for k, _ := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, _ = c.Next() {
			keys = append(keys, string(k))
		}

		return nil
	})
	return keys, err
}

func main() {
	dbPath := "/Users/ujin/Desktop/Blockchain" // 데이터베이스 생성 경로
	db, err := createDatabase(dbPath)          // 데이터베이스 생성
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bucketName := "exampleBucket"

	err = putData(db, bucketName, "key1", "value1") // 키- 값 쌍 삽입
	if err != nil {
		log.Fatal(err)
	}

	err = putData(db, bucketName, "key2", "value2")
	if err != nil {
		log.Fatal(err)
	}

	value, err := getData(db, bucketName, "tx1") // 데이터 쿼리
	// "blockchain_data" 버킷에서 "tx1" 키와 연결된 값을 검색하여 프린트
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Value for key 'tx1': %s\n", value)

	err = iterateBucket(db, "blockchain_data")
	// "blockchain_data" 버킷의 모든 키-값 쌍을 반복하여 프린트
	if err != nil {
		log.Fatal(err)
	}

	keys, err := getKeysWithPrefix(db, "blockchain_data", "tx")
	//  접두사 "tx"가 있는 "blockchain_data" 버킷의 모든 키를 검색하고 프린트
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Keys with prefix 'tx': %v\n", keys)
}
