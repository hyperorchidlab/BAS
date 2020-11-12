package querySrv

import (
	"encoding/json"
	"fmt"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/go-miner-pool/common"
	"net"
)

type UDPBAS struct {
	book *dbSrv.BASTable
	srv  *net.UDPConn
}

func UDPSrv(db *dbSrv.BASTable) *UDPBAS {
	srv, err := net.ListenUDP("udp4", &net.UDPAddr{Port: dbSrv.BASQueryPort})
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

		common.NewThread(func(sig chan struct{}) {
			tb.answer(buf[:n], addr)
		}, func(err interface{}) {
			_, _ = tb.srv.WriteToUDP(dbSrv.EmptyData, addr)
		}).Start()
	}
}

func (tb *UDPBAS) answer(data []byte, from *net.UDPAddr) {
	var (
		req            = &dbSrv.BasQuery{}
		resData []byte = nil
	)

	if err := json.Unmarshal(data, req); err != nil {
		panic(err)
	}

	record := tb.book.Find(req)
	if record == nil {
		_, _ = tb.srv.WriteToUDP(dbSrv.EmptyData, from)
		return
	}
	fmt.Println(string(record.NAddr))

	answer := &dbSrv.BasAnswer{}
	answer.ExtData = record.ExtData
	answer.Sig = record.Sig
	answer.SignData.NetworkAddr = &dbSrv.NetworkAddr{
		BTyp:    record.BType,
		NTyp:    record.NType,
		NetAddr: record.NAddr,
	}
	resData, _ = json.Marshal(*answer)
	_, _ = tb.srv.WriteToUDP(resData, from)
}
