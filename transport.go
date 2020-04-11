package community_infrastructure

import (
	"fmt"
	cd "github.com/214alphadev/community-bl"
	vo "github.com/214alphadev/community-bl/value_objects"
	email "github.com/214alphadev/email-delivery-go"
)

type Transport struct {
	emailDelivery email.ESPSenderInterface
	sendFrom      vo.EmailAddress
}

func (t *Transport) SendConfirmationCode(confirmationCode cd.ConfirmationCode) error {

	confirmationCodeEmail := newS33DConfirmationEmail(t.sendFrom, confirmationCode.EmailAddress, confirmationCode.ConfirmationCode)

	confirmationCodeEmail.CompanyInfo.ButtonText = confirmationCode.ConfirmationCode.String()
	confirmationCodeEmail.Email.HermesTheme = confirmationCodeEmail.HermesTheme()
	confirmationCodeEmail.Email.HermesEmail = confirmationCodeEmail.HermesEmail()

	go func() {
		_, err := t.emailDelivery.Send(confirmationCodeEmail.Email)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func NewTransport(emailTransport email.ESPSenderInterface, email vo.EmailAddress) *Transport {
	return &Transport{
		emailDelivery: emailTransport,
		sendFrom:      email,
	}
}
