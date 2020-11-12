package basc

import (
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/btcsuite/goleveldb/leveldb/filter"
	"github.com/btcsuite/goleveldb/leveldb/opt"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/go-miner-pool/common"
	"github.com/hyperorchidlab/go-miner-pool/network"
	"net"
	"time"
)

const SimpleBasCacheTime = time.Hour

type basCacheItem struct {
	ExtData string
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
	_, naddr, err := c.QueryExtend(ba)
	return naddr, err
}

func NewCacheBasCli2(basip string, db *leveldb.DB) BASClient {
	return &cachedClient{basIP: basip, database: db}
}

func (c *cachedClient) QueryExtend(ba []byte) (extData string, naddr *dbSrv.NetworkAddr, err error) {
	res := &basCacheItem{}
	if err := common.GetJsonObj(c.database, ba, res); err == nil && !res.expired() {
		return res.ExtData, res.NetworkAddr, nil
	}

	extdata, ntAddr, err := QueryExtendBySrvIP(ba, c.basIP)
	if err != nil {
		return "", nil, err
	}
	if ntAddr.NTyp == dbSrv.NoItem {
		return "", nil, fmt.Errorf("no such block chain address's[%s] ip address", ba)
	}

	res = &basCacheItem{
		ExtData:     extdata,
		When:        time.Now().Add(SimpleBasCacheTime),
		NetworkAddr: ntAddr,
	}
	_ = common.SaveJsonObj(c.database, ba, res)
	return extdata, ntAddr, nil
}

func (c *cachedClient) QueryByConn(conn *network.JsonConn, ba []byte) (*dbSrv.NetworkAddr, error) {
	_, naddr, err := c.QueryExtendByConn(conn, ba)

	return naddr, err
}

func (c *cachedClient) QueryExtendByConn(conn *network.JsonConn, ba []byte) (extData string, naddr *dbSrv.NetworkAddr, err error) {
	res := &basCacheItem{}
	if err := common.GetJsonObj(c.database, ba, res); err == nil && !res.expired() {
		return res.ExtData, res.NetworkAddr, nil
	}

	extdata, ntAddr, err := QueryExtendByConn(conn, ba)
	if err != nil {
		return "", nil, err
	}
	if ntAddr.NTyp == dbSrv.NoItem {
		return "", nil, fmt.Errorf("no such block chain address's[%s] ip address", ba)
	}

	res = &basCacheItem{
		ExtData:     extdata,
		When:        time.Now().Add(SimpleBasCacheTime),
		NetworkAddr: ntAddr,
	}
	_ = common.SaveJsonObj(c.database, ba, res)
	return res.ExtData, ntAddr, nil
}

func (c *cachedClient) Register(req *dbSrv.RegRequest) error {
	return RegisterBySrvIP(req, c.basIP)
}

func (c *cachedClient) String() string {
	return fmt.Sprintf("\n----------BAS Cached Client----------"+
		"->BAS IP:%s\n"+
		"\n-------------------",
		c.basIP)
}

func (c *cachedClient) BaseAddr() string {
	addr := &net.UDPAddr{IP: net.ParseIP(c.basIP),
		Port: dbSrv.BASQueryPort}
	return addr.String()
}
