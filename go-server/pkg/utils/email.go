package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func GetEmailConfig() EmailConfig {
	return EmailConfig{
		Host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		Port:     587, // TLS port
		Username: getEnvOrDefault("SMTP_USER", os.Getenv("EMAIL_USER")),
		Password: strings.ReplaceAll(getEnvOrDefault("SMTP_PASSWORD", os.Getenv("EMAIL_PASS")), " ", ""),
		From:     getEnvOrDefault("FROM_EMAIL", getEnvOrDefault("SMTP_USER", os.Getenv("EMAIL_USER"))),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func SendInvitationEmail(to, inviteLink string) error {
	config := GetEmailConfig()

	if config.Username == "" || config.Password == "" {
		return fmt.Errorf("SMTP credentials not configured. Set SMTP_USER and SMTP_PASSWORD in .env")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Action Required: RBAC System Invitation")

	// Extract token from inviteLink
	var token string
	if len(inviteLink) > 0 {
		for i := len(inviteLink) - 1; i >= 0; i-- {
			if inviteLink[i] == '=' {
				token = inviteLink[i+1:]
				break
			}
		}
	}

	declineLink := fmt.Sprintf("http://localhost:5173/decline-invitation?token=%s", token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; background-color: #f4f4f5; margin: 0; padding: 0; }
        .container { max-width: 500px; margin: 40px auto; background: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06); }
        .content { padding: 40px; text-align: center; }
        .h1 { font-size: 24px; font-weight: 700; color: #18181b; margin-bottom: 8px; }
        .p { color: #52525b; margin-bottom: 32px; font-size: 16px; }
        .btn-container { display: flex; justify-content: center; gap: 16px; margin-bottom: 32px; }
        .btn { display: inline-block; padding: 12px 24px; font-size: 14px; font-weight: 500; text-decoration: none; border-radius: 8px; transition: opacity 0.2s; }
        .btn-accept { background-color: #18181b; color: #ffffff; }
        .btn-decline { background-color: #ef4444; color: #ffffff; }
        .footer { background-color: #fafafa; padding: 20px; text-align: center; border-top: 1px solid #e4e4e7; }
        .footer-text { font-size: 12px; color: #71717a; }
    </style>
</head>
<body>
    <div class="container">
        <div class="content">
            <div class="h1">Welcome to RBAC System</div>
            <div class="p">You have been invited to join the team.</div>
            
            <div class="btn-container">
                <a href="%s" class="btn btn-accept">Accept Invitation</a>
                <a href="%s" class="btn btn-decline">Decline Invitation</a>
            </div>
        </div>
        <div class="footer">
            <div class="footer-text">Links expire in 24 hours.</div>
            <div class="footer-text" style="margin-top: 8px;">If you didn't expect this invitation, please ignore this email.</div>
        </div>
    </div>
</body>
</html>
	`, inviteLink, declineLink)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
