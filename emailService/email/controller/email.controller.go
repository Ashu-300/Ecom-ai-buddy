package controller

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(receiverMail string, receiverName string) {
	senderMail := os.Getenv("SENDER_MAIL")
	sendgridApiKey := os.Getenv("SENDGRID_API_KEY")

	if senderMail == "" || sendgridApiKey == "" {
		log.Print("❌ SENDGRID_API_KEY or SENDER_MAIL is empty")
		return
	}

	from := mail.NewEmail("SUPERNOVA Team", senderMail)
	subject := fmt.Sprintf("Welcome to SUPERNOVA, %s!", receiverName)
	to := mail.NewEmail(receiverName, receiverMail)

	// Plain text version for clients that don't support HTML
	plainTextContent := fmt.Sprintf(
		"Hello %s,\n\nWelcome to SUPERNOVA! We're excited to have you on board.\n"+
			"You can now explore our platform and enjoy our services.\n\n"+
			"Best regards,\nThe SUPERNOVA Team",
		receiverName,
	)

	// HTML version for richer formatting
	htmlContent := fmt.Sprintf(
		`<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6;">
				<h2>Hello %s,</h2>
				<p>Welcome to <strong>SUPERNOVA</strong>! We're thrilled to have you on board.</p>
				<p>You can now explore our platform and enjoy our services.</p>
				<p>Feel free to reach out to us anytime at <a href="mailto:support@supernova.com">support@supernova.com</a>.</p>
				<br>
				<p>Best regards,<br><strong>The SUPERNOVA Team</strong></p>
			</body>
		</html>`,
		receiverName,
	)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(sendgridApiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println("❌ Error sending email:", err)
	} else {
		log.Printf("📧 Email sent to %s | Status: %d\n", receiverMail, response.StatusCode)
	}
}
