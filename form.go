// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

type FormClient struct {
	credentials *Credentials
	renderer    Renderer
	hasher      Hasher
}

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

func (p *FormClient) BuildAuthorizationFormButton(amount Amount, orderID, clientID, description string, htmlOptions, options Options) string {
	params := options.copy()

	params[ParamAmount] = amount

	return p.buildProcessButton(
		OperationTypeAuthorization,
		orderID,
		clientID,
		description,
		htmlOptions,
		params,
	)
}

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
