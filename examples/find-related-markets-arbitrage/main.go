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

	fmt.Println("🔍 Finding related markets arbitrage opportunities...")
	fmt.Println(strings.Repeat("=", 82))
	fmt.Println("\nScanning events with multiple markets where probabilities should sum to 1.0...")

	// Configuration
	targetCount := 5
	limit := 50
	offset := 0
	maxAttempts := 10
	minDeviation := 0.02   // 2% minimum deviation to consider as arbitrage
	minLiquidity := 1000.0 // Minimum liquidity for execution

	var arbitrageOpportunities []ArbitrageOpportunity
	closed := false

	fmt.Println("\n🔄 Fetching events...")

	for attempt := 0; attempt < maxAttempts && len(arbitrageOpportunities) < targetCount; attempt++ {
		params := &polymarketgamma.GetEventsParams{
			Limit:  limit,
			Offset: offset,
			Closed: &closed,
		}

		events, err := client.GetEvents(ctx, params)
		if err != nil {
			log.Fatalf("Failed to fetch events: %v", err)
		}

		if len(events) == 0 {
			fmt.Printf("   No more events to fetch (reached end)\n")
			break
		}

		fmt.Printf("   Fetched %d events (offset: %d)...\n", len(events), offset)

		// Analyze each event
		for _, event := range events {
			// Skip events with less than 2 markets
			if len(event.Markets) < 2 {
				continue
			}

			// Calculate probability sum and check for arbitrage
			opportunity := analyzeEvent(&event, minDeviation, minLiquidity)
			if opportunity != nil {
				arbitrageOpportunities = append(arbitrageOpportunities, *opportunity)
				fmt.Printf("   ✓ Found #%d: %s (deviation: %.2f%%, %d markets)\n",
					len(arbitrageOpportunities), truncateString(event.Title, 50),
					opportunity.Deviation*100, len(event.Markets))

				if len(arbitrageOpportunities) >= targetCount {
					break
				}
			}
		}

		offset += limit
	}

	if len(arbitrageOpportunities) == 0 {
		fmt.Println("\n❌ No arbitrage opportunities found")
		return
	}

	fmt.Printf("\n✅ Found %d arbitrage opportunities:\n\n", len(arbitrageOpportunities))

	// Print detailed analysis
	for i, opp := range arbitrageOpportunities {
		printArbitrageAnalysis(i+1, &opp)
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
	fmt.Println("\n⚠️  Important Notes:")
	fmt.Println("   • This analysis assumes markets are mutually exclusive and collectively exhaustive")
	fmt.Println("   • Verify market outcomes are truly mutually exclusive before trading")
	fmt.Println("   • Consider transaction fees and slippage in actual execution")
	fmt.Println("   • Markets may resolve to multiple outcomes or have special rules")
	fmt.Println("   • Always read market descriptions and resolution criteria carefully")
}

type ArbitrageOpportunity struct {
	Event          *polymarketgamma.Event
	ProbabilitySum float64
	Deviation      float64
	ArbitrageType  string // "underpriced" or "overpriced"
	TotalLiquidity float64
	ExecutableSize float64
	ExpectedReturn float64
}

func analyzeEvent(event *polymarketgamma.Event, minDeviation float64, minLiquidity float64) *ArbitrageOpportunity {
	if len(event.Markets) < 2 {
		return nil
	}

	var sum float64
	var totalLiquidity float64
	minMarketLiquidity := math.MaxFloat64
	allAcceptingOrders := true

	// Calculate sum of probabilities (last trade prices)
	for _, market := range event.Markets {
		// Skip markets with no price data
		if market.LastTradePrice <= 0 {
			return nil
		}

		// Skip closed markets
		if market.Closed {
			return nil
		}

		sum += market.LastTradePrice
		totalLiquidity += market.LiquidityNum

		if market.LiquidityNum < minMarketLiquidity {
			minMarketLiquidity = market.LiquidityNum
		}

		if !market.AcceptingOrders {
			allAcceptingOrders = false
		}
	}

	// Check if all markets accept orders
	if !allAcceptingOrders {
		return nil
	}

	// Check minimum liquidity
	if minMarketLiquidity < minLiquidity {
		return nil
	}

	// Calculate deviation from 1.0
	deviation := math.Abs(sum - 1.0)

	// Check if deviation is significant enough
	if deviation < minDeviation {
		return nil
	}

	// Determine arbitrage type
	arbType := "overpriced"
	if sum < 1.0 {
		arbType = "underpriced"
	}

	// Calculate executable size (limited by smallest market liquidity)
	executableSize := minMarketLiquidity * 0.1 // Conservative: 10% of smallest market

	// Calculate expected return (simplified)
	expectedReturn := deviation * executableSize

	return &ArbitrageOpportunity{
		Event:          event,
		ProbabilitySum: sum,
		Deviation:      deviation,
		ArbitrageType:  arbType,
		TotalLiquidity: totalLiquidity,
		ExecutableSize: executableSize,
		ExpectedReturn: expectedReturn,
	}
}

