package config

type configer interface {
	AppName() string
	SQLDsn() string
}
