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
	Verify.Flags().StringP("id-token", "i", "", "ID Token is to verify")
	Verify.Flags().StringP("audience", "a", "", "Audience that this ID Token is intended for (optional)")

}

// Verify represents the verify command
var Verify = &cobra.Command{
	Use:   "verify",
	Short: "verify the token",
	Run: func(cmd *cobra.Command, args []string) {
		timeout, _ := cmd.Flags().GetInt("timeout")
		idToken, _ := cmd.Flags().GetString("id-token")
		audience, _ := cmd.Flags().GetString("audience")

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		verifier := tokenlib.NewVerifier()
		err := verifier.VerifyIDToken(ctx, idToken, audience)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
