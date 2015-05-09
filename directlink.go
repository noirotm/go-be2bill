package be2bill

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type BasicResponse struct {
	OperationType string
	TransactionID string
	ExecCode      string
	Message       string
}

type TransactionResponse struct {
	BasicResponse
	Descriptor string
}

type RefundResponse struct {
	BasicResponse
	Amount string
}

type DirectLinkClient interface {
	Payment(cardPan, cardDate, cardCryptogram, cardFullName string, amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string, options Options) (*TransactionResponse, error)
	Authorization(cardPan, cardDate, cardCryptogram, cardFullName string, amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string, options Options) (*TransactionResponse, error)
	Credit(cardPan, cardDate, cardCryptogram, cardFullName string, amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string, options Options) (*TransactionResponse, error)
	OneClickPayment(alias string, amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string, options Options) (*TransactionResponse, error)
	Refund(transactionID, orderID, description string, options Options) (*RefundResponse, error)
	Capture(transactionID, orderID, description string, options Options) (*BasicResponse, error)
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

type directLinkClientImpl struct {
	credentials *Credentials
	urls        []string
	hasher      Hasher
}

func (p *directLinkClientImpl) getURLs(path string) []string {
	urls := make([]string, len(p.urls))
	for i, url := range p.urls {
		urls[i] = url + path
	}
	return urls
}

func (p *directLinkClientImpl) getDirectLinkURLs() []string {
	return p.getURLs(directLinkPath)
}

func (p *directLinkClientImpl) doPostRequest(url string, params Options) ([]byte, error) {
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

func (p *directLinkClientImpl) requests(urls []string, params Options, result interface{}) error {
	for _, url := range urls {
		buf, err := p.doPostRequest(url, params)

		if err != nil {
			// break if a timeout occured, otherwise try next URL
			if err == ErrTimeout {
				return err
			} else {
				continue
			}
		}

		// decode result
		return json.Unmarshal(buf, result)
	}

	// we can reach this statement only if the URLs slice is empty
	return ErrURLMissing
}

func (p *directLinkClientImpl) transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent string, options Options, result interface{}) error {
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

	return p.requests(p.getDirectLinkURLs(), params, result)
}

func (p *directLinkClientImpl) Payment(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options) (*TransactionResponse, error) {
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

	result := &TransactionResponse{}
	if err := p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *directLinkClientImpl) Authorization(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options) (*TransactionResponse, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeAuthorization
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName
	params[ParamAmount] = amount.(SingleAmount)

	result := &TransactionResponse{}
	if err := p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *directLinkClientImpl) Credit(
	cardPan, cardDate, cardCryptogram, cardFullName string,
	amount Amount,
	orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options) (*TransactionResponse, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypeCredit
	params[ParamCardCode] = cardPan
	params[ParamCardValidityDate] = cardDate
	params[ParamCardCVV] = cardCryptogram
	params[ParamCardFullName] = cardFullName
	params[ParamAmount] = amount.(SingleAmount)

	result := &TransactionResponse{}
	if err := p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *directLinkClientImpl) OneClickPayment(
	alias string,
	amount Amount, orderID, clientID, clientEmail, clientIP, description, clientUserAgent string,
	options Options) (*TransactionResponse, error) {
	params := options.copy()

	params[ParamOperationType] = OperationTypePayment
	params[ParamAlias] = alias
	params[ParamAliasMode] = AliasModeOneClick
	params[ParamAmount] = amount.(SingleAmount)

	result := &TransactionResponse{}
	if err := p.transaction(orderID, clientID, clientEmail, clientIP, description, clientUserAgent, params, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *directLinkClientImpl) Refund(transactionID, orderID, description string, options Options) (*RefundResponse, error) {
	params := options.copy()

	params[ParamIdentifier] = p.credentials.identifier
	params[ParamOperationType] = OperationTypeRefund
	params[ParamDescription] = description
	params[ParamTransactionID] = transactionID
	params[ParamVersion] = APIVersion
	params[ParamOrderID] = orderID

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	result := &RefundResponse{}
	if err := p.requests(p.getDirectLinkURLs(), params, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *directLinkClientImpl) Capture(transactionID, orderID, description string, options Options) (*BasicResponse, error) {
	params := options.copy()

	params[ParamIdentifier] = p.credentials.identifier
	params[ParamOperationType] = OperationTypeCapture
	params[ParamVersion] = APIVersion
	params[ParamDescription] = description
	params[ParamTransactionID] = transactionID
	params[ParamOrderID] = orderID

	params[ParamHash] = p.hasher.ComputeHash(p.credentials.password, params)

	result := &BasicResponse{}
	if err := p.requests(p.getDirectLinkURLs(), params, result); err != nil {
		return nil, err
	}

	return result, nil
}
