package be2bill

import "testing"

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
