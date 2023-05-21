package provider

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"hash/fnv"

	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

// https://github.com/aws/amazon-eks-pod-identity-webhook/blob/master/hack/self-hosted/main.go

// copied from kubernetes/kubernetes#78502
func keyIDFromPublicKey(publicKey interface{}) (string, error) {
	publicKeyDERBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to serialize public key to DER format: %v", err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(publicKeyDERBytes)
	publicKeyDERHash := hasher.Sum(nil)

	keyID := base64.RawURLEncoding.EncodeToString(publicKeyDERHash)

	return keyID, nil
}

type KeyResponse struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

func readKey(content string) ([]byte, error) {
	var response []byte

	block, _ := pem.Decode([]byte(content))
	if block == nil {
		return response, errors.Errorf("Error decoding PEM content")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return response, errors.Wrapf(err, "Error parsing key content")
	}
	switch pubKey.(type) {
	case *rsa.PublicKey:
	case *ecdsa.PublicKey:
	default:
		return response, fmt.Errorf("Public key was not RSA or ECDSA, but %T", pubKey)
	}

	var alg jose.SignatureAlgorithm
	switch pubKey.(type) {
	case *rsa.PublicKey:
		alg = jose.RS256
	case *ecdsa.PublicKey:
		alg = jose.ES256
	default:
		return response, fmt.Errorf("invalid public key type %T, must be *rsa.PrivateKey or *ecdsa.PublicKey", pubKey)
	}

	kid, err := keyIDFromPublicKey(pubKey)
	if err != nil {
		return response, err
	}

	var keys []jose.JSONWebKey
	keys = append(keys, jose.JSONWebKey{
		Key:       pubKey,
		KeyID:     kid,
		Algorithm: string(alg),
		Use:       "sig",
	})

	keyResponse := KeyResponse{Keys: keys}
	return json.Marshal(keyResponse)
}

// https://stackoverflow.com/a/13582881
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
