# Find Wide Spread Markets

This example demonstrates how to find markets with spreads that are wider than 3x the minimum tick size on Polymarket.

## What is Spread?

The spread is the difference between the best bid and best ask prices in a market. A wider spread typically indicates:
- Lower liquidity
- Higher uncertainty
- Less active market makers
- Potential opportunities for providing liquidity

## What is Tick Size?

The tick size (OrderPriceMinTickSize) is the minimum price increment allowed for orders in the market. For example, if the tick size is 0.01, prices must be in increments of 0.01 (e.g., 0.45, 0.46, 0.47).

## Usage

```bash
cd examples/find-wide-spread-markets
go run main.go
```

## What the Example Does

1. Fetches 100 active markets from Polymarket
2. Analyzes each market's spread and tick size
3. Identifies markets where spread > 3x tick size
4. Prints detailed information for up to 3 such markets

## Output Information

For each market found, the example prints:

### Basic Information
- Market question and ID
- Status (active/closed)

### Spread & Tick Information
- Tick size
- Current spread
- Spread ratio (how many times the tick size)

### Price Information
- Best bid and ask prices
- Last trade price

### Liquidity & Volume
- Total liquidity (broken down by AMM and CLOB)
- Total volume and 24-hour volume

### Order Book Settings
- Minimum order size
- Maker and taker fees
- Whether the market is accepting orders

### Market Details
- Market type and category
- Possible outcomes
- Description

### Analysis
- Interpretation of what the wide spread might indicate
- Potential opportunities

## Use Cases

This example is useful for:
- **Market Makers**: Identify markets with wide spreads where providing liquidity could be profitable
- **Traders**: Find markets with inefficient pricing
- **Researchers**: Analyze market microstructure and liquidity patterns
- **Arbitrageurs**: Spot potential arbitrage opportunities

## Example Output

```
🔍 Finding markets with spread > 3x tick size...
============================================================

📊 Fetched 100 active markets
Analyzing spreads...

✅ Found 3 market(s) with wide spreads:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Market #1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📌 Basic Information:
   Question:     Will X happen before Y?
   Market ID:    12345
   ...

💰 Spread & Tick Information:
   Tick Size:    0.010000
   Spread:       0.045000
   Spread Ratio: 4.50x tick size ⚠️
   ...
```

## Notes

- Markets with very low liquidity may have even wider spreads
- The tick size is set per market and can vary
- Wide spreads don't necessarily indicate bad markets - they might represent legitimate uncertainty
- Always consider transaction costs (fees) when evaluating opportunities
