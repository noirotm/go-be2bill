package be2bill

type Credentials struct {
	identifier string
	password   string
	production bool
}

func User(identifier string, password string, production bool) *Credentials {
	return &Credentials{identifier, password, production}
}

func ProductionUser(identifier string, password string) *Credentials {
	return &Credentials{identifier, password, true}
}

func SandboxUser(identifier string, password string) *Credentials {
	return &Credentials{identifier, password, false}
}
