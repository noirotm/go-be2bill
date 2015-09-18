// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// A Hasher is used to sign a request to the API.
//
// All be2bill API requests, in the Form and DirectLink clients, are
// represented as a map of parameters stored in an Options object.
//
// Before a request is sent to the Be2bill server, the parameters
// are hashed using your Be2bill account's password as a salt, and the
// resulting hash is then inserted in the parameters so any further modification
// of the parameters will render the request invalid.
//
// The default Be2bill hasher uses the SHA-256 algorithm.
type Hasher interface {
	// ComputeHash returns a hash string computed from the given password and options.
	ComputeHash(password string, options Options) string
}

type defaultHasher struct{}

func (defaultHasher) ComputeHash(password string, params Options) string {
	var clearString bytes.Buffer
	_, _ = clearString.WriteString(password)

	for _, k := range params.sortedKeys() {
		value := params[k]
		if k == ParamHash {
			continue
		}

		if valueMap, ok := value.(Options); ok {
			for _, vk := range valueMap.sortedKeys() {
				_, _ = clearString.WriteString(fmt.Sprintf("%s[%s]=%v%s", k, vk, valueMap[vk], password))
			}
		} else {
			_, _ = clearString.WriteString(fmt.Sprintf("%s=%v%s", k, value, password))
		}
	}

	return fmt.Sprintf("%x", sha256.Sum256(clearString.Bytes()))
}

// CheckHash extracts a parameter named HASH from the given options,
// computes a hash using the given hasher for the options and password,
// and compares the two strings.
// It then returns true if both strings are identical, false otherwise.
func CheckHash(hasher Hasher, password string, params Options) bool {
	receivedHash, ok := params[ParamHash].(string)
	if !ok {
		receivedHash = ""
	}

	computedHash := hasher.ComputeHash(password, params)

	return receivedHash == computedHash
}
