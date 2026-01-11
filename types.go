package polymarketgamma

import (
	"encoding/json"
	"strings"
	"time"
)

// NormalizedTime is a custom time type that can parse multiple time formats
type NormalizedTime time.Time

// UnmarshalJSON implements the json.Unmarshaler interface
func (ct *NormalizedTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		*ct = NormalizedTime(time.Time{})
		return nil
	}

	// Normalize timezone offset format: convert "+00" to "+00:00" for Go's time.Parse
	// API sometimes returns "2020-11-02 16:31:01+00" instead of "2020-11-02 16:31:01+00:00"
	// Only normalize if it ends with +XX or -XX (2 digits) and not already +XX:00 format
	if len(s) > 10 {
		last3 := s[len(s)-3:]
		// Check if it ends with +XX or -XX (2 digits)
		if (last3[0] == '+' || last3[0] == '-') && last3[1] >= '0' && last3[1] <= '9' && last3[2] >= '0' && last3[2] <= '9' {
			// Check if it's NOT already in +XX:00 format (check last 6 chars)
			if len(s) < 6 || s[len(s)-6:] != last3+":00" {
				// Not in +XX:00 format, convert "+00" to "+00:00"
				s = s[:len(s)-3] + last3 + ":00"
			}
		}
	}

	// Try different time formats
	// Note: RFC3339 format (e.g., "2024-11-06T15:17:41Z") should be tried first
	formats := []string{
		time.RFC3339,                     // "2006-01-02T15:04:05Z07:00" or "2006-01-02T15:04:05Z"
		time.RFC3339Nano,                 // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",           // Explicit Z format (e.g., "2024-11-06T15:17:41Z")
		"2006-01-02T15:04:05.999999999Z", // With nanoseconds and Z
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999+00:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05+00:00", // Format like "2020-11-02 16:31:01+00:00" (normalized)
		"2006-01-02",                // Simple date format (YYYY-MM-DD)
		"January 2, 2006",           // Long month name format (e.g., "November 1, 2022")
	}

	var err error
	var t time.Time
	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			*ct = NormalizedTime(t)
			return nil
		}
	}

	// If all formats fail, return the last error
	return err
}

// MarshalJSON implements the json.Marshaler interface
func (ct NormalizedTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format(time.RFC3339))
}

// Time returns the underlying time.Time value
func (ct NormalizedTime) Time() time.Time {
	return time.Time(ct)
}

// IsZero reports whether ct represents the zero time instant
func (ct NormalizedTime) IsZero() bool {
	return time.Time(ct).IsZero()
}

// Format returns a textual representation of the time value
func (ct NormalizedTime) Format(layout string) string {
	return time.Time(ct).Format(layout)
}

func (ct NormalizedTime) String() string {
	return ct.Time().String()
}

// StringOrArray is a custom type that can unmarshal both string and array from JSON
// API sometimes returns these fields as strings, sometimes as arrays
type StringOrArray []string

