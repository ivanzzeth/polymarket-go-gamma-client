# Markets Closing Soon Trading Example

## Overview

This example identifies **information arbitrage opportunities** by finding markets approaching resolution where outcomes may be predictable, yet prices haven't fully adjusted.

## The Opportunity

### Why Trade Markets Closing Soon?

```
Market closes in 6 hours: "Will it rain in NYC today?"
Current price: 0.30 (30% probability)
Your information: It's currently raining heavily in NYC

Opportunity:
- Outcome is already determined/highly predictable
- Price hasn't adjusted to reality
- Information asymmetry = profit potential
```

**Key insight**: As resolution approaches, markets should converge to 0.00 or 1.00 based on observable reality. If they haven't, there may be mispricing.

### Types of Closing-Soon Opportunities

#### 1. Already Determined (Information Lag)
```
Market: "Will Team A win the championship?"
Closes: In 2 hours
Current Price: 0.65

Reality: Team A won 3 hours ago!
Opportunity: Buy at 0.65, resolves to 1.00
Profit: 35% (0.35 / 0.65 = 54% ROI)
Reason: Price hasn't updated to reflect known outcome
```

#### 2. Highly Predictable (Low Remaining Uncertainty)
```
Market: "Will unemployment rate be above 4.5% in January?"
Closes: In 12 hours
Current Price: 0.15 (15%)

Reality: Official data released showing 5.2%!
Opportunity: Buy at 0.15, resolves to 1.00
Profit: 85% (0.85 / 0.15 = 567% ROI)
Reason: Official data available, outcome certain
```

#### 3. Mispriced Near Resolution
```
Market: "Will bill pass by end of month?"
Closes: Tonight (12 hours)
Current Price: 0.82

Reality: Bill failed yesterday, but price hasn't dropped
Opportunity: Sell at 0.82, resolves to 0.00
Profit: 82% (assuming you can sell/short)
Reason: Market slow to react
```

## How to Use This Example

### 1. Run the Example

```bash
cd examples/find-closing-soon-markets
go run main.go
```

### 2. Configuration

Edit `main.go` to adjust:

```go
hoursUntilClose := 48.0    // Find markets closing within 48 hours
minVolume := 1000.0        // Minimum trading volume
targetCount := 5           // Number of opportunities
```

### 3. Understanding the Output

The example finds markets with:
- ✅ Closing within specified hours (default 48h)
- ✅ Still open and active
- ✅ Sufficient volume (>$1,000)
- ✅ Clear end date

For each opportunity, it displays:

#### Timing Information
```
End Date:             2025-01-20 18:00:00 EST
Time Until Close:     6 hours 23 minutes ⚠️ CLOSING SOON!

Categories by urgency:
< 12 hours: ⚠️ CLOSING SOON! (Act fast)
12-24 hours: ⏳ (Time to research)
24-48 hours: (Standard timeframe)
```

#### Current Pricing
```
Current Price:        0.7500 (Likely YES)

Price interpretations:
> 0.95: Very likely YES
0.80-0.95: Likely YES
0.60-0.80: Somewhat likely YES
0.40-0.60: Uncertain
0.20-0.40: Somewhat likely NO
0.05-0.20: Likely NO
< 0.05: Very likely NO
```

#### Resolution Information
```
Automatically Resolved: true ✅ (Lower resolution risk)

Auto resolution: Usually based on objective data source
Manual resolution: Requires human judgment (higher risk)
```

#### Mispricing Assessment

The tool analyzes potential mispricing:

**Fairly Priced:**
```
Mispricing Assessment: fairly_priced

The market shows a high probability (0.96) with 5 hours until close.
Market appears to have converged to likely outcome.

Suggested Action:
• Only trade if you have strong contradicting evidence
• Be very careful - the market may be correct
• Consider the cost of being wrong
```

**Potentially Overpriced:**
```
Mispricing Assessment: potentially_overpriced

The market shows a high probability (0.92) with 36 hours until close.
This seems early for such confidence.

Suggested Action:
• Research if the high probability is justified
• If outcome is less certain, consider SELLING/SHORTING
• Verify through multiple independent sources
```

