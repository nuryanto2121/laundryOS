package bot

import (
	"fmt"
	"strings"
	"time"

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
	// If dynamic services are provided, format using the canonical 5 groups
	return formatPriceResponse()
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

// FormatIDR formats a numeric amount with thousands dots (e.g. 25000 -> "25.000").
func FormatIDR(amount float64) string {
	return formatIDR(amount)
}

func formatPriceResponse() string {
	return `🧺 *DAFTAR HARGA & LAYANAN FAN LAUNDRY* 🧺

📍 Jl. Kalibaru Barat I kel. Kalibaru kec. Cilincing, Jakarta Utara 14110
📞 Hotline: +62 852-8597-7834

Berikut adalah daftar harga resmi layanan laundry kami:

*1. CUCI SETRIKA (Min. 4 Kg)*
• Regular (3 Hari): *Rp 7.000 / Kg*
• Express (2 Hari): *Rp 9.000 / Kg*
• Kilat (1 Hari / 24 Jam): *Rp 11.000 / Kg*
• Super Kilat (12 Jam): *Rp 13.000 / Kg*
• Turbo (6 Jam): *Rp 16.000 / Kg*

*2. SETRIKA SAJA (Min. 4 Kg)*
• Regular (3 Hari): *Rp 6.000 / Kg*
• Express (2 Hari): *Rp 8.000 / Kg*
• Kilat (1 Hari / 24 Jam): *Rp 10.000 / Kg*
• Super Kilat (12 Jam): *Rp 12.000 / Kg*
• Turbo (6 Jam): *Rp 15.000 / Kg*

*3. KARPET PERMADANI (Pengerjaan 7-8 Hari)*
• Permadani Sedang: *Rp 20.000 / Mtr*
• Permadani Rumbai: *Rp 22.000 / Mtr*
• Permadani Tebal: *Rp 30.000 / Mtr*
• Kasur Lantai Anti Slip / Busa: *Rp 40.000 - Rp 60.000*

*4. BAD COVER*
• Regular (3 Hari): *Rp 30.000 / Pcs*
• Express (2 Hari): *Rp 35.000 / Pcs*
• Kilat (1 Hari): *Rp 40.000 / Pcs*

*5. BONEKA (Pengerjaan 3 Hari)*
• Mini: *Rp 15.000 / Pcs*
• Kecil: *Rp 25.000 / Pcs*
• Sedang: *Rp 35.000 / Pcs*
• Besar: *Rp 60.000 / Pcs*

✨ *Mau pakaian kotor dijemput sekarang?*
Cukup balas paket pilihan Anda (contoh: *'Cuci Reguler'*, *'Kilat'*, atau *'Bad Cover'*), atau ketik *2* untuk request penjemputan! 🛵💨`
}

func formatServiceResponse() string {
	return `🧺 *PILIHAN LAYANAN FAN LAUNDRY* 🧺

Kami melayani 5 kategori utama cucian dengan deterjen higienis dan pewangi premium wangi tahan lama:

1️⃣ *Cuci Setrika (Min. 4 Kg)*: Regular (3 hari), Express (2 hari), Kilat (24 jam), Super Kilat (12 jam), Turbo (6 jam)
2️⃣ *Setrika Saja (Min. 4 Kg)*: Regular, Express, Kilat, Super Kilat, Turbo
3️⃣ *Karpet Permadani*: Sedang, Rumbai, Tebal, dan Kasur Busa/Lantai (7-8 hari)
4️⃣ *Bad Cover*: Regular, Express, Kilat
5️⃣ *Boneka*: Mini, Kecil, Sedang, Besar

✨ 1 Mesin 1 Pelanggan (Tidak pernah dicampur)

Cukup balas nama paket (contoh: *'Cuci Reguler'* atau *'Bad Cover'*) untuk memesan jemput cucian Anda! 🛵`
}

func formatHoursResponse() string {
	return `⏰ *JAM OPERASIONAL FAN LAUNDRY* ⏰

Kami buka setiap hari untuk melayani kebutuhan cucian Anda:

📅 *Hari:* Senin - Minggu (Buka Setiap Hari)
🕐 *Outlet Fisik:* 07.00 - 21.00 WIB
🚚 *Layanan Antar-Jemput:* 08.00 - 19.00 WIB
💬 *Customer Service WhatsApp:* 07.00 - 22.00 WIB

_Pemesanan pickup di luar jam operasional akan dijadwalkan pada pagi hari berikutnya._

Ketik *2* jika ingin request jemput cucian! 🛵`
}

func formatAddressResponse() string {
	return `📍 *LOKASI OUTLET & ALAMAT FAN LAUNDRY* 📍

🏠 *Alamat Outlet:*
Jl. Kalibaru Barat I kel. Kalibaru kec. Cilincing, Jakarta Utara 14110
(Patokan: Kalibaru Barat I)

🗺️ *Google Maps:*
https://maps.google.com/?q=FAN+Laundry+Cilincing

📞 *Telepon / WhatsApp:*
+62 852-8597-7834

⏰ *Jam Buka:*
• Outlet Fisik: *07.00 - 21.00 WIB* (Buka Setiap Hari)
• Antar-Jemput Kurir: *08.00 - 19.00 WIB*

🚗 Tersedia area parkir untuk motor & mobil.

Ketik *2* atau balas: *'Mau jemput cucian'* jika ingin pakaian kotor dijemput kurir kami! 🛵`
}

func formatCoverageResponse() string {
	return `🎁 *PROMO, JANGKAUAN AREA & KETENTUAN* 🎁

🚚 *Layanan Antar-Jemput FAN Laundry:*
• *GRATIS ONGKIR* untuk area Kalibaru, Cilincing dan sekitarnya (min. 5 Kg)
• Di bawah 5 Kg atau jarak lebih jauh: Ongkir flat hemat Rp 5.000 - Rp 10.000

✨ *Jaminan Standar Mutu FAN Laundry:*
• 1 Mesin 1 Pelanggan (Tidak pernah dicampur dengan cucian orang lain)
• Menggunakan deterjen higienis dan pewangi premium wangi tahan lama
• Minimal 4 Kg untuk Cuci Setrika dan Setrika Saja
• Setrika uap rapi tidak merusak serat kain

🛵 *Mau request jemput cucian sekarang?*
Balas langsung paket pilihan Anda (contoh: *'Cuci Reguler'*) atau ketik *2* untuk mulai penjemputan!`
}

func formatPickupResponse() string {
	return `🛵 *REQUEST PENJEMPUTAN FAN LAUNDRY* 🛵

Kurir kami siap menjemput pakaian kotor langsung ke rumah/kos/apartemen Anda.

Boleh share *Alamat Lengkap Penjemputan* dan perkiraan jam berapa mau dijemput?
Contoh: _Jl. Mawar No. 12, jemput jam 2 siang ya_

_(Ketik BATAL jika ingin membatalkan)_`
}

func formatHelpResponse(nameGreeting string) string {
	return fmt.Sprintf(`🧺 *Halo%s! Selamat Datang di FAN Laundry* 👋

Silakan ketik angka berikut untuk info cepat:

1️⃣ *Info Layanan & Harga* ➔ Tarif kiloan, satuan & bed cover
2️⃣ *Order / Penjemputan Laundry* ➔ Minta kurir jemput pakaian kotor
3️⃣ *Alamat & Jam Operasional* ➔ Lokasi outlet & jam buka
4️⃣ *Promo & Syarat* ➔ Gratis ongkir & area jangkauan

Atau Anda juga bisa langsung chat santai seperti:
• *"Mba saya mau laundry"*
• *"Cuci regular ya"*
• *"Neng ambil baju kotor ke rumah ya"* 🙏`, nameGreeting)
}

func formatUnknownResponse(nameGreeting string) string {
	return fmt.Sprintf(`🧺 *Halo%s dari FAN Laundry!* 👋

Maaf, kami belum memahami pertanyaan Anda. Untuk respon otomatis yang cepat, silakan ketik salah satu pilihan angka berikut:

1️⃣ *Info Layanan & Harga*
2️⃣ *Order / Penjemputan Laundry*
3️⃣ *Alamat & Jam Buka*
4️⃣ *Promo & Syarat*

Atau Anda bisa langsung ketik: *"Mba mau laundry"* untuk kami jemput pakaian kotornya. Terima kasih! 🙏`, nameGreeting)
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
		return fmt.Sprintf(`🧺 *CEK STATUS CUCIAN FAN LAUNDRY* 🧺

Halo%s! Saat ini tidak ditemukan transaksi cucian aktif atas nomor WhatsApp Anda.

• Jika baru saja menitipkan pakaian, staf kami sedang memproses penimbangan & antrean.
• Anda juga dapat mengecek dengan mengetik nomor nota secara manual (contoh: *STATUS #261118000001*).
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

	receivedAt := time.Now()
	if order.CreatedAt != "" {
		for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05"} {
			if parsed, err := time.Parse(layout, order.CreatedAt); err == nil {
				receivedAt = parsed
				break
			}
		}
	}

	slaEstimate := CalculateCompletionEstimate(order.ServiceName, duration, receivedAt)

	return fmt.Sprintf("🧺 *STATUS CUCIAN ANDA (FAN Laundry)* 🧺\n\n"+
		"Halo%s! Berikut adalah rincian terkini pengerjaan cucian Anda:\n\n"+
		"• *No. Order:* %s\n"+
		"• *Paket Layanan:* %s\n"+
		"• *Berat Cucian:* %s\n"+
		"• *Total Biaya:* %s\n"+
		"• *Status Terkini:* %s\n"+
		"• *Estimasi Selesai:* *%s*\n\n"+
		"Ketik *MENU* untuk pilihan info lainnya. Terima kasih! ✨",
		greeting, FormatDisplayOrderNumber(order.OrderNumber), order.ServiceName, weightStr, priceStr, statusDesc, slaEstimate,
	)
}

// FormatDisplayOrderNumber formats an internal order number like ORD261119000012 to #261119000012
// which is single-tap copy-friendly in WhatsApp.
func FormatDisplayOrderNumber(orderNumber string) string {
	clean := strings.TrimSpace(orderNumber)
	if strings.HasPrefix(strings.ToUpper(clean), "ORD") {
		return "#" + clean[3:]
	}
	if strings.HasPrefix(clean, "#") {
		return clean
	}
	if clean == "" {
		return ""
	}
	return "#" + clean
}

func getIndonesianDay(t time.Time) string {
	days := [...]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	return days[t.Weekday()]
}

// CalculateCompletionEstimate calculates exact estimated completion time with Indonesian day and time.
func CalculateCompletionEstimate(serviceName string, durationDays int, receivedAt time.Time) string {
	lower := strings.ToLower(serviceName)
	wib := time.FixedZone("WIB", 7*3600)
	t := receivedAt.In(wib)

	// 1. Self Service (Cuci Basah or Cuci Kering) -> 3 Jam SLA
	if strings.Contains(lower, "self service") || strings.Contains(lower, "cuci basah") || strings.Contains(lower, "cuci kering") {
		target := t.Add(3 * time.Hour)
		cutoff := time.Date(t.Year(), t.Month(), t.Day(), 21, 0, 0, 0, wib)
		if target.Before(cutoff) || target.Equal(cutoff) {
			return fmt.Sprintf("Hari ini, %s WIB (Estimasi 3 jam)", target.Format("02-01-2006 15:04"))
		}
		nextDay := t.AddDate(0, 0, 1)
		return fmt.Sprintf("Besok pagi, %s pukul 09:00 - 10:00 WIB", nextDay.Format("02-01-2006"))
	}

	// 2. Express Turbo (6 Jam)
	if strings.Contains(lower, "turbo") {
		target := t.Add(6 * time.Hour)
		return fmt.Sprintf("%s, %s WIB (Estimasi 6 jam)", getIndonesianDay(target), target.Format("02-01-2006 15:04"))
	}

	// 3. Super Kilat (12 Jam)
	if strings.Contains(lower, "super kilat") {
		target := t.Add(12 * time.Hour)
		return fmt.Sprintf("%s, %s WIB (Estimasi 12 jam)", getIndonesianDay(target), target.Format("02-01-2006 15:04"))
	}

	// 4. Daily SLA (Cuci Setrika, Setrika, Bad Cover, Karpet, Boneka)
	days := durationDays
	if days <= 0 {
		days = 3
	}
	target := t.AddDate(0, 0, days)
	return fmt.Sprintf("%s, %s WIB (Estimasi %d hari)", getIndonesianDay(target), target.Format("02-01-2006 15:04"), days)
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
