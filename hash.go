// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Hasher interface {
	ComputeHash(string, Options) string
	CheckHash(string, Options) bool
}

type defaultHasher struct {
}

func newHasher() Hasher {
	return new(defaultHasher)
}

func (p *defaultHasher) ComputeHash(password string, params Options) string {
	var clearString bytes.Buffer
	clearString.WriteString(password)

	for _, k := range params.sortedKeys() {
		value := params[k]
		if k == ParamHash {
			continue
		}

		if valueMap, ok := value.(Options); ok {
			for _, vk := range valueMap.sortedKeys() {
				clearString.WriteString(fmt.Sprintf("%s[%s]=%v%s", k, vk, valueMap[vk], password))
			}
		} else {
			clearString.WriteString(fmt.Sprintf("%s=%v%s", k, value, password))
		}
	}

	return fmt.Sprintf("%x", sha256.Sum256(clearString.Bytes()))
}

func (p *defaultHasher) CheckHash(password string, params Options) bool {
	receivedHash, ok := params[ParamHash].(string)
	if !ok {
		receivedHash = ""
	}

	computedHash := p.ComputeHash(password, params)

	return receivedHash == computedHash
}
