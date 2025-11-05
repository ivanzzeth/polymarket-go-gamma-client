# New Active Markets - Early Market Making Example

## Overview

This example identifies **early market making opportunities** by finding recently launched markets where you can establish positions before competition intensifies and spreads tighten.

## The Early Mover Advantage

### Why Trade New Markets?

```
Day 1: Market launches
- Liquidity: $2,000
- Spread: 6% (0.06)
- Market makers: 1-2
- Your opportunity: High margins, low competition

Day 7: Market established
- Liquidity: $15,000
- Spread: 2% (0.02)
- Market makers: 5-7
- Opportunity: Still okay, but compressed

Day 30: Mature market
- Liquidity: $50,000
- Spread: 0.5% (0.005)
- Market makers: 15+
- Opportunity: Highly competitive, thin margins
```

**Key insight**: Early market makers capture wide spreads before competition drives them down, and can establish dominant positions at favorable prices.

### The First-Mover Advantages

#### 1. Wide Spreads
```
New market (Day 1):
- Bid: 0.42
- Ask: 0.58
- Spread: 0.16 (16%!)
- Your quotes: 0.45 / 0.55 (10% spread)
- Still very profitable, much better than later

Established market (Day 30):
- Bid: 0.485
- Ask: 0.495
- Spread: 0.01 (1%)
- Your quotes: Must be 0.487 / 0.493 (0.6% spread)
- Thin margins, high competition
```

#### 2. Less Competition
```
New market:
- Total liquidity: $3,000
- Your capital: $500 (17% market share)
- Other MMs: 2-3 participants
- You're a major player

Established market:
- Total liquidity: $50,000
- Your capital: $500 (1% market share)
- Other MMs: 15+ participants
- You're a small fish
```

#### 3. Price Discovery Participation
```
New market:
- Price is uncertain: Could be 0.30 or 0.60
- You help discover "fair value"
- Your quotes influence the market
- Early positions at good prices
- Profit from volatility

Established market:
- Fair value known: ~0.48-0.52
- You're price-taker, not maker
- Little influence on market
- Late positions at crowded prices
- Less volatility to profit from
```

#### 4. Reputation Building
```
Early presence:
- Become "the market maker" for this market
- Traders recognize your quotes
- Build trust and relationships
- Preferential order flow
- Information advantages

Late entry:
- Just another participant
- No special recognition
- Compete equally with everyone
- Standard order flow
```

## How to Use This Example

### 1. Run the Example

```bash
cd examples/find-new-active-markets
go run main.go
```

### 2. Configuration

Edit `main.go` to adjust parameters:

```go
daysOld := 7.0             // Markets less than 7 days old
maxLiquidity := 10000.0    // Maximum liquidity (early stage)
targetCount := 5           // Number of opportunities
```

### 3. Understanding the Output

The example finds markets with:
- ✅ Age < 7 days (newly launched)
- ✅ Liquidity < $10,000 (early stage)
- ✅ Minimum liquidity > $2,000 (enough to start)
- ✅ Open and accepting orders
- ✅ Active trading

For each opportunity, it displays:

#### Market Age
```
Created:              2025-01-15 14:23:00
Age:                  1.2 days 🆕 VERY NEW!

Age categories:
≤ 1 day: 🆕 VERY NEW! (Best opportunity)
1-3 days: 🆕 NEW! (Excellent opportunity)
3-7 days: (Good opportunity)
> 7 days: (Passed early stage)
```

#### Opportunity Score (0-100)
```
⭐ Opportunity Score:  85/100 🔥 EXCELLENT!

Score components:
- Age: Newer = higher score (up to +20)
- Liquidity: Lower = higher score (up to +15)
- Volume: Some activity = higher score (up to +10)
- Spread: Wider = higher score (up to +10)
- Active orders: Accepting = +5
- Featured: -5 (more competition)

Ratings:
80-100: 🔥 EXCELLENT!
70-80: ✨ VERY GOOD
60-70: ✓ GOOD
< 60: - MODERATE
```

#### Liquidity Analysis
```
Current Liquidity:    $3,500 💎 (Very low - great opportunity!)

Assessment:
< $1,000: 💎 Very low - great opportunity!
$1k-3k: ✓ Low - good opportunity
$3k-5k: Moderate
> $5k: Higher (more competition)
```

