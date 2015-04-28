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
		t.Errorf("expected %s, got %s", OperationTypeCapture, r.OperationType)
	}

	if r2.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}
