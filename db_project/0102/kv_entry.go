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
	//then construct Entry struct: entry := Entry{key: data[8:len(key)], val: data[8+len(key):]}
	//comments: still need to find out how to use read 100%. But i think that's it.
}

// QzBQWVJJOUhU https://trialofcode.org/
