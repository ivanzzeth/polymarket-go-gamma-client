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

	fmt.Println("🔍 Finding NegRisk (Negative Risk) market opportunities...")
	fmt.Println(strings.Repeat("=", 82))
	fmt.Println("\nNegRisk markets allow combined positions with reduced collateral requirements.")
	fmt.Println("This creates unique arbitrage and hedging opportunities.")

	// Configuration
	minLiquidity := 1000.0 // Minimum liquidity
	targetCount := 5       // Find 5 markets
	limit := 50
	offset := 0
	maxAttempts := 10

	var opportunities []*NegRiskOpportunity
	closed := false

	fmt.Println("\n🔄 Searching for negRisk events...")

	for attempt := 0; attempt < maxAttempts && len(opportunities) < targetCount; attempt++ {
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

		for _, event := range events {
			// Look for events with negRisk enabled
			if !event.NegRisk && !event.EnableNegRisk {
				continue
			}

			// Skip events with less than 2 markets
			if len(event.Markets) < 2 {
				continue
			}

			// Check if markets have sufficient liquidity
			totalLiquidity := 0.0
			allAcceptingOrders := true
			negRiskMarketCount := 0

			for _, market := range event.Markets {
				totalLiquidity += market.LiquidityNum

				if market.Closed || !market.AcceptingOrders {
					allAcceptingOrders = false
				}

				// Count markets with negRisk
				if market.NegRiskOther {
					negRiskMarketCount++
				}
			}

			if !allAcceptingOrders || totalLiquidity < minLiquidity {
				continue
			}

			opportunity := &NegRiskOpportunity{
				Event:              &event,
				TotalLiquidity:     totalLiquidity,
				MarketCount:        len(event.Markets),
				NegRiskMarketCount: negRiskMarketCount,
				NegRiskFeeBips:     event.NegRiskFeeBips,
			}

			// Calculate probability sum for arbitrage check
			opportunity.ProbabilitySum = calculateProbabilitySum(&event)

			opportunities = append(opportunities, opportunity)
			fmt.Printf("   ✓ Found #%d: %s (%d markets, fee: %d bps)\n",
				len(opportunities), truncateString(event.Title, 50),
				len(event.Markets), event.NegRiskFeeBips)

			if len(opportunities) >= targetCount {
				break
			}
		}

		offset += limit
	}

	if len(opportunities) == 0 {
		fmt.Println("\n❌ No negRisk opportunities found")
		return
	}

	fmt.Printf("\n✅ Found %d negRisk opportunities:\n\n", len(opportunities))

	// Print detailed analysis
	for i, opp := range opportunities {
		printNegRiskAnalysis(i+1, opp)
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("✨ Analysis complete!")
	fmt.Println("\n📚 Understanding NegRisk:")
	fmt.Println("   NegRisk (Negative Risk) markets allow traders to take combined positions")
	fmt.Println("   across mutually exclusive outcomes with reduced collateral requirements.")
	fmt.Println("   For example, if you buy 'Team A wins' and 'Team B wins' in the same event,")
	fmt.Println("   you only need to post 1 unit of collateral instead of 2, since only one can win.")
	fmt.Println("\n⚠️  Important Notes:")
	fmt.Println("   • NegRisk mechanics require understanding of collateral optimization")
	fmt.Println("   • Special fees may apply (negRiskFeeBips)")
	fmt.Println("   • Requires CLOB API for actual trading")
	fmt.Println("   • Always verify collateral requirements before trading")
}

type NegRiskOpportunity struct {
	Event              *polymarketgamma.Event
	TotalLiquidity     float64
	MarketCount        int
	NegRiskMarketCount int
	NegRiskFeeBips     int
	ProbabilitySum     float64
}

func calculateProbabilitySum(event *polymarketgamma.Event) float64 {
	sum := 0.0
	for _, market := range event.Markets {
		if market.LastTradePrice > 0 {
			sum += market.LastTradePrice
		}
	}
	return sum
}

