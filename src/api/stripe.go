package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	stripeWebhook "github.com/stripe/stripe-go/v76/webhook"
)

// Plan prices (in cents)
const (
	StarterPriceID = "price_starter_monthly" // You'll create this in Stripe Dashboard
	ProPriceID     = "price_pro_monthly"     // You'll create this in Stripe Dashboard
)

// CheckoutRequest represents the checkout request from frontend
type CheckoutRequest struct {
	Plan string `json:"plan"` // "starter" or "pro"
}

// createCheckoutSessionHandler creates a Stripe checkout session
func createCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate plan
	if req.Plan != "starter" && req.Plan != "pro" {
		respondError(w, http.StatusBadRequest, "Invalid plan. Must be 'starter' or 'pro'")
		return
	}

	// Set Stripe API key from environment
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Println("ERROR: STRIPE_SECRET_KEY environment variable not set")
		respondError(w, http.StatusInternalServerError, "Payment system not configured")
		return
	}
	stripe.Key = stripeKey

	// Determine price ID based on plan
	var priceID string
	if req.Plan == "starter" {
		priceID = os.Getenv("STRIPE_STARTER_PRICE_ID")
		if priceID == "" {
			priceID = StarterPriceID // Fallback to constant
		}
	} else if req.Plan == "pro" {
		priceID = os.Getenv("STRIPE_PRO_PRICE_ID")
		if priceID == "" {
			priceID = ProPriceID // Fallback to constant
		}
	}

	// Create Stripe checkout session
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("https://authgrid.org/pricing?success=true&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://authgrid.org/pricing?canceled=true"),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			TrialPeriodDays: stripe.Int64(14), // 14-day free trial
		},
		AllowPromotionCodes:      stripe.Bool(true),
		BillingAddressCollection: stripe.String("auto"),
		// CustomerEmail can be added here if user is logged in
	}

	sess, err := session.New(params)
	if err != nil {
		log.Printf("Stripe checkout session creation failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create checkout session")
		return
	}

	// Return session ID to frontend
	respondJSON(w, http.StatusOK, map[string]string{
		"id": sess.ID,
	})
}

// stripeWebhookHandler handles Stripe webhook events
// This is called when payments succeed, subscriptions are canceled, etc.
func stripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set Stripe API key
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Println("ERROR: STRIPE_SECRET_KEY not set")
		respondError(w, http.StatusInternalServerError, "Payment system not configured")
		return
	}
	stripe.Key = stripeKey

	// Verify webhook signature (important for security!)
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	var event stripe.Event

	if webhookSecret != "" {
		// Verify signature
		signatureHeader := r.Header.Get("Stripe-Signature")
		event, err = stripeWebhook.ConstructEvent(payload, signatureHeader, webhookSecret)
		if err != nil {
			log.Printf("Webhook signature verification failed: %v", err)
			respondError(w, http.StatusBadRequest, "Invalid signature")
			return
		}
	} else {
		// No signature verification (not recommended for production!)
		log.Println("WARNING: Webhook signature verification skipped - STRIPE_WEBHOOK_SECRET not set")
		err = json.Unmarshal(payload, &event)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
	}

	// Handle different event types
	switch event.Type {
	case "checkout.session.completed":
		handleCheckoutCompleted(event)
	case "customer.subscription.deleted":
		handleSubscriptionDeleted(event)
	case "invoice.payment_succeeded":
		handlePaymentSucceeded(event)
	case "invoice.payment_failed":
		handlePaymentFailed(event)
	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
	}

	respondJSON(w, http.StatusOK, map[string]string{"received": "true"})
}

// handleCheckoutCompleted processes successful checkouts
func handleCheckoutCompleted(event stripe.Event) {
	var session stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &session)
	if err != nil {
		log.Printf("Error parsing checkout.session.completed: %v", err)
		return
	}

	log.Printf("✅ Payment successful!")
	log.Printf("   Customer: %s", session.CustomerDetails.Email)
	log.Printf("   Subscription: %s", session.Subscription.ID)
	log.Printf("   Amount: $%.2f", float64(session.AmountTotal)/100)

	// TODO: Implement fulfillment logic
	// 1. Create customer record in database
	// 2. Generate API keys or provision managed instance
	// 3. Send welcome email

	// For now, log what should happen
	log.Printf("🎉 NEW CUSTOMER! Next steps:")
	log.Printf("   1. Send welcome email to: %s", session.CustomerDetails.Email)
	log.Printf("   2. Provision managed hosting instance")
	log.Printf("   3. Generate API credentials")
	log.Printf("   4. Send onboarding guide")
}

// handleSubscriptionDeleted processes subscription cancellations
func handleSubscriptionDeleted(event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		log.Printf("Error parsing customer.subscription.deleted: %v", err)
		return
	}

	log.Printf("❌ Subscription canceled: %s", subscription.ID)
	log.Printf("   Customer: %s", subscription.Customer.Email)

	// TODO: Deprovision resources, send cancellation confirmation email
}

// handlePaymentSucceeded processes successful recurring payments
func handlePaymentSucceeded(event stripe.Event) {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		log.Printf("Error parsing invoice.payment_succeeded: %v", err)
		return
	}

	log.Printf("💰 Payment succeeded: $%.2f", float64(invoice.AmountPaid)/100)
	log.Printf("   Customer: %s", invoice.CustomerEmail)

	// TODO: Send receipt email
}

// handlePaymentFailed processes failed payments
func handlePaymentFailed(event stripe.Event) {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		log.Printf("Error parsing invoice.payment_failed: %v", err)
		return
	}

	log.Printf("⚠️  Payment failed!")
	log.Printf("   Customer: %s", invoice.CustomerEmail)
	log.Printf("   Amount: $%.2f", float64(invoice.AmountDue)/100)

	// TODO: Send payment failed email, retry logic
}
