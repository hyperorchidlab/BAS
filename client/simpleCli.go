package basc

import (
	"github.com/hyperorchidlab/BAS/dbSrv"
)

type client struct {
	basIP string
}

func NewBasCli(basIP string) BASClient {
	c := &client{basIP: basIP}
	return c
}

func (c *client) Query(ba []byte) (*dbSrv.NetworkAddr, error) {
	return QueryBySrvIP(ba, c.basIP)
}

func (c *client) Register(req *dbSrv.RegRequest) error {
	return RegisterBySrvIP(req, c.basIP)
}
