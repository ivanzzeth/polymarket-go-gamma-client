# Low Liquidity / High Volume Market Making Example

## Overview

This example identifies **market making opportunities** by finding markets with high trading volume but insufficient liquidity - a classic sign that market makers can earn significant spreads.

## What Makes This a Good Market Making Opportunity?

### The Market Making Opportunity Formula

```
High Volume + Low Liquidity = Wide Spreads = Market Making Profits
```

**Why this combination works:**

1. **High Volume** = Strong demand and active traders
   - Many participants want to trade
   - Frequent transactions = more spread capture
   - Consistent order flow

2. **Low Liquidity** = Insufficient supply of market makers
   - Wider bid-ask spreads
   - Less competition
   - Higher profit per trade

3. **Result** = Profitable market making
   - Capture spreads repeatedly
   - High turnover on capital
   - Lower inventory risk (volume helps rebalance)

## Market Making Fundamentals

### What is Market Making?

Market making is the practice of:
1. Quoting both **bid** (buy) and **ask** (sell) prices
2. Capturing the **spread** between them
3. **Rebalancing** inventory as needed
4. Earning profits through high turnover

### Example Trade Cycle

```
Market: "Will it rain tomorrow?"
Current spread: 0.48 bid / 0.52 ask (4¢ spread)

Your quotes:
- Bid: 0.49 (willing to buy)
- Ask: 0.51 (willing to sell)

Trade sequence:
1. Someone sells to you: Buy at 0.49
2. Someone buys from you: Sell at 0.51
3. Profit: 0.02 per share (2¢ or 4% on capital)

If this happens 50 times per day:
Daily profit: 50 × 0.02 = 1.00 per share
On $1,000 position: $100/day potential
Monthly: $3,000 (if sustained)
```

### Key Metrics

**Turnover Ratio** = Volume / Liquidity
```
High turnover (>5x): Excellent for market making
Medium turnover (2-5x): Good opportunity
Low turnover (<2x): May not be worth it

Example:
Volume: $50,000
Liquidity: $5,000
Turnover: 10x
Interpretation: Capital turns over 10 times, high profit potential
```

**Spread as % of Price**
```
Wide spread (>3%): Excellent margins
Moderate spread (1.5-3%): Good margins
Tight spread (<1.5%): Competitive, harder to profit

Example:
Price: 0.50
Spread: 0.04
Spread %: 8%
Interpretation: Very profitable if you can capture it
```

## How to Use This Example

### 1. Run the Example

```bash
cd examples/find-low-liquidity-high-volume
go run main.go
```

### 2. Understanding the Output

The example finds markets with:
- ✅ 24-hour volume > $10,000 (strong activity)
- ✅ Liquidity < $20,000 (room for more market makers)
- ✅ Minimum liquidity > $2,000 (enough to start)
- ✅ Open and accepting orders
- ✅ Active trading

For each opportunity, it displays:

#### Market Information
- Question and market ID
- Category and visibility (featured/trending)
- Direct link to market

#### Volume Metrics
- 24-hour volume (activity level)
- Total volume (overall interest)
- Volume trend (increasing/stable/declining)

#### Liquidity Analysis
- Current total liquidity
- CLOB vs AMM liquidity breakdown
- Volume-to-liquidity ratio (turnover)
- Liquidity assessment

#### Spread Analysis
- Current spread (profit opportunity)
- Ticks in spread (competition room)
- Spread quality rating
- Comparative assessment

#### Market Making Score
A 0-100 score based on:
- Volume levels (higher = better)
- Volume-to-liquidity ratio (higher = better)
- Spread width (wider = better, to a point)
- Market activity (more recent = better)
- Accepting orders (required)

### 3. Interpreting Scores

#### Excellent Opportunity (80-100)
```
Characteristics:
- Very high volume (>$50k/day)
- Low liquidity (<$10k)
- Wide spread (>2%)
- High turnover (>8x)

Example:
24h Volume: $125,000
Liquidity: $7,500
Turnover: 16.7x
Spread: 3.5%
Score: 92

Action: Strong candidate for market making
```

#### Very Good Opportunity (70-80)
```
Characteristics:
- High volume ($25k-50k/day)
- Low-moderate liquidity ($10k-15k)
- Good spread (1.5-2%)
- Good turnover (5-8x)

Example:
24h Volume: $35,000
Liquidity: $12,000
Turnover: 2.9x
Spread: 2.2%
Score: 75

Action: Solid opportunity
```

