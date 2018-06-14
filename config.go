package mail4go

import (
	"crypto/tls"
)

////////////////////////////////////////////////////////////////////////////////
type MailConfig struct {
	username string
	password string
	host     string
	port     string
	TLS      *tls.Config
}

func NewMailConfig(username string, password string, host string, port string) *MailConfig {
	var config = &MailConfig{}
	config.username = username
	config.password = password
	config.host = host
	config.port = port
	return config
}

func (this *MailConfig) Address() string {
	return this.host + ":" + this.port
}
