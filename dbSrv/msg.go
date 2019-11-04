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

func ConvertIP(ip string) (uint8, []byte, error) {

	netIP := net.ParseIP(ip)
	if netIP == nil {
		return NoItem, nil, fmt.Errorf("parse ip[%s] failed", ip)
	}

	fmt.Println(netIP, len(netIP))

	if len(netIP) == net.IPv4len {
		return IPV4, netIP[:], nil
	} else if len(netIP) == net.IPv6len {
		return IPV6, netIP[:], nil
	}

	return NoItem, nil, fmt.Errorf("invalid ip string[%s]", ip)
}
