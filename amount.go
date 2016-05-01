// Copyright 2016 Marc Noirot. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

// An Amount represents a sum of money as understood by the platform.
// Amounts are usually immediate but can also be fragmented, for example
// if several payments are scheduled for future dates.
type Amount interface {
	// Immediate returns true if the amount is due immediately.
	Immediate() bool

	// Options returns the amount as an Options object suitable
	// for use as parameter for server calls.
	// This is only relevant if the amount is not immediate.
	Options() Options
}

// A SingleAmount is a simple amount expressed in cents.
// It will be typically initialized like this:
//
//   amount := be2bill.SingleAmount(2350)
//
type SingleAmount int

// Immediate always returns true for single amounts.
func (SingleAmount) Immediate() bool {
	return true
}

// Options always returns nil for single amounts.
func (SingleAmount) Options() Options {
	return nil
}

// A FragmentedAmount is a map of future dates, expressed as "YYYY-MM-DD"
// strings, to numeric amounts in cents.
// It will be typically initialized like this:
//
//   amount := be2bill.FragmentedAmount{"2016-05-14": 15235, "2016-06-14": 14723}
//
type FragmentedAmount Options

// Immediate always returns false for fragmented amounts.
func (FragmentedAmount) Immediate() bool {
	return false
}

// Options returns the fragmented amount as an Options instance.
func (p FragmentedAmount) Options() Options {
	return Options(p)
}
