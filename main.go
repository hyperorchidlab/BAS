package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperorchidlab/BAS/crypto"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/BAS/querySrv"
	"github.com/hyperorchidlab/BAS/regSrv"
	"github.com/hyperorchidlab/go-miner-pool/network"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const Version = "1.0.0_gr"

var param struct {
	addr  string
	typ   uint8
	basIP string
}

var queryCmd = &cobra.Command{
	Use: "query",

	Short: "query ip address of block chain address",

	Long: `usage description`,

	Run: queryAction,
}

var rootCmd = &cobra.Command{
	Use: "BAS",

	Short: "BlockChain Address Service",

	Long: `usage description`,

	Run: mainRun,
}
var Conf struct {
	Version bool
	DBPath  string
	PidPath string
}

func BaseDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	baseDir := filepath.Join(usr.HomeDir, string(filepath.Separator), ".bass")
	return baseDir
}

func init() {
	base := BaseDir()
	defaultDB := filepath.Join(base, string(filepath.Separator), "baseBook")
	defaultPid := filepath.Join(base, string(filepath.Separator), ".pid")

	rootCmd.Flags().BoolVarP(&Conf.Version, "Version", "v", false, "show current Version")
	rootCmd.Flags().StringVarP(&Conf.DBPath, "database", "b", defaultDB, "BAS -b [DATA BASE DIR]")
	Conf.PidPath = defaultPid

	queryCmd.Flags().StringVarP(&param.basIP, "basip", "b", "", "BAS query -b [BAS IP ADDRESS]")
	queryCmd.Flags().StringVarP(&param.addr, "address", "a", "", "BAS query -a [BLOCK CHAIN ADDRESS]")
	queryCmd.Flags().Uint8VarP(&param.typ, "netType", "t", 0, "BAS query -t [1:ETH, 2:HOP]")

	rootCmd.AddCommand(queryCmd)
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainRun(_ *cobra.Command, _ []string) {
	if Conf.Version {
		fmt.Println(Version)
		return
	}

	db := dbSrv.InitTable(Conf.DBPath)
	searchSrv := querySrv.UDPSrv(db)

	saveSrv := regSrv.NewReg(db)
	done := make(chan bool, 1)

	go searchSrv.Run(done)
	go saveSrv.Serve(done)
	go waitSignal(done)
	<-done
}

func waitSignal(done chan bool) {
	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>BAS start at pid(%s)<<<<<<<<<<\n", pid)
	if err := ioutil.WriteFile(Conf.PidPath, []byte(pid), 0644); err != nil {
		fmt.Print("failed to write running pid", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sigCh

	fmt.Printf("\n>>>>>>>>>>process finished(%s)<<<<<<<<<<\n", sig)
	done <- true
}

func queryAction(_ *cobra.Command, _ []string) {

	ip := param.basIP
	if ip == "" {
		ip = "127.0.0.1"
	}

	conn, err := net.DialUDP("udp", nil,
		&net.UDPAddr{IP: net.ParseIP(ip),
			Port: dbSrv.BASQueryPort})
	if err != nil {
		panic(err)
	}

	if param.addr == "" {
		fmt.Println("please input address")
		return
	}

	jConn := &network.JsonConn{Conn: conn}

	var key []byte

	typ := param.typ
	if typ == 0 {
		//guess type
		if param.addr[0:2] == "HO" {
			typ = crypto.HOP
		}

		if param.addr[0:2] == "0x" {
			typ = crypto.BTETH
		}
	}
	if typ == crypto.BTETH {
		key = common.HexToAddress(param.addr).Bytes()
	} else if typ == crypto.HOP {
		key = []byte(param.addr)
	} else {
		panic("unknown crypt type!")
	}
	req := &dbSrv.BasQuery{
		BlockAddr: key,
	}

	if err := jConn.WriteJsonMsg(req); err != nil {
		panic(err)
	}

	res := &dbSrv.BasAnswer{}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 6))
	if err := jConn.ReadJsonMsg(res); err != nil {
		panic(err)
	}

	if res.NTyp == dbSrv.NoItem {
		fmt.Println("No such bas item")
		return
	}

	if !dbSrv.Verify(res.BTyp, key, &res.SignData, res.Sig) {
		panic("this is a polluted address")
	}

	fmt.Println(res.SignData.String())
}
