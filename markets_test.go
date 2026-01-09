package polymarketgamma

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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

// TestStringOrArray_Bug_ArrayContainingJSONString tests the bug where
// an array containing a JSON string like ["[\"Up\", \"Down\"]"] is not
// properly parsed. The StringOrArray should detect that the array element
// is a JSON string and parse it, resulting in ["Up", "Down"] instead of
// ["[\"Up\", \"Down\"]"].
func TestStringOrArray_Bug_ArrayContainingJSONString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StringOrArray
		bug      bool // true if this is the bug case
	}{
		{
			name:     "normal array - should work",
			input:    `["Up", "Down"]`,
			expected: StringOrArray{"Up", "Down"},
			bug:      false,
		},
		{
			name:     "JSON string in array - BUG CASE",
			input:    `["[\"Up\", \"Down\"]"]`,
			expected: StringOrArray{"Up", "Down"}, // Should parse the JSON string
			bug:      true,
		},
		{
			name:     "JSON string directly - should work",
			input:    `"[\"Up\", \"Down\"]"`,
			expected: StringOrArray{"Up", "Down"},
			bug:      false,
		},
		{
			name:     "array with single JSON string element - BUG CASE",
			input:    `["[\"Yes\", \"No\"]"]`,
			expected: StringOrArray{"Yes", "No"}, // Should parse the JSON string
			bug:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result StringOrArray
			err := json.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}

			t.Logf("Input: %s", tt.input)
			t.Logf("Expected: %v", tt.expected)
			t.Logf("Got: %v", result)
			t.Logf("Expected length: %d", len(tt.expected))
			t.Logf("Got length: %d", len(result))

			// Check if any element looks like a JSON array string (bug indicator)
			hasJSONStringBug := false
			for i, v := range result {
				t.Logf("  [%d]: %q (len=%d)", i, v, len(v))
				if len(v) > 0 && v[0] == '[' && v[len(v)-1] == ']' {
					hasJSONStringBug = true
					t.Errorf("BUG DETECTED: Element at index %d is a JSON array string: %q", i, v)
				}
			}

			// Verify the result matches expected
			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(result), len(tt.expected))
			} else {
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("Element [%d] mismatch: got %q, want %q", i, result[i], tt.expected[i])
					}
				}
			}

			// If this is a bug case and we detected the bug, mark test as failed
			if tt.bug && hasJSONStringBug {
				t.Errorf("BUG CONFIRMED: Array containing JSON string was not properly parsed")
			}
		})
	}
}

// TestMarket_OutcomesParsing tests the actual Market struct parsing
// to verify that outcomes are correctly parsed from API responses
func TestMarket_OutcomesParsing(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected int // expected number of outcomes
	}{
		{
			name:     "normal array format",
			jsonData: `{"outcomes": ["Up", "Down"]}`,
			expected: 2,
		},
		{
			name:     "JSON string format",
			jsonData: `{"outcomes": "[\"Up\", \"Down\"]"}`,
			expected: 2,
		},
		{
			name:     "BUG: array containing JSON string",
			jsonData: `{"outcomes": ["[\"Up\", \"Down\"]"]}`,
			expected: 2, // Should parse to 2 outcomes, not 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var market struct {
				Outcomes StringOrArray `json:"outcomes"`
			}

			err := json.Unmarshal([]byte(tt.jsonData), &market)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			t.Logf("Input: %s", tt.jsonData)
			t.Logf("Outcomes: %v", market.Outcomes)
			t.Logf("Outcomes length: %d", len(market.Outcomes))
			t.Logf("Expected length: %d", tt.expected)

			// Check for bug: any outcome that looks like a JSON array string
			for i, outcome := range market.Outcomes {
				t.Logf("  Outcome[%d]: %q", i, outcome)
				if len(outcome) > 0 && outcome[0] == '[' && outcome[len(outcome)-1] == ']' {
					t.Errorf("BUG DETECTED: Outcome[%d] is a JSON array string: %q", i, outcome)
				}
			}

			if len(market.Outcomes) != tt.expected {
				t.Errorf("Outcomes count mismatch: got %d, want %d", len(market.Outcomes), tt.expected)
			}
		})
	}
}

// TestMarketByConditionID queries a market by conditionID from environment variable
// and prints all fields to help debug market status issues
func TestMarketByConditionID(t *testing.T) {
	conditionID := os.Getenv("CONDITION_ID")
	if conditionID == "" {
		t.Skip("CONDITION_ID environment variable not set, skipping test")
	}

	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	t.Logf("Querying market with conditionID: %s", conditionID)

	// Query market by conditionID
	markets, err := client.GetMarkets(ctx, &GetMarketsParams{
		ConditionIDs: []string{conditionID},
		Limit:        1,
	})
	if err != nil {
		t.Fatalf("GetMarkets failed: %v", err)
	}

	if len(markets) == 0 {
		t.Fatalf("No market found with conditionID: %s", conditionID)
	}

	market := markets[0]
	t.Logf("\n=== Market Details for ConditionID: %s ===\n", conditionID)

	// Print all fields using JSON marshaling for complete visibility
	marketJSON, err := json.MarshalIndent(market, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal market to JSON: %v", err)
	}
	t.Logf("\n=== Complete Market JSON ===\n%s\n", string(marketJSON))

	// Print key status fields explicitly
	t.Logf("\n=== Key Status Fields ===\n")
	t.Logf("ID: %s", market.ID)
	t.Logf("Question: %s", market.Question)
	t.Logf("ConditionID: %s", market.ConditionID)
	t.Logf("Slug: %s", market.Slug)
	t.Logf("Active: %v", market.Active)
	t.Logf("Closed: %v", market.Closed)
	t.Logf("Archived: %v", market.Archived)
	t.Logf("ClosedTime: %v (IsZero: %v)", market.ClosedTime, market.ClosedTime.IsZero())
	t.Logf("StartDate: %v (IsZero: %v)", market.StartDate, market.StartDate.IsZero())
	t.Logf("EndDate: %v (IsZero: %v)", market.EndDate, market.EndDate.IsZero())
	t.Logf("UMAEndDate: %v (IsZero: %v)", market.UMAEndDate, market.UMAEndDate.IsZero())
	t.Logf("Outcomes: %v", market.Outcomes)
	t.Logf("OutcomePrices: %v", market.OutcomePrices)
	t.Logf("UMAResolutionStatus: %s", market.UMAResolutionStatus)
	t.Logf("AutomaticallyResolved: %v", market.AutomaticallyResolved)
	t.Logf("ResolvedBy: %s", market.ResolvedBy)

	// Check Events if available
	if len(market.Events) > 0 {
		t.Logf("\n=== Events ===\n")
		for i, event := range market.Events {
			eventJSON, _ := json.MarshalIndent(event, "", "  ")
			t.Logf("Event %d:\n%s\n", i, string(eventJSON))
		}
	}
}
