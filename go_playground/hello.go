package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Entry struct {
	key []byte
	val []byte
}

func main() {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, 1523)

	ent := Entry{key: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), val: b}

	data := make([]byte, 4+4+len(ent.key)+len(ent.val))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.val)))
	copy(data[8:], ent.key)
	copy(data[8+len(ent.key):], ent.val)
	fmt.Println("data:", data)
	fmt.Println("length of data: ", len(data))
	fmt.Println("length of a byte", binary.LittleEndian.Uint32([]byte{243, 5, 0, 0}))

	reader := bytes.NewBuffer(data)
	n, _ := reader.Read(data[8:])
	fmt.Println("n: ", n)
}
