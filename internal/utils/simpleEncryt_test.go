/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "test-passphrase"
	plaintext := []byte("hello world, this is sensitive data")

	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted text does not match original: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptWithShortCiphertext(t *testing.T) {
	_, err := Decrypt([]byte("short"), "passphrase")
	if err == nil {
		t.Fatal("Decrypt should fail with ciphertext shorter than nonce size")
	}
}

func TestDecryptWithWrongPassphrase(t *testing.T) {
	passphrase := "correct-passphrase"
	wrongPassphrase := "wrong-passphrase"
	plaintext := []byte("secret data")

	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(ciphertext, wrongPassphrase)
	if err == nil {
		t.Fatal("Decrypt should fail with wrong passphrase")
	}
}