#### Good Opportunity (60-70)
```
Characteristics:
- Moderate volume ($10k-25k/day)
- Moderate liquidity ($15k-20k)
- Adequate spread (1-1.5%)
- Moderate turnover (3-5x)

Example:
24h Volume: $18,000
Liquidity: $15,500
Turnover: 1.2x
Spread: 1.5%
Score: 65

Action: Consider if you have experience
```

## Step-by-Step Market Making Guide

### Phase 1: Opportunity Assessment

#### Step 1: Evaluate the Market

```
Questions to answer:

📊 Volume Analysis:
✅ Is volume consistently high or one-time spike?
✅ Check multiple days of history
✅ Identify if event-driven or sustained interest

💰 Liquidity Analysis:
✅ Current total liquidity
✅ CLOB vs AMM (CLOB is more competitive)
✅ Your potential market share
✅ Other active market makers?

📈 Spread Analysis:
✅ Current spread width
✅ Recent spread history
✅ Number of ticks (room for competition)
✅ Typical order sizes

🎯 Risk Analysis:
✅ Price stability (volatile = more risk)
✅ Resolution timeline (longer = more opportunity)
✅ Information asymmetry risk
✅ Event complexity
```

#### Step 2: Calculate Potential Returns

```
Example Market:
- 24h Volume: $50,000
- Current Spread: 2% (0.02)
- Your planned liquidity: $5,000
- Estimated market share: 20%

Daily volume capture: $50,000 × 20% = $10,000
Spread captured: $10,000 × 2% = $200/day
Monthly projection: $200 × 30 = $6,000
ROI on $5,000: 120% per month

Reality checks:
- Fees: ~2% of volume = $200/month = -$200
- Net projection: $5,800
- ROI: 116% per month (excellent!)
- Risk: Inventory risk, adverse selection

Conservative estimate: 50% of theoretical
- Expected monthly: $2,900
- ROI: 58% per month (still excellent)
```

#### Step 3: Size Your Position

```
Conservative Approach:
- Start with 5-10% of current liquidity
- Market liquidity: $10,000
- Your capital: $500-1,000
- Quotes: $250-500 per side

Rationale:
- Test the market first
- Learn the dynamics
- Minimize risk
- Scale up if successful

Moderate Approach:
- Use 10-20% of liquidity
- Market liquidity: $10,000
- Your capital: $1,000-2,000
- Quotes: $500-1,000 per side

Aggressive Approach (not recommended initially):
- Use 25%+ of liquidity
- Risk: Market impact, harder to rebalance
- Only if very confident
```

### Phase 2: Initial Setup

#### Step 1: Analyze Current Order Book

```
Check:
- Current best bid
- Current best ask
- Depth on each side
- Recent trade prices
- Order sizes (typical)

Example order book:
Bids:
  0.48: $500
  0.47: $1,200
  0.46: $800

Asks:
  0.52: $600
  0.53: $1,000
  0.54: $1,500

Spread: 0.04 (4¢)
Mid: 0.50

Your strategy:
Bid: 0.49 ($400) - inside current best
Ask: 0.51 ($400) - inside current best
Your spread: 0.02 (2¢) - tighter than market
```

#### Step 2: Set Your Quotes

```
Quote Placement Strategy:

Option A: Inside the spread (aggressive)
- Bid: 0.49 (vs market 0.48)
- Ask: 0.51 (vs market 0.52)
- Pros: Get filled first, high turnover
- Cons: Lower margin, may move market

Option B: Join the best (moderate)
- Bid: 0.48 (match current best)
- Ask: 0.52 (match current best)
- Pros: Still competitive, better margin
- Cons: May not get filled first

Option C: Slightly wider (conservative)
- Bid: 0.47 (below current best)
- Ask: 0.53 (above current best)
- Pros: Higher margin, less risk
- Cons: Lower fill rate, less turnover

Recommended: Start with Option B, adjust based on results
```

#### Step 3: Place Initial Orders

```
Execution checklist:
✅ Place both bid and ask simultaneously
✅ Equal sizes on both sides (balanced)
✅ Use limit orders (not market)
✅ Set reasonable expiration (not immediate-or-cancel)
✅ Monitor for fills

Example using CLOB API:
Buy order:
  Side: BUY
  Price: 0.49
  Size: 1000 shares ($490)
  Type: LIMIT

Sell order:
  Side: SELL
  Price: 0.51
  Size: 1000 shares ($510 collateral needed)
  Type: LIMIT

Total capital needed: ~$1,000
```

