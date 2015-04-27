package be2bill

type FormClient interface {
	BuildPaymentFormButton(amount Amount, orderID, clientID, description string, htmlOptions, options Options) string
	BuildAuthorizationFormButton(amount Amount, orderID, clientID, description string, htmlOptions, options Options) string
}

type formClientImpl struct {
	credentials *Credentials
	renderer    Renderer
	hasher      Hasher
}

func (p *formClientImpl) BuildPaymentFormButton(amount Amount, orderID, clientID, description string, htmlOptions, options Options) string {
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

func (p *formClientImpl) BuildAuthorizationFormButton(amount Amount, orderID, clientID, description string, htmlOptions, options Options) string {
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

func (p *formClientImpl) buildProcessButton(operationType, orderID, clientID, description string, htmlOptions, options Options) string {
	options[ParamIdentifier] = p.credentials.identifier
	options[ParamOperationType] = operationType
	options[ParamOrderID] = orderID
	options[ParamClientIdent] = clientID
	options[ParamDescription] = description
	options[ParamVersion] = APIVersion

	options[ParamHash] = p.hasher.ComputeHash(p.credentials.password, options)

	return p.renderer.Render(options, htmlOptions)
}
