package notifications

import (
	"fmt"
	"log"

	cfg "github.com/akhilckenshi/notification/pkg/settings"

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
func SendWhatsAppMessage(to, sub, message string) error {
	accountSid := cfg.Config.WhatsAppProviderKey
	authToken := cfg.Config.WhatsAppProviderSecret
	fromNumber := "whatsapp:" + cfg.Config.WhatsAppFromNumber
	toNumber := "whatsapp:" + to

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(fromNumber)
	params.SetBody(fmt.Sprintf("%s: %s", sub, message))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
	}

	fmt.Printf("Message sent! SID: %s", *resp.Sid)

	return nil
}

// // SendWhatsAppMessage sends a WhatsApp notification with a PDF attachment.
// func SendWhatsAppMessage(to, sub, htmlContent string) error {
// 	// Get Twilio credentials from the config
// 	accountSid := cfg.Config.WhatsAppProviderKey
// 	authToken := cfg.Config.WhatsAppProviderSecret
// 	fromNumber := "whatsapp:" + cfg.Config.WhatsAppFromNumber
// 	toNumber := "whatsapp:" + to

// 	// Initialize the Twilio client
// 	client := twilio.NewRestClientWithParams(twilio.ClientParams{
// 		Username: accountSid,
// 		Password: authToken,
// 	})
// 	fmt.Println("747474773")
// 	// Convert HTML to PDF and store it locally
// 	pdfPath, err := ConvertHTMLToPDF(htmlContent)
// 	if err != nil {
// 		return fmt.Errorf("error converting HTML to PDF: %v", err)
// 	}
// 	defer os.Remove(pdfPath) // Ensure file is removed after sending

// 	// Send the PDF as a WhatsApp message
// 	params := &openapi.CreateMessageParams{}
// 	params.SetTo(toNumber)
// 	params.SetFrom(fromNumber)
// 	params.SetBody(fmt.Sprintf("%s: %s", sub, "Please find the PDF attached."))
// 	params.SetMediaUrl([]string{fmt.Sprintf("file://%s", pdfPath)}) // Passing the URL as a slice of strings

// 	// Send the message through Twilio
// 	resp, err := client.Api.CreateMessage(params)
// 	if err != nil {
// 		log.Printf("Error sending WhatsApp message: %v", err)
// 		return err
// 	}

// 	// Log the response SID for tracking
// 	fmt.Printf("Message sent! SID: %s\n", *resp.Sid)
// 	return nil
// }

// // ConvertHTMLToPDF converts HTML string content to a PDF and returns the file path.
// func ConvertHTMLToPDF(htmlContent string) (string, error) {
// 	fmt.Println("0000000")
// 	// Set the wkhtmltopdf binary path if needed (optional)
// 	// os.Setenv("WKHTMLTOPDF_BIN", "/path/to/wkhtmltopdf")

// 	// Create a new PDF generator instance
// 	pdf, _ := wkhtmltopdf.NewPDFGenerator()
// 	fmt.Println("ddddd", pdf)
// 	// if err != nil {
// 	// 	fmt.Println("first")
// 	// 	return "", fmt.Errorf("failed to create PDF generator: %v", err)
// 	// }

// 	// Convert HTML content to io.Reader using bytes.NewReader
// 	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent)))
// 	fmt.Println("eee")
// 	// Add the page to the PDF generator
// 	pdf.AddPage(page)

// 	// Create a temporary file for the PDF
// 	fmt.Println("111")
// 	pdfDir := "pdf"
// 	if err := os.MkdirAll(pdfDir, os.ModePerm); err != nil {
// 		return "", fmt.Errorf("failed to create PDF directory: %v", err)
// 	}
// 	tempFile := fmt.Sprintf("%s/temp_pdf_%d.pdf", pdfDir, os.Getpid())
// 	fmt.Println("22222")
// 	// Create the PDF
// 	err := pdf.Create()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create PDF: %v", err)
// 	}
// 	fmt.Println("3333")
// 	// Write the generated PDF to the file
// 	err = pdf.WriteFile(tempFile)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to write PDF to file: %v", err)
// 	}
// 	fmt.Println("444")
// 	// Return the path to the temporary PDF file
// 	return tempFile, nil
// }
