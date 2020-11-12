package basc

import (
	"fmt"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/go-miner-pool/network"
	"net"
)

type client struct {
	basIP string
}

func NewBasCli(basIP string) BASClient {
	c := &client{basIP: basIP}
	return c
}

func (c *client) QueryExtend(ba []byte) (extData string, naddr *dbSrv.NetworkAddr, err error) {
	return QueryExtendBySrvIP(ba, c.basIP)
}

func (c *client) Query(ba []byte) (*dbSrv.NetworkAddr, error) {
	return QueryBySrvIP(ba, c.basIP)
}

func (c *client) QueryExtendByConn(conn *network.JsonConn, ba []byte) (extData string, naddr *dbSrv.NetworkAddr, err error) {
	return QueryExtendByConn(conn, ba)
}

func (c *client) QueryByConn(conn *network.JsonConn, ba []byte) (*dbSrv.NetworkAddr, error) {
	return QueryByConn(conn, ba)
}

func (c *client) Register(req *dbSrv.RegRequest) error {
	return RegisterBySrvIP(req, c.basIP)
}

func (c *client) String() string {
	return fmt.Sprintf("\n----------Bas Simple client----------"+
		"\n->BAS IP:%s"+
		"\n---------------------------------------",
		c.basIP)
}

func (c *client) BaseAddr() string {
	addr := &net.UDPAddr{IP: net.ParseIP(c.basIP),
		Port: dbSrv.BASQueryPort}
	return addr.String()
}
