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

	fmt.Println("🔍 Finding markets closing soon with potential mispricing...")
	fmt.Println(strings.Repeat("=", 82))
	fmt.Println("\nScanning for markets about to close where outcome may be predictable...")

	// Configuration
	hoursUntilClose := 48.0 // Find markets closing within 48 hours
	minVolume := 1000.0     // Minimum volume to filter out dead markets
	targetCount := 5        // Find 5 markets
	limit := 100
	offset := 0
	maxAttempts := 10

	var opportunities []*ClosingSoonOpportunity
	closed := false

	now := time.Now()
	fmt.Printf("\n🔄 Current time: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Searching for markets closing within %.0f hours...\n\n", hoursUntilClose)

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
			// Skip markets without end date
			if market.EndDate.IsZero() {
				continue
			}

			// Skip already closed markets
			if market.Closed || !market.Active {
				continue
			}

			// Skip low volume markets
			if market.VolumeNum < minVolume {
				continue
			}

			endTime := market.EndDate.Time()
			timeUntilClose := endTime.Sub(now)

			// Check if closing within our window
			if timeUntilClose > 0 && timeUntilClose <= time.Duration(hoursUntilClose)*time.Hour {
				hoursRemaining := timeUntilClose.Hours()

				opportunity := &ClosingSoonOpportunity{
					Market:                market,
					TimeUntilClose:        timeUntilClose,
					HoursRemaining:        hoursRemaining,
					CurrentPrice:          market.LastTradePrice,
					PotentialMispricing:   analyzePotentialMispricing(market, hoursRemaining),
					AutomaticallyResolved: market.AutomaticallyResolved,
				}

				opportunities = append(opportunities, opportunity)
				fmt.Printf("   ✓ Found #%d: %s (%.1f hours left, price: %.3f)\n",
					len(opportunities), truncateString(market.Question, 50),
					hoursRemaining, market.LastTradePrice)

				if len(opportunities) >= targetCount {
					break
				}
			}
		}

		offset += limit
	}

	if len(opportunities) == 0 {
		fmt.Println("\n❌ No markets closing soon found")
		return
	}

	fmt.Printf("\n✅ Found %d markets closing soon:\n\n", len(opportunities))

	// Print detailed analysis
	for i, opp := range opportunities {
		printClosingSoonAnalysis(i+1, opp)
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
	fmt.Println("\n⚠️  Important Notes:")
	fmt.Println("   • Always verify the outcome yourself before trading")
	fmt.Println("   • Check official sources for event results")
	fmt.Println("   • Consider that other traders may have the same information")
	fmt.Println("   • Resolution risk: markets may resolve differently than expected")
	fmt.Println("   • Time decay: ensure you have time to exit if wrong")
}

type ClosingSoonOpportunity struct {
	Market                *polymarketgamma.Market
	TimeUntilClose        time.Duration
	HoursRemaining        float64
	CurrentPrice          float64
	PotentialMispricing   string // "overpriced", "underpriced", "fairly_priced", "uncertain"
	AutomaticallyResolved bool
}

func analyzePotentialMispricing(market *polymarketgamma.Market, hoursRemaining float64) string {
	price := market.LastTradePrice

	// Markets very close to closing with extreme prices might be fairly priced
	if hoursRemaining < 6 {
		if price > 0.95 || price < 0.05 {
			return "fairly_priced"
		}
	}

	// Markets with moderate prices close to closing might have uncertainty
	if hoursRemaining < 12 && price > 0.3 && price < 0.7 {
		return "uncertain"
	}

	// Markets with extreme prices far from closing might be mispriced
	if hoursRemaining > 24 {
		if price > 0.90 {
			return "potentially_overpriced"
		}
		if price < 0.10 {
			return "potentially_underpriced"
		}
	}

	return "uncertain"
}

