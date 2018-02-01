package token

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// ClaimSet represents claim set
type ClaimSet struct {
	Iss           string `json:"iss"`             // email address of the client_id of the application making the access token request
	Scope         string `json:"scope,omitempty"` // space-delimited list of the permissions the application requests
	Aud           string `json:"aud"`             // descriptor of the intended target of the assertion (Optional).
	Azp           string `json:"azp"`
	Exp           int64  `json:"exp"`           // the expiration time of the assertion (seconds since Unix epoch)
	Iat           int64  `json:"iat"`           // the time the assertion was issued (seconds since Unix epoch)
	Typ           string `json:"typ,omitempty"` // token type (Optional).
	Sub           string `json:"sub,omitempty"` // Email for which the application is requesting delegated access (Optional).
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// IDToken represents JWT token
type IDToken struct {
	jwt      string
	header   string
	payload  string
	claimSet *ClaimSet
}

// NewIDToken creates a new IDToken from the jwt string
func NewIDToken(jwt string) (*IDToken, error) {
	args := strings.Split(jwt, ".")
	if len(args) != 3 {
		return nil, errors.New("invalid token received")
	}
	header, err := base64.RawURLEncoding.DecodeString(args[0])
	if err != nil {
		return nil, err
	}
	payload, err := base64.RawURLEncoding.DecodeString(args[1])
	if err != nil {
		return nil, err
	}
	claimSet, err := getClaimSet(payload)
	if err != nil {
		return nil, err
	}
	return &IDToken{
		jwt:      jwt,
		header:   string(header),
		payload:  string(payload),
		claimSet: claimSet,
	}, nil
}

func checkTime(claimSet *ClaimSet) bool {
	if (time.Now().Unix() < claimSet.Iat) || (time.Now().Unix() > claimSet.Exp) {
		return false
	}
	return true
}

func getClaimSet(payload []byte) (*ClaimSet, error) {
	claimSet := &ClaimSet{}
	err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(claimSet)
	if err != nil {
		return nil, err
	}
	return claimSet, nil
}

// CheckTime checks if the expiry time (exp) of the ID token has not passed
func (t *IDToken) CheckTime() bool {
	if (time.Now().Unix() < t.claimSet.Iat) || (time.Now().Unix() > t.claimSet.Exp) {
		return false
	}
	return true
}

// GetJWT returns raw token
func (t *IDToken) GetJWT() string {
	return t.jwt
}

// GetHeader returns header part of the token
func (t *IDToken) GetHeader() string {
	return t.header
}

// GetPayload returns payload part of the token
func (t *IDToken) GetPayload() string {
	return t.payload
}

// GetClaimSet returns parsed claim set
func (t *IDToken) GetClaimSet() ClaimSet {
	return *t.claimSet
}
