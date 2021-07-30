package models

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/badoux/checkmail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	gomail "gopkg.in/mail.v2"
)
  
type SendMail struct {
	Email string
}

func (sm *SendMail) SendEmail(password string) error {

	if len(sm.Email) < 1 {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(sm.Email); err != nil {
		return errors.New("Invalid Email")
	}

	templateData := struct {
		Password string
	 }{
		Password: password,
	 }

	m := gomail.NewMessage()
  
	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("SYSTEM_EMAIL"))
  
	// Set E-Mail receivers
	m.SetHeader("To", sm.Email)
  
	// Set E-Mail subject
	m.SetHeader("Subject", "Set/Reset User Login Password")

	var err error
	t, err := template.ParseFiles("./html/email_template.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, templateData); err != nil {
		return err
	}
	result := tpl.String()
  
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", result)

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
    if err != nil {
        return err
    }
  
	// Settings for SMTP server
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SYSTEM_EMAIL"), os.Getenv("SYSTEM_PASSWORD"))
  
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
  
	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
	  fmt.Println(err)
	}
  
	return nil
}

func (sm *SendMail) SendGridMail() error {
	from := mail.NewEmail(os.Getenv("SYSTEM_ADMIN_USERNAME"), os.Getenv("SYSTEM_EMAIL"))
	subject := "Reset your credentials"
	to := mail.NewEmail("Sender User", sm.Email)
	plainTextContent := "Hi, Please set your credentials by following steps below:"

	var err error
	t, err := template.ParseFiles("./html/email_template.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, nil); err != nil {
		return err
	}
	htmlContent := tpl.String()

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)
	return nil
}
