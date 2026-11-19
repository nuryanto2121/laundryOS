# Sprint 09 - Multi-Branch, Inventory & Logistics Scale

> **Perubahan Roadmap**:  
> Lingkup asli Sprint 09 (Dashboard Analitik Dasar & Export Laporan) telah diserap dan diselesaikan pada **Sprint 05 (Milestone 5)** dan dikembangkan lebih lanjut menjadi AI Analytics pada **Sprint 07**.  
> Sprint 09 dialihkan menjadi **Manajemen Inventori Bahan Baku & Skalabilitas Multi-Cabang / Armada Kurir**.

---

## Goal

Menyediakan skalabilitas operasional saat FAN Laundry membuka cabang/outlet baru di wilayah lain, mencakup kontrol stok deterjen/parfum dan manajemen armada kurir terpusat.

## Business Problem

Ketika volume transaksi membesar dan cabang bertambah, pencatatan stok deterjen/gas/plastik dan rute kurir membutuhkan pengawasan terpusat agar tidak terjadi kebocoran bahan baku.

## Scope

- **Manajemen Inventori & Stok Bahan Baku**:
  - Pelacakan stok deterjen, softener, parfum, gas dryer, dan plastik packing.
  - Notifikasi peringatan otomatis saat stok bahan baku mendekati batas minimum.
- **Multi-Branch Centralized Management**:
  - Dashboard terpusat untuk memantau performa beberapa outlet (Outlet Kalibaru, Outlet 2, dst).
  - Transfer stok antar outlet.
- **Advanced Logistics & Courier Route Dispatch**:
  - Penugasan order pickup/delivery ke beberapa kurir dengan rute teroptimal.

## Acceptance Criteria

- [ ] Penggunaan bahan baku terhitung otomatis berdasarkan volume cucian harian.
- [ ] Multi-outlet switcher dan pelaporan terpisah/konsolidasi per outlet berjalan sempurna.
