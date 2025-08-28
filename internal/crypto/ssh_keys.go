package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

// GenerateED25519KeyPair gera um par de chaves ED25519 para o Beszel
func GenerateED25519KeyPair() (publicKey string, privateKey string, err error) {
	// Gerar par de chaves ED25519
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate ED25519 key pair: %w", err)
	}

	// Converter chave pública para formato SSH
	sshPublicKey := formatSSHPublicKey(pubKey)

	// Converter chave privada para base64 (formato que o Beszel usa internamente)
	privateKeyB64 := base64.StdEncoding.EncodeToString(privKey)

	return sshPublicKey, privateKeyB64, nil
}

// formatSSHPublicKey converte uma chave pública ED25519 para formato SSH correto
func formatSSHPublicKey(pubKey ed25519.PublicKey) string {
	// Criar o formato wire protocol SSH para ed25519
	keyType := "ssh-ed25519"

	// Construir payload SSH wire format:
	// string "ssh-ed25519" + string <32-byte-public-key>
	var payload []byte

	// Adicionar o tipo da chave
	keyTypeBytes := []byte(keyType)
	payload = append(payload, encodeString(keyTypeBytes)...)

	// Adicionar a chave pública
	payload = append(payload, encodeString(pubKey)...)

	// Codificar em base64
	keyData := base64.StdEncoding.EncodeToString(payload)

	return fmt.Sprintf("ssh-ed25519 %s", keyData)
}

// encodeString codifica uma string no formato SSH wire protocol
func encodeString(data []byte) []byte {
	// SSH wire format: 4 bytes (length) + data
	result := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(result[0:4], uint32(len(data)))
	copy(result[4:], data)
	return result
}

// GenerateBeszelToken gera um token simples para o Beszel
func GenerateBeszelToken() (string, error) {
	// Gerar 32 bytes aleatórios para o token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Converter para base64
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}
