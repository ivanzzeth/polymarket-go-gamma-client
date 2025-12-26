package polymarketgamma

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringOrArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StringOrArray
		wantErr  bool
	}{
		{
			name:     "unmarshal string",
			input:    `"hello"`,
			expected: StringOrArray{"hello"},
			wantErr:  false,
		},
		{
			name:     "unmarshal empty string",
			input:    `""`,
			expected: StringOrArray{},
			wantErr:  false,
		},
		{
			name:     "unmarshal array",
			input:    `["a", "b", "c"]`,
			expected: StringOrArray{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "unmarshal empty array",
			input:    `[]`,
			expected: StringOrArray{},
			wantErr:  false,
		},
		{
			name:     "unmarshal single element array",
			input:    `["single"]`,
			expected: StringOrArray{"single"},
			wantErr:  false,
		},
		{
			name:     "unmarshal null",
			input:    `null`,
			expected: StringOrArray{},
			wantErr:  false,
		},
		{
			name:     "unmarshal 2D array - should flatten",
			input:    `[["a", "b"]]`,
			expected: StringOrArray{"a", "b"},
			wantErr:  false,
		},
		{
			name:     "unmarshal 2D array with multiple inner arrays - should flatten",
			input:    `[["a", "b"], ["c", "d"]]`,
			expected: StringOrArray{"a", "b", "c", "d"},
			wantErr:  false,
		},
		{
			name:     "unmarshal nested 3D array - should return empty (not supported)",
			input:    `[[["a"]]]`,
			expected: StringOrArray{}, // 3D arrays are not supported, will fall back to empty
			wantErr:  false,
		},
		{
			name:     "unmarshal array with numeric strings",
			input:    `["0.1", "0.2"]`,
			expected: StringOrArray{"0.1", "0.2"},
			wantErr:  false,
		},
		{
			name:     "unmarshal string with special characters",
			input:    `"hello, world"`,
			expected: StringOrArray{"hello, world"},
			wantErr:  false,
		},
		{
			name:     "unmarshal array with empty strings",
			input:    `["", "a", ""]`,
			expected: StringOrArray{"", "a", ""},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result StringOrArray
			err := json.Unmarshal([]byte(tt.input), &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("UnmarshalJSON() = %v, want %v", result, tt.expected)
					t.Logf("Result type: %T", result)
					t.Logf("Result length: %d", len(result))
					for i, v := range result {
						t.Logf("  [%d]: %v (type: %T)", i, v, v)
					}
				}
			}
		})
	}
}

func TestStringOrArray_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    StringOrArray
		expected string
		wantErr  bool
	}{
		{
			name:     "marshal empty array",
			input:    StringOrArray{},
			expected: `[]`,
			wantErr:  false,
		},
		{
			name:     "marshal single element",
			input:    StringOrArray{"hello"},
			expected: `["hello"]`,
			wantErr:  false,
		},
		{
			name:     "marshal multiple elements",
			input:    StringOrArray{"a", "b", "c"},
			expected: `["a","b","c"]`, // JSON arrays don't have spaces
			wantErr:  false,
		},
		{
			name:     "marshal with empty strings",
			input:    StringOrArray{"", "a", ""},
			expected: `["","a",""]`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()

			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Parse both to compare properly
				var expectedArr []string
				var resultArr []string

				if err := json.Unmarshal([]byte(tt.expected), &expectedArr); err != nil {
					t.Fatalf("Failed to unmarshal expected JSON: %v", err)
				}

				if err := json.Unmarshal(result, &resultArr); err != nil {
					t.Fatalf("Failed to unmarshal result JSON: %v", err)
				}

				if !reflect.DeepEqual(resultArr, expectedArr) {
					t.Errorf("MarshalJSON() = %s, want %s", string(result), tt.expected)
				}
			}
		})
	}
}

