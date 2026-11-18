package bot

// FAQ Categories supported in Sprint 03 & Sprint 04
const (
	CategoryPrice    = "PRICE"
	CategoryService  = "SERVICE"
	CategoryHours    = "HOURS"
	CategoryAddress  = "ADDRESS"
	CategoryCoverage = "COVERAGE"
	CategoryPickup   = "PICKUP"
	CategoryStatus   = "STATUS"
	CategoryHelp     = "HELP"
	CategoryUnknown  = "UNKNOWN"
)

// FAQItem represents an entry in the FAQ knowledge base.
type FAQItem struct {
	Category    string
	Title       string
	Keywords    []string
	Description string
}

// GetKnowledgeBase returns the default FAQ knowledge base for LaundryOS.
func GetKnowledgeBase() []FAQItem {
	return []FAQItem{
		{
			Category: CategoryPrice,
			Title:    "Daftar Harga & Tarif",
			Keywords: []string{
				"harga", "tarif", "biaya", "biayanya", "ongkos", "pricelist", "price",
				"kiloan", "satuan", "berapaan", "per kilo", "biaya cuci",
			},
			Description: "Informasi lengkap tarif kiloan, satuan, bed cover, dan express.",
		},
		{
			Category: CategoryService,
			Title:    "Layanan Laundry",
			Keywords: []string{
				"layanan", "service", "cuci", "setrika", "dry clean", "express",
				"kilat", "reguler", "apa saja", "jenis", "paket", "laundry",
			},
			Description: "Informasi jenis layanan yang tersedia di LaundryOS.",
		},
		{
			Category: CategoryHours,
			Title:    "Jam Operasional",
			Keywords: []string{
				"jam", "buka", "tutup", "operasional", "jadwal", "hari apa", "kapan buka",
				"waktu", "buka jam", "tutup jam", "masih buka",
			},
			Description: "Informasi hari dan jam operasional outlet serta customer support.",
		},
		{
			Category: CategoryAddress,
			Title:    "Alamat & Lokasi Outlet",
			Keywords: []string{
				"alamat", "lokasi", "tempat", "posisi", "dimana", "outlet",
				"toko", "cabang", "maps", "gmaps", "patokan",
			},
			Description: "Informasi lokasi fisik toko dan tautan Google Maps.",
		},
		{
			Category: CategoryCoverage,
			Title:    "Area Jangkauan Antar-Jemput",
			Keywords: []string{
				"area", "jangkauan", "antar", "coverage", "daerah",
				"wilayah", "bisa antar", "radius", "ongkir", "free ongkir",
			},
			Description: "Wilayah operasional antar-jemput laundry dan syarat gratis ongkir.",
		},
		{
			Category: CategoryPickup,
			Title:    "Informasi & Cara Request Pickup",
			Keywords: []string{
				"pickup", "jemput", "ambil", "request pickup", "minta jemput",
				"ambil cucian", "cara jemput", "mau dijemput", "jemput pakaian",
			},
			Description: "Tata cara meminta penjemputan cucian langsung ke rumah/kos.",
		},
		{
			Category: CategoryStatus,
			Title:    "Cek Status Cucian & Jadwal Selesai",
			Keywords: []string{
				"status", "cek status", "cek cucian", "status cucian", "selesai",
				"kapan selesai", "sudah selesai", "selesai belum", "progres", "lacak",
				"tracking", "nota", "order", "jadwal selesai",
			},
			Description: "Informasi status pengerjaan cucian dan estimasi waktu selesai.",
		},
		{
			Category: CategoryHelp,
			Title:    "Menu Bantuan",
			Keywords: []string{
				"halo", "hai", "hello", "hi", "p", "menu", "bantuan", "help",
				"info", "bot", "start", "mulai", "assalamualaikum",
			},
			Description: "Daftar pilihan informasi FAQ LaundryOS.",
		},
	}
}
