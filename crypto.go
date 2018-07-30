package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

type Crypto struct {
	SecretKey [32]byte
}

func (c Crypto) Init() {
	secretKeyBytes, err := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	if err != nil {
		panic(err)
	}
	copy(c.SecretKey[:], secretKeyBytes)
}

func (c Crypto) Encrypt(message []byte) []byte {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	encrypted := secretbox.Seal(nonce[:], []byte(message), &nonce, &c.SecretKey)

	return encrypted
}

func (c Crypto) Decrypt(encrypted []byte) []byte {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, _ := secretbox.Open(nil, encrypted[24:], &decryptNonce, &c.SecretKey)
	return decrypted
}
