package main

import (
	"fmt"
	"os"

	"github.com/grepplabs/google-id-token/cmd/token"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().IntP("timeout", "t", 5, "Timeout in seconds")
	rootCmd.PersistentFlags().String("client-id", "", "Client ID (optional)")
}

var rootCmd = &cobra.Command{
	Use:   "google-id-token",
	Short: "Retrieve or verify google ID token",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(token.New)
	rootCmd.AddCommand(token.Get)
	rootCmd.AddCommand(token.Print)
	rootCmd.AddCommand(token.Verify)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