func printNegRiskAnalysis(index int, opp *NegRiskOpportunity) {
	event := opp.Event

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("NegRisk Opportunity #%d\n", index)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Event Information
	fmt.Printf("📌 Event Information:\n")
	fmt.Printf("   Title:                %s\n", event.Title)
	fmt.Printf("   Event ID:             %s\n", event.ID)
	fmt.Printf("   Markets:              %d\n", opp.MarketCount)
	fmt.Printf("   NegRisk Enabled:      %t ✓\n", event.NegRisk || event.EnableNegRisk)
	fmt.Printf("   NegRisk Market ID:    %s\n", event.NegRiskMarketID)
	if event.Slug != "" {
		fmt.Printf("   URL:                  https://polymarket.com/event/%s\n", event.Slug)
	}

	// NegRisk Configuration
	fmt.Printf("\n⚙️  NegRisk Configuration:\n")
	fmt.Printf("   NegRisk Fee:          %d bps", opp.NegRiskFeeBips)
	if opp.NegRiskFeeBips > 0 {
		fmt.Printf(" (%.2f%%)\n", float64(opp.NegRiskFeeBips)/100)
	} else {
		fmt.Printf(" (No additional fee)\n")
	}
	fmt.Printf("   Markets with NegRisk: %d/%d\n", opp.NegRiskMarketCount, opp.MarketCount)

	// Liquidity
	fmt.Printf("\n💰 Liquidity:\n")
	fmt.Printf("   Total Liquidity:      $%s\n", formatNumber(opp.TotalLiquidity))
	fmt.Printf("   Event 24h Volume:     $%s\n", formatNumber(event.Volume24hr))

	// Probability Analysis
	fmt.Printf("\n📊 Probability Analysis:\n")
	fmt.Printf("   Probability Sum:      %.4f", opp.ProbabilitySum)

	deviation := 0.0
	if opp.ProbabilitySum > 1.0 {
		deviation = opp.ProbabilitySum - 1.0
		fmt.Printf(" (%.2f%% > 1.0) ⚠️  OVERPRICED\n", deviation*100)
	} else if opp.ProbabilitySum < 1.0 {
		deviation = 1.0 - opp.ProbabilitySum
		fmt.Printf(" (%.2f%% < 1.0) ✓ UNDERPRICED\n", deviation*100)
	} else {
		fmt.Printf(" (Fair)\n")
	}

	// Market Breakdown
	fmt.Printf("\n📋 Market Breakdown:\n")
	for i, market := range event.Markets {
		fmt.Printf("   %d. %s\n", i+1, truncateString(market.Question, 65))
		fmt.Printf("      Price:         %.4f (%.2f%% implied probability)\n", market.LastTradePrice, market.LastTradePrice*100)
		fmt.Printf("      Liquidity:     $%s\n", formatNumber(market.LiquidityNum))
		fmt.Printf("      Spread:        %.4f (%.2f%%)\n", market.Spread, market.Spread*100)
		fmt.Printf("      24h Volume:    $%s\n", formatNumber(market.Volume24hr))
		fmt.Printf("      NegRisk:       %t\n", market.NegRiskOther)
		if i < len(event.Markets)-1 {
			fmt.Println()
		}
	}

	// NegRisk Opportunities
	fmt.Printf("\n💡 NegRisk Trading Strategies:\n")

	// Determine if NegRisk has advantage
	isNegRiskAdvantage := opp.ProbabilitySum > 1.0

	// Strategy 1: Reduced Collateral Portfolio
	fmt.Printf("\n   1️⃣  Reduced Collateral Portfolio Strategy:\n")
	if isNegRiskAdvantage {
		// NegRisk is better when sum > 1.0
		savings := opp.ProbabilitySum - 1.0
		savingsPercent := (savings / opp.ProbabilitySum) * 100
		fmt.Printf("      Concept:      Buy all outcomes with NegRisk optimization\n")
		fmt.Printf("      Standard cost:    $%.4f (sum of prices)\n", opp.ProbabilitySum)
		fmt.Printf("      NegRisk cost:     $1.0000 (reduced collateral)\n")
		fmt.Printf("      ✓ Savings:        $%.4f (%.2f%% reduction) 💰\n", savings, savingsPercent)
		fmt.Printf("\n      How it works:\n")
		fmt.Printf("      • Standard markets: costs %.4f to buy all outcomes\n", opp.ProbabilitySum)
		fmt.Printf("      • NegRisk markets: only need 1.0 collateral\n")
		fmt.Printf("      • NegRisk is %.0f%% more capital efficient ✓\n", (savings/1.0)*100)
	} else {
		// Standard markets are better when sum < 1.0
		extraCost := 1.0 - opp.ProbabilitySum
		extraCostPercent := (extraCost / 1.0) * 100
		fmt.Printf("      ⚠️  Analysis:     NegRisk is NOT advantageous here!\n")
		fmt.Printf("      Standard cost:    $%.4f (sum of prices) ✓ BETTER\n", opp.ProbabilitySum)
		fmt.Printf("      NegRisk cost:     $1.0000 (required collateral)\n")
		fmt.Printf("      ❌ Extra cost:    $%.4f (%.2f%% more expensive)\n", extraCost, extraCostPercent)
		fmt.Printf("\n      Why standard is better:\n")
		fmt.Printf("      • Standard markets: costs only %.4f\n", opp.ProbabilitySum)
		fmt.Printf("      • NegRisk markets: requires 1.0 collateral\n")
		fmt.Printf("      • Standard saves $%.4f (%.0f%% cheaper)\n", extraCost, extraCostPercent)
	}

	// Strategy 2: Arbitrage
	if opp.ProbabilitySum < 0.98 {
		fmt.Printf("\n   2️⃣  Arbitrage Strategy (Use STANDARD Markets):\n")
		fmt.Printf("      ⚠️  IMPORTANT:   Do NOT use negRisk for this arbitrage!\n")
		fmt.Printf("      Opportunity:  Probabilities sum to %.4f (< 1.0)\n", opp.ProbabilitySum)
		fmt.Printf("      Strategy:     Buy all outcomes in STANDARD markets\n")
		fmt.Printf("      Cost:         $%.4f (standard) vs $1.00 (negRisk)\n", opp.ProbabilitySum)
		fmt.Printf("      Payout:       $1.00 guaranteed when event resolves\n")
		fmt.Printf("      Profit:       $%.4f (%.2f%% return)\n", 1.0-opp.ProbabilitySum, (1.0-opp.ProbabilitySum)*100)
		fmt.Printf("\n      Why not negRisk?\n")
		fmt.Printf("      • NegRisk requires $1.00 collateral = NO profit\n")
		fmt.Printf("      • Standard costs $%.4f = $%.4f profit ✓\n", opp.ProbabilitySum, 1.0-opp.ProbabilitySum)
		fmt.Printf("      • Standard is clearly superior for this trade\n")
	} else if opp.ProbabilitySum > 1.02 {
		fmt.Printf("\n   2️⃣  NegRisk Advantage (Overpriced Markets):\n")
		fmt.Printf("      Opportunity:  Probabilities sum to %.4f (> 1.0)\n", opp.ProbabilitySum)
		fmt.Printf("      Strategy:     Use NegRisk to reduce capital requirements\n")
		fmt.Printf("      Standard cost:    $%.4f (would lose money)\n", opp.ProbabilitySum)
		fmt.Printf("      NegRisk cost:     $1.00 (capital efficient)\n")
		fmt.Printf("      Benefit:      Save $%.4f in capital requirements\n", opp.ProbabilitySum-1.0)
		fmt.Printf("\n      Note: This is not an arbitrage (sum > 1.0)\n")
		fmt.Printf("      But negRisk allows you to hold all positions with less capital\n")
	}

	// Strategy 3: Hedging
	fmt.Printf("\n   3️⃣  Capital-Efficient Hedging:\n")
	fmt.Printf("      Use case:     Hedge positions across mutually exclusive outcomes\n")
	if isNegRiskAdvantage {
		fmt.Printf("      Advantage:    Reduce capital tied up in hedges ✓\n")
	} else {
		fmt.Printf("      Note:         Limited advantage when sum < 1.0\n")
	}
	fmt.Printf("      Example:      Hold position in Market A, hedge with Market B\n")
	fmt.Printf("                    Only post net exposure as collateral\n")

	// Capital Requirements Comparison
	fmt.Printf("\n💵 Capital Requirements Comparison:\n")
	fmt.Printf("   To buy 1 share of each market:\n")
	fmt.Printf("   • Standard markets:    $%.4f (sum of prices)", opp.ProbabilitySum)
	if !isNegRiskAdvantage {
		fmt.Printf(" ✓ CHEAPER\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("   • NegRisk markets:     $1.0000 (required collateral)")
	if isNegRiskAdvantage {
		fmt.Printf(" ✓ CHEAPER\n")
		fmt.Printf("   • Savings with NegRisk: $%.4f (%.1f%% reduction)\n", opp.ProbabilitySum-1.0, (opp.ProbabilitySum-1.0)/opp.ProbabilitySum*100)
	} else {
		fmt.Printf("\n")
		fmt.Printf("   • Extra cost w/ NegRisk: $%.4f (%.1f%% more expensive)\n", 1.0-opp.ProbabilitySum, (1.0-opp.ProbabilitySum)/1.0*100)
		fmt.Printf("\n   💡 Recommendation: Use STANDARD markets for this event\n")
	}

	// Fee Impact
	if opp.NegRiskFeeBips > 0 {
		fmt.Printf("\n💸 Fee Considerations:\n")
		fmt.Printf("   NegRisk Fee Impact:   %.4f per $1 trade\n", float64(opp.NegRiskFeeBips)/10000)
		fmt.Printf("   Break-even:           Need > %.4f mispricing to profit\n", float64(opp.NegRiskFeeBips)/10000)
	}

	// Risks
	fmt.Printf("\n⚠️  NegRisk-Specific Risks:\n")
	fmt.Printf("   • Complexity: Requires understanding of collateral mechanics\n")
	fmt.Printf("   • Execution: Need to trade on CLOB with negRisk support\n")
	fmt.Printf("   • Slippage: Multiple trades required may face slippage\n")
	fmt.Printf("   • Resolution: All markets must resolve correctly\n")
	if opp.NegRiskFeeBips > 0 {
		fmt.Printf("   • Fees: Additional negRisk fees (%.2f%%) reduce profit\n", float64(opp.NegRiskFeeBips)/100)
	}
	fmt.Printf("   • Liquidity: Need sufficient liquidity across all markets\n")

	// Prerequisites
	fmt.Printf("\n📋 Trading Prerequisites:\n")
	fmt.Printf("   ✓ Understanding of negRisk collateral rules\n")
	fmt.Printf("   ✓ Access to CLOB API for negRisk trading\n")
	fmt.Printf("   ✓ Sufficient capital (minimum $1.00 per combined position)\n")
	fmt.Printf("   ✓ Ability to execute across multiple markets simultaneously\n")
	fmt.Printf("   ✓ Risk management for partial fills\n")

	// Example Calculation
	if len(event.Markets) >= 2 {
		fmt.Printf("\n🧮 Example Calculation:\n")
		fmt.Printf("   Scenario: Buy 100 shares of each outcome\n")

		standardCost := opp.ProbabilitySum * 100
		negRiskCost := 100.0
		payout := 100.0

		if isNegRiskAdvantage {
			// NegRisk is better (sum > 1.0)
			fmt.Printf("\n   Standard Market Approach:\n")
			fmt.Printf("   • Capital needed:     $%.2f (100 × %.4f)\n", standardCost, opp.ProbabilitySum)
			fmt.Printf("   • Payout:             $%.2f (one outcome wins)\n", payout)
			if standardCost > payout {
				fmt.Printf("   • Loss:               $%.2f ❌\n", standardCost-payout)
			} else {
				fmt.Printf("   • Profit:             $%.2f\n", payout-standardCost)
			}

			fmt.Printf("\n   NegRisk Approach (RECOMMENDED):\n")
			fmt.Printf("   • Capital needed:     $%.2f (negRisk optimization) ✓\n", negRiskCost)
			fmt.Printf("   • Payout:             $%.2f (one outcome wins)\n", payout)
			negRiskProfit := payout - negRiskCost
			if opp.NegRiskFeeBips > 0 {
				fees := negRiskCost * float64(opp.NegRiskFeeBips) / 10000
				fmt.Printf("   • Gross profit:       $%.2f\n", negRiskProfit)
				fmt.Printf("   • NegRisk fees:       $%.2f\n", fees)
				fmt.Printf("   • Net profit:         $%.2f\n", negRiskProfit-fees)
			} else {
				fmt.Printf("   • Profit/Loss:        $%.2f\n", negRiskProfit)
			}
			savings := standardCost - negRiskCost
			fmt.Printf("\n   💰 NegRisk Advantage: Saves $%.2f (%.1f%% better)\n", savings, (savings/standardCost)*100)

		} else {
			// Standard is better (sum < 1.0)
			fmt.Printf("\n   Standard Market Approach (RECOMMENDED):\n")
			fmt.Printf("   • Capital needed:     $%.2f (100 × %.4f) ✓\n", standardCost, opp.ProbabilitySum)
			fmt.Printf("   • Payout:             $%.2f (one outcome wins)\n", payout)
			fmt.Printf("   • Profit:             $%.2f 💰\n", payout-standardCost)

			fmt.Printf("\n   NegRisk Approach (NOT recommended):\n")
			fmt.Printf("   • Capital needed:     $%.2f (required collateral)\n", negRiskCost)
			fmt.Printf("   • Payout:             $%.2f (one outcome wins)\n", payout)
			negRiskProfit := payout - negRiskCost
			if opp.NegRiskFeeBips > 0 {
				fees := negRiskCost * float64(opp.NegRiskFeeBips) / 10000
				fmt.Printf("   • Gross profit:       $%.2f\n", negRiskProfit)
				fmt.Printf("   • NegRisk fees:       $%.2f\n", fees)
				fmt.Printf("   • Net profit:         $%.2f\n", negRiskProfit-fees)
			} else {
				fmt.Printf("   • Profit/Loss:        $%.2f\n", negRiskProfit)
			}
			extraCost := negRiskCost - standardCost
			fmt.Printf("\n   ❌ NegRisk Disadvantage: Costs $%.2f MORE (%.1f%% worse)\n", extraCost, (extraCost/negRiskCost)*100)
			fmt.Printf("      You LOSE $%.2f in potential profit by using negRisk!\n", extraCost)
		}
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
