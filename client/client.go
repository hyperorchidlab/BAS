package basc

import (
	"fmt"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/go-miner-pool/network"
	"net"
	"time"
)

type BASClient interface {
	QueryExtendByConn(conn *network.JsonConn,ba []byte) (extData string, naddr *dbSrv.NetworkAddr,err error)
	QueryByConn(conn *network.JsonConn, ba []byte) (*dbSrv.NetworkAddr, error)
	QueryExtend(ba []byte) (extData string, naddr *dbSrv.NetworkAddr,err error)
	Query([]byte) (*dbSrv.NetworkAddr, error)
	Register(*dbSrv.RegRequest) error
	String() string
	BaseAddr() string
}

func RegisterBySrvIP(req *dbSrv.RegRequest, ip string) error {
	addr := &net.TCPAddr{IP: net.ParseIP(ip),
		Port: dbSrv.BASRegPort}
	nc, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}
	defer nc.Close()

	conn:=&network.JsonConn{Conn:&network.LVConn{Conn:nc}}

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

func QueryExtendBySrvIP(ba []byte,ip string) (extData string, naddr *dbSrv.NetworkAddr,err error)  {
	addr := &net.UDPAddr{IP: net.ParseIP(ip),
		Port: dbSrv.BASQueryPort}
	conn, err := network.DialJson("udp", addr.String())
	if err != nil {
		return "",nil, err
	}

	defer conn.Close()

	req := &dbSrv.BasQuery{
		BlockAddr: ba,
	}

	if err := conn.WriteJsonMsg(req); err != nil {
		return "",nil, err
	}

	res := &dbSrv.BasAnswer{}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err := conn.ReadJsonMsg(res); err != nil {
		return "",nil, err
	}

	if res.NTyp == dbSrv.NoItem {
		return "",nil, fmt.Errorf("no such BAS item")
	}

	if !dbSrv.Verify(res.BTyp, ba, &res.SignData, res.Sig) {
		return "",nil, fmt.Errorf("this is a polluted address:\n%s", res.String())
	}
	return res.ExtData,res.NetworkAddr, nil
}

func QueryBySrvIP(ba []byte, ip string) (*dbSrv.NetworkAddr, error) {
	_,nddr,err:=QueryExtendBySrvIP(ba,ip)
	return nddr,err
}

func QueryExtendByConn(conn *network.JsonConn,ba []byte) (extData string, naddr *dbSrv.NetworkAddr,err error)  {
	defer conn.Close()

	req := &dbSrv.BasQuery{
		BlockAddr: ba,
	}

	if err := conn.WriteJsonMsg(req); err != nil {
		return "",nil, err
	}

	res := &dbSrv.BasAnswer{}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err := conn.ReadJsonMsg(res); err != nil {
		return "",nil, err
	}

	if res.NTyp == dbSrv.NoItem {
		return "",nil, fmt.Errorf("no such BAS item")
	}

	if !dbSrv.Verify(res.BTyp, ba, &res.SignData, res.Sig) {
		return "",nil, fmt.Errorf("this is a polluted address:\n%s", res.String())
	}
	return res.ExtData,res.NetworkAddr, nil
}

func QueryByConn(conn *network.JsonConn, ba []byte) (*dbSrv.NetworkAddr, error) {
	_,naddr,err:=QueryExtendByConn(conn,ba)

	return naddr,err
}
