package mail4go

import (
	"crypto/tls"
	"errors"
	"net/mail"
	"net/smtp"
)

type sender func(addr string, a smtp.Auth, t *tls.Config, m *Message) error

func SendWithConfig(config *MailConfig, m *Message) error {
	if config == nil {
		return errors.New("config 不能为空")
	}
	if m.From == "" {
		if config.From != "" {
			m.From = config.From
		} else {
			m.From = config.username
		}
	}
	if config.TLS != nil {
		if config.StartTLS {
			return SendWithStartTLS(config.username, config.password, config.host, config.port, config.TLS, m)
		}
		return SendWithTLS(config.username, config.password, config.host, config.port, config.TLS, m)
	}
	return Send(config.username, config.password, config.host, config.port, m)
}

func _send(username, password, host, port string, t *tls.Config, m *Message, s sender) error {
	if m.From == "" {
		m.From = username
	}
	var auth = smtp.PlainAuth("", username, password, host)
	var addr = host + ":" + port
	return s(addr, auth, t, m)
}

func Send(username, password, host, port string, m *Message) error {
	return _send(username, password, host, port, nil, m, send)
}

func send(addr string, a smtp.Auth, t *tls.Config, m *Message) error {
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
	sender, err := m.parseSender()
	if err != nil {
		return err
	}
	raw, err := m.Bytes()
	if err != nil {
		return err
	}
	return smtp.SendMail(addr, a, sender, to, raw)
}

func SendWithTLS(username, password, host, port string, t *tls.Config, m *Message) error {
	return _send(username, password, host, port, t, m, sendWithTLS)
}

func sendWithTLS(addr string, a smtp.Auth, t *tls.Config, m *Message) error {
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
	sender, err := m.parseSender()
	if err != nil {
		return err
	}
	raw, err := m.Bytes()
	if err != nil {
		return err
	}

	conn, err := tls.Dial("tcp", addr, t)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, t.ServerName)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(sender); err != nil {
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
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func SendWithStartTLS(username, password, host, port string, t *tls.Config, m *Message) error {
	return _send(username, password, host, port, t, m, sendWithStartTLS)
}

func sendWithStartTLS(addr string, a smtp.Auth, t *tls.Config, m *Message) error {
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
	sender, err := m.parseSender()
	if err != nil {
		return err
	}
	raw, err := m.Bytes()
	if err != nil {
		return err
	}

	// Taken from the standard library
	// https://github.com/golang/go/blob/master/src/net/smtp/smtp.go#L328
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}
	// Use TLS if available
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(t); err != nil {
			return err
		}
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(sender); err != nil {
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
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
