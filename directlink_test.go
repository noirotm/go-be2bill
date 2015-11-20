// Copyright 2015 Dalenys. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
)

func requestParameters(form url.Values) Options {
	params := make(Options)

	// parse every params[*] key
	re := regexp.MustCompile(`^params\[([^\]]+)\](?:\[([^\]]+)\])?$`)
	for name, value := range form {
		s := re.FindStringSubmatch(name)

		if len(s) < 3 {
			continue
		}

		if s[2] == "" {
			params[s[1]] = value[0]
		} else {
			params[s[1]] = Options{s[2]: value[0]}
		}
	}

	return params
}

func checkParams(params Options, t *testing.T) {
	// check hash
	if ok := CheckHash(&defaultHasher{}, "bar", params); !ok {
		t.Error("invalid hash")
	}

	// check for VERSION field
	ver, ok := params[ParamVersion]
	if !ok {
		t.Error("missing VERSION parameter")
	}
	if ver != APIVersion {
		t.Errorf("invalid VERSION parameter, want %s, got %s", APIVersion, ver)
	}
}

func TestResult(t *testing.T) {
	result := &Result{}
	if result.StringValue("foo") != "" {
		t.Error("non-empty result")
	}

	result = &Result{ResultParamExecCode: ExecCodeSuccess}
	if !result.Success() {
		t.Error("invalid success status")
	}
}

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

func TestEmptyEnvironment(t *testing.T) {
	env := Environment{}

	c := NewDirectLinkClient(User("foo", "bar", env))
	date := time.Now().AddDate(1, 1, 0)
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
	if err != ErrURLMissing {
		t.Errorf("got error: %v", err)
	}
	if r != nil {
		t.Error("r should be nil")
	}
}

