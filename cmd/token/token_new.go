package token

import (
	"context"
	"fmt"
	tokenlib "github.com/grepplabs/google-id-token/pkg/token"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// New represents the new command
var New = &cobra.Command{
	Use:   "new",
	Short: "get a new token",
	Run: func(cmd *cobra.Command, args []string) {
		timeout, _ := cmd.Flags().GetInt("timeout")
		clientID, _ := cmd.Flags().GetString("client-id")

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		retriever := tokenlib.NewRetriever(clientID)

		idToken, err := retriever.GetNewIDToken(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(idToken.GetJWT())
	},
}