#### Current Spread
```
Current Spread:       0.0450 (4.50%) 💰 (Wide - great for MM!)

Assessment:
> 5%: 💰 Wide - great for market making!
3-5%: ✓ Good for market making
1.5-3%: Moderate
< 1.5%: Tight (competitive)
```

#### Market Making Strategy
The tool provides:
- Early mover advantages specific to this market
- Recommended approach (conservative sizing initially)
- Position management guidelines
- Quote sizing suggestions
- Rebalancing frequency

#### Early Stage Indicators
```
📊 Early Stage Indicators:
✓ Has 24h trading volume ($850)
✓ Low liquidity ($3,500) - less competition
✓ Wide spread (4.50%) - good profit potential
✓ High turnover ratio (1.2x) - active market
```

## Step-by-Step Market Making Guide

### Phase 1: Opportunity Assessment (First 30 Minutes)

#### Step 1: Evaluate the Market

```
Critical evaluation checklist:

📊 Market Quality:
✅ Question is clear and unambiguous?
✅ Resolution criteria are objective?
✅ Resolution date is reasonable (weeks/months away)?
✅ Topic has sustained interest potential?
✅ Category is actively traded?

⚠️ Red Flags:
❌ Vague question or criteria
❌ Subjective resolution
❌ Resolves in < 1 week (not enough time)
❌ Very niche topic (low interest)
❌ First market from new creator (resolution risk)

🎯 Competition Analysis:
✅ Current liquidity level (lower = less competition)
✅ Number of active market makers (fewer = better)
✅ Spread width (wider = less competition)
✅ Featured status (not featured = less visible)

💰 Profit Potential:
✅ Spread > 3% (good margins)
✅ Some volume (activity exists)
✅ Growing interest (topic is trending)
✅ Time until resolution (more time = more profit)
```

**Example assessment:**

```
Market: "Will inflation exceed 3% in Q2?"
Age: 1.5 days old
Liquidity: $4,200
Spread: 4.8%
Volume 24h: $680

Evaluation:
✅ Clear question (objective metric)
✅ Official data source (CPI from BLS)
✅ Resolution: 4 months away (good timeframe)
✅ Topic: Economics (actively traded category)
✅ Spread: 4.8% (excellent margins)
✅ Activity: $680/day (some interest)
⚠️ Not featured (good - less competition)

Competition:
- Current MMs: 2 based on order book
- Your potential share: 20-30% with $500-1000
- First mover advantage: Strong

Decision: EXCELLENT OPPORTUNITY
Score: 82/100
```

#### Step 2: Analyze the Order Book

```
Check current state:

Bids:
  0.45: $300 (probably the main MM)
  0.43: $800 (probably AMM)
  0.40: $400

Asks:
  0.50: $350 (probably same main MM)
  0.52: $750 (probably AMM)
  0.55: $600

Analysis:
- Main MM: $300-350 per side, 5% spread
- Spread: 0.45 to 0.50 (5%)
- Mid: 0.475
- Your opportunity: Quote inside at 0.46/0.49 (3% spread)
- Or match at 0.45/0.50 (5% spread, less competition)

Strategy decision:
- Option A: Compete (0.46/0.49) - more fills, lower margin
- Option B: Match (0.45/0.50) - fewer fills, higher margin
- Recommendation for Day 1: Option B (match, observe)
```

#### Step 3: Determine Position Size

```
Conservative approach for new markets:

Your capital: $5,000 available
Market liquidity: $4,200

Position sizing:
- Day 1: $300-500 (7-12% of liquidity)
- Week 1: $500-800 (scale up if successful)
- Week 2+: $800-1,200 (established presence)

Rationale:
- Start small (learning period)
- Test market dynamics
- Minimize risk
- Room to scale

Example:
- Start: $400 ($200 per side)
- Quote: 200 shares at 0.45 bid, 200 at 0.50 ask
- Capital needed: ~$400 total
- Leaves $4,600 for other opportunities
```

### Phase 2: Initial Entry (First 3 Days)

#### Step 1: Place Conservative Quotes

