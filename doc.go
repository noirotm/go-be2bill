// Copyright 2015 Dalenys. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package be2bill implements client access to the Be2bill merchant API defined
at http://developer.be2bill.com/API.

Form Client

Every initial transaction is made using a secure form.
The Form Client API exposes methods that return the HTML code
for payment or authorization buttons.

The first step is to build a client for the given environment, using
your credentials:

	client := be2bill.BuildSandboxFormClient("test", "password")

To build a payment form button, call:

	button := client.BuildPaymentFormButton(
		be2bill.SingleAmount(15235),   // amount
		"order_1412327697",            // order ID
		"6328_john.smith@example.org", // user ID
		"Fashion jacket",              // order description
		be2bill.Options{
			be2bill.HTMLOptionSubmit: be2bill.Options{
				"value": "Pay with be2bill",
				"class": "flatButton",
			},
			be2bill.HTMLOptionForm: be2bill.Options{"id": "myform"},
		},                             // HTML options for the form
		be2bill.Options{
			be2bill.ParamClientEmail: "john.smith@example.org",
			be2bill.Param3DSecure:    "yes",
		},                             // additional platform options
	)

Authorization form buttons are created similarily, except that the
method to call is `BuildAuthorizationFormButton` that takes the same
parameters.
An authorization must be captured using the `Capture` method of the
Direct Link Client API.

Direct Link Client

All operations that do not require interactive data input from the client
can be made using HTTP POST calls to the be2bill servers.
The Direct Link Client API is an abstraction layer for these calls.

Just like the Form Client API, the first step is to build a client for
the given environment, using your credentials:

	client := be2bill.BuildSandboxDirectLinkClient("test", "password")

Then, for example to capture a previously authorized transaction, call:

	result, err := client.Capture(
		"A151621",
		"order_1423675675",
		"capture_transaction_A151621",
		be2bill.Options{},
	)

Please note that access to the Direct Link Client API is not enabled by default.
This service can only be activated by your account manager based on specific
criteria. Please contact him or the support team for more information.
*/
package be2bill
