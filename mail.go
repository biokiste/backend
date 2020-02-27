package main

import (
	"net/smtp"

	"github.com/spf13/viper"
)

// SendEMail send an email to recipient
func (h Handlers) SendEMail(mailRecipient string) error {

	smtpHost := viper.GetString("smtpHost")
	smtpPort := viper.GetString("smtpPort")
	smtpPassword := viper.GetString("smtpPassword")
	smtpUsername := viper.GetString("smtpUsername")

	msg := "From: " + smtpUsername + "\n" +
		"To: " + mailRecipient + "\n" +
		"Subject: Hello there\n\n" +
		"Biokiste"

	err := smtp.SendMail(smtpHost+":"+smtpPort,
		smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost),
		smtpUsername, []string{mailRecipient}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
