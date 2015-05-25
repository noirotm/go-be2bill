// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"fmt"
	"testing"
)

func TestAmountSingle(t *testing.T) {
	m := Amount(SingleAmount(350))

	if !m.Immediate() {
		t.Error("amount must be immediate")
	}

	if m.Options() != nil {
		t.Error("single amount cannot be represented as options")
	}

	if fmt.Sprint(m) != "350" {
		t.Error("invalid amount")
	}
}

func TestAmountFragmented(t *testing.T) {
	m := Amount(FragmentedAmount{
		"2010-10-21": "2100",
		"2010-11-21": "1120",
	})

	if m.Immediate() {
		t.Error("amount must not be immediate")
	}

	if m.Options() == nil {
		t.Error("single amount must be represented as options")
	}

	if _, ok := m.(FragmentedAmount); !ok {
		t.Error("amount is not an FragmentedAmount")
	}
}
