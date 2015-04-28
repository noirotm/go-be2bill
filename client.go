package be2bill

import "sort"

const (
	APIVersion = "2.0"
)

const (
	OperationTypeAuthorization = "authorization"
	OperationTypeCapture       = "capture"
	OperationTypeCredit        = "credit"
	OperationTypePayment       = "payment"
	OperationTypeRefund        = "refund"
	OperationTypeStopNTimes    = "stopntimes"
)

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

var (
	productionURLs = []string{
		"https://secure-magenta1.be2bill.com",
		"https://secure-magenta2.be2bill.com",
	}

	sandboxURLs = []string{
		"https://secure-test.be2bill.com",
	}
)

func NewFormClient(credentials *Credentials) FormClient {
	var urls []string
	if credentials.production {
		urls = productionURLs
	} else {
		urls = sandboxURLs
	}

	return &formClientImpl{
		credentials,
		newHTMLRenderer(urls[0]),
		newHasher(),
	}
}

func NewDirectLinkClient(credentials *Credentials) DirectLinkClient {
	var urls []string
	if credentials.production {
		urls = productionURLs
	} else {
		urls = sandboxURLs
	}

	return &directLinkClientImpl{
		credentials,
		urls,
		newHasher(),
	}
}

func SwitchProductionURLs() {
	sort.Sort(sort.Reverse(sort.StringSlice(productionURLs)))
}

func BuildSandboxFormClient(identifier, password string) FormClient {
	return NewFormClient(SandboxUser(identifier, password))
}

func BuildProductionFormClient(identifier, password string) FormClient {
	return NewFormClient(ProductionUser(identifier, password))
}

func BuildSandboxDirectLinkClient(identifier, password string) DirectLinkClient {
	return NewDirectLinkClient(SandboxUser(identifier, password))
}

func BuildProductionDirectLinkClient(identifier, password string) DirectLinkClient {
	return NewDirectLinkClient(ProductionUser(identifier, password))
}
