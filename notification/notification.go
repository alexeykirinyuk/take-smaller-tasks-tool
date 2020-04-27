package notification

import (
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/config"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/history"
	"net/smtp"
	"strconv"
)

// Notificator is a structure for sending notifications
type Notificator struct {
	config config.SMTPConfiguration
}

func CreteNotificator(configuration config.SMTPConfiguration) *Notificator {
	return &Notificator{config: configuration}
}

// SendNotification is a method for sending notifications
func (n *Notificator) Notify(h *history.History) error {
	html, err := h.Html()
	if err != nil {
		return err
	}

	c := n.config
	auth := createLoginAuth(c.UserName, c.Password)

	subject := "Subject: Key Result: Take Smaller Tasks\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := []byte(subject + mime + html)

	err = smtp.SendMail(
		c.Domain+":"+strconv.Itoa(c.Port),
		auth,
		c.UserName,
		[]string{c.UserName},
		body)

	if err != nil {
		return fmt.Errorf("Error when send email: %s", err)
	}

	return nil
}
