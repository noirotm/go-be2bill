package be2bill

import "testing"

func setupSandboxClient() DirectLinkClient {
	// build client
	return BuildSandboxDirectLinkClient("IDENTIFIER", "PASSWORD")
}

func TestAuthorization(t *testing.T) {
	c := setupSandboxClient()

	r, err := c.Authorization(
		"1111222233334444",
		"01-12",
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
		t.Fatal("Got error: ", err)
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

	r, err := c.Capture("test1", "order_21", "Capture test 01", DefaultOptions)
	if err != nil {
		t.Fatal("Got error: ", err)
	}

	if r.OperationType != OperationTypeCapture {
		t.Errorf("expected %s, got %s", OperationTypeCapture, r.OperationType)
	}

	if r.ExecCode != ExecCodeSuccess {
		t.Errorf("exec code %s, message: %s", r.ExecCode, r.Message)
	}
}
