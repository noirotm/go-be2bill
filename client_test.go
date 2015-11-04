package be2bill

import "testing"

func TestCredentials(t *testing.T) {
	user := User("foo", "bar", EnvSandbox)
	if user.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", user.identifier)
	}
	if user.password != "bar" {
		t.Errorf("unexpected password: %s", user.password)
	}
	if user.environment[0] != EnvSandbox[0] {
		t.Errorf("unexpected environment: %v", user.environment)
	}

	user = SandboxUser("foo", "bar")
	if user.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", user.identifier)
	}
	if user.password != "bar" {
		t.Errorf("unexpected password: %s", user.password)
	}
	if user.environment[0] != EnvSandbox[0] {
		t.Errorf("unexpected environment: %v", user.environment)
	}

	user = ProductionUser("foo", "bar")
	if user.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", user.identifier)
	}
	if user.password != "bar" {
		t.Errorf("unexpected password: %s", user.password)
	}
	if user.environment[0] != EnvProduction[0] {
		t.Errorf("unexpected environment: %v", user.environment)
	}
}

func TestFormClient(t *testing.T) {
	user := SandboxUser("foo", "bar")

	client := NewFormClient(user)
	if client.credentials != user {
		t.Errorf("unexpected credentials: %s", client.credentials)
	}

	client = BuildSandboxFormClient("foo", "bar")
	if client.credentials.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", client.credentials.identifier)
	}
	if client.credentials.password != "bar" {
		t.Errorf("unexpected password: %s", client.credentials.password)
	}
	if client.credentials.environment[0] != EnvSandbox[0] {
		t.Errorf("unexpected environment: %v", client.credentials.environment)
	}

	client = BuildProductionFormClient("foo", "bar")
	if client.credentials.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", client.credentials.identifier)
	}
	if client.credentials.password != "bar" {
		t.Errorf("unexpected password: %s", client.credentials.password)
	}
	if client.credentials.environment[0] != EnvProduction[0] {
		t.Errorf("unexpected environment: %v", client.credentials.environment)
	}
}

func TestDirectLinkClient(t *testing.T) {
	user := SandboxUser("foo", "bar")

	client := NewDirectLinkClient(user)
	if client.credentials != user {
		t.Errorf("unexpected credentials: %s", client.credentials)
	}

	client = BuildSandboxDirectLinkClient("foo", "bar")
	if client.credentials.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", client.credentials.identifier)
	}
	if client.credentials.password != "bar" {
		t.Errorf("unexpected password: %s", client.credentials.password)
	}
	if client.credentials.environment[0] != EnvSandbox[0] {
		t.Errorf("unexpected environment: %v", client.credentials.environment)
	}

	client = BuildProductionDirectLinkClient("foo", "bar")
	if client.credentials.identifier != "foo" {
		t.Errorf("unexpected identifier: %s", client.credentials.identifier)
	}
	if client.credentials.password != "bar" {
		t.Errorf("unexpected password: %s", client.credentials.password)
	}
	if client.credentials.environment[0] != EnvProduction[0] {
		t.Errorf("unexpected environment: %v", client.credentials.environment)
	}
}

func TestSwitchURLs(t *testing.T) {
	if EnvProduction[0] != "https://secure-magenta1.be2bill.com" {
		t.Errorf("unexpected primary production URL: %s", EnvProduction[0])
	}
	if EnvProduction[1] != "https://secure-magenta2.be2bill.com" {
		t.Errorf("unexpected secondary production URL: %s", EnvProduction[1])
	}

	EnvProduction.SwitchURLs()

	if EnvProduction[0] != "https://secure-magenta2.be2bill.com" {
		t.Errorf("unexpected primary production URL: %s", EnvProduction[0])
	}
	if EnvProduction[1] != "https://secure-magenta1.be2bill.com" {
		t.Errorf("unexpected secondary production URL: %s", EnvProduction[1])
	}

	// restore normal order
	EnvProduction.SwitchURLs()

	if EnvProduction[0] != "https://secure-magenta1.be2bill.com" {
		t.Errorf("unexpected primary production URL: %s", EnvProduction[0])
	}
	if EnvProduction[1] != "https://secure-magenta2.be2bill.com" {
		t.Errorf("unexpected secondary production URL: %s", EnvProduction[1])
	}
}
