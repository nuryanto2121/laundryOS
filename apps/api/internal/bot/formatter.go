package bot

import (
	"fmt"
	"strings"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// FormatResponse formats the response text based on the detected FAQ category.
func FormatResponse(category string, pushName string) string {
	return FormatResponseWithServices(category, pushName, nil)
}

// FormatResponseWithServices formats the response text, using dynamic database services when available.
func FormatResponseWithServices(category string, pushName string, services []model.ServiceItem) string {
	nameGreeting := ""
	if strings.TrimSpace(pushName) != "" {
		nameGreeting = fmt.Sprintf(" Kak *%s*", strings.TrimSpace(pushName))
	}

	switch category {
	case CategoryPrice:
		if len(services) > 0 {
			return formatDynamicPriceResponse(services)
		}
		return formatPriceResponse()
	case CategoryService:
		return formatServiceResponse()
	case CategoryHours:
		return formatHoursResponse()
	case CategoryAddress:
		return formatAddressResponse()
	case CategoryCoverage:
		return formatCoverageResponse()
	case CategoryPickup:
		return formatPickupResponse()
	case CategoryStatus:
		return formatHelpResponse(nameGreeting)
	case CategoryHelp:
		return formatHelpResponse(nameGreeting)
	default:
		return formatUnknownResponse(nameGreeting)
	}
}

func formatDynamicPriceResponse(services []model.ServiceItem) string {
	var sb strings.Builder
	sb.WriteString("🧺 *DAFTAR HARGA LAUNDRYOS (UPDATE RESMI)* 🧺\n\n")
	sb.WriteString("Berikut adalah tarif layanan laundry aktif kami saat ini:\n\n")

	for _, s := range services {
		durationStr := ""
		if s.DurationDay > 0 {
			durationStr = fmt.Sprintf(" (Selesai %d hari)", s.DurationDay)
		}
		sb.WriteString(fmt.Sprintf("• *%s*: Rp %s / Kg%s\n", s.ServiceName, formatIDR(s.PricePerKg), durationStr))
	}

	sb.WriteString("\n*Layanan Satuan & Bedding:*\n")
	sb.WriteString("• Bed Cover Single: *Rp 20.000 - Rp 25.000*\n")
	sb.WriteString("• Bed Cover King/Queen: *Rp 30.000 - Rp 35.000*\n")
	sb.WriteString("• Selimut: *Rp 15.000 - Rp 20.000*\n")
	sb.WriteString("• Jas / Blazer: *Rp 25.000*\n")
	sb.WriteString("• Sepatu: *Rp 25.000 - Rp 35.000 / pasang*\n\n")
	sb.WriteString("_Ketik *PICKUP* untuk pesan jemput cucian atau *MENU* untuk pilihan lainnya._")

	return sb.String()
}

func formatIDR(amount float64) string {
	n := int64(amount)
	in := fmt.Sprintf("%d", n)
	if len(in) <= 3 {
		return in
	}
	var res []string
	for len(in) > 3 {
		res = append([]string{in[len(in)-3:]}, res...)
		in = in[:len(in)-3]
	}
	if len(in) > 0 {
		res = append([]string{in}, res...)
	}
	return strings.Join(res, ".")
}

func formatPriceResponse() string {
	return `🧺 *DAFTAR HARGA LAUNDRYOS* 🧺

Berikut adalah daftar harga layanan laundry kami:

*1. Cuci Kiloan (Min. 3 Kg)*
• Reguler (2 Hari): *Rp 7.000 / Kg* (Cuci + Kering + Setrika Rapi)
• Kilat (24 Jam): *Rp 10.000 / Kg*
• Express (6 Jam): *Rp 15.000 / Kg*
• Cuci Kering Lipat: *Rp 5.000 / Kg*
• Setrika Saja: *Rp 4.000 / Kg*

*2. Satuan & Bedding*
• Bed Cover Single: *Rp 20.000 - Rp 25.000*
• Bed Cover King/Queen: *Rp 30.000 - Rp 35.000*
• Selimut: *Rp 15.000 - Rp 20.000*
• Jas / Blazer: *Rp 25.000*
• Sepatu: *Rp 25.000 - Rp 35.000 / pasang*

_Ketik *PICKUP* untuk pesan jemput cucian atau *MENU* untuk pilihan lainnya._`
}

func formatServiceResponse() string {
	return `🧺 *LAYANAN LAUNDRYOS* 🧺

Kami melayani berbagai kebutuhan cucian Anda dengan deterjen higienis dan pewangi premium:

✨ *Cuci Komplit Kiloan*
Pencucian higienis 1 mesin 1 pelanggan (tidak dicampur), dikeringkan mesin, dan disetrika uap wangi rapi.

⚡ *Pilihan Kecepatan:*
• *Reguler*: Selesai 2 hari kerja
• *Kilat (24 Jam)*: Cocok untuk pakaian kerja mendesak
• *Express (6 Jam)*: Siap pakai di hari yang sama

👔 *Layanan Khusus / Satuan:*
• Cuci & Treatment Sepatu (Deep Clean)
• Cuci Bed Cover, Selimut, Sprei, Gorden
• Dry Clean Jas, Gaun, Jaket Kulit / Parasut

_Ketik *HARGA* untuk daftar tarif atau *PICKUP* untuk minta jemput cucian._`
}

func formatHoursResponse() string {
	return `⏰ *JAM OPERASIONAL LAUNDRYOS* ⏰

Kami buka setiap hari untuk melayani kebutuhan cucian Anda:

📅 *Hari:* Senin - Minggu (Buka Setiap Hari)
🕐 *Outlet Fisik:* 07.00 - 21.00 WIB
🚚 *Layanan Antar-Jemput:* 08.00 - 19.00 WIB
💬 *Customer Service WhatsApp:* 07.00 - 22.00 WIB

_Pemesanan pickup di luar jam operasional akan dijadwalkan pada pagi hari berikutnya._

_Ketik *ALAMAT* untuk lokasi kami atau *PICKUP* untuk request jemput._`
}

func formatAddressResponse() string {
	return `📍 *LOKASI OUTLET LAUNDRYOS* 📍

Silakan kunjungi outlet kami di:

🏠 *Alamat:*
Jl. Merdeka No. 45, Kebayoran Baru, Jakarta Selatan, 12160
(Patokan: 50m sebelah Indomaret Point, seberang Apotek Sehat)

🗺️ *Google Maps:*
https://maps.google.com/?q=LaundryOS+Outlet

🚗 Tersedia area parkir untuk motor & mobil.

_Ketik *JAM* untuk melihat jam buka atau *AREA* untuk info antar-jemput._`
}

func formatCoverageResponse() string {
	return `🚚 *AREA JANGKAUAN ANTAR-JEMPUT* 🚚

Layanan pickup & delivery LaundryOS mencakup area berikut:

📍 *Wilayah Jangkauan:*
• Kebayoran Baru, Gandaria, Blok M
• Senopati, SCBD, Senayan
• Mampang Prapatan, Bangka, Kemang
• Pondok Indah & Fatmawati Utara

💰 *Ketentuan Ongkir:*
• *GRATIS ONGKIR* untuk minimal cucian 5 Kg (radius s/d 5 Km)
• Di bawah 5 Kg atau jarak 5-10 Km: Biaya flat *Rp 5.000 - Rp 10.000*

_Ketik *PICKUP* untuk memesan jemput cucian Anda sekarang._`
}

func formatPickupResponse() string {
	return `🛵 *CARA REQUEST PICKUP LAUNDRYOS* 🛵

Ingin cucian dijemput kurir kami langsung ke rumah/kos/kantor? Sangat mudah!

Silakan kirim pesan dengan format berikut:

*Nama:* [Nama Anda]
*Alamat Lengkap:* [Jalan, No Rumah, Kos/Apartemen, Patokan]
*Estimasi Berat / Item:* [Contoh: 5 kg kemeja & celana]
*Pilihan Layanan:* [Reguler / Kilat 24 Jam / Express]
*Waktu Pickup:* [Sekarang / Jam yang diinginkan]

Kurir kami akan segera mengonfirmasi dan meluncur ke lokasi Anda! 🛵💨

_Ketik *HARGA* untuk tarif atau *AREA* untuk cek jangkauan pengantaran._`
}

func formatHelpResponse(nameGreeting string) string {
	return fmt.Sprintf(`🧺 *Halo%s! Selamat Datang di LaundryOS* 👋

Ada yang bisa kami bantu? Anda dapat mengetik salah satu kata kunci di bawah ini untuk info cepat:

• *HARGA* ➔ Info daftar tarif kiloan & satuan
• *LAYANAN* ➔ Pilihan paket cuci (Reguler, Kilat, Express)
• *JAM* ➔ Jam operasional outlet & jadwal pickup
• *ALAMAT* ➔ Lokasi outlet fisik & link Google Maps
• *AREA* ➔ Wilayah jangkauan antar-jemput & ongkir
• *PICKUP* ➔ Format dan cara memesan jemput cucian

Silakan ketik pilihan Anda! 🙏`, nameGreeting)
}

func formatUnknownResponse(nameGreeting string) string {
	return fmt.Sprintf(`🧺 *Halo%s dari LaundryOS!* 👋

Maaf, kami belum memahami pertanyaan Anda. Untuk respon otomatis yang cepat, silakan ketik salah satu pilihan menu berikut:

• *HARGA* ➔ Daftar tarif kiloan & satuan
• *LAYANAN* ➔ Jenis layanan cuci & paket
• *JAM* ➔ Jam buka outlet kami
• *ALAMAT* ➔ Lokasi outlet LaundryOS
• *AREA* ➔ Jangkauan wilayah antar-jemput
• *PICKUP* ➔ Cara request jemput cucian

Staff kami juga akan meninjau pesan Anda secepatnya jika membutuhkan bantuan khusus. Terima kasih! 🙏`, nameGreeting)
}

// FormatOrderStatusResponse formats the tracking message for a customer's laundry order.
func FormatOrderStatusResponse(order *model.OrderTrackingDetail, pushName string) string {
	name := strings.TrimSpace(pushName)
	if name == "" && order != nil && order.CustomerName != "" {
		name = order.CustomerName
	}
	greeting := ""
	if name != "" {
		greeting = fmt.Sprintf(" Kak *%s*", name)
	}

	if order == nil || order.OrderNumber == "" {
		return fmt.Sprintf(`🧺 *CEK STATUS CUCIAN LAUNDRYOS* 🧺

Halo%s! Saat ini tidak ditemukan transaksi cucian aktif atas nomor WhatsApp Anda.

• Jika baru saja menitipkan pakaian, staf kami sedang memproses penimbangan & antrean.
• Anda juga dapat mengecek dengan mengetik nomor nota secara manual (contoh: *STATUS ORD261118000001*).
• Ketik *PICKUP* jika ingin meminta penjemputan pakaian kotor. 🙏`, greeting)
	}

	statusDesc := humanizeStatus(order.Status)

	weightStr := "-"
	if order.TotalWeight > 0 {
		weightStr = fmt.Sprintf("%.2f Kg", order.TotalWeight)
	}

	priceStr := "-"
	if order.TotalPrice > 0 {
		priceStr = fmt.Sprintf("Rp %s", formatIDR(order.TotalPrice))
	}

	duration := order.DurationDay
	if duration <= 0 {
		duration = 2
	}

	return fmt.Sprintf(`🧺 *STATUS CUCIAN ANDA (LaundryOS)* 🧺

Halo%s! Berikut adalah rincian terkini pengerjaan cucian Anda:

• *No. Order:* %s
• *Paket Layanan:* %s
• *Berat Cucian:* %s
• *Total Biaya:* %s
• *Status Terkini:* %s

_Estimasi pengerjaan: %d hari kerja._

Ketik *MENU* untuk pilihan info lainnya. Terima kasih! ✨`,
		greeting, order.OrderNumber, order.ServiceName, weightStr, priceStr, statusDesc, duration,
	)
}

func humanizeStatus(status model.OrderStatus) string {
	switch status {
	case model.OrderStatusWaitingPickup:
		return "🛵 *MENUNGGU PENJEMPUTAN*\n  _(Kurir segera menuju lokasi Anda)_"
	case model.OrderStatusNew:
		return "📋 *PESANAN DITERIMA*\n  _(Masuk antrean penimbangan & sortir)_"
	case model.OrderStatusConfirmed:
		return "✅ *DIKONFIRMASI / DALAM ANTREAN CUCI*"
	case model.OrderStatusInProgress:
		return "🧼 *SEDANG PROSES PENGERJAAN*\n  _(Cuci / Kering / Setrika Uap Rapi)_"
	case model.OrderStatusCompleted:
		return "✨ *SELESAI & WANGI RAPI*\n  _(Siap diambil di outlet atau menunggu pengantaran)_"
	case model.OrderStatusCancelled:
		return "❌ *DIBATALKAN*"
	default:
		return fmt.Sprintf("*%s*", string(status))
	}
}
