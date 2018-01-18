package mail4go

import (
	"errors"
	"net/mail"
	"net/smtp"
	"crypto/tls"
	"net"
)

////////////////////////////////////////////////////////////////////////////////
type MailConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Secure   bool
}

func NewMailConfig(username string, password string, host string, port string, secure bool) *MailConfig {
	var config = &MailConfig{}
	config.Username = username
	config.Password = password
	config.Host = host
	config.Port = port
	config.Secure = secure
	return config
}

func (this *MailConfig) Address() string {
	return this.Host + ":" + this.Port
}

////////////////////////////////////////////////////////////////////////////////
func SendMail(config *MailConfig, m *Message) error {
	if config == nil {
		return errors.New("config 不能为空")
	}
	if len(m.From) == 0 {
		m.From = config.Username
	}

	return send(config.Address(), smtp.PlainAuth("", config.Username, config.Password, config.Host), m, config.Secure)
}

// Send an email using the given host and SMTP auth (optional), returns any error thrown by smtp.SendMail
// This function merges the To, Cc, and Bcc fields and calls the smtp.SendMail function using the Email.Bytes() output as the message
func send(addr string, a smtp.Auth, m *Message, secure bool) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(m.To)+len(m.Cc)+len(m.Bcc))
	to = append(append(append(to, m.To...), m.Cc...), m.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if m.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	from, err := mail.ParseAddress(m.From)
	if err != nil {
		return err
	}
	raw, err := m.Bytes()
	if err != nil {
		return err
	}

//	if secure {
	return sendMail(addr, a, from.Address, to, raw, secure)
//	}
//	return smtp.SendMail(addr, a, from.Address, to, raw)
}

func sendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte, secure bool) error {

	host, _, _ := net.SplitHostPort(addr)

	var conn net.Conn
	var err error
	if secure {
		tlsConfig := &tls.Config {
			InsecureSkipVerify: true,
			ServerName: host,
		}
		conn, err = tls.Dial("tcp", addr, tlsConfig)
	} else {
//		conn, err = net.Dial("tcp", addr)
		return smtp.SendMail(addr, auth, from, to, msg)
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	defer c.Close()

	if err != nil {
		return err
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}