```
Day 1 strategy: Join the best quotes

Existing market:
- Best bid: 0.45
- Best ask: 0.50

Your quotes:
- Bid: 0.45 (match, don't lead yet)
- Ask: 0.50 (match)
- Size: $200 per side

Rationale:
- Observe market dynamics first
- See who else is participating
- Learn typical order sizes
- Understand fill patterns
- Lower risk while learning
```

#### Step 2: Monitor Price Discovery

```
Price discovery observations:

Hour 1-4:
- Trades at 0.47, 0.48, 0.46
- Conclusion: Fair value ~0.47
- Action: Adjust mid to 0.47

Hour 4-8:
- News breaks, trades at 0.52, 0.54
- Conclusion: Value increased
- Action: Move quotes up to 0.52/0.57

Hour 8-24:
- Settles at 0.50-0.53
- Conclusion: New fair value ~0.515
- Action: Quote around 0.51/0.52

Learning:
- Market is volatile (wide moves)
- News-sensitive
- Need to adjust quotes frequently
- Good for market making (volatility = opportunity)
```

#### Step 3: Establish Your Presence

```
Days 1-3 goals:

✅ Be consistently present:
   - Update quotes every 2-4 hours minimum
   - Respond to market moves
   - Maintain competitive quotes

✅ Build volume:
   - Get filled 10-20 times
   - Turn over capital 2-3x
   - Establish track record

✅ Learn the market:
   - Who are other MMs?
   - What times are most active?
   - How does news affect it?
   - What's typical order size?

✅ Maintain discipline:
   - Stick to risk limits
   - Don't over-size
   - Keep inventory balanced
   - Take profits regularly

Results after 3 days:
- Filled: 25 times
- Turnover: 3x your capital
- Gross profit: $60-120 (15-30% on capital)
- Lessons learned: Ready to scale
```

### Phase 3: Expansion (Days 4-7)

#### Step 1: Scale Up Gradually

```
If successful in Days 1-3:

Day 1-3 metrics:
- Position size: $400
- Gross profit: $90
- ROI: 22.5%
- Fill rate: Good
- Competition: Manageable

Day 4-7 scaling:
- New position: $700 (75% increase)
- Expected profit: $150-180
- ROI target: 20-25%
- Monitor: Competition changes

Rationale:
- Proven strategy works
- Market dynamics understood
- Still early (low competition)
- Scale while opportunity exists
```

#### Step 2: Tighten Your Spreads

```
As market matures, tighten spreads:

Day 1-3:
- Your spread: 5% (0.45 / 0.50)
- Strategy: Match existing market
- Reason: Learning, conservative

Day 4-7:
- Your spread: 3-4% (0.46 / 0.49)
- Strategy: More competitive
- Reason: Established, want more flow

Day 7+:
- Your spread: 2-3% (0.47 / 0.49)
- Strategy: Competitive quotes
- Reason: Market maturing, need fills

Trade-off:
- Tighter spread = more fills, lower margin
- Wider spread = fewer fills, higher margin
- Optimize for maximum absolute profit
```

#### Step 3: Manage Competition

```
As competitors enter:

Week 1: Just you and 1-2 others
- Strategy: Wide spreads (4-5%)
- Position: $400-700
- Approach: Relaxed

Week 2: 3-4 competitors appear
- Strategy: Moderate spreads (3-4%)
- Position: $700-1,000
- Approach: Competitive but sustainable

Week 3: 5+ competitors
- Strategy: Tight spreads (2-3%)
- Position: Consider stabilizing
- Approach: Optimize efficiency

Decision point (Week 3-4):
- If still profitable (>1% spread captured): Continue
- If spreads <1%: Scale down or exit
- Consider: Opportunity cost vs other new markets
```

### Phase 4: Maturity Decision (Day 14+)

#### Option A: Continue (If Profitable)

```
Continue if:
✅ Spreads still >2%
✅ Good turnover (5x+ per week)
✅ Competition manageable (5-7 MMs)
✅ You're established (recognized participant)
✅ Still profitable (>15% monthly ROI)

Approach:
- Maintain presence
- Optimize quote placement
- Manage inventory efficiently
- Harvest ongoing profits
```

