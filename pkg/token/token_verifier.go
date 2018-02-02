package token

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"sync"
)

// Verifier represents a Google Token verifier which caches retrieved Google's public keys.
// As these keys are regularly rotated, the user should determine when to retrieve them again
type Verifier struct {
	certs *certs

	l sync.Mutex
}

// NewVerifier creates a new Google Token verifier
func NewVerifier() *Verifier {
	return &Verifier{}
}

// VerifyIDToken verifies the google ID token.
// See also https://developers.google.com/identity/sign-in/web/backend-auth
func (v *Verifier) VerifyIDToken(ctx context.Context, jwt string, audience string) error {
	certs, err := v.getOrFetchCerts(ctx)
	if err != nil {
		return err
	}
	return v.verifyIDToken(jwt, audience, certs)
}

func (v *Verifier) verifyIDToken(jwt string, audience string, certs *certs) error {
	if jwt == "" {
		return errors.New("token is empty")
	}

	header, payload, signature, checksum, err := v.splitToken(jwt)
	if err != nil {
		return err
	}

	err = v.checkClaims(payload, audience)
	if err != nil {
		return err
	}
	err = v.checkSignature(certs, header, signature, checksum)
	if err != nil {
		return err
	}
	return nil
}

// ResetCerts resets the caches Google's public keys. They will be retrieved during next ID token verification.
func (v *Verifier) ResetCerts() {
	v.l.Lock()
	defer v.l.Unlock()

	v.certs = nil
}

func (v *Verifier) getOrFetchCerts(ctx context.Context) (*certs, error) {
	v.l.Lock()
	defer v.l.Unlock()

	if v.certs == nil {
		certs, err := getGoogleCerts(ctx)
		if err != nil {
			return nil, err
		}
		v.certs = certs
	}
	return v.certs, nil
}

func (v *Verifier) checkSignature(certs *certs, header, signature, checksum []byte) error {
	signKey, err := v.getSignKey(certs.Keys, v.getTokenKid(header))
	if err != nil {
		return err
	}

	pKey, err := v.getPublicKey(signKey)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pKey, crypto.SHA256, checksum, signature)
	if err != nil {
		return errors.New("token is not valid, rsa verification error")
	}
	return nil
}

func (v *Verifier) checkClaims(payload []byte, audience string) error {
	claimSet, err := getClaimSet(payload)
	if err != nil {
		return err
	}
	if (claimSet.Iss != "accounts.google.com") && (claimSet.Iss != "https://accounts.google.com") {
		return errors.New("token is not valid, ISS from token and certificate don't match")
	}
	if !checkTime(claimSet) {
		return errors.New("token is not valid, token is expired")
	}
	if audience != "" && audience != claimSet.Aud {
		return errors.New("token is not valid, AUD from token and audience don't match")
	}
	return nil
}

func (v *Verifier) getPublicKey(key *keys) (*rsa.PublicKey, error) {
	b, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(b)

	b, err = base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}
	e := new(big.Int).SetBytes(b)

	pub := &rsa.PublicKey{N: n, E: int(e.Uint64())}
	return pub, nil
}

func (v *Verifier) getSignKey(certsKeys []keys, tokenKid string) (*keys, error) {
	for _, v := range certsKeys {
		if v.Kid == tokenKid {
			return &v, nil
		}
	}
	return nil, errors.New("token is not valid, kid from token and certificate don't match")
}

func (v *Verifier) getTokenKid(bt []byte) string {
	var a keys
	json.Unmarshal(bt, &a)
	return a.Kid
}

func (v *Verifier) splitToken(jwt string) ([]byte, []byte, []byte, []byte, error) {
	args := strings.Split(jwt, ".")
	if len(args) != 3 {
		return nil, nil, nil, nil, errors.New("invalid token received")
	}
	header, err := base64.RawURLEncoding.DecodeString(args[0])
	if err != nil {
		return nil, nil, nil, nil, err
	}
	payload, err := base64.RawURLEncoding.DecodeString(args[1])
	if err != nil {
		return nil, nil, nil, nil, err
	}
	signature, err := base64.RawURLEncoding.DecodeString(args[2])
	if err != nil {
		return nil, nil, nil, nil, err
	}
	checksum := v.calcSum(args[0] + "." + args[1])
	return header, payload, signature, checksum, nil
}

func (v *Verifier) calcSum(str string) []byte {
	a := sha256.New()
	a.Write([]byte(str))
	return a.Sum(nil)
}
