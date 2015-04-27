package be2bill

import "sort"

const (
	APIVersion = "2.0"

	OperationTypeAuthorization = "authorization"
	OperationTypeCapture       = "capture"
	OperationTypePayment       = "payment"

	ExecCodeSuccess     = "0000"
	ExecCodeInvalidHash = "1003"
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
