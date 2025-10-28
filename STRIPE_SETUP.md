# Stripe Integration Setup Guide

Complete guide to setting up Stripe payment processing for Authgrid.

---

## Prerequisites

- Stripe account created at https://stripe.com
- Stripe API keys (you have the publishable key already)
- Access to Stripe Dashboard

---

## Step 1: Get Your Stripe API Keys

You already have your **publishable key** (starts with `pk_live_`):
```
pk_live_51SNFHRGWvG7njsH8o9vJcSxz9ADDGXt2yT1FFmzUXO3iH1dvvEpOyKOhTQNe4pAdsd7ygsuY1EH7fk7YdO47wyOT00smiNWsHU
```

Now get your **secret key**:

1. Go to https://dashboard.stripe.com/apikeys
2. Copy your **Secret key** (starts with `sk_live_`)
3. Keep it secure - you'll add it as an environment variable

---

## Step 2: Create Products and Prices in Stripe

### Option A: Via Stripe Dashboard (Recommended for beginners)

1. Go to https://dashboard.stripe.com/products
2. Click **"+ Add product"**

#### Create Starter Plan:

- **Name:** Authgrid Starter
- **Description:** Managed hosting with 10,000 MAU, custom domain, and email support
- **Pricing:**
  - Recurring: **Monthly**
  - Price: **$29.00 USD**
  - **Free trial:** 14 days
- Click **Save product**
- **Copy the Price ID** (starts with `price_...`) - you'll need this

#### Create Pro Plan:

- **Name:** Authgrid Pro
- **Description:** Everything in Starter plus 100,000 MAU, SSO, priority support, white-label
- **Pricing:**
  - Recurring: **Monthly**
  - Price: **$99.00 USD**
  - **Free trial:** 14 days
- Click **Save product**
- **Copy the Price ID** (starts with `price_...`)

---

### Option B: Via Stripe CLI (for advanced users)

Install Stripe CLI:
```bash
# macOS
brew install stripe/stripe-cli/stripe

# Login
stripe login
```

Create products:
```bash
# Starter plan
stripe products create \
  --name="Authgrid Starter" \
  --description="Managed hosting with 10,000 MAU"

stripe prices create \
  --product=<PRODUCT_ID_FROM_ABOVE> \
  --unit-amount=2900 \
  --currency=usd \
  --recurring[interval]=month

# Pro plan
stripe products create \
  --name="Authgrid Pro" \
  --description="Everything in Starter plus 100,000 MAU and advanced features"

stripe prices create \
  --product=<PRODUCT_ID_FROM_ABOVE> \
  --unit-amount=9900 \
  --currency=usd \
  --recurring[interval]=month
```

---

## Step 3: Configure Environment Variables on Fly.io

Set the Stripe keys and price IDs on your Fly.io app:

```bash
# Set Stripe secret key (REQUIRED)
flyctl secrets set STRIPE_SECRET_KEY="sk_live_YOUR_SECRET_KEY_HERE" -a authgrid-api

# Set price IDs (REQUIRED)
flyctl secrets set STRIPE_STARTER_PRICE_ID="price_YOUR_STARTER_PRICE_ID" -a authgrid-api
flyctl secrets set STRIPE_PRO_PRICE_ID="price_YOUR_PRO_PRICE_ID" -a authgrid-api

# Set webhook secret (you'll get this in Step 5)
flyctl secrets set STRIPE_WEBHOOK_SECRET="whsec_YOUR_WEBHOOK_SECRET" -a authgrid-api
```

**Note:** Setting secrets will automatically restart your app.

---

## Step 4: Deploy Updated Code

Deploy the Stripe-integrated backend:

```bash
cd /home/k/Desktop/authgrid
flyctl deploy
```

Wait for deployment to complete and verify:
```bash
curl https://authgrid-api.fly.dev/health
```

---

## Step 5: Set Up Stripe Webhooks

Webhooks notify your backend when payments succeed, subscriptions cancel, etc.

### Create Webhook Endpoint:

1. Go to https://dashboard.stripe.com/webhooks
2. Click **"+ Add endpoint"**
3. **Endpoint URL:** `https://authgrid-api.fly.dev/stripe-webhook`
4. **Events to send:**
   - `checkout.session.completed` - Customer completed checkout
   - `customer.subscription.created` - New subscription started
   - `customer.subscription.deleted` - Subscription canceled
   - `invoice.payment_succeeded` - Monthly payment succeeded
   - `invoice.payment_failed` - Payment failed
