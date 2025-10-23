package controller

import (
	"fmt"
	"log"
	"os"
	"strings"
	"supernova/emailService/email/dto"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func AuthEmail(receiverMail string, receiverName string) {
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


func PaymentInitiatedEmail(body dto.PaymentData) {
	senderMail := os.Getenv("SENDER_MAIL")
	sendgridApiKey := os.Getenv("SENDGRID_API_KEY")

	if senderMail == "" || sendgridApiKey == "" {
		log.Print("❌ SENDGRID_API_KEY or SENDER_MAIL is empty")
		return
	}

	receiverName := strings.Split(body.ReceiverMail, "@")[0]

	from := mail.NewEmail("SUPERNOVA Payments", senderMail)
	subject := fmt.Sprintf("Payment Initiated for Order #%s", body.OrderID)
	to := mail.NewEmail(receiverName, body.ReceiverMail)

	// Plain text version
	plainTextContent := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Your payment process has been initiated.\n\n"+
			"Order ID: %s\nPayment ID: %s\nAmount: %.2f %s\n\n"+
			"You will receive a confirmation once the payment is successfully processed.\n\n"+
			"Thank you for choosing SUPERNOVA!\n\n"+
			"Best regards,\nSUPERNOVA Payments Team",
		receiverName, body.OrderID, body.PaymentID, body.Amount, body.Currency,
	)

	// HTML version
	htmlContent := fmt.Sprintf(
		`<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<h2>Payment Initiated 💳</h2>
				<p>Hi <strong>%s</strong>,</p>
				<p>Your payment process has been initiated. Please wait for confirmation once the transaction is completed.</p>

				<table style="border-collapse: collapse; margin-top: 10px;">
					<tr><td><strong>Order ID:</strong></td><td>%s</td></tr>
					<tr><td><strong>Payment ID:</strong></td><td>%s</td></tr>
					<tr><td><strong>Amount:</strong></td><td>%.2f %s</td></tr>
				</table>

				<p style="margin-top: 15px;">
					Thank you for choosing <strong>SUPERNOVA</strong>! We appreciate your trust in our services.
				</p>

				<p>If you have any questions, feel free to reach out to our support team at 
					<a href="mailto:support@supernova.com">support@supernova.com</a>.
				</p>

				<br>
				<p>Warm regards,<br><strong>The SUPERNOVA Payments Team</strong></p>
			</body>
		</html>`,
		receiverName, body.OrderID, body.PaymentID, body.Amount, body.Currency,
	)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridApiKey)

	response, err := client.Send(message)
	if err != nil {
		log.Println("❌ Error sending payment initiated email:", err)
	} else {
		log.Printf("💳 Payment initiated email sent to %s | Status: %d\n", body.ReceiverMail, response.StatusCode)
	}
}