#### Option B: Scale Down (If Compression)

```
Scale down if:
⚠️ Spreads compressed <1.5%
⚠️ Too many competitors (>10 MMs)
⚠️ Your fills decreasing
⚠️ ROI declining (<10% monthly)
⚠️ Better opportunities elsewhere

Approach:
- Reduce position 50%
- Maintain small presence
- Focus on high-turnover times
- Allocate capital to newer markets
```

#### Option C: Exit (If Uneconomic)

```
Exit if:
❌ Spreads <1%
❌ Heavy competition (15+ MMs)
❌ Barely break even after fees
❌ ROI <5% monthly
❌ Market resolving soon

Approach:
- Close all positions
- Cancel orders
- Redeploy to new opportunities
- Document lessons learned
```

## Real-World Example Analysis

### Example from Output:

```
Market: "Will Fed cut rates in March?"
Age: 2.3 days
Current Liquidity: $5,800
Current Volume: $12,400
24h Volume: $1,850
Spread: 3.8%
Opportunity Score: 78/100 (✨ VERY GOOD)
```

### Week-by-Week Strategy:

**Week 1 (Days 1-7):**
```
Your plan:
- Entry: Day 2 (you're very early!)
- Position: $500 ($250 per side)
- Quotes: Match existing spread (3.8%)
- Target: 2-3 turnovers per week
- Goal: Learn market, build presence

Expected results:
- Gross profit: $90-120
- ROI: 18-24% for week
- Fills: 15-20 times
- Lessons: Price discovery patterns

Actual monitoring:
- Check every 3-4 hours
- Rebalance when >$100 skew
- Adjust for news
- Track competition
```

**Week 2 (Days 8-14):**
```
Assessment:
- Week 1 went well: $105 profit
- Competition: Now 4 MMs total
- Spread: Compressed to 2.8%
- Your adaptation needed

Your plan:
- Position: $700 (scale up 40%)
- Quotes: 2.5-3% spread (more competitive)
- Target: 3-4 turnovers per week
- Goal: Maintain edge as competition grows

Expected results:
- Gross profit: $120-150
- ROI: 17-21% for week
- More fills: 25-30 times
- Observation: Still profitable
```

**Week 3 (Days 15-21):**
```
Assessment:
- Week 2: $135 profit (good!)
- Competition: Now 6-7 MMs
- Spread: Further compressed to 1.8%
- Decision time approaching

Your options:

A) Stay and optimize ($700):
   - Projected profit: $90-120/week
   - ROI: ~15%/week (still good)
   - Time: Same monitoring needed

B) Scale back to $400:
   - Projected profit: $50-70/week
   - ROI: ~15%/week
   - Time: Reduced commitment
   - Deploy $300 to newer market

C) Full exit:
   - Close positions
   - Redeploy $700 to new market
   - Opportunity cost consideration

Recommendation: Option B
- Keep small presence
- Free up capital for new markets
- Maintain optionality
- Better overall portfolio return
```

## Advanced Techniques for New Markets

### 1. Reputation Building

```
Become "the market maker":

Consistency tactics:
- Always quote reasonable spreads
- Respond to all order flow
- Quick rebalancing (reduce inventory risk)
- Fair pricing (earn trust)

Benefits:
- Traders prefer your quotes
- Better fill rates
- Information flow (see order interest)
- Sustainable advantage

Example:
Market regular: "I always trade with Alice, she's reliable"
= You (Alice) get preferential order flow
= Higher volumes, better profits
```

### 2. Multi-Market Strategy

```
Instead of one new market, do 3-5:

Portfolio approach:
- Market A (Day 1): $300
- Market B (Day 2): $300
- Market C (Day 3): $300
- Market D (Day 5): $400
- Market E (Day 6): $400

Benefits:
- Diversification (not all succeed)
- Better capital utilization
- More learning opportunities
- Compounding advantage

Management:
- Standardized approach
- Consistent position sizing
- Regular portfolio review
- Replace underperformers
```

### 3. Price Discovery Leadership

```
In new markets, early MMs influence price:

Strategy:
- Research fair value estimate
- Quote around your estimate
- Adjust based on order flow
- Lead price discovery

Example:
New market: "Will X happen?"
Your research: 65% probability
Your quotes: 0.63 / 0.67 (around 0.65)

Order flow:
- More buying than selling
- Adjust up: 0.68 / 0.72
- Update fair value estimate

Benefit:
- Your quotes help "make the market"
- Early positioning at good prices
- Profit from information edge
- Establish thought leadership
```

### 4. Monitoring Automation

```
As you manage more markets:

Automation needs:
- Price alerts (market moves)
- Fill notifications (orders executed)
- Inventory alerts (getting unbalanced)
- Competition monitoring (spread compression)

Tools:
- API connections
- Custom dashboards
- Automated rebalancing rules
- Performance tracking

Benefits:
- Scale to more markets
- Faster reactions
- Less manual monitoring
- Better results
```

## Risk Management for New Markets

### Unique Risks of New Markets

#### 1. Resolution Risk
```
Problem: New market creator might have unclear resolution

Protection:
- Check creator's history
- Read criteria 3 times
- Look for ambiguities
- Reduce size for new creators
- Skip if criteria unclear

Example:
New creator, first market, vague criteria
= Skip OR very small position (50% normal size)
```

#### 2. Volume Risk
```
Problem: Volume may not materialize

Protection:
- Monitor daily volume trends
- Set minimum volume thresholds
- Exit if volume doesn't develop
- Time limit (exit if dead after 1 week)

Example:
Day 1: $500 volume - okay
Day 3: $200 volume - declining
Day 7: $50 volume - time to exit
```

#### 3. Price Discovery Volatility
```
Problem: Price can swing wildly in early days

Protection:
- Wider stops (allow for volatility)
- Smaller position sizes
- More frequent rebalancing
- Don't chase moves

Example:
Price: 0.50 → 0.70 → 0.45 (first day)
Your reaction: Stay calm, keep quoting, rebalance
Don't panic: This is normal price discovery
```

#### 4. Competition Surge
```
Problem: Market suddenly gets crowded

Protection:
- Monitor liquidity growth
- Watch spread compression
- Have exit plan ready
- Re-evaluate daily

Example:
Day 1: 2 MMs, 4% spread
Day 5: 8 MMs, 1.5% spread (sudden surge)
Action: Reassess, possibly exit
```

### Position Limits for New Markets

```
Conservative limits:

Per market:
- Maximum: 10% of your capital
- Start: 5% of your capital
- Scale up: Only after 3 days success

Total new markets:
- Maximum: 40% of capital in new markets
- Reason: Higher risk profile
- Balance: 60% in established markets

Example with $10,000:
- New market A: $500 (5%)
- New market B: $400 (4%)
- New market C: $600 (6%)
- Established markets: $6,500 (65%)
- Reserve: $2,000 (20%)
- Total new: $1,500 (15%) ✅ Within limit
```

## Performance Tracking

### Key Metrics for New Markets

```
Track separately from established:

Per market:
☐ Entry date and market age
☐ Initial liquidity and spread
☐ Your position size over time
☐ Weekly turnover achieved
☐ Gross profit by week
☐ When competition arrived
☐ When you exited/scaled
☐ Total ROI
☐ Time invested

Portfolio level:
- Success rate (% of markets profitable)
- Average ROI per market
- Average holding period
- Best/worst performers
- Common success factors
- Common failure factors
```

### Evaluate Your Approach

```
After 10 new markets:

Questions to answer:
1. What's your success rate?
   Target: >70% profitable

2. Average ROI?
   Target: >50% over holding period

3. How long do you stay?
   Typical: 2-4 weeks

4. When do you enter?
   Best: Days 1-3
   Okay: Days 4-7
   Late: Days 8+

5. What factors predict success?
   - Category? (some better than others)
   - Initial liquidity level?
   - Spread width?
   - Featured status?

Optimize:
- Focus on best categories
- Refine entry timing
- Improve position sizing
- Better exit discipline
```

## Common Mistakes and Solutions

### ❌ Mistake 1: Oversizing Too Early

```
WRONG:
- Day 1, new market
- "This looks great!"
- Deploy $2,000 (40% of capital)
- Market doesn't develop
- Capital tied up, can't redeploy

CORRECT:
- Day 1: $300-500 (start small)
- Day 4: Scale to $600-800 (if successful)
- Week 2: Scale to $1,000+ (if still working)
- Always leave room to pivot
```

### ❌ Mistake 2: Fighting Price Discovery

```
WRONG:
- You think fair value is 0.60
- Market trading at 0.75
- You keep selling at 0.60
- Accumulate huge long position
- Market stays at 0.75
- You were wrong about value

CORRECT:
- Have opinion (0.60)
- Market disagrees (0.75)
- Quote around market (0.73/0.77)
- Respect market's wisdom
- Don't fight with large position
```

### ❌ Mistake 3: Ignoring Competition

```
WRONG:
- Week 1: Great (4% spread, $200/week)
- Week 2: Okay (2.5% spread, $150/week)
- Week 3: Declining (1.5% spread, $80/week)
- Week 4: Break even (0.8% spread, $20/week)
- Still there, wasting time and capital

CORRECT:
- Week 1: Great, continue
- Week 2: Monitor closely
- Week 3: Scale down 50%
- Week 4: Exit completely
- Redeploy to new opportunities
```

### ❌ Mistake 4: Poor Time Management

```
WRONG:
- Have 8 new markets
- Each needs monitoring every 2 hours
- 8 × 8 checks/day = 64 checks!
- Overwhelming, mistakes happen
- Some markets neglected

CORRECT:
- Start with 2-3 new markets
- Manageable monitoring (16-24 checks/day)
- Quality over quantity
- Scale only when efficient
- Use automation when possible
```

### ❌ Mistake 5: No Exit Plan

```
WRONG:
- Enter Day 1
- Do well Weeks 1-2
- Competition grows Week 3
- "Maybe it'll get better"
- Week 6, barely profitable
- Opportunity cost = huge loss

CORRECT:
- Set exit criteria upfront
  - Exit if spread <1.5%
  - Exit if competition >10 MMs
  - Exit if ROI <10% monthly
  - Review every 2 weeks
- Stick to criteria
- Don't get emotionally attached
```

## Success Checklist

### Before Entering New Market:

```
☐ Market Quality:
  ☐ Clear, unambiguous question
  ☐ Objective resolution criteria
  ☐ Reasonable resolution timeline (weeks+)
  ☐ Category has active interest
  ☐ Creator has good history

☐ Opportunity Assessment:
  ☐ Age < 7 days
  ☐ Liquidity < $10,000
  ☐ Spread > 2.5%
  ☐ Some volume ($500+ daily)
  ☐ Accepting orders

☐ Competition:
  ☐ Counted current MMs (< 5 ideal)
  ☐ Analyzed order book depth
  ☐ Checked for dominant MM
  ☐ Assessed your potential share

☐ Your Readiness:
  ☐ Capital available (5-10%)
  ☐ Time to monitor (first 3 days critical)
  ☐ Understand market dynamics
  ☐ Exit criteria defined
```

## Conclusion

New active markets offer the best market making opportunities through:

**Early mover advantages:**
- Wide spreads (3-6%+)
- Low competition (2-3 initial MMs)
- Price discovery participation
- Establishing dominant position

**Optimal approach:**
- Enter very early (Days 1-3)
- Start small ($300-500)
- Scale gradually (50-100% per week)
- Monitor competition closely
- Exit when no longer profitable

**Keys to success:**
- Be first (or at least early)
- Start conservative, scale with success
- Adapt as competition arrives
- Exit disciplined when compressed
- Portfolio approach (multiple markets)

The best profits are made in the first 2-3 weeks. After that, competition typically reduces opportunities. Stay nimble, track results, and always be looking for the next new market.

## Related Examples

- [Low Liquidity / High Volume](../find-low-liquidity-high-volume/) - For established market making
- [Wide Spread Markets](../find-wide-spread-markets/) - Another market making approach
- [Related Markets Arbitrage](../find-related-markets-arbitrage/) - Complementary strategy

## Additional Resources

- [Polymarket CLOB API](https://docs.polymarket.com/#clob-api) - For placing orders
- [Market Making Guide](https://docs.polymarket.com/#market-making)
- Main repository README for other strategies
