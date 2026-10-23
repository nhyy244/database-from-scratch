package db0103

import "bytes"

type KV struct {
	log Log
	mem map[string][]byte
}

func (kv *KV) Open() error {
	if err := kv.log.Open(); err != nil {
		return err
	}
	kv.mem = map[string][]byte{}
	for {
		ent := Entry{}
		eof, err := kv.log.Read(&ent)
		if eof {
			break
		}
		if err != nil {
			return err
		}
		if ent.deleted {
			delete(kv.mem, string(ent.key))
		} else {
			kv.mem[string(ent.key)] = ent.val
		}
	}
	return nil
}
func (kv *KV) Close() error { return kv.log.Close() }

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, ok = kv.mem[string(key)]
	return
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	prev, exists := kv.mem[string(key)]
	updated = !bytes.Equal(prev, val) || !exists
	if updated {
		kv.mem[string(key)] = val
		ent := Entry{key: key, val: val, deleted: false}
		err = kv.log.Write(&ent)
		if err != nil {
			return false, err
		}
	}
	return updated, nil

}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	_, exists := kv.mem[string(key)]
	if !exists {
		return false, nil
	}

	delete(kv.mem, string(key))

	ent := Entry{key: key, deleted: true}
	err = kv.log.Write(&ent)
	if err != nil {
		return false, err
	}
	return true, nil
}

// QzBQWVJJOUhU https://trialofcode.org/
