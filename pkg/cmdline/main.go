package main

import (

	"github.com/spf13/cobra"
)


func main() {
	cmd := &cobra.Command{
		Use: "middlewarectl",
		Short: "Command line tool to interact with middleware",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}




	cmd.Execute()
}
