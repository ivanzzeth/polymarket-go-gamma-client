package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	polymarketgamma "github.com/ivanzzeth/polymarket-go-gamma-client"
)

func main() {
	client := polymarketgamma.NewClient(http.DefaultClient)
	ctx := context.Background()

	fmt.Println("🔍 Finding newly launched active markets...")
	fmt.Println(strings.Repeat("=", 82))
	fmt.Println("\nSearching for recently created markets with early market maker opportunities...")

	// Configuration
	daysOld := 7.0          // Markets launched within last 7 days
	maxLiquidity := 10000.0 // Maximum liquidity (early stage markets)
	targetCount := 5        // Find 5 markets
	limit := 100
	offset := 0
	maxAttempts := 10

	var opportunities []*NewMarketOpportunity
	closed := false

	now := time.Now()
	cutoffDate := now.Add(-time.Duration(daysOld*24) * time.Hour)

	fmt.Printf("\n🔄 Current time: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Searching for markets created after: %s\n\n", cutoffDate.Format("2006-01-02 15:04:05"))

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
			// Skip markets without proper metadata
			if market.StartDate.IsZero() {
				continue
			}

			// Check if market is new enough
			startTime := market.StartDate.Time()
			if startTime.Before(cutoffDate) {
				continue
			}

			// Skip closed or inactive markets
			if market.Closed || !market.Active || !market.AcceptingOrders {
				continue
			}

			// Check liquidity constraint
			if market.LiquidityNum > maxLiquidity {
				continue
			}

			age := now.Sub(startTime)
			daysSinceCreation := age.Hours() / 24

			opportunity := &NewMarketOpportunity{
				Market:            market,
				Age:               age,
				DaysSinceCreation: daysSinceCreation,
				CurrentLiquidity:  market.LiquidityNum,
				CurrentVolume:     market.VolumeNum,
				OpportunityScore:  calculateOpportunityScore(market, daysSinceCreation),
			}

			opportunities = append(opportunities, opportunity)
			fmt.Printf("   ✓ Found #%d: %s (%.1f days old, liq: $%.0f)\n",
				len(opportunities), truncateString(market.Question, 50),
				daysSinceCreation, market.LiquidityNum)

			if len(opportunities) >= targetCount {
				break
			}
		}

		offset += limit
	}

	if len(opportunities) == 0 {
		fmt.Println("\n❌ No new active markets found")
		return
	}

	fmt.Printf("\n✅ Found %d new market opportunities:\n\n", len(opportunities))

	// Print detailed analysis
	for i, opp := range opportunities {
		printNewMarketAnalysis(i+1, opp)
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
	fmt.Println("\n⚠️  Important Notes:")
	fmt.Println("   • New markets often have wider spreads and less competition")
	fmt.Println("   • Early market makers can establish positions before spreads tighten")
	fmt.Println("   • Higher risk due to less historical data and price discovery")
	fmt.Println("   • Monitor for increasing trading activity as market gains attention")
	fmt.Println("   • Be prepared to adjust quotes frequently as price discovery occurs")
}

type NewMarketOpportunity struct {
	Market            *polymarketgamma.Market
	Age               time.Duration
	DaysSinceCreation float64
	CurrentLiquidity  float64
	CurrentVolume     float64
	OpportunityScore  float64 // 0-100, higher is better
}

func calculateOpportunityScore(market *polymarketgamma.Market, daysSinceCreation float64) float64 {
	score := 50.0 // Base score

	// Newer is better (up to +20 points)
	if daysSinceCreation <= 1 {
		score += 20
	} else if daysSinceCreation <= 3 {
		score += 15
	} else if daysSinceCreation <= 5 {
		score += 10
	} else {
		score += 5
	}

	// Lower liquidity means less competition (+15 points)
	if market.LiquidityNum < 1000 {
		score += 15
	} else if market.LiquidityNum < 3000 {
		score += 10
	} else if market.LiquidityNum < 5000 {
		score += 5
	}

	// Some volume is good (shows interest) (+10 points)
	if market.Volume24hr > 100 && market.Volume24hr < 5000 {
		score += 10
	} else if market.Volume24hr >= 5000 {
		score += 5
	}

	// Wide spread is good for market makers (+10 points)
	if market.Spread > 0.05 {
		score += 10
	} else if market.Spread > 0.03 {
		score += 5
	}

	// Accepting orders is required (+5 points)
	if market.AcceptingOrders {
		score += 5
	}

	// Featured markets have more visibility (-5 points, more competition)
	if market.Featured {
		score -= 5
	}

	return score
}

