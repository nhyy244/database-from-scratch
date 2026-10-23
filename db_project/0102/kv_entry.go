package db0102

import (
	"encoding/binary"
	"io"
)

type Entry struct {
	key []byte
	val []byte
}

func (ent *Entry) Encode() []byte {

	data := make([]byte, 4+4+len(ent.key)+len(ent.val))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.val)))
	copy(data[8:], ent.key)
	copy(data[8+len(ent.key):], ent.val)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	//find the length of key and value => binary.LittleEndian.Uint32(b). Converts bytes to uint32
	//then construct Entry struct: entry := Entry{key: data[8:8+len(key)], val: data[8+len(key):]}
	//caller doesn't know how many bytes to allocate, so we use Read(). we know that the serialized data must have at least 8 bytes
	//I would just give data as the argument instead of the reader... but I guess io.Reader is passed here for flexibility? to handle more data sources

	data := make([]byte, 8)
	read := 0
	for read < len(data) { // could also use io.ReadFull()
		n, err := r.Read(data[read:]) //Read() reads UP TO len(header), so no guarantee that will read everything in one call
		if err != nil {
			return err
		}
		read += n
	}
	keySize := binary.LittleEndian.Uint32(data[0:4])
	valSize := binary.LittleEndian.Uint32(data[4:8])
	//fmt.Printf("header=%v keySize=%d valSize=%d\n", data, keySize, valSize)

	data = append(data, make([]byte, keySize+valSize)...)
	for read < len(data) {
		n, err := r.Read(data[read:])
		if err != nil {
			return err
		}
		read += n
	}
	ent.key = data[8 : 8+keySize]
	ent.val = data[8+keySize:]
	//fmt.Printf("data=%v", data)

	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
