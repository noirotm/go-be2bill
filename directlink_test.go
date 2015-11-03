// Copyright 2015 Dalenys. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"os"
	"regexp"
	"testing"
	"time"
)

func TestIsHttpUrl(t *testing.T) {
	cases := []struct {
		str      string
		expected bool
	}{
		{"test", false},
		{"http://toto.te/test", true},
		{"https://prolps.reez.com/", true},
		{"httpfoobar", false},
		{"mailto:tzze@zef.org", false},
	}

	for _, tc := range cases {
		result := isHTTPURL(tc.str)
		if result != tc.expected {
			t.Errorf("isHttpUrl: %s", tc.str)
			t.Errorf("want %+v, got %+v", tc.expected, result)
		}
	}
}

func setupSandboxClient(t *testing.T) *DirectLinkClient {
	if testing.Short() {
		t.Skip("skipping remote tests in short mode.")
	}

	// build client with identifiers from the environment
	identifier := os.Getenv("BE2BILL_IDENTIFIER")
	password := os.Getenv("BE2BILL_PASSWORD")

	if len(identifier) == 0 {
		t.Fatal("identifier not set")
	}

	if len(password) == 0 {
		t.Fatal("password not set")
	}

	return BuildSandboxDirectLinkClient(identifier, password)
}

func TestPayment(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().Add((365 + 30) * 24 * time.Hour)
	r, err := c.Payment(
		"1111222233334444",
		date.Format("01-06"),
		"123",
		"john doe",
		SingleAmount(100),
		"42",
		"ident",
		"test@test.com",
		"1.1.1.1",
		"desc",
		"Firefox",
		Options{},
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypePayment {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}

func TestAuthorization(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().Add((365 + 30) * 24 * time.Hour)
	r, err := c.Authorization(
		"1111222233334444",
		date.Format("01-06"),
		"123",
		"john doe",
		100,
		"42",
		"ident",
		"test@test.com",
		"1.1.1.1",
		"desc",
		"Firefox",
		Options{},
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypeAuthorization {
		t.Errorf("expected %s, got %s", OperationTypeAuthorization, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}

func TestOneClickPayment(t *testing.T) {
	c := setupSandboxClient(t)

	r, err := c.OneClickPayment(
		"A142429",
		SingleAmount(100),
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"onelick_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		Options{},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypePayment {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}

func TestRefund(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().Add((365 + 30) * 24 * time.Hour)
	r, err := c.Payment(
		"1111222233334444",
		date.Format("01-06"),
		"123",
		"john doe",
		SingleAmount(5000),
		"42",
		"ident",
		"test@test.com",
		"1.1.1.1",
		"desc",
		"Firefox",
		Options{},
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	r2, err := c.Refund(
		r.TransactionID(),
		"refund_transaction_test",
		"refund transaction test",
		Options{
			ParamAmount: SingleAmount(5000),
		},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r2.OperationType() != OperationTypeRefund {
		t.Errorf("expected %s, got %s", OperationTypeRefund, r.OperationType())
	}
	if r2.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r2.ExecCode(), r2.Message())
	}
	if r2.StringValue(ResultParamAmount) != "5000" {
		t.Errorf("expected %s, got %s", "5000", r2.StringValue(ResultParamAmount))
	}
	if r2.StringValue(ResultParamTransactionID) == "" {
		t.Error("empty transactionID")
	}
	if r2.Message() == "" {
		t.Error("empty message")
	}
}

func TestCapture(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().Add((365 + 30) * 24 * time.Hour)
	r, err := c.Authorization(
		"1111222233334444",
		date.Format("01-06"),
		"123",
		"john doe",
		100,
		"42",
		"ident",
		"test@test.com",
		"1.1.1.1",
		"desc",
		"Firefox",
		Options{},
	)

	r2, err := c.Capture(r.TransactionID(), "order_21", "Capture test 01", Options{})
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r2.OperationType() != OperationTypeCapture {
		t.Errorf("expected %s, got %s", OperationTypeCapture, r2.OperationType())
	}
	if r2.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r2.ExecCode(), r2.Message())
	}
	if r2.StringValue(ResultParamTransactionID) == "" {
		t.Error("empty transactionID")
	}
	if r2.Message() == "" {
		t.Error("empty message")
	}
}

func TestOneClickAuthorization(t *testing.T) {
	c := setupSandboxClient(t)

	r, err := c.OneClickAuthorization(
		"A142429",
		100,
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"onelick_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		Options{},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypeAuthorization {
		t.Errorf("expected %s, got %s", OperationTypeAuthorization, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}

func TestSubscriptionAuthorization(t *testing.T) {
	c := setupSandboxClient(t)

	r, err := c.SubscriptionAuthorization(
		"A142429",
		100,
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"subscription_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		Options{},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypeAuthorization {
		t.Errorf("expected %s, got %s", OperationTypeAuthorization, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}

func TestSubscriptionPayment(t *testing.T) {
	c := setupSandboxClient(t)

	r, err := c.SubscriptionPayment(
		"A142429",
		SingleAmount(100),
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"subscription_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		Options{},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypePayment {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}

func TestStopNTimes(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().Add((365 + 30) * 24 * time.Hour)

	amount := FragmentedAmount{}
	amount[time.Now().Format("2006-01-02")] = 5000
	amount[time.Now().Add(20*24*time.Hour).Format("2006-01-02")] = 5000
	amount[time.Now().Add(30*24*time.Hour).Format("2006-01-02")] = 5000

	r, err := c.Payment(
		"1111222233334444",
		date.Format("01-06"),
		"123",
		"john doe",
		amount,
		"42",
		"ident",
		"test@test.com",
		"1.1.1.1",
		"desc",
		"Firefox",
		Options{},
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.ExecCode() != ExecCodeSuccess {
		t.Fatalf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}

	r2, err := c.StopNTimes(r.TransactionID(), Options{})

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r2.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r2.ExecCode(), r2.Message())
	}
	if r2.Message() == "" {
		t.Error("empty message")
	}
}

func TestRedirectForPayment(t *testing.T) {
	t.Skip("special user account needed")

	c := setupSandboxClient(t)

	r, err := c.RedirectForPayment(
		SingleAmount(10000),
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"subscription_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		Options{},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.ExecCode() != ExecCodeSuccess {
		t.Fatalf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.OperationType() != OperationTypePayment {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.Message() == "" {
		t.Error("empty message")
	}

	html := r.StringValue(ResultParamRedirectHTML)
	re := regexp.MustCompile("^[A-Za-z0-9]+={0,3}")
	if !re.MatchString(html) {
		t.Errorf("want Base64 data, got %s", html)
	}
}