5. Click **Add endpoint**
6. **Copy the Signing secret** (starts with `whsec_...`)
7. Set it as an environment variable:
   ```bash
   flyctl secrets set STRIPE_WEBHOOK_SECRET="whsec_YOUR_WEBHOOK_SECRET" -a authgrid-api
   ```

---

## Step 6: Test in Test Mode (Recommended)

Before going live, test with Stripe test mode:

1. Switch pricing page to test mode:
   - Use test publishable key: `pk_test_...` instead of `pk_live_...`
2. Use test price IDs in your backend
3. Test checkout with these test cards:
   - **Success:** `4242 4242 4242 4242`
   - **Decline:** `4000 0000 0000 0002`
   - Any future expiry date, any 3-digit CVC

---

## Step 7: Go Live

Once testing is complete:

1. ✅ Pricing page already uses live publishable key
2. Set live secret key and price IDs on Fly.io (see Step 3)
3. Update webhook endpoint to use live mode
4. Deploy: `flyctl deploy`

---

## Step 8: Upload Pricing Page to authgrid.org

Your pricing page is ready at `examples/pricing/index.html`. Upload it to your web host:

**Option 1: Manual Upload (FTP/File Manager)**
- Upload `examples/pricing/index.html` to `https://authgrid.org/pricing/index.html`

**Option 2: GitHub Pages**
```bash
# Create gh-pages branch
git checkout --orphan gh-pages
cp examples/pricing/index.html index.html
git add index.html
git commit -m "Deploy pricing page"
git push origin gh-pages

# Configure custom domain in GitHub settings
```

**Option 3: Vercel/Netlify**
- Deploy the `examples/pricing/` folder
- Configure custom domain: authgrid.org

---

## Testing the Complete Flow

1. Visit https://authgrid.org/pricing
2. Click **"Start 14-Day Trial"** on Starter plan
3. Fill in test card: `4242 4242 4242 4242`
4. Complete checkout
5. Verify:
   - Redirected to success page
   - Subscription created in Stripe Dashboard
   - Webhook event received (check Fly.io logs: `flyctl logs -a authgrid-api`)

---

## Monitoring Payments

### View Subscriptions:
https://dashboard.stripe.com/subscriptions

### View Customers:
https://dashboard.stripe.com/customers

### View Revenue:
https://dashboard.stripe.com/balance

### View Webhook Logs:
https://dashboard.stripe.com/webhooks

### View API Logs:
```bash
flyctl logs -a authgrid-api
```

---

## Handling Customer Emails

When a customer signs up, you'll want to:

1. **Collect their email** during checkout (already configured in `stripe.go`)
2. **Send welcome email** via webhook handler
3. **Provision their account** (create database entry, generate API keys, etc.)

Next step: Implement webhook handler in `src/api/stripe.go` to:
- Create customer account in your database
- Send welcome email
- Provision managed hosting instance

---

## Pricing Summary

| Plan | Price | MAU | Free Trial |
|------|-------|-----|------------|
| Free (Self-Hosted) | $0 | Unlimited | N/A |
| Starter | $29/mo | 10,000 | 14 days |
| Pro | $99/mo | 100,000 | 14 days |
| Enterprise | Custom | Unlimited | Custom |

---

## Next Steps After First Payment

1. **Send personal thank you email** to first customer
2. **Offer lifetime discount** (50% off forever - $14.50/mo)
3. **Schedule onboarding call**
4. **Ask for testimonial** after they're happy
5. **Use as case study** in marketing

---

## Troubleshooting

### "Payment system not configured"
- Check that `STRIPE_SECRET_KEY` is set: `flyctl secrets list -a authgrid-api`

### Webhook events not received
- Verify webhook URL is correct: `https://authgrid-api.fly.dev/stripe-webhook`
- Check webhook signing secret is set
- View webhook logs in Stripe Dashboard

### Checkout button does nothing
- Check browser console for errors
- Verify API endpoint is reachable: `curl https://authgrid-api.fly.dev/health`
- Check CORS is configured for authgrid.org

---

## Security Best Practices

✅ **DO:**
- Keep secret keys in environment variables (never commit to git)
- Verify webhook signatures
- Use HTTPS everywhere
- Enable Stripe Radar for fraud detection

❌ **DON'T:**
- Hardcode secret keys in source code
- Skip webhook signature verification
- Store credit card numbers (Stripe handles this)

---

## Resources

- **Stripe Dashboard:** https://dashboard.stripe.com
- **Stripe Docs:** https://stripe.com/docs
- **Stripe API Reference:** https://stripe.com/docs/api
- **Stripe Testing:** https://stripe.com/docs/testing

---

**You're ready to accept your first payment! 🚀**

Next: Set your Stripe secret key and price IDs, then deploy.
