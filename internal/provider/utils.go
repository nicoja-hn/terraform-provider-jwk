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

	var alg jose.SignatureAlgorithm
	switch pubKey.(type) {
	case *rsa.PublicKey:
		alg = jose.RS256
		alg = jose.RS384
		alg = jose.RS512
		// Note: Could potentially allow for different signing algorithms
		// 		 needs bit-length checking
		// var bitLen = pubKey.N.BitLen()
	case *ecdsa.PublicKey:
		// Note: Possible alternative to go via bit-size
		// pubKey.Curve.Params().BitSize()
		// Ref: https://stackoverflow.com/a/42718174
		// var bitLen = pubKey.(*ecdsa.PublicKey).Curve.Params().BitSize

		// Canonical names: P-224, P-256, P-384, P-521
		// pubKey.Curve.Params().Name()
		switch pubKey.(*ecdsa.PublicKey).Curve.Params().Name {
		case "P-256":
			alg = jose.ES256
		case "P-384":
			alg = jose.ES384
		case "P-521":
			alg = jose.ES512
		default:
			// Note: P-224 can't be used due to too short length
			return response, fmt.Errorf("ECDSA public key should be either P256, P382, P521, but is %T", pubKey.(*ecdsa.PublicKey).Curve.Params().Name)
		}
	default:
		return response, fmt.Errorf("Public key was not RSA or ECDSA, but %T", pubKey)
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
