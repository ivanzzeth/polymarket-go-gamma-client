# Trading Opportunities Analysis

This document outlines various trading opportunities that can be discovered and analyzed using the Polymarket Gamma API.

## Table of Contents

1. [Market Making Opportunities](#market-making-opportunities)
2. [Arbitrage Opportunities](#arbitrage-opportunities)
3. [Price Anomaly Detection](#price-anomaly-detection)
4. [Statistical Arbitrage](#statistical-arbitrage)
5. [Special Situations](#special-situations)
6. [Implementation Examples](#implementation-examples)

---

## Market Making Opportunities

Market making involves providing liquidity to markets by placing both buy and sell orders, profiting from the bid-ask spread.

### 1.1 Wide Spread Markets

**Description**: Markets where the spread (difference between best bid and best ask) is significantly larger than the minimum tick size.

**Identification Criteria**:
- `market.Spread > 3 × market.OrderPriceMinTickSize`
- Filter out: `market.Spread ≈ 1.0` (indicates closed markets with stale data)

**Key API Fields**:
```go
market.Spread                  // Bid-ask spread
market.OrderPriceMinTickSize   // Minimum price increment
market.BestBid                 // Current best bid price
market.BestAsk                 // Current best ask price
market.LiquidityClob          // CLOB liquidity (lower = wider spreads)
```

**Profit Potential**: High spread allows market makers to profit from the difference while providing liquidity.

**Example Implementation**: See `examples/find-wide-spread-markets/`

---

### 1.2 Low Liquidity, High Volume Markets

**Description**: Markets with significant trading activity but insufficient liquidity, creating opportunities for liquidity providers.

**Identification Criteria**:
- High ratio: `market.Volume24hr / market.LiquidityNum`
- `market.Volume24hr > threshold` (e.g., $10,000)
- `market.LiquidityClob < threshold` (e.g., $5,000)

**Key API Fields**:
```go
market.Volume24hr      // 24-hour trading volume
market.VolumeNum       // Total volume
market.LiquidityNum    // Total liquidity
market.LiquidityClob   // CLOB liquidity
market.LiquidityAmm    // AMM liquidity
```

**Profit Potential**: High volume indicates active trading; low liquidity means wider spreads and better market making profits.

---

### 1.3 New Active Markets

**Description**: Recently launched markets where liquidity has not yet been established, providing early mover advantages for market makers.

**Identification Criteria**:
- `market.StartDate` recent (e.g., within last 7 days)
- `market.Active = true`
- `market.AcceptingOrders = true`
- Low `market.LiquidityClob`

**Key API Fields**:
```go
market.StartDate        // Market launch date
market.Active           // Market is active
market.AcceptingOrders  // Currently accepting orders
market.LiquidityClob    // Current CLOB liquidity
```

**Profit Potential**: Early market makers can establish positions before competition increases and spreads tighten.

---

## Arbitrage Opportunities

Arbitrage involves exploiting price discrepancies across related markets for risk-free profit.

### 2.1 Related Markets Arbitrage

**Description**: Markets within the same event where probabilities should sum to 100% (or 1.0), but don't due to market inefficiencies.

**Identification Criteria**:
- Sum of probabilities: `Σ(market.LastTradePrice) ≠ 1.0` for mutually exclusive outcomes
- Sufficient liquidity for execution
- All markets accepting orders

**Key API Fields**:
```go
event.Markets[]           // All markets in the event
market.LastTradePrice     // Latest trade price (probability)
market.Question           // Market outcome description
market.Outcomes           // Possible outcomes
market.AcceptingOrders    // Can execute trades
```

**Example Scenario**:
```
Event: "2024 Presidential Election Winner"
- Market A: "Trump wins" = 0.48
- Market B: "Biden wins" = 0.47
- Market C: "Other wins" = 0.03
Sum = 0.98 (should be 1.0)
```

**Arbitrage Strategy**: If sum < 1.0, buy all outcomes; if sum > 1.0, sell all outcomes (requires sufficient capital).

---

### 2.2 Correlated Events Arbitrage

**Description**: Logically related events across different event groups that show pricing inconsistencies.

**Identification Criteria**:
- Identify correlated events using tags or categories
- Compare pricing for logically dependent outcomes
- Look for contradictions in implied probabilities

**Key API Fields**:
```go
event.Tags[]              // Event categorization
market.Category           // Market category
market.Description        // Detailed description
event.Title               // Event title for correlation analysis
```

**Example Scenario**:
```
Event A: "Trump wins 2024 election" = 0.55
Event B: "Republican wins 2024 election" = 0.50
```
Logical issue: P(Trump wins) should be ≤ P(Republican wins)

**Search Strategy**: Use Search API with category filters to find related events.

---

### 2.3 Series Markets Arbitrage

**Description**: Pricing inconsistencies across recurring markets in a series (e.g., daily, weekly markets).

**Identification Criteria**:
- Markets in the same series with inconsistent pricing
- Compare similar outcomes across time periods
- Look for anomalies in price progression

**Key API Fields**:
```go
series.Markets[]          // All markets in series
market.SeriesSlug         // Series identifier
market.EventWeek          // Week number in series
market.Recurrence         // Recurrence type (daily, weekly)
event.Series[]            // Event series information
```

**Example Scenario**:
```
Series: "Daily Bitcoin Price Over $50k"
- Monday market: 0.65
- Tuesday market: 0.40 (anomaly)
- Wednesday market: 0.62
```

**Arbitrage Strategy**: If Tuesday's price is unjustifiably low, buy Tuesday and sell adjacent days.

---

### 2.4 NEG Risk Arbitrage

**Description**: Markets with negative risk (negRisk) enabled have special pricing rules that may create arbitrage opportunities.

**Identification Criteria**:
- `market.EnableNegRisk = true` or `event.EnableNegRisk = true`
- Understanding of negRisk mechanics
- Price discrepancies across negRisk vs. standard markets

**Key API Fields**:
```go
market.EnableNegRisk      // Market uses negative risk
event.EnableNegRisk       // Event uses negative risk
market.NegRisk           // Template negRisk setting
```

**Note**: NegRisk markets allow combined positions that reduce required collateral, creating unique arbitrage possibilities.

---

## Price Anomaly Detection

Identifying markets with unusual price movements that may indicate opportunities.

### 3.1 Rapid Price Movement

**Description**: Markets experiencing sudden, significant price changes that may represent overreaction or underreaction to new information.

**Identification Criteria**:
- `|market.OneDayPriceChange| > threshold` (e.g., 0.15 = 15%)
- Compare to `market.OneWeekPriceChange` for context
- High `market.Volume24hr` indicates active trading

**Key API Fields**:
```go
market.OneDayPriceChange    // 24-hour price change
market.OneWeekPriceChange   // 7-day price change
market.OneMonthPriceChange  // 30-day price change
market.LastTradePrice       // Current price
market.Volume24hr           // Recent volume
```

**Trading Strategy**:
- **Mean Reversion**: If price moved too far, expect pullback
- **Momentum**: If new information is significant, trend may continue

---

### 3.2 Stale Prices

**Description**: Markets where prices have not updated despite new relevant information becoming available.

**Identification Criteria**:
- Low recent volume despite significant related news
- Price hasn't changed significantly but similar markets have
- Time since last trade is unusually long

**Key API Fields**:
```go
market.Volume24hr           // Recent trading activity (low = stale)
market.LastTradePrice       // May not reflect current info
market.Active               // Market still active
market.AcceptingOrders      // Can trade
```

**Detection Strategy**: Compare market to related markets or external news sources.

---

## Statistical Arbitrage

Using statistical methods to identify mispriced markets.

### 4.1 Mean Reversion

**Description**: Markets whose prices have deviated significantly from historical averages may revert to the mean.

**Identification Criteria**:
- Track historical price data over time
- Calculate standard deviations from mean
- Identify extreme price movements

**Key API Fields**:
```go
market.OneDayPriceChange
market.OneWeekPriceChange
market.OneMonthPriceChange
market.LastTradePrice
```

**Implementation**: Requires historical data collection and statistical analysis.

---

### 4.2 Volume/Liquidity Imbalance

**Description**: Unusual ratios between trading volume and available liquidity may indicate price pressure or opportunities.

**Identification Criteria**:
- High `market.Volume24hr / market.LiquidityNum` ratio
- Sudden changes in liquidity distribution
- Imbalance between AMM and CLOB liquidity

**Key API Fields**:
```go
market.Volume24hr      // Recent volume
market.VolumeNum       // Total volume
market.VolumeAmm       // AMM volume
market.VolumeClob      // CLOB volume
market.LiquidityAmm    // AMM liquidity
market.LiquidityClob   // CLOB liquidity
```

**Analysis**:
- High volume/liquidity ratio → price pressure, potential for price movement
- AMM vs CLOB imbalance → arbitrage between venues

---

## Special Situations

Unique market conditions that create trading opportunities.

### 5.1 About to Close Markets

**Description**: Markets approaching their resolution deadline where outcomes may be clear but prices haven't fully adjusted.

**Identification Criteria**:
- `market.EndDate` within 24-48 hours
- `market.Closed = false` and `market.Active = true`
- Outcome is predictable or already determined
- Price does not reflect certain outcome

**Key API Fields**:
```go
market.EndDate          // Market close time
market.Closed           // Currently closed?
market.Active           // Currently active?
market.Ended            // Event has ended?
market.LastTradePrice   // Current price
```

**Strategy**: If outcome is clear (e.g., game already finished), trade toward certain outcome (0.0 or 1.0).

---

### 5.2 Automatic Resolution Markets

**Description**: Markets that resolve automatically based on oracle data, reducing resolution risk.

**Identification Criteria**:
- `market.AutomaticallyResolved = true`
- `event.AutomaticallyActive = true`
- Clear data source for resolution

**Key API Fields**:
```go
market.AutomaticallyResolved   // Auto-resolves based on data
event.AutomaticallyActive      // Event auto-resolves
market.ResolvedBy              // Resolution source
market.ResolutionSource        // Data source description
```

**Advantage**: Lower resolution risk means more confident trading near boundaries.

---

### 5.3 Featured vs Non-Featured Markets

**Description**: Featured markets receive more visibility and liquidity; non-featured markets may be less efficiently priced.

**Identification Criteria**:
- `event.Featured = false`
- Lower liquidity than comparable featured markets
- Similar underlying uncertainty

**Key API Fields**:
```go
event.Featured          // Market is featured on platform
market.FeaturedOrder    // Featured ranking order
market.LiquidityNum     // Total liquidity
```

**Strategy**: Non-featured markets may have wider spreads and less competition for market making.

---

### 5.4 Deploying/Pending Markets

**Description**: Markets in the process of being deployed or scheduled for deployment.

**Identification Criteria**:
- `market.PendingDeployment = true` or `market.Deploying = true`
- `market.ScheduledDeploymentTimestamp` set

**Key API Fields**:
```go
market.PendingDeployment           // Awaiting deployment
market.Deploying                    // Currently deploying
market.DeployingTimestamp           // When deployment started
market.ScheduledDeploymentTimestamp // Scheduled deployment time
```

**Strategy**: Monitor for early entry when market goes live.

---

## Implementation Examples

### Example 1: Low Liquidity, High Volume Scanner

```go
// Find markets with high volume but low liquidity
func findLowLiquidityHighVolume(markets []*Market, minVolume, maxLiquidity float64) []*Market {
    var opportunities []*Market

    for _, market := range markets {
        if market.Volume24hr > minVolume &&
           market.LiquidityClob < maxLiquidity &&
           market.AcceptingOrders {
            opportunities = append(opportunities, market)
        }
    }

    return opportunities
}
```

### Example 2: Related Markets Probability Sum Check

```go
// Check if probabilities in an event sum to ~1.0
func checkProbabilitySum(event *Event) (float64, bool) {
    sum := 0.0

    for _, market := range event.Markets {
        sum += market.LastTradePrice
    }

    // Allow 2% deviation
    isArbitrage := sum < 0.98 || sum > 1.02

    return sum, isArbitrage
}
```

### Example 3: Price Movement Alert

```go
// Find markets with significant price changes
func findRapidPriceMovement(markets []*Market, threshold float64) []*Market {
    var alerts []*Market

    for _, market := range markets {
        if math.Abs(market.OneDayPriceChange) > threshold &&
           market.Volume24hr > 1000 { // Minimum volume
            alerts = append(alerts, market)
        }
    }

    return alerts
}
```

### Example 4: Closing Soon Markets

```go
// Find markets closing within the next 24 hours
func findClosingSoon(markets []*Market) []*Market {
    var closing []*Market
    now := time.Now()

    for _, market := range markets {
        if !market.Closed &&
           !market.EndDate.IsZero() &&
           market.EndDate.Time().Sub(now) < 24*time.Hour &&
           market.EndDate.Time().After(now) {
            closing = append(closing, market)
        }
    }

    return closing
}
```

---

## Important Considerations

### Risk Management

1. **Execution Risk**: Prices may move before orders are filled
2. **Liquidity Risk**: Insufficient liquidity may prevent execution
3. **Resolution Risk**: Market may resolve unexpectedly
4. **Fee Impact**: Consider `market.MakerBaseFee` and `market.TakerBaseFee`
5. **Capital Requirements**: Many strategies require significant capital

### API Usage

1. **Rate Limiting**: Respect API rate limits
2. **Data Freshness**: Some opportunities require real-time data
3. **Historical Data**: Statistical strategies need historical tracking
4. **External Data**: May need external news/data sources

### Market Mechanics

1. **Tick Size**: `market.OrderPriceMinTickSize` affects minimum profit
2. **Minimum Order Size**: `market.OrderMinSize` affects capital needs
3. **Order Types**: May need to use CLOB API for actual trading
4. **Neg Risk**: Special collateral rules for negRisk markets

---

## Recommended Tools

To implement these strategies, consider building:

1. **Market Scanner**: Continuously scan for opportunities using Gamma API
2. **Alert System**: Notify when opportunities meet criteria
3. **Historical Tracker**: Store price/volume data for statistical analysis
4. **Position Manager**: Track positions and P&L across markets
5. **Risk Calculator**: Assess risk before entering positions

---

## API Endpoints Used

- `GET /markets` - Fetch markets with filters
- `GET /events` - Fetch events with related markets
- `GET /events/{id}` - Get specific event details
- `GET /series/{id}` - Get series market information
- `GET /public-search` - Search for related markets/events

---

## Disclaimer

This analysis is for educational purposes only. Trading prediction markets involves substantial risk. Always conduct thorough research, understand the mechanics of each market, and never risk more capital than you can afford to lose. Past opportunities do not guarantee future results.
