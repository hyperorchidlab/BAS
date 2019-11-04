package regSrv

import (
	"fmt"
	"github.com/hyperorchid/go-miner-pool/network"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"net"
)

type Register struct {
	db  *dbSrv.BASTable
	srv net.Listener
}

func NewReg(db *dbSrv.BASTable) *Register {
	srv, err := net.ListenTCP("tcp", &net.TCPAddr{Port: dbSrv.BASRegPort})
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
		go r.Register(conn)
	}
}

func (r *Register) Register(conn net.Conn) {

	defer conn.Close()
	jsonCon := network.JsonConn{Conn: conn}
	req := &dbSrv.RegRequest{}
	if err := jsonCon.ReadJsonMsg(req); err != nil {
		fmt.Println(err)
		return
	}

	if !req.Verify() {
		fmt.Println("verify failed->:", req.String())
		conn.Write(dbSrv.NewRegResponse(false, 1, "verify failed"))
		return
	}

	if err := r.db.Save(req); err != nil {
		e := fmt.Errorf("save data base err:%s", err)
		fmt.Println(e)
		conn.Write(dbSrv.NewRegResponse(false, 2, e.Error()))
		return
	}

	conn.Write(dbSrv.NewRegResponse(true, 0, ""))
}
