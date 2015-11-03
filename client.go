// Copyright 2015 Dalenys. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

// APIVersion represents the currently supported version of the API.
// It is sent to the server as part of the various API calls in order
// to ensure compatibility.
const APIVersion = "2.0"

// These constants represent the available keys for the htmlOptions parameters.
const (
	HTMLOptionForm   = "FORM"
	HTMLOptionSubmit = "SUBMIT"
)

// These constants represent the possible keys for the options parameters.
// See https://developer.be2bill.com/annexes/parameters for more information.
const (
	Param3DSecure            = "3DSECURE"
	Param3DSecureDisplayMode = "3DSECUREDISPLAYMODE"
	ParamAlias               = "ALIAS"
	ParamAliasMode           = "ALIASMODE"
	ParamAmount              = "AMOUNT"
	ParamAmounts             = "AMOUNTS"
	ParamBillingAddress      = "BILLINGADDRESS"
	ParamBillingCountry      = "BILLINGCOUNTRY"
	ParamBillingFirstName    = "BILLINGFIRSTNAME"
	ParamBillingLastName     = "BILLINGLASTNAME"
	ParamBillingPhone        = "BILLINGPHONE"
	ParamBillingPostalCode   = "BILLINGPOSTALCODE"
	ParamCallbackURL         = "CALLBACKURL"
	ParamCardCode            = "CARDCODE"
	ParamCardCVV             = "CARDCVV"
	ParamCardFullName        = "CARDFULLNAME"
	ParamCardValidityDate    = "CARDVALIDITYDATE"
	ParamClientAddress       = "CLIENTADDRESS"
	ParamClientDOB           = "CLIENTDOB"
	ParamClientEmail         = "CLIENTEMAIL"
	ParamClientIdent         = "CLIENTIDENT"
	ParamClientIP            = "CLIENTIP"
	ParamClientReferrer      = "CLIENTREFERRER"
	ParamClientUserAgent     = "CLIENTUSERAGENT"
	ParamCompression         = "COMPRESSION"
	ParamCreateAlias         = "CREATEALIAS"
	ParamDate                = "DATE"
	ParamDay                 = "DAY"
	ParamDescription         = "DESCRIPTION"
	ParamDisplayCreateAlias  = "DISPLAYCREATEALIAS"
	ParamEndDate             = "ENDDATE"
	ParamExtraData           = "EXTRADATA"
	ParamHash                = "HASH"
	ParamHideCardFullName    = "HIDECARDFULLNAME"
	ParamHideClientEmail     = "HIDECLIENTEMAIL"
	ParamIdentifier          = "IDENTIFIER"
	ParamLanguage            = "LANGUAGE"
	ParamMailTo              = "MAILTO"
	ParamMetadata            = "METADATA"
	ParamOperationType       = "OPERATIONTYPE"
	ParamOrderID             = "ORDERID"
	ParamScheduleID          = "SCHEDULEID"
	ParamShipToAddress       = "SHIPTOADDRESS"
	ParamShipToCountry       = "SHIPTOCOUNTRY"
	ParamShipToFirstName     = "SHIPTOFIRSTNAME"
	ParamShipToLastName      = "SHIPTOLASTNAME"
	ParamShipToPhone         = "SHIPTOPHONE"
	ParamShipToPostalCode    = "SHIPTOPOSTALCODE"
	ParamStartDate           = "STARTDATE"
	ParamTimeZone            = "TIMEZONE"
	ParamTransactionID       = "TRANSACTIONID"
	ParamVersion             = "VERSION"
	ParamVME                 = "VME"
)

// These constants represent the operation codes supported by the be2bill API calls.
const (
	OperationTypeAuthorization                = "authorization"
	OperationTypeCapture                      = "capture"
	OperationTypeCredit                       = "credit"
	OperationTypePayment                      = "payment"
	OperationTypeRefund                       = "refund"
	OperationTypeStopNTimes                   = "stopntimes"
	OperationTypeGetTransactions              = "getTransactions"
	OperationTypeExportTransactions           = "exportTransactions"
	OperationTypeExportChargebacks            = "exportChargebacks"
	OperationTypeExportReconciliation         = "exportReconciliation"
	OperationTypeExportReconciledTransactions = "exportReconciledTransactions"
)

// These constants represent the compression formats supported by the export
// methods.
const (
	CompressionZip  = "ZIP"
	CompressionGzip = "GZIP"
	CompressionBzip = "BZIP"
)

// These constants represent the possible keys for the API calls' result maps.
const (
	ResultParamOperationType = "OPERATIONTYPE"
	ResultParamTransactionID = "TRANSACTIONID"
	ResultParamExecCode      = "EXECCODE"
	ResultParamMessage       = "MESSAGE"
	ResultParamDescriptor    = "DESCRIPTOR"
	ResultParamAmount        = "AMOUNT"
	ResultParamRedirectHTML  = "REDIRECTHTML"
)

