# Rapid Price Movement Trading Example

## Overview

This example identifies **price anomaly trading opportunities** by detecting markets with significant 24-hour price changes, helping you capitalize on momentum or mean reversion patterns.

## Trading Strategies for Price Movements

### Two Main Approaches

#### 1. Mean Reversion Strategy
**Thesis**: Overreactions return to fair value

```
Market suddenly moves:
Price: 0.30 → 0.50 (+67% in 24h)

Assumption: Overreaction, will reverse
Action: SELL at 0.50
Target: Price returns to 0.35-0.40
Profit: Capture the reversion
```

#### 2. Momentum Strategy
**Thesis**: Trends continue in the same direction

```
Market starts moving:
Price: 0.30 → 0.50 (+67% in 24h)

Assumption: Strong trend, will continue
Action: BUY at 0.50
Target: Price continues to 0.60-0.70
Profit: Ride the momentum
```

### Which Strategy to Use?

The tool analyzes **short-term (1-day) vs longer-term (1-week)** trends to suggest the appropriate strategy:

```
Scenario A: Mean Reversion Signal
- 1-day change: +20% (rapid spike)
- 1-week change: +5% (small overall change)
- Interpretation: Recent spike, likely overreaction
- Strategy: MEAN REVERSION (sell the spike)

Scenario B: Momentum Signal
- 1-day change: +20% (rapid increase)
- 1-week change: +40% (sustained uptrend)
- Interpretation: Part of larger trend
- Strategy: MOMENTUM (buy the trend)

Scenario C: Unclear/Uncertain
- 1-day change: +20%
- 1-week change: +15%
- Interpretation: Mixed signals
- Strategy: WAIT or use smaller position
```

## How to Use This Example

### 1. Run the Example

```bash
cd examples/find-rapid-price-movement
go run main.go
```

### 2. Configuration Options

Edit `main.go` to adjust parameters:

```go
minPriceChange := 0.15    // 15% minimum price change
targetCount := 5          // Number of opportunities to find
minVolume := 5000.0       // Minimum 24h volume
```

### 3. Understanding the Output

The example finds markets with:
- ✅ 24-hour price change ≥ 15% (significant movement)
- ✅ Sufficient volume (>$5,000/day)
- ✅ Open and accepting orders
- ✅ Clear price history data

For each opportunity, it displays:

#### Market Information
- Question and category
- Market ID and slug
- Link to Polymarket
- Featured/trending status

#### Price Movement Analysis
- Current price
- 24-hour price change (absolute and %)
- 1-week price change (for context)
- Direction (UP or DOWN)

#### Movement Characterization
```
Strong Movement (>30%):
   "SIGNIFICANT movement" or "VERY STRONG"

Moderate Movement (20-30%):
   "SUBSTANTIAL" or "STRONG"

Notable Movement (15-20%):
   "SIGNIFICANT" or "NOTABLE"
```

#### Trading Strategy Recommendation

Based on comparative analysis:

**Mean Reversion Signal:**
```
💡 Trading Strategy:
   Type:             MEAN REVERSION (Potential Overreaction)

   Analysis:
   • 24h change (20%) significantly exceeds 1w change (5%)
   • Suggests recent spike may be temporary overreaction
   • Price may revert toward longer-term trend

   Suggested Action:
   • Consider SELLING / SHORT at current price (0.65)
   • Set target: 0.55-0.60 (reversion to pre-spike levels)
   • Stop loss: 0.70 (if trend actually continues)
   • Position size: Conservative (this is counter-trend)
```

**Momentum Signal:**
```
💡 Trading Strategy:
   Type:             MOMENTUM (Trend Continuation)

   Analysis:
   • 24h change (20%) aligns with 1w uptrend (35%)
   • Suggests strong, sustained directional movement
   • Trend shows no signs of reversing

   Suggested Action:
   • Consider BUYING / LONG at current price (0.65)
   • Set target: 0.75-0.80 (continuation of trend)
   • Stop loss: 0.60 (if trend reverses)
   • Position size: Moderate (trend-following)
```

