package basc

import (
	"github.com/hyperorchidlab/BAS/dbSrv"
	"net"
)

type client struct {
	basAddr string
}

func NewBasCli(basIP string) BASClient {
	addr := &net.UDPAddr{IP: net.ParseIP(basIP),
		Port: dbSrv.BASQueryPort}
	c := &client{basAddr: addr.String()}
	return c
}

func (c *client) Query(ba []byte) (*dbSrv.NetworkAddr, error) {
	return QueryBySrvIP(ba, c.basAddr)
}

func (c *client) Register(req *dbSrv.RegRequest) error {
	return RegisterBySrvIP(req, c.basAddr)
}
