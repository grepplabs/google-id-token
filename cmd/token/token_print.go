package token

import (
	"context"
	"fmt"
	tokenlib "github.com/grepplabs/google-id-token/pkg/token"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	Print.Flags().StringP("id-token", "i", "", "ID Token is to print (optional)")
	Print.Flags().Bool("header", false, "Print the header")
}

// Print represents the print command
var Print = &cobra.Command{
	Use:   "print",
	Short: "print the decoded token",
	Run: func(cmd *cobra.Command, args []string) {
		timeout, _ := cmd.Flags().GetInt("timeout")
		header, _ := cmd.Flags().GetBool("header")
		jwt, _ := cmd.Flags().GetString("id-token")
		clientID, _ := cmd.Flags().GetString("client-id")

		var idToken *tokenlib.IDToken
		var err error
		if jwt == "" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			defer cancel()

			retriever := tokenlib.NewRetriever(clientID)
			idToken, err = retriever.GetIDToken(ctx)
		} else {
			idToken, err = tokenlib.NewIDToken(jwt)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if header {
			fmt.Print(idToken.GetHeader())
		} else {
			fmt.Print(idToken.GetPayload())
		}
	},
}
