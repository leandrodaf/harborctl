package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func GenerateED25519KeyPair() (publicKey string, privateKey string, err error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate ED25519 key pair: %w", err)
	}

	sshPublicKey := formatSSHPublicKey(pubKey)
	privateKeyB64 := base64.StdEncoding.EncodeToString(privKey)

	return sshPublicKey, privateKeyB64, nil
}

func formatSSHPublicKey(pubKey ed25519.PublicKey) string {
	keyType := "ssh-ed25519"
	var payload []byte

	keyTypeBytes := []byte(keyType)
	payload = append(payload, encodeString(keyTypeBytes)...)
	payload = append(payload, encodeString(pubKey)...)

	keyData := base64.StdEncoding.EncodeToString(payload)
	return fmt.Sprintf("ssh-ed25519 %s", keyData)
}

func encodeString(data []byte) []byte {
	result := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(result[0:4], uint32(len(data)))
	copy(result[4:], data)
	return result
}

func GenerateBeszelToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}