func ProductCreatedEmail(body dto.ProductData) {
	senderMail := os.Getenv("SENDER_MAIL")
	sendgridApiKey := os.Getenv("SENDGRID_API_KEY")

	if senderMail == "" || sendgridApiKey == "" {
		log.Print("❌ SENDGRID_API_KEY or SENDER_MAIL is empty")
		return
	}

	receiverName := strings.Split(body.ReceiverMail, "@")[0]

	from := mail.NewEmail("SUPERNOVA Marketplace", senderMail)
	subject := fmt.Sprintf("New Product Created: %s", body.ProductName)
	to := mail.NewEmail(receiverName, body.ReceiverMail)

	// Plain text content
	plainTextContent := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Your new product has been successfully created in the marketplace.\n\n"+
			"Product Name: %s\nProduct ID: %s\nPrice: %.2f %s\nCategory: %s\n\n"+
			"Your product is now live and visible to customers.\n\n"+
			"Thank you for using SUPERNOVA Marketplace!\n\n"+
			"Best regards,\nSUPERNOVA Marketplace Team",
		receiverName, body.ProductName, body.ProductID, body.Price, body.Currency,
	)

	// HTML content
	htmlContent := fmt.Sprintf(
		`<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<h2>New Product Created 🛒</h2>
				<p>Hi <strong>%s</strong>,</p>
				<p>Your new product has been successfully added to the marketplace.</p>

				<table style="border-collapse: collapse; margin-top: 10px;">
					<tr><td><strong>Product Name:</strong></td><td>%s</td></tr>
					<tr><td><strong>Product ID:</strong></td><td>%s</td></tr>
					<tr><td><strong>Price:</strong></td><td>%.2f %s</td></tr>
					<tr><td><strong>Category:</strong></td><td>%s</td></tr>
				</table>

				<p style="margin-top: 15px;">
					Your product is now live and visible to customers!
				</p>

				<p>If you have any questions, feel free to contact our support team at 
					<a href="mailto:support@supernova.com">support@supernova.com</a>.
				</p>

				<br>
				<p>Warm regards,<br><strong>The SUPERNOVA Marketplace Team</strong></p>
			</body>
		</html>`,
		receiverName, body.ProductName, body.ProductID, body.Price, body.Currency,
	)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridApiKey)

	response, err := client.Send(message)
	if err != nil {
		log.Println("❌ Error sending product created email:", err)
	} else {
		log.Printf("🛒 Product created email sent to %s | Status: %d\n", body.ReceiverMail, response.StatusCode)
	}
}

func OrderPlacedEmail(body dto.OrderData) {
	senderMail := os.Getenv("SENDER_MAIL")
	sendgridApiKey := os.Getenv("SENDGRID_API_KEY")

	if senderMail == "" || sendgridApiKey == "" {
		log.Print("❌ SENDGRID_API_KEY or SENDER_MAIL is empty")
		return
	}

	receiverName := strings.Split(body.ReceiverMail, "@")[0]

	from := mail.NewEmail("SUPERNOVA Marketplace", senderMail)
	subject := fmt.Sprintf("Order Placed Successfully: #%s", body.OrderID)
	to := mail.NewEmail(receiverName, body.ReceiverMail)

	// Plain text content
	plainTextContent := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Your order has been successfully placed!\n\n"+
			"Order ID: %s\nProduct Name: %s\nQuantity: %d\nTotal Amount: %.2f %s\n\n"+
			"You will receive a confirmation once the order is processed.\n\n"+
			"Thank you for shopping with SUPERNOVA Marketplace!\n\n"+
			"Best regards,\nSUPERNOVA Marketplace Team",
		receiverName, body.OrderID, body.TotalAmount, body.Currency,
	)

	// HTML content
	htmlContent := fmt.Sprintf(
		`<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<h2>Order Placed ✅</h2>
				<p>Hi <strong>%s</strong>,</p>
				<p>Thank you for your order! Your order has been successfully placed.</p>

				<table style="border-collapse: collapse; margin-top: 10px;">
					<tr><td><strong>Order ID:</strong></td><td>%s</td></tr>
					<tr><td><strong>Product Name:</strong></td><td>%s</td></tr>
					<tr><td><strong>Quantity:</strong></td><td>%d</td></tr>
					<tr><td><strong>Total Amount:</strong></td><td>%.2f %s</td></tr>
				</table>

				<p style="margin-top: 15px;">
					You will receive a confirmation once the order is processed.
				</p>

				<p>If you have any questions, feel free to contact our support team at 
					<a href="mailto:support@supernova.com">support@supernova.com</a>.
				</p>

				<br>
				<p>Warm regards,<br><strong>The SUPERNOVA Marketplace Team</strong></p>
			</body>
		</html>`,
		receiverName, body.OrderID, body.TotalAmount, body.Currency,
	)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridApiKey)

	response, err := client.Send(message)
	if err != nil {
		log.Println("❌ Error sending order placed email:", err)
	} else {
		log.Printf("🛒 Order placed email sent to %s | Status: %d\n", body.ReceiverMail, response.StatusCode)
	}
}