// These constants represent the possible values for the exec code result field.
// See https://developer.be2bill.com/annexes/execcodes for more information.
const (
	ExecCodeSuccess                   = "0000"
	ExecCode3DSecureRequired          = "0001"
	ExecCodeAlternateRedirectRequired = "0002"

	ExecCodeMissingParameter    = "1001"
	ExecCodeInvalidParameter    = "1002"
	ExecCodeInvalidHash         = "1003"
	ExecCodeUnsupportedProtocol = "1004"

	ExecCodeAliasNotFound              = "2001"
	ExecCodeTransactionNotFound        = "2002"
	ExecCodeUnsuccessfulTransaction    = "2003"
	ExecCodeTransactionNotRefundable   = "2004"
	ExecCodeAuthorizationNotCapturable = "2005"
	ExecCodeIncompleteTransaction      = "2006"
	ExecCodeInvalidCaptureAmount       = "2007"
	ExecCodeInvalidRefundAmount        = "2008"
	ExecCodeAuthorizationTimeout       = "2009"
	ExecCodeScheduleNotFound           = "2010"
	ExecCodeInterruptedSchedule        = "2011"
	ExecCodeScheduleFinished           = "2012"

	ExecCodeAccountDeactivated      = "3001"
	ExecCodeUnauthorizedServerIP    = "3002"
	ExecCodeUnauthorizedTransaction = "3003"

	ExecCodeTransactionRefusedBank        = "4001"
	ExecCodeUnsufficientFunds             = "4002"
	ExecCodeCardRefused                   = "4003"
	ExecCodeTransactionAbandoned          = "4004"
	ExecCodeSuspectedFraud                = "4005"
	ExecCodeCardLost                      = "4006"
	ExecCodeCardStolen                    = "4007"
	ExecCode3DSecureAuthenticationFailed  = "4008"
	ExecCode3DSecureAuthenticationTimeout = "4009"
	ExecCodeInvalidTransaction            = "4010"
	ExecCodeDuplicateTransaction          = "4011"
	ExecCodeInvalidCardData               = "4012"
	ExecCodeTransactionNotAuthorized      = "4013"
	ExecCodeCard3DSecureNotSupported      = "4014"
	ExecCodeTransactionTimeout            = "4015"
	ExecCodeTransactionRefusedByTerminal  = "4016"

	ExecCodeExchangeProtocolError = "5001"
	ExecCodeBankNetworkError      = "5002"
	ExecCodeHandlerTimeout        = "5004"
	ExecCode3DSecureDisplayError  = "5005"

	ExecCodeTransactionRefusedMerchant      = "6001"
	ExecCodeTransactionRefusedUnknown       = "6002"
	ExecCodeTransactionChallenged           = "6003"
	ExecCodeTransactionRefusedMerchantRules = "6004"
)

// An Environment is the set of URLs that represent a be2bill endpoint.
// There are currently two existing environments, production and sandbox.
//
// The production environment is where real transactions and sales take place.
//
// The sandbox can be used for testing and no real operations occur.
type Environment []string

var (
	// EnvProduction represents the production servers where real operations take place.
	EnvProduction Environment
	// EnvSandbox represents the testing server where operations are simulated.
	EnvSandbox Environment
)

func init() {
	EnvProduction = Environment{
		"https://secure-magenta1.be2bill.com",
		"https://secure-magenta2.be2bill.com",
	}
	EnvSandbox = Environment{
		"https://secure-test.be2bill.com",
	}
}

// SwitchURLs reverses the order of URLs defined in the environment.
// This is useful to test the failover URLs.
func (p Environment) SwitchURLs() {
	for i := len(p)/2 - 1; i >= 0; i-- {
		opp := len(p) - 1 - i
		p[i], p[opp] = p[opp], p[i]
	}
}

// Credentials represent the information that is necessary for
// the clients to identify themselves to the API.
type Credentials struct {
	identifier  string
	password    string
	environment Environment
}

// User returns new Credentials using the given client identifiers and
// environment.
func User(identifier string, password string, environment Environment) *Credentials {
	return &Credentials{identifier, password, environment}
}

// ProductionUser returns new Credentials using the given client identifiers
// and the production environment.
func ProductionUser(identifier string, password string) *Credentials {
	return &Credentials{identifier, password, EnvProduction}
}

// SandboxUser returns new Credentials using the given client identifiers
// and the sandbox environment.
func SandboxUser(identifier string, password string) *Credentials {
	return &Credentials{identifier, password, EnvSandbox}
}

// NewFormClient returns a new FormClient using the given credentials.
func NewFormClient(credentials *Credentials) *FormClient {
	return &FormClient{
		credentials,
		newHTMLRenderer(credentials.environment[0]),
		&defaultHasher{},
	}
}

// NewDirectLinkClient returns a new DirectLinkClient using the given
// credentials.
func NewDirectLinkClient(credentials *Credentials) *DirectLinkClient {
	return &DirectLinkClient{
		credentials,
		credentials.environment,
		&defaultHasher{},
	}
}

// BuildSandboxFormClient returns a new FormClient using the given
// client identifier and password for the sandbox environment.
func BuildSandboxFormClient(identifier, password string) *FormClient {
	return NewFormClient(SandboxUser(identifier, password))
}

// BuildProductionFormClient returns a new FormClient using the given
// client identifier and password for the production environment.
func BuildProductionFormClient(identifier, password string) *FormClient {
	return NewFormClient(ProductionUser(identifier, password))
}

// BuildSandboxDirectLinkClient returns a new DirectLinkClient using the given
// client identifier and password for the sandbox environment.
func BuildSandboxDirectLinkClient(identifier, password string) *DirectLinkClient {
	return NewDirectLinkClient(SandboxUser(identifier, password))
}

// BuildProductionDirectLinkClient returns a new DirectLinkClient using the given
// client identifier and password for the production environment.
func BuildProductionDirectLinkClient(identifier, password string) *DirectLinkClient {
	return NewDirectLinkClient(ProductionUser(identifier, password))
}
