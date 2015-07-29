// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import "testing"

func TestSimpleHash(t *testing.T) {
	opts := Options{"c": 3, "a": "1", "b": "2"}
	hash := &defaultHasher{}

	if h := hash.ComputeHash("password", opts); h != "77c71c1e70ea28525cf078537d22d1932922e3741ed83287b0dc0a117bf77999" {
		t.Error("Got", h)
	}
}

func TestHashWithParameter(t *testing.T) {
	opts := Options{"c": 3, "a": "1", "b": "2", "HASH": "shouldnotimpact"}
	hash := &defaultHasher{}

	if h := hash.ComputeHash("password", opts); h != "77c71c1e70ea28525cf078537d22d1932922e3741ed83287b0dc0a117bf77999" {
		t.Error("Got", h)
	}
}

func TestHashAmount(t *testing.T) {
	a := Amount(SingleAmount(2510))
	opts := Options{"a": a}
	opts2 := Options{"a": "2510"}
	hash := &defaultHasher{}

	h := hash.ComputeHash("password", opts)
	h2 := hash.ComputeHash("password", opts2)

	if h != h2 {
		t.Errorf("invalid hash, expected %s, got %s", h2, h)
	}
}

func TestHashAmountFragmented(t *testing.T) {
	a := Amount(FragmentedAmount{
		"2010-10-21": "2100",
		"2010-11-21": "1120",
	})
	opts := Options{"a": a.Options()}
	opts2 := Options{
		"a": Options{
			"2010-10-21": "2100",
			"2010-11-21": "1120",
		},
	}
	hash := &defaultHasher{}

	h := hash.ComputeHash("password", opts)
	h2 := hash.ComputeHash("password", opts2)

	if h != h2 {
		t.Errorf("invalid hash, expected %s, got %s", h2, h)
	}
}

func TestHashRecursive(t *testing.T) {
	opts := Options{
		"c": 3,
		"a": "1",
		"b": "2",
		"d": Options{
			"y": 43,
			"x": 42,
		},
	}
	hash := &defaultHasher{}

	if h := hash.ComputeHash("password", opts); h != "376383093261372eb97909ed1a44b1adb5e8f2687f7a64f1c41d5a0c8cc0b0fa" {
		t.Error("Got", h)
	}
}

func TestHashCheck(t *testing.T) {
	hash := &defaultHasher{}
	ok := hash.CheckHash("password", Options{
		"c":    3,
		"a":    "1",
		"b":    "2",
		"HASH": "77c71c1e70ea28525cf078537d22d1932922e3741ed83287b0dc0a117bf77999",
	})
	if !ok {
		t.Error("Invalid hash")
	}
}

func TestHashCheckCapture(t *testing.T) {
	hash := &defaultHasher{}
	o := Options{
		"DESCRIPTION":   "Capture test 01",
		"HASH":          "ea22191f962b6fc708f48baff8600f8caaea3afa7e8be3f2d2dafb3249396e72",
		"IDENTIFIER":    "IDENTIFIER",
		"OPERATIONTYPE": "capture",
		"ORDERID":       "order_21",
		"TRANSACTIONID": "test1",
		"VERSION":       "2.0",
	}
	ok := hash.CheckHash("PASSWORD", o)
	if !ok {
		t.Error("Invalid hash")
	}
}
