package cmd

import (
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/BAS/pbs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type cmdService struct{
	book *dbSrv.BASTable
}

func StartCmdService( book *dbSrv.BASTable) {

	l, err := net.Listen("tcp", "127.0.0.1:50998")
	if err != nil {
		panic(err)
	}

	cmdServer := grpc.NewServer()

	pbs.RegisterCmdServiceServer(cmdServer, &cmdService{book: book})

	reflection.Register(cmdServer)
	if err := cmdServer.Serve(l); err != nil {
		panic(err)
	}
}

func DialToCmdService() pbs.CmdServiceClient {
	var address = "127.0.0.1:50998"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pbs.NewCmdServiceClient(conn)

	return client
}