// UnmarshalJSON implements the json.Unmarshaler interface
func (sa *StringOrArray) UnmarshalJSON(b []byte) error {
	// Handle null explicitly
	if string(b) == "null" {
		*sa = StringOrArray([]string{})
		return nil
	}

	// Try to unmarshal as 1D array first
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		// Ensure non-nil slice (in case of null)
		if arr == nil {
			*sa = StringOrArray([]string{})
			return nil
		}

		// Check if any element is a JSON array string (e.g., "[\"Up\", \"Down\"]")
		// This handles cases where the API returns an array containing a JSON-encoded array string
		var flattened []string
		needsFlattening := false
		for _, s := range arr {
			// Check if the string looks like a JSON array
			if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
				// Try to unmarshal the string as a JSON array
				var innerArr []string
				if err := json.Unmarshal([]byte(s), &innerArr); err == nil {
					flattened = append(flattened, innerArr...)
					needsFlattening = true
					continue
				}

				// Try as 2D array and flatten
				var innerArr2D [][]string
				if err := json.Unmarshal([]byte(s), &innerArr2D); err == nil {
					for _, innerArr := range innerArr2D {
						flattened = append(flattened, innerArr...)
					}
					needsFlattening = true
					continue
				}
			}

			// Not a JSON array string, keep as-is
			flattened = append(flattened, s)
		}

		if needsFlattening {
			*sa = StringOrArray(flattened)
		} else {
			*sa = StringOrArray(arr)
		}
		return nil
	}

	// Try to unmarshal as 2D array and flatten it
	var arr2D [][]string
	if err := json.Unmarshal(b, &arr2D); err == nil {
		// Flatten the 2D array into 1D
		var flattened []string
		for _, innerArr := range arr2D {
			flattened = append(flattened, innerArr...)
		}
		*sa = StringOrArray(flattened)
		return nil
	}

	// If that fails, try as string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if s == "" {
			*sa = StringOrArray([]string{})
			return nil
		}

		// Check if the string is a JSON array (e.g., "[\"Yes\", \"No\"]")
		// This handles cases where the API returns a JSON-encoded array as a string
		if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
			// Try to unmarshal the string as a JSON array
			var arr []string
			if err := json.Unmarshal([]byte(s), &arr); err == nil {
				*sa = StringOrArray(arr)
				return nil
			}

			// Try as 2D array and flatten
			var arr2D [][]string
			if err := json.Unmarshal([]byte(s), &arr2D); err == nil {
				var flattened []string
				for _, innerArr := range arr2D {
					flattened = append(flattened, innerArr...)
				}
				*sa = StringOrArray(flattened)
				return nil
			}
		}

		// If not a JSON array, treat as a single string value
		*sa = StringOrArray([]string{s})
		return nil
	}

	// Default to empty array
	*sa = StringOrArray([]string{})
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (sa StringOrArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(sa))
}

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Data string `json:"data"`
}

// GetEventsParams represents parameters for fetching events list
type GetEventsParams struct {
	Limit           int             `json:"limit,omitempty"`  // Maximum number of events to return
	Offset          int             `json:"offset,omitempty"` // Pagination offset
	Order           string          `json:"order,omitempty"`  // Comma-separated list of fields to order by
	Ascending       *bool           `json:"ascending,omitempty"`
	ID              []int           `json:"id,omitempty"`
	Slug            []string        `json:"slug,omitempty"`
	TagID           *int            `json:"tag_id,omitempty"`
	ExcludeTagID    []int           `json:"exclude_tag_id,omitempty"`
	RelatedTags     *bool           `json:"related_tags,omitempty"`
	Featured        *bool           `json:"featured,omitempty"`
	CYOM            *bool           `json:"cyom,omitempty"`
	IncludeChat     *bool           `json:"include_chat,omitempty"`
	IncludeTemplate *bool           `json:"include_template,omitempty"`
	Recurrence      string          `json:"recurrence,omitempty"`
	Closed          *bool           `json:"closed,omitempty"`
	StartDateMin    *NormalizedTime `json:"start_date_min,omitempty"` // ISO 8601 date-time
	StartDateMax    *NormalizedTime `json:"start_date_max,omitempty"` // ISO 8601 date-time
	EndDateMin      *NormalizedTime `json:"end_date_min,omitempty"`   // ISO 8601 date-time
	EndDateMax      *NormalizedTime `json:"end_date_max,omitempty"`   // ISO 8601 date-time
}

// GetMarketsParams represents parameters for fetching markets list
type GetMarketsParams struct {
	Limit               int             `json:"limit,omitempty"`
	Offset              int             `json:"offset,omitempty"`
	Order               string          `json:"order,omitempty"`          // Comma-separated list of fields to order by
	Ascending           *bool           `json:"ascending,omitempty"`      // Use pointer to distinguish between false and unset
	ID                  []int           `json:"id,omitempty"`             // Market IDs
	Slug                []string        `json:"slug,omitempty"`           // Market slugs
	ClobTokenIDs        []string        `json:"clob_token_ids,omitempty"` // CLOB token IDs
	ConditionIDs        []string        `json:"condition_ids,omitempty"`  // Condition IDs
	MarketMakerAddress  []string        `json:"market_maker_address,omitempty"`
	LiquidityNumMin     *float64        `json:"liquidity_num_min,omitempty"`
	LiquidityNumMax     *float64        `json:"liquidity_num_max,omitempty"`
	VolumeNumMin        *float64        `json:"volume_num_min,omitempty"`
	VolumeNumMax        *float64        `json:"volume_num_max,omitempty"`
	StartDateMin        *NormalizedTime `json:"start_date_min,omitempty"` // ISO 8601 date-time
	StartDateMax        *NormalizedTime `json:"start_date_max,omitempty"` // ISO 8601 date-time
	EndDateMin          *NormalizedTime `json:"end_date_min,omitempty"`   // ISO 8601 date-time
	EndDateMax          *NormalizedTime `json:"end_date_max,omitempty"`   // ISO 8601 date-time
	TagID               *int            `json:"tag_id,omitempty"`
	RelatedTags         *bool           `json:"related_tags,omitempty"`
	CYOM                *bool           `json:"cyom,omitempty"`
	UMAResolutionStatus string          `json:"uma_resolution_status,omitempty"`
	GameID              string          `json:"game_id,omitempty"`
	SportsMarketTypes   []string        `json:"sports_market_types,omitempty"`
	RewardsMinSize      *float64        `json:"rewards_min_size,omitempty"`
	QuestionIDs         []string        `json:"question_ids,omitempty"`
	IncludeTag          *bool           `json:"include_tag,omitempty"`
	Closed              *bool           `json:"closed,omitempty"`
}