**Potentially Underpriced:**
```
Mispricing Assessment: potentially_underpriced

The market shows a low probability (0.08) with 30 hours until close.
This seems early for such pessimism.

Suggested Action:
• Research if the low probability is justified
• If outcome is more certain, consider BUYING
• Verify through multiple independent sources
```

**Uncertain:**
```
Mispricing Assessment: uncertain

The market outcome appears genuinely uncertain (price: 0.52).
With 18 hours remaining, significant information may still emerge.

Suggested Action:
• Research the event thoroughly
• Look for information asymmetry opportunities
• Consider if you have better information than the market
• Monitor for new developments
```

## Step-by-Step Trading Guides

### Strategy 1: Already Determined Outcomes

The most straightforward opportunity - outcome has occurred but price hasn't adjusted.

#### Step 1: Verify the Outcome

```
Critical checklist:

✅ Official sources (required):
   - Government websites
   - Official announcements
   - Verified news sources (multiple)
   - Primary data sources

✅ Timing verification:
   - When did event occur?
   - When does market close?
   - Is there time to trade?

✅ Resolution criteria:
   - Read market description carefully
   - Check resolution source specified
   - Verify your information matches criteria
   - Look for edge cases or exceptions

❌ DO NOT trade based on:
   - Social media rumors
   - Unofficial sources
   - Partial information
   - Assumptions
```

**Example verification:**

```
Market: "Will unemployment rate be above 4% in January?"
Closes: Tonight 11:59 PM
Current Price: 0.25

Step 1: Find official source
✅ Visit: fred.stlouisfed.org or bls.gov
✅ Check: Official January unemployment data
✅ Find: 4.8% (above 4%)
✅ Verify: Data is official and final

Step 2: Check resolution criteria
✅ Read market description
✅ Confirms: Uses official BLS data
✅ Confirms: January data (correct month)
✅ Confirms: Threshold is 4% (not 4.5%)

Step 3: Verify timing
✅ Data released: 9am today
✅ Market closes: Tonight 11:59pm
✅ Time to trade: Yes (14 hours remaining)

Conclusion: STRONG BUY at 0.25 (should be 1.00)
```

#### Step 2: Calculate Opportunity Size

```
Market price: 0.25
Expected resolution: 1.00
Gross profit: 0.75 per share
ROI: 0.75 / 0.25 = 300%

Position sizing:
Account size: $10,000
Risk: What if you're wrong?
  - Lose: 100% of position
  - Probability: Very low (<5% if verified correctly)

Conservative: $500 (5% of account)
- Shares: $500 / 0.25 = 2,000
- Profit if correct: 2,000 × 0.75 = $1,500
- Loss if wrong: $500

Moderate: $1,000 (10% of account)
- Shares: $1,000 / 0.25 = 4,000
- Profit if correct: $3,000
- Loss if wrong: $1,000

Aggressive: $2,000 (20% of account)
- Only if VERY confident
- Triple-checked sources
- Zero doubt
```

#### Step 3: Execute Quickly

```
Time sensitivity is critical!

Execution steps:
1. Complete verification (10-15 minutes)
2. Size position (2 minutes)
3. Place order IMMEDIATELY
4. Monitor fill
5. Set alert for resolution

Why speed matters:
- Price may correct quickly as others notice
- Opportunity window might be minutes
- Other traders doing same analysis
- Algorithm detection

Order type:
- Limit order at current ask (0.26)
- Small premium acceptable (0.27-0.28)
- Avoid market orders (slippage)
- But don't miss trade over 1-2 ticks
```

#### Step 4: Monitor Until Resolution

```
After entry at 0.25:

Scenario A: Price corrects to 0.95
- Opportunity: Take partial profits?
- Consideration: 4 hours until close
- Decision: Hold (last 5% not worth risk of exit)

Scenario B: Price stays at 0.25-0.30
- Market seems "broken" or illiquid
- Decision: Hold, you have verified outcome
- Patience: Will resolve correctly

Scenario C: Price drops to 0.10
- Concern: Did you miss something?
- Action: RE-VERIFY immediately
  - Check sources again
  - Read resolution criteria again
  - Look for new information
  - If still confident: Hold
  - If doubt: Exit and reassess

Resolution (market closes):
- Wait for resolution (usually 0-48h)
- Receive 1.00 per share
- Profit: 0.75 per share = $1,500-6,000
```

