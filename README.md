# polymarket-go-gamma-client

Go SDK for the Polymarket Gamma RESTful API.

All market data necessary for market resolution is available on-chain (ie ancillaryData in UMA 00 request), but Polymarket also provides a hosted service, Gamma, that indexes this data and provides additional market metadata (ie categorization, indexed volume, etc). This service is made available through a REST API. For public users, this resource read only and can be used to fetch useful information about markets for things like non-profit research projects, alternative trading interfaces, automated trading systems etc.

**API Endpoint:** https://gamma-api.polymarket.com

## Features

- Complete Go SDK for Polymarket Gamma API
- Type-safe API with comprehensive market data structures
- Support for Markets, Events, Series, Search, Sports, and Tags endpoints
- Real-world examples for trading opportunity discovery

## Installation

Requires Go 1.24 or later.

```bash
go get github.com/ivanzzeth/polymarket-go-gamma-client
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    polymarketgamma "github.com/ivanzzeth/polymarket-go-gamma-client"
)

func main() {
    client := polymarketgamma.NewClient(http.DefaultClient)

    // Fetch markets
    closed := false
    params := &polymarketgamma.GetMarketsParams{
        Limit:  10,
        Closed: &closed,
    }

    markets, err := client.GetMarkets(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    for _, market := range markets {
        fmt.Printf("%s: $%.2f (24h vol: $%.2f)\n",
            market.Question, market.LastTradePrice, market.Volume24hr)
    }
}
```

## API Coverage

### Markets
- `GetMarkets()` - List markets with filtering
- `GetMarketByID()` - Get single market by ID
- `GetMarketBySlug()` - Get single market by slug

### Events
- `GetEvents()` - List events with filtering
- `GetEventByID()` - Get single event by ID
- `GetEventBySlug()` - Get single event by slug

### Series
- `GetSeries()` - List series with filtering
- `GetSeriesByID()` - Get single series by ID

### Search
- `Search()` - Search markets, events, and profiles

### Health
- `HealthCheck()` - Check API health status

### Sports
- `GetTeams()` - List sports teams with filtering
- `GetSportsMetadata()` - Get sports metadata including images and resolution sources

### Tags
- `GetTags()` - List tags with filtering
- `GetTagByID()` - Get single tag by ID
- `GetTagBySlug()` - Get single tag by slug
- `GetRelatedTagsByID()` - Get related tags by tag ID
- `GetRelatedTagsBySlug()` - Get related tags by tag slug
- `GetRelatedTagsDetailByID()` - Get detailed tag information for related tags by ID
- `GetRelatedTagsDetailBySlug()` - Get detailed tag information for related tags by slug

## Examples

This repository includes comprehensive examples demonstrating various trading opportunity detection strategies:

### 1. Market Making Opportunities

#### [Wide Spread Markets](./examples/find-wide-spread-markets/)
Find markets where the bid-ask spread is significantly larger than the minimum tick size.
- Identifies inefficient markets with wide spreads
- Great for market makers to capture spread profits
- Filters out false positives (closed markets)

```bash
cd examples/find-wide-spread-markets
go run main.go
```

#### [Low Liquidity, High Volume Markets](./examples/find-low-liquidity-high-volume/)
Discover markets with high trading activity but insufficient liquidity.
- High volume indicates strong interest
- Low liquidity means wider spreads
- Excellent opportunity for liquidity providers
- Provides detailed market making analysis

```bash
cd examples/find-low-liquidity-high-volume
go run main.go
```

#### [New Active Markets](./examples/find-new-active-markets/)
Find recently launched markets with early mover advantages.
- Markets less than 7 days old
- Low competition from other market makers
- Wide spreads and early positioning opportunities
- Includes opportunity scoring system

```bash
cd examples/find-new-active-markets
go run main.go
```

### 2. Arbitrage Opportunities

#### [Related Markets Arbitrage](./examples/find-related-markets-arbitrage/)
Identify events where probabilities don't sum to 100%.
- Exploits pricing inefficiencies across related markets
- Detects both underpriced and overpriced scenarios
- Calculates expected returns and ROI
- Provides detailed execution strategies

```bash
cd examples/find-related-markets-arbitrage
go run main.go
```

#### [NegRisk Opportunities](./examples/find-negrisk-opportunities/)
Discover capital-efficient trading using Negative Risk markets.
- Reduced collateral requirements
- NegRisk-specific arbitrage strategies
- Capital efficiency analysis
- Detailed fee impact calculations

```bash
cd examples/find-negrisk-opportunities
go run main.go
```

### 3. Price Anomaly Detection

#### [Rapid Price Movement](./examples/find-rapid-price-movement/)
Detect markets with significant 24-hour price changes.
- Identifies potential mean reversion opportunities
- Detects momentum trading signals
- Compares short-term vs long-term trends
- Suggests trading strategies based on movement type

```bash
cd examples/find-rapid-price-movement
go run main.go
```

#### [Markets Closing Soon](./examples/find-closing-soon-markets/)
Find markets approaching resolution where outcomes may be predictable.
- Markets closing within 48 hours
- Identifies potential mispricing near expiration
- Distinguishes automatically vs manually resolved markets
- Includes resolution risk assessment

```bash
cd examples/find-closing-soon-markets
go run main.go
```

## Trading Strategy Analysis

See [ANALYSIS.md](./ANALYSIS.md) for comprehensive documentation on:
- Market making strategies
- Arbitrage techniques
- Statistical analysis approaches
- Risk management considerations
- Implementation guides

## Key Data Structures

### Market
Contains comprehensive market data including:
- Pricing (bid, ask, last trade, spread)
- Volume (24h, 1w, 1m, total)
- Liquidity (CLOB, AMM, total)
- Order book settings (tick size, min size, fees)
- Price changes (1h, 1d, 1w, 1m)
- Resolution information
- Metadata and categorization

### Event
Represents trading events with:
- Multiple markets
- Series information
- Tags and categories
- Volume and liquidity aggregates
- NegRisk configuration

### Search
Unified search across:
- Markets
- Events
- Tags
- Profiles

## Best Practices

1. **Rate Limiting**: Respect API rate limits
2. **Error Handling**: Always check for errors
3. **Data Validation**: Verify market status and data quality
4. **Risk Management**: Use appropriate position sizing
5. **Fee Awareness**: Consider maker/taker fees in strategies

## Disclaimer

This SDK and examples are for educational and research purposes only. Trading prediction markets involves substantial risk. Always:
- Conduct thorough research
- Understand market mechanics
- Never risk more than you can afford to lose
- Verify all information independently
- Comply with local regulations

Past opportunities do not guarantee future results.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Resources

- [Polymarket Documentation](https://docs.polymarket.com/)
- [Gamma API Documentation](https://docs.polymarket.com/#gamma-markets-api)
- [CLOB API Documentation](https://docs.polymarket.com/#clob-api) (for actual trading)

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing examples for reference
- Review ANALYSIS.md for strategy details