// GetMarketByIDQueryParams represents query parameters for fetching a single market by ID
type GetMarketByIDQueryParams struct {
	IncludeTag *bool `json:"include_tag,omitempty"` // Include tag information in the response
}

// GetSeriesParams represents parameters for fetching series list
type GetSeriesParams struct {
	Limit            int      `json:"limit,omitempty"`  // Maximum number of series to return
	Offset           int      `json:"offset,omitempty"` // Pagination offset
	Order            string   `json:"order,omitempty"`  // Comma-separated list of fields to order by
	Ascending        *bool    `json:"ascending,omitempty"`
	Slug             []string `json:"slug,omitempty"`
	CategoriesIDs    []int    `json:"categories_ids,omitempty"`
	CategoriesLabels []string `json:"categories_labels,omitempty"`
	Closed           *bool    `json:"closed,omitempty"`
	IncludeChat      *bool    `json:"include_chat,omitempty"`
	Recurrence       string   `json:"recurrence,omitempty"`
}

// GetEventByIDQueryParams represents query parameters for fetching a single event by ID
type GetEventByIDQueryParams struct {
	IncludeChat     *bool `json:"include_chat,omitempty"`     // Include chat information in the response
	IncludeTemplate *bool `json:"include_template,omitempty"` // Include template information in the response
}

// GetEventBySlugQueryParams represents query parameters for fetching a single event by slug
type GetEventBySlugQueryParams struct {
	IncludeChat     *bool `json:"include_chat,omitempty"`     // Include chat information in the response
	IncludeTemplate *bool `json:"include_template,omitempty"` // Include template information in the response
}

// GetSeriesByIDQueryParams represents query parameters for fetching a single series by ID
type GetSeriesByIDQueryParams struct {
	IncludeChat *bool `json:"include_chat,omitempty"` // Include chat information in the response
}

// GetTeamsParams represents query parameters for fetching teams
type GetTeamsParams struct {
	Limit        int      `json:"limit,omitempty"`        // Number of results to return (x >= 0)
	Offset       int      `json:"offset,omitempty"`       // Number of results to skip (x >= 0)
	Order        string   `json:"order,omitempty"`        // Comma-separated list of fields to order by
	Ascending    *bool    `json:"ascending,omitempty"`    // Sort order
	League       []string `json:"league,omitempty"`       // Filter by league(s)
	Name         []string `json:"name,omitempty"`         // Filter by name(s)
	Abbreviation []string `json:"abbreviation,omitempty"` // Filter by abbreviation(s)
}

// GetTagsParams represents query parameters for fetching tags
type GetTagsParams struct {
	Limit           int    `json:"limit,omitempty"`            // Number of results to return (x >= 0)
	Offset          int    `json:"offset,omitempty"`           // Number of results to skip (x >= 0)
	Order           string `json:"order,omitempty"`            // Comma-separated list of fields to order by
	Ascending       *bool  `json:"ascending,omitempty"`        // Sort order
	IncludeTemplate *bool  `json:"include_template,omitempty"` // Include template information
	IsCarousel      *bool  `json:"is_carousel,omitempty"`      // Filter by carousel status
}

// GetTagByIDQueryParams represents query parameters for fetching a single tag by ID
type GetTagByIDQueryParams struct {
	IncludeTemplate *bool `json:"include_template,omitempty"` // Include template information
}

