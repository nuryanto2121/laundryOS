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

		// Real Customer Everyday Phrases
		{
			name:         "Mba saya mau laundry",
			input:        "mba saya mau laundry",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Mba ambil laundry",
			input:        "mba ambil laundry",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Neng kerumah ya ambil baju",
			input:        "Neng kerumah ya ambil baju",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Umi baju kotor yaa udh disiapin",
			input:        "Umi baju kotor yaa udh disiapin",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Mba laundry donk",
			input:        "mba laundry donk",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Numeric menu 1",
			input:        "1",
			wantCategory: CategoryPrice,
		},
		{
			name:         "Numeric menu 2",
			input:        "2",
			wantCategory: CategoryPickup,
		},
		{
			name:         "Numeric menu 3",
			input:        "3",
			wantCategory: CategoryAddress,
		},
		{
			name:         "Numeric menu 4",
			input:        "4",
			wantCategory: CategoryCoverage,
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
			if !strings.Contains(res, "FAN Laundry") && !strings.Contains(res, "FAN LAUNDRY") {
				t.Errorf("FormatResponse(%q) missing FAN Laundry branding", cat)
			}
		})
	}
}

func TestExtractCustomerName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"saya nana teman'a nini mba, boleh liat daftar harga?", "Nana"},
		{"nama saya budi santoso", "Budi"},
		{"halo aku rina dari blok b", "Rina"},
		{"saya mau laundry mba", ""}, // "mau" is a stopword
		{"harganya berapa?", ""},
	}

	for _, tt := range tests {
		got := ExtractCustomerName(tt.input)
		if got != tt.expected {
			t.Errorf("ExtractCustomerName(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestMatchServiceIntent(t *testing.T) {
	tests := []struct {
		input             string
		expectedCanonical string
	}{
		{"cuci regular iah", "Cuci Komplit Reguler"},
		{"mau yang reguler mba", "Cuci Komplit Reguler"},
		{"ambil paket kilat ya", "Cuci Setrika Kilat"},
		{"tolong cuci express ya", "Cuci Setrika Express"},
		{"mau cuci turbo 6 jam", "Cuci Setrika Turbo"},
		{"setrika saja ya mba", "Setrika Saja"},
		{"ada cuci bed cover?", "Cuci Bed Cover"},
		{"mau cuci bad cover", "Cuci Bed Cover"},
		{"bisa cuci karpet permadani?", "Karpet Permadani"},
		{"cuci boneka besar", "Cuci Boneka"},
	}

	for _, tt := range tests {
		matched := MatchServiceIntent(tt.input, nil)
		if matched == nil {
			t.Fatalf("MatchServiceIntent(%q) returned nil", tt.input)
		}
		if matched.ServiceName != tt.expectedCanonical {
			t.Errorf("MatchServiceIntent(%q) = %q; want %q", tt.input, matched.ServiceName, tt.expectedCanonical)
		}
	}
}

func TestExtractOrderNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"• *No. Order:* #261118000001", "ORD261118000001"},
		{"timbang #261119000012 9.5", "ORD261119000012"},
		{"saya ambil ORD261118000001 ya kak", "ORD261118000001"},
		{"otw ORD261118123456", "ORD261118123456"},
		{"selesai #261119000015", "ORD261119000015"},
		{"261119000016", "ORD261119000016"},
		{"halo apa kabar?", ""},
		{"ORD123", ""},
	}

	for _, tt := range tests {
		got := ExtractOrderNumber(tt.input)
		if got != tt.expected {
			t.Errorf("ExtractOrderNumber(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestParseWeightCommand(t *testing.T) {
	tests := []struct {
		input       string
		wantWeight  float64
		wantSuccess bool
	}{
		{"TIMBANG #261119000012 9.5", 9.5, true},
		{"timbang 5,5 kg", 5.5, true},
		{"berat 8.2", 8.2, true},
		{"12.4 kg", 12.4, true},
		{"hanya teks biasa", 0, false},
	}

	for _, tt := range tests {
		gotWeight, gotSuccess := ParseWeightCommand(tt.input)
		if gotSuccess != tt.wantSuccess || gotWeight != tt.wantWeight {
			t.Errorf("ParseWeightCommand(%q) = (%v, %v); want (%v, %v)", tt.input, gotWeight, gotSuccess, tt.wantWeight, tt.wantSuccess)
		}
	}
}

func TestCommandRecognizers(t *testing.T) {
	if !IsWeighCommand("TIMBANG #261119000012 9.5") {
		t.Errorf("expected IsWeighCommand to be true")
	}
	if !IsWeighCommand("berat cucian 5kg") {
		t.Errorf("expected IsWeighCommand to be true")
	}
	if !IsCompleteCommand("SELESAI #261119000012") {
		t.Errorf("expected IsCompleteCommand to be true")
	}
	if !IsCompleteCommand("order done ya") {
		t.Errorf("expected IsCompleteCommand to be true")
	}
}
