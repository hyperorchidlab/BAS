package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Register struct {
	db  *BASTable
	srv net.Listener
}

func newReg(db *BASTable) *Register {
	srv, err := net.ListenTCP("tcp", &net.TCPAddr{Port: DNSSPort})
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

	buf := make([]byte, BufSize)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	req := &RegRequest{}
	if err := json.Unmarshal(buf[:n], req); err != nil {
		fmt.Println(err)
		return
	}

	if !req.Verify() {
		fmt.Println("verify failed->:", req.String())
		conn.Write(newRegResponse(false, 1, "verify failed"))
		return
	}

	if err := r.db.save(req); err != nil {
		e := fmt.Errorf("save data base err:%s", err)
		fmt.Println(e)
		conn.Write(newRegResponse(false, 2, e.Error()))
		return
	}

	conn.Write(newRegResponse(true, 0, ""))
}
