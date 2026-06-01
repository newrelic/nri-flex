/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func deriveKey(passphrase string) ([]byte, error) {
	// Static salt is acceptable here because nri-flex encrypts ephemeral config
	// values (e.g. passwords passed via CLI). The passphrase itself provides the
	// entropy, and a per-ciphertext random salt would require a format change to
	// store it alongside the ciphertext, breaking the existing API contract.
	salt := []byte("nri-flex-static-salt")
	return pbkdf2.Key(sha256.New, passphrase, salt, 100000, 32)
}

// Encrypt string
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	key, err := deriveKey(passphrase)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt string
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	key, err := deriveKey(passphrase)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
