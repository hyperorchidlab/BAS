package cmd

import (
	"context"
	"encoding/json"

	"github.com/hyperorchidlab/BAS/pbs"
)

func (cs *cmdService) ShowAllMiners(ctx context.Context, r *pbs.EmptyRequest) (*pbs.CommonResponse, error) {

	msg := "no miners"

	miners:=cs.book.FindAll(21,60)
	for i:=0;i<len(miners);i++{
		j,_:=json.MarshalIndent(miners[i]," ","\t")
		msg += string(j)
		msg += "================================================"
	}

	return &pbs.CommonResponse{
		Msg: msg,
	}, nil

}

func (cs *cmdService) ShowAllPools(ctx context.Context, r *pbs.EmptyRequest) (*pbs.CommonResponse, error) {
	msg := "no pools"

	pools:=cs.book.FindAll(0,21)
	for i:=0;i<len(pools);i++{
		j,_:=json.MarshalIndent(pools[i]," ","\t")
		msg += string(j)
		msg += "================================================"
	}

	return &pbs.CommonResponse{
		Msg: msg,
	}, nil

}
