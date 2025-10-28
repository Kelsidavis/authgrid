package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ResendEmail represents an email to send via Resend API
type ResendEmail struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// sendWelcomeEmail sends a welcome email to new customers
func sendWelcomeEmail(customerEmail, customerName, subscriptionID string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: RESEND_API_KEY not set - email not sent")
		return fmt.Errorf("email service not configured")
	}

	// Create welcome email
	// Using Resend's onboarding domain temporarily (works immediately)
	// After domain verification in Resend, switch to: "Authgrid <hello@authgrid.net>"
	email := ResendEmail{
		From:    "Authgrid <onboarding@resend.dev>",
		To:      []string{customerEmail},
		Subject: "🎉 Welcome to Authgrid!",
		HTML:    generateWelcomeEmailHTML(customerName, subscriptionID),
	}

	// Send via Resend API
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("email send failed with status: %d", resp.StatusCode)
	}

	log.Printf("✅ Welcome email sent to: %s", customerEmail)
	return nil
}

// generateWelcomeEmailHTML creates the welcome email HTML
func generateWelcomeEmailHTML(customerName, subscriptionID string) string {
	if customerName == "" {
		customerName = "there"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; border-radius: 8px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; }
        .content { padding: 30px 0; }
        .button { display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { border-top: 1px solid #eee; padding-top: 20px; margin-top: 30px; font-size: 14px; color: #666; }
        .highlight { background: #f0f4ff; padding: 15px; border-left: 4px solid #667eea; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎉 Welcome to Authgrid!</h1>
        <p>You're officially customer #1 in our community</p>
    </div>

    <div class="content">
        <p>Hi %s! 👋</p>

        <p>Thank you for being an <strong>early adopter</strong> of Authgrid! You're helping us build the future of passwordless authentication.</p>

        <div class="highlight">
            <strong>🎁 As a founding customer, you get:</strong>
            <ul>
                <li><strong>14-day free trial</strong> - No charges during trial period</li>
                <li><strong>Lifetime 50%% discount</strong> - $14.50/mo forever (instead of $29)</li>
                <li><strong>Direct access to me</strong> - Email me anytime at hello@authgrid.net</li>
                <li><strong>Priority support</strong> - I'll personally help you get set up</li>
                <li><strong>Feature requests</strong> - Your input shapes the product roadmap</li>
            </ul>
        </div>

        <h2>🚀 Getting Started</h2>

        <p><strong>Option 1: Self-Hosted (Free Forever)</strong></p>
        <p>Deploy Authgrid to your own infrastructure in minutes:</p>
        <a href="https://github.com/Kelsidavis/authgrid" class="button">View GitHub Repo →</a>

        <p><strong>Option 2: Managed Hosting (Coming Soon)</strong></p>
        <p>We're building fully managed hosting so you don't have to worry about DevOps. You'll be the first to know when it's ready!</p>

        <h2>📚 Resources</h2>
        <ul>
            <li><a href="https://authgrid.org">Demo & Documentation</a></li>
            <li><a href="https://github.com/Kelsidavis/authgrid/blob/main/README.md">Quick Start Guide</a></li>
            <li><a href="https://github.com/Kelsidavis/authgrid/tree/main/examples">Code Examples</a></li>
        </ul>

        <h2>💬 Let's Chat!</h2>
        <p>I'd love to learn more about your use case and help you get set up. Reply to this email or schedule a 15-minute onboarding call with me:</p>
        <a href="mailto:hello@authgrid.net?subject=Onboarding Call Request&body=Hi! I'd like to schedule an onboarding call." class="button">Schedule Call →</a>

        <p>Questions? Just reply to this email - I read every message personally.</p>

        <p>Thanks again for your support! 🙏</p>

        <p>
            <strong>Kelsi Davis</strong><br>
            Founder, Authgrid<br>
            <a href="mailto:hello@authgrid.net">hello@authgrid.net</a>
        </p>
    </div>

    <div class="footer">
        <p><strong>Subscription Details:</strong></p>
        <p>
            Plan: Starter ($29/mo → $14.50/mo with founder discount)<br>
            Subscription ID: %s<br>
            Trial ends: 14 days from today<br>
        </p>
        <p>
            <a href="https://billing.stripe.com/p/login/test_YOUR_PORTAL_LINK">Manage Subscription</a> |
            <a href="https://authgrid.org/pricing">View Pricing</a> |
            <a href="https://github.com/Kelsidavis/authgrid">GitHub</a>
        </p>
        <p style="font-size: 12px; color: #999;">
            You're receiving this because you subscribed to Authgrid.<br>
            Authgrid · Passwordless authentication made simple
        </p>
    </div>
</body>
</html>
`, customerName, subscriptionID)
}
