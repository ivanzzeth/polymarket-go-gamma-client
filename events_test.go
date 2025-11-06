package polymarketgamma

import (
	"context"
	"net/http"
	"testing"
)

func TestGetEvents(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Test getting events with nil params (all events)
	events, err := client.GetEvents(ctx, nil)
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events returned, skipping remaining tests")
	}

	t.Logf("Successfully fetched %d events", len(events))

	// Test with specific params
	featured := true
	eventsWithParams, err := client.GetEvents(ctx, &GetEventsParams{
		Limit:    5,
		Featured: &featured,
	})
	if err != nil {
		t.Fatalf("GetEvents with params failed: %v", err)
	}

	t.Logf("Successfully fetched %d featured events with limit", len(eventsWithParams))
}

func TestGetEventByID(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 10})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available to test GetEventByID")
	}

	// Test GetEventByID for each event
	for i, event := range events {
		if i >= 3 { // Only test first 3 events to save time
			break
		}

		t.Run("EventID_"+event.ID, func(t *testing.T) {
			fetchedEvent, err := client.GetEventByID(ctx, event.ID, nil)
			if err != nil {
				t.Errorf("GetEventByID failed for ID %s: %v", event.ID, err)
				return
			}

			if fetchedEvent.ID != event.ID {
				t.Errorf("Expected event ID %s, got %s", event.ID, fetchedEvent.ID)
			}

			t.Logf("Successfully fetched event by ID: %s - %s", fetchedEvent.ID, fetchedEvent.Title)
		})
	}
}

func TestGetEventByIDWithParams(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 5})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available to test GetEventByID with params")
	}

	// Test GetEventByID with query parameters
	includeChat := true
	includeTemplate := true

	for i, event := range events {
		if i >= 2 { // Only test first 2 events to save time
			break
		}

		t.Run("EventID_"+event.ID+"_WithParams", func(t *testing.T) {
			fetchedEvent, err := client.GetEventByID(ctx, event.ID, &GetEventByIDQueryParams{
				IncludeChat:     &includeChat,
				IncludeTemplate: &includeTemplate,
			})
			if err != nil {
				t.Errorf("GetEventByID with params failed for ID %s: %v", event.ID, err)
				return
			}

			if fetchedEvent.ID != event.ID {
				t.Errorf("Expected event ID %s, got %s", event.ID, fetchedEvent.ID)
			}

			t.Logf("Successfully fetched event by ID with params: %s - %s", fetchedEvent.ID, fetchedEvent.Title)

			// Log if chat or template data is included
			if len(fetchedEvent.Chats) > 0 {
				t.Logf("  Event has %d chat(s)", len(fetchedEvent.Chats))
			}
			if len(fetchedEvent.Templates) > 0 {
				t.Logf("  Event has %d template(s)", len(fetchedEvent.Templates))
			}
		})
	}
}

func TestGetEventBySlug(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 10})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available to test GetEventBySlug")
	}

	// Test GetEventBySlug for each event
	for i, event := range events {
		if i >= 3 { // Only test first 3 events to save time
			break
		}

		if event.Slug == "" {
			t.Logf("Skipping event %s - no slug available", event.ID)
			continue
		}

		t.Run("Slug_"+event.Slug, func(t *testing.T) {
			fetchedEvent, err := client.GetEventBySlug(ctx, event.Slug, nil)
			if err != nil {
				t.Errorf("GetEventBySlug failed for slug %s: %v", event.Slug, err)
				return
			}

			if fetchedEvent.Slug != event.Slug {
				t.Errorf("Expected event slug %s, got %s", event.Slug, fetchedEvent.Slug)
			}

			t.Logf("Successfully fetched event by slug: %s - %s", fetchedEvent.Slug, fetchedEvent.Title)
		})
	}
}