### Strategy 2: Highly Predictable Outcomes

Outcome hasn't occurred yet, but is highly likely based on available information.

#### Step 1: Assess Probability

```
Market: "Will bill pass Congress by Friday?"
Closes: Friday 11:59 PM (48 hours)
Current Price: 0.45

Available information:
- Bill scheduled for vote Thursday
- Current vote count: 250 committed YES (need 218)
- No expected defections
- Procedural hurdles cleared

Your assessment:
- Probability of passage: 90-95%
- Current price: 0.45 (45%)
- Mispricing: Significant
- Confidence: High (but not certain)

Opportunity assessment:
- Expected value: 0.95 × 1.00 + 0.05 × 0.00 = 0.95
- Current price: 0.45
- Edge: 0.50 (50% underpriced)
- Risk: 5% chance of loss (bill fails)

Decision: STRONG BUY
```

#### Step 2: Research Thoroughly

```
Research checklist:

✅ Primary sources:
   - Official schedules
   - Vote counts (from official sources)
   - Procedural status
   - Amendment status

✅ Risk factors:
   - Could vote be postponed?
   - Any last-minute opposition?
   - Procedural tricks possible?
   - Time sufficient for passage?

✅ Resolution criteria:
   - What exactly triggers YES?
   - Just House or Senate too?
   - Simple majority or supermajority?
   - Does President signature matter?

✅ Historical context:
   - Similar bills in past?
   - Surprises in similar situations?
   - This Congress's patterns?

Warning signs:
❌ Conflicting vote counts
❌ Unclear procedural status
❌ Major amendments pending
❌ Time constraints
```

#### Step 3: Position with Proper Risk

```
Even with 90-95% confidence, there's 5-10% risk!

Position sizing for probabilistic outcomes:

Conservative (recommended):
- Risk: 2% of account
- If wrong, lose 2%
- Account: $10,000
- Risk: $200
- Position: $200 / 0.45 = $445 worth
- Shares: ~990

Rationale:
- 90% chance: Profit = 990 × 0.55 = $545
- 10% chance: Loss = -$445
- Expected value: (0.90 × 545) + (0.10 × -445) = $446

Aggressive:
- Position: $2,000 (20% of account)
- Only if extremely confident
- Still have 5-10% chance of $2k loss!
```

#### Step 4: Monitor Developments

```
Active monitoring required:

Every 2-4 hours until resolution:
✅ Check for postponements
✅ Monitor vote counts
✅ Watch for surprise amendments
✅ Track procedural progress
✅ Review price movements

Set news alerts for:
- Bill name
- Key representatives
- Related legislation
- Congress floor activity

If new information emerges:
- Reassess probability
- Adjust position if needed
- Be willing to exit if confidence drops
```

#### Step 5: Manage Exit

```
Exit scenarios:

✅ Bill passes (expected):
- Shares: 990
- Cost: $445 @ 0.45
- Payout: 990 @ 1.00 = $990
- Profit: $545 (122% ROI)

❌ Bill fails (unexpected):
- Shares: 990
- Cost: $445
- Payout: $0
- Loss: -$445 (100% loss on position)

⚠️ Vote postponed:
- Market may not resolve
- Check resolution criteria
- May need to hold longer
- Or exit at market price
```

### Strategy 3: Information Asymmetry

You have information or analysis the broader market hasn't incorporated.

#### Step 1: Identify Your Edge

```
Types of information edges:

1. Local knowledge:
   Market: "Will it snow in Denver tomorrow?"
   Edge: You live in Denver, see weather firsthand

2. Domain expertise:
   Market: "Will Fed raise rates?"
   Edge: You're an economist who analyzes Fed closely

3. Research-based:
   Market: "Will movie gross $100M opening weekend?"
   Edge: You've analyzed advance ticket sales data

4. Timing advantage:
   Market: "Will team win championship?"
   Edge: Game ended 10 minutes ago, you watched it

Important: Edge must be real and verifiable
```

**Example edge:**

