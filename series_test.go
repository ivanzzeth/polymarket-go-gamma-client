package polymarketgamma

import (
	"context"
	"net/http"
	"testing"
)

func TestAllSeriesFunctions(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Step 1: Get series
	t.Log("Step 1: Fetching series...")
	seriesList, err := client.GetSeries(ctx, &GetSeriesParams{Limit: 3})
	if err != nil {
		t.Fatalf("GetSeries failed: %v", err)
	}

	if len(seriesList) == 0 {
		t.Skip("No series available for comprehensive test")
	}

	t.Logf("Fetched %d series", len(seriesList))

	// Print details of first series from GetSeries
	if len(seriesList) > 0 {
		t.Logf("\n=== Sample Series from GetSeries ===\n%+v", seriesList[0])
	}

	// Step 2: Traverse all series and test GetSeriesByID
	for i, series := range seriesList {
		t.Logf("\n=== Testing Series %d/%d ===", i+1, len(seriesList))
		t.Logf("Series ID: %s, Title: %s, Slug: %s", series.ID, series.Title, series.Slug)

		// Test GetSeriesByID without chat
		t.Run("GetSeriesByID_"+series.ID, func(t *testing.T) {
			fetchedSeries, err := client.GetSeriesByID(ctx, series.ID, nil)
			if err != nil {
				t.Errorf("GetSeriesByID failed: %v", err)
				return
			}
			t.Logf("✓ GetSeriesByID successful: %s", fetchedSeries.Title)

			// Print detailed fields for first series
			if i == 0 {
				t.Logf("\n=== Detailed Series Fields (GetSeriesByID) ===\n%+v", fetchedSeries)
			}
		})

		// Test GetSeriesByID with chat included
		t.Run("GetSeriesByID_WithChat_"+series.ID, func(t *testing.T) {
			includeChat := true
			fetchedSeries, err := client.GetSeriesByID(ctx, series.ID, &GetSeriesByIDQueryParams{
				IncludeChat: &includeChat,
			})
			if err != nil {
				t.Errorf("GetSeriesByID with chat failed: %v", err)
				return
			}
			t.Logf("✓ GetSeriesByID with chat successful: %s (Chats: %d)", fetchedSeries.Title, len(fetchedSeries.Chats))
		})
	}

	t.Log("\n=== All tests completed ===")
}
