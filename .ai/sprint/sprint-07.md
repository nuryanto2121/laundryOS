# Sprint 07 - AI Conversational Analytics & Business Intelligence

## Goal

Memberikan kemampuan kecerdasan buatan (*AI-Powered Analytics*) pada Web/Mobile Dashboard agar Owner dan Manajemen dapat mengajukan pertanyaan analitik bisnis laundry dalam bahasa sehari-hari (*Natural Language Query / Text-to-Insight*) dan AI langsung menyajikan data riil dari database.

## Business Problem

Owner dan staf manajemen sering kali kesulitan menemukan insight cepat saat sedang mobile jika harus membuka banyak menu tabel dan laporan manual. Owner ingin bisa bertanya langsung layaknya memiliki asisten manajer bisnis pribadi.

## Scope

- **AI Text-to-Data Assistant**:
  - Chat interface khusus Owner/Admin di Dashboard Web & Mobile.
  - Pemrosesan bahasa alami untuk mengubah pertanyaan owner menjadi query analitik database (Read-Only) yang aman.
- **Kasus Pertanyaan Analitik yang Didukung**:
  - *"Berapa omset dan laba bersih 7 hari terakhir dibanding minggu lalu?"*
  - *"Siapa saja 5 pelanggan dengan piutang paling lama dan berapa totalnya?"*
  - *"Layanan apa yang paling banyak menyumbang omset di bulan ini?"*
  - *"Berapa pengeluaran terbesar kita bulan ini dan untuk kategori apa?"*
  - *"Hari apa saja yang orderannya paling ramai di FAN Laundry?"*
- **Visualisasi Dinamis**:
  - AI tidak hanya menjawab teks angka, melainkan otomatis menampilkan kartu ringkasan, tabel breakdown, dan visual grafik batang/garis.
- **Rekomendasi Bisnis Cerdas**:
  - AI memberikan saran berbasis data (misal: mengingatkan untuk follow-up pelanggan loyal yang sudah 2 minggu tidak mencuci, atau evaluasi pembelian deterjen).

## Security & Guardrails

- **Strict Access Isolation**: Fitur ini terisolasi hanya untuk user berhak akses **Owner** (tidak bisa diakses pelanggan umum atau kasir biasa).
- **Read-Only Data Access**: AI hanya diberikan akses eksekusi query analitik `SELECT` terproteksi, tidak dapat mengubah atau menghapus data transaksi.

## Acceptance Criteria

- [ ] Owner dapat bertanya seputar omset, piutang, pengeluaran, dan pelanggan dalam bahasa Indonesia santai.
- [ ] AI menyajikan angka akurat yang bersumber langsung dari database PostgreSQL.
- [ ] Tersedia visualisasi grafik atau tabel otomatis dari hasil jawaban AI.
- [ ] Waktu respon cepat (< 3 detik) menggunakan Gemini Function Calling / Tool Calling.
