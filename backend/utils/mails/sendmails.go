package mails

import (
	"backend/serializers"
	"backend/state"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(header string, message *gomail.Message, receiver []string) error {
	from := state.AppConfig.EmailUser
	password := state.AppConfig.EmailPassword
	message.SetHeader("From", from)
	message.SetHeader("To", receiver...)
	message.SetHeader("Subject", header)
	// Create a new dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 465, from, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}

func SendForgetPasswordMail(receiver []string, name, token string) error {

	message := gomail.NewMessage()
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
		PasswordResetLink: fmt.Sprintf("%s?token=%s", state.AppConfig.PasswordResetLink, token),
	})

	// Attach the HTML body to the email
	message.SetBody("text/html", body.String())
	if err := SendMail("Password Reset Request", message, receiver); err != nil {
		return err
	}

	return nil
}

func AdminOnRampMail(receiver []string, data serializers.AdminOnRampSerializer) error {
	message := gomail.NewMessage()
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(dir + "/templates/confirm-onramp.html")
	if err != nil {
		return err
	}
	var body bytes.Buffer
	t.Execute(&body, data)

	message.SetBody("text/html", body.String())
	if err := SendMail("OnRamp Request", message, receiver); err != nil {
		return err
	}

	return nil
}

func AdminOffRampMail(receiver []string, data serializers.AdminOffRampSerializer) error {
	message := gomail.NewMessage()
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(dir + "/templates/verify-offramp.html")
	if err != nil {
		return err
	}
	var body bytes.Buffer
	t.Execute(&body, data)

	message.SetBody("text/html", body.String())
	if err := SendMail("OffRamp Confirmation", message, receiver); err != nil {
		return err
	}

	return nil
}

func UserOffRampMail(receiver []string, data serializers.UserOffRampMail) error {
	message := gomail.NewMessage()
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(dir + "/templates/user-offramp.html")
	if err != nil {
		return err
	}
	var body bytes.Buffer
	t.Execute(&body, data)

	message.SetBody("text/html", body.String())
	if err := SendMail("OffRamp Notification", message, receiver); err != nil {
		return err
	}

	return nil
}
