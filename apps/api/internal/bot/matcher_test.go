package bot

import (
	"strings"
	"testing"
)

func TestMatcher(t *testing.T) {
	matcher := NewMatcher(GetKnowledgeBase())

	tests := []struct {
		name         string
		input        string
		wantCategory string
	}{
		// Price Information
		{
			name:         "Exact keyword harga",
			input:        "harga",
			wantCategory: CategoryPrice,
		},
		{
			name:         "Sentence asking for price",
			input:        "Halo admin, berapa tarif cuci kiloan per kg ya?",
			wantCategory: CategoryPrice,
		},
		{
			name:         "Price asking with slang and punctuation",
			input:        "Berapaan ongkos cucian disini kak???",
			wantCategory: CategoryPrice,
		},

		// Service Information
		{
			name:         "Exact keyword layanan",
			input:        "layanan",
			wantCategory: CategoryService,
		},
		{
			name:         "Multi-word phrase dry clean",
			input:        "Apakah ada paket dry clean untuk jas?",
			wantCategory: CategoryService,
		},
		{
			name:         "Asking for express service",
			input:        "Bisa cuci kilat express selesai 6 jam gak?",
			wantCategory: CategoryService,
		},

		// Operating Hours
		{
			name:         "Exact keyword jam",
			input:        "jam",
			wantCategory: CategoryHours,
		},
		{
			name:         "Asking for operational hours",
			input:        "Toko buka jam berapa dan tutup jam berapa?",
			wantCategory: CategoryHours,
		},
		{
			name:         "Asking if open today",
			input:        "Hari Minggu gini outlet buka gak ya?",
			wantCategory: CategoryHours,
		},

		// Address & Location
		{
			name:         "Exact keyword alamat",
			input:        "alamat",
			wantCategory: CategoryAddress,
		},
		{
			name:         "Asking for location and maps",
			input:        "Bisa minta shareloc atau google maps outletnya dimana?",
			wantCategory: CategoryAddress,
		},
		{
			name:         "Asking where is the store",
			input:        "Posisi tokonya sebelah mana ya patokannya?",
			wantCategory: CategoryAddress,
		},

		// Delivery Coverage
		{
			name:         "Exact keyword area",
			input:        "area",
			wantCategory: CategoryCoverage,
		},
		{
			name:         "Asking for delivery coverage",
			input:        "Apakah wilayah Senopati masuk jangkauan antar jemput?",
			wantCategory: CategoryCoverage,
		},
		{
			name:         "Asking for free shipping",
			input:        "Ada syarat free ongkir untuk jemput cucian?",
			wantCategory: CategoryCoverage,
		},

		// Pickup Information
		{
			name:         "Exact keyword pickup",
			input:        "pickup",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Request pickup",
			input:        "Kak saya mau minta tolong jemput cucian di kosan",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Asking how to pickup",
			input:        "Bagaimana cara jemput pakaian kotor?",
			wantCategory: CategoryPickup,
		},

		// Help / Greetings
		{
			name:         "Greeting halo",
			input:        "Halo",
			wantCategory: CategoryHelp,
		},
		{
			name:         "Greeting P",
			input:        "p",
			wantCategory: CategoryHelp,
		},
		{
			name:         "Command menu",
			input:        "menu",
			wantCategory: CategoryHelp,
		},

		// Unknown Queries
		{
			name:         "Gibberish text",
			input:        "xyz123abc random text",
			wantCategory: CategoryUnknown,
		},
		{
			name:         "Empty string",
			input:        "",
			wantCategory: CategoryUnknown,
		},
		{
			name:         "Punctuation only",
			input:        "??? !!! ...",
			wantCategory: CategoryUnknown,
		},
		{
			name:         "Irrelevant question",
			input:        "Apakah jual pulsa kuota telkomsel?",
			wantCategory: CategoryUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat, _ := matcher.Match(tt.input)
			if cat != tt.wantCategory {
				t.Errorf("Match(%q) = %v; want %v", tt.input, cat, tt.wantCategory)
			}
		})
	}
}

func TestFormatResponse(t *testing.T) {
	categories := []string{
		CategoryPrice,
		CategoryService,
		CategoryHours,
		CategoryAddress,
		CategoryCoverage,
		CategoryPickup,
		CategoryHelp,
		CategoryUnknown,
	}

	for _, cat := range categories {
		t.Run(cat, func(t *testing.T) {
			res := FormatResponse(cat, "Budi")
			if strings.TrimSpace(res) == "" {
				t.Fatalf("FormatResponse(%q) returned empty response", cat)
			}
			if !strings.Contains(res, "LaundryOS") && !strings.Contains(res, "LAUNDRYOS") {
				t.Errorf("FormatResponse(%q) missing LaundryOS branding", cat)
			}
		})
	}
}
