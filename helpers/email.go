package helpers

import (
	"errors"
	"html/template"
	"path/filepath"
	"raihpeduli/config"
	"runtime"
	"strconv"
	"strings"

	"github.com/wneessen/go-mail"
)

func EmailService(fullname, email, otp, status string) error {
	user := config.LoadSMTPConfig().SMTP_USER
	password := config.LoadSMTPConfig().SMTP_PASS
	port := config.LoadSMTPConfig().SMTP_PORT
	convPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	m := mail.NewMsg()
	if err := m.From(user); err != nil {
		return err
	}
	if err := m.To(email); err != nil {
		return err
	}

	m.Subject("Verifikasi Raih Peduli: Kode Keamanan")
	emailTemplate := struct {
		Otp    string
		Fullname string
	}{
		Otp:    otp,
		Fullname: fullname,
	}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("Failed to get the current file path")
	}

	var templateEmail string

	if status == "1" {
		templateEmail = "email-template.html"
	}else{
		templateEmail = "email-template-forgot.html"
	}

	templatePath := filepath.Join(filepath.Dir(filename), templateEmail)
	tmpl, err := template.New("emailTemplate").ParseFiles(templatePath)
	if err != nil {
		return err
	}
	var bodyContent strings.Builder
	if err := tmpl.ExecuteTemplate(&bodyContent, templateEmail, emailTemplate); err != nil {
		return err
	}
	m.SetBodyString(mail.TypeTextHTML, bodyContent.String())

	c, err := mail.NewClient("smtp.gmail.com", mail.WithPort(convPort), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(user), mail.WithPassword(password))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(m); err != nil {
		return err
	}
	return nil
}