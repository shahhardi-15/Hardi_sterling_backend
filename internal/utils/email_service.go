package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailService struct {
	SMTPHost       string
	SMTPPort       string
	SenderEmail    string
	SenderPassword string
	FromName       string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SenderEmail:    getEnv("SENDER_EMAIL", "noreply@sterlinghms.com"),
		SenderPassword: getEnv("SENDER_PASSWORD", ""),
		FromName:       getEnv("FROM_NAME", "Sterling HMS"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SendPasswordResetEmail sends password reset link to user email
func (es *EmailService) SendPasswordResetEmail(toEmail string, userFirstName string, resetLink string) error {
	subject := "Password Reset Request - Sterling HMS"
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #1e40af; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
			.content { background-color: #f3f4f6; padding: 20px; border-radius: 0 0 5px 5px; }
			.button { text-align: center; margin: 30px 0; }
			.reset-btn { background-color: #1e40af; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold; }
			.reset-btn:hover { background-color: #1e3a8a; }
			.warning { color: #dc2626; font-size: 12px; margin-top: 20px; }
			.footer { margin-top: 20px; font-size: 12px; color: #6b7280; text-align: center; }
			.link-text { word-break: break-all; color: #1e40af; font-size: 12px; margin-top: 10px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Sterling HMS - Password Reset</h1>
			</div>
			<div class="content">
				<p>Hello %s,</p>
				<p>You requested a password reset for your Sterling HMS account. Click the button below to reset your password:</p>
				<div class="button">
					<a href="%s" class="reset-btn">Reset Password</a>
				</div>
				<p>Or copy and paste this link in your browser:</p>
				<div class="link-text">%s</div>
				<p>This link will expire in <strong>1 hour</strong>.</p>
				<div class="warning">
					⚠️ <strong>Security Note:</strong>
					<ul>
						<li>Never share this link with anyone</li>
						<li>Sterling HMS staff will never ask for this link</li>
						<li>If you didn't request this, please ignore this email</li>
					</ul>
				</div>
				<p style="margin-top: 30px; color: #6b7280;">Best regards,<br>Sterling HMS Team</p>
				<div class="footer">
					<p>This is an automated email. Please do not reply.</p>
				</div>
			</div>
		</div>
	</body>
	</html>
	`, userFirstName, resetLink, resetLink)

	return es.sendEmail(toEmail, subject, htmlBody)
}

// SendPasswordResetSuccessEmail sends confirmation email after successful password reset
func (es *EmailService) SendPasswordResetSuccessEmail(toEmail string, userFirstName string) error {
	subject := "Password Reset Successful - Sterling HMS"
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #059669; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
			.content { background-color: #f3f4f6; padding: 20px; border-radius: 0 0 5px 5px; }
			.success-icon { text-align: center; font-size: 48px; margin: 20px 0; }
			.footer { margin-top: 20px; font-size: 12px; color: #6b7280; text-align: center; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Sterling HMS - Password Reset Successful</h1>
			</div>
			<div class="content">
				<div class="success-icon">✓</div>
				<p>Hello %s,</p>
				<p>Your password has been successfully reset. You can now log in to your Sterling HMS account with your new password.</p>
				<p><strong>Next Steps:</strong></p>
				<ul>
					<li>Log in with your new password</li>
					<li>If you notice any suspicious activity, contact our support team immediately</li>
				</ul>
				<p style="margin-top: 30px; color: #6b7280;">Best regards,<br>Sterling HMS Team</p>
				<div class="footer">
					<p>This is an automated email. Please do not reply.</p>
				</div>
			</div>
		</div>
	</body>
	</html>
	`, userFirstName)

	return es.sendEmail(toEmail, subject, htmlBody)
}

// SendEmail sends raw email
func (es *EmailService) sendEmail(toEmail string, subject string, htmlBody string) error {
	// In development, log instead of sending
	if os.Getenv("ENV") == "development" {
		log.Printf("[EMAIL] To: %s\nSubject: %s\nBody: %s\n", toEmail, subject, htmlBody)
		return nil
	}

	// Production: send via SMTP
	from := fmt.Sprintf("%s <%s>", es.FromName, es.SenderEmail)
	to := []string{toEmail}

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		from,
		toEmail,
		subject,
		htmlBody,
	)

	auth := smtp.PlainAuth("", es.SenderEmail, es.SenderPassword, es.SMTPHost)
	addr := fmt.Sprintf("%s:%s", es.SMTPHost, es.SMTPPort)

	err := smtp.SendMail(addr, auth, es.SenderEmail, to, []byte(message))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}
