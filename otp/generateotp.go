package otp

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/smtp"
)

const otp = "1234567890"

func Generateotp(email string) string {
	b := make([]byte, 6)
	_, err := io.ReadAtLeast(rand.Reader, b, 6)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = otp[int(b[i])%len(otp)]
	}
	//	fmt.Print(string(b))
	from := "scgamer1401@gmail.com"
	password := "Cs@8722552149"
	to := []string{
		email,
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	str := fmt.Sprintf("The Generated OTP is %s", string(b))
	message := []byte(str)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
