package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	crusadoCmd = &cobra.Command{
		Use:   "crusado",
		Short: "TODO",
		Long:  `TODO`,
		Run:   Crusado,
	}
)

func Execute() error {
	return crusadoCmd.Execute()
}

func Crusado(cmd *cobra.Command, args []string) {
	// TODO make this nice
	fmt.Println("crusado root command - run with --help to get an overview!")
}
