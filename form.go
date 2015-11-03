// Copyright 2015 Dalenys. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

// A FormClient builds various forms to be embedded on a merchant website
// to use Be2bill to process payments or authorizations.
type FormClient struct {
	credentials *Credentials
	renderer    Renderer
	hasher      Hasher
}

// BuildPaymentFormButton returns a payment form ready to be embedded on
// a merchant website.
//
// The amount parameter can either be immediate or fragmented.
//
// See https://developer.be2bill.com/functions/buildPaymentFormButton.
func (p *FormClient) BuildPaymentFormButton(amount Amount, orderID, clientID, description string, htmlOptions, options Options) string {
	params := options.copy()

	// Handle N-Time payments
	if amount.Immediate() {
		params[ParamAmount] = amount
	} else {
		params[ParamAmounts] = amount.Options()
	}

	return p.buildProcessButton(
		OperationTypePayment,
		orderID,
		clientID,
		description,
		htmlOptions,
		params,
	)
}

// BuildAuthorizationFormButton returns an authorization form ready to be embedded on
// a merchant website.
//
// As opposed to BuildPaymentFormButton, the amount parameter is an integer,
// because it can only be immediate.
//
// See https://developer.be2bill.com/functions/buildAuthorizationFormButton.
func (p *FormClient) BuildAuthorizationFormButton(amount int, orderID, clientID, description string, htmlOptions, options Options) string {
	params := options.copy()

	params[ParamAmount] = SingleAmount(amount)

	return p.buildProcessButton(
		OperationTypeAuthorization,
		orderID,
		clientID,
		description,
		htmlOptions,
		params,
	)
}

// General form builder
func (p *FormClient) buildProcessButton(operationType, orderID, clientID, description string, htmlOptions, options Options) string {
	options[ParamIdentifier] = p.credentials.identifier
	options[ParamOperationType] = operationType
	options[ParamOrderID] = orderID
	options[ParamClientIdent] = clientID
	options[ParamDescription] = description
	options[ParamVersion] = APIVersion

	options[ParamHash] = p.hasher.ComputeHash(p.credentials.password, options)

	return p.renderer.Render(options, htmlOptions)
}
