# Related Markets Arbitrage Example

## Overview

This example identifies **true arbitrage opportunities** in Polymarket events where the sum of probabilities across mutually exclusive markets deviates significantly from 1.0 (100%).

## What is Related Markets Arbitrage?

In prediction markets with mutually exclusive outcomes, the sum of all probabilities should theoretically equal 1.0, because exactly one outcome must occur. When this sum deviates from 1.0, arbitrage opportunities may exist.

### The Theory

For an event with mutually exclusive outcomes (e.g., "Who will win the election?"):
- Market A (Candidate 1): 0.40 (40%)
- Market B (Candidate 2): 0.35 (35%)
- Market C (Candidate 3): 0.20 (20%)
- **Theoretical Sum**: 1.00 (100%)

### When Arbitrage Exists

#### Underpriced Scenario (Sum < 1.0) - TRUE ARBITRAGE ✅
```
Market A: 0.40
Market B: 0.32
Market C: 0.18
Sum: 0.90 (< 1.0)

Strategy: BUY all markets
Cost: $0.90
Payout: $1.00 (one market will resolve to YES)
Profit: $0.10 (11.1% return) - GUARANTEED!
```

#### Overpriced Scenario (Sum > 1.0) - NOT TRUE ARBITRAGE ⚠️
```
Market A: 0.40
Market B: 0.38
Market C: 0.24
Sum: 1.02 (> 1.0)

Standard approach:
  Cost: $1.02
  Payout: $1.00
  Loss: $0.02 (not arbitrage!)

NegRisk approach (if available):
  Cost: $1.00 (fixed collateral)
  Payout: $1.00
  Profit: $0.00 (capital efficiency, not arbitrage)
```

## Key Difference: Arbitrage vs Capital Efficiency

