package db0105

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
	crc32   []byte
}

func (ent *Entry) Encode() []byte {
	data := make([]byte, 3*4+1+len(ent.key)+len(ent.val))

	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.key)))
	copy(data[13:], ent.key)
	if ent.deleted {
		data[12] = 1
	} else {
		binary.LittleEndian.PutUint32(data[8:12], uint32(len(ent.val)))
		copy(data[13+len(ent.key):], ent.val)
	}

	recordChecksum := crc32.ChecksumIEEE(data[4:])
	binary.LittleEndian.PutUint32(data[0:4], recordChecksum)
	return data
}

var ErrBadSum = errors.New("bad checksum")

func (ent *Entry) Decode(r io.Reader) error {
	var header [13]byte
	if _, err := io.ReadFull(r, header[:]); err != nil { //read 1
		return err
	}
	orignalChecksum := binary.LittleEndian.Uint32(header[:4])
	klen := int(binary.LittleEndian.Uint32(header[4:8]))
	vlen := int(binary.LittleEndian.Uint32(header[8:12]))
	deleted := header[12]

	record := make([]byte, klen+vlen) //read 2. the reader keeps track of how much it has read
	if _, err := io.ReadFull(r, record); err != nil {
		return err
	}
	currentChecksum := crc32.ChecksumIEEE(append(header[4:], record[:]...))
	if orignalChecksum != currentChecksum {
		return ErrBadSum
	}
	ent.key = record[:klen]
	if deleted != 0 {
		ent.deleted = true
	} else {
		ent.deleted = false
		ent.val = record[klen:]
	}

	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
