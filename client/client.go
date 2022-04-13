package client

type logger interface {
	Infof(string, ...any)
	Errorf(string, ...any)
}
