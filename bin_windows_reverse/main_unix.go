//go:build !windows

package main

import (
	"crypto/rand"
	"flag"
	"log"
	"os"

	"golang.org/x/crypto/chacha20poly1305"
)

// Main function for CLI required to initialize the windows payload
func main() {
	shellcode := os.Args[1]

	flag.Parse()

	if shellcode == "" {
		log.Fatal("No shellcode was passed; go run . [SHELLCODE]")
	}

	buf := []byte(shellcode)
	encrypted, key, err := encryptShellcode(buf)
	if err != nil {
		log.Fatalf("Failed to encrypt shellcode: %v", err)
	}

	err = os.WriteFile("./shell.bin", encrypted, 0600)
	if err != nil {
		log.Fatalf("Failed to save shellcode: %v", err)
	}
	err = os.WriteFile("./shell.key", key, 0600)
	if err != nil {
		log.Fatalf("Failed to save encryption key: %v", err)
	}
}

// Encrypts passed shellcode; returns encrypted, key, error
func encryptShellcode(buf []byte) (encryptedMsg []byte, key []byte, err error) {
	key = make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, nil, err
	}

	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(buf)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg = aead.Seal(nonce, nonce, buf, nil)

	return encryptedMsg, key, nil
}
