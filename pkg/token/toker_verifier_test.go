package token

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testGoogleCert = `
{
 "keys": [
  {
   "kty": "RSA",
   "alg": "RS256",
   "use": "sig",
   "kid": "ac7ebbdff9e77669785f4c530fe2d4a6408bc98d",
   "n": "1L2jYqXcdvdxtY10zT3PTZyTxG_gIScRcSheHYsuRMfdsh40xl_fBhpSfAIIlCwgyMz8_A03SueTw_jQE0z0JGezXi-1WiGpEJBY7Dmm78t5W6U7J-ktT3N7PtJgzS0gza8XVarGSkPc_J78PiRKmI4HJmw7d7nhHAEPJpH4RiMXIsOm6-HKEb_j_BdZJbiQRof67L8eAvIKNWi-8NsWKDK4UFE-8uddFo_E1_vjDdAfr99Rp45802_oN4sJNlgjU6s4vequyfFq-LQc5wXlt68ja3r_0TIaMCTFqO8ui-t1JGm1rS1qfCYLfxd3ZkhtVVDf2TtLrfC1rOTSzLxMUQ",
   "e": "AQAB"
  },
  {
   "kty": "RSA",
   "alg": "RS256",
   "use": "sig",
   "kid": "978ca4118bf1883b316bbca6ce9044d9977f2027",
   "n": "qpe-lPi7HVP8_SRqodC19iWDcYJ-5-wZbBxxxgszoPbphgN8cUdcwOYuPoTT7BmDvezKhHq_JPjqxkJWO5_GESPw_ijMnXE3PddO1nNmWIBOxUBSE34LUf_GDsyXL6DmiiPsJtSdPgW4BzxkSf4VU-obP-K1BEyxmWwUJdUhNpUM7aj7aC-pCZJZyNF_OBjY5mq1lKn9kJvuy_EiSRvyCySR149lJW86K7VnbLGguu1pOo4s2JXf7nWGeccNydeJznY5FOi4tAxTiaGpM2gzXtS7gUDKgEKufE5V1fq2MtY-pYypRObZsesRit9CA3fQHQ5hrHhA4_uwLjhsVK0Z0w",
   "e": "AQAB"
  },
  {
   "kty": "RSA",
   "alg": "RS256",
   "use": "sig",
   "kid": "3405d0ec4edf60539acf73be64604d49a097189a",
   "n": "vBNfb9rmZLTwVpjoeT9lsLvzwl5rAVWGius9n2AFdibXlTaA_orGvSXL7l7SYLFcoxVGwNGrXDlAqwvpytyvOyRKcIepjbgRwADOAMbn4B4iQFlwI90dS_xqfGa6Ye6B-M6B802m0M43MJmZeEP9b81s2ExVPouE_zQz6-Pu1ZpABB2X7NaReVOzJAdOboRMQbZVh-X-HnYbM9PTWV-4fecQqE9sD_Qi8NSXiN1aC2n8DaIMjeEkDH5PJCPO6wmDUkolb2dIb3jryr19dGV0_Z-jRgzl8vdNw-u0o3mm0X8nj3cJSajjVEsTdbH35SdyFeu9Ob5G0oxwbPIBUr-2jw",
   "e": "AQAB"
  }
 ]
}
`
var testHeader = `
{
 "alg": "RS256",
 "kid": "978ca4118bf1883b316bbca6ce9044d9977f2027"
}`

// exp 2114380800 - 2037-01-01T00:00:00+00:00
// iat 1516304351 - 2018-01-18T19:39:11+00:00
var testClaims = `
{
 "azp": "4711.apps.googleusercontent.com",
 "aud": "4711.apps.googleusercontent.com",
 "sub": "100004711",
 "hd": "grepplabs.com",
 "email": "info@grepplabs.com",
 "email_verified": true,
 "exp": 2114380800,
 "iss": "accounts.google.com",
 "iat": 1516304351
}
`

func TestVerifyNotGoogleKey(t *testing.T) {
	a := assert.New(t)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	a.Nil(err)

	testToken, err := encodeTestToken(testHeader, testClaims, privateKey)
	a.Nil(err)

	certs, err := toCerts([]byte(testGoogleCert))
	a.Nil(err)

	err = NewVerifier().verifyIDToken(testToken, "", certs)
	a.EqualError(err, "token is not valid, rsa verification error")
}

func TestVerifyNoKid(t *testing.T) {
	a := assert.New(t)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	a.Nil(err)

	testToken, err := encodeTestToken(testHeader, testClaims, privateKey)
	a.Nil(err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = NewVerifier().VerifyIDToken(ctx, testToken, "")
	a.NotNil(err)
	a.EqualError(err, "token is not valid, kid from token and certificate don't match")
}

func encodeTestToken(headerJSON string, claimsJSON string, key *rsa.PrivateKey) (string, error) {
	sg := func(data []byte) (sig []byte, err error) {
		h := sha256.New()
		h.Write(data)
		return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, h.Sum(nil))
	}
	ss := fmt.Sprintf("%s.%s", base64.RawURLEncoding.EncodeToString([]byte(headerJSON)), base64.RawURLEncoding.EncodeToString([]byte(claimsJSON)))
	sig, err := sg([]byte(ss))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", ss, base64.RawURLEncoding.EncodeToString(sig)), nil
}
