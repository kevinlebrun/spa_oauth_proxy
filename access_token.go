package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expiration   int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (at *AccessToken) String() string {
	return fmt.Sprintf("%s|%s|%d|%s", at.AccessToken, at.TokenType, at.Expiration, at.RefreshToken)
}

func (at *AccessToken) Encode(key []byte) (string, error) {
	return encrypt(key, at.String())
}

func DecodeAccessToken(key []byte, value string) (*AccessToken, error) {
	token, err := decrypt(key, value)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(token, "|")
	at := &AccessToken{
		AccessToken:  parts[0],
		TokenType:    parts[1],
		RefreshToken: parts[3],
	}
	at.Expiration, _ = strconv.Atoi(parts[2])

	return at, nil
}

func encrypt(key []byte, text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, text string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("Ciphertext too short: %s\n", ciphertext)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
