package maill

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)

		if err != nil {
			return fmt.Errorf("Failed to attch file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func SendEmailWithGmail(recipientEmail string, receipientCode string, recipientName string, recipientId int) error {

	/*
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file:", err)
			return
		}
	*/

	senderName := "Hello"
	senderAddress := "ttechcodebeelab@gmail.com"
	senderPassword := "bgmcyapqpbcafeqi"

	// fmt.Println("sender name:", senderName)
	// fmt.Println("sender Address:", senderAddress)
	// fmt.Println("sender Password:", senderPassword)

	sender := NewGmailSender(senderName, senderAddress, senderPassword)
	subject := "Welcome to the Go lang Crud App"
	verifyUrl := fmt.Sprintf("http://localhost/verify_email?id=%d&secret_code=%s",
		recipientId, receipientCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering with us! <br/>
	Please <a href="%s">Click Here <a/> to verify your email address. <br/>
	`, recipientName, verifyUrl)
	to := []string{recipientEmail}
	//attachFiles := []string{"./sender.go"}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		fmt.Println("Failed to send the mail")
	}
	//fmt.Println(recipientId, recipientEmail, receipientCode, recipientName)
	return nil
}