func TestGetEventBySlugWithParams(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events with slugs
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 10})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available to test GetEventBySlug with params")
	}

	// Find events with slugs
	var eventsWithSlugs []Event
	for _, event := range events {
		if event.Slug != "" {
			eventsWithSlugs = append(eventsWithSlugs, event)
		}
	}

	if len(eventsWithSlugs) == 0 {
		t.Skip("No events with slugs available")
	}

	// Test GetEventBySlug with query parameters
	includeChat := true
	includeTemplate := true

	for i, event := range eventsWithSlugs {
		if i >= 2 { // Only test first 2 events to save time
			break
		}

		t.Run("Slug_"+event.Slug+"_WithParams", func(t *testing.T) {
			fetchedEvent, err := client.GetEventBySlug(ctx, event.Slug, &GetEventBySlugQueryParams{
				IncludeChat:     &includeChat,
				IncludeTemplate: &includeTemplate,
			})
			if err != nil {
				t.Errorf("GetEventBySlug with params failed for slug %s: %v", event.Slug, err)
				return
			}

			if fetchedEvent.Slug != event.Slug {
				t.Errorf("Expected event slug %s, got %s", event.Slug, fetchedEvent.Slug)
			}

			t.Logf("Successfully fetched event by slug with params: %s - %s", fetchedEvent.Slug, fetchedEvent.Title)

			// Log if chat or template data is included
			if len(fetchedEvent.Chats) > 0 {
				t.Logf("  Event has %d chat(s)", len(fetchedEvent.Chats))
				for j, chat := range fetchedEvent.Chats {
					if j >= 2 { // Only log first 2 chats
						break
					}
					t.Logf("    Chat #%d: %s (Channel: %s, Live: %t)", j+1, chat.ChannelName, chat.ChannelID, chat.Live)
				}
			}
			if len(fetchedEvent.Templates) > 0 {
				t.Logf("  Event has %d template(s)", len(fetchedEvent.Templates))
				for j, tmpl := range fetchedEvent.Templates {
					if j >= 2 { // Only log first 2 templates
						break
					}
					t.Logf("    Template #%d: %s (Slug: %s)", j+1, tmpl.EventTitle, tmpl.EventSlug)
				}
			}
		})
	}
}

