package otp

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"net/smtp"
	"os"
)

const otp = "1234567890"

func Generateotp() (string, bool) {
	//generating random numbers of len 6
	b := make([]byte, 6)
	_, err := io.ReadAtLeast(rand.Reader, b, 6)
	if err != nil {
		return err.Error(), true
	}
	for i := 0; i < len(b); i++ {
		b[i] = otp[int(b[i])%len(otp)]
	}
	return string(b), false
}
func Emailgenerate(email, b string) error {
	file, _ := os.Open("credentials.csv")
	defer file.Close()
	csvfile := csv.NewReader(file)
	det, _ := csvfile.Read()
	from := det[0]
	password := det[1]
	to := []string{
		email,
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	str := fmt.Sprintf("The Generated OTP is %s", b)
	message := []byte(str)
	auth := smtp.PlainAuth("", from, password, smtpHost) //

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}
