package marketsurge_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/major/marketsurge-go/marketsurge"
)

func TestReportScreensReturnsExpectedCount(t *testing.T) {
	t.Parallel()

	const want = 50
	got := marketsurge.ReportScreens()
	if len(got) != want {
		t.Fatalf("ReportScreens() length = %d, want %d", len(got), want)
	}
}

func TestReportScreensReturnsCopy(t *testing.T) {
	t.Parallel()

	first := marketsurge.ReportScreens()
	second := marketsurge.ReportScreens()

	original := first[0]
	first[0].Name = "MUTATED"

	if second[0].Name == "MUTATED" {
		t.Fatal("ReportScreens() returned same backing array, not a copy")
	}

	third := marketsurge.ReportScreens()
	if diff := cmp.Diff(original, third[0]); diff != "" {
		t.Errorf("ReportScreens() mutation leaked (-want +got):\n%s", diff)
	}
}

func TestReportScreensNoDuplicateIDs(t *testing.T) {
	t.Parallel()

	screens := marketsurge.ReportScreens()
	seen := make(map[int]string, len(screens))
	for _, screen := range screens {
		if name, ok := seen[screen.ID]; ok {
			t.Errorf("ReportScreens() duplicate ID %d for %q and %q", screen.ID, name, screen.Name)
		}

		seen[screen.ID] = screen.Name
	}
}

func TestReportScreensAllHaveMetadata(t *testing.T) {
	t.Parallel()

	screens := marketsurge.ReportScreens()
	for i, screen := range screens {
		if screen.ID == 0 {
			t.Errorf("ReportScreens()[%d].ID = %d, want non-zero", i, screen.ID)
		}

		if screen.Name == "" {
			t.Errorf("ReportScreens()[%d].Name = %q, want non-empty", i, screen.Name)
		}

		if screen.Description == "" {
			t.Errorf("ReportScreens()[%d].Description = %q, want non-empty", i, screen.Description)
		}
	}
}

func TestReportScreensContainsKnownEntries(t *testing.T) {
	t.Parallel()

	screens := marketsurge.ReportScreens()
	byID := make(map[int]marketsurge.ReportScreenInfo, len(screens))
	for _, screen := range screens {
		byID[screen.ID] = screen
	}

	tests := []struct {
		name string
		want marketsurge.ReportScreenInfo
	}{
		{
			name: "IBD Big Cap 20",
			want: marketsurge.ReportScreenInfo{
				ID:          40,
				Name:        "IBD Big Cap 20",
				Description: "The IBD Big Cap 20 is a computer-generated ranking of the leading large-capitalization companies trading in the U.S. The list is featured in Investor's Business Daily every Tuesday. Rankings are based on a combination of each company's profit growth; IBD's Composite Rating, which includes key measures such as return on equity, sales growth and profit margins; and relative price strength in the past 12 months. Some of the companies also appear on the IBD 100, which features a blend of small- mid- and large-cap names.    For investors who prefer less volatile investments with greater liquidity, the Big Cap 20 offers more mature yet still vibrant enterprises. Companies must have a minimum market value of $15 billion and trade more than 300,000 shares a day.",
			},
		},
		{
			name: "Minervini Trend - 5 Months",
			want: marketsurge.ReportScreenInfo{
				ID:          120,
				Name:        "Minervini Trend - 5 Months",
				Description: "Developed by MarketSurge partner Mark Minervini, the Trend Template lists screen for stocks with moving averages in an uptrend.  The 1-Month version requires the 200-day moving average to be trending up for 1 month while the 5-Month list requires the 200-day moving average to be trending up for 5 months.",
			},
		},
		{
			name: "Top 25 Industry or Sector Funds Over 3 Years",
			want: marketsurge.ReportScreenInfo{
				ID:          51,
				Name:        "Top 25 Industry or Sector Funds Over 3 Years",
				Description: "Snapshot of top 25 Industry or Sector Funds over the last 3 years. The list is sorted by performance in descending order.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := byID[tt.want.ID]
			if !ok {
				t.Fatalf("ReportScreens() missing ID %d for %q", tt.want.ID, tt.want.Name)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ReportScreens() entry %d mismatch (-want +got):\n%s", tt.want.ID, diff)
			}
		})
	}
}