func printArbitrageAnalysis(index int, opp *ArbitrageOpportunity) {
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Arbitrage Opportunity #%d\n", index)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Event Information
	fmt.Printf("📌 Event Information:\n")
	fmt.Printf("   Title:            %s\n", opp.Event.Title)
	fmt.Printf("   Event ID:         %s\n", opp.Event.ID)
	fmt.Printf("   Markets:          %d\n", len(opp.Event.Markets))
	if opp.Event.Slug != "" {
		fmt.Printf("   URL:              https://polymarket.com/event/%s\n", opp.Event.Slug)
	}

	// NegRisk Status Check
	hasNegRisk := opp.Event.NegRisk || opp.Event.EnableNegRisk

	// Arbitrage Analysis
	fmt.Printf("\n🎯 Arbitrage Analysis:\n")
	fmt.Printf("   Probability Sum:  %.4f", opp.ProbabilitySum)
	if opp.ArbitrageType == "underpriced" {
		fmt.Printf(" (< 1.0, markets underpriced) ✓\n")
	} else {
		fmt.Printf(" (> 1.0, markets overpriced) ⚠️\n")
	}
	fmt.Printf("   Deviation:        %.2f%% from theoretical 1.0\n", opp.Deviation*100)
	fmt.Printf("   Arbitrage Type:   %s\n", opp.ArbitrageType)

	// NegRisk Information
	fmt.Printf("\n🔍 NegRisk Status:\n")
	fmt.Printf("   NegRisk Enabled:  %t", hasNegRisk)
	if hasNegRisk {
		fmt.Printf(" (Event has NegRisk support)\n")
		if opp.Event.NegRiskMarketID != "" {
			fmt.Printf("   NegRisk Market:   %s\n", opp.Event.NegRiskMarketID)
		}
		if opp.Event.NegRiskFeeBips > 0 {
			fmt.Printf("   NegRisk Fee:      %d bps (%.2f%%)\n", opp.Event.NegRiskFeeBips, float64(opp.Event.NegRiskFeeBips)/100)
		}
	} else {
		fmt.Printf("\n")
	}

	// Market Breakdown
	fmt.Printf("\n📊 Market Breakdown:\n")
	for i, market := range opp.Event.Markets {
		fmt.Printf("   %d. %s\n", i+1, truncateString(market.Question, 70))
		fmt.Printf("      Price:         %.4f (%.2f%%)\n", market.LastTradePrice, market.LastTradePrice*100)
		fmt.Printf("      Liquidity:     $%.2f\n", market.LiquidityNum)
		fmt.Printf("      Spread:        %.4f\n", market.Spread)
		fmt.Printf("      Accepting:     %t\n", market.AcceptingOrders)
		if i < len(opp.Event.Markets)-1 {
			fmt.Println()
		}
	}

	// Liquidity & Execution
	fmt.Printf("\n💰 Liquidity & Execution:\n")
	fmt.Printf("   Total Liquidity:  $%.2f\n", opp.TotalLiquidity)
	fmt.Printf("   Executable Size:  $%.2f (conservative estimate)\n", opp.ExecutableSize)
	fmt.Printf("   Expected Return:  $%.2f (before fees)\n", opp.ExpectedReturn)
	if opp.ExpectedReturn > 0 {
		roi := (opp.ExpectedReturn / opp.ExecutableSize) * 100
		fmt.Printf("   ROI:              %.2f%%\n", roi)
	}

	// Strategy
	fmt.Printf("\n💡 Arbitrage Strategy:\n")
	if opp.ArbitrageType == "underpriced" {
		fmt.Printf("   Strategy:         BUY all markets (sum < 1.0)\n")
		fmt.Printf("   Rationale:        Probabilities sum to %.4f < 1.0\n", opp.ProbabilitySum)
		fmt.Printf("   Execution:        Buy equal amounts of all %d markets\n", len(opp.Event.Markets))
		fmt.Printf("   Profit Source:    When event resolves, one market pays 1.0, you paid %.4f\n", opp.ProbabilitySum)
		fmt.Printf("   Expected Profit:  %.4f per unit (%.2f%%)\n", 1.0-opp.ProbabilitySum, opp.Deviation*100)

		// NegRisk consideration for underpriced
		if hasNegRisk {
			fmt.Printf("\n   ⚠️  NegRisk Warning:\n")
			fmt.Printf("      • This event HAS NegRisk support, but DO NOT use it for this arbitrage!\n")
			fmt.Printf("      • Why? With NegRisk you pay 1.0 collateral, but markets only cost %.4f\n", opp.ProbabilitySum)
			fmt.Printf("      • Use STANDARD markets: Pay %.4f, receive 1.0 = profit %.4f\n",
				opp.ProbabilitySum, 1.0-opp.ProbabilitySum)
			fmt.Printf("      • Using NegRisk: Pay 1.0, receive 1.0 = profit 0.0 (NO ARBITRAGE!)\n")
		}
	} else {
		fmt.Printf("   Strategy:         SELL all markets (sum > 1.0)\n")
		fmt.Printf("   Rationale:        Probabilities sum to %.4f > 1.0\n", opp.ProbabilitySum)
		fmt.Printf("   Execution:        Sell equal amounts of all %d markets\n", len(opp.Event.Markets))
		fmt.Printf("   Profit Source:    Collect %.4f, pay out 1.0 when event resolves\n", opp.ProbabilitySum)
		fmt.Printf("   Expected Profit:  %.4f per unit (%.2f%%)\n", opp.ProbabilitySum-1.0, opp.Deviation*100)
		fmt.Printf("   ⚠️  Note:         Requires sufficient capital and ability to sell/write options\n")

		// NegRisk consideration for overpriced
		if hasNegRisk {
			fmt.Printf("\n   💡 NegRisk Alternative (Capital Efficiency, NOT Arbitrage):\n")
			fmt.Printf("      • This event HAS NegRisk support - can reduce capital requirements\n")
			fmt.Printf("      • Standard approach: Buy all markets for %.4f collateral\n", opp.ProbabilitySum)
			fmt.Printf("      • NegRisk approach: Only need 1.0 collateral (saves %.4f)\n", opp.ProbabilitySum-1.0)
			fmt.Printf("      • However, this is NOT true arbitrage - no guaranteed profit\n")
			fmt.Printf("      • NegRisk advantage: Lower capital lock-up (%.2f%% reduction)\n", (opp.ProbabilitySum-1.0)/opp.ProbabilitySum*100)
			fmt.Printf("      • Use case: Holding positions with less capital, not arbitrage\n")
		}
	}

	// Risks & Considerations
	fmt.Printf("\n⚠️  Risks & Considerations:\n")
	fmt.Printf("   • Verify markets are truly mutually exclusive\n")
	fmt.Printf("   • Check if event can have multiple outcomes or 'none of the above'\n")
	fmt.Printf("   • Consider transaction fees (maker/taker fees)\n")
	if hasNegRisk && opp.Event.NegRiskFeeBips > 0 {
		fmt.Printf("   • NegRisk fee: %d bps (%.2f%%) on NegRisk transactions\n",
			opp.Event.NegRiskFeeBips, float64(opp.Event.NegRiskFeeBips)/100)
	}
	fmt.Printf("   • Account for slippage on larger orders\n")
	fmt.Printf("   • Markets may move before execution completes\n")
	if opp.ArbitrageType == "overpriced" {
		fmt.Printf("   • Selling requires margin/collateral\n")
		fmt.Printf("   • May need to hold position until resolution\n")
	}
	if hasNegRisk && opp.ArbitrageType == "underpriced" {
		fmt.Printf("   • CRITICAL: Use standard markets, NOT NegRisk for underpriced arbitrage!\n")
	}

	// Market Descriptions
	if len(opp.Event.Description) > 0 {
		fmt.Printf("\n📋 Event Description:\n")
		fmt.Printf("   %s\n", truncateString(opp.Event.Description, 200))
	}

	fmt.Println()
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
