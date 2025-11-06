package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	polymarketgamma "github.com/ivanzzeth/polymarket-go-gamma-client"
)

func main() {
	client := polymarketgamma.NewClient(http.DefaultClient)
	ctx := context.Background()

	fmt.Println("🔍 Finding markets with rapid price movements...")
	fmt.Println(strings.Repeat("=", 82))
	fmt.Println("\nScanning for markets with significant price changes that may present trading opportunities...")

	// Configuration thresholds
	minPriceChange := 0.15 // 15% minimum price change in 24h
	minVolume := 5000.0    // Minimum $5k volume to filter out illiquid markets
	targetCount := 5       // Find 5 markets
	limit := 100
	offset := 0
	maxAttempts := 10

	var opportunities []*PriceMovementOpportunity
	closed := false

	fmt.Printf("\n🔄 Searching markets...\n")
	fmt.Printf("   Criteria: 24h price change > %.0f%%, Volume > $%.0f\n\n", minPriceChange*100, minVolume)

	for attempt := 0; attempt < maxAttempts && len(opportunities) < targetCount; attempt++ {
		params := &polymarketgamma.GetMarketsParams{
			Limit:  limit,
			Offset: offset,
			Closed: &closed,
		}

		markets, err := client.GetMarkets(ctx, params)
		if err != nil {
			log.Fatalf("Failed to fetch markets: %v", err)
		}

		if len(markets) == 0 {
			fmt.Printf("   No more markets to fetch (reached end)\n")
			break
		}

		fmt.Printf("   Fetched %d markets (offset: %d)...\n", len(markets), offset)

		for _, market := range markets {
			// Skip markets without sufficient data
			if market.Volume24hr < minVolume {
				continue
			}

			// Skip closed markets
			if market.Closed || !market.AcceptingOrders {
				continue
			}

			// Check for significant price change
			absChange := math.Abs(market.OneDayPriceChange)
			if absChange >= minPriceChange {
				opportunity := &PriceMovementOpportunity{
					Market:         market,
					PriceChange24h: market.OneDayPriceChange,
					Direction:      "up",
				}

				if market.OneDayPriceChange < 0 {
					opportunity.Direction = "down"
				}

				// Calculate momentum vs mean reversion signals
				opportunity.MovementType = analyzePriceMovement(market)

				opportunities = append(opportunities, opportunity)
				fmt.Printf("   ✓ Found #%d: %s (%s %.1f%%, vol: $%.0f)\n",
					len(opportunities), truncateString(market.Question, 55),
					opportunity.Direction, absChange*100, market.Volume24hr)

				if len(opportunities) >= targetCount {
					break
				}
			}
		}

		offset += limit
	}

	if len(opportunities) == 0 {
		fmt.Println("\n❌ No rapid price movement opportunities found")
		return
	}

	fmt.Printf("\n✅ Found %d rapid price movement opportunities:\n\n", len(opportunities))

	// Print detailed analysis
	for i, opp := range opportunities {
		printPriceMovementAnalysis(i+1, opp)
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
	fmt.Println("\n⚠️  Important Notes:")
	fmt.Println("   • Rapid price movements can indicate new information or market overreaction")
	fmt.Println("   • Always verify the reason for price movement before trading")
	fmt.Println("   • Mean reversion strategies work best in stable conditions")
	fmt.Println("   • Momentum strategies work when new information is significant")
	fmt.Println("   • Consider market liquidity and spread before entering positions")
}

type PriceMovementOpportunity struct {
	Market         *polymarketgamma.Market
	PriceChange24h float64
	Direction      string // "up" or "down"
	MovementType   string // "momentum" or "mean_reversion" or "uncertain"
}

func analyzePriceMovement(market *polymarketgamma.Market) string {
	// Compare 24h change to 1-week change to determine if this is a new trend or continuation
	oneDayAbs := math.Abs(market.OneDayPriceChange)
	oneWeekAbs := math.Abs(market.OneWeekPriceChange)

	// If 1-day change is much larger than 1-week average, it's likely a sudden move (mean reversion candidate)
	if oneWeekAbs < oneDayAbs*0.5 {
		return "mean_reversion"
	}

	// If 1-day and 1-week changes are in the same direction and similar magnitude, it's momentum
	if market.OneDayPriceChange*market.OneWeekPriceChange > 0 {
		return "momentum"
	}

	return "uncertain"
}

func printPriceMovementAnalysis(index int, opp *PriceMovementOpportunity) {
	market := opp.Market

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Price Movement Opportunity #%d\n", index)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Market Information
	fmt.Printf("📌 Market Information:\n")
	fmt.Printf("   Question:     %s\n", market.Question)
	fmt.Printf("   Market ID:    %s\n", market.ID)
	fmt.Printf("   Category:     %s\n", market.Category)
	fmt.Printf("   Status:       Active: %t, Accepting Orders: %t\n", market.Active, market.AcceptingOrders)

	// Price Movement Analysis
	fmt.Printf("\n📊 Price Movement Analysis:\n")
	fmt.Printf("   24h Change:           %+.2f%% ", opp.PriceChange24h*100)
	if opp.Direction == "up" {
		fmt.Printf("📈\n")
	} else {
		fmt.Printf("📉\n")
	}

	if market.OneWeekPriceChange != 0 {
		fmt.Printf("   1 Week Change:        %+.2f%%\n", market.OneWeekPriceChange*100)
	}
	if market.OneMonthPriceChange != 0 {
		fmt.Printf("   1 Month Change:       %+.2f%%\n", market.OneMonthPriceChange*100)
	}

	fmt.Printf("   Current Price:        %.4f\n", market.LastTradePrice)
	fmt.Printf("   Movement Type:        %s\n", opp.MovementType)

	// Volume & Liquidity
	fmt.Printf("\n💰 Volume & Liquidity:\n")
	fmt.Printf("   24h Volume:           $%s\n", formatNumber(market.Volume24hr))
	fmt.Printf("   Total Volume:         $%s\n", formatNumber(market.VolumeNum))
	fmt.Printf("   Liquidity:            $%s\n", formatNumber(market.LiquidityNum))
	fmt.Printf("   Spread:               %.4f (%.2f%%)\n", market.Spread, market.Spread*100)

	// Order Book Info
	fmt.Printf("\n📋 Order Book:\n")
	fmt.Printf("   Best Bid:             %.4f\n", market.BestBid)
	fmt.Printf("   Best Ask:             %.4f\n", market.BestAsk)
	fmt.Printf("   Maker Fee:            %d bps\n", market.MakerBaseFee)
	fmt.Printf("   Taker Fee:            %d bps\n", market.TakerBaseFee)

	// Market Timeline
	if !market.EndDate.IsZero() {
		fmt.Printf("\n📅 Market Timeline:\n")
		fmt.Printf("   End Date:             %s\n", market.EndDate.Format("2006-01-02 15:04:05"))
	}

	// Link
	fmt.Printf("\n🔗 Link:\n")
	if len(market.Events) > 0 {
		fmt.Printf("   https://polymarket.com/event/%s\n", market.Events[0].Slug)
	} else {
		fmt.Printf("   https://polymarket.com/event/%s\n", market.Slug)
	}

	// Trading Strategy
	fmt.Printf("\n💡 Trading Strategy Suggestion:\n")

	switch opp.MovementType {
	case "mean_reversion":
		fmt.Printf("   Strategy Type:        Mean Reversion 🔄\n")
		fmt.Printf("   Rationale:            24h price change (%.1f%%) is much larger than recent trend\n", opp.PriceChange24h*100)
		fmt.Printf("   Opportunity:          Price may have overreacted to news\n")
		fmt.Printf("\n   Suggested Approach:\n")
		if opp.Direction == "up" {
			fmt.Printf("   • Consider SELLING at current elevated price (%.4f)\n", market.LastTradePrice)
			fmt.Printf("   • Price moved up %.1f%% in 24h, may pull back\n", math.Abs(opp.PriceChange24h)*100)
		} else {
			fmt.Printf("   • Consider BUYING at current depressed price (%.4f)\n", market.LastTradePrice)
			fmt.Printf("   • Price moved down %.1f%% in 24h, may recover\n", math.Abs(opp.PriceChange24h)*100)
		}
		fmt.Printf("   • Set stop-loss in case movement continues\n")
		fmt.Printf("   • Monitor for reversal signals\n")

	case "momentum":
		fmt.Printf("   Strategy Type:        Momentum Trading 🚀\n")
		fmt.Printf("   Rationale:            Price showing sustained movement over multiple periods\n")
		fmt.Printf("   Opportunity:          Trend may continue\n")
		fmt.Printf("\n   Suggested Approach:\n")
		if opp.Direction == "up" {
			fmt.Printf("   • Consider BUYING to ride the upward momentum\n")
			fmt.Printf("   • Price up %.1f%% (24h) and %.1f%% (1w) - strong trend\n",
				math.Abs(opp.PriceChange24h)*100, math.Abs(market.OneWeekPriceChange)*100)
		} else {
			fmt.Printf("   • Consider SELLING to profit from downward momentum\n")
			fmt.Printf("   • Price down %.1f%% (24h) and %.1f%% (1w) - strong trend\n",
				math.Abs(opp.PriceChange24h)*100, math.Abs(market.OneWeekPriceChange)*100)
		}
		fmt.Printf("   • Use trailing stop to lock in profits\n")
		fmt.Printf("   • Watch for trend reversal signals\n")

	case "uncertain":
		fmt.Printf("   Strategy Type:        Cautious/Research 🔍\n")
		fmt.Printf("   Rationale:            Price movement pattern is unclear\n")
		fmt.Printf("   Opportunity:          Requires further analysis\n")
		fmt.Printf("\n   Suggested Approach:\n")
		fmt.Printf("   • Research the reason for the price movement\n")
		fmt.Printf("   • Check external news sources\n")
		fmt.Printf("   • Wait for clearer signals before entering\n")
		fmt.Printf("   • Consider paper trading to test hypothesis\n")
	}

	// Risks
	fmt.Printf("\n⚠️  Risk Considerations:\n")
	fmt.Printf("   • Verify the cause of price movement (news, events, etc.)\n")
	fmt.Printf("   • Check if movement is justified by new information\n")
	fmt.Printf("   • Consider transaction costs (fees, spread)\n")
	fmt.Printf("   • Be prepared for continued volatility\n")
	if market.Spread > 0.02 {
		fmt.Printf("   • Wide spread (%.2f%%) may reduce profit potential\n", market.Spread*100)
	}

	fmt.Println()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func formatNumber(n float64) string {
	return fmt.Sprintf("%.2f", n)
}
