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

// GetKnowledgeBase returns the default FAQ knowledge base for FAN Laundry.
func GetKnowledgeBase() []FAQItem {
	return []FAQItem{
		{
			Category: CategoryPrice,
			Title:    "Daftar Harga & Tarif",
			Keywords: []string{
				"1", "1️⃣", "satu", "menu 1", "pilihan 1",
				"harga", "tarif", "biaya", "biayanya", "ongkos", "pricelist", "price",
				"kiloan", "satuan", "berapaan", "per kilo", "biaya cuci",
				"daftar harga", "liat harga", "lihat harga", "minta pricelist",
				"info harga", "info layanan dan harga", "info layanan & harga",
			},
			Description: "Informasi lengkap tarif kiloan, satuan, bed cover, dan express.",
		},
		{
			Category: CategoryService,
			Title:    "Layanan Laundry",
			Keywords: []string{
				"layanan", "service", "cuci", "setrika", "dry clean", "express",
				"kilat", "reguler", "apa saja", "jenis", "paket", "paket laundry", "jenis layanan",
			},
			Description: "Informasi jenis layanan yang tersedia di FAN Laundry.",
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
				"3", "3️⃣", "tiga", "menu 3", "pilihan 3",
				"alamat", "lokasi", "tempat", "posisi", "dimana", "outlet",
				"toko", "cabang", "maps", "gmaps", "patokan",
				"alamat dan jam operasional", "alamat & jam operasional",
			},
			Description: "Informasi lokasi fisik toko dan tautan Google Maps.",
		},
		{
			Category: CategoryCoverage,
			Title:    "Area Jangkauan, Promo & Syarat",
			Keywords: []string{
				"4", "4️⃣", "empat", "menu 4", "pilihan 4",
				"promo", "syarat", "promo & syarat", "promo dan syarat", "syarat dan ketentuan",
				"area", "jangkauan", "antar", "coverage", "daerah",
				"wilayah", "bisa antar", "radius", "ongkir", "free ongkir", "gratis ongkir",
				"bisa ambil di daerah", "bisa ambil di", "antar jemput ke",
			},
			Description: "Wilayah operasional antar-jemput laundry dan syarat gratis ongkir.",
		},
		{
			Category: CategoryPickup,
			Title:    "Informasi & Cara Request Pickup",
			Keywords: []string{
				"2", "2️⃣", "dua", "menu 2", "pilihan 2",
				"order", "penjemputan", "penjemputan laundry", "order / penjemputan laundry",
				"pickup", "jemput", "ambil", "request pickup", "minta jemput",
				"ambil cucian", "cara jemput", "mau dijemput", "jemput pakaian", "jemput baju",
				"mau laundry", "saya mau laundry", "mau laundry mba", "mau laundry donk",
				"mba laundry donk", "mba saya mau laundry", "mba ambil laundry", "ambil laundry",
				"ambil laundry mba", "ambil baju", "neng kerumah", "kerumah ya",
				"kerumah ya ambil baju", "neng kerumah ya ambil baju",
				"umi baju kotor", "baju kotor", "baju kotor yaa udh disiapin", "udh disiapin",
				"sudah disiapkan", "mau nyuci", "mau nyuci mba", "mau cuci", "bisa laundry",
				"bisa ambil laundry", "order laundry", "pesan jemput",
			},
			Description: "Tata cara meminta penjemputan cucian langsung ke rumah/kos.",
		},
		{
			Category: CategoryStatus,
			Title:    "Cek Status Cucian & Jadwal Selesai",
			Keywords: []string{
				"status", "cek status", "cek cucian", "status cucian", "selesai",
				"kapan selesai", "sudah selesai", "selesai belum", "progres", "lacak",
				"tracking", "nota", "jadwal selesai",
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
			Description: "Daftar pilihan informasi FAQ FAN Laundry.",
		},
	}
}
