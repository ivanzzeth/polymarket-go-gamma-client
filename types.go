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
	if len(s) >= 10 {
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
		time.RFC3339,                    // "2006-01-02T15:04:05Z07:00" or "2006-01-02T15:04:05Z"
		time.RFC3339Nano,                // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",          // Explicit Z format (e.g., "2024-11-06T15:17:41Z")
		"2006-01-02T15:04:05.999999999Z", // With nanoseconds and Z
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999+00:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05+00:00",     // Format like "2020-11-02 16:31:01+00:00" (normalized)
		"2006-01-02",                     // Simple date format (YYYY-MM-DD)
		"January 2, 2006",               // Long month name format (e.g., "November 1, 2022")
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
	// Core market identifiers
	ID          string `json:"id"`
	Question    string `json:"question"`
	ConditionID string `json:"conditionId"`
	Slug        string `json:"slug"`
	QuestionID  string `json:"questionID"`

	// Display and branding
	TwitterCardImage string `json:"twitterCardImage"`
	Image            string `json:"image"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`

	// Resolution and timing
	ResolutionSource string         `json:"resolutionSource"`
	EndDate          NormalizedTime `json:"endDate"`
	StartDate        NormalizedTime `json:"startDate"`
	EndDateISO       string         `json:"endDateIso"`
	StartDateISO     string         `json:"startDateIso"`
	UMAEndDate       NormalizedTime `json:"umaEndDate"`
	UMAEndDateISO    string         `json:"umaEndDateIso"`
	ClosedTime       NormalizedTime `json:"closedTime"`

	// Market mechanics
	Category          string  `json:"category"`
	AmmType           string  `json:"ammType"`
	Liquidity         string  `json:"liquidity"`
	LiquidityNum      float64 `json:"liquidityNum"`
	Volume            string  `json:"volume"`
	VolumeNum         float64 `json:"volumeNum"`
	Fee               string  `json:"fee"`
	DenominationToken string  `json:"denominationToken"`

	// Sponsor information
	SponsorName  string `json:"sponsorName"`
	SponsorImage string `json:"sponsorImage"`

	// Chart configuration
	XAxisValue string `json:"xAxisValue"`
	YAxisValue string `json:"yAxisValue"`
	LowerBound string `json:"lowerBound"`
	UpperBound string `json:"upperBound"`

	// Outcomes and pricing
	// Outcomes represents the possible outcomes for a market.
	// API returns this as a JSON-encoded string (e.g., "[\"Yes\", \"No\"]") which is automatically parsed into a string array.
	// Example: ["Yes", "No"] for binary markets, or ["Option A", "Option B", "Option C"] for categorical markets.
	Outcomes StringOrArray `json:"outcomes"`
	// OutcomePrices represents the current prices for each outcome, typically as decimal strings.
	// API returns this as a JSON-encoded string (e.g., "[\"0.52\", \"0.48\"]") which is automatically parsed into a string array.
	// Example: ["0.52", "0.48"] for binary markets, where prices sum to 1.0.
	OutcomePrices StringOrArray `json:"outcomePrices"`
	ShortOutcomes StringOrArray `json:"shortOutcomes"`

	// Status flags
	Active     bool `json:"active"`
	Closed     bool `json:"closed"`
	Archived   bool `json:"archived"`
	New        bool `json:"new"`
	Featured   bool `json:"featured"`
	Restricted bool `json:"restricted"`
	WideFormat bool `json:"wideFormat"`
	Ready      bool `json:"ready"`
	Funded     bool `json:"funded"`

	// Market type and format
	MarketType string `json:"marketType"`
	FormatType string `json:"formatType"`

	// Date boundaries
	LowerBoundDate NormalizedTime `json:"lowerBoundDate"`
	UpperBoundDate NormalizedTime `json:"upperBoundDate"`

	// Contract and exchange
	MarketMakerAddress string `json:"marketMakerAddress"`

	// Metadata
	CreatedBy int            `json:"createdBy"`
	UpdatedBy int            `json:"updatedBy"`
	CreatedAt NormalizedTime `json:"createdAt"`
	UpdatedAt NormalizedTime `json:"updatedAt"`

	// Marketing
	MailchimpTag string `json:"mailchimpTag"`
	ResolvedBy   string `json:"resolvedBy"`

	// Grouping
	MarketGroup        int    `json:"marketGroup"`
	GroupItemTitle     string `json:"groupItemTitle"`
	GroupItemThreshold string `json:"groupItemThreshold"`
	GroupItemRange     string `json:"groupItemRange"`

	// UMA resolution
	UMAResolutionStatus   string `json:"umaResolutionStatus"`
	UMAResolutionStatuses string `json:"umaResolutionStatuses"`
	UMABond               string `json:"umaBond"`
	UMAReward             string `json:"umaReward"`

	// Order book configuration
	EnableOrderBook       bool    `json:"enableOrderBook"`
	OrderPriceMinTickSize float64 `json:"orderPriceMinTickSize"`
	OrderMinSize          float64 `json:"orderMinSize"`
	MakerBaseFee          int     `json:"makerBaseFee"`
	TakerBaseFee          int     `json:"takerBaseFee"`
	AcceptingOrders       bool    `json:"acceptingOrders"`
	NotificationsEnabled  bool    `json:"notificationsEnabled"`

	// Curation and scoring
	CurationOrder int     `json:"curationOrder"`
	Score         float64 `json:"score"`

	// Review status
	HasReviewedDates bool `json:"hasReviewedDates"`
	ReadyForCron     bool `json:"readyForCron"`
	CommentsEnabled  bool `json:"commentsEnabled"`

	// Volume metrics
	Volume24hr     float64 `json:"volume24hr"`
	Volume1wk      float64 `json:"volume1wk"`
	Volume1mo      float64 `json:"volume1mo"`
	Volume1yr      float64 `json:"volume1yr"`
	Volume24hrAmm  float64 `json:"volume24hrAmm"`
	Volume1wkAmm   float64 `json:"volume1wkAmm"`
	Volume1moAmm   float64 `json:"volume1moAmm"`
	Volume1yrAmm   float64 `json:"volume1yrAmm"`
	Volume24hrClob float64 `json:"volume24hrClob"`
	Volume1wkClob  float64 `json:"volume1wkClob"`
	Volume1moClob  float64 `json:"volume1moClob"`
	Volume1yrClob  float64 `json:"volume1yrClob"`
	VolumeAmm      float64 `json:"volumeAmm"`
	VolumeClob     float64 `json:"volumeClob"`

	// Liquidity breakdown
	LiquidityAmm  float64 `json:"liquidityAmm"`
	LiquidityClob float64 `json:"liquidityClob"`

	// Gaming/sports specific
	GameStartTime    NormalizedTime `json:"gameStartTime"`
	SecondsDelay     int            `json:"secondsDelay"`
	ClobTokenIDs     string         `json:"clobTokenIds"`
	TeamAID          string         `json:"teamAID"`
	TeamBID          string         `json:"teamBID"`
	GameID           string         `json:"gameId"`
	SportsMarketType string         `json:"sportsMarketType"`
	Line             float64        `json:"line"`

	// Discussions
	DisqusThread string `json:"disqusThread"`

	// AMM status
	FPMMLive bool `json:"fpmmLive"`

	// Custom settings
	CustomLiveness int `json:"customLiveness"`

	// Rewards
	RewardsMinSize   float64 `json:"rewardsMinSize"`
	RewardsMaxSpread float64 `json:"rewardsMaxSpread"`

	// Image optimization
	ImageOptimized *ImageOptimized `json:"imageOptimized,omitempty"`
	IconOptimized  *ImageOptimized `json:"iconOptimized,omitempty"`

	// Related entities
	Events     []Event    `json:"events,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Tags       []Tag      `json:"tags,omitempty"`

	// Creator and metadata
	Creator                  string         `json:"creator"`
	PastSlugs                string         `json:"pastSlugs"`
	ReadyTimestamp           NormalizedTime `json:"readyTimestamp"`
	FundedTimestamp          NormalizedTime `json:"fundedTimestamp"`
	AcceptingOrdersTimestamp NormalizedTime `json:"acceptingOrdersTimestamp"`

	// Competition
	Competitive float64 `json:"competitive"`

	// Spread information
	Spread float64 `json:"spread"`

	// Resolution flags
	AutomaticallyResolved bool `json:"automaticallyResolved"`

	// Price changes
	OneDayPriceChange   float64 `json:"oneDayPriceChange"`
	OneHourPriceChange  float64 `json:"oneHourPriceChange"`
	OneWeekPriceChange  float64 `json:"oneWeekPriceChange"`
	OneMonthPriceChange float64 `json:"oneMonthPriceChange"`
	OneYearPriceChange  float64 `json:"oneYearPriceChange"`

	// Current prices
	LastTradePrice float64 `json:"lastTradePrice"`
	BestBid        float64 `json:"bestBid"`
	BestAsk        float64 `json:"bestAsk"`

	// Activation
	AutomaticallyActive bool `json:"automaticallyActive"`
	ClearBookOnStart    bool `json:"clearBookOnStart"`
	ManualActivation    bool `json:"manualActivation"`

	// Chart styling
	ChartColor     string `json:"chartColor"`
	SeriesColor    string `json:"seriesColor"`
	ShowGmpSeries  bool   `json:"showGmpSeries"`
	ShowGmpOutcome bool   `json:"showGmpOutcome"`

	// Negative risk
	NegRiskOther bool `json:"negRiskOther"`

	// Deployment status
	PendingDeployment            bool           `json:"pendingDeployment"`
	Deploying                    bool           `json:"deploying"`
	DeployingTimestamp           NormalizedTime `json:"deployingTimestamp"`
	ScheduledDeploymentTimestamp NormalizedTime `json:"scheduledDeploymentTimestamp"`

	// RFQ
	RFQEnabled bool `json:"rfqEnabled"`

	// Event timing (may be empty/null for backward compatibility)
	EventStartTime NormalizedTime `json:"eventStartTime"`
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
	ID                string         `json:"id"`
	Ticker            string         `json:"ticker"`
	Slug              string         `json:"slug"`
	Title             string         `json:"title"`
	Subtitle          string         `json:"subtitle"`
	Description       string         `json:"description"`
	ResolutionSource  string         `json:"resolutionSource"`
	StartDate         NormalizedTime `json:"startDate"`
	CreationDate      NormalizedTime `json:"creationDate"`
	EndDate           NormalizedTime `json:"endDate"`
	Image             string         `json:"image"`
	Icon              string         `json:"icon"`
	Active            bool           `json:"active"`
	Closed            bool           `json:"closed"`
	Archived          bool           `json:"archived"`
	New               bool           `json:"new"`
	Featured          bool           `json:"featured"`
	Restricted        bool           `json:"restricted"`
	Liquidity         float64        `json:"liquidity"`
	Volume            float64        `json:"volume"`
	OpenInterest      float64        `json:"openInterest"`
	SortBy            string         `json:"sortBy"`
	Category          string         `json:"category"`
	Subcategory       string         `json:"subcategory"`
	IsTemplate        bool           `json:"isTemplate"`
	TemplateVariables string         `json:"templateVariables"`
	PublishedAt       NormalizedTime `json:"published_at"`
	CreatedBy         string         `json:"createdBy"`
	UpdatedBy         string         `json:"updatedBy"`
	CreatedAt         NormalizedTime `json:"createdAt"`
	UpdatedAt         NormalizedTime `json:"updatedAt"`
	CommentsEnabled   bool           `json:"commentsEnabled"`
	Competitive       float64        `json:"competitive"`
	Volume24hr        float64        `json:"volume24hr"`
	Volume1wk         float64        `json:"volume1wk"`
	Volume1mo         float64        `json:"volume1mo"`
	Volume1yr         float64        `json:"volume1yr"`
	FeaturedImage     string         `json:"featuredImage"`
	DisqusThread      string         `json:"disqusThread"`
	ParentEvent       string         `json:"parentEvent"`
	EnableOrderBook   bool           `json:"enableOrderBook"`
	LiquidityAmm      float64        `json:"liquidityAmm"`
	LiquidityClob     float64        `json:"liquidityClob"`
	NegRisk           bool           `json:"negRisk"`
	NegRiskMarketID   string         `json:"negRiskMarketID"`
	NegRiskFeeBips    int            `json:"negRiskFeeBips"`
	CommentCount      int            `json:"commentCount"`

	// Optimized images
	ImageOptimized         *ImageOptimized `json:"imageOptimized,omitempty"`
	IconOptimized          *ImageOptimized `json:"iconOptimized,omitempty"`
	FeaturedImageOptimized *ImageOptimized `json:"featuredImageOptimized,omitempty"`

	// Related entities
	SubEvents   []string     `json:"subEvents,omitempty"`
	Markets     []Market     `json:"markets,omitempty"`
	Series      []Series     `json:"series,omitempty"`
	Categories  []Category   `json:"categories,omitempty"`
	Collections []Collection `json:"collections,omitempty"`
	Tags        []Tag        `json:"tags,omitempty"`

	// Additional fields
	CYOM                         bool           `json:"cyom"`
	ClosedTime                   NormalizedTime `json:"closedTime"`
	ShowAllOutcomes              bool           `json:"showAllOutcomes"`
	ShowMarketImages             bool           `json:"showMarketImages"`
	AutomaticallyResolved        bool           `json:"automaticallyResolved"`
	EnableNegRisk                bool           `json:"enableNegRisk"`
	AutomaticallyActive          bool           `json:"automaticallyActive"`
	EventDate                    NormalizedTime `json:"eventDate"`
	StartTime                    NormalizedTime `json:"startTime"`
	EventWeek                    int            `json:"eventWeek"`
	SeriesSlug                   string         `json:"seriesSlug"`
	Score                        string         `json:"score"`
	Elapsed                      string         `json:"elapsed"`
	Period                       string         `json:"period"`
	Live                         bool           `json:"live"`
	Ended                        bool           `json:"ended"`
	FinishedTimestamp            NormalizedTime `json:"finishedTimestamp"`
	GMPChartMode                 string         `json:"gmpChartMode"`
	EventCreators                []EventCreator `json:"eventCreators,omitempty"`
	TweetCount                   int            `json:"tweetCount"`
	Chats                        []Chat         `json:"chats,omitempty"`
	FeaturedOrder                int            `json:"featuredOrder"`
	EstimateValue                bool           `json:"estimateValue"`
	CantEstimate                 bool           `json:"cantEstimate"`
	EstimatedValue               string         `json:"estimatedValue"`
	Templates                    []Template     `json:"templates,omitempty"`
	SpreadsMainLine              float64        `json:"spreadsMainLine"`
	TotalsMainLine               float64        `json:"totalsMainLine"`
	CarouselMap                  string         `json:"carouselMap"`
	PendingDeployment            bool           `json:"pendingDeployment"`
	Deploying                    bool           `json:"deploying"`
	DeployingTimestamp           NormalizedTime `json:"deployingTimestamp"`
	ScheduledDeploymentTimestamp NormalizedTime `json:"scheduledDeploymentTimestamp"`
	GameStatus                   string         `json:"gameStatus"`
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
	ID          string         `json:"id"`
	Label       string         `json:"label"`
	Slug        string         `json:"slug"`
	ForceShow   bool           `json:"forceShow"`
	PublishedAt NormalizedTime `json:"publishedAt"`
	CreatedBy   int            `json:"createdBy"`
	UpdatedBy   int            `json:"updatedBy"`
	CreatedAt   NormalizedTime `json:"createdAt"`
	UpdatedAt   NormalizedTime `json:"updatedAt"`
	ForceHide   bool           `json:"forceHide"`
	IsCarousel  bool           `json:"isCarousel"`
}

// Series represents a series of events
type Series struct {
	ID                string         `json:"id"`
	Ticker            string         `json:"ticker"`
	Slug              string         `json:"slug"`
	Title             string         `json:"title"`
	Subtitle          string         `json:"subtitle"`
	SeriesType        string         `json:"seriesType"`
	Recurrence        string         `json:"recurrence"`
	Description       string         `json:"description"`
	Image             string         `json:"image"`
	Icon              string         `json:"icon"`
	Layout            string         `json:"layout"`
	Active            bool           `json:"active"`
	Closed            bool           `json:"closed"`
	Archived          bool           `json:"archived"`
	New               bool           `json:"new"`
	Featured          bool           `json:"featured"`
	Restricted        bool           `json:"restricted"`
	IsTemplate        bool           `json:"isTemplate"`
	TemplateVariables bool           `json:"templateVariables"`
	PublishedAt       NormalizedTime `json:"publishedAt"`
	CreatedBy         string         `json:"createdBy"`
	UpdatedBy         string         `json:"updatedBy"`
	CreatedAt         NormalizedTime `json:"createdAt"`
	UpdatedAt         NormalizedTime `json:"updatedAt"`
	CommentsEnabled   bool           `json:"commentsEnabled"`
	Competitive       string         `json:"competitive"`
	Volume24hr        float64        `json:"volume24hr"`
	Volume            float64        `json:"volume"`
	Liquidity         float64        `json:"liquidity"`
	StartDate         NormalizedTime `json:"startDate"`
	PythTokenID       string         `json:"pythTokenID"`
	CGAssetName       string         `json:"cgAssetName"`
	Score             float64        `json:"score"`
	Events            []Event        `json:"events,omitempty"`
	Collections       []Collection   `json:"collections,omitempty"`
	Categories        []Category     `json:"categories,omitempty"`
	Tags              []Tag          `json:"tags,omitempty"`
	CommentCount      int            `json:"commentCount"`
	Chats             []Chat         `json:"chats,omitempty"`
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
