package polymarketgamma

import (
	"net/http"
	"testing"
)

func TestGetTeams(t *testing.T) {
	client := NewClient(http.DefaultClient)

	// Test fetching teams without parameters
	t.Run("FetchTeamsNoParams", func(t *testing.T) {
		teams, err := client.GetTeams(nil)
		if err != nil {
			t.Fatalf("Failed to fetch teams: %v", err)
		}

		if len(teams) == 0 {
			t.Log("No teams found (this might be expected if the API returns empty data)")
		} else {
			t.Logf("Successfully fetched %d teams", len(teams))
			// Print first team details
			if len(teams) > 0 {
				team := teams[0]
				t.Logf("First team: ID=%d, Name=%s, League=%s, Abbreviation=%s",
					team.ID, team.Name, team.League, team.Abbreviation)
			}
		}
	})

	// Test fetching teams with limit
	t.Run("FetchTeamsWithLimit", func(t *testing.T) {
		limit := 5
		params := &GetTeamsParams{
			Limit: limit,
		}

		teams, err := client.GetTeams(params)
		if err != nil {
			t.Fatalf("Failed to fetch teams with limit: %v", err)
		}

		if len(teams) > limit {
			t.Errorf("Expected at most %d teams, got %d", limit, len(teams))
		}

		t.Logf("Successfully fetched %d teams with limit=%d", len(teams), limit)
	})

	// Test fetching teams with filter
	t.Run("FetchTeamsWithLeagueFilter", func(t *testing.T) {
		params := &GetTeamsParams{
			Limit:  10,
			League: []string{"NBA"},
		}

		teams, err := client.GetTeams(params)
		if err != nil {
			t.Fatalf("Failed to fetch teams with league filter: %v", err)
		}

		t.Logf("Successfully fetched %d NBA teams", len(teams))

		// Verify all teams are from NBA
		for _, team := range teams {
			if team.League != "NBA" {
				t.Errorf("Expected NBA team, got league=%s for team=%s", team.League, team.Name)
			}
		}
	})
}

func TestGetSportsMetadata(t *testing.T) {
	client := NewClient(http.DefaultClient)

	metadata, err := client.GetSportsMetadata()
	if err != nil {
		t.Fatalf("Failed to fetch sports metadata: %v", err)
	}

	if len(metadata) == 0 {
		t.Log("No sports metadata found (this might be expected if the API returns empty data)")
	} else {
		t.Logf("Successfully fetched metadata for %d sports", len(metadata))

		// Print details for each sport
		for i, sport := range metadata {
			t.Logf("Sport %d: %s (Image: %s, Resolution: %s)",
				i+1, sport.Sport, sport.Image, sport.Resolution)
		}
	}
}