### Phase 3: Active Management

#### Step 1: Monitor and Rebalance

```
Imbalanced inventory example:

Starting: 1000 shares, balanced

After 2 hours:
- Bought: 500 more (now 1500 shares)
- Sold: 200 (now 1300 shares)
- Net position: +300 shares (too long)

Risk: If price drops, you lose on inventory

Rebalancing actions:
1. Widen bid quote (slow down buying)
   - Was: 0.49, now: 0.48
2. Tighten ask quote (encourage selling)
   - Was: 0.51, now: 0.505
3. Or: Take liquidity to rebalance
   - Sell 300 shares at market bid

Goal: Stay near neutral (±10% of capital)
```

#### Step 2: Adjust to Market Conditions

```
Scenario A: Market moving up
- Price: 0.50 → 0.55
- Action: Move quotes up
  - New bid: 0.54
  - New ask: 0.56
- Risk: Don't fight the trend

Scenario B: Spread tightening
- Spread: 0.04 → 0.02
- Reason: More market makers entering
- Action:
  - Accept lower margin, or
  - Reduce position size, or
  - Exit opportunity (no longer profitable)

Scenario C: Volume declining
- Volume: $50k → $20k/day
- Action:
  - Reduce position size
  - Widen quotes (lower turnover = need more margin)
  - Consider exiting

Scenario D: Volatility increasing
- Price swings: ±2% → ±5%
- Action:
  - Widen quotes for protection
  - Reduce position size
  - Rebalance more frequently
```

#### Step 3: Manage Fills and Inventory

```
Fill management:

Every 15-30 minutes:
✅ Check which orders filled
✅ Calculate current position
✅ Assess inventory skew
✅ Replace filled orders
✅ Adjust for market movement

Example:
Time: 10:00 AM
- Bid filled: Bought 500 @ 0.49
- Position: +500 shares (long)
- Action:
  - Replace bid: 0.48 (wider, slow down)
  - Keep ask: 0.51
  - Or sell 500 at market to rebalance

Time: 10:30 AM
- Ask filled: Sold 700 @ 0.51
- Position: -200 shares (short from previous +500)
- Profit: 500 × (0.51-0.49) = $10
- Action:
  - Replace ask: 0.51
  - Tighten bid: 0.495 (encourage buying)
```

### Phase 4: Profit Taking and Exit

#### Daily Reconciliation

```
End of Day:
- Total buys: 2,500 shares @ avg 0.485
- Total sells: 2,300 shares @ avg 0.515
- Net position: +200 shares @ avg 0.485
- Gross profit: 2,300 × (0.515-0.485) = $69
- Unrealized: 200 × (current-0.485)
- Fees: ~$25 (2% of $1,250 volume)
- Net profit: $44/day

Current position value: 200 × 0.50 = $100
Available capital: $900
```

#### Exit Strategies

```
Strategy 1: Full exit
- Close all positions
- Cancel all orders
- Realize all P&L
- When: End of opportunity, going on break

Strategy 2: Partial exit
- Reduce position size
- Keep presence in market
- When: Competition increasing, spreads tightening

Strategy 3: Wind down
- Stop placing new orders
- Let existing orders get filled
- Gradually exit inventory
- When: Approaching resolution

Strategy 4: Scale out
- Reduce size over multiple days
- Take profits gradually
- When: Very large position, avoid market impact
```

## Real-World Example Analysis

### Example from Output:

```
Market: "Will Trump deport 1.5M+ by...?"
Category: Politics
24h Volume: $125,000 (HIGH activity ✅)
Liquidity: $7,500 (LOW liquidity ✅)

Metrics:
- Volume/Liquidity: 16.7x (EXCELLENT turnover)
- Spread: 3.5% (WIDE spread)
- Best Bid: 0.12
- Best Ask: 0.155
- Spread: 0.035 (3.5¢)

Score: 92/100 (EXCELLENT)
```

### Opportunity Assessment:

**Pros:**
✅ Extremely high turnover (16.7x)
✅ Wide spread (3.5%) = good margins
✅ High volume ($125k/day) = many fills
✅ Low competition (only $7.5k liquidity)
✅ Active market (accepting orders)

**Risks:**
⚠️ Political market (information events)
⚠️ May have informed traders
⚠️ Could be volatile
⚠️ Resolution timeline unclear

### Position Sizing:

