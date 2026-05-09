package marketsurge_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/major/marketsurge-go/marketsurge"
)

func TestColumnsReturnsNonEmpty(t *testing.T) {
	t.Parallel()

	cols := marketsurge.Columns()
	if len(cols) == 0 {
		t.Fatal("Columns() returned empty slice")
	}
}

func TestColumnsReturnsCopy(t *testing.T) {
	t.Parallel()

	first := marketsurge.Columns()
	second := marketsurge.Columns()

	// Mutate first slice and verify second is unaffected.
	original := first[0]
	first[0].Name = "MUTATED"

	if second[0].Name == "MUTATED" {
		t.Fatal("Columns() returned same backing array, not a copy")
	}

	// Verify original value is still intact in a fresh call.
	third := marketsurge.Columns()
	if third[0].Name != original.Name {
		t.Errorf("mutation leaked: got %q, want %q", third[0].Name, original.Name)
	}
}

func TestColumnsConsistentLength(t *testing.T) {
	t.Parallel()

	first := marketsurge.Columns()
	second := marketsurge.Columns()

	if len(first) != len(second) {
		t.Errorf("inconsistent length: first=%d, second=%d", len(first), len(second))
	}
}

func TestColumnsNoDuplicateNames(t *testing.T) {
	t.Parallel()

	cols := marketsurge.Columns()
	seen := make(map[string]bool, len(cols))

	for _, col := range cols {
		if seen[col.Name] {
			t.Errorf("duplicate wire name: %q", col.Name)
		}

		seen[col.Name] = true
	}
}

func TestColumnsAllHaveNames(t *testing.T) {
	t.Parallel()

	cols := marketsurge.Columns()
	for i, col := range cols {
		if col.Name == "" {
			t.Errorf("entry %d has empty Name", i)
		}

		if col.DisplayName == "" {
			t.Errorf("entry %d (%s) has empty DisplayName", i, col.Name)
		}
	}
}

func TestColumnsContainsKnownEntries(t *testing.T) {
	t.Parallel()

	cols := marketsurge.Columns()
	byName := make(map[string]marketsurge.ColumnInfo, len(cols))
	for _, col := range cols {
		byName[col.Name] = col
	}

	tests := []struct {
		name     string
		wantInfo marketsurge.ColumnInfo
	}{
		{
			name: "EPSRating constant",
			wantInfo: marketsurge.ColumnInfo{
				Name:        marketsurge.ColumnEPSRating,
				DisplayName: "EPS Rating",
				Category:    "SmartSelect Rating",
			},
		},
		{
			name: "Symbol constant",
			wantInfo: marketsurge.ColumnInfo{
				Name:        marketsurge.ColumnSymbol,
				DisplayName: "Symbol",
			},
		},
		{
			name: "CompanyName constant",
			wantInfo: marketsurge.ColumnInfo{
				Name:        marketsurge.ColumnCompanyName,
				DisplayName: "Name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := byName[tt.wantInfo.Name]
			if !ok {
				t.Fatalf("column %q not found in catalog", tt.wantInfo.Name)
			}

			if got.Name != tt.wantInfo.Name {
				t.Errorf("Name: got %q, want %q", got.Name, tt.wantInfo.Name)
			}

			if got.DisplayName != tt.wantInfo.DisplayName {
				t.Errorf("DisplayName: got %q, want %q", got.DisplayName, tt.wantInfo.DisplayName)
			}

			if tt.wantInfo.Description != "" && got.Description != tt.wantInfo.Description {
				t.Errorf("Description: got %q, want %q", got.Description, tt.wantInfo.Description)
			}

			if tt.wantInfo.Category != "" && got.Category != tt.wantInfo.Category {
				t.Errorf("Category: got %q, want %q", got.Category, tt.wantInfo.Category)
			}
		})
	}
}

func TestColumnConstantsMatchCatalog(t *testing.T) {
	t.Parallel()

	cols := marketsurge.Columns()
	byName := make(map[string]bool, len(cols))
	for _, col := range cols {
		byName[col.Name] = true
	}

	// Verify a representative set of constants appear in the catalog.
	// These span different categories to catch structural issues.
	constants := []string{
		marketsurge.ColumnEPSRating,
		marketsurge.ColumnRSRating,
		marketsurge.ColumnSMRRating,
		marketsurge.ColumnAccDisRating,
		marketsurge.ColumnCompositeRating,
		marketsurge.ColumnSymbol,
		marketsurge.ColumnCompanyName,
		marketsurge.ColumnPrice,
		marketsurge.ColumnVolume,
		marketsurge.ColumnIndustryName,
		marketsurge.ColumnExchange,
		marketsurge.ColumnFundYTDReturn,
		marketsurge.ColumnAlertPrice,
		marketsurge.ColumnIBD50Flag,
	}

	for _, c := range constants {
		if !byName[c] {
			t.Errorf("constant %q not found in Columns() catalog", c)
		}
	}
}

func TestColumnsSmartSelectHasDescription(t *testing.T) {
	t.Parallel()

	cols := marketsurge.Columns()
	for _, col := range cols {
		if col.Name == marketsurge.ColumnEPSRating {
			if col.Description == "" {
				t.Error("EPSRating entry has empty Description")
			}

			return
		}
	}

	t.Fatal("EPSRating not found in catalog")
}

func TestColumnsIncludesSpaceNames(t *testing.T) {
	t.Parallel()

	// Wire names with spaces can't be Go constants but must appear
	// in the catalog for completeness.
	cols := marketsurge.Columns()
	byName := make(map[string]bool, len(cols))
	for _, col := range cols {
		byName[col.Name] = true
	}

	spaceNames := []string{
		"Timeliness Rating",
		"Options Exchange",
		"Add To List Date",
	}

	for _, sn := range spaceNames {
		if !byName[sn] {
			t.Errorf("space-name column %q not found in catalog", sn)
		}
	}
}

func TestColumnInfoStruct(t *testing.T) {
	t.Parallel()

	info := marketsurge.ColumnInfo{
		Name:        "TestCol",
		DisplayName: "Test Column",
		Description: "A test column.",
		Category:    "Testing",
	}

	want := marketsurge.ColumnInfo{
		Name:        "TestCol",
		DisplayName: "Test Column",
		Description: "A test column.",
		Category:    "Testing",
	}

	if diff := cmp.Diff(want, info); diff != "" {
		t.Errorf("ColumnInfo mismatch (-want +got):\n%s", diff)
	}
}
