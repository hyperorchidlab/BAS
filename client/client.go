package basc

import (
	"fmt"
	"github.com/hyperorchid/go-miner-pool/network"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"time"
)

type BASClient interface {
	Query([]byte) (*dbSrv.NetworkAddr, error)
	Register(*dbSrv.RegRequest) error
}

func RegisterBySrvIP(req *dbSrv.RegRequest, srvAddr string) error {
	conn, err := network.DialJson("tcp", srvAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.WriteJsonMsg(req); err != nil {
		return err
	}

	res := &dbSrv.RegResponse{}
	if err := conn.ReadJsonMsg(res); err != nil {
		return err
	}

	if res.Success {
		return nil
	}
	return fmt.Errorf("reg failed because errno:[%d] msg:[%s]", res.ENO, res.MSG)
}

func QueryBySrvIP(ba []byte, srvAddr string) (*dbSrv.NetworkAddr, error) {
	conn, err := network.DialJson("udp", srvAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &dbSrv.BasQuery{
		BlockAddr: ba,
	}

	if err := conn.WriteJsonMsg(req); err != nil {
		return nil, err
	}

	res := &dbSrv.BasAnswer{}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err := conn.ReadJsonMsg(res); err != nil {
		return nil, err
	}
	if !dbSrv.Verify(res.BTyp, ba, res.NetworkAddr, res.Sig) {
		return nil, fmt.Errorf("this is a polluted address")
	}
	return res.NetworkAddr, nil
}
