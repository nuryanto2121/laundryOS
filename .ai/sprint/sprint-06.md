# Sprint 06 - AI Customer Assistant & Guardrails

> **Perubahan Roadmap**:  
> Lingkup asli Sprint 06 (Workflow Operasional Cucian) telah diserap dan diselesaikan pada **Sprint 05 (Milestone 1 & 3)**.  
> Sprint 06 dialihkan menjadi pengembangan **Asisten AI Percakapan & Sistem Guardrail Pelanggan**.

---

## Goal

Mengintegrasikan AI Cerdas (Google Gemini) pada WhatsApp Gateway dan Web Chat untuk melayani percakapan alami pelanggan dengan batasan proteksi (*guardrails*) ketat agar fokus hanya pada konteks bisnis laundry.

## Business Problem

1. FAQ otomatis berbasis pencocokan kata kunci (*keyword matcher*) kaku dan tidak mampu merespons pertanyaan pelanggan yang menggunakan bahasa santai atau kalimat panjang.
2. Tanpa batasan ketat (*guardrails*), AI percakapan publik berisiko disalahgunakan pelanggan untuk bertanya hal-hal di luar bisnis (politik, resep masakan, koding, tugas sekolah), yang menghabiskan token dan mengurangi profesionalitas brand FAN Laundry.

## Scope

- **AI Engine Integration**: Menghubungkan WhatsApp Baileys Gateway & Web Chat dengan Google Gemini Flash API.
- **Laundry Domain Knowledge Base**: Konteks lengkap FAN Laundry (daftar harga, estimasi waktu pengerjaan 3 jam vs harian, aturan tabung max 8.5 kg, wangi parfum, alamat Kalibaru, rekening bayar).
- **Strict Context Boundary & Guardrails**:
  - Filter intent dan prompt boundary: AI hanya menjawab pertanyaan seputar layanan, status, perawatan pakaian, harga, dan operasional FAN Laundry.
  - Respon sopan penolakan konteks luar: *"Mohon maaf Kak, saya hanya dapat membantu informasi dan layanan seputar FAN Laundry. Ada yang bisa kami bantu terkait cucian Kakak? 😊🧺"*.
- **Eskalasi ke Admin/Human Handover**: Jika pelanggan komplain keras atau meminta bicara langsung dengan pemilik, bot memberi tahu owner dan mengalihkan sesi.

## Acceptance Criteria

- [ ] AI mampu menjawab pertanyaan non-baku pelanggan seputar cucian dengan ramah dan akurat.
- [ ] Pertanyaan di luar konteks laundry ditolak secara sopan dan konsisten.
- [ ] Sesi percakapan tetap menjaga efisiensi token dan keamanan data pelanggan.
- [ ] Terintegrasi dengan database untuk pengecekan status pesanan secara alami.
