package basc

import (
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/btcsuite/goleveldb/leveldb/filter"
	"github.com/btcsuite/goleveldb/leveldb/opt"
	"github.com/hyperorchid/go-miner-pool/common"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"time"
)

const SimpleBasCacheTime = time.Hour

type basCacheItem struct {
	*dbSrv.NetworkAddr
	When time.Time `json:"expire"`
}

func (bci *basCacheItem) expired() bool {
	return bci.When.Before(time.Now())
}

type cachedClient struct {
	basIP    string
	database *leveldb.DB
}

func NewCachedBasCli(basIP, dbPath string) (BASClient, error) {
	opts := opt.Options{
		Strict:      opt.DefaultStrict,
		Compression: opt.NoCompression,
		Filter:      filter.NewBloomFilter(10),
	}

	db, err := leveldb.OpenFile(dbPath, &opts)
	if err != nil {
		return nil, err
	}

	c := &cachedClient{basIP: basIP, database: db}
	return c, nil
}

func (c *cachedClient) Query(ba []byte) (*dbSrv.NetworkAddr, error) {
	res := &basCacheItem{}
	if err := common.GetJsonObj(c.database, ba, res); err == nil && !res.expired() {
		return res.NetworkAddr, nil
	}

	ntAddr, err := QueryBySrvIP(ba, c.basIP)
	if err != nil {
		return nil, err
	}
	if ntAddr.NTyp == dbSrv.NoItem {
		return nil, fmt.Errorf("no such block chain address's[%s] ip address", ba)
	}

	res = &basCacheItem{
		When:        time.Now().Add(SimpleBasCacheTime),
		NetworkAddr: ntAddr,
	}
	_ = common.SaveJsonObj(c.database, ba, res)
	return ntAddr, nil
}

func (c *cachedClient) Register(req *dbSrv.RegRequest) error {
	return RegisterBySrvIP(req, c.basIP)
}
