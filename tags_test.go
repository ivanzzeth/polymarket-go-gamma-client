package polymarketgamma

import (
	"context"
	"net/http"
	"testing"
)

func TestGetTags(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Test fetching tags without parameters
	t.Run("FetchTagsNoParams", func(t *testing.T) {
		tags, err := client.GetTags(ctx, nil)
		if err != nil {
			t.Fatalf("Failed to fetch tags: %v", err)
		}

		if len(tags) == 0 {
			t.Log("No tags found (this might be expected if the API returns empty data)")
		} else {
			t.Logf("Successfully fetched %d tags", len(tags))
			// Print first few tags with full details
			for i := 0; i < len(tags) && i < 3; i++ {
				tag := tags[i]
				t.Logf("Tag %d: %+v", i+1, tag)
			}
		}
	})

	// Test fetching tags with limit
	t.Run("FetchTagsWithLimit", func(t *testing.T) {
		limit := 5
		params := &GetTagsParams{
			Limit: limit,
		}

		tags, err := client.GetTags(ctx, params)
		if err != nil {
			t.Fatalf("Failed to fetch tags with limit: %v", err)
		}

		if len(tags) > limit {
			t.Errorf("Expected at most %d tags, got %d", limit, len(tags))
		}

		t.Logf("Successfully fetched %d tags with limit=%d", len(tags), limit)
	})

	// Test fetching carousel tags
	t.Run("FetchCarouselTags", func(t *testing.T) {
		isCarousel := true
		params := &GetTagsParams{
			Limit:      10,
			IsCarousel: &isCarousel,
		}

		tags, err := client.GetTags(ctx, params)
		if err != nil {
			t.Fatalf("Failed to fetch carousel tags: %v", err)
		}

		t.Logf("Successfully fetched %d carousel tags", len(tags))

		// Verify all tags are carousel tags
		for _, tag := range tags {
			if !tag.IsCarousel {
				t.Errorf("Expected carousel tag, got IsCarousel=false for tag=%s", tag.Label)
			}
		}
	})
}

func TestGetTagByID(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get a list of tags to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test GetTagByID")
	}

	testTagID := tags[0].ID
	testTagLabel := tags[0].Label

	t.Run("FetchTagByID", func(t *testing.T) {
		tag, err := client.GetTagByID(ctx, testTagID, nil)
		if err != nil {
			t.Fatalf("Failed to fetch tag by ID: %v", err)
		}

		if tag.ID != testTagID {
			t.Errorf("Expected tag ID=%s, got ID=%s", testTagID, tag.ID)
		}

		t.Logf("Successfully fetched tag by ID: %+v", tag)
	})

	t.Run("FetchTagByIDWithTemplate", func(t *testing.T) {
		includeTemplate := true
		params := &GetTagByIDQueryParams{
			IncludeTemplate: &includeTemplate,
		}

		tag, err := client.GetTagByID(ctx, testTagID, params)
		if err != nil {
			t.Fatalf("Failed to fetch tag by ID with template: %v", err)
		}

		t.Logf("Successfully fetched tag with template: %s", tag.Label)
	})

	_ = testTagLabel // Use variable to avoid unused warning
}

func TestGetTagBySlug(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// First get a list of tags to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test GetTagBySlug")
	}

	testTagSlug := tags[0].Slug
	testTagID := tags[0].ID

	t.Run("FetchTagBySlug", func(t *testing.T) {
		tag, err := client.GetTagBySlug(ctx, testTagSlug, nil)
		if err != nil {
			t.Fatalf("Failed to fetch tag by slug: %v", err)
		}

		if tag.Slug != testTagSlug {
			t.Errorf("Expected tag Slug=%s, got Slug=%s", testTagSlug, tag.Slug)
		}

		t.Logf("Successfully fetched tag by slug: %+v", tag)
	})

	t.Run("FetchTagBySlugWithTemplate", func(t *testing.T) {
		includeTemplate := true
		params := &GetTagBySlugQueryParams{
			IncludeTemplate: &includeTemplate,
		}

		tag, err := client.GetTagBySlug(ctx, testTagSlug, params)
		if err != nil {
			t.Fatalf("Failed to fetch tag by slug with template: %v", err)
		}

		t.Logf("Successfully fetched tag with template: %s", tag.Label)
	})

	_ = testTagID // Use variable to avoid unused warning
}

func TestGetTagByIDAndSlugConsistency(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Get a tag to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test consistency")
	}

	testTag := tags[0]

	// Fetch by ID
	tagByID, err := client.GetTagByID(ctx, testTag.ID, nil)
	if err != nil {
		t.Fatalf("Failed to fetch tag by ID: %v", err)
	}

	// Fetch by slug
	tagBySlug, err := client.GetTagBySlug(ctx, testTag.Slug, nil)
	if err != nil {
		t.Fatalf("Failed to fetch tag by slug: %v", err)
	}

	// Verify they return the same tag
	if tagByID.ID != tagBySlug.ID {
		t.Errorf("GetTagByID and GetTagBySlug returned different tags: ID=%s vs ID=%s",
			tagByID.ID, tagBySlug.ID)
	}

	if tagByID.Slug != tagBySlug.Slug {
		t.Errorf("GetTagByID and GetTagBySlug returned different slugs: %s vs %s",
			tagByID.Slug, tagBySlug.Slug)
	}

	t.Logf("✓ GetTagByID and GetTagBySlug return consistent results for tag: %s", tagByID.Label)
}

