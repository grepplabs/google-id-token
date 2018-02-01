package token

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"hash/fnv"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

// Retriever is a Google ID Token retriever.
type Retriever struct {
	clientID string
}

// NewRetriever creates a new Google ID Token retriever
func NewRetriever(clientID string) *Retriever {
	return &Retriever{
		clientID: clientID,
	}
}

func (r *Retriever) newIDToken(ctx context.Context) (*IDToken, error) {
	tokenSource, err := google.DefaultTokenSource(ctx, oauth2.UserinfoEmailScope)
	if err != nil {
		return nil, err
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	if token.Extra("id_token") == nil {
		return nil, fmt.Errorf("id_token not present")
	}
	idToken := token.Extra("id_token").(string)
	return NewIDToken(idToken)
}

func (r *Retriever) cachedIDToken() (*IDToken, error) {
	filename := r.getTokenCacheFile()
	idToken, err := r.getTokenFromFile(filename)
	if err == nil && idToken != "" {
		return NewIDToken(idToken)
	}
	return nil, err
}

func (r *Retriever) cacheIDToken(idToken *IDToken) error {
	filename := r.getTokenCacheFile()
	return r.saveTokenToFile(filename, idToken.jwt)
}

// GetNewIDToken retrieves a new ID token and stores it in the cache
func (r *Retriever) GetNewIDToken(ctx context.Context) (*IDToken, error) {
	result, err := r.newIDToken(ctx)
	if err == nil {
		r.cacheIDToken(result)
	}
	return result, err
}

// GetIDToken retrieves a cached token. If the ID token is not found or is expired a new one is retrieved and cached
func (r *Retriever) GetIDToken(ctx context.Context) (*IDToken, error) {
	result, err := r.cachedIDToken()
	if err == nil && result != nil && result.CheckTime() {
		return result, nil
	}
	return r.GetNewIDToken(ctx)
}

func (r *Retriever) getTokenCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return "."
}

func (r *Retriever) getTokenCacheFile() string {
	hash := fnv.New32a()
	hash.Write([]byte(r.clientID))
	fn := fmt.Sprintf("grepplabs-google-id-token-%v", hash.Sum32())
	return filepath.Join(r.getTokenCacheDir(), url.QueryEscape(fn))
}

func (r *Retriever) getTokenFromFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func (r *Retriever) saveTokenToFile(filename string, idToken string) error {
	data := []byte(idToken)
	return ioutil.WriteFile(filename, data, 0644)
}
