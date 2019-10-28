package dbSrv

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperorchid/go-miner-pool/account"
)

const (
	_ = iota
	BTECDSA
	BTEd25519

	_ = iota
	NoItem
	IPV4
	IPV6

	BufSize  = 1024
	DNSGPort = 53
	DNSSPort = 54
)

type BlockChainAddr struct {
	BlockAddr []byte `json:"ba"`
}

type NetworkAddr struct {
	NTyp    uint8  `json:"nt"`
	BTyp    uint8  `json:"bt"`
	NetAddr []byte `json:"na"`
}

var Empty = &NetworkAddr{
	NTyp: NoItem,
}

var EmptyData, _ = json.Marshal(Empty)

type RegRequest struct {
	Sig          []byte `json:"sig"`
	BlockAddr    []byte `json:"ba"`
	*NetworkAddr `json:"data"`
}

type RegResponse struct {
	Success bool   `json:"success"`
	ENO     uint8  `json:"eno"`
	MSG     string `json:"msg"`
}

func NewRegResponse(success bool, eno uint8, msg string) []byte {
	var err = &RegResponse{
		Success: success,
		ENO:     eno,
		MSG:     msg,
	}
	b, _ := json.Marshal(err)
	return b
}

func (req *RegRequest) Verify() bool {

	switch req.BTyp {
	case BTECDSA:
		addr := common.BytesToAddress(req.BlockAddr)
		return account.VerifyJsonSig(addr, req.Sig, req.NetworkAddr)
	case BTEd25519:
		id := account.ID(string(req.BlockAddr))
		return account.VerifySubSig(id, req.Sig, req.NetworkAddr)
	case NoItem:
		return false
	default:
		return false
	}
}

func (req *RegRequest) String() string {
	return fmt.Sprintf("*********************"+
		"\n*BlockChainType:\t%d"+
		"\n*NetworkType:\t%d"+
		"\n*NetworkAddr:\t%s"+
		"\n*********************\n", req.BTyp, req.NTyp, req.NetAddr)
}