// GetTagBySlugQueryParams represents query parameters for fetching a single tag by slug
type GetTagBySlugQueryParams struct {
	IncludeTemplate *bool `json:"include_template,omitempty"` // Include template information
}

// GetRelatedTagsParams represents query parameters for fetching related tags
type GetRelatedTagsParams struct {
	OmitEmpty *bool     `json:"omit_empty,omitempty"` // Omit empty relationships
	Status    TagStatus `json:"status,omitempty"`     // Status filter (active, closed, all)
}

// SearchParams represents parameters for searching markets, events, and profiles
type SearchParams struct {
	Q                 string   `json:"q"`                             // Search query (required)
	Cache             *bool    `json:"cache,omitempty"`               // Use cache
	EventsStatus      string   `json:"events_status,omitempty"`       // Events status filter
	LimitPerType      *int     `json:"limit_per_type,omitempty"`      // Limit results per type
	Page              *int     `json:"page,omitempty"`                // Page number for pagination
	EventsTag         []string `json:"events_tag,omitempty"`          // Filter by event tags
	KeepClosedMarkets *int     `json:"keep_closed_markets,omitempty"` // Keep closed markets (0 or 1)
	Sort              string   `json:"sort,omitempty"`                // Sort field
	Ascending         *bool    `json:"ascending,omitempty"`           // Sort order
	SearchTags        *bool    `json:"search_tags,omitempty"`         // Include tags in search
	SearchProfiles    *bool    `json:"search_profiles,omitempty"`     // Include profiles in search
	Recurrence        string   `json:"recurrence,omitempty"`          // Recurrence filter
	ExcludeTagID      []int    `json:"exclude_tag_id,omitempty"`      // Exclude tag IDs
	Optimized         *bool    `json:"optimized,omitempty"`           // Return optimized images
}

// SearchResponse represents the response from the search endpoint
type SearchResponse struct {
	Events     []Event     `json:"events"`
	Tags       []SearchTag `json:"tags"`
	Profiles   []Profile   `json:"profiles"`
	Pagination Pagination  `json:"pagination"`
}

// SearchTag represents a tag in search results
type SearchTag struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Slug       string `json:"slug"`
	EventCount int    `json:"event_count"`
}

// Profile represents a user profile
type Profile struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	User                  int             `json:"user"`
	Referral              string          `json:"referral"`
	CreatedBy             int             `json:"createdBy"`
	UpdatedBy             int             `json:"updatedBy"`
	CreatedAt             NormalizedTime  `json:"createdAt"`
	UpdatedAt             NormalizedTime  `json:"updatedAt"`
	UTMSource             string          `json:"utmSource"`
	UTMMedium             string          `json:"utmMedium"`
	UTMCampaign           string          `json:"utmCampaign"`
	UTMContent            string          `json:"utmContent"`
	UTMTerm               string          `json:"utmTerm"`
	WalletActivated       bool            `json:"walletActivated"`
	Pseudonym             string          `json:"pseudonym"`
	DisplayUsernamePublic bool            `json:"displayUsernamePublic"`
	ProfileImage          string          `json:"profileImage"`
	Bio                   string          `json:"bio"`
	ProxyWallet           string          `json:"proxyWallet"`
	ProfileImageOptimized *ImageOptimized `json:"profileImageOptimized,omitempty"`
	IsCloseOnly           bool            `json:"isCloseOnly"`
	IsCertReq             bool            `json:"isCertReq"`
	CertReqDate           NormalizedTime  `json:"certReqDate"`
}

// Pagination represents pagination information
type Pagination struct {
	HasMore      bool `json:"hasMore"`
	TotalResults int  `json:"totalResults"`
}

