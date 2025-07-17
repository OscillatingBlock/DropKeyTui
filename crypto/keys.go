package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func GenerateKeyPair() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key pair : %w", err)
	}
	return base64.StdEncoding.EncodeToString(pub), base64.StdEncoding.EncodeToString(priv), nil
}