```
Market: "Will Company X hit revenue target?"
Closes: In 24 hours
Current Price: 0.60

Your edge:
- You follow this company closely
- Earnings reported 2 hours ago
- Showed: $1.2B revenue (target was $1.0B)
- Market hasn't updated yet (low volume, evening)

Verification:
✅ Check official earnings report (investor relations)
✅ Verify numbers match resolution criteria
✅ Confirm no adjustments/restatements pending
✅ Check resolution source matches your source

Confidence: Very high (99%)
Action: Buy aggressively at 0.60
Target: 1.00 (should update overnight)
```

#### Step 2: Verify Information Quality

```
Information quality checklist:

✅ Primary source (not secondhand)
✅ Official or authoritative
✅ Recent/current
✅ Directly relevant to market criteria
✅ Independently verifiable
✅ Not ambiguous

Warning signs of bad information:
❌ "I heard from someone..."
❌ Unofficial sources only
❌ Can't verify independently
❌ Interpretation required
❌ Contradicted by other sources
❌ Too good to be true
```

#### Step 3: Act Before Edge Disappears

```
Information edges are temporary!

Typical edge lifetime:
- Breaking news: 1-30 minutes
- Data releases: 5-60 minutes
- Analysis-based: 1-24 hours
- Local knowledge: Variable

Execution speed:
1. Verify information (5-10 min) - CAN'T SKIP
2. Check resolution criteria (2 min) - CAN'T SKIP
3. Size position (1 min)
4. Place order (immediate)

Total time: 10-15 minutes max

After this, assume others will notice:
- Other traders
- Algorithms
- Market makers
- Opportunity may close
```

#### Step 4: Adjust for Resolution Risk

```
Even with perfect information, resolution can surprise:

Resolution risks:
- Interpretation differences
- Data source mismatches
- Unexpected criteria application
- Resolution delays/disputes
- Platform issues

Protection:
- Read resolution criteria 3 times
- Verify your source matches resolution source
- Check historical resolutions for this creator
- Reduce position size by 20-30% for resolution risk
- Be prepared for unexpected outcomes

Example:
Perfect information suggests 100% chance of YES
Your confidence: 100%
Resolution risk: 5-10%
Actual position: Size for 90-95% confidence
Reason: Account for resolution uncertainty
```

## Risk Management Framework

### Pre-Trade Checklist

```
Before entering ANY closing-soon trade:

☐ Outcome verification:
   ☐ Checked minimum 2 independent sources
   ☐ Sources are official/primary
   ☐ Information is current
   ☐ No conflicting information

☐ Resolution criteria:
   ☐ Read full market description
   ☐ Understand resolution source
   ☐ Timing requirements clear
   ☐ Edge cases considered

☐ Timing:
   ☐ Sufficient time to trade
   ☐ Market still accepting orders
   ☐ Time until resolution known
   ☐ No expected delays

☐ Position sizing:
   ☐ Calculated based on confidence
   ☐ Within risk limits
   ☐ Account for fees
   ☐ Can afford to be wrong

☐ Exit plan:
   ☐ Know when resolution occurs
   ☐ Set calendar alerts
   ☐ Planned response to scenarios
   ☐ Time to monitor trade
```

### Position Sizing by Confidence

```
Confidence Level | Position Size | Typical Scenario
-----------------|---------------|------------------
99%+ | 10-20% | Already determined, verified
90-95% | 5-10% | Highly predictable, good info
80-90% | 3-5% | Solid prediction, some uncertainty
70-80% | 1-2% | Moderate confidence
<70% | Skip | Too uncertain for closing-soon trade

Important: These are maximums! Start smaller.
```

### Red Flags - When NOT to Trade

```
❌ Automatic skip conditions:

1. Information uncertainty:
   - Only one source
   - Unofficial sources
   - Can't verify independently
   - Contradictory information

2. Resolution concerns:
   - Ambiguous criteria
   - Subjective judgment required
   - Historical disputes
   - Manual resolution with bias risk

3. Timing issues:
   - < 30 minutes to close (too rushed)
   - Closed to new orders
   - Resolution date uncertain
   - Platform issues

4. Market conditions:
   - No liquidity (can't enter/exit)
   - Very wide spread (>10%)
   - Recent resolution disputes
   - Platform trading issues

5. Personal factors:
   - Can't verify properly
   - Don't understand criteria
   - Too emotionally invested
   - At risk limits
```