// Market represents a market from the Gamma API /markets endpoint
// This contains extensive metadata beyond what's in the CLOB API Market type
type Market struct {
	ID                           string    `json:"id"`
	Question                     string    `json:"question"`
	ConditionID                  string    `json:"conditionId"`
	Slug                         string    `json:"slug"`
	ResolutionSource             string    `json:"resolutionSource"`
	EndDate                      time.Time `json:"endDate"`
	Liquidity                    string    `json:"liquidity"`
	StartDate                    time.Time `json:"startDate"`
	Image                        string    `json:"image"`
	Icon                         string    `json:"icon"`
	Description                  string    `json:"description"`
	Outcomes                     string    `json:"outcomes"`
	Volume                       string    `json:"volume"`
	Active                       bool      `json:"active"`
	Closed                       bool      `json:"closed"`
	MarketMakerAddress           string    `json:"marketMakerAddress"`
	CreatedAt                    time.Time `json:"createdAt"`
	New                          bool      `json:"new"`
	Featured                     bool      `json:"featured"`
	Archived                     bool      `json:"archived"`
	Restricted                   bool      `json:"restricted"`
	GroupItemThreshold           string    `json:"groupItemThreshold"`
	QuestionID                   string    `json:"questionID"`
	EnableOrderBook              bool      `json:"enableOrderBook"`
	OrderPriceMinTickSize        float64   `json:"orderPriceMinTickSize"`
	OrderMinSize                 int       `json:"orderMinSize"`
	VolumeNum                    int       `json:"volumeNum"`
	LiquidityNum                 int       `json:"liquidityNum"`
	EndDateIso                   string    `json:"endDateIso"`
	StartDateIso                 string    `json:"startDateIso"`
	HasReviewedDates             bool      `json:"hasReviewedDates"`
	Volume24Hr                   int       `json:"volume24hr"`
	Volume1Wk                    int       `json:"volume1wk"`
	Volume1Mo                    int       `json:"volume1mo"`
	Volume1Yr                    int       `json:"volume1yr"`
	ClobTokenIds                 string    `json:"clobTokenIds"`
	Volume24HrAmm                int       `json:"volume24hrAmm"`
	Volume1WkAmm                 int       `json:"volume1wkAmm"`
	Volume1MoAmm                 int       `json:"volume1moAmm"`
	Volume1YrAmm                 int       `json:"volume1yrAmm"`
	Volume24HrClob               int       `json:"volume24hrClob"`
	Volume1WkClob                int       `json:"volume1wkClob"`
	Volume1MoClob                int       `json:"volume1moClob"`
	Volume1YrClob                int       `json:"volume1yrClob"`
	VolumeAmm                    int       `json:"volumeAmm"`
	VolumeClob                   int       `json:"volumeClob"`
	LiquidityAmm                 int       `json:"liquidityAmm"`
	LiquidityClob                int       `json:"liquidityClob"`
	NegRisk                      bool      `json:"negRisk"`
	Ready                        bool      `json:"ready"`
	Funded                       bool      `json:"funded"`
	Cyom                         bool      `json:"cyom"`
	Competitive                  int       `json:"competitive"`
	PagerDutyNotificationEnabled bool      `json:"pagerDutyNotificationEnabled"`
	Approved                     bool      `json:"approved"`
	RewardsMinSize               int       `json:"rewardsMinSize"`
	RewardsMaxSpread             int       `json:"rewardsMaxSpread"`
	Spread                       int       `json:"spread"`
	OneDayPriceChange            int       `json:"oneDayPriceChange"`
	OneHourPriceChange           int       `json:"oneHourPriceChange"`
	OneWeekPriceChange           int       `json:"oneWeekPriceChange"`
	OneMonthPriceChange          int       `json:"oneMonthPriceChange"`
	OneYearPriceChange           int       `json:"oneYearPriceChange"`
	LastTradePrice               int       `json:"lastTradePrice"`
	BestBid                      int       `json:"bestBid"`
	BestAsk                      int       `json:"bestAsk"`
	AutomaticallyActive          bool      `json:"automaticallyActive"`
	ClearBookOnStart             bool      `json:"clearBookOnStart"`
	ShowGmpSeries                bool      `json:"showGmpSeries"`
	ShowGmpOutcome               bool      `json:"showGmpOutcome"`
	ManualActivation             bool      `json:"manualActivation"`
	NegRiskOther                 bool      `json:"negRiskOther"`
	UmaResolutionStatuses        string    `json:"umaResolutionStatuses"`
	PendingDeployment            bool      `json:"pendingDeployment"`
	Deploying                    bool      `json:"deploying"`
	RfqEnabled                   bool      `json:"rfqEnabled"`
	EventStartTime               time.Time `json:"eventStartTime"`
	HoldingRewardsEnabled        bool      `json:"holdingRewardsEnabled"`
	FeesEnabled                  bool      `json:"feesEnabled"`
	RequiresTranslation          bool      `json:"requiresTranslation"`
}

