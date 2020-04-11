package community_infrastructure

import (
	"github.com/matcornic/hermes/v2"
	vo "github.com/214alphadev/community-bl/value_objects"
	email "github.com/214alphadev/email-delivery-go"
)

type confirmationEmail struct {
	email.Email
	CompanyInfo
}

type CompanyInfo struct {
	Name            string
	Greeting        string
	URL             string
	LogoURL         string
	Introduction    string
	Instructions    string
	ButtonColourHex string
	ButtonText      string
	ClickURL        string
	Outro           string
	Copyright       string
	TroubleText     string
}

func newS33DConfirmationEmail(fromEmail, recipientEmail vo.EmailAddress, confirmationCode vo.ConfirmationCode) confirmationEmail {
	confirmationEmail := confirmationEmail{
		Email: email.Email{
			FromName:       "S33D",
			FromEmail:      fromEmail.String(),
			RecipientEmail: recipientEmail.String(),
			Subject:        "Confirmation Code",
		},
		CompanyInfo: CompanyInfo{
			Name:            "S33D",
			Greeting:        "For increased security, 2FA is required",
			URL:             "https://s33d.life/",
			LogoURL:         "https://s33d.life/wp-content/uploads/2019/03/S33D.life_Logo.png",
			Introduction:    "",
			Instructions:    "Enter your confirmation code in the app to authenticate.",
			Outro:           "If you encounter any issues, contact us.",
			Copyright:       "Copyright © 2019 S33D.  All rights reserved.",
			ButtonColourHex: "#FABB4D",
			TroubleText:     " ",
		},
	}
	confirmationEmail.CompanyInfo.ButtonText = confirmationCode.String()
	confirmationEmail.Email.HermesTheme = confirmationEmail.HermesTheme()
	confirmationEmail.Email.HermesEmail = confirmationEmail.HermesEmail()
	return confirmationEmail
}

func (c confirmationEmail) HermesTheme() hermes.Hermes {
	return hermes.Hermes{
		Product: hermes.Product{
			Name: c.CompanyInfo.Name,
			Link: c.CompanyInfo.URL,
			Logo:        c.CompanyInfo.LogoURL,
			Copyright:   c.CompanyInfo.Copyright,
			TroubleText: c.CompanyInfo.TroubleText,
		},
	}
}

func (c confirmationEmail) HermesEmail() hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Greeting: c.Greeting,
			Intros: []string{
				c.CompanyInfo.Introduction,
			},
			Actions: []hermes.Action{
				{
					Instructions: c.CompanyInfo.Instructions,
					Button: hermes.Button{
						Color: c.CompanyInfo.ButtonColourHex,
						Text:  c.CompanyInfo.ButtonText,
						Link:  c.CompanyInfo.ClickURL,
					},
				},
			},
			Outros: []string{
				c.CompanyInfo.Outro,
			},
		},
	}
}
