package common

type Server interface {
	Run() error
}

type ServerBase struct {
	Server
	IP   string
	Port int
}
