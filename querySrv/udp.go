package querySrv

import (
	"encoding/json"
	"fmt"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"net"
)

type UDPBAS struct {
	book *dbSrv.BASTable
	srv  *net.UDPConn
}

func UDPSrv(db *dbSrv.BASTable) *UDPBAS {
	srv, err := net.ListenUDP("udp", &net.UDPAddr{Port: dbSrv.BASQueryPort})
	if err != nil {
		panic(err)
	}
	return &UDPBAS{
		book: db,
		srv:  srv,
	}
}

func (tb *UDPBAS) Run(done chan bool) {
	defer tb.srv.Close()

	for {
		buf := make([]byte, dbSrv.BufSize)
		n, addr, err := tb.srv.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			done <- false
			return
		}
		fmt.Println(addr)
		go tb.answer(buf[:n], addr)
	}
}

func (tb *UDPBAS) answer(data []byte, from *net.UDPAddr) {
	var (
		req            = &dbSrv.BlockChainAddr{}
		resData []byte = nil
	)

	if err := json.Unmarshal(data, req); err != nil {
		fmt.Println(err)
		return
	}

	record := tb.book.Find(req)

	if record == nil {
		resData = dbSrv.EmptyData
	} else {
		resData, _ = json.Marshal(&dbSrv.NetworkAddr{
			BTyp:    record.BType,
			NTyp:    record.NType,
			NetAddr: record.NAddr,
		})
	}

	if _, err := tb.srv.WriteToUDP(resData, from); err != nil {
		fmt.Println(err)
	}
}
