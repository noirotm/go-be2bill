// Copyright 2016 Marc Noirot. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"fmt"
	"net/url"
	"sort"
)

// Options is a map of name/value parameters used to represent a request to
// the Be2bill API.
//
// Options are also used to pass optional parameters for various methods.
// Internally, all method calls will use this type so a hash can be computed
// to protect the request from modification.
//
// The be2bill.Param* constants represent the names supported by the plateform.
// Each operation supports a different set of parameters, as described in the
// platform documentation.
type Options map[string]interface{}

func (p Options) sortedKeys() []string {
	keys := make([]string, len(p))
	i := 0
	for k := range p {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (p Options) copy() Options {
	c := make(Options)
	for k, v := range p {
		c[k] = v
	}
	return c
}

func recurseFlatten(name string, options, result Options) {
	for k, v := range options {
		key := fmt.Sprintf("%s[%s]", name, k)
		if opts, ok := v.(Options); ok {
			recurseFlatten(key, opts, result)
		} else {
			result[key] = fmt.Sprint(v)
		}
	}
}

func (p Options) flatten() Options {
	result := Options{}
	for k, v := range p {
		if opts, ok := v.(Options); ok {
			recurseFlatten(k, opts, result)
		} else {
			result[k] = fmt.Sprint(v)
		}
	}
	return result
}

func (p Options) urlValues() url.Values {
	values := url.Values{}
	opts := p.flatten()

	for k, v := range opts {
		values.Set(k, fmt.Sprint(v))
	}

	return values
}
