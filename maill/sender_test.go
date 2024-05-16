package maill

import (
	"fmt"

	//"github.com/joho/godotenv"

	//"github.com/stretchr/testify/require"
)

// func init(){
// 	initializers.LoadVariable()
// }

func TestSendEmailWithGmail(recipientEmail string) {

	/*
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	*/

	senderName :="Hello"
	senderAddress := "ttechcodebeelab@gmail.com"
	senderPassword := "bgmcyapqpbcafeqi"

	fmt.Println("sender name:", senderName)
	fmt.Println("sender Address:", senderAddress)
	fmt.Println("sender Password:", senderPassword)

	sender := NewGmailSender(senderName, senderAddress, senderPassword)
	subject := "A test email"
	content := `
	<h1>Hello World</h1>
	<p>This is a test message</p>
	`
	to := []string{recipientEmail}
	//attachFiles := []string{"./sender.go"}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)

	if err!=nil{
		fmt.Println("error",err)
	}
}
