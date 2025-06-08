package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"learninggo/db"
	"learninggo/email"
)

func ListenForUserEvents() {
	sub := db.RedisClient.Subscribe(context.Background(), "userEvent")
	ch := sub.Channel()

	fmt.Println("✅ Listening for userEvent channel...")

	for msg := range ch {
		var event map[string]string
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			fmt.Println("❌ Failed to unmarshal event:", err)
			continue
		}

		switch event["type"] {
		case "user":
			err := email.SendWelcomeEmail(event["email"], event["name"])
			if err != nil {
				fmt.Printf("❌ Failed to send welcome email: %v\n", err)
			} else {
				fmt.Printf("📬 New user registered: Name = %s, Email = %s\n", event["name"], event["email"])
			}
		default:
			fmt.Println("⚠️ Unknown event type received:", event["type"])
		}
	}
}
