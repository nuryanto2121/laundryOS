# Sprint 08 - Loyalty, Membership & Customer Retention

> **Perubahan Roadmap**:  
> Lingkup asli Sprint 08 (POS MVP Kasir) telah diserap dan diselesaikan pada **Sprint 05 (Milestone 1 & 3)**.  
> Sprint 08 dialihkan menjadi pengembangan sistem **Loyalitas Pelanggan, Membership & Otomasi Retensi**.

---

## Goal

Meningkatkan frekuensi transaksi pelanggan berulang (*repeat order*) FAN Laundry melalui program poin digital, membership, voucher otomatis, dan pengingat jadwal cuci via WhatsApp.

## Business Problem

Pelanggan sering berpindah ke kompetitor atau lupa mencuci jika tidak ada pengingat rutin dan program penghargaan kesetiaan pelanggan.

## Scope

- **Program Loyalitas Digital (Cuci 9x Gratis 1x)**:
  - Pelacakan akumulasi transaksi pelanggan secara otomatis di database.
  - Saat pelanggan mencapai transaksi ke-10, diskon 1x cuci otomatis teraplikasikan di POS kasir dan bot WA memberikan ucapan selamat.
- **Voucher Promo & Diskon Ulang Tahun / Event**:
  - Generator kupon promo unik untuk disebarkan via WhatsApp blast broadcast.
- **Otomasi Pengingat Rutin (Smart Retention Reminder)**:
  - Cron scheduler mendeteksi pelanggan rutin yang belum mencuci kembali setelah 7-14 hari: *"Halo Kak, pakaian kotor sudah menumpuk? Kurir FAN Laundry siap jemput hari ini ya! 🛵✨"*.
- **Member Tiering (Silver, Gold, VIP)**:
  - Kategori pelanggan loyal dengan fasilitas prioritas antrian atau diskon langganan bulanan.

## Acceptance Criteria

- [ ] Akumulasi transaksi loyalitas terhitung akurat per nomor WhatsApp pelanggan.
- [ ] Pengingat otomatis terkirim tanpa dianggap spam oleh pelanggan.
- [ ] Kasir dapat dengan mudah memotong voucher/loyalitas saat input transaksi di Web POS.
