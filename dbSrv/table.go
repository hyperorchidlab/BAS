package dbSrv

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/btcsuite/goleveldb/leveldb/filter"
	"github.com/btcsuite/goleveldb/leveldb/opt"
	"github.com/btcsuite/goleveldb/leveldb/util"
	"sync"
)

type BASTable struct {
	lock sync.RWMutex
	database *leveldb.DB
}
type Record struct {
	BAddr   []byte `json:"blockAddr"`
	Sig     []byte `json:"signature"`
	BType   uint8  `json:"blockType"`
	NType   uint8  `json:"netType"`
	NAddr   []byte `json:"netAddr"`
	ExtData string `json:"ext_data"`
}

func InitTable(path string) *BASTable {
	opts := opt.Options{
		Strict:      opt.DefaultStrict,
		Compression: opt.NoCompression,
		Filter:      filter.NewBloomFilter(10),
	}

	db, err := leveldb.OpenFile(path, &opts)
	if err != nil {
		panic(err)
	}

	return &BASTable{database: db}
}

func (book *BASTable) Find(ba *BasQuery) *Record {
	book.lock.RLock()
	defer book.lock.RUnlock()

	if has, err := book.database.Has(ba.BlockAddr, nil); !has || err != nil {
		fmt.Println(err)
		return nil
	}

	data, err := book.database.Get(ba.BlockAddr, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := &Record{}
	if err := json.Unmarshal(data, r); err != nil {
		fmt.Println(err)
		return nil
	}
	return r
}

func (book *BASTable) Save(req *RegRequest) error {

	book.lock.Lock()
	defer book.lock.Unlock()

	//if _, err := book.database.Has(req.BlockAddr, nil); err != nil {
	//	return err
	//}

	r := &Record{
		BAddr:   req.BlockAddr,
		Sig:     req.Sig,
		BType:   req.BTyp,
		NType:   req.NTyp,
		NAddr:   req.NetAddr,
		ExtData: req.ExtData,
	}
	b, e := json.Marshal(r)
	if e != nil {
		return e
	}

	wo := &opt.WriteOptions{
		Sync: true,
	}
	return book.database.Put(req.BlockAddr, b, wo)
}

func (book *BASTable)FindAll(keyl1,keyl2 int) []*Record {
	book.lock.RLock()
	defer book.lock.RUnlock()

	r:=&util.Range{Start: nil,Limit: nil}

	iter := book.database.NewIterator(r,nil)
	defer iter.Release()

	var rs []*Record

	for iter.Next(){
		k:=iter.Key()
		if len(k) <= keyl2 && len(k) > keyl1{
			rtmp:= &Record{}
			if err := json.Unmarshal(iter.Value(), rtmp); err != nil {
				continue
			}
			rs = append(rs,rtmp)
		}
	}

	return rs
}
