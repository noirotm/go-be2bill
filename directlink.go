// Copyright 2016 Marc Noirot. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// These strings represent the alias modes supported by the payment and
// authorization methods of the Direct Link API.
const (
	aliasModeOneClick     = "oneclick"
	aliasModeSubscription = "subscription"
)

// these strings represent the search modes supported by the getTransactions
// operation.
const (
	searchByOrderID       = "ORDERID"
	searchByTransactionID = "TRANSACTIONID"
)

// A Result represents the return value from a call to an operation of the
// DirectLink API.
// It is a map of key/values where keys are strings and values are generic.
//
// See the Notification Parameters at https://developer.be2bill.com/annexes/parameters
// for the supported keys.
type Result map[string]interface{}

// StringValue returns the value for the given property of a Result object.
// The value must be of string type, otherwise an empty string is returned instead.
func (r Result) StringValue(name string) string {
	val, ok := r[name]
	if !ok {
		return ""
	}
	return val.(string)
}

// OperationType returns the name of the operation that returned this object.
func (r Result) OperationType() string {
	return r.StringValue(ResultParamOperationType)
}

// ExecCode returns the execution code that represents the success status
// of the operation.
//
// See https://developer.be2bill.com/annexes/execcodes for a list of supported
// execution codes.
func (r Result) ExecCode() string {
	return r.StringValue(ResultParamExecCode)
}

// Message returns the textual message associated with the result's execution
// code.
func (r Result) Message() string {
	return r.StringValue(ResultParamMessage)
}

// TransactionID returns the identifier of the transaction associated with
// the current operation.
func (r Result) TransactionID() string {
	return r.StringValue(ResultParamTransactionID)
}

// Success returns true if the operation succeeded, false otherwise.
func (r Result) Success() bool {
	return r.ExecCode() == ExecCodeSuccess
}

const (
	directLinkPath     = "/front/service/rest/process"
	exportPath         = "/front/service/rest/export"
	reconciliationPath = "/front/service/rest/reconciliation"

	defaultRequestTimeout = 30 * time.Second
)

var (
	// ErrTimeout is returned by DirectClient operations if the request
	// encounters a timeout and cannot finish.
	ErrTimeout = errors.New("timeout")
	// ErrURLMissing is returned by DirectClient operations if the current
	// environment has no URL specified.
	ErrURLMissing = errors.New("no URL provided")
	// ErrServerError is returned by DirectClient operations if the request
	// encounters a server-side error.
	ErrServerError = errors.New("server error")
)

// A DirectLinkClient represent an access to the Direct Link Be2bill API
// that allows a merchant to perform direct calls to the Be2bill servers
// without the need for a graphical interface such as a web page.
// It supports a variety of operations used to perform payments, captures,
// or getting informations about past transactions.
type DirectLinkClient struct {
	credentials *Credentials
	urls        []string
	hasher      Hasher
	// RequestTimeout is the duration after which requests time out
	// and return an ErrTimeout error.
	// The default timeout is 30 seconds.
	RequestTimeout time.Duration
}

// NewDirectLinkClient returns a new DirectLinkClient using the given
// credentials.
func NewDirectLinkClient(credentials *Credentials) *DirectLinkClient {
	return &DirectLinkClient{
		credentials,
		credentials.environment,
		&defaultHasher{},
		defaultRequestTimeout,
	}
}

func (p *DirectLinkClient) getURLs(path string) []string {
	urls := make([]string, len(p.urls))
	for i, url := range p.urls {
		urls[i] = url + path
	}
	return urls
}

func (p *DirectLinkClient) getDirectLinkURLs() []string {
	return p.getURLs(directLinkPath)
}

