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
	client := polymarketgamma.NewClient(http.DefaultClient)
	ctx := context.Background()

	fmt.Println("🔍 Finding high volume, low liquidity markets (market making opportunities)...")
	fmt.Println(strings.Repeat("=", 82))

	// Configuration thresholds
	minVolume24hr := 10000.0 // Minimum $10k daily volume
	maxLiquidity := 5000.0   // Maximum $5k liquidity
	minVolumeRatio := 2.0    // Volume should be at least 2x liquidity
	targetCount := 5         // Find 5 markets
	limit := 100
	offset := 0
	maxAttempts := 10

	var opportunities []*polymarketgamma.Market
	closed := false

	fmt.Println("\n🔄 Scanning markets...")
	fmt.Printf("   Criteria: 24h Volume > $%.0f, Liquidity < $%.0f, Volume/Liquidity > %.1fx\n\n",
		minVolume24hr, maxLiquidity, minVolumeRatio)

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
			if market.Volume24hr <= 0 || market.LiquidityNum <= 0 {
				continue
			}

			// Skip markets not accepting orders
			if !market.AcceptingOrders {
				continue
			}

			// Calculate volume to liquidity ratio
			volumeRatio := market.Volume24hr / market.LiquidityNum

			// Check if market meets criteria
			if market.Volume24hr > minVolume24hr &&
				market.LiquidityNum < maxLiquidity &&
				volumeRatio > minVolumeRatio {

				opportunities = append(opportunities, market)
				fmt.Printf("   ✓ Found #%d: %s (V/L ratio: %.2fx, 24h vol: $%.0f, liq: $%.0f)\n",
					len(opportunities), truncateString(market.Question, 60),
					volumeRatio, market.Volume24hr, market.LiquidityNum)

				if len(opportunities) >= targetCount {
					break
				}
			}
		}

		offset += limit
	}

	if len(opportunities) == 0 {
		fmt.Println("\n❌ No opportunities found matching the criteria")
		return
	}

	fmt.Printf("\n✅ Found %d market making opportunities:\n\n", len(opportunities))

	// Print detailed analysis for each opportunity
	for i, market := range opportunities {
		volumeRatio := market.Volume24hr / market.LiquidityNum

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Opportunity #%d\n", i+1)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Basic Information
		fmt.Printf("📌 Market Information:\n")
		fmt.Printf("   Question:     %s\n", market.Question)
		fmt.Printf("   Market ID:    %s\n", market.ID)
		fmt.Printf("   Category:     %s\n", market.Category)
		fmt.Printf("   Status:       Active: %t, Accepting Orders: %t\n", market.Active, market.AcceptingOrders)

		// Volume & Liquidity Analysis
		fmt.Printf("\n💰 Volume & Liquidity Analysis:\n")
		fmt.Printf("   24h Volume:           $%s\n", formatNumber(market.Volume24hr))
		fmt.Printf("   Total Liquidity:      $%s\n", formatNumber(market.LiquidityNum))
		fmt.Printf("   CLOB Liquidity:       $%s\n", formatNumber(market.LiquidityClob))
		fmt.Printf("   AMM Liquidity:        $%s\n", formatNumber(market.LiquidityAmm))
		fmt.Printf("   Volume/Liquidity:     %.2fx ⚠️  (High turnover!)\n", volumeRatio)
		fmt.Printf("   Total Volume:         $%s\n", formatNumber(market.VolumeNum))

		// Spread Information
		fmt.Printf("\n📊 Spread & Pricing:\n")
		fmt.Printf("   Current Spread:       %.4f (%.2f%%)\n", market.Spread, market.Spread*100)
		fmt.Printf("   Best Bid:             %.4f\n", market.BestBid)
		fmt.Printf("   Best Ask:             %.4f\n", market.BestAsk)
		fmt.Printf("   Last Trade Price:     %.4f\n", market.LastTradePrice)
		fmt.Printf("   Min Tick Size:        %.6f\n", market.OrderPriceMinTickSize)
		if market.OrderPriceMinTickSize > 0 {
			spreadRatio := market.Spread / market.OrderPriceMinTickSize
			fmt.Printf("   Spread/Tick Ratio:    %.2fx\n", spreadRatio)
		}

		// Fee Structure
		fmt.Printf("\n💸 Fee Structure:\n")
		fmt.Printf("   Maker Fee:            %d bps\n", market.MakerBaseFee)
		fmt.Printf("   Taker Fee:            %d bps\n", market.TakerBaseFee)
		fmt.Printf("   Min Order Size:       %.2f\n", market.OrderMinSize)

		// Price Movement
		if market.OneDayPriceChange != 0 || market.OneWeekPriceChange != 0 {
			fmt.Printf("\n📈 Recent Price Changes:\n")
			if market.OneDayPriceChange != 0 {
				fmt.Printf("   1 Day:                %+.2f%%\n", market.OneDayPriceChange*100)
			}
			if market.OneWeekPriceChange != 0 {
				fmt.Printf("   1 Week:               %+.2f%%\n", market.OneWeekPriceChange*100)
			}
		}

		// Market Details
		fmt.Printf("\n📅 Market Timeline:\n")
		if !market.StartDate.IsZero() {
			fmt.Printf("   Start Date:           %s\n", market.StartDate.Format("2006-01-02 15:04:05"))
		}
		if !market.EndDate.IsZero() {
			fmt.Printf("   End Date:             %s\n", market.EndDate.Format("2006-01-02 15:04:05"))
		}

		// Link
		fmt.Printf("\n🔗 Link:\n")
		if len(market.Events) > 0 {
			fmt.Printf("   https://polymarket.com/event/%s\n", market.Events[0].Slug)
		} else {
			fmt.Printf("   https://polymarket.com/event/%s\n", market.Slug)
		}

		// Analysis
		fmt.Printf("\n💡 Market Making Analysis:\n")
		fmt.Printf("   This market shows high trading activity (24h volume: $%.0f) relative to\n", market.Volume24hr)
		fmt.Printf("   available liquidity ($%.0f), with a turnover ratio of %.2fx.\n", market.LiquidityNum, volumeRatio)
		fmt.Printf("\n   Opportunity Indicators:\n")
		fmt.Printf("   • High volume suggests active trading and consistent order flow\n")
		fmt.Printf("   • Low liquidity means wider spreads and better profit margins\n")
		fmt.Printf("   • Volume/Liquidity ratio of %.2fx indicates frequent rebalancing needs\n", volumeRatio)
		if market.Spread > 0.02 {
			fmt.Printf("   • Spread of %.2f%% provides good profit potential for market makers\n", market.Spread*100)
		}
		fmt.Printf("\n   Market Making Strategy:\n")
		fmt.Printf("   1. Place competitive bid/ask orders within the spread\n")
		fmt.Printf("   2. Adjust positions frequently due to high turnover\n")
		fmt.Printf("   3. Monitor for price changes and rebalance as needed\n")
		fmt.Printf("   4. Consider fees (maker: %dbps, taker: %dbps) in pricing\n", market.MakerBaseFee, market.TakerBaseFee)

		fmt.Println()
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
	fmt.Println("\n⚠️  Note: This analysis is for informational purposes only.")
	fmt.Println("   Always conduct your own research and risk assessment before trading.")
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// formatNumber formats a float with thousands separators
func formatNumber(n float64) string {
	s := fmt.Sprintf("%.2f", n)
	// Simple formatting - for production use a proper number formatting library
	return s
}
