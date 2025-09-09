// Package email provides email sending functionality for the application
package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
)

// Service handles email operations
type Service struct {
	config Config
	logger *slog.Logger
}

// Config contains email service configuration
type Config struct {
	Enabled  bool
	Service  string // smtp service provider
	Host     string
	Port     int
	User     string
	Password string
	From     string // sender email address
	FromName string // sender name
}

// NewService creates a new email service
func NewService(config Config, logger *slog.Logger) *Service {
	// Set defaults based on service provider
	if config.Service == "gmail" && config.Host == "" {
		config.Host = "smtp.gmail.com"
		config.Port = 587
	} else if config.Service == "sendgrid" && config.Host == "" {
		config.Host = "smtp.sendgrid.net"
		config.Port = 587
	}

	// Default from address if not set
	if config.From == "" {
		config.From = config.User
	}

	if config.FromName == "" {
		config.FromName = "ArchesAI"
	}

	return &Service{
		config: config,
		logger: logger,
	}
}

// Email represents an email message
type Email struct {
	To      string
	Subject string
	Body    string
	HTML    bool
}

// SendEmail sends an email message
func (s *Service) SendEmail(_ context.Context, email Email) error {
	if !s.config.Enabled {
		s.logger.Info("email service disabled, skipping email", "to", email.To, "subject", email.Subject)
		return nil
	}

	// Build message
	msg := s.buildMessage(email)

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)

	// Connect to server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Send email with TLS
	err := s.sendEmailWithTLS(addr, auth, s.config.From, []string{email.To}, msg)
	if err != nil {
		s.logger.Error("failed to send email", "error", err, "to", email.To)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("email sent successfully", "to", email.To, "subject", email.Subject)
	return nil
}

// buildMessage builds the email message
func (s *Service) buildMessage(email Email) []byte {
	var msg strings.Builder

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	// Content type
	if email.HTML {
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	}

	// Body
	msg.WriteString("\r\n")
	msg.WriteString(email.Body)

	return []byte(msg.String())
}

// sendEmailWithTLS sends email using TLS
func (s *Service) sendEmailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Connect to the SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			s.logger.Warn("failed to close SMTP client", "error", err)
		}
	}()

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set sender and recipients
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection
	return client.Quit()
}

// SendVerificationEmail sends a verification email to a user
func (s *Service) SendVerificationEmail(ctx context.Context, toEmail, userName, verificationToken string) error {
	// Build verification URL - this should be configured
	baseURL := "http://localhost:3000" // TODO: Get from config
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", baseURL, verificationToken)

	subject := "Verify your ArchesAI account"

	// HTML email body
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify your email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f4f4f4; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; }
        .footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #ddd; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to ArchesAI!</h1>
        </div>
        <div class="content">
            <p>Hi %s,</p>
            <p>Thank you for signing up with ArchesAI. Please verify your email address by clicking the button below:</p>
            <p style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">Verify Email Address</a>
            </p>
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #007bff;">%s</p>
            <p>This link will expire in 24 hours.</p>
            <p>If you didn't create an account with ArchesAI, you can safely ignore this email.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The ArchesAI Team</p>
        </div>
    </div>
</body>
</html>
`, userName, verificationURL, verificationURL)

	return s.SendEmail(ctx, Email{
		To:      toEmail,
		Subject: subject,
		Body:    htmlBody,
		HTML:    true,
	})
}

// SendPasswordResetEmail sends a password reset email
func (s *Service) SendPasswordResetEmail(ctx context.Context, toEmail, userName, resetToken string) error {
	// Build reset URL
	baseURL := "http://localhost:3000" // TODO: Get from config
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", baseURL, resetToken)

	subject := "Reset your ArchesAI password"

	// HTML email body
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset your password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f4f4f4; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #dc3545; color: white; text-decoration: none; border-radius: 4px; }
        .footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #ddd; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>Hi %s,</p>
            <p>We received a request to reset your password. Click the button below to create a new password:</p>
            <p style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #dc3545;">%s</p>
            <p>This link will expire in 1 hour for security reasons.</p>
            <p>If you didn't request a password reset, please ignore this email. Your password won't be changed.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The ArchesAI Team</p>
        </div>
    </div>
</body>
</html>
`, userName, resetURL, resetURL)

	return s.SendEmail(ctx, Email{
		To:      toEmail,
		Subject: subject,
		Body:    htmlBody,
		HTML:    true,
	})
}

// SendWelcomeEmail sends a welcome email after successful verification
func (s *Service) SendWelcomeEmail(ctx context.Context, toEmail, userName string) error {
	subject := "Welcome to ArchesAI!"

	// HTML email body
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to ArchesAI</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #ddd; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Email Verified Successfully!</h1>
        </div>
        <div class="content">
            <p>Hi %s,</p>
            <p>Your email has been successfully verified. You now have full access to all ArchesAI features!</p>
            <p>Here are some things you can do next:</p>
            <ul>
                <li>Complete your profile</li>
                <li>Create your first organization</li>
                <li>Explore our workflow automation features</li>
                <li>Connect with AI providers</li>
            </ul>
            <p>If you have any questions, don't hesitate to reach out to our support team.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The ArchesAI Team</p>
        </div>
    </div>
</body>
</html>
`, userName)

	return s.SendEmail(ctx, Email{
		To:      toEmail,
		Subject: subject,
		Body:    htmlBody,
		HTML:    true,
	})
}