### True Arbitrage (Sum < 1.0)
- **Guaranteed profit** regardless of outcome
- **Risk-free** (assuming proper execution)
- **Use standard markets** (lower cost than NegRisk's 1.0)
- Example: Pay $0.90, receive $1.00 → $0.10 profit

### Capital Efficiency (Sum > 1.0 with NegRisk)
- **No guaranteed profit**
- **Reduces capital requirements** only
- **Use NegRisk if available** (1.0 vs 1.02+ standard cost)
- Example: Pay $1.00, receive $1.00 → $0.00 profit, but saves capital

## How to Use This Example

### 1. Run the Example

```bash
cd examples/find-related-markets-arbitrage
go run main.go
```

### 2. Understanding the Output

The example searches for events with:
- ✅ Multiple mutually exclusive markets (2+)
- ✅ Probability sum deviating from 1.0 by at least 2%
- ✅ All markets are open and accepting orders
- ✅ Sufficient liquidity (minimum $1,000 per market)

For each opportunity, it displays:

#### Event Information
- Event title and description
- Number of markets
- Link to Polymarket event page

#### Arbitrage Analysis
- Probability sum (< 1.0 or > 1.0)
- Deviation percentage from 1.0
- Arbitrage type classification
- Expected return calculation

#### NegRisk Status
**NEW**: The tool now checks if events have NegRisk support and provides appropriate guidance:

For **underpriced** markets (sum < 1.0):
```
⚠️  NegRisk Warning:
   • This event HAS NegRisk support, but DO NOT use it for this arbitrage!
   • Why? With NegRisk you pay 1.0 collateral, but markets only cost 0.9520
   • Use STANDARD markets: Pay 0.9520, receive 1.0 = profit 0.0480
   • Using NegRisk: Pay 1.0, receive 1.0 = profit 0.0 (NO ARBITRAGE!)
```

For **overpriced** markets (sum > 1.0):
```
💡 NegRisk Alternative (Capital Efficiency, NOT Arbitrage):
   • This event HAS NegRisk support - can reduce capital requirements
   • Standard approach: Buy all markets for 1.0260 collateral
   • NegRisk approach: Only need 1.0 collateral (saves 0.0260)
   • However, this is NOT true arbitrage - no guaranteed profit
```

#### Market Breakdown
- Individual market prices
- Liquidity for each market
- Bid-ask spreads
- Order acceptance status

#### Execution Strategy
Detailed step-by-step instructions based on arbitrage type

### 3. Interpreting Opportunities

#### Excellent Arbitrage (Deviation > 5%)
```
Sum: 0.90 (10% deviation)
Expected Return: $0.10 per $0.90 invested = 11.1%
Risk: Very low (if markets truly mutually exclusive)
Action: STRONG BUY opportunity
```

#### Good Arbitrage (Deviation 3-5%)
```
Sum: 0.96 (4% deviation)
Expected Return: $0.04 per $0.96 invested = 4.2%
Risk: Low (check fees don't exceed profit)
Action: Consider after fee calculation
```

#### Marginal Arbitrage (Deviation 2-3%)
```
Sum: 0.98 (2% deviation)
Expected Return: $0.02 per $0.98 invested = 2.0%
Risk: May be eliminated by fees
Action: Calculate total costs carefully
```

## Step-by-Step Trading Guide

### Strategy 1: Underpriced Arbitrage (Sum < 1.0) - TRUE ARBITRAGE

**Example**: Event with 3 markets summing to 0.90

#### Step 1: Verify Market Conditions ✅

Critical checks before executing:
```
✅ Markets are truly mutually exclusive?
   - Only ONE outcome can occur
   - Read resolution criteria carefully
   - Check for "none of above" possibility

✅ All markets are active?
   - All accepting orders
   - Not closed or halted
   - Sufficient liquidity

✅ Calculate position size
   - Available capital
   - Liquidity constraints
   - Risk tolerance
```

#### Step 2: Calculate Exact Costs and Returns

```
Position size: $1,000
Markets sum: 0.90

Buying Strategy:
- Allocate: $1,000 / 0.90 = $1,111 worth of positions
- Market A (40%): Buy $444 worth ($444 / 0.40 = 1,111 shares)
- Market B (32%): Buy $356 worth ($356 / 0.32 = 1,111 shares)
- Market C (18%): Buy $200 worth ($200 / 0.18 = 1,111 shares)

At Resolution:
- One market pays: 1,111 shares × $1.00 = $1,111
- Other markets pay: $0
- Total payout: $1,111

Profit Calculation:
- Investment: $1,000
- Payout: $1,111
- Gross profit: $111 (11.1%)
- Fees (2% avg): $20
- Net profit: $91 (9.1%)
```

#### Step 3: Execute the Arbitrage

**IMPORTANT**: Use **STANDARD markets**, NOT NegRisk!

1. **Place Orders Simultaneously**
   ```
   Why simultaneous?
   - Prevents prices from moving against you
   - Maintains the arbitrage spread
   - Reduces execution risk
   ```

2. **Equal Dollar Amounts**
   ```
   Market A: $444 (40% × $1,111)
   Market B: $356 (32% × $1,111)
   Market C: $200 (18% × $1,111)

   This ensures equal SHARES (1,111) in each market
   ```

3. **Use Market Orders for Speed** (if liquidity allows)
   ```
   Pros: Immediate execution, captures arbitrage
   Cons: Slippage may reduce profit

   Alternative: Limit orders at current prices
   Pros: Price protection
   Cons: May not fill if market moves
   ```

4. **Monitor Execution**
   ```
   Check: All orders filled?
   Check: Final average prices match expectations?
   Check: Total cost close to predicted?
   ```

#### Step 4: Hold Until Resolution

```
Do:
✅ Monitor market for any resolution issues
✅ Track any news about the event
✅ Verify resolution criteria being met

Don't:
❌ Try to exit early (loses arbitrage)
❌ Sell any positions (breaks the hedge)
❌ Add to positions (changes risk profile)
```

#### Step 5: Collect Profits

```
At Resolution:
- Winning market: Receives $1.00 per share
- Losing markets: Receive $0.00 per share
- Total: 1,111 × $1.00 = $1,111

Profit Realized:
- Payout: $1,111
- Cost: $1,000
- Fees: $20
- Net Profit: $91 (9.1% return)
```

### Strategy 2: Overpriced Markets (Sum > 1.0) - NOT ARBITRAGE

**Example**: Event with 3 markets summing to 1.03

#### Understanding the Situation

```
This is NOT true arbitrage because:
- Standard cost: $1.03 > Payout: $1.00 (you lose money!)
- Selling strategy: Complex, requires margin, not risk-free

Options:
1. Skip opportunity (recommended for most)
2. Use NegRisk for capital efficiency (advanced)
3. Sell strategy (very advanced, requires margin)
```

#### Option A: Skip (Recommended)

```
Why skip?
- No guaranteed profit
- Requires additional capital
- More complex execution
- Higher risk

When to skip:
- You're seeking arbitrage only
- You don't have NegRisk access
- You don't understand selling mechanics
```

#### Option B: NegRisk Capital Efficiency (If Available)

**Only if NegRisk is enabled for the event!**

```
Use Case: You want exposure with less capital

Standard Approach:
- Buy all markets: $1.03 collateral needed
- Payout: $1.00
- Loss: $0.03
- Why? Not actually doing this

NegRisk Approach:
- Collateral: $1.00 (fixed)
- Payout: $1.00
- Profit: $0.00
- Advantage: 2.9% less capital locked

Benefit: Free up capital for other trades
```

#### Option C: Selling Strategy (Advanced)

**WARNING**: This is complex and risky!

```
Theory:
- Sell all markets: Collect $1.03
- Payout at resolution: $1.00
- Profit: $0.03

Reality - Why It's Hard:
1. Requires margin/collateral for each market
2. Collateral needed: $1.00 per market × 3 = $3.00
3. Not capital efficient
4. Risk of assignment/early resolution
5. Margin requirements may change

Recommendation: Skip unless experienced
```

## Real-World Examples from Output

### Example 1: Excellent Arbitrage
```
Event: "NATO/EU troops fighting in Ukraine by...?"
Markets: 2
Probability Sum: 0.1230 (87.7% deviation!)

Analysis:
- Market 1: 0.023 ($16,012 liquidity)
- Market 2: 0.100 ($5,058 liquidity)
- Sum: 0.123 (far below 1.0)

Opportunity Assessment:
✅ HUGE deviation (87.7%)
⚠️  BUT: Very low sum suggests market uncertainty
⚠️  Likely both outcomes considered very unlikely
⚠️  Check for "other outcomes" not listed

Risk Check:
- Are these truly mutually exclusive?
- Is there a "none of above" scenario?
- Could event resolve differently?

Verdict:
- High percentage return IF executed correctly
- MUST verify markets are truly exclusive
- May be uncertainty priced in (neither happens)
```

### Example 2: Good Arbitrage with NegRisk Warning
```
Event: "How many people will Trump deport in 2025?"
Markets: 9 mutually exclusive ranges
Probability Sum: 0.9520 (4.8% deviation)
NegRisk Enabled: YES

Analysis:
- 9 ranges covering all possibilities
- Sum: 0.952 (underpriced by 4.8%)
- Liquidity: $64,647 total
- Min market liquidity: $3,718

Strategy:
✅ Buy all 9 markets using STANDARD interface
❌ DO NOT use NegRisk (would cost 1.0, eliminating profit)

Execution:
- Per $1,000 investment:
  - Cost: $952
  - Payout: $1,000
  - Gross profit: $48 (5%)
  - Fees: ~$19 (2%)
  - Net profit: ~$29 (3.0%)

Why Standard, Not NegRisk:
- Standard: Pay $0.952, get $1.00 = +$0.048
- NegRisk: Pay $1.000, get $1.00 = $0.000
- Difference: $0.048 per dollar = 4.8% profit lost!
```

### Example 3: Marginal Arbitrage
```
Event: "How much spending will DOGE cut in 2025?"
Markets: 6 mutually exclusive ranges
Probability Sum: 0.9770 (2.3% deviation)
NegRisk Enabled: YES

Analysis:
- 6 ranges covering possibilities
- Sum: 0.977 (underpriced by 2.3%)
- Total liquidity: $66,717

Profit Calculation:
- Per $1,000: Profit $23 (2.3%)
- Fees (~2%): $20
- Net: $3 (0.3%)

Assessment:
⚠️  Marginal opportunity
⚠️  Fees eat most profit
⚠️  Slippage could eliminate gains
⚠️  May not be worth the capital lock-up

Decision Framework:
- Skip if fees > 2%
- Skip if you need capital soon
- Consider only for large size with low slippage
```

## Risk Management

### Critical Verification Steps

#### 1. Market Exclusivity Verification
```
Questions to ask:
✅ Can only ONE outcome occur?
✅ Do markets cover ALL possibilities?
✅ Is there a "none of above" option?
✅ Are resolution criteria clear?
✅ Have similar markets resolved as expected?

Example Red Flags:
❌ "Market resolves if X happens" - what if X partially happens?
❌ Vague resolution criteria
❌ Historical disputes on similar markets
❌ Ambiguous outcome definitions
```

#### 2. Liquidity Analysis
```
Check:
✅ Total liquidity > 10x your position
✅ Each market liquidity > 5x your allocation
✅ Recent trading volume (market is active)
✅ Spread is reasonable (< 2%)

Example:
Your position: $1,000
Market A: $10,000 liquidity ✅
Market B: $3,000 liquidity ⚠️ (only 3x)
Market C: $500 liquidity ❌ (only 0.5x)

Verdict: Reduce position size or skip
```

#### 3. Fee Calculation
```
Typical fees:
- Maker fee: 2% (if providing liquidity)
- Taker fee: 2-4% (if taking liquidity)
- NegRisk fee: 0-1% (if applicable)
- Gas fees: Negligible on Polygon

Example:
Profit: 2.3%
Taker fees: 2%
Net: 0.3%
Verdict: Probably not worth it
```

#### 4. Slippage Estimation
```
Formula: Position size / Market liquidity

Example:
Market A liquidity: $10,000
Your buy: $444
Slippage estimate: 444/10,000 = 4.4% of depth
Expected slippage: ~0.2-0.5%

Total slippage across 3 markets: ~1-2%
Impact on 2.3% profit: May eliminate it entirely
```

### Position Sizing Guidelines

```
Conservative:
- Use 5-10% of available capital per opportunity
- Limit to 1% of each market's liquidity
- Keep 50%+ capital for other opportunities

Moderate:
- Use 10-20% of capital
- Limit to 5% of market liquidity
- Diversify across 3-5 opportunities

Aggressive (not recommended):
- Using >25% of capital
- Taking >10% of market liquidity
- All-in on single opportunity
```

## Common Mistakes to Avoid

### ❌ Mistake 1: Using NegRisk for Underpriced Arbitrage
```
Event: Sum = 0.95

WRONG:
  Choose NegRisk because it's available
  Cost: $1.00
  Payout: $1.00
  Profit: $0.00 ❌

CORRECT:
  Use standard markets
  Cost: $0.95
  Payout: $1.00
  Profit: $0.05 ✅
```

### ❌ Mistake 2: Ignoring "None of Above" Scenarios
```
Event: "Winner of Best Picture Oscar"
Markets:
- Movie A: 0.40
- Movie B: 0.30
- Movie C: 0.25
Sum: 0.95 (looks like arbitrage!)

Reality:
- Movie D could win (not listed): 0.05 probability
- If Movie D wins: ALL your positions lose
- Not true arbitrage after all

Lesson: Verify ALL possible outcomes are covered
```

### ❌ Mistake 3: Sequential Execution
```
WRONG:
  1. Buy Market A ($400)
  2. Market moves, now A = 0.45
  3. Buy Market B ($350)
  4. Market moves, now B = 0.37
  5. Buy Market C ($220)
  Sum now: 1.02 (no longer arbitrage!)

CORRECT:
  Place all orders simultaneously
  Use limit orders at current prices
  Execute within seconds
```

### ❌ Mistake 4: Ignoring Fees
```
Opportunity: 2.5% deviation
Gross profit: $25 per $1,000
Fees: 2% × $1,000 = $20
Net profit: $5 (0.5%)

Mistake: Thinking "2.5% is good!"
Reality: After fees, barely profitable
Risk: Slippage could eliminate remaining profit
```

### ❌ Mistake 5: Misunderstanding Overpriced Markets
```
Event: Sum = 1.05

Wrong thinking:
"I can sell all markets, collect $1.05, pay $1.00, profit $0.05!"

Reality:
- Selling requires $1.00 collateral PER market
- Total collateral: $1.00 × 3 markets = $3.00
- You collect: $1.05
- You pay back: $1.00
- Profit: $0.05
- But you locked up: $3.00 (ROI = 1.67%, very poor)

Better: Skip or use NegRisk for capital efficiency only
```

## Advanced Techniques

### 1. Event Relationship Analysis
```
Look for related events:
Event A: "Trump wins 2024?" - 0.60
Event B: "Democrat wins 2024?" - 0.42
Sum: 1.02

But wait - these should be exact opposites!
Potential arbitrage:
- Buy "Democrat wins": $0.42
- Sell "Trump wins": Collect $0.60
- Profit: $0.18?

Reality check:
- Are there third party candidates?
- Could "neither" happen?
- Are resolution criteria identical?
```

### 2. Time-Based Arbitrage
```
Monitor for arbitrage appearing/disappearing:
- Morning: Sum = 1.02 (no arbitrage)
- News breaks
- Afternoon: Sum = 0.94 (arbitrage!)
- Execute quickly
- Evening: Sum back to 1.00 (arbitrage closed)

Key: Speed of detection and execution
```

### 3. Cross-Event Hedging
```
Event 1: "Team A wins championship" - 0.30
Event 2: "Team A makes finals" - 0.25

Logic problem: Can't win championship without making finals!
Arbitrage: If Team A wins championship (0.30), then definitely made finals (should be at least 0.30)

Strategy:
- Buy "makes finals" at 0.25
- If championship wins: Finals definitely YES
- If championship loses but makes finals: Finals YES
- Profit opportunity from logical inconsistency
```

## Monitoring and Adjustment

### During Holding Period

```
Daily checks:
✅ Monitor for resolution news
✅ Watch market prices (should converge to 1.0)
✅ Check for any disputes/questions
✅ Verify all positions still held correctly

Weekly checks:
✅ Review resolution timeline
✅ Check for any rule changes
✅ Monitor platform announcements
✅ Calculate unrealized P&L
```

### Warning Signs

```
🚨 Red Flags to Watch:
- Prices diverging further from 1.0
- Questions about resolution criteria
- Unusual trading volume
- Markets being halted/reopened
- Ambiguous outcomes emerging
- Similar markets resolving unexpectedly

Action: Research immediately, consider exiting
```

## Real-World Performance Expectations

### Realistic Returns

```
Excellent opportunity (5%+ deviation):
- Gross profit: 5-10%
- Fees: 2%
- Slippage: 0.5-1%
- Net profit: 2.5-7.5%

Good opportunity (3-5% deviation):
- Gross profit: 3-5%
- Fees: 2%
- Slippage: 0.5-1%
- Net profit: 0.5-2.5%

Marginal opportunity (2-3% deviation):
- Gross profit: 2-3%
- Fees: 2%
- Slippage: 0.5-1%
- Net profit: -0.5 to +0.5%
- Verdict: Usually not worth it
```

### Frequency of Opportunities

```
Expected frequency:
- Excellent (5%+): Rare, maybe 1-2 per month
- Good (3-5%): Occasional, 2-5 per month
- Marginal (2-3%): Common, 10+ per month

Reality:
- Most "opportunities" aren't true arbitrage
- Markets are usually efficient
- Best opportunities are taken quickly
- You need to act fast when found
```

## Conclusion

Related markets arbitrage provides TRUE risk-free profit opportunities, but only when:

✅ Probability sum < 1.0 (underpriced markets)
✅ Markets are verified mutually exclusive
✅ Liquidity supports your position size
✅ Fees + slippage don't eliminate profit
✅ Execution is simultaneous and fast
✅ You use STANDARD markets (not NegRisk for underpriced!)

This tool helps you identify candidates, but successful execution requires:
- Careful verification
- Quick execution
- Proper position sizing
- Understanding of risks
- Recognition of when to skip

Start small, verify thoroughly, and scale up as you gain experience.

## Related Examples

- [NegRisk Opportunities](../find-negrisk-opportunities/) - For capital efficiency strategies (sum > 1.0)
- [Low Liquidity / High Volume](../find-low-liquidity-high-volume/) - For market making
- [Closing Soon Markets](../find-closing-soon-markets/) - For information arbitrage

## Additional Resources

- [Polymarket Documentation](https://docs.polymarket.com/)
- [Market Resolution Criteria](https://docs.polymarket.com/#market-resolution)
- Main repository README for other strategies