// ImageOptimized represents an optimized image resource
type ImageOptimized struct {
	ID                        string         `json:"id"`
	ImageURLSource            string         `json:"imageUrlSource"`
	ImageURLOptimized         string         `json:"imageUrlOptimized"`
	ImageSizeKBSource         int            `json:"imageSizeKbSource"`
	ImageSizeKBOptimized      int            `json:"imageSizeKbOptimized"`
	ImageOptimizedComplete    bool           `json:"imageOptimizedComplete"`
	ImageOptimizedLastUpdated NormalizedTime `json:"imageOptimizedLastUpdated"`
	RelID                     int            `json:"relID"`
	Field                     string         `json:"field"`
	Relname                   string         `json:"relname"`
}

// Event represents an event from the Gamma API
type Event struct {
	ID                  string    `json:"id"`
	Ticker              string    `json:"ticker"`
	Slug                string    `json:"slug"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	ResolutionSource    string    `json:"resolutionSource"`
	StartDate           time.Time `json:"startDate"`
	CreationDate        time.Time `json:"creationDate"`
	EndDate             time.Time `json:"endDate"`
	Image               string    `json:"image"`
	Icon                string    `json:"icon"`
	Active              bool      `json:"active"`
	Closed              bool      `json:"closed"`
	Archived            bool      `json:"archived"`
	New                 bool      `json:"new"`
	Featured            bool      `json:"featured"`
	Restricted          bool      `json:"restricted"`
	OpenInterest        int       `json:"openInterest"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	EnableOrderBook     bool      `json:"enableOrderBook"`
	NegRisk             bool      `json:"negRisk"`
	CommentCount        int       `json:"commentCount"`
	Markets             []Market `json:"markets"`
	Series              []Series  `json:"series"`
	Tags                []Tag     `json:"tags"`
	Cyom                bool      `json:"cyom"`
	ShowAllOutcomes     bool      `json:"showAllOutcomes"`
	ShowMarketImages    bool      `json:"showMarketImages"`
	EnableNegRisk       bool      `json:"enableNegRisk"`
	StartTime           time.Time `json:"startTime"`
	SeriesSlug          string    `json:"seriesSlug"`
	NegRiskAugmented    bool      `json:"negRiskAugmented"`
	PendingDeployment   bool      `json:"pendingDeployment"`
	Deploying           bool      `json:"deploying"`
	RequiresTranslation bool      `json:"requiresTranslation"`
}

// Category represents a category or subcategory
type Category struct {
	ID             string         `json:"id"`
	Label          string         `json:"label"`
	ParentCategory string         `json:"parentCategory"`
	Slug           string         `json:"slug"`
	PublishedAt    NormalizedTime `json:"publishedAt"`
	CreatedBy      string         `json:"createdBy"`
	UpdatedBy      string         `json:"updatedBy"`
	CreatedAt      NormalizedTime `json:"createdAt"`
	UpdatedAt      NormalizedTime `json:"updatedAt"`
}

// Tag represents a tag associated with markets or events
type Tag struct {
	ID                  string    `json:"id"`
	Label               string    `json:"label"`
	Slug                string    `json:"slug"`
	ForceShow           bool      `json:"forceShow,omitempty"`
	PublishedAt         string    `json:"publishedAt,omitempty"`
	CreatedBy           int       `json:"createdBy,omitempty"`
	UpdatedBy           int       `json:"updatedBy,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	RequiresTranslation bool      `json:"requiresTranslation"`
	IsCarousel          bool      `json:"isCarousel,omitempty"`
}

// Series represents a series of events
type Series struct {
	ID                  string    `json:"id"`
	Ticker              string    `json:"ticker"`
	Slug                string    `json:"slug"`
	Title               string    `json:"title"`
	SeriesType          string    `json:"seriesType"`
	Recurrence          string    `json:"recurrence"`
	Image               string    `json:"image"`
	Icon                string    `json:"icon"`
	Active              bool      `json:"active"`
	Closed              bool      `json:"closed"`
	Archived            bool      `json:"archived"`
	Featured            bool      `json:"featured"`
	Restricted          bool      `json:"restricted"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	CommentCount        int       `json:"commentCount"`
	RequiresTranslation bool      `json:"requiresTranslation"`
}