func TestGetEventBySlugWithIndividualParams(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events with slugs
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 5})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	// Find an event with a slug
	var testEvent *Event
	for _, event := range events {
		if event.Slug != "" {
			testEvent = &event
			break
		}
	}

	if testEvent == nil {
		t.Skip("No events with slugs available")
	}

	t.Run("OnlyIncludeChat", func(t *testing.T) {
		includeChat := true
		fetchedEvent, err := client.GetEventBySlug(ctx, testEvent.Slug, &GetEventBySlugQueryParams{
			IncludeChat: &includeChat,
		})
		if err != nil {
			t.Errorf("GetEventBySlug with include_chat failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("OnlyIncludeTemplate", func(t *testing.T) {
		includeTemplate := true
		fetchedEvent, err := client.GetEventBySlug(ctx, testEvent.Slug, &GetEventBySlugQueryParams{
			IncludeTemplate: &includeTemplate,
		})
		if err != nil {
			t.Errorf("GetEventBySlug with include_template failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("BothParams", func(t *testing.T) {
		includeChat := true
		includeTemplate := true
		fetchedEvent, err := client.GetEventBySlug(ctx, testEvent.Slug, &GetEventBySlugQueryParams{
			IncludeChat:     &includeChat,
			IncludeTemplate: &includeTemplate,
		})
		if err != nil {
			t.Errorf("GetEventBySlug with both params failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("WithFalseValues", func(t *testing.T) {
		includeChat := false
		includeTemplate := false
		fetchedEvent, err := client.GetEventBySlug(ctx, testEvent.Slug, &GetEventBySlugQueryParams{
			IncludeChat:     &includeChat,
			IncludeTemplate: &includeTemplate,
		})
		if err != nil {
			t.Errorf("GetEventBySlug with false params failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})
}

func TestGetEventByIDWithIndividualParams(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 5})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available")
	}

	testEvent := events[0]

	t.Run("OnlyIncludeChat", func(t *testing.T) {
		includeChat := true
		fetchedEvent, err := client.GetEventByID(ctx, testEvent.ID, &GetEventByIDQueryParams{
			IncludeChat: &includeChat,
		})
		if err != nil {
			t.Errorf("GetEventByID with include_chat failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("OnlyIncludeTemplate", func(t *testing.T) {
		includeTemplate := true
		fetchedEvent, err := client.GetEventByID(ctx, testEvent.ID, &GetEventByIDQueryParams{
			IncludeTemplate: &includeTemplate,
		})
		if err != nil {
			t.Errorf("GetEventByID with include_template failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("BothParams", func(t *testing.T) {
		includeChat := true
		includeTemplate := true
		fetchedEvent, err := client.GetEventByID(ctx, testEvent.ID, &GetEventByIDQueryParams{
			IncludeChat:     &includeChat,
			IncludeTemplate: &includeTemplate,
		})
		if err != nil {
			t.Errorf("GetEventByID with both params failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("WithFalseValues", func(t *testing.T) {
		includeChat := false
		includeTemplate := false
		fetchedEvent, err := client.GetEventByID(ctx, testEvent.ID, &GetEventByIDQueryParams{
			IncludeChat:     &includeChat,
			IncludeTemplate: &includeTemplate,
		})
		if err != nil {
			t.Errorf("GetEventByID with false params failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})

	t.Run("NilParams", func(t *testing.T) {
		fetchedEvent, err := client.GetEventByID(ctx, testEvent.ID, nil)
		if err != nil {
			t.Errorf("GetEventByID with nil params failed: %v", err)
			return
		}

		t.Logf("Event: %s", fetchedEvent.Title)
		t.Logf("Chats count: %d", len(fetchedEvent.Chats))
		t.Logf("Templates count: %d", len(fetchedEvent.Templates))
	})
}

func TestGetEventTags(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get some events
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 10})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available to test GetEventTags")
	}

	// Test GetEventTags for each event
	taggedEventCount := 0
	for i, event := range events {
		if i >= 5 { // Test up to 5 events
			break
		}

		t.Run("Tags_"+event.ID, func(t *testing.T) {
			tags, err := client.GetEventTags(ctx, event.ID)
			if err != nil {
				t.Errorf("GetEventTags failed for event ID %s: %v", event.ID, err)
				return
			}

			if len(tags) > 0 {
				taggedEventCount++
				t.Logf("Event %s (%s) has %d tags", event.ID, event.Title, len(tags))

				// Log first few tags
				for j, tag := range tags {
					if j >= 3 { // Only log first 3 tags
						break
					}
					t.Logf("  Tag: %s (slug: %s, created: %s)", tag.Label, tag.Slug, tag.CreatedAt.Format("2006-01-02"))
				}
			} else {
				t.Logf("Event %s (%s) has no tags", event.ID, event.Title)
			}
		})
	}

	if taggedEventCount == 0 {
		t.Log("Note: No events with tags found in the test sample")
	} else {
		t.Logf("Found %d events with tags", taggedEventCount)
	}
}

func TestAllEventsFunctions(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Step 1: Get events
	t.Log("Step 1: Fetching events...")
	events, err := client.GetEvents(ctx, &GetEventsParams{Limit: 3})
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Skip("No events available for comprehensive test")
	}

	t.Logf("Fetched %d events", len(events))

	// Print details of first event from GetEvents
	if len(events) > 0 {
		t.Logf("\n=== Sample Event from GetEvents ===\n%+v", events[0])
	}

	// Step 2: Traverse all events and test each function
	for i, event := range events {
		t.Logf("\n=== Testing Event %d/%d ===", i+1, len(events))
		t.Logf("Event ID: %s, Title: %s, Slug: %s", event.ID, event.Title, event.Slug)

		// Test GetEventByID
		t.Run("GetEventByID_"+event.ID, func(t *testing.T) {
			fetchedEvent, err := client.GetEventByID(ctx, event.ID, nil)
			if err != nil {
				t.Errorf("GetEventByID failed: %v", err)
				return
			}
			t.Logf("✓ GetEventByID successful: %s", fetchedEvent.Title)

			// Print detailed fields for first event
			if i == 0 {
				t.Logf("\n=== Detailed Event Fields (GetEventByID) ===\n%+v", fetchedEvent)
			}
		})

		// Test GetEventBySlug (if slug exists)
		if event.Slug != "" {
			t.Run("GetEventBySlug_"+event.Slug, func(t *testing.T) {
				fetchedEvent, err := client.GetEventBySlug(ctx, event.Slug, nil)
				if err != nil {
					t.Errorf("GetEventBySlug failed: %v", err)
					return
				}
				t.Logf("✓ GetEventBySlug successful: %s", fetchedEvent.Title)

				// Print detailed fields for first event
				if i == 0 {
					t.Logf("\n=== Detailed Event Fields (GetEventBySlug) ===\n%+v", fetchedEvent)
				}
			})
		} else {
			t.Logf("⊘ Skipping GetEventBySlug (no slug available)")
		}

		// Test GetEventTags
		t.Run("GetEventTags_"+event.ID, func(t *testing.T) {
			tags, err := client.GetEventTags(ctx, event.ID)
			if err != nil {
				t.Errorf("GetEventTags failed: %v", err)
				return
			}
			t.Logf("✓ GetEventTags successful: %d tags found", len(tags))

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
