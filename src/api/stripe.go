package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
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

	payload, err := json.Marshal(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify webhook signature
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("WARNING: STRIPE_WEBHOOK_SECRET not set - webhook signature verification skipped")
	}

	// TODO: Implement webhook event handling
	// Example events to handle:
	// - checkout.session.completed: User completed payment
	// - customer.subscription.deleted: User canceled subscription
	// - invoice.payment_succeeded: Monthly payment succeeded
	// - invoice.payment_failed: Payment failed

	log.Printf("Received Stripe webhook: %s", string(payload))

	respondJSON(w, http.StatusOK, map[string]string{"received": "true"})
}
