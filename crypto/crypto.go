package crypto

import (
	"crypto/ed25519"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperorchid/go-miner-pool/account"
)

const (
	_     = iota
	BTETH = 1
	HOP   = 2
)

type Verifier interface {
	Verify(pub, sig []byte, msg interface{}) bool
}

var CurVerifier = map[uint8]Verifier{
	BTETH: &ETHSigner{},
	HOP:   &HOPSigner{},
}

type ETHSigner struct {
}

func (eth *ETHSigner) Verify(pub, sig []byte, msg interface{}) bool {
	data, err := json.Marshal(msg)
	if err != nil {
		return false
	}
	hash := crypto.Keccak256(data)
	signer, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false
	}
	return common.BytesToAddress(pub) == crypto.PubkeyToAddress(*signer)
}

type HOPSigner struct {
}

func (ed *HOPSigner) Verify(pubData, sig []byte, v interface{}) bool {

	id, err := account.ConvertToID(string(pubData))
	if err != nil {
		return false
	}

	data, err := json.Marshal(v)
	if err != nil {
		return false
	}

	return ed25519.Verify(id.ToPubKey(), data, sig)
}