## Real-World Example Analysis

### Example 1: Clear Mispricing

```
From Output:
Market: "Will unemployment be above X%?"
Current Price: 0.15 (15%)
Time Until Close: 18 hours

Research:
✅ Official BLS data released this morning
✅ Shows: 5.2% (above threshold)
✅ Data is final (no revisions expected)
✅ Resolution criteria: Uses official BLS data
✅ Source matches

Analysis:
- Outcome: Already determined
- Current price: 0.15 (WAY underpriced)
- Should be: 0.99-1.00
- Reason for mispricing: Low volume market, evening
- Opportunity quality: EXCELLENT

Position sizing:
- Confidence: 99%
- Position: 15% of account = $1,500
- Cost: $1,500 / 0.15 = $10,000 shares
- Expected profit: 10,000 × 0.85 = $8,500
- Risk: Misread criteria or source = -$1,500

Action: STRONG BUY
```

### Example 2: Predictable but Not Certain

```
From Output:
Market: "Will bill pass Congress by Friday?"
Current Price: 0.52 (uncertain)
Time Until Close: 36 hours

Research:
✅ Vote scheduled for Thursday 2pm
✅ Current whip count: 235 YES, 198 NO (need 218)
✅ No expected switches
⚠️ Vote could be postponed
⚠️ Last-minute amendments possible
✅ Resolution: Based on official Congressional record

Analysis:
- Outcome: Highly likely but not certain
- Current price: 0.52 (seems fair to underpriced)
- Your estimate: 85% passes
- Expected value: 0.85 vs price 0.52 = 33% edge
- Risk factors: Timing, amendments, surprises

Position sizing:
- Confidence: 85%
- Position: 5% of account = $500
- Cost: $500 / 0.52 = 960 shares
- Expected profit: 960 × 0.48 × 0.85 = $391
- Expected loss: 960 × 0.52 × 0.15 = -$75
- Net EV: +$316

Action: MODERATE BUY (with monitoring)
```

### Example 3: False Opportunity

```
From Output:
Market: "Will Team A win championship?"
Current Price: 0.85 (likely YES)
Time Until Close: 6 hours

Initial thought:
"Game is today at 5pm, team A is heavily favored!"

Research reveals:
❌ Game already happened (ended 1 hour ago)
❌ Team A LOST!
❌ But price is still 0.85?

Wait, something's wrong:
✅ Re-read criteria: "Win championship SERIES" (not this game)
✅ Series is best-of-7
✅ Team A won this game, leads series 3-2
✅ Price of 0.85 makes sense now

Lesson:
- Almost made terrible trade!
- Misread resolution criteria
- ALWAYS verify criteria carefully
- "Too good to be true" often is

Action: SKIP (misunderstood market)
```

## Advanced Techniques

### 1. Resolution Source Monitoring

```
Set up automated monitoring:

For data-driven markets:
- RSS feeds for official sources
- Email alerts for publications
- API access to data sources
- Scheduled checks (cron jobs)

Benefits:
- Earlier detection
- Faster execution
- Better prices
- More opportunities
```

### 2. Historical Resolution Analysis

```
Before trading:
Check market creator's history:

Red flags:
- Frequent disputes
- Ambiguous resolutions
- Delays in resolution
- Controversial decisions

Green flags:
- Clear, consistent resolutions
- Fast resolution times
- Uses official sources
- Transparent process
```

### 3. Arbitrage Across Platforms

```
Same event on multiple platforms:

Platform A: Price 0.25, closes in 12h
Platform B: Price 0.85, closes in 12h
Official data: Shows outcome is YES

Opportunity:
- Buy on Platform A @ 0.25
- Sell on Platform B @ 0.85
- Profit: 0.60 per share (240% ROI)

Challenges:
- Capital needed on both platforms
- Different resolution processes
- Execution timing
- Withdrawal times
```

## Performance Tracking

### Track These Metrics

