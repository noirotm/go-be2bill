package be2bill

import (
	"os"
	"testing"
	"time"
)

func setupSandboxClient() DirectLinkClient {
	// build client with identifiers from the environment
	identifier := os.Getenv("BE2BILL_IDENTIFIER")
	password := os.Getenv("BE2BILL_PASSWORD")

	return BuildSandboxDirectLinkClient(identifier, password)
}

func TestPayment(t *testing.T) {
	c := setupSandboxClient()

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
		DefaultOptions,
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType != OperationTypePayment {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType)
	}

	if r.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}

func TestAuthorization(t *testing.T) {
	c := setupSandboxClient()

	date := time.Now().Add((365 + 30) * 24 * time.Hour)
	r, err := c.Authorization(
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
		DefaultOptions,
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType != OperationTypeAuthorization {
		t.Errorf("expected %s, got %s", OperationTypeAuthorization, r.OperationType)
	}

	if r.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}

func TestOneClickPayment(t *testing.T) {
	c := setupSandboxClient()

	r, err := c.OneClickPayment(
		"A142429",
		SingleAmount(100),
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"onelick_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		DefaultOptions,
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType != OperationTypePayment {
		t.Errorf("expected %s, got %s", OperationTypePayment, r.OperationType)
	}

	if r.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}

func TestRefund(t *testing.T) {
	c := setupSandboxClient()

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
		DefaultOptions,
	)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	r2, err := c.Refund(
		r.TransactionID,
		"refund_transaction_test",
		"refund transaction test",
		Options{
			ParamAmount: SingleAmount(5000),
		},
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r2.OperationType != OperationTypeRefund {
		t.Errorf("expected %s, got %s", OperationTypeRefund, r.OperationType)
	}

	if r2.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r2.ExecCode, r2.Message)
	}

	if r2.Amount != "5000" {
		t.Errorf("expected %s, got %s", "5000", r2.Amount)
	}
}

func TestCapture(t *testing.T) {
	c := setupSandboxClient()

	date := time.Now().Add((365 + 30) * 24 * time.Hour)
	r, err := c.Authorization(
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
		DefaultOptions,
	)

	r2, err := c.Capture(r.TransactionID, "order_21", "Capture test 01", DefaultOptions)
	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r2.OperationType != OperationTypeCapture {
		t.Errorf("expected %s, got %s", OperationTypeCapture, r2.OperationType)
	}

	if r2.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r2.ExecCode, r2.Message)
	}
}

func TestOneClickAuthorization(t *testing.T) {
	c := setupSandboxClient()

	r, err := c.OneClickAuthorization(
		"A142429",
		SingleAmount(100),
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"onelick_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		DefaultOptions,
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType != OperationTypeAuthorization {
		t.Errorf("expected %s, got %s", OperationTypeAuthorization, r.OperationType)
	}

	if r.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}

func TestSubscriptionAuthorization(t *testing.T) {
	c := setupSandboxClient()

	r, err := c.SubscriptionAuthorization(
		"A142429",
		SingleAmount(100),
		"order_1431181407",
		"6328_john.smith",
		"6328_john.smith@gmail.com",
		"123.123.123.123",
		"onelick_transaction",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36",
		DefaultOptions,
	)

	if err != nil {
		t.Fatal("got error: ", err)
	}

	if r.OperationType != OperationTypeAuthorization {
		t.Errorf("expected %s, got %s", OperationTypeAuthorization, r.OperationType)
	}

	if r.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}