// Collection represents a collection of events
type Collection struct {
	ID                   string          `json:"id"`
	Ticker               string          `json:"ticker"`
	Slug                 string          `json:"slug"`
	Title                string          `json:"title"`
	Subtitle             string          `json:"subtitle"`
	CollectionType       string          `json:"collectionType"`
	Description          string          `json:"description"`
	Tags                 string          `json:"tags"`
	Image                string          `json:"image"`
	Icon                 string          `json:"icon"`
	HeaderImage          string          `json:"headerImage"`
	Layout               string          `json:"layout"`
	Active               bool            `json:"active"`
	Closed               bool            `json:"closed"`
	Archived             bool            `json:"archived"`
	New                  bool            `json:"new"`
	Featured             bool            `json:"featured"`
	Restricted           bool            `json:"restricted"`
	IsTemplate           bool            `json:"isTemplate"`
	TemplateVariables    string          `json:"templateVariables"`
	PublishedAt          NormalizedTime  `json:"publishedAt"`
	CreatedBy            string          `json:"createdBy"`
	UpdatedBy            string          `json:"updatedBy"`
	CreatedAt            NormalizedTime  `json:"createdAt"`
	UpdatedAt            NormalizedTime  `json:"updatedAt"`
	CommentsEnabled      bool            `json:"commentsEnabled"`
	ImageOptimized       *ImageOptimized `json:"imageOptimized,omitempty"`
	IconOptimized        *ImageOptimized `json:"iconOptimized,omitempty"`
	HeaderImageOptimized *ImageOptimized `json:"headerImageOptimized,omitempty"`
}

// EventCreator represents a creator of an event
type EventCreator struct {
	ID            string         `json:"id"`
	CreatorName   string         `json:"creatorName"`
	CreatorHandle string         `json:"creatorHandle"`
	CreatorURL    string         `json:"creatorUrl"`
	CreatorImage  string         `json:"creatorImage"`
	CreatedAt     NormalizedTime `json:"createdAt"`
	UpdatedAt     NormalizedTime `json:"updatedAt"`
}

// Chat represents a chat channel associated with an event
type Chat struct {
	ID           string         `json:"id"`
	ChannelID    string         `json:"channelId"`
	ChannelName  string         `json:"channelName"`
	ChannelImage string         `json:"channelImage"`
	Live         bool           `json:"live"`
	StartTime    NormalizedTime `json:"startTime"`
	EndTime      NormalizedTime `json:"endTime"`
}

// Template represents an event template
type Template struct {
	ID               string `json:"id"`
	EventTitle       string `json:"eventTitle"`
	EventSlug        string `json:"eventSlug"`
	EventImage       string `json:"eventImage"`
	MarketTitle      string `json:"marketTitle"`
	Description      string `json:"description"`
	ResolutionSource string `json:"resolutionSource"`
	NegRisk          bool   `json:"negRisk"`
	SortBy           string `json:"sortBy"`
	ShowMarketImages bool   `json:"showMarketImages"`
	SeriesSlug       string `json:"seriesSlug"`
	Outcomes         string `json:"outcomes"`
}

// Team represents a sports team
type Team struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	League       string         `json:"league"`
	Record       string         `json:"record"`
	Logo         string         `json:"logo"`
	Abbreviation string         `json:"abbreviation"`
	Alias        string         `json:"alias"`
	CreatedAt    NormalizedTime `json:"createdAt"`
	UpdatedAt    NormalizedTime `json:"updatedAt"`
}

// SportMetadata represents metadata information for a sport
type SportMetadata struct {
	Sport      string `json:"sport"`
	Image      string `json:"image"`
	Resolution string `json:"resolution"`
	Ordering   string `json:"ordering"`
	Tags       string `json:"tags"`
	Series     string `json:"series"`
}

// TagRelationship represents a relationship between two tags
type TagRelationship struct {
	ID           string `json:"id"`
	TagID        int    `json:"tagID"`
	RelatedTagID int    `json:"relatedTagID"`
	Rank         int    `json:"rank"`
}

// TagStatus represents the status filter for related tags
type TagStatus string

const (
	TagStatusActive TagStatus = "active"
	TagStatusClosed TagStatus = "closed"
	TagStatusAll    TagStatus = "all"
)
