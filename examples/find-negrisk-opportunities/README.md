# NegRisk Opportunities Example

## Overview

This example demonstrates how to find and analyze **NegRisk-enabled events** on Polymarket where the probability sum exceeds 1.0, allowing for capital-efficient position taking.

## What is NegRisk?

NegRisk (Negative Risk) is a Polymarket feature that significantly reduces collateral requirements when trading mutually exclusive outcomes:

- **Standard Markets**: To buy all outcomes, you need collateral equal to the sum of all prices
- **NegRisk Markets**: To buy all outcomes, you only need 1.0 collateral (the maximum payout)

## Capital Efficiency Example

Consider an event with 3 mutually exclusive markets priced at:
- Market A: 0.40 (40%)
- Market B: 0.35 (35%)
- Market C: 0.28 (28%)
- **Sum**: 1.03 (103%)

### Standard Approach
- Collateral needed: $1.03
- Payout when resolved: $1.00
- Capital efficiency: Locked $1.03 to receive $1.00

### NegRisk Approach
- Collateral needed: $1.00
- Payout when resolved: $1.00
- **Savings**: $0.03 (2.9% reduction in capital requirements)

## IMPORTANT: This is NOT Arbitrage!

### Why NegRisk ≠ Arbitrage

Many traders mistakenly believe NegRisk provides arbitrage opportunities. **This is incorrect**:

1. **No guaranteed profit**: You still pay 1.0 and receive 1.0 at resolution
2. **Not risk-free**: You still need one outcome to occur
3. **Cannot mix markets**: NegRisk collateral is locked; you cannot buy via NegRisk then sell in standard markets

### When is NegRisk Useful?

NegRisk is valuable for:
- **Capital efficiency**: Lock up less capital for the same position
- **Portfolio optimization**: Free up capital for other trades
- **Long-term positions**: Hold positions with lower capital cost
- **Hedging strategies**: Reduce total collateral requirements

## When Does NegRisk Have Advantage?

### ✅ Probability Sum > 1.0 (Overpriced Markets)
- Standard cost: Sum of all prices (e.g., 1.03)
- NegRisk cost: 1.0
- **NegRisk saves capital**: 3% reduction in this example

### ❌ Probability Sum < 1.0 (Underpriced Markets - TRUE ARBITRAGE)
- Standard cost: Sum of all prices (e.g., 0.98)
- NegRisk cost: 1.0
- **Standard is better**: True arbitrage opportunity (pay 0.98, receive 1.0)
- **NEVER use NegRisk for underpriced arbitrage!**

## How to Use This Example

### 1. Run the Example

```bash
cd examples/find-negrisk-opportunities
go run main.go
```

### 2. Understanding the Output

The example finds events with:
- ✅ NegRisk support enabled
- ✅ Probability sum > 1.0 (where NegRisk provides advantage)
- ✅ Multiple mutually exclusive markets
- ✅ Sufficient liquidity for execution

For each opportunity, it displays:

#### Event Information
- Event title and ID
- Number of markets
- Direct link to Polymarket

#### NegRisk Configuration
- NegRisk Market ID
- Fee structure (in basis points)
- Collateral requirements

#### Probability Analysis
- Sum of all market prices
- Deviation from 1.0
- Individual market prices

#### Capital Efficiency Comparison

**When Sum > 1.0** (NegRisk Advantage):
```
Standard Markets:
   Total Cost:       $1.0300 (sum of all prices)
   ✓ Savings:        $0.0300 (2.91% reduction)

NegRisk Markets:
   Total Cost:       $1.0000 (fixed collateral)
   Capital Efficiency: 2.91% better ✓
```

**When Sum < 1.0** (Standard is Better):
```
⚠️  Analysis:     NegRisk is NOT advantageous here!
❌ Extra cost:    $0.0200 (2.04% more expensive)

For underpriced markets, use STANDARD markets for true arbitrage!
```

### 3. Interpreting Opportunities

#### Excellent Opportunity (Score 80-100)
- High deviation (> 5%)
- Good liquidity (> $10,000)
- Multiple markets (3-5+)
- Reasonable fees
- **Action**: Consider taking position with NegRisk

#### Good Opportunity (Score 60-80)
- Moderate deviation (3-5%)
- Adequate liquidity
- Multiple markets
- **Action**: Evaluate based on your capital efficiency goals

#### Moderate Opportunity (Score < 60)
- Small deviation (< 3%)
- Lower liquidity
- **Action**: May not be worth the transaction costs

## Step-by-Step Trading Guide

### Strategy 1: Capital-Efficient Position Taking (Sum > 1.0)

**Example**: Event with sum = 1.03

#### Step 1: Verify Market Conditions
```
✅ Check: Are markets truly mutually exclusive?
✅ Check: Do all markets accept orders?
✅ Check: Is liquidity sufficient for your size?
✅ Check: Calculate fees impact
```

#### Step 2: Calculate Position Size
```
Example with $1,000 capital:
- Without NegRisk: Can take ~$971 position (1000/1.03)
- With NegRisk: Can take $1,000 position
- Advantage: 2.9% more exposure with same capital
```

#### Step 3: Execute via NegRisk
1. Access the NegRisk Market ID (shown in output)
2. Use NegRisk-enabled interface/API
3. Place orders across all markets
4. Monitor until resolution

#### Step 4: Resolution
- One market resolves to YES (pays $1.00)
- Other markets resolve to NO (pay $0.00)
- Net result: $1.00 received per $1.00 collateral

### Strategy 2: Underpriced Arbitrage (Sum < 1.0)

**CRITICAL**: If the example shows "NegRisk is NOT advantageous", this is TRUE ARBITRAGE!

#### Step 1: Confirm Underpricing
```
Example: Sum = 0.98
✅ Standard cost: $0.98
✅ Payout: $1.00
✅ Profit: $0.02 (2% guaranteed)
```

