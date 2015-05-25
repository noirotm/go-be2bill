// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

type Amount interface {
	Immediate() bool
	Options() Options
}

type SingleAmount int

func (SingleAmount) Immediate() bool {
	return true
}

func (SingleAmount) Options() Options {
	return nil
}

type FragmentedAmount Options

func (FragmentedAmount) Immediate() bool {
	return false
}

func (p FragmentedAmount) Options() Options {
	return Options(p)
}
