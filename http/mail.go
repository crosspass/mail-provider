package http

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/toolkits/web/param"
	"gopkg.in/gomail.v2"
	"vnote.club/mail-provider/config"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)

		// s := smtp.NewSMTP(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.TLS, cfg.Smtp.Anonymous, cfg.Smtp.SkipVerify)
		// err := s.SendMail(cfg.Smtp.From, tos, subject, content)

		m := gomail.NewMessage()
		m.SetHeader("From", cfg.Smtp.From)
		m.SetHeader("To", tos)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", content)

		d := gomail.NewDialer(cfg.Smtp.Addr, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password)

		d.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.Smtp.SkipVerify}

		// Send the email to Bob, Cora and Dan.
		err := d.DialAndSend(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
