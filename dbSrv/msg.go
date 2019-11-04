package dbSrv

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperorchid/go-miner-pool/account"
	"net"
)

const (
	_         = iota
	BTECDSA   = 1
	BTEd25519 = 2

	_      = iota
	NoItem = 1
	IPV4   = 2
	IPV6   = 3

	BufSize      = 1024
	BASQueryPort = 53
	BASRegPort   = 54
)

type BlockChainAddr struct {
	BlockAddr []byte `json:"ba"`
}

type NetworkAddr struct {
	NTyp    uint8  `json:"networkType"`
	BTyp    uint8  `json:"blockChainType"`
	NetAddr []byte `json:"networkAddr"`
}

func (na *NetworkAddr) String() string {
	return fmt.Sprintf("\n------------------\n"+
		"\n network type:%d"+
		"\n blockChain type:%d"+
		"\n network address:%s"+
		"\n------------------\n",
		na.NTyp,
		na.BTyp,
		string(na.NetAddr))
}

var Empty = &NetworkAddr{
	NTyp: NoItem,
}

var EmptyData, _ = json.Marshal(Empty)

type RegRequest struct {
	Sig          []byte `json:"sig"`
	BlockAddr    []byte `json:"blockChainAddr"`
	*NetworkAddr `json:"data"`
}

type RegResponse struct {
	Success bool   `json:"success"`
	ENO     uint8  `json:"eno"`
	MSG     string `json:"msg"`
}

func NewRegResponse(success bool, eno uint8, msg string) *RegResponse {
	return &RegResponse{
		Success: success,
		ENO:     eno,
		MSG:     msg,
	}
}

func (req *RegRequest) Verify() bool {

	switch req.BTyp {
	case BTECDSA:
		addr := common.BytesToAddress(req.BlockAddr)
		return account.VerifyJsonSig(addr, req.Sig, req.NetworkAddr)
	case BTEd25519:
		id := account.ID(string(req.BlockAddr))
		return account.VerifySubSig(id, req.Sig, req.NetworkAddr)
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

func CheckIPType(ip string) (uint8, error) {

	netIP := net.ParseIP(ip)
	if netIP == nil {
		return NoItem, fmt.Errorf("parse ip[%s] failed", ip)
	}
	fmt.Println(netIP, len(netIP))
	if netIP.To4() != nil {
		return IPV4, nil
	} else if netIP.To16() != nil {
		return IPV6, nil
	}
	return NoItem, fmt.Errorf("invalid ip string[%s]", ip)
}
