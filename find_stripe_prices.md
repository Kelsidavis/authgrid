# Finding Your Stripe Price IDs

You created the products, but we need the **price IDs** (not product IDs).

## Option 1: Find Price IDs in Dashboard

### For Starter Plan:
1. Go to: https://dashboard.stripe.com/products/prod_TJthGQvpiFpAH4
2. Look under the "Pricing" section
3. You'll see your $29/month price listed
4. Click on the price
5. Copy the **API ID** field (starts with `price_`)

### For Pro Plan:
1. Go to: https://dashboard.stripe.com/products/prod_TJtkYsnecNB8hH
2. Look under "Pricing" section
3. Click on the $99/month price
4. Copy the **API ID**

---

## Option 2: Create New Prices with Free Trials

If you want to add free trials (recommended!), create new prices:

### For Starter Plan ($29/mo):
1. Go to: https://dashboard.stripe.com/products/prod_TJthGQvpiFpAH4
2. Under "Pricing" section, click **"+ Add another price"**
3. Fill in:
   - **Price:** 29.00
   - **Billing period:** Recurring → Monthly
   - **Currency:** USD
4. **Click "Additional options"** (expand this section)
5. Under **"Free trial"**, enter: 14 days
6. Click **"Add price"**
7. Copy the new **price ID** (starts with `price_`)

### For Pro Plan ($99/mo):
1. Go to: https://dashboard.stripe.com/products/prod_TJtkYsnecNB8hH
2. Click **"+ Add another price"**
3. Fill in:
   - **Price:** 99.00
   - **Billing period:** Recurring → Monthly
   - **Currency:** USD
4. **Click "Additional options"**
5. Under **"Free trial"**, enter: 14 days
6. Click **"Add price"**
7. Copy the new **price ID**

---

## Note About Free Trials

Free trials can be set in two places:
1. **On the price** (when creating/editing the price) - this is a default trial
2. **At checkout time** (in the code) - we already configured this in stripe.go

If you don't see the free trial option when creating prices, that's OK - our backend code already adds 14-day trials at checkout time (see `src/api/stripe.go` line 67).

---

## Once You Have the Price IDs

They'll look like: `price_1234567890AbCdEf`

Reply with both price IDs and I'll set them up!