```
Conservative:
- Capital: $500 (6.7% of liquidity)
- Per side: $250
- Shares: ~2,000 per side

Moderate:
- Capital: $1,000 (13.3% of liquidity)
- Per side: $500
- Shares: ~4,000 per side

Expected monthly P&L (moderate):
- Daily capture: 13.3% × $125k = $16,625
- Spread capture: $16,625 × 3.5% = $582
- Days/month: 30
- Gross: $17,460
- Fees: ~$333
- Net: $17,127
- ROI: 1,713% (likely too optimistic!)

Realistic (conservative estimate at 10% of theoretical):
- Monthly net: $1,713
- ROI: 171% per month
- Still excellent if achievable!
```

### Execution Plan:

**Week 1 (Testing):**
```
Capital: $250
Goal: Learn market dynamics
Monitor: Fill rates, rebalancing needs
Adjust: Quote placement strategy
Target: Break even to +20%
```

**Week 2-3 (Scaling):**
```
Capital: $500 (if Week 1 successful)
Goal: Optimize quote placement
Monitor: Competition changes
Adjust: Position size, spread
Target: +40-60% on capital
```

**Week 4+ (Optimized):**
```
Capital: $1,000 (if sustained success)
Goal: Maximum sustainable size
Monitor: Market conditions closely
Adjust: As needed
Target: +80-120% on capital
```

## Risk Management

### Inventory Risk

```
Problem: Getting stuck long or short

Prevention:
- Set max position limits (±20% of capital)
- Rebalance regularly (every 30 min)
- Widen quotes when skewed
- Take liquidity if needed

Example limits:
Capital: $1,000
Max long: +$200 (200 shares if price ~1.0)
Max short: -$200
Alert at: ±$150
Forced rebalance at: ±$200
```

### Adverse Selection Risk

```
Problem: Informed traders hitting your quotes

Signs:
- Your quotes get hit immediately
- Market moves against you after fill
- Consistent losses on inventory

Protection:
- Don't be the tightest quote (join, don't lead)
- Widen quotes after fast fills
- Reduce size after losing trades
- Watch for news/information events
- Cancel quotes during major announcements
```

### Spread Compression Risk

```
Problem: More market makers enter, spreads tighten

Signs:
- Spread: 4% → 2% → 1%
- More competing quotes
- Your fill rate drops

Response:
- Accept lower margins if still profitable
- Increase turnover to compensate
- Reduce position size (lower risk = accept lower return)
- Or exit opportunity if no longer profitable

Example:
Previous: 4% spread, 10x turnover = 40% monthly
New: 2% spread, 12x turnover = 24% monthly
Decision: Still good? Or find better opportunity?
```

### Resolution Risk

```
Problem: Market resolves while you hold inventory

Prevention:
- Monitor resolution timeline
- Reduce exposure as resolution approaches
- Close positions 1-3 days before resolution
- Set calendar reminders

Example:
Resolution date: Dec 31
Current date: Dec 28 (3 days)
Current position: +500 shares

Action:
- Dec 28: Reduce to +200 shares
- Dec 29: Reduce to +50 shares
- Dec 30: Close completely
- Reason: Don't want inventory risk at resolution
```

## Advanced Techniques

### 1. Dynamic Spread Adjustment

```
Based on volatility:

Low volatility (stable prices):
- Use tighter spreads (1-1.5%)
- Higher turnover
- Lower margin per trade

High volatility (price swings):
- Use wider spreads (2-4%)
- Lower turnover
- Higher margin per trade

Example:
Normal: 0.49 / 0.51 (2% spread)
Volatile: 0.47 / 0.53 (6% spread)

Trade-off: Fewer fills but safer
```

### 2. Layered Quoting

```
Instead of single bid/ask:

Multiple levels:
Bids:
  0.49: $200 (tight)
  0.48: $300 (moderate)
  0.47: $500 (wide)

Asks:
  0.51: $200 (tight)
  0.52: $300 (moderate)
  0.53: $500 (wide)

Benefits:
- Capture more volume
- Better average prices
- More information on order flow
- Smooths rebalancing
```

### 3. Automated Rebalancing

```
Rules-based approach:

If position > +$150:
  - Widen bid by 1 tick
  - Tighten ask by 1 tick

If position < -$150:
  - Tighten bid by 1 tick
  - Widen ask by 1 tick

If position > +$200:
  - Cancel bid
  - Hit best bid with 100 shares

If position < -$200:
  - Cancel ask
  - Hit best ask with 100 shares

Benefit: Maintains neutrality automatically
```

