package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	polymarketgamma "github.com/ivanzzeth/polymarket-go-gamma-client"
)

func main() {
	// Create a new client
	client := polymarketgamma.NewClient(http.DefaultClient)
	ctx := context.Background()

	fmt.Println("🔍 Finding markets with spread > 3x tick size...")
	fmt.Println(strings.Repeat("=", 62))

	// Find markets with spread > 3x tick size
	var wideSpreadMarkets []*polymarketgamma.Market
	closed := false
	limit := 100
	offset := 0
	maxAttempts := 10 // Maximum number of pages to fetch
	targetCount := 3  // We want to find 3 markets

	fmt.Println("\n🔄 Searching through markets...")

	for attempt := 0; attempt < maxAttempts && len(wideSpreadMarkets) < targetCount; attempt++ {
		params := &polymarketgamma.GetMarketsParams{
			Limit:  limit,
			Offset: offset,
			Closed: &closed,
		}

		markets, err := client.GetMarkets(ctx, params)
		if err != nil {
			log.Fatalf("Failed to fetch markets: %v", err)
		}

		// If no more markets, stop
		if len(markets) == 0 {
			fmt.Printf("   No more markets to fetch (reached end)\n")
			break
		}

		fmt.Printf("   Fetched %d markets (offset: %d)...\n", len(markets), offset)

		// Analyze markets in this batch
		for _, market := range markets {
			// Skip markets without proper tick size or spread data
			if market.OrderPriceMinTickSize <= 0 || market.Spread <= 0 {
				continue
			}

			// Skip markets with spread == 1.0
			// Note: Polymarket API sometimes returns spread=1.0 for markets that are actually closed
			// but still marked as active (API data inconsistency). These are not genuine wide spread markets.
			// A spread of 1.0 means the entire probability range (0 to 1), which typically indicates
			// no active liquidity or a closed market with stale data.
			if market.Spread >= 0.99 && market.Spread <= 1.01 {
				continue
			}

			// Check if spread > 3x tick size
			spreadRatio := market.Spread / market.OrderPriceMinTickSize
			if spreadRatio > 3.0 {
				wideSpreadMarkets = append(wideSpreadMarkets, market)
				fmt.Printf("   ✓ Found market #%d: %s (spread ratio: %.2fx)\n",
					len(wideSpreadMarkets), market.Question, spreadRatio)

				// Stop if we found enough markets
				if len(wideSpreadMarkets) >= targetCount {
					break
				}
			}
		}

		// Move to next page
		offset += limit
	}

	if len(wideSpreadMarkets) == 0 {
		fmt.Println("❌ No markets found with spread > 3x tick size")
		return
	}

	fmt.Printf("✅ Found %d market(s) with wide spreads:\n\n", len(wideSpreadMarkets))

	// Print detailed information for each market
	for i, market := range wideSpreadMarkets {
		spreadRatio := market.Spread / market.OrderPriceMinTickSize

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Market #%d\n", i+1)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Basic Information
		fmt.Printf("📌 Basic Information:\n")
		fmt.Printf("   Question:     %s\n", market.Question)
		fmt.Printf("   Market ID:    %s\n", market.ID)
		fmt.Printf("   Slug:         %s\n", market.Slug)
		fmt.Printf("   Status:       Active: %t, Closed: %t\n", market.Active, market.Closed)

		// Spread and Tick Size Information
		fmt.Printf("\n💰 Spread & Tick Information:\n")
		fmt.Printf("   Tick Size:    %.6f\n", market.OrderPriceMinTickSize)
		fmt.Printf("   Spread:       %.6f\n", market.Spread)
		fmt.Printf("   Spread Ratio: %.2fx tick size ⚠️\n", spreadRatio)

		// Price Information
		fmt.Printf("\n💵 Price Information:\n")
		fmt.Printf("   Best Bid:     %.4f\n", market.BestBid)
		fmt.Printf("   Best Ask:     %.4f\n", market.BestAsk)
		fmt.Printf("   Last Trade:   %.4f\n", market.LastTradePrice)

		// Liquidity and Volume
		fmt.Printf("\n📈 Liquidity & Volume:\n")
		fmt.Printf("   Total Liquidity:  $%.2f\n", market.LiquidityNum)
		fmt.Printf("   - AMM Liquidity:  $%.2f\n", market.LiquidityAmm)
		fmt.Printf("   - CLOB Liquidity: $%.2f\n", market.LiquidityClob)
		fmt.Printf("   Total Volume:     $%.2f\n", market.VolumeNum)
		fmt.Printf("   - AMM Volume:     $%.2f\n", market.VolumeAmm)
		fmt.Printf("   - CLOB Volume:    $%.2f\n", market.VolumeClob)
		fmt.Printf("   24h Volume:       $%.2f\n", market.Volume24hr)

		// Order Book Settings
		fmt.Printf("\n⚙️  Order Book Settings:\n")
		fmt.Printf("   Min Order Size:   %.6f\n", market.OrderMinSize)
		fmt.Printf("   Maker Base Fee:   %d bps\n", market.MakerBaseFee)
		fmt.Printf("   Taker Base Fee:   %d bps\n", market.TakerBaseFee)
		fmt.Printf("   Accepting Orders: %t\n", market.AcceptingOrders)

		// Market Details
		fmt.Printf("\n📋 Market Details:\n")
		fmt.Printf("   Market Type:      %s\n", market.MarketType)
		fmt.Printf("   Category:         %s\n", market.Category)
		fmt.Printf("   Outcomes:         %s\n", market.Outcomes)
		fmt.Printf("   Description:      %s\n", truncateString(market.Description, 100))

		// Dates
		fmt.Printf("\n📅 Important Dates:\n")
		if !market.StartDate.IsZero() {
			fmt.Printf("   Start Date:   %s\n", market.StartDate.Format("2006-01-02 15:04:05"))
		}
		if !market.EndDate.IsZero() {
			fmt.Printf("   End Date:     %s\n", market.EndDate.Format("2006-01-02 15:04:05"))
		}

		// Price Changes
		if market.OneDayPriceChange != 0 || market.OneWeekPriceChange != 0 {
			fmt.Printf("\n📊 Price Changes:\n")
			if market.OneDayPriceChange != 0 {
				fmt.Printf("   1 Day:        %+.2f%%\n", market.OneDayPriceChange*100)
			}
			if market.OneWeekPriceChange != 0 {
				fmt.Printf("   1 Week:       %+.2f%%\n", market.OneWeekPriceChange*100)
			}
			if market.OneMonthPriceChange != 0 {
				fmt.Printf("   1 Month:      %+.2f%%\n", market.OneMonthPriceChange*100)
			}
		}

		// URLs
		fmt.Printf("\n🔗 Links:\n")
		if len(market.Events) > 0 {
			fmt.Printf("   Market URL:   https://polymarket.com/event/%s\n", market.Events[0].Slug)
		} else {
			fmt.Printf("   Market URL:   https://polymarket.com/event/%s\n", market.Slug)
		}

		// Analysis
		fmt.Printf("\n💡 Analysis:\n")
		fmt.Printf("   This market has a spread of %.6f, which is %.2fx the minimum tick size.\n", market.Spread, spreadRatio)
		fmt.Printf("   This wide spread may indicate:\n")
		fmt.Printf("   - Low liquidity or market maker activity\n")
		fmt.Printf("   - High uncertainty or volatility\n")
		fmt.Printf("   - Potential opportunity for market making\n")

		fmt.Println()
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
