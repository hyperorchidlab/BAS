package main

import (
	"fmt"
	"github.com/hyperorchidlab/BAS/dbSrv"
	"github.com/hyperorchidlab/BAS/querySrv"
	"github.com/hyperorchidlab/BAS/regSrv"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const Version = "0.1"

var rootCmd = &cobra.Command{
	Use: "minerPool",

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

	baseDir := filepath.Join(usr.HomeDir, string(filepath.Separator), ".bas")
	return baseDir
}

func init() {
	base := BaseDir()
	defaultDB := filepath.Join(base, string(filepath.Separator), "baseBook")
	defaultPid := filepath.Join(base, string(filepath.Separator), ".pid")

	rootCmd.Flags().BoolVarP(&Conf.Version, "Version", "v", false, "show current Version")
	rootCmd.Flags().StringVarP(&Conf.DBPath, "database", "b", defaultDB, "BAS -b [DATA BASE DIR]")
	Conf.PidPath = defaultPid
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
