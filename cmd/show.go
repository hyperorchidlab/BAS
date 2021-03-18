package cmd

import (
	"context"
	"fmt"

	"github.com/hyperorchidlab/BAS/pbs"

	"github.com/spf13/cobra"
)

var ShowMinerCmd = &cobra.Command{
	Use: "allminer",

	Short: "show all miners",

	Long: `usage description`,

	Run: showAllMiners,
}

var ShowPoolCmd = &cobra.Command{
	Use: "allpool",

	Short: "show all pools",

	Long: `usage description`,

	Run: showAllPools,
}


func showAllMiners(_ *cobra.Command, _ []string)  {
	c := DialToCmdService()
	msg, err := c.ShowAllMiners(context.TODO(), &pbs.EmptyRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(msg.Msg)
}


func showAllPools(_ *cobra.Command, _ []string)  {
	c := DialToCmdService()
	msg, err := c.ShowAllPools(context.TODO(), &pbs.EmptyRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(msg.Msg)
}