func (p *DirectLinkClient) doPostRequest(url string, params Options) (Result, error) {
	requestParams := Options{
		"method": params[ParamOperationType],
		"params": params,
	}

	responseChan := make(chan Result, 1)
	errChan := make(chan error, 1)

	go func() {
		client := &http.Client{
			Timeout: p.RequestTimeout + 5*time.Second,
		}

		resp, err := client.PostForm(url, requestParams.urlValues())
		if err != nil {
			errChan <- err
			return
		}

		if resp.StatusCode != 200 {
			errChan <- ErrServerError
			return
		}

		defer func() { _ = resp.Body.Close() }()

		r := json.NewDecoder(resp.Body)
		result := make(Result)
		err = r.Decode(&result)
		if err != nil {
			errChan <- err
			return
		}

		responseChan <- result
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-time.After(p.RequestTimeout):
		return nil, ErrTimeout
	case result := <-responseChan:
		return result, nil
	}
}

func (p *DirectLinkClient) requests(urls []string, params Options) (Result, error) {
	if len(urls) == 0 {
		return nil, ErrURLMissing
	}

	var errRet error
	for _, url := range urls {
		result, err := p.doPostRequest(url, params)
		if err != nil {
			// break if a timeout occurred, otherwise try next URL
			if err == ErrTimeout {
				return nil, err
			}
			errRet = err
			continue
		}

		return result, err
	}

	return nil, errRet
}

func (p *DirectLinkClient) transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOrderID] = orderID
	params[ParamClientIdent] = clientID
	params[ParamClientEmail] = clientEmail
	params[ParamDescription] = description
	params[ParamClientUserAgent] = clientUserAgent
	params[ParamClientIP] = clientIP
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getDirectLinkURLs(), params)
}

func isHTTPURL(str string) bool {
	url, err := url.Parse(str)
	return err == nil && (url.Scheme == "http" || url.Scheme == "https")
}

func (p *DirectLinkClient) getTransactions(searchBy string, idList []string, destination, compression string) (Result, error) {
	params := Options{}
	params[ParamOperationType] = OperationTypeGetTransactions
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	id := strings.Join(idList, ";")

	if searchBy == searchByOrderID {
		params[ParamOrderID] = id
	} else if searchBy == searchByTransactionID {
		params[ParamTransactionID] = id
	}

	params[ParamCompression] = compression
	if isHTTPURL(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(exportPath), params)
}

