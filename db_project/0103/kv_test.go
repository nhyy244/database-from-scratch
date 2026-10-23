package db0103

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVBasic(t *testing.T) {
	kv := KV{}
	kv.log.FileName = ".test_db"
	defer os.Remove(kv.log.FileName)

	os.Remove(kv.log.FileName)
	err := kv.Open()
	assert.Nil(t, err)
	defer kv.Close()

	updated, err := kv.Set([]byte("k1"), []byte("v1"))
	assert.True(t, updated && err == nil)

	val, ok, err := kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v1" && ok && err == nil)

	_, ok, err = kv.Get([]byte("xxx"))
	assert.True(t, !ok && err == nil)

	updated, err = kv.Del([]byte("xxx"))
	assert.True(t, !updated && err == nil)

	updated, err = kv.Del([]byte("k1"))
	assert.True(t, updated && err == nil)

	_, ok, err = kv.Get([]byte("xxx"))
	assert.True(t, !ok && err == nil)

	updated, err = kv.Set([]byte("k2"), []byte("v2"))
	assert.True(t, updated && err == nil)

	// reopen
	kv.Close()
	err = kv.Open()
	assert.Nil(t, err)

	_, ok, err = kv.Get([]byte("k1"))
	assert.True(t, !ok && err == nil)
	val, ok, err = kv.Get([]byte("k2"))
	assert.True(t, string(val) == "v2" && ok && err == nil)
}

func TestEntryEncode(t *testing.T) {
	ent := Entry{key: []byte("k1"), val: []byte("xxx")}
	data := []byte{2, 0, 0, 0, 3, 0, 0, 0, 0, 'k', '1', 'x', 'x', 'x'}

	assert.Equal(t, data, ent.Encode())

	decoded := Entry{}
	err := decoded.Decode(bytes.NewBuffer(data))
	assert.Nil(t, err)
	assert.Equal(t, ent, decoded)

	ent = Entry{key: []byte("k1"), deleted: true}
	data = []byte{2, 0, 0, 0, 0, 0, 0, 0, 1, 'k', '1'}

	assert.Equal(t, data, ent.Encode())

	decoded = Entry{}
	err = decoded.Decode(bytes.NewBuffer(data))
	assert.Nil(t, err)
	assert.Equal(t, ent, decoded)
}
func TestPlayground(t *testing.T) {
	log := Log{}
	log.FileName = "db_logs"
	kv := KV{log: log}
	kv.Open()
	kv.Set([]byte("k1"), []byte("v1"))
	kv.Set([]byte("k2"), []byte("v2"))
	updated, err := kv.Set([]byte("k2"), []byte("v2"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("updated should be false: ", updated)
	assert.False(t, updated)

	val, _, _ := kv.Get([]byte("k1"))
	fmt.Println("value of k1: ", string(val)) //need to automate the conversion from bytes

	deleted, err := kv.Del([]byte("k1"))
	fmt.Println("deleted should be true: ", deleted)
	assert.True(t, deleted)
	kv.Set([]byte("k3"), []byte("v3"))
	fmt.Println("DB state before closing: ")
	for k, v := range kv.mem {
		fmt.Println(k, ":", v)
	}
	kv.Close()
	updated, err = kv.Set([]byte("k2"), []byte("v3"))
	if err != nil {
		fmt.Println("DB is now closed. Should print error: ", err) //
	}
	kv2 := KV{log: log}
	kv2.Open()
	fmt.Println("DB after recovering state: ")
	for k, v := range kv2.mem {
		fmt.Println(k, ":", v)
	}
	kv2.Close()

}

// QzBQWVJJOUhU https://trialofcode.org/