func printClosingSoonAnalysis(index int, opp *ClosingSoonOpportunity) {
	market := opp.Market

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Closing Soon Opportunity #%d\n", index)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Market Information
	fmt.Printf("📌 Market Information:\n")
	fmt.Printf("   Question:     %s\n", market.Question)
	fmt.Printf("   Market ID:    %s\n", market.ID)
	fmt.Printf("   Category:     %s\n", market.Category)

	// Timing Information
	fmt.Printf("\n⏰ Timing:\n")
	fmt.Printf("   End Date:             %s\n", market.EndDate.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("   Time Until Close:     ")

	hours := int(opp.HoursRemaining)
	minutes := int((opp.HoursRemaining - float64(hours)) * 60)

	if hours > 0 {
		fmt.Printf("%d hours %d minutes", hours, minutes)
	} else {
		fmt.Printf("%d minutes", minutes)
	}

	if opp.HoursRemaining < 12 {
		fmt.Printf(" ⚠️  CLOSING SOON!\n")
	} else if opp.HoursRemaining < 24 {
		fmt.Printf(" ⏳\n")
	} else {
		fmt.Printf("\n")
	}

	if !market.StartDate.IsZero() {
		fmt.Printf("   Start Date:           %s\n", market.StartDate.Format("2006-01-02 15:04:05"))
	}

	// Price & Resolution
	fmt.Printf("\n💰 Current Pricing:\n")
	fmt.Printf("   Current Price:        %.4f", opp.CurrentPrice)

	if opp.CurrentPrice > 0.95 {
		fmt.Printf(" (Very likely YES)\n")
	} else if opp.CurrentPrice > 0.80 {
		fmt.Printf(" (Likely YES)\n")
	} else if opp.CurrentPrice > 0.60 {
		fmt.Printf(" (Somewhat likely YES)\n")
	} else if opp.CurrentPrice > 0.40 {
		fmt.Printf(" (Uncertain)\n")
	} else if opp.CurrentPrice > 0.20 {
		fmt.Printf(" (Somewhat likely NO)\n")
	} else if opp.CurrentPrice > 0.05 {
		fmt.Printf(" (Likely NO)\n")
	} else {
		fmt.Printf(" (Very likely NO)\n")
	}

	fmt.Printf("   Best Bid:             %.4f\n", market.BestBid)
	fmt.Printf("   Best Ask:             %.4f\n", market.BestAsk)
	fmt.Printf("   Spread:               %.4f (%.2f%%)\n", market.Spread, market.Spread*100)

	// Resolution Info
	fmt.Printf("\n🎯 Resolution:\n")
	fmt.Printf("   Automatically Resolved: %t", opp.AutomaticallyResolved)
	if opp.AutomaticallyResolved {
		fmt.Printf(" ✅ (Lower resolution risk)\n")
	} else {
		fmt.Printf(" ⚠️  (Manual resolution)\n")
	}
	if market.ResolutionSource != "" {
		fmt.Printf("   Resolution Source:    %s\n", market.ResolutionSource)
	}
	if market.ResolvedBy != "" {
		fmt.Printf("   Resolved By:          %s\n", market.ResolvedBy)
	}

	// Volume & Liquidity
	fmt.Printf("\n📊 Trading Activity:\n")
	fmt.Printf("   Total Volume:         $%s\n", formatNumber(market.VolumeNum))
	fmt.Printf("   24h Volume:           $%s\n", formatNumber(market.Volume24hr))
	fmt.Printf("   Liquidity:            $%s\n", formatNumber(market.LiquidityNum))
	fmt.Printf("   Accepting Orders:     %t\n", market.AcceptingOrders)

	// Fees
	fmt.Printf("\n💸 Fees:\n")
	fmt.Printf("   Maker Fee:            %d bps\n", market.MakerBaseFee)
	fmt.Printf("   Taker Fee:            %d bps\n", market.TakerBaseFee)

	// Link
	fmt.Printf("\n🔗 Link:\n")
	if len(market.Events) > 0 {
		fmt.Printf("   https://polymarket.com/event/%s\n", market.Events[0].Slug)
	} else {
		fmt.Printf("   https://polymarket.com/event/%s\n", market.Slug)
	}

	// Analysis
	fmt.Printf("\n💡 Opportunity Analysis:\n")
	fmt.Printf("   Mispricing Assessment: %s\n", opp.PotentialMispricing)

	switch opp.PotentialMispricing {
	case "potentially_overpriced":
		fmt.Printf("\n   The market shows a high probability (%.2f) with %.1f hours until close.\n", opp.CurrentPrice, opp.HoursRemaining)
		fmt.Printf("   If you have information that the outcome is uncertain, this may be overpriced.\n")
		fmt.Printf("\n   Suggested Action:\n")
		fmt.Printf("   • Research if the high probability is justified\n")
		fmt.Printf("   • If outcome is less certain, consider SELLING/SHORTING\n")
		fmt.Printf("   • Verify through multiple independent sources\n")

	case "potentially_underpriced":
		fmt.Printf("\n   The market shows a low probability (%.2f) with %.1f hours until close.\n", opp.CurrentPrice, opp.HoursRemaining)
		fmt.Printf("   If you have information that the outcome is likely, this may be underpriced.\n")
		fmt.Printf("\n   Suggested Action:\n")
		fmt.Printf("   • Research if the low probability is justified\n")
		fmt.Printf("   • If outcome is more certain, consider BUYING\n")
		fmt.Printf("   • Verify through multiple independent sources\n")

	case "fairly_priced":
		fmt.Printf("\n   The market appears fairly priced given the time remaining (%.1f hours).\n", opp.HoursRemaining)
		fmt.Printf("   The extreme price (%.2f) close to resolution suggests consensus.\n", opp.CurrentPrice)
		fmt.Printf("\n   Suggested Action:\n")
		fmt.Printf("   • Only trade if you have strong contradicting evidence\n")
		fmt.Printf("   • Be very careful - the market may be correct\n")
		fmt.Printf("   • Consider the cost of being wrong\n")

	case "uncertain":
		fmt.Printf("\n   The market outcome appears genuinely uncertain (price: %.2f).\n", opp.CurrentPrice)
		fmt.Printf("   With %.1f hours remaining, significant information may still emerge.\n", opp.HoursRemaining)
		fmt.Printf("\n   Suggested Action:\n")
		fmt.Printf("   • Research the event thoroughly\n")
		fmt.Printf("   • Look for information asymmetry opportunities\n")
		fmt.Printf("   • Consider if you have better information than the market\n")
		fmt.Printf("   • Monitor for new developments\n")
	}

	// Strategy
	fmt.Printf("\n📋 Trading Checklist:\n")
	fmt.Printf("   ☐ Verify the current status/outcome of the event\n")
	fmt.Printf("   ☐ Check multiple reliable sources\n")
	fmt.Printf("   ☐ Confirm the market is still accepting orders\n")
	fmt.Printf("   ☐ Calculate potential profit minus fees\n")
	fmt.Printf("   ☐ Consider the resolution criteria carefully\n")
	fmt.Printf("   ☐ Set appropriate position size given time risk\n")
	if !opp.AutomaticallyResolved {
		fmt.Printf("   ☐ Understand manual resolution process and timeline\n")
	}

	// Risks
	fmt.Printf("\n⚠️  Specific Risks:\n")
	fmt.Printf("   • Time decay: Limited time to exit if position goes against you\n")
	fmt.Printf("   • Resolution risk: Market may resolve unexpectedly\n")
	if !opp.AutomaticallyResolved {
		fmt.Printf("   • Manual resolution: May take time and be subject to interpretation\n")
	}
	if opp.HoursRemaining < 6 {
		fmt.Printf("   • VERY SHORT TIME: Less than 6 hours to resolution!\n")
	}
	fmt.Printf("   • Information asymmetry: Other traders may know something you don't\n")
	fmt.Printf("   • Spread cost: %.2f%% spread reduces profit potential\n", market.Spread*100)

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
