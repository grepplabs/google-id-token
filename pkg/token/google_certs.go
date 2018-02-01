package token

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	googleOpenidConfigurationURL = "https://accounts.google.com/.well-known/openid-configuration"
)

type certs struct {
	Keys []keys `json:"keys"`
}

// https://tools.ietf.org/html/rfc7517#appendix-A
type keys struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// https://developers.google.com/identity/protocols/OpenIDConnect#discovery
type openidConfiguration struct {
	Issuer      string `json:"issuer"`
	AuthURL     string `json:"authorization_endpoint"`
	TokenURL    string `json:"token_endpoint"`
	JWKSURL     string `json:"jwks_uri"`
	UserInfoURL string `json:"userinfo_endpoint"`
}

func getGoogleCerts(ctx context.Context) (*certs, error) {
	bConfig, err := getFromURL(ctx, googleOpenidConfigurationURL)
	if err != nil {
		return nil, err
	}
	openidConfiguration, err := toOpenIDConfiguration(bConfig)
	if err != nil {
		return nil, err
	}
	if openidConfiguration.JWKSURL == "" {
		return nil, fmt.Errorf("jwks_uri from %s is empty", googleOpenidConfigurationURL)
	}
	bCerts, err := getFromURL(ctx, openidConfiguration.JWKSURL)
	if err != nil {
		return nil, err
	}
	return toCerts(bCerts)
}

func getFromURL(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	res, err := ctxhttp.Get(ctx, client, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get from %s failed: %s %s", url, res.Status, strings.TrimSpace(string(body)))
	}

	return body, nil
}

func toCerts(bt []byte) (*certs, error) {
	var certs *certs
	err := json.Unmarshal(bt, &certs)
	if err != nil {
		return nil, err
	}
	return certs, nil
}

func toOpenIDConfiguration(bt []byte) (*openidConfiguration, error) {
	var config *openidConfiguration
	err := json.Unmarshal(bt, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
