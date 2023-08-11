package token

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"strings"
)

// Error constants
var (
	ErrInvalidKey       = errors.New("key is invalid")
	ErrInvalidKeyType   = errors.New("key is of invalid type")
	ErrHashUnavailable  = errors.New("the requested hash function is unavailable")
	ErrSignatureInvalid = errors.New("signature is invalid")
)

// HS256Verify implements token verification for the SigningMethod. Returns nil if the signature is valid.
func HS256Verify(signingString, signature string, key interface{}) error {
	// Verify the key is the right type
	keyBytes, ok := key.([]byte)
	if !ok {
		return ErrInvalidKeyType
	}

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := hmac.New(crypto.SHA256.New, keyBytes)
	hasher.Write([]byte(signingString))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	// No validation errors.  Signature is good.
	return nil
}

// HS256Sign implements token signing for the SigningMethod.
// Key must be []byte
func HS256Sign(signing []byte, key interface{}) (string, error) {
	if keyBytes, ok := key.([]byte); ok {
		hasher := hmac.New(crypto.SHA256.New, keyBytes)
		hasher.Write(signing)
		return EncodeSegment(hasher.Sum(nil)), nil
	}

	return "", ErrInvalidKey
}

var DecodePaddingAllowed bool

// EncodeSegment encodes a JWT specific base64url encoding with padding stripped
// should only be used internally
func EncodeSegment(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}

// DecodeSegment decodes a JWT specific base64url encoding with padding stripped
// should only be used internally
func DecodeSegment(seg string) ([]byte, error) {
	if DecodePaddingAllowed {
		if l := len(seg) % 4; l > 0 {
			seg += strings.Repeat("=", 4-l)
		}
		return base64.URLEncoding.DecodeString(seg)
	}

	return base64.RawURLEncoding.DecodeString(seg)
}
