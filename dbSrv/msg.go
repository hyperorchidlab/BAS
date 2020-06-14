package dbSrv

import (
	"encoding/json"
	"fmt"
	"github.com/hyperorchidlab/BAS/crypto"
	"net"
)

const (
	_      = iota
	NoItem = 1
	IPV4   = 2
	IPV6   = 3

	BufSize      = 1024
	BASQueryPort = 8853
	BASRegPort   = 8854
)

type BasQuery struct {
	BlockAddr []byte `json:"ba"`
}

type BasAnswer struct {
	Sig []byte `json:"signature"`
	*NetworkAddr
}

type NetworkAddr struct {
	NTyp    uint8  `json:"networkType"`
	BTyp    uint8  `json:"blockChainType"`
	NetAddr []byte `json:"networkAddr"`
}

func (na *NetworkAddr) String() string {
	return fmt.Sprintf("\n----------------------------------------------"+
		"\n network type:\t%d [1:invalid, 2:IPV4, 3:IPV6]"+
		"\n blockChain type:\t%d [1:ETH, 2:HOP]"+
		"\n network address:\t%s"+
		"\n ----------------------------------------------",
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

func Verify(typ uint8, BAddr []byte, nAddr *NetworkAddr, sig []byte) bool {
	v, ok := crypto.CurVerifier[typ]
	if !ok {
		return false
	}
	return v.Verify(BAddr, sig, nAddr)
}

func (req *RegRequest) String() string {
	return fmt.Sprintf("***********[Bas Register Reuest]**********"+
		"\n*BlockChainType:\t%d"+
		"\n*NetworkType:\t%d"+
		"\n*NetworkAddr:\t%s"+
		"\n************************************",
		req.BTyp, req.NTyp, req.NetAddr)
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