// Payment performs a payment operation using the given card holder information.
// Immediate and fragmented amounts are supported for this method.
//
// See https://developer.be2bill.com/functions/payment
func (p *DirectLinkClient) Payment(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	// Handle N-Time payments
	if amount.Immediate() {
		params[ParamAmount] = amount
	} else {
		params[ParamAmounts] = amount.Options()
	}

	params[ParamOperationType] = OperationTypePayment
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// Authorization performs an authorization operation using the given card holder information.
//
// See https://developer.be2bill.com/functions/authorization
func (p *DirectLinkClient) Authorization(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount int,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName
	params[ParamAmount] = amount

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// Credit performs a credit operation using the given card holder information.
func (p *DirectLinkClient) Credit(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount int,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeCredit
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName
	params[ParamAmount] = amount

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// OneClickPayment performs a payment operation for an already registered client.
// Card data must be present for the given client ID using an alias created
// in a previous payment or authorization operation.
// Immediate and fragmented amounts are supported for this method.
//
// See https://developer.be2bill.com/functions/oneClickPayment
func (p *DirectLinkClient) OneClickPayment(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	// Handle N-Time payments
	if amount.Immediate() {
		params[ParamAmount] = amount
	} else {
		params[ParamAmounts] = amount.Options()
	}

	params[ParamOperationType] = OperationTypePayment
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeOneClick

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// Refund performs a refund operation over a given previous transaction.
// If a refund is done the same day as the initial payment, it will
// be a simple payment cancellation and no remote collection will be made.
//
// See https://developer.be2bill.com/functions/refund
func (p *DirectLinkClient) Refund(transactionID, orderID, description string, options Options) (Result, error) {
	params := options.copy()

	params[ParamIdentifier] = p.credentials.identifier
	params[ParamOperationType] = OperationTypeRefund
	params[ParamDescription] = description
	params[ParamTransactionID] = transactionID
	params[ParamVersion] = APIVersion
	params[ParamOrderID] = orderID

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getDirectLinkURLs(), params)
}

// Capture performs a capture operation on a previous authorization.
// The capture allows to collect the carholder's funds up to 7 days after
// the authorization has been made.
//
// An optional amount can be specified in the options parameter. It will
// replace the authorized amount so the capture will be partial.
// Only an amount inferior to the original can be used.
//
// See https://developer.be2bill.com/functions/capture
func (p *DirectLinkClient) Capture(transactionID, orderID, description string, options Options) (Result, error) {
	params := options.copy()

	params[ParamIdentifier] = p.credentials.identifier
	params[ParamOperationType] = OperationTypeCapture
	params[ParamVersion] = APIVersion
	params[ParamDescription] = description
	params[ParamTransactionID] = transactionID
	params[ParamOrderID] = orderID

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getDirectLinkURLs(), params)
}

// OneClickAuthorization performs an authorization operation for an already registered client.
// Card data must be present for the given client ID using an alias created
// in a previous payment or authorization operation.
// Only immediate amounts are supported for this method.
//
// See https://developer.be2bill.com/functions/oneClickAuthorization
func (p *DirectLinkClient) OneClickAuthorization(
	alias string,
	amount int, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeOneClick
	params[ParamAmount] = amount

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// SubscriptionAuthorization performs an authorization operation for an already registered client.
// This is used for periodic payments like monthly billing for a continuous service.
// Card data must be present for the given client ID using an alias created
// in a previous payment or authorization operation.
// Only immediate amounts are supported for this method.
//
// See https://developer.be2bill.com/functions/subscriptionAuthorization
func (p *DirectLinkClient) SubscriptionAuthorization(
	alias string,
	amount int, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeSubscription
	params[ParamAmount] = amount

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// SubscriptionPayment performs a payment operation for an already registered client.
// This is used for periodic payments like monthly billing for a continuous service.
// Card data must be present for the given client ID using an alias created
// in a previous payment or authorization operation.
// Only immediate amounts are supported for this method.
//
// See https://developer.be2bill.com/functions/subscriptionPayment
func (p *DirectLinkClient) SubscriptionPayment(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	// Handle N-Time payments
	if amount.Immediate() {
		params[ParamAmount] = amount
	} else {
		params[ParamAmounts] = amount.Options()
	}

	params[ParamOperationType] = OperationTypePayment
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeSubscription

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// StopNTimes cancels future scheduled payments when a transaction has been made
// using a fragmented amount.
// The initial transaction identifier is used as parameter.
//
// See https://developer.be2bill.com/functions/stopNTimes
func (p *DirectLinkClient) StopNTimes(scheduleID string, options Options) (Result, error) {
	params := options.copy()

	params[ParamIdentifier] = p.credentials.identifier
	params[ParamOperationType] = OperationTypeStopNTimes
	params[ParamScheduleID] = scheduleID
	params[ParamVersion] = APIVersion

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getDirectLinkURLs(), params)
}

// RedirectForPayment returns HTML code used to  to redirect the customer to
// an alternative payment service such as PayPal once his cart is validated.
// The result object contains a field named `be2bill.ResultParamRedirectHTML`
// representing a Base64 representation of the HTML code that must be inserted
// into a page to perform the redirection.
//
// For example:
//
//    result, err := client.RedirectForPayment(
//    	100,
//    	"order_1446456185",
//   	"6328_john.smith",
//    	"6328_john.smith@gmail.com",
//    	"123.123.123.123",
//    	"be2bill_transaction_processed_by_PayPal",
//    	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
//    );
//    str := result.StringValue(be2bill.ResultParamRedirectHTML)
//    htmlCode, err := base64.StdEncoding.DecodeString(str)
//
// See https://developer.be2bill.com/functions/redirectForPayment
func (p *DirectLinkClient) RedirectForPayment(
	amount int,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypePayment
	params[ParamAmount] = amount

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

// GetTransactionsByTransactionID retrieves a list of transactions given a list
// of transaction identifiers, and sends the given list as a file to a destination
// which can be a HTTP URL or an email address.
// The file is compressed used the specified compression scheme, as specified
// in the be2bill.Compression* constants.
//
// See https://developer.be2bill.com/functions/getTransactionsByTransactionId
func (p *DirectLinkClient) GetTransactionsByTransactionID(transactionIDs []string, destination, compression string) (Result, error) {
	return p.getTransactions(searchByTransactionID, transactionIDs, destination, compression)
}

// GetTransactionsByOrderID retrieves a list of transactions given a list
// of order identifiers, and sends the given list as a file to a destination
// which can be a HTTP URL or an email address.
// The file is compressed used the specified compression scheme, as specified
// in the be2bill.Compression* constants.
//
// See https://developer.be2bill.com/functions/getTransactionsByOrderId
func (p *DirectLinkClient) GetTransactionsByOrderID(orderIDs []string, destination, compression string) (Result, error) {
	return p.getTransactions(searchByOrderID, orderIDs, destination, compression)
}

// ExportTransactions retrieves a list of transactions given a date or
// interval of dates, and sends the given list as a file to a destination
// which can be a HTTP URL or an email address.
// The file is compressed used the specified compression scheme, as specified
// in the be2bill.Compression* constants.
//
// Dates can be specified as a month (YYYY-MM), or as a day (YYYY-MM-DD).
// If endDate is an empty string, transactions at startDate will be retrieved,
// otherwise transactions between startDate and endDate will be retrieved.
//
// See https://developer.be2bill.com/functions/exportTransactions
func (p *DirectLinkClient) ExportTransactions(startDate, endDate, destination, compression string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeExportTransactions
	params[ParamCompression] = compression
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	// date can be either an interval or a single value
	if len(endDate) > 0 {
		params[ParamStartDate] = startDate
		params[ParamEndDate] = endDate
	} else {
		params[ParamDate] = startDate
	}

	if isHTTPURL(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(exportPath), params)
}

// ExportChargebacks retrieves a list of chargebacks given a date or
// interval of dates, and sends the given list as a file to a destination
// which can be a HTTP URL or an email address.
// The file is compressed used the specified compression scheme, as specified
// in the be2bill.Compression* constants.
//
// Dates can be specified as a month (YYYY-MM), or as a day (YYYY-MM-DD).
// If endDate is an empty string, transactions at startDate will be retrieved,
// otherwise transactions between startDate and endDate will be retrieved.
//
// See https://developer.be2bill.com/functions/exportChargebacks
func (p *DirectLinkClient) ExportChargebacks(startDate, endDate, destination, compression string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeExportChargebacks
	params[ParamCompression] = compression
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	// date can be either an interval or a single value
	if len(endDate) > 0 {
		params[ParamStartDate] = startDate
		params[ParamEndDate] = endDate
	} else {
		params[ParamDate] = startDate
	}

	if isHTTPURL(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(exportPath), params)
}

// ExportReconciliation retrieves the final reconciliation for a given a date
// and sends it as a file to a destination which can be a HTTP URL or an email address.
// The file is compressed used the specified compression scheme, as specified
// in the be2bill.Compression* constants.
//
// The date can be specified as a month (YYYY-MM), or as a day (YYYY-MM-DD).
//
// See https://developer.be2bill.com/functions/exportReconciliation
func (p *DirectLinkClient) ExportReconciliation(date, destination, compression string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeExportReconciliation
	params[ParamCompression] = compression
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	// date can only be a single value
	params[ParamDate] = date

	if isHTTPURL(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(reconciliationPath), params)
}

// ExportReconciledTransactions retrieves the collected transactions for a given day
// and sends it as a file to a destination which can be a HTTP URL or an email address.
// The file is compressed used the specified compression scheme, as specified
// in the be2bill.Compression* constants.
//
// The date can only be specified as a day (YYYY-MM-DD).
//
// See https://developer.be2bill.com/functions/exportReconciledTransactions
func (p *DirectLinkClient) ExportReconciledTransactions(date, destination, compression string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeExportReconciledTransactions
	params[ParamCompression] = compression
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	// date can only be a single value
	params[ParamDate] = date

	if isHTTPURL(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(reconciliationPath), params)
}