### 4. Time-Based Strategies

```
Time of day effects:

Morning (8am-12pm):
- Higher volume
- Wider spreads
- More opportunities
- Strategy: More aggressive

Afternoon (12pm-6pm):
- Moderate volume
- Tighter spreads
- Strategy: Balanced

Evening (6pm-12am):
- Lower volume
- Can widen spreads
- Strategy: Conservative

Night (12am-8am):
- Very low volume
- Wider spreads okay
- Strategy: Wider quotes, smaller size
```

## Performance Tracking

### Daily Metrics to Track

```
Trading metrics:
- Total volume captured
- Average spread captured
- Number of round trips
- Fill rates (bid vs ask)
- Average inventory held
- Max inventory deviation

Example log:
Date: 2025-01-15
Volume captured: $15,000
Avg spread: 2.3%
Round trips: 12
Gross profit: $345
Fees: $30
Net profit: $315
ROI: 31.5% (on $1,000 capital for 1 day)
```

### Weekly Review

```
Questions to answer:
- Is profitability consistent?
- Are fill rates acceptable?
- Is inventory management working?
- Has competition increased?
- Are spreads tightening?
- Should position size change?
- Is the opportunity still attractive?

Actions:
- Adjust position size
- Modify quote strategy
- Rebalance parameters
- Consider exit if degrading
```

## When to Exit

### Exit Signals

```
✅ Time to exit when:
- Spread < 1% (too competitive)
- Volume dropping significantly
- Multiple new market makers entering
- Consistent losses on inventory
- Approaching resolution
- Better opportunities elsewhere
- Strategy not working after 1 week

Exit checklist:
☐ Close all inventory
☐ Cancel all orders
☐ Verify no pending trades
☐ Calculate final P&L
☐ Document lessons learned
```

## Common Mistakes to Avoid

### ❌ Mistake 1: Oversizing Position
```
Market liquidity: $10,000
Your capital: $8,000 (80%!)

Problems:
- You become the market
- Hard to rebalance
- High inventory risk
- Market impact on every trade

Correct: Start with $500-1,000 (5-10%)
```

### ❌ Mistake 2: Fighting the Trend
```
Scenario:
- You're long 500 shares @ 0.50
- Price drops to 0.45
- You keep buying "to average down"
- Price drops to 0.40
- You're now long 1,500 @ 0.45 avg
- Massive unrealized loss

Correct: Cut losses, rebalance, or exit
```

### ❌ Mistake 3: Ignoring Fees
```
Bad calculation:
- Spread captured: 2%
- "I make 2% per trade!"

Reality:
- Spread: 2%
- Maker fee: 2%
- Net: 0% (or negative with taker fees)

Correct: Factor all fees into strategy
```

### ❌ Mistake 4: Set and Forget
```
Wrong:
- Place orders at 9am
- Go to lunch
- Come back at 5pm
- Wonder why you're down money

Correct:
- Monitor every 15-30 minutes
- Rebalance regularly
- Adjust for market conditions
- Active management required
```

### ❌ Mistake 5: Not Validating Volume
```
Trap:
- Sees "24h volume: $100k"
- Assumes it will continue
- Volume next day: $5k
- Opportunity gone

Correct:
- Check multiple days history
- Identify volume patterns
- Understand why volume is high
- Prepared for changes
```

## Conclusion

Low liquidity / high volume markets offer excellent market making opportunities when:

✅ Volume is consistently high (not one-time spike)
✅ Liquidity is low enough for wide spreads
✅ Spreads are wide enough to profit after fees
✅ You can actively manage the position
✅ Market conditions are stable enough

Success requires:
- Proper position sizing (start small!)
- Active monitoring and rebalancing
- Quick response to market changes
- Disciplined risk management
- Tracking and optimization

Start conservatively, learn the dynamics, and scale up gradually. Market making can be highly profitable but requires attention and discipline.

## Related Examples

- [Wide Spread Markets](../find-wide-spread-markets/) - Another market making opportunity type
- [New Active Markets](../find-new-active-markets/) - Early market maker advantages
- [Related Markets Arbitrage](../find-related-markets-arbitrage/) - Different strategy (arbitrage vs market making)

## Additional Resources

- [Polymarket CLOB API](https://docs.polymarket.com/#clob-api) - For placing orders
- [Market Making Basics](https://docs.polymarket.com/#market-making) - General concepts
- Main repository README for other strategies