func TestServerFallback(t *testing.T) {
	// first server, returns error 500 immediatly
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer ts.Close()
	// second server, normal operation
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"OPERATIONTYPE":"payment","TRANSACTIONID":"ABCDE01","EXECCODE":"0000","MESSAGE":"ok","DESCRIPTOR":"descr"}`)
	}))
	defer ts2.Close()

	env := Environment{ts.URL, ts2.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	date := time.Now().AddDate(1, 1, 0)
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

func TestConnectionError(t *testing.T) {
	// arbitrary url
	env := Environment{"http://127.0.0.1:61256"}

	c := NewDirectLinkClient(User("foo", "bar", env))
	date := time.Now().AddDate(1, 1, 0)
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
	if err == nil {
		t.Error("err should not be nil")
	}
	if r != nil {
		t.Error("r should be nil")
	}
}

func TestServerError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	})

	// first server, returns error 500
	ts := httptest.NewServer(h)
	defer ts.Close()
	// second server, returns error 500
	ts2 := httptest.NewServer(h)
	defer ts2.Close()

	env := Environment{ts.URL, ts2.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	date := time.Now().AddDate(1, 1, 0)
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
	if err != ErrServerError {
		t.Errorf("got error: %v", err)
	}
	if r != nil {
		t.Error("r should be nil")
	}
}

func TestServerTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long running test in short mode.")
	}

	// test server that replies after allowed timeout
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprint(w, `{"OPERATIONTYPE":"payment","TRANSACTIONID":"ABCDE01","EXECCODE":"0000","MESSAGE":"ok","DESCRIPTOR":"descr"}`)
	}))
	defer ts.Close()

	env := Environment{ts.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	c.RequestTimeout = 2 * time.Second

	date := time.Now().AddDate(1, 1, 0)
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
	if err != ErrTimeout {
		t.Errorf("got error: %v", err)
	}
	if r != nil {
		t.Error("r should be nil")
	}
}

func TestServerCloseConnection(t *testing.T) {
	// test server that forcefully closes all connections when handling a request
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// write partial response
		fmt.Fprint(w, `{"OPERATIONTYPE":"payment","TRANSA`)
		// close connections to simulate network failure
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	env := Environment{ts.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	c.RequestTimeout = 2 * time.Second

	date := time.Now().AddDate(1, 1, 0)
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
	if err == nil {
		t.Error("err should not be nil")
	}
	if r != nil {
		t.Error("r should be nil")
	}
}

func TestInvalidJSON(t *testing.T) {
	// test server that sends a PHP-like error message instead of JSON
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<b>Fatal error</b>: Uncaught exception 'Exception' with message 'Unable to connect to SQL server'`)
	}))
	defer ts.Close()

	env := Environment{ts.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	c.RequestTimeout = 2 * time.Second

	date := time.Now().AddDate(1, 1, 0)
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
	if err == nil {
		t.Error("err should not be nil")
	}
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("err should be a json.SyntaxError, got %+v", err)
	}
	if r != nil {
		t.Error("r should be nil")
	}
}

func TestRedirectForPayment(t *testing.T) {
	htmlCode := []byte(`<a href="http://example.org/">Link</a>`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b64 := base64.StdEncoding.EncodeToString(htmlCode)
		fmt.Fprintf(w, `{"OPERATIONTYPE":"payment","TRANSACTIONID":"ABCDE01","EXECCODE":"0002","MESSAGE":"ok","REDIRECTHTML":"%s"}`, b64)
	}))
	defer ts.Close()

	env := Environment{ts.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	r, err := c.RedirectForPayment(
		10000,
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
	if r.TransactionID() == "" {
		t.Error("empty transactionID")
	}
	if r.ExecCode() != ExecCodeAlternateRedirectRequired {
		t.Fatalf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.Message() == "" {
		t.Error("empty message")
	}

	str := r.StringValue(ResultParamRedirectHTML)
	if str == "" {
		t.Error("empty HTML code")
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		t.Errorf("invalid base64 data: %s", str)
	}

	if !bytes.Equal(data, htmlCode) {
		t.Error("invalid HTML code")
	}
}

func TestCredit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// test URI
		if r.URL.Path != directLinkPath {
			t.Errorf("invalid URI: %s", r.URL.Path)
		}

		// test method
		if r.Method != "POST" {
			t.Errorf("invalid HTTP method: %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}

		// check for METHOD field
		method, ok := r.Form["method"]
		if !ok {
			t.Error("missing method")
		}
		if method[0] != OperationTypeCredit {
			t.Errorf("invalid method: %s", method)
		}

		params := requestParameters(r.Form)

		// check request
		checkParams(params, t)

		fmt.Fprint(w, `{"OPERATIONTYPE":"credit","TRANSACTIONID":"ABCDE01","EXECCODE":"0000","MESSAGE":"ok","DESCRIPTOR":"descr"}`)
	}))
	defer ts.Close()

	env := Environment{ts.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))
	date := time.Now().AddDate(1, 1, 0)
	r, err := c.Credit(
		"1111222233334444",
		date.Format("01-06"),
		"123",
		"john doe",
		10000,
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

	if r.OperationType() != OperationTypeCredit {
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

func TestGetTransactions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// test URI
		if r.URL.Path != exportPath {
			t.Errorf("invalid URI: %s", r.URL.Path)
		}

		// test method
		if r.Method != "POST" {
			t.Errorf("invalid HTTP method: %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}

		// check for METHOD field
		method, ok := r.Form["method"]
		if !ok {
			t.Error("missing method")
		}
		if method[0] != OperationTypeGetTransactions {
			t.Errorf("invalid method: %s", method)
		}

		params := requestParameters(r.Form)

		// check request
		checkParams(params, t)

		fmt.Fprint(w, `{"OPERATIONTYPE":"getTransactions","EXECCODE":"0000","MESSAGE":"ok","DESCRIPTOR":"descr"}`)
	}))
	defer ts.Close()

	env := Environment{ts.URL}

	c := NewDirectLinkClient(User("foo", "bar", env))

	r, err := c.GetTransactionsByTransactionID(
		[]string{"1", "2", "3"},
		"exports@example.org",
		CompressionZip,
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypeGetTransactions {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}

	r, err = c.GetTransactionsByOrderID(
		[]string{"1", "2", "3"},
		"http://example.org/transaction.php",
		CompressionZip,
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType() != OperationTypeGetTransactions {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType())
	}
	if r.ExecCode() != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode(), r.Message())
	}
	if r.Message() == "" {
		t.Error("empty message")
	}
	if r.StringValue(ResultParamDescriptor) == "" {
		t.Error("empty descriptor")
	}
}
