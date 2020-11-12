package iosLib

import "C"
import basc "github.com/hyperorchidlab/BAS/client"

type iosClient struct {
	basIP string
}

var cliInst *iosClient = nil

func initBas(ip string) {
	cliInst = &iosClient{
		basIP: ip,
	}
}

func Query(ba []byte) (string, error) {

	ret, err := basc.QueryBySrvIP(ba, cliInst.basIP)
	if err != nil {
		return "", err
	}

	return string(ret.NetAddr), nil
}

func QueryExtend(ba []byte) (string, string, error) {
	extd, ret, err := basc.QueryExtendBySrvIP(ba, cliInst.basIP)

	return extd, string(ret.NetAddr), err
}
