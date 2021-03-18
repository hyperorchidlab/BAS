package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperorchidlab/BAS/dbSrv"

	"github.com/hyperorchidlab/BAS/pbs"
)

func toMinerString(record *dbSrv.Record) string {
	msg := ""
	msg += fmt.Sprintf("address: %s \r\n",string(record.BAddr))
	msg += fmt.Sprintf("ip addr: %s \r\n",string(record.NAddr))
	msg += fmt.Sprintf("extData: %s\r\n",record.ExtData)
	return msg
}

func (cs *cmdService) ShowAllMiners(ctx context.Context, r *pbs.EmptyRequest) (*pbs.CommonResponse, error) {

	msg := ""

	miners:=cs.book.FindAll(21,60)
	for i:=0;i<len(miners);i++{

		msg += toMinerString(miners[i])
		msg += "================================================\r\n"
	}

	if msg == ""{
		msg = "no miners"
	}

	return &pbs.CommonResponse{
		Msg: msg,
	}, nil

}

func toPoolsString(record *dbSrv.Record) string  {
	msg := ""
	msg += fmt.Sprintf("address: %s \r\n",common.BytesToAddress(record.BAddr).String())
	msg += fmt.Sprintf("ip addr: %s \r\n",string(record.NAddr))
	msg += fmt.Sprintf("extData: %s\r\n",record.ExtData)

	return msg
}

func (cs *cmdService) ShowAllPools(ctx context.Context, r *pbs.EmptyRequest) (*pbs.CommonResponse, error) {
	msg := ""

	pools:=cs.book.FindAll(0,21)
	for i:=0;i<len(pools);i++{

		msg += toPoolsString(pools[i])
		msg += "================================================\r\n"
	}

	if msg == ""{
		msg = "no pools"
	}

	return &pbs.CommonResponse{
		Msg: msg,
	}, nil

}
