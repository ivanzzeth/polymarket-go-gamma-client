package polymarketgamma

import (
	"context"
	"net/http"
	"testing"
)

func TestAllMarketsFunctions(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Step 1: Get markets
	t.Log("Step 1: Fetching markets...")
	markets, err := client.GetMarkets(ctx, &GetMarketsParams{Limit: 3})
	if err != nil {
		t.Fatalf("GetMarkets failed: %v", err)
	}

	if len(markets) == 0 {
		t.Skip("No markets available for comprehensive test")
	}

	t.Logf("Fetched %d markets", len(markets))

	// Print details of first market from GetMarkets
	if len(markets) > 0 {
		t.Logf("\n=== Sample Market from GetMarkets ===\n%+v", *markets[0])
	}

	// Step 2: Traverse all markets and test each function
	for i, market := range markets {
		t.Logf("\n=== Testing Market %d/%d ===", i+1, len(markets))
		t.Logf("Market ID: %s, Question: %s, Slug: %s", market.ID, market.Question, market.Slug)

		// Test GetMarketByID without tags
		t.Run("GetMarketByID_"+market.ID, func(t *testing.T) {
			fetchedMarket, err := client.GetMarketByID(ctx, market.ID, nil)
			if err != nil {
				t.Errorf("GetMarketByID failed: %v", err)
				return
			}
			t.Logf("✓ GetMarketByID successful: %s", fetchedMarket.Question)

			// Print detailed fields for first market
			if i == 0 {
				t.Logf("\n=== Detailed Market Fields (GetMarketByID) ===\n%+v", fetchedMarket)
			}
		})

		// Test GetMarketByID with tags included
		t.Run("GetMarketByID_WithTags_"+market.ID, func(t *testing.T) {
			includeTag := true
			fetchedMarket, err := client.GetMarketByID(ctx, market.ID, &GetMarketByIDQueryParams{
				IncludeTag: &includeTag,
			})
			if err != nil {
				t.Errorf("GetMarketByID with tags failed: %v", err)
				return
			}
			t.Logf("✓ GetMarketByID with tags successful: %s (Tags: %d)", fetchedMarket.Question, len(fetchedMarket.Tags))
		})

		// Test GetMarketBySlug (if slug exists)
		if market.Slug != "" {
			t.Run("GetMarketBySlug_"+market.Slug, func(t *testing.T) {
				fetchedMarket, err := client.GetMarketBySlug(ctx, market.Slug, nil)
				if err != nil {
					t.Errorf("GetMarketBySlug failed: %v", err)
					return
				}
				t.Logf("✓ GetMarketBySlug successful: %s", fetchedMarket.Question)

				// Print detailed fields for first market
				if i == 0 {
					t.Logf("\n=== Detailed Market Fields (GetMarketBySlug) ===\n%+v", fetchedMarket)
				}
			})

			// Test GetMarketBySlug with tags included
			t.Run("GetMarketBySlug_WithTags_"+market.Slug, func(t *testing.T) {
				includeTag := true
				fetchedMarket, err := client.GetMarketBySlug(ctx, market.Slug, &GetMarketByIDQueryParams{
					IncludeTag: &includeTag,
				})
				if err != nil {
					t.Errorf("GetMarketBySlug with tags failed: %v", err)
					return
				}
				t.Logf("✓ GetMarketBySlug with tags successful: %s (Tags: %d)", fetchedMarket.Question, len(fetchedMarket.Tags))
			})
		} else {
			t.Logf("⊘ Skipping GetMarketBySlug (no slug available)")
		}

		// Test GetMarketTags
		t.Run("GetMarketTags_"+market.ID, func(t *testing.T) {
			tags, err := client.GetMarketTags(ctx, market.ID)
			if err != nil {
				t.Errorf("GetMarketTags failed: %v", err)
				return
			}
			t.Logf("✓ GetMarketTags successful: %d tags found", len(tags))

			if len(tags) > 0 {
				t.Log("\n=== Tag Details ===")
				for j, tag := range tags {
					if j >= 3 { // Print up to 3 tags
						t.Logf("  ... and %d more tags", len(tags)-3)
						break
					}
					t.Logf("\nTag #%d: %+v", j+1, tag)
				}
			}
		})
	}

	t.Log("\n=== All tests completed ===")
}
