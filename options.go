package be2bill

import (
	"fmt"
	"net/url"
	"sort"
)

const (
	HTMLOptionForm   = "FORM"
	HTMLOptionSubmit = "SUBMIT"

	Param3DSecure            = "3DSECURE"
	Param3DSecureDisplayMode = "3DSECUREDISPLAYMODE"
	ParamAmount              = "AMOUNT"
	ParamAmounts             = "AMOUNTS"
	ParamBillingAddress      = "BILLINGADDRESS"
	ParamBillingCountry      = "BILLINGCOUNTRY"
	ParamBillingFirstName    = "BILLINGFIRSTNAME"
	ParamBillingLastName     = "BILLINGLASTNAME"
	ParamBillingPhone        = "BILLINGPHONE"
	ParamBillingPostalCode   = "BILLINGPOSTALCODE"
	ParamCardCode            = "CARDCODE"
	ParamCardCVV             = "CARDCVV"
	ParamCardFullName        = "CARDFULLNAME"
	ParamCardValidityDate    = "CARDVALIDITYDATE"
	ParamClientAddress       = "CLIENTADDRESS"
	ParamClientDOB           = "CLIENTDOB"
	ParamClientEmail         = "CLIENTEMAIL"
	ParamClientIdent         = "CLIENTIDENT"
	ParamClientIP            = "CLIENTIP"
	ParamClientUserAgent     = "CLIENTUSERAGENT"
	ParamCreateAlias         = "CREATEALIAS"
	ParamDescription         = "DESCRIPTION"
	ParamDisplayCreateAlias  = "DISPLAYCREATEALIAS"
	ParamExtraData           = "EXTRADATA"
	ParamHash                = "HASH"
	ParamHideCardFullName    = "HIDECARDFULLNAME"
	ParamHideClientEmail     = "HIDECLIENTEMAIL"
	ParamIdentifier          = "IDENTIFIER"
	ParamLanguage            = "LANGUAGE"
	ParamOperationType       = "OPERATIONTYPE"
	ParamOrderID             = "ORDERID"
	ParamShipToAddress       = "SHIPTOADDRESS"
	ParamShipToFirstName     = "SHIPTOFIRSTNAME"
	ParamShipToLastName      = "SHIPTOLASTNAME"
	ParamShipToPhone         = "SHIPTOPHONE"
	ParamShipToPostalCode    = "SHIPTOPOSTALCODE"
	ParamSubmit              = "SUBMIT"
	ParamTransactionID       = "TRANSACTIONID"
	ParamVersion             = "VERSION"
	ParamVME                 = "VME"
)

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

func (p Options) urlValues() url.Values {
	values := url.Values{}
	for k, v := range p {
		if opts, ok := v.(Options); ok {
			for subkey, subval := range opts {
				values.Set(fmt.Sprint(k, "[", subkey, "]"), fmt.Sprint(subval))
			}
		} else {
			values.Set(k, fmt.Sprint(v))
		}
	}

	return values
}

var DefaultOptions = make(Options)
