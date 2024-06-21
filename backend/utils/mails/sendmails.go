package mails

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

func SendForgetPasswordMail(receiver []string, name, token string) error {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	// Create a new message
	message := gomail.NewMessage()

	// Set sender and recipient
	message.SetHeader("From", from)
	message.SetHeader("To", receiver...)
	message.SetHeader("Subject", "Password Change Request")

	// Load the email template
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(dir + "/templates/resetpassword.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	t.Execute(&body, struct {
		Name              string
		PasswordResetLink string
	}{
		Name:              name,
		PasswordResetLink: "http://localhost:3000/reset-password?token=" + token,
	})

	// Attach the HTML body to the email
	message.SetBody("text/html", body.String())

	// Create a new dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 465, from, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Email Sent Successfully!")

	return nil
}
