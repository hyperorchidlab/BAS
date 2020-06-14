package regSrv

import (
	"fmt"
	"github.com/hyperorchid/go-miner-pool/common"
	"github.com/hyperorchid/go-miner-pool/network"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"net"
)

type Register struct {
	db  *dbSrv.BASTable
	srv net.Listener
}

func NewReg(db *dbSrv.BASTable) *Register {
	srv, err := net.ListenTCP("tcp4", &net.TCPAddr{Port: dbSrv.BASRegPort})
	if err != nil {
		panic(err)
	}

	return &Register{
		db:  db,
		srv: srv,
	}
}

func (r *Register) Serve(done chan bool) {

	defer r.srv.Close()
	for {
		conn, err := r.srv.Accept()
		if err != nil {
			fmt.Println(err)
			done <- false
			return
		}

		common.NewThread(func(sig chan struct{}) {
			r.Register(&network.JsonConn{Conn: conn})
		}, func(err interface{}) {
			conn.Close()
		}).Start()
	}
}

func (r *Register) Register(jsonCon *network.JsonConn) {

	req := &dbSrv.RegRequest{}
	if err := jsonCon.ReadJsonMsg(req); err != nil {
		fmt.Println(err)
		return
	}

	if !dbSrv.Verify(req.BTyp, req.BlockAddr, req.NetworkAddr, req.Sig) {
		fmt.Println("verify failed->:", req.String())
		_ = jsonCon.WriteJsonMsg(dbSrv.NewRegResponse(false, 1, "verify failed"))
		return
	}

	fmt.Println(string(req.BlockAddr), len(req.BlockAddr))
	fmt.Println(req.String())

	if err := r.db.Save(req); err != nil {
		e := fmt.Errorf("save data base err:%s", err)
		fmt.Println(e)
		_ = jsonCon.WriteJsonMsg(dbSrv.NewRegResponse(false, 2, e.Error()))
		return
	}

	_ = jsonCon.WriteJsonMsg(dbSrv.NewRegResponse(true, 0, ""))
}
