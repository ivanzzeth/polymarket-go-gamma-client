package polymarketgamma

import (
	"context"
	"net/http"
	"testing"
)

func TestSearch(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Test basic search
	t.Run("BasicSearch", func(t *testing.T) {
		result, err := client.Search(ctx, &SearchParams{
			Q: "election",
		})
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}

		t.Logf("Search results for 'election':")
		t.Logf("  Events: %d", len(result.Events))
		t.Logf("  Tags: %d", len(result.Tags))
		t.Logf("  Profiles: %d", len(result.Profiles))
		t.Logf("  Pagination - HasMore: %t, TotalResults: %d", result.Pagination.HasMore, result.Pagination.TotalResults)

		// Log first few events
		if len(result.Events) > 0 {
			t.Log("\nFirst few events:")
			for i, event := range result.Events {
				if i >= 3 {
					break
				}
				t.Logf("  %d. %s (ID: %s, Slug: %s)", i+1, event.Title, event.ID, event.Slug)
			}
		}

		// Log tags
		if len(result.Tags) > 0 {
			t.Log("\nTags:")
			for i, tag := range result.Tags {
				if i >= 5 {
					break
				}
				t.Logf("  %d. %s (Event count: %d)", i+1, tag.Label, tag.EventCount)
			}
		}

		// Log profiles
		if len(result.Profiles) > 0 {
			t.Log("\nProfiles:")
			for i, profile := range result.Profiles {
				if i >= 3 {
					break
				}
				t.Logf("  %d. %s (ID: %s)", i+1, profile.Name, profile.ID)
			}
		}
	})

	// Test search with limit
	t.Run("SearchWithLimit", func(t *testing.T) {
		limit := 5
		result, err := client.Search(ctx, &SearchParams{
			Q:            "trump",
			LimitPerType: &limit,
		})
		if err != nil {
			t.Fatalf("Search with limit failed: %v", err)
		}

		t.Logf("Search results for 'trump' with limit %d:", limit)
		t.Logf("  Events: %d", len(result.Events))
		t.Logf("  Tags: %d", len(result.Tags))
		t.Logf("  Profiles: %d", len(result.Profiles))
	})

	// Test search with tags and profiles
	t.Run("SearchWithTagsAndProfiles", func(t *testing.T) {
		searchTags := true
		searchProfiles := true
		result, err := client.Search(ctx, &SearchParams{
			Q:              "crypto",
			SearchTags:     &searchTags,
			SearchProfiles: &searchProfiles,
		})
		if err != nil {
			t.Fatalf("Search with tags and profiles failed: %v", err)
		}

		t.Logf("Search results for 'crypto' with tags and profiles:")
		t.Logf("  Events: %d", len(result.Events))
		t.Logf("  Tags: %d", len(result.Tags))
		t.Logf("  Profiles: %d", len(result.Profiles))
	})

	// Test search with sorting
	t.Run("SearchWithSorting", func(t *testing.T) {
		ascending := false
		result, err := client.Search(ctx, &SearchParams{
			Q:         "sports",
			Sort:      "volume",
			Ascending: &ascending,
		})
		if err != nil {
			t.Fatalf("Search with sorting failed: %v", err)
		}

		t.Logf("Search results for 'sports' sorted by volume (descending):")
		t.Logf("  Events: %d", len(result.Events))
		if len(result.Events) > 0 {
			t.Log("\nTop events by volume:")
			for i, event := range result.Events {
				if i >= 3 {
					break
				}
				t.Logf("  %d. %s (Volume: %.2f)", i+1, event.Title, event.Volume)
			}
		}
	})

	// Test search with pagination
	t.Run("SearchWithPagination", func(t *testing.T) {
		limit := 10
		page := 1
		result, err := client.Search(ctx, &SearchParams{
			Q:            "bitcoin",
			LimitPerType: &limit,
			Page:         &page,
		})
		if err != nil {
			t.Fatalf("Search with pagination failed: %v", err)
		}

		t.Logf("Search results for 'bitcoin' (page %d, limit %d):", page, limit)
		t.Logf("  Events: %d", len(result.Events))
		t.Logf("  Pagination - HasMore: %t, TotalResults: %d", result.Pagination.HasMore, result.Pagination.TotalResults)
	})

	// Test search with event tags filter
	t.Run("SearchWithEventTags", func(t *testing.T) {
		result, err := client.Search(ctx, &SearchParams{
			Q:         "election",
			EventsTag: []string{"politics"},
		})
		if err != nil {
			t.Fatalf("Search with event tags failed: %v", err)
		}

		t.Logf("Search results for 'election' with politics tag:")
		t.Logf("  Events: %d", len(result.Events))
	})

	// Test search with exclude tag
	t.Run("SearchWithExcludeTag", func(t *testing.T) {
		result, err := client.Search(ctx, &SearchParams{
			Q:            "sports",
			ExcludeTagID: []int{1, 2},
		})
		if err != nil {
			t.Fatalf("Search with exclude tags failed: %v", err)
		}

		t.Logf("Search results for 'sports' excluding tag IDs [1, 2]:")
		t.Logf("  Events: %d", len(result.Events))
	})

	// Test search with optimized images
	t.Run("SearchWithOptimizedImages", func(t *testing.T) {
		optimized := true
		result, err := client.Search(ctx, &SearchParams{
			Q:         "nfl",
			Optimized: &optimized,
		})
		if err != nil {
			t.Fatalf("Search with optimized images failed: %v", err)
		}

		t.Logf("Search results for 'nfl' with optimized images:")
		t.Logf("  Events: %d", len(result.Events))

		// Check if optimized images are included
		if len(result.Events) > 0 && result.Events[0].ImageOptimized != nil {
			t.Log("  Optimized image data is included")
		}
	})

	// Test search with all parameters
	t.Run("SearchWithAllParams", func(t *testing.T) {
		cache := true
		limit := 5
		page := 1
		keepClosed := 1
		ascending := false
		searchTags := true
		searchProfiles := true
		optimized := true

		result, err := client.Search(ctx, &SearchParams{
			Q:                 "election 2024",
			Cache:             &cache,
			EventsStatus:      "active",
			LimitPerType:      &limit,
			Page:              &page,
			EventsTag:         []string{"politics", "usa"},
			KeepClosedMarkets: &keepClosed,
			Sort:              "volume",
			Ascending:         &ascending,
			SearchTags:        &searchTags,
			SearchProfiles:    &searchProfiles,
			Recurrence:        "daily",
			ExcludeTagID:      []int{99},
			Optimized:         &optimized,
		})
		if err != nil {
			t.Fatalf("Search with all params failed: %v", err)
		}

		t.Logf("Search results for 'election 2024' with all parameters:")
		t.Logf("  Events: %d", len(result.Events))
		t.Logf("  Tags: %d", len(result.Tags))
		t.Logf("  Profiles: %d", len(result.Profiles))
		t.Logf("  Pagination - HasMore: %t, TotalResults: %d", result.Pagination.HasMore, result.Pagination.TotalResults)
	})

	// Test error case: empty query
	t.Run("ErrorEmptyQuery", func(t *testing.T) {
		_, err := client.Search(ctx, &SearchParams{
			Q: "",
		})
		if err == nil {
			t.Error("Expected error for empty query, got nil")
		}
		t.Logf("Got expected error: %v", err)
	})

	// Test error case: nil params
	t.Run("ErrorNilParams", func(t *testing.T) {
		_, err := client.Search(ctx, nil)
		if err == nil {
			t.Error("Expected error for nil params, got nil")
		}
		t.Logf("Got expected error: %v", err)
	})
}

func TestSearchDifferentQueries(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	queries := []string{
		"trump",
		"biden",
		"bitcoin",
		"ethereum",
		"nfl",
		"nba",
		"weather",
		"stock market",
	}

	for _, query := range queries {
		t.Run("Query_"+query, func(t *testing.T) {
			limit := 3
			result, err := client.Search(ctx, &SearchParams{
				Q:            query,
				LimitPerType: &limit,
			})
			if err != nil {
				t.Errorf("Search failed for query '%s': %v", query, err)
				return
			}

			t.Logf("Query: '%s' - Events: %d, Tags: %d, Profiles: %d",
				query, len(result.Events), len(result.Tags), len(result.Profiles))

			if len(result.Events) > 0 {
				t.Logf("  First event: %s", result.Events[0].Title)
			}
		})
	}
}
