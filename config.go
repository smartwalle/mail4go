package mail4go

import (
	"crypto/tls"
	"net/smtp"
)

////////////////////////////////////////////////////////////////////////////////
type MailConfig struct {
	username string
	password string
	host     string
	port     string
	TLS      *tls.Config
	auth     smtp.Auth
}

func NewMailConfig(username string, password string, host string, port string) *MailConfig {
	var config = &MailConfig{}
	config.username = username
	config.password = password
	config.host = host
	config.port = port
	config.auth = smtp.PlainAuth("", username, password, host)
	return config
}

func (this *MailConfig) Address() string {
	return this.host + ":" + this.port
}