#### Step 2: Use STANDARD Markets (NOT NegRisk!)
```
❌ WRONG: Use NegRisk (costs $1.00, no profit)
✅ CORRECT: Use standard markets (costs $0.98, profit $0.02)
```

#### Step 3: Execute Arbitrage
1. Buy equal dollar amounts in ALL markets using standard interface
2. Hold until resolution
3. Collect $1.00 payout
4. Realize $0.02 profit per dollar invested

## Risk Considerations

### General Risks
- **Market exclusivity**: Verify outcomes are truly mutually exclusive
- **Multiple winners**: Some events may resolve multiple markets to YES
- **None-of-above**: Check if event can resolve to outcomes not listed
- **Fees**: Calculate maker/taker fees (typically 2-4%)
- **NegRisk fees**: Additional fees may apply (shown in output)

### NegRisk-Specific Risks
- **Collateral lock-up**: Capital is locked until resolution
- **No early exit arbitrage**: Cannot convert to standard positions for profit
- **Smart contract risk**: NegRisk uses different contract logic
- **Lower liquidity**: NegRisk markets may have less depth

### Resolution Risks
- **Manual resolution**: May take time after event concludes
- **Disputes**: Resolution may be contested
- **Ambiguous outcomes**: Interpretation may differ from expectations

## Common Mistakes to Avoid

### ❌ Mistake 1: Using NegRisk for Underpriced Arbitrage
```
Event sum = 0.98

WRONG:
  Use NegRisk → Pay $1.00 → Receive $1.00 → No profit

CORRECT:
  Use Standard → Pay $0.98 → Receive $1.00 → Profit $0.02
```

### ❌ Mistake 2: Thinking NegRisk is "Free Money"
- NegRisk only provides **capital efficiency**, not profit
- You still need correct outcome to occur
- You still pay transaction fees

### ❌ Mistake 3: Trying to Cross-Market Arbitrage
```
Idea: Buy via NegRisk ($1.00) → Sell in standard markets (sum $1.03) → Profit $0.03?

Why it doesn't work:
  - NegRisk collateral is LOCKED
  - Cannot sell positions separately
  - Converting to standard requires additional collateral
  - The mechanism prevents this by design
```

### ❌ Mistake 4: Ignoring Non-Exclusive Outcomes
- Always verify markets are truly mutually exclusive
- Check resolution criteria carefully
- Some events may have "none of above" possibility

## Advanced Strategies

### 1. Capital Rotation
Use NegRisk to free up capital for other trades:
```
Portfolio: $10,000

Without NegRisk:
  Trade A: $1,030 locked (sum 1.03)
  Trade B: $1,050 locked (sum 1.05)
  Trade C: $980 locked (sum 0.98)
  Remaining: $6,940

With NegRisk (for A and B only):
  Trade A: $1,000 locked (NegRisk)
  Trade B: $1,000 locked (NegRisk)
  Trade C: $980 locked (standard - true arb!)
  Remaining: $7,020

Extra capital: $80 (0.8% more efficient)
```

### 2. Fee Optimization
Calculate total costs including fees:
```
Position: $1,000
Maker fee: 2% = $20
NegRisk fee: 0.5% = $5
Total cost: $1,025

Compare to standard:
Standard cost: $1,030
Standard fee: 2% = $20.60
Total: $1,050.60

NegRisk savings: $25.60
```

### 3. Liquidity Provision
Provide liquidity on NegRisk markets:
- Lower capital requirements
- Capture spread on both sides
- Manage inventory more efficiently

## Real-World Example

From the output:
```
Event: "How many people will Trump deport in 2025?"
Markets: 9 mutually exclusive outcomes
Probability Sum: 1.018 (1.8% overpriced)
NegRisk Enabled: Yes
```

### Analysis
- **Standard approach**: Need $1.018 collateral per position
- **NegRisk approach**: Need $1.000 collateral per position
- **Savings**: $0.018 per dollar (1.8% capital efficiency)

### Decision Framework
1. **Is 1.8% efficiency worth it?**
   - For $10,000 position: Saves $180 in locked capital
   - Consider holding period and alternative uses

2. **Fee impact**:
   - Trading fees: ~2% = $200
   - NegRisk fees: Check output
   - Net benefit after fees

3. **Liquidity check**:
   - Total liquidity: Check output
   - Your size vs market depth
   - Slippage estimate

## Monitoring and Exit

### During Holding Period
1. **Monitor market movements**: Prices may converge toward 1.0
2. **Watch for new information**: May want to adjust position
3. **Track resolution timeline**: Plan for capital lock-up
4. **Check for disputes**: Stay informed on resolution process

### At Resolution
1. **Verify outcome**: Confirm which market resolves YES
2. **Claim payout**: Receive 1.0 per position
3. **Calculate realized efficiency**: Compare to standard approach
4. **Document results**: Learn for future opportunities

## Conclusion

NegRisk opportunities provide **capital efficiency**, not arbitrage profit. Use this tool to:

✅ Identify events where NegRisk reduces capital requirements
✅ Optimize portfolio capital allocation
✅ Avoid mistakenly using NegRisk for true arbitrage (sum < 1.0)
✅ Understand the trade-offs and risks

Remember: The goal is capital efficiency, not guaranteed profit. Always verify market conditions and understand the mechanics before trading.

## Additional Resources

- [Polymarket NegRisk Documentation](https://docs.polymarket.com/)
- [Related Markets Arbitrage Example](../find-related-markets-arbitrage/) - For true arbitrage opportunities
- Main README for other trading strategies

## Support

For questions or issues:
- Check the main repository README
- Review Polymarket documentation
- Test with small positions first
- Understand the risks before trading
