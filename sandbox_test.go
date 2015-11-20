package be2bill

import (
	"os"
	"testing"
	"time"
)

func setupSandboxClient(t *testing.T) *DirectLinkClient {
	if testing.Short() {
		t.Skip("skipping remote tests in short mode")
	}

	// build client with identifiers from the environment
	identifier := os.Getenv("BE2BILL_IDENTIFIER")
	password := os.Getenv("BE2BILL_PASSWORD")

	if len(identifier) == 0 {
		t.Skip("BE2BILL_IDENTIFIER environment variable not set")
	}

	if len(password) == 0 {
		t.Skip("BE2BILL_PASSWORD environment variable not set")
	}

	return BuildSandboxDirectLinkClient(identifier, password)
}

func TestSandboxPayment(t *testing.T) {
	c := setupSandboxClient(t)

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

func TestSandboxAuthorization(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().AddDate(1, 1, 0)
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

func TestSandboxOneClickPayment(t *testing.T) {
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

func TestSandboxOneClickPaymentFragmented(t *testing.T) {
	c := setupSandboxClient(t)

	a := make(FragmentedAmount)
	a[time.Now().Format("2006-01-02")] = 15235
	a[time.Now().AddDate(0, 1, 0).Format("2006-01-02")] = 14723

	r, err := c.OneClickPayment(
		"A142429",
		a,
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

func TestSandboxRefund(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().AddDate(1, 1, 0)
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

func TestSandboxCapture(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().AddDate(1, 1, 0)
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

func TestSandboxOneClickAuthorization(t *testing.T) {
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

func TestSandboxSubscriptionAuthorization(t *testing.T) {
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

func TestSandboxSubscriptionPayment(t *testing.T) {
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

func TestSandboxSubscriptionPaymentFragmented(t *testing.T) {
	c := setupSandboxClient(t)

	a := make(FragmentedAmount)
	a[time.Now().Format("2006-01-02")] = 15235
	a[time.Now().AddDate(0, 1, 0).Format("2006-01-02")] = 14723

	r, err := c.SubscriptionPayment(
		"A142429",
		a,
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

func TestSandboxStopNTimes(t *testing.T) {
	c := setupSandboxClient(t)

	date := time.Now().AddDate(1, 1, 0)

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