```
For each trade:
☐ Entry date/time
☐ Entry price
☐ Confidence level (%)
☐ Position size ($ and % of account)
☐ Information source
☐ Time until close
☐ Resolution date
☐ Final outcome
☐ Profit/loss
☐ Post-trade notes (what went right/wrong)

Monthly summary:
- Total trades
- Win rate
- Average ROI
- Best/worst trades
- Common mistakes
- Lessons learned
```

### Evaluate Your Edge

```
After 20 trades:

Calibration check:
- 99% confidence trades: Should win 19-20 (95%+)
- 90% confidence trades: Should win 17-19 (85-95%)
- 80% confidence trades: Should win 15-17 (75-85%)

If you're overconfident:
- 99% confidence but only 80% win rate
- Reduce position sizes
- Be more conservative in assessment
- Improve verification process

If you're underconfident:
- 80% confidence but 95% win rate
- Can increase position sizes
- May be missing opportunities
- Can be more aggressive
```

## Common Mistakes and Solutions

### ❌ Mistake 1: Rushing Verification

```
WRONG:
- See opportunity, trade immediately
- Skip verification steps
- "It's obvious!"

RESULT:
- Misread resolution criteria
- Used wrong data source
- Lost money on "sure thing"

CORRECT:
- Always verify (10-15 minutes)
- Check multiple sources
- Read criteria carefully
- Even if "obvious"
```

### ❌ Mistake 2: Confirmation Bias

```
WRONG:
- Want market to resolve YES
- Only look for confirming information
- Ignore contradictory data
- See what you want to see

RESULT:
- Miss key information
- Overconfident in wrong direction
- Surprised by resolution

CORRECT:
- Actively seek contradicting information
- Steel-man the opposite case
- Be willing to skip if uncertain
- Objectivity over wishful thinking
```

### ❌ Mistake 3: Oversizing "Sure Things"

```
WRONG:
- "This is 100% certain!"
- Bet 50% of account
- Resolution goes against you
- Devastating loss

RESULT:
- Even "certainties" have 1-5% resolution risk
- One mistake wipes out many wins
- Emotional damage

CORRECT:
- Max 20% even on best opportunities
- Account for resolution risk
- Remember: you can be wrong
- Preserve capital
```

### ❌ Mistake 4: Ignoring Time Zones

```
WRONG:
Market: "Will X happen by end of Tuesday?"
Current: Tuesday 11pm your time
Trade: Buy at 0.15 (thinking plenty of time)
Resolution: Uses ET time, already Wednesday ET
Outcome: Market resolves NO
Loss: 100%

RESULT:
- Didn't check time zone
- Thought had more time
- Costly mistake

CORRECT:
- Always check time zone in criteria
- Convert to your local time
- Add buffer for safety
- When in doubt, ask/skip
```

### ❌ Mistake 5: Trading Illiquid Markets

```
WRONG:
- Perfect setup, great price
- Try to buy $5,000 worth
- Only $500 liquidity available
- Can't enter full position

RESULT:
- Partial fill only
- Worse average price (slippage)
- Opportunity cost (too small)

CORRECT:
- Check liquidity BEFORE planning position
- Size position to available liquidity
- If insufficient, skip or reduce size
- Don't force oversized trades
```

## Conclusion

Markets closing soon offer unique opportunities through:

**Information arbitrage:**
- Outcomes already determined
- Prices haven't adjusted
- Fast, high-ROI trades

**Predictable outcomes:**
- Strong indicators near resolution
- Lower uncertainty
- Edge from research/analysis

**Information asymmetry:**
- Your knowledge vs market
- Timing advantages
- Domain expertise

Keys to success:
- Rigorous verification (never skip!)
- Proper position sizing
- Fast but careful execution
- Resolution criteria understanding
- Accept small losses on mistakes

Start conservatively, verify thoroughly, and build confidence over time.

## Related Examples

- [Rapid Price Movement](../find-rapid-price-movement/) - For momentum/reversion trading
- [Related Markets Arbitrage](../find-related-markets-arbitrage/) - For structural arbitrage
- [NegRisk Opportunities](../find-negrisk-opportunities/) - For capital efficiency

## Additional Resources

- [Polymarket Resolution Guidelines](https://docs.polymarket.com/#resolution)
- [Official Data Sources List]
- Main repository README for other strategies