func TestGetRelatedTagsByID(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Get a tag to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test GetRelatedTagsByID")
	}

	testTagID := tags[0].ID

	t.Run("FetchRelatedTags", func(t *testing.T) {
		relationships, err := client.GetRelatedTagsByID(ctx, testTagID, nil)
		if err != nil {
			t.Fatalf("Failed to fetch related tags: %v", err)
		}

		t.Logf("Successfully fetched %d related tag relationships", len(relationships))
		for i, rel := range relationships {
			if i >= 3 {
				break
			}
			t.Logf("  Relationship %d: %+v", i+1, rel)
		}
	})

	t.Run("FetchRelatedTagsWithStatus", func(t *testing.T) {
		params := &GetRelatedTagsParams{
			Status: TagStatusActive,
		}

		relationships, err := client.GetRelatedTagsByID(ctx, testTagID, params)
		if err != nil {
			t.Fatalf("Failed to fetch related tags with status: %v", err)
		}

		t.Logf("Successfully fetched %d active related tag relationships", len(relationships))
	})
}

func TestGetRelatedTagsBySlug(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Get a tag to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test GetRelatedTagsBySlug")
	}

	testTagSlug := tags[0].Slug

	t.Run("FetchRelatedTagsBySlug", func(t *testing.T) {
		relationships, err := client.GetRelatedTagsBySlug(ctx, testTagSlug, nil)
		if err != nil {
			t.Fatalf("Failed to fetch related tags by slug: %v", err)
		}

		t.Logf("Successfully fetched %d related tag relationships", len(relationships))
	})
}

func TestGetRelatedTagsDetailByID(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Get a tag to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test GetRelatedTagsDetailByID")
	}

	testTagID := tags[0].ID

	t.Run("FetchRelatedTagsDetail", func(t *testing.T) {
		relatedTags, err := client.GetRelatedTagsDetailByID(ctx, testTagID, nil)
		if err != nil {
			t.Fatalf("Failed to fetch related tags detail: %v", err)
		}

		t.Logf("Successfully fetched %d related tags with details", len(relatedTags))
		for i, tag := range relatedTags {
			if i >= 3 {
				break
			}
			t.Logf("  Related tag %d: %+v", i+1, tag)
		}
	})

	t.Run("FetchRelatedTagsDetailWithStatusAll", func(t *testing.T) {
		params := &GetRelatedTagsParams{
			Status: TagStatusAll,
		}

		relatedTags, err := client.GetRelatedTagsDetailByID(ctx, testTagID, params)
		if err != nil {
			t.Fatalf("Failed to fetch related tags detail with status: %v", err)
		}

		t.Logf("Successfully fetched %d related tags (all statuses)", len(relatedTags))
	})
}

func TestGetRelatedTagsDetailBySlug(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Get a tag to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test GetRelatedTagsDetailBySlug")
	}

	testTagSlug := tags[0].Slug

	t.Run("FetchRelatedTagsDetailBySlug", func(t *testing.T) {
		relatedTags, err := client.GetRelatedTagsDetailBySlug(ctx, testTagSlug, nil)
		if err != nil {
			t.Fatalf("Failed to fetch related tags detail by slug: %v", err)
		}

		t.Logf("Successfully fetched %d related tags with details", len(relatedTags))
	})
}

func TestRelatedTagsConsistency(t *testing.T) {
	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	// Get a tag to test with
	tags, err := client.GetTags(ctx, &GetTagsParams{Limit: 1})
	if err != nil {
		t.Fatalf("Failed to fetch tags for test setup: %v", err)
	}

	if len(tags) == 0 {
		t.Skip("No tags available to test consistency")
	}

	testTag := tags[0]

	// Fetch relationships by ID
	relsByID, err := client.GetRelatedTagsByID(ctx, testTag.ID, nil)
	if err != nil {
		t.Fatalf("Failed to fetch relationships by ID: %v", err)
	}

	// Fetch relationships by slug
	relsBySlug, err := client.GetRelatedTagsBySlug(ctx, testTag.Slug, nil)
	if err != nil {
		t.Fatalf("Failed to fetch relationships by slug: %v", err)
	}

	// Verify they return the same number of relationships
	if len(relsByID) != len(relsBySlug) {
		t.Logf("Warning: Different number of relationships: ID=%d, Slug=%d",
			len(relsByID), len(relsBySlug))
	} else {
		t.Logf("✓ GetRelatedTagsByID and GetRelatedTagsBySlug return consistent count: %d", len(relsByID))
	}

	// Fetch detailed tags by ID
	tagsByID, err := client.GetRelatedTagsDetailByID(ctx, testTag.ID, nil)
	if err != nil {
		t.Fatalf("Failed to fetch related tags detail by ID: %v", err)
	}

	// Fetch detailed tags by slug
	tagsBySlug, err := client.GetRelatedTagsDetailBySlug(ctx, testTag.Slug, nil)
	if err != nil {
		t.Fatalf("Failed to fetch related tags detail by slug: %v", err)
	}

	// Verify they return the same number of tags
	if len(tagsByID) != len(tagsBySlug) {
		t.Logf("Warning: Different number of detailed tags: ID=%d, Slug=%d",
			len(tagsByID), len(tagsBySlug))
	} else {
		t.Logf("✓ GetRelatedTagsDetailByID and GetRelatedTagsDetailBySlug return consistent count: %d", len(tagsByID))
	}
}