**Unclear Signal:**
```
💡 Trading Strategy:
   Type:             UNCERTAIN (Mixed Signals)

   Analysis:
   • 24h change (20%) conflicts with 1w change (-5%)
   • Short-term and long-term trends don't align
   • Could be reversal or consolidation

   Suggested Action:
   • WAIT for clearer confirmation
   • Monitor for additional signals
   • If trading, use very small position
   • Consider analyzing longer timeframes
```

#### Volume and Liquidity Context
- Trading volume (24h, 1w, total)
- Liquidity levels
- Assessment of tradability

## Step-by-Step Trading Guides

### Strategy 1: Mean Reversion Trading

Mean reversion assumes prices that move rapidly away from their average will return.

#### When to Use Mean Reversion

✅ **Good signals:**
- Large 1-day spike (>20%) with small 1-week change (<10%)
- No fundamental news justifying the move
- Historical pattern of volatility and reversals
- Price at extreme levels (near 0.05 or 0.95)
- Volume spike accompanying price spike

❌ **Avoid when:**
- Major news event causing the move
- Sustained trend over multiple timeframes
- Near resolution date (fundamentals matter more)
- Low liquidity (can't enter/exit easily)

#### Step 1: Validate the Signal

```
Example Market:
"Will X happen by March?"

Price history:
- 1 week ago: 0.45
- 3 days ago: 0.43
- Yesterday: 0.42
- Today: 0.60 (+43% in 24h!)

Validation:
✅ Check news: No major developments
✅ Check volume: Spike from $5k to $25k (suspicious?)
✅ Check fundamentals: Nothing changed
✅ Conclusion: Likely overreaction

Signal quality: STRONG mean reversion candidate
```

#### Step 2: Plan Your Entry

```
Current price: 0.60
Target (reversion): 0.45-0.50
Stop loss (wrong): 0.65
Risk/reward: (0.60-0.65) / (0.60-0.47.5) = 5:12.5 = 2.5:1

Entry strategy:
Option A: Sell immediately at 0.60
- Pro: Don't miss the opportunity
- Con: Price could go to 0.65 first

Option B: Wait for confirmation
- Entry: Sell at 0.62 if it spikes more
- Pro: Better price
- Con: Might miss the trade

Option C: Scale in
- Sell 1/3 at 0.60
- Sell 1/3 at 0.62 (if reaches)
- Sell 1/3 at 0.64 (if reaches)
- Pro: Better average, less timing risk
- Con: More complex
```

#### Step 3: Size Your Position

```
Account size: $10,000
Risk tolerance: 2% per trade = $200

Calculation:
Entry: 0.60 (selling)
Stop: 0.65
Risk per share: 0.05
Position size: $200 / 0.05 = $4,000
Shares: $4,000 / 0.60 = 6,667 shares

Collateral needed (for selling):
6,667 shares × $1.00 = $6,667

Check: Can you afford this?
Available: $10,000
Needed: $6,667
Remaining: $3,333 ✅

Position: 66.7% of capital (aggressive)
Consider: Reduce to 3,000-4,000 shares (30-40%)
```

#### Step 4: Execute the Trade

```
Selling (shorting) steps:
1. Verify you can sell on platform
2. Check collateral requirements
3. Place sell order (SELL 4,000 shares @ 0.60)
4. Confirm order filled
5. Set stop loss (BUY 4,000 @ 0.65)
6. Set target (BUY 4,000 @ 0.47-0.50)

Alternative: If can't sell
- Wait for reversion, then BUY low
- Miss the downside, but catch the upside
- Lower risk, lower reward
```

#### Step 5: Manage the Position

```
Scenario A: Price drops immediately to 0.50
- Your entry: 0.60
- Current: 0.50
- Profit: 0.10 per share = $400
- Action:
  - Take 50% profit: Buy 2,000 @ 0.50
  - Let 2,000 ride to 0.45 target
  - Move stop to breakeven (0.60)

Scenario B: Price rises to 0.64
- Your entry: 0.60
- Current: 0.64
- Loss: 0.04 per share = -$160
- Stop loss: 0.65 (close)
- Action:
  - Monitor closely
  - If hits 0.65: Exit (loss $200)
  - If reverses: Stay in trade

Scenario C: Price oscillates 0.58-0.62
- Action:
  - Hold position
  - Monitor for trend development
  - Patience (mean reversion takes time)
  - Re-evaluate after 2-3 days
```

#### Step 6: Exit Strategy

```
Exit scenarios:

✅ Target hit (0.47):
- Close position: Buy 4,000 @ 0.47
- Profit: 0.13 × 4,000 = $520
- ROI: 13% on $4,000 position

✅ Partial target (0.52):
- Close 50%: Buy 2,000 @ 0.52
- Profit on half: 0.08 × 2,000 = $160
- Let rest run to 0.47

❌ Stop loss hit (0.65):
- Close position: Buy 4,000 @ 0.65
- Loss: 0.05 × 4,000 = -$200
- ROI: -5%
- Accept loss, move on

⏱️ Time stop (3 days):
- If no movement after 3 days
- Close position at market
- Free up capital for other opportunities
```

### Strategy 2: Momentum Trading

Momentum trading assumes prices that are rising will continue to rise (and vice versa).

#### When to Use Momentum

✅ **Good signals:**
- 1-day change aligns with 1-week trend (both up or down)
- Strong, sustained directional movement
- Increasing volume on each leg up/down
- Breaking through resistance/support levels
- Catalyst or news driving the move

❌ **Avoid when:**
- Price at extremes (0.95 or 0.05)
- Volume declining as price moves
- Conflicting timeframe signals
- Near resolution (less room to move)
- Overbought/oversold indicators

#### Step 1: Validate the Trend

```
Example Market:
"Will Y happen by June?"

Price history:
- 1 month ago: 0.30
- 2 weeks ago: 0.38 (+27%)
- 1 week ago: 0.45 (+44% from start)
- Today: 0.60 (+100% from start, +33% this week)

Volume history:
- 1 month: $10k/day
- 2 weeks: $25k/day
- 1 week: $50k/day
- Today: $100k/day

Validation:
✅ Consistent uptrend across timeframes
✅ Accelerating price movement
✅ Increasing volume (healthy)
✅ News catalyst (major development)
✅ Conclusion: Strong momentum

Signal quality: STRONG momentum candidate
```

#### Step 2: Plan Your Entry

```
Current price: 0.60
Trend direction: UP
Next resistance: 0.70
Next support: 0.55
Target: 0.75-0.80
Stop loss: 0.55 (below support)

Entry strategy:
Option A: Buy now (aggressive)
- Entry: 0.60
- Pro: Don't miss the move
- Con: Could pullback first

Option B: Buy on pullback (preferred)
- Wait for: 0.57-0.58 (pullback to support)
- Entry: Buy if bounces from 0.57
- Pro: Better price, confirmed support
- Con: Might not pullback

Option C: Buy on breakout
- Wait for: 0.65 (breaks resistance)
- Entry: Buy at 0.65-0.66
- Pro: Confirmed continuation
- Con: Higher entry price

Recommendation: Option B (buy pullback)
```

#### Step 3: Position Sizing

```
Account: $10,000
Risk: 2% = $200 per trade

Entry: 0.57 (buying on pullback)
Stop: 0.55 (support level)
Risk per share: 0.02
Position size: $200 / 0.02 = $10,000 worth

Shares: $10,000 / 0.57 = 17,544 shares
Cost: 17,544 × 0.57 = $10,000

Problem: This is 100% of capital!

Solution: Reduce position
Shares: 5,000 (more reasonable)
Cost: 5,000 × 0.57 = $2,850
Risk: 5,000 × 0.02 = $100 (1% of account) ✅

This allows:
- Multiple positions
- Buffer for adverse moves
- Capital for rebalancing
```

#### Step 4: Execute the Trade

```
Buying steps:
1. Wait for pullback to 0.57-0.58
2. Place limit buy order: BUY 5,000 @ 0.58
3. Once filled, immediately set stop: SELL 5,000 @ 0.55
4. Set target: SELL 5,000 @ 0.75-0.80
5. Monitor for trend continuation

Alternative: If no pullback
- Market keeps rising: 0.60 → 0.63 → 0.65
- Decision: Enter at higher price or skip
- Risk: Chasing (buying at top)
- Recommendation: Let this one go, wait for next
```

#### Step 5: Trail Your Stop

```
Initial:
- Entry: 0.57
- Stop: 0.55
- Risk: $100

After move to 0.65:
- Current: 0.65
- New stop: 0.60 (trailing)
- Risk: $0 (now profitable)
- Profit protected: 0.03 × 5,000 = $150

After move to 0.72:
- Current: 0.72
- New stop: 0.68 (trailing)
- Profit protected: 0.11 × 5,000 = $550

Trailing strategy:
- Move stop under each new support
- Lock in profits as trend continues
- Let winners run
- Exit only if trend breaks
```

#### Step 6: Exit Strategy

```
Exit scenarios:

✅ Target hit (0.78):
- Sell 5,000 @ 0.78
- Profit: 0.21 × 5,000 = $1,050
- ROI: 37% on $2,850 position
- Time: 1-2 weeks

✅ Partial profits:
- At 0.70: Sell 2,500 (half) = $325 profit
- Let 2,500 run to 0.80+
- Trail stop on remaining half
- Lock in some gains, let rest run

❌ Stop hit (0.55):
- Sell 5,000 @ 0.55 (initial stop)
- Loss: 0.02 × 5,000 = -$100
- ROI: -3.5%
- Quick exit, accept small loss

❌ Trailing stop hit (0.68):
- After reaching 0.72
- Sell 5,000 @ 0.68
- Profit: 0.11 × 5,000 = $550
- ROI: 19%
- Trend is breaking, exit with profit
```

## Real-World Example Analysis

### Example 1: Strong Mean Reversion Signal

```
Market: "Will Fed cut rates in March?"
Current Price: 0.75
24h Change: +25% (from 0.60)
1 Week Change: +5% (from 0.71)

Analysis:
- Short-term spike: 0.60 → 0.75 in 24h
- Longer-term: Only slightly up from 0.71
- Pattern: Spike from lower level
- Interpretation: Overreaction to recent news

Volume:
- 24h volume: $150,000 (high)
- Normal volume: $30,000/day
- Spike: 5x normal (emotional buying?)

Strategy: MEAN REVERSION
Action: SELL at 0.75
Target: 0.68-0.70 (reversion to week average)
Stop: 0.78 (if genuinely new information)
Position: Moderate
Timeframe: 1-3 days

Risk/Reward:
- Risk: 0.03 (0.75 → 0.78)
- Reward: 0.06 (0.75 → 0.69)
- Ratio: 1:2 (acceptable)
```

### Example 2: Strong Momentum Signal

```
Market: "Will inflation exceed 3% in Q1?"
Current Price: 0.65
24h Change: +15% (from 0.57)
1 Week Change: +30% (from 0.50)

Analysis:
- Consistent uptrend: 0.50 → 0.57 → 0.65
- Accelerating: Getting steeper
- Pattern: Sustained directional move
- Interpretation: Strong momentum

Volume:
- 24h volume: $80,000
- Week ago: $20,000/day
- Growing: 4x increase over week
- Healthy volume increase

Strategy: MOMENTUM
Action: BUY on pullback to 0.63
Target: 0.75-0.80 (continuation)
Stop: 0.59 (below recent support at 0.60)
Position: Moderate to Aggressive
Timeframe: 1-2 weeks

Risk/Reward:
- Risk: 0.04 (0.63 → 0.59)
- Reward: 0.14 (0.63 → 0.77)
- Ratio: 1:3.5 (excellent)
```

### Example 3: Uncertain Signal

```
Market: "Will stock market hit new high by end of month?"
Current Price: 0.55
24h Change: +22% (from 0.45)
1 Week Change: -8% (from 0.60)

Analysis:
- Conflicting signals
- Down trend: 0.60 → 0.45 over week
- Sudden spike: 0.45 → 0.55 in day
- Pattern: Unclear (reversal or dead cat bounce?)

Questions:
- Is this trend reversal (new uptrend starting)?
- Or is this bounce within downtrend?
- Not enough information to decide

Strategy: WAIT or VERY SMALL POSITION
Options:
A) Wait for clearer signal (recommended)
B) Small long if breaks 0.60 (trend reversal confirmed)
C) Small short if fails at 0.58 (downtrend continues)
Position: Very small (25% normal size)
```

## Risk Management

### Position Sizing Formula

```
Kelly Criterion (simplified):
Position % = Edge / Odds

Example:
Win rate: 60%
Avg win: +10%
Avg loss: -5%
Edge: (0.60 × 0.10) - (0.40 × 0.05) = 0.04 (4%)
Odds: 0.10 / 0.05 = 2:1
Kelly: 0.04 / 2 = 0.02 = 2% position size

Half-Kelly (more conservative): 1% position size
```

### Stop Loss Placement

```
For mean reversion:
- Set stop above entry + 1.5x typical volatility
- Example: Entry 0.60, volatility 0.03
- Stop: 0.60 + (0.03 × 1.5) = 0.645

For momentum:
- Set stop below recent support
- Or: Entry - 2x ATR (Average True Range)
- Example: Entry 0.57, ATR 0.025
- Stop: 0.57 - (0.025 × 2) = 0.52
```

### Portfolio Risk

```
Total risk limits:
- Max per trade: 2% of account
- Max open risk: 6% of account
- Max positions: 3-5 concurrent

Example with $10,000:
Trade 1: Risk $150 (1.5%)
Trade 2: Risk $200 (2.0%)
Trade 3: Risk $175 (1.75%)
Total risk: $525 (5.25%) ✅

Don't add Trade 4 if it would exceed 6% total risk
```

### Time-Based Stops

```
Mean reversion trades:
- Exit if no reversion within 3-5 days
- Thesis likely wrong if taking too long

Momentum trades:
- Let trend run as long as it's working
- But review every 7-10 days
- Exit if volume declining or weakening

General:
- Don't let any trade "rot"
- If nothing happening, close and move on
- Free up capital and mental space
```

## Advanced Analysis Techniques

### 1. Volume Profile Analysis

```
Healthy price movement:
Price up 20% with:
- Volume up 3-5x
- Consistent across the move
- No single large trade
= Suggests many participants, sustainable

Unhealthy price movement:
Price up 20% with:
- Volume same or lower
- One or two large trades
- Most action in short period
= Suggests manipulation or thin market
```

### 2. Time-of-Day Patterns

```
Morning (8am-12pm):
- Higher volume
- More news-driven
- Better for momentum

Afternoon (12pm-6pm):
- Moderate volume
- Often consolidation
- Good for entries on pullbacks

Evening/Night:
- Lower volume
- More volatile (thin)
- Prices less reliable
```

### 3. Multiple Timeframe Analysis

```
Check alignment across timeframes:

Bullish alignment (momentum buy):
- 1 month: Up
- 1 week: Up
- 1 day: Up
- 4 hour: Up
= Strong momentum signal

Bearish divergence (mean reversion):
- 1 month: Flat
- 1 week: Down slightly
- 1 day: Up sharply
- 4 hour: Spiking
= Overreaction, reversion candidate
```

### 4. Correlation Analysis

```
Check related markets:

Example: Election markets
- "Candidate A wins presidency": Up 15%
- "Party A wins Senate": Up 5%
- "Party A wins House": Flat

Analysis:
- Presidency move seems overdone
- Congressional markets haven't moved much
- Suggests presidency price may be overreaction
- Strategy: Mean reversion on presidency market
```

## Performance Tracking

### Track These Metrics

```
For each trade:
- Entry price and date
- Exit price and date
- Strategy used (momentum or reversion)
- Signal quality (strong/moderate/weak)
- Position size
- Profit/loss
- Hold time

Monthly summary:
- Total trades
- Win rate
- Average win
- Average loss
- Profit factor = (Avg Win × Win Rate) / (Avg Loss × Loss Rate)
- Best strategy (momentum vs reversion)
```

### Evaluate Strategy Performance

```
After 20 trades, analyze:

Mean reversion trades:
- Win rate: 70%
- Avg win: +8%
- Avg loss: -4%
- Verdict: Working well

Momentum trades:
- Win rate: 45%
- Avg win: +18%
- Avg loss: -5%
- Verdict: Lower win rate but good R:R, keep it

Adjustments:
- Focus more on mean reversion (higher win rate)
- Be more selective on momentum (only strongest signals)
- Increase position size on mean reversion
- Reduce position size on momentum
```

## Common Mistakes to Avoid

### ❌ Mistake 1: Chasing Price

```
WRONG:
- Market moves from 0.40 → 0.60
- "I need to get in!"
- Buy at 0.60
- Immediately drops to 0.55
- Bought at top

CORRECT:
- Wait for pullback
- Buy at 0.57 (pullback to support)
- Or skip and wait for next opportunity
```

### ❌ Mistake 2: Fighting the Trend

```
WRONG:
- Market trending up: 0.40 → 0.50 → 0.60
- "It's too high, must drop"
- Short at 0.60
- Continues to 0.75
- Stopped out for big loss

CORRECT:
- "Trend is up, respect it"
- Look for opportunities to buy pullbacks
- Or skip and trade with the trend elsewhere
```

### ❌ Mistake 3: No Stop Loss

```
WRONG:
- Enter mean reversion sell at 0.70
- Expected drop to 0.60
- Rises to 0.75, 0.80, 0.85
- "It will come back down"
- Never does, huge loss

CORRECT:
- Set stop at 0.73 (before entry)
- If hit, accept small loss and exit
- Preserve capital for better opportunities
```

### ❌ Mistake 4: Overtrading

```
WRONG:
- Take every 15%+ move
- 10 trades per day
- 50% hit stops (overtrading)
- Death by 1000 cuts

CORRECT:
- Be selective, quality over quantity
- 2-3 good setups per week
- Better win rate with patience
```

### ❌ Mistake 5: Ignoring Fundamentals

```
WRONG:
- See 30% spike in "Will war start?" market
- "Mean reversion! Sell it!"
- War actually starts
- Market goes to 0.95
- Thesis was wrong (news-driven, not overreaction)

CORRECT:
- Check WHY price moved
- If fundamental news, respect it
- Mean reversion works on emotional moves, not fundamental
- Skip if major news
```

## When to Skip Trading

### Don't Trade When:

```
❌ Insufficient data:
- No 1-week comparison
- Limited price history
- First day of market

❌ Major news pending:
- Awaiting key announcement
- Resolution date approaching
- Known scheduled events

❌ Low confidence:
- Mixed signals
- Contradictory timeframes
- Unclear trend direction

❌ Poor liquidity:
- < $5,000 daily volume
- Wide spreads (>5%)
- Can't enter/exit easily

❌ Personal factors:
- Can't monitor position
- At risk limits
- Emotional/tilted
- Tired/distracted
```

## Conclusion

Rapid price movement trading offers opportunities through two main strategies:

**Mean Reversion:**
- Best for overreactions
- Short-term spikes diverging from longer trend
- Higher win rate, smaller moves
- Quick trades (1-3 days)

**Momentum:**
- Best for sustained trends
- Multi-timeframe alignment
- Lower win rate, larger moves
- Longer holds (1-2 weeks)

Success requires:
- Proper signal identification
- Appropriate strategy selection
- Strict stop losses
- Position sizing discipline
- Patience and selectivity

Start conservatively, track your results, and refine your approach over time.

## Related Examples

- [Closing Soon Markets](../find-closing-soon-markets/) - For predictable outcomes trading
- [Related Markets Arbitrage](../find-related-markets-arbitrage/) - For risk-free opportunities
- [Low Liquidity / High Volume](../find-low-liquidity-high-volume/) - For market making

## Additional Resources

- [Polymarket Documentation](https://docs.polymarket.com/)
- Technical Analysis Basics
- Main repository README for other strategies
