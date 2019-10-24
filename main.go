package main
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use: "minerPool",

	Short: "BlockChain Address Service",

	Long: `usage description`,

	Run: mainRun,
}
var param struct {
	version  bool
}

func init() {

	rootCmd.Flags().BoolVarP(&param.version, "version",
		"v", false, "show current version")

	rootCmd.Flags().BoolVarP(&common.SysDebugMode, "debug",
		"d", false, "run in debug model")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainRun(_ *cobra.Command, _ []string) {
	if param.version {
		fmt.Println(common.CurrentVersion)
		return
	}
	

	done := make(chan bool, 1)
	go waitSignal(done)
	<-done
}

func waitSignal(done chan bool) {
	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>BAS start at pid(%s)<<<<<<<<<<\n", pid)
	if err := ioutil.WriteFile(".pid", []byte(pid), 0644); err != nil {
		fmt.Print("failed to write running pid", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sigCh

	core.PoolInstance().Finish()
	fmt.Printf("\n>>>>>>>>>>process finished(%s)<<<<<<<<<<\n", sig)

	done <- true
}
