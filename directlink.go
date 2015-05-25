// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

const (
	searchByOrderID       = "ORDERID"
	searchByTransactionID = "TRANSACTIONID"
)

type Result map[string]interface{}

func (r Result) StringValue(name string) string {
	val, ok := r[name]
	if !ok {
		return ""
	}
	return val.(string)
}

func (r Result) OperationType() string {
	return r.StringValue(ResultParamOperationType)
}

func (r Result) ExecCode() string {
	return r.StringValue(ResultParamExecCode)
}

func (r Result) Message() string {
	return r.StringValue(ResultParamMessage)
}

func (r Result) TransactionID() string {
	return r.StringValue(ResultParamTransactionID)
}

func (r Result) Success() bool {
	return r.ExecCode() == ExecCodeSuccess
}

const (
	directLinkPath     = "/front/service/rest/process"
	exportPath         = "/front/service/rest/export"
	reconciliationPath = "/front/service/rest/reconciliation"

	requestTimeout = 30 * time.Second
)

var (
	ErrTimeout    = errors.New("timeout")
	ErrURLMissing = errors.New("no URL provided")
)

type DirectLinkClient struct {
	credentials *Credentials
	urls        []string
	hasher      Hasher
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

func (p *DirectLinkClient) doPostRequest(url string, params Options) ([]byte, error) {
	requestParams := Options{
		"method": params[ParamOperationType],
		"params": params,
	}

	responseChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		resp, err := http.PostForm(url, requestParams.urlValues())
		if err != nil {
			errChan <- err
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			errChan <- err
			return
		}

		responseChan <- body
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-time.After(requestTimeout):
		return nil, ErrTimeout
	case response := <-responseChan:
		return response, nil
	}
}

func (p *DirectLinkClient) requests(urls []string, params Options) (Result, error) {
	for _, url := range urls {
		buf, err := p.doPostRequest(url, params)

		if err != nil {
			// break if a timeout occured, otherwise try next URL
			if err == ErrTimeout {
				return nil, err
			} else {
				continue
			}
		}

		// decode result
		result := make(Result)
		err = json.Unmarshal(buf, &result)

		return result, err
	}

	// we can reach this statement only if the URLs slice is empty
	return nil, ErrURLMissing
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

func isHttpUrl(str string) bool {
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
	if isHttpUrl(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(exportPath), params)
}

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

func (p *DirectLinkClient) Authorization(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

func (p *DirectLinkClient) Credit(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeCredit
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

func (p *DirectLinkClient) OneClickPayment(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypePayment
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeOneClick
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

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

func (p *DirectLinkClient) OneClickAuthorization(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeOneClick
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

func (p *DirectLinkClient) SubscriptionAuthorization(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeSubscription
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

func (p *DirectLinkClient) SubscriptionPayment(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypePayment
	params[ParamAlias] = alias
	params[ParamAliasMode] = aliasModeSubscription
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

func (p *DirectLinkClient) StopNTimes(scheduleID string, options Options) (Result, error) {
	params := options.copy()

	params[ParamIdentifier] = p.credentials.identifier
	params[ParamOperationType] = OperationTypeStopNTimes
	params[ParamScheduleID] = scheduleID
	params[ParamVersion] = APIVersion

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getDirectLinkURLs(), params)
}

func (p *DirectLinkClient) RedirectForPayment(
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options,
) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypePayment
	params[ParamAmount] = amount.(SingleAmount)

	return p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params)
}

func (p *DirectLinkClient) GetTransactionsByTransactionID(transactionIDs []string, destination, compression string) (Result, error) {
	return p.getTransactions(searchByTransactionID, transactionIDs, destination, compression)
}

func (p *DirectLinkClient) GetTransactionsByOrderID(orderIDs []string, destination, compression string) (Result, error) {
	return p.getTransactions(searchByOrderID, orderIDs, destination, compression)
}

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

	if isHttpUrl(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(exportPath), params)
}

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

	if isHttpUrl(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(exportPath), params)
}

func (p *DirectLinkClient) ExportReconciliation(startDate, endDate, destination, compression string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeExportReconciliation
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

	if isHttpUrl(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(reconciliationPath), params)
}

func (p *DirectLinkClient) ExportReconciledTransactions(date, destination, compression string, options Options) (Result, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeExportReconciledTransactions
	params[ParamCompression] = compression
	params[ParamIdentifier] = p.credentials.identifier
	params[ParamVersion] = APIVersion

	// date can only be a single value
	params[ParamDate] = date

	if isHttpUrl(destination) {
		params[ParamCallbackURL] = destination
	} else {
		params[ParamMailTo] = destination
	}

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	return p.requests(p.getURLs(reconciliationPath), params)
}