func printNewMarketAnalysis(index int, opp *NewMarketOpportunity) {
	market := opp.Market

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("New Market Opportunity #%d\n", index)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Market Information
	fmt.Printf("📌 Market Information:\n")
	fmt.Printf("   Question:     %s\n", market.Question)
	fmt.Printf("   Market ID:    %s\n", market.ID)
	fmt.Printf("   Category:     %s\n", market.Category)
	fmt.Printf("   Featured:     %t", market.Featured)
	if market.Featured {
		fmt.Printf(" ⭐\n")
	} else {
		fmt.Printf(" (Less competition)\n")
	}

	// Age & Timing
	fmt.Printf("\n⏰ Market Age:\n")
	fmt.Printf("   Created:              %s\n", market.StartDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Age:                  %.1f days", opp.DaysSinceCreation)
	if opp.DaysSinceCreation <= 1 {
		fmt.Printf(" 🆕 VERY NEW!\n")
	} else if opp.DaysSinceCreation <= 3 {
		fmt.Printf(" 🆕 NEW!\n")
	} else {
		fmt.Printf("\n")
	}

	if !market.EndDate.IsZero() {
		fmt.Printf("   End Date:             %s\n", market.EndDate.Format("2006-01-02 15:04:05"))
		daysUntilEnd := market.EndDate.Time().Sub(time.Now()).Hours() / 24
		fmt.Printf("   Time Until Close:     %.1f days\n", daysUntilEnd)
	}

	// Opportunity Score
	fmt.Printf("\n⭐ Opportunity Score:  %.0f/100", opp.OpportunityScore)
	if opp.OpportunityScore >= 80 {
		fmt.Printf(" 🔥 EXCELLENT!\n")
	} else if opp.OpportunityScore >= 70 {
		fmt.Printf(" ✨ VERY GOOD\n")
	} else if opp.OpportunityScore >= 60 {
		fmt.Printf(" ✓ GOOD\n")
	} else {
		fmt.Printf(" - MODERATE\n")
	}

	// Liquidity & Volume
	fmt.Printf("\n💰 Liquidity & Volume:\n")
	fmt.Printf("   Current Liquidity:    $%s", formatNumber(opp.CurrentLiquidity))
	if opp.CurrentLiquidity < 1000 {
		fmt.Printf(" 💎 (Very low - great opportunity!)\n")
	} else if opp.CurrentLiquidity < 3000 {
		fmt.Printf(" ✓ (Low - good opportunity)\n")
	} else {
		fmt.Printf("\n")
	}

	fmt.Printf("   - CLOB Liquidity:     $%s\n", formatNumber(market.LiquidityClob))
	fmt.Printf("   - AMM Liquidity:      $%s\n", formatNumber(market.LiquidityAmm))
	fmt.Printf("   Total Volume:         $%s\n", formatNumber(opp.CurrentVolume))
	fmt.Printf("   24h Volume:           $%s", formatNumber(market.Volume24hr))

	if market.Volume24hr > 1000 {
		fmt.Printf(" 📈 (Active trading!)\n")
	} else if market.Volume24hr > 100 {
		fmt.Printf(" ✓ (Some activity)\n")
	} else {
		fmt.Printf(" (Low activity)\n")
	}

	// Spread & Pricing
	fmt.Printf("\n📊 Spread & Pricing:\n")
	fmt.Printf("   Current Spread:       %.4f (%.2f%%)", market.Spread, market.Spread*100)
	if market.Spread > 0.05 {
		fmt.Printf(" 💰 (Wide - great for MM!)\n")
	} else if market.Spread > 0.03 {
		fmt.Printf(" ✓ (Good for MM)\n")
	} else {
		fmt.Printf("\n")
	}

	fmt.Printf("   Best Bid:             %.4f\n", market.BestBid)
	fmt.Printf("   Best Ask:             %.4f\n", market.BestAsk)
	fmt.Printf("   Last Trade Price:     %.4f\n", market.LastTradePrice)
	fmt.Printf("   Min Tick Size:        %.6f\n", market.OrderPriceMinTickSize)

	if market.OrderPriceMinTickSize > 0 {
		ticksInSpread := market.Spread / market.OrderPriceMinTickSize
		fmt.Printf("   Ticks in Spread:      %.0f", ticksInSpread)
		if ticksInSpread > 10 {
			fmt.Printf(" ✓ (Room for competition)\n")
		} else {
			fmt.Printf("\n")
		}
	}

	// Order Book Settings
	fmt.Printf("\n⚙️  Order Book:\n")
	fmt.Printf("   Accepting Orders:     %t", market.AcceptingOrders)
	if market.AcceptingOrders {
		fmt.Printf(" ✓\n")
	} else {
		fmt.Printf(" ❌\n")
	}
	fmt.Printf("   Min Order Size:       %.6f\n", market.OrderMinSize)
	fmt.Printf("   Maker Fee:            %d bps\n", market.MakerBaseFee)
	fmt.Printf("   Taker Fee:            %d bps\n", market.TakerBaseFee)

	// Price Movement (if any)
	if market.OneDayPriceChange != 0 {
		fmt.Printf("\n📈 Recent Price Movement:\n")
		fmt.Printf("   24h Change:           %+.2f%%\n", market.OneDayPriceChange*100)
	}

	// Link
	fmt.Printf("\n🔗 Link:\n")
	if len(market.Events) > 0 {
		fmt.Printf("   https://polymarket.com/event/%s\n", market.Events[0].Slug)
	} else {
		fmt.Printf("   https://polymarket.com/event/%s\n", market.Slug)
	}

	// Market Making Strategy
	fmt.Printf("\n💡 Market Making Strategy:\n")
	fmt.Printf("   Early Mover Advantage:\n")
	fmt.Printf("   • This market is only %.1f days old with limited liquidity\n", opp.DaysSinceCreation)
	fmt.Printf("   • Competition is likely limited at this stage\n")
	fmt.Printf("   • Wide spreads (%.2f%%) provide good profit margins\n", market.Spread*100)

	fmt.Printf("\n   Recommended Approach:\n")
	fmt.Printf("   1. Start with conservative position sizing\n")
	fmt.Printf("   2. Place competitive quotes within the spread\n")
	fmt.Printf("   3. Monitor for other market makers entering\n")
	fmt.Printf("   4. Tighten spreads gradually as liquidity increases\n")
	fmt.Printf("   5. Be prepared for higher volatility during price discovery\n")

	fmt.Printf("\n   Position Management:\n")
	fmt.Printf("   • Initial capital: Start with 5-10%% of current liquidity\n")
	fmt.Printf("   • Quote size: $%.0f - $%.0f per side\n", opp.CurrentLiquidity*0.05, opp.CurrentLiquidity*0.10)
	fmt.Printf("   • Spread target: %.2f%% - %.2f%%\n", market.Spread*0.6*100, market.Spread*0.8*100)
	fmt.Printf("   • Rebalance frequency: Every 30-60 minutes initially\n")

	// Early Stage Indicators
	fmt.Printf("\n📊 Early Stage Indicators:\n")

	liquidityToVolume := 0.0
	if market.VolumeNum > 0 {
		liquidityToVolume = market.LiquidityNum / market.VolumeNum
	}

	if market.Volume24hr > 0 {
		fmt.Printf("   ✓ Has 24h trading volume ($%.0f)\n", market.Volume24hr)
	} else {
		fmt.Printf("   ⚠️  No 24h volume yet\n")
	}

	if opp.CurrentLiquidity < 3000 {
		fmt.Printf("   ✓ Low liquidity ($%.0f) - less competition\n", opp.CurrentLiquidity)
	}

	if market.Spread > 0.03 {
		fmt.Printf("   ✓ Wide spread (%.2f%%) - good profit potential\n", market.Spread*100)
	}

	if liquidityToVolume > 0 && liquidityToVolume < 1 {
		fmt.Printf("   ✓ High turnover ratio (%.2fx) - active market\n", 1/liquidityToVolume)
	}

	// Risks
	fmt.Printf("\n⚠️  Early Market Risks:\n")
	fmt.Printf("   • Price discovery: True fair value is still being established\n")
	fmt.Printf("   • Low volume: May have difficulty executing large orders\n")
	fmt.Printf("   • Information risk: Less public information and analysis available\n")
	fmt.Printf("   • Volatility: Expect wider price swings than established markets\n")
	fmt.Printf("   • Market may not gain traction and remain illiquid\n")
	if !market.Featured {
		fmt.Printf("   • Non-featured: May receive less visibility and volume\n")
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
