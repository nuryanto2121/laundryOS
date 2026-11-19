# Sprint 05 - Web Dashboard & POS (Point of Sale)

> **Catatan Konsolidasi Arsitektur**:  
> Sprint 05 telah diekspansi secara strategis untuk menggabungkan dan menyelesaikan kebutuhan inti dari:
> - **Sprint 06 (Workflow & Status Siklus Cucian)**
> - **Sprint 08 (POS Kasir Walk-in & Cetak Nota)**
> - **Sprint 09 (Laporan Keuangan Dasar & CRM Operasional)**
> Dengan konsolidasi ini, FAN Laundry langsung memiliki sistem operasional mandiri yang 100% menggantikan aplikasi pihak ketiga (**Edobi**).

---

## Goal

Membangun Web Dashboard & POS modern, responsif, dan terintegrasi 2-arah dengan WhatsApp Bot untuk kasir, staf operasional, dan owner FAN Laundry.

## Scope & 6 Milestones

1. **Milestone 1 (Selesai ✅)**: Skema Database & Backend REST APIs
   - Migration `000003_pos_and_expenses.up.sql`
   - Multi-layanan per order (`order_items`), kas kecil (`expenses`), foto dokumentasi (`order_photos`), audit status, master parfum, diskon, metode bayar & outlet
   - REST API `/pos/...`, `/expenses/...`, `/reports/...`, `/master/...` dengan CORS
   - Notifikasi otomatis WhatsApp saat order dibuat & status berganti
   - 112 unit tests Go passing

2. **Milestone 2 (Sedang Berjalan ⏳)**: Setup Frontend, Fresh Modern Layout & Halaman Home
   - Setup Tailwind CSS & Lucide icons di `apps/dashboard`
   - Fresh Modern UI Layout (Mobile-first bottom navigation + Desktop header)
   - Real-time KPI Cards (Omset Hari Ini, Masuk, Harus Selesai, Terlambat)
   - Floating Quick Actions (`+ Transaksi`, `+ Pengeluaran`, `Cari Order`) & Outlet Switcher

3. **Milestone 3**: Walk-in POS Kasir & Order Lifecycle
   - Kasir walk-in multi-layanan (timbang desimal, parfum, diskon, opsi bayar)
   - 5 Tab Status Antrian (`Antrian`, `Proses`, `Siap Ambil`, `Selesai`, `Batal`)
   - Tombol progresi 1-klik (`PROSES`, `SIAP AMBIL`, `SELESAI`)
   - Pelunasan pembayaran (Tunai, Transfer BCA, QRIS AstraPay)
   - Cetak 3 nota thermal (Nota Konsumen, Pertinggal, Produksi) & Upload foto dokumentasi

4. **Milestone 4**: Manajemen Pengeluaran Kas Kecil (Petty Cash)
   - Catat biaya operasional harian (deterjen, plastik, listrik, air, sewa, gaji)
   - Rekap pengeluaran untuk perhitungan laba bersih

5. **Milestone 5**: Laporan Keuangan, Operasional & CRM
   - Grafik garis Omset & Pendapatan Kas harian + Export to Excel
   - Laporan Piutang & Laporan Laba/Rugi
   - Rekap volume layanan (Kg, Pcs, Tabung, Meter) & Rak belum diambil
   - CRM Pelanggan Loyal (Top 10 leaderboard & riwayat 50 transaksi seumur hidup)

6. **Milestone 6**: Master Settings & RBAC Pegawai
   - Privasi Owner: Saklar toggle untuk sembunyikan omset/keuangan dari kasir
   - Kustomisasi printer thermal & template nota (logo, header, footer)
   - Master layanan, parfum, diskon loyalitas, dan toggle *Skip Proses*