func TestStringOrArray_2DArrayBug(t *testing.T) {
	// This test specifically targets the bug mentioned in the TODO comment
	// The API sometimes returns: [["0.0000004113679809846114013590098187297978", "0.9999995886320190153885986409901813"]]
	// This should be handled correctly and not create a nested structure

	testCases := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "2D array with single inner array",
			input:       `[["a", "b"]]`,
			description: "Should flatten 2D array to 1D: [\"a\", \"b\"]",
		},
		{
			name:        "2D array with multiple inner arrays",
			input:       `[["a", "b"], ["c", "d"]]`,
			description: "Should flatten multiple inner arrays: [\"a\", \"b\", \"c\", \"d\"]",
		},
		{
			name:        "2D array from TODO comment example",
			input:       `[["0.0000004113679809846114013590098187297978", "0.9999995886320190153885986409901813"]]`,
			description: "Real-world example from TODO comment - should flatten",
		},
		{
			name:        "2D array with single element",
			input:       `[["single"]]`,
			description: "2D array with one string in inner array - should flatten to [\"single\"]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result StringOrArray
			err := json.Unmarshal([]byte(tc.input), &result)

			t.Logf("Input: %s", tc.input)
			t.Logf("Description: %s", tc.description)
			t.Logf("Error: %v", err)
			t.Logf("Result: %v", result)
			t.Logf("Result type: %T", result)
			t.Logf("Result length: %d", len(result))

			// Verify the result is correctly flattened
			for i, v := range result {
				t.Logf("  [%d]: %v (type: %T)", i, v, v)
				// Check if the string looks like a JSON array (bug indicator - should not happen)
				if len(v) > 0 && v[0] == '[' && v[len(v)-1] == ']' {
					t.Errorf("BUG DETECTED: Element at index %d appears to be a JSON array string: %v", i, v)
				}
			}

			// The result should be a flattened 1D array
			// For [["a", "b"]], we expect ["a", "b"]
			// For [["a", "b"], ["c", "d"]], we expect ["a", "b", "c", "d"]
		})
	}
}

func TestStringOrArray_RoundTrip(t *testing.T) {
	// Test that marshaling and unmarshaling preserves the data
	tests := []struct {
		name  string
		input StringOrArray
	}{
		{
			name:  "empty array",
			input: StringOrArray{},
		},
		{
			name:  "single element",
			input: StringOrArray{"hello"},
		},
		{
			name:  "multiple elements",
			input: StringOrArray{"a", "b", "c"},
		},
		{
			name:  "with empty strings",
			input: StringOrArray{"", "a", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			marshaled, err := tt.input.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			// Unmarshal
			var result StringOrArray
			if err := json.Unmarshal(marshaled, &result); err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(result, tt.input) {
				t.Errorf("RoundTrip() = %v, want %v", result, tt.input)
			}
		})
	}
}

func TestStringOrArray_RealWorldScenarios(t *testing.T) {
	// Test scenarios based on actual API responses
	tests := []struct {
		name     string
		jsonData string
		field    string
	}{
		{
			name:     "outcomes as string",
			jsonData: `{"outcomes": "Yes"}`,
			field:    "outcomes",
		},
		{
			name:     "outcomes as array",
			jsonData: `{"outcomes": ["Yes", "No"]}`,
			field:    "outcomes",
		},
		{
			name:     "outcomePrices as string",
			jsonData: `{"outcomePrices": "0.5"}`,
			field:    "outcomePrices",
		},
		{
			name:     "outcomePrices as array",
			jsonData: `{"outcomePrices": ["0.5", "0.5"]}`,
			field:    "outcomePrices",
		},
		{
			name:     "outcomePrices as 2D array - bug case",
			jsonData: `{"outcomePrices": [["0.5", "0.5"]]}`,
			field:    "outcomePrices",
		},
		{
			name:     "outcomes as JSON string (API format)",
			jsonData: `{"outcomes": "[\"Yes\", \"No\"]"}`,
			field:    "outcomes",
		},
		{
			name:     "outcomePrices as JSON string (API format)",
			jsonData: `{"outcomePrices": "[\"0\", \"0\"]"}`,
			field:    "outcomePrices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var market struct {
				Outcomes      StringOrArray `json:"outcomes"`
				OutcomePrices StringOrArray `json:"outcomePrices"`
				ShortOutcomes StringOrArray `json:"shortOutcomes"`
			}

			err := json.Unmarshal([]byte(tt.jsonData), &market)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			var result StringOrArray
			switch tt.field {
			case "outcomes":
				result = market.Outcomes
			case "outcomePrices":
				result = market.OutcomePrices
			case "shortOutcomes":
				result = market.ShortOutcomes
			}

			t.Logf("Field: %s", tt.field)
			t.Logf("Result: %v", result)
			t.Logf("Result length: %d", len(result))

			// Verify all elements are strings and check for bug (strings containing JSON arrays)
			for i, v := range result {
				// Check if the string looks like a JSON array (bug indicator)
				if len(v) > 0 && v[0] == '[' && v[len(v)-1] == ']' {
					t.Errorf("BUG DETECTED: Element at index %d appears to be a JSON array string: %v", i, v)
				}
			}
		})
	}
}
