package notifications

import (
	"fmt"
	"log"
	cfg "smartsme-notificationservice/pkg/settings"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioConfig struct {
	AccountSID string `mapstructure:"account_sid"` // Twilio Account SID
	AuthToken  string `mapstructure:"auth_token"`  // Twilio Auth Token
	FromNumber string `mapstructure:"from_number"` // Phone number to send SMS from
	Url        string `mapstructure:"url"`
}

// SendWhatsAppMessage sends a WhatsApp notification using a third-party API like Twilio.
func SendWhatsAppMessage(to string, message string) error {
	accountSid := cfg.Config.WhatsappConfig.Key
	authToken := cfg.Config.WhatsappConfig.Secret
	fromNumber := "whatsapp:" + cfg.Config.WhatsappConfig.Number
	toNumber := "whatsapp:" + to

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(fromNumber)
	params.SetBody(fmt.Sprintf("Your verification code is: %s", message))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
	}

	fmt.Printf("Message sent! SID: %s", *resp.Sid)

	return nil
}
