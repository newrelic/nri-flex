/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" // nolint
	"crypto/rand"
	"encoding/hex"
	"io"
)

func createHash(key string) string {
	hasher := md5.New() // nolint
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Encrypt string
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	ciphertext := []byte("")
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err == nil {

		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
			ciphertext = gcm.Seal(nonce, nonce, data, nil)
		}
	}

	return ciphertext, err
}

//Decrypt string
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	var err error
	var block cipher.Block
	var gcm cipher.AEAD
	plaintext := []byte("")
	key := []byte(createHash(passphrase))
	block, err = aes.NewCipher(key)
	if err == nil {

		gcm, err = cipher.NewGCM(block)
		if err == nil {
			nonceSize := gcm.NonceSize()
			nonce, ciphertext := data[:nonceSize], data[nonceSize:]
			plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
		}
	}
	return plaintext, err
}
