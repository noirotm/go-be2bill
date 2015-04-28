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
	ParamTimeZone            = "TIMEZONE"
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
