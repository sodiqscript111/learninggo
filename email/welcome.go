package email

import (
	"fmt"
	"github.com/resend/resend-go/v2"
)

func SendWelcomeEmail(to, name string) error {
	apiKey := "re_aENqruXp_9eoxbdzctoZm6aiaLXTC4QhN"
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable is not set")
	}

	client := resend.NewClient(apiKey)

	resp, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    "EventManager <onboarding@resend.dev>",
		To:      []string{to},
		Subject: "Welcome to EventManager",
		Html:    fmt.Sprintf("<strong>Welcome, %s!</strong><p>Thanks for signing up 🎉</p>", name),
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w, response: %v", err, resp)
	}

	return nil
}
