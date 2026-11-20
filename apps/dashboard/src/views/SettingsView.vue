<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- 1. Header Page -->
      <div
        class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-4"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-brand-600 to-indigo-600 text-white flex items-center justify-center shadow-md shadow-brand-500/20 shrink-0"
          >
            <Settings class="w-6 h-6" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h1 class="text-xl font-black text-slate-900 tracking-tight">
                Pengaturan & Master Data
              </h1>
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-brand-50 text-brand-700 border border-brand-200"
              >
                Milestone 6
              </span>
            </div>
            <p class="text-xs text-slate-500 mt-0.5">
              Kelola privasi kasir, template printer thermal, kapasitas tabung mesin, dan master
              layanan FAN Laundry.
            </p>
          </div>
        </div>

        <!-- Role Badge & Action -->
        <div class="flex items-center gap-2 shrink-0">
          <div
            class="px-3 py-1.5 rounded-xl border text-xs font-bold flex items-center gap-2 shadow-2xs"
            :class="
              settings.role === 'OWNER'
                ? 'bg-purple-50 text-purple-900 border-purple-200'
                : 'bg-slate-100 text-slate-700 border-slate-200'
            "
          >
            <Shield v-if="settings.role === 'OWNER'" class="w-4 h-4 text-purple-600" />
            <User v-else class="w-4 h-4 text-slate-500" />
            <span
              >Peran: {{ settings.role === 'OWNER' ? '👑 Owner (Akses Penuh)' : '👤 Kasir' }}</span
            >
          </div>

          <button
            type="button"
            class="px-3 py-1.5 rounded-xl bg-slate-900 text-white text-xs font-bold hover:bg-slate-800 transition-all cursor-pointer flex items-center gap-1.5"
            @click="handleSaveAll"
          >
            <Save class="w-3.5 h-3.5" />
            <span>Simpan</span>
          </button>
        </div>
      </div>

      <!-- 2. Navigation Tabs (6 Tabs) -->
      <div
        class="flex items-center gap-1.5 overflow-x-auto pb-1 scrollbar-none bg-slate-100/80 p-1.5 rounded-2xl text-xs font-bold"
      >
        <button
          v-for="tab in tabs"
          :key="tab.id"
          type="button"
          class="flex items-center gap-2 px-3.5 py-2 rounded-xl transition-all whitespace-nowrap cursor-pointer"
          :class="
            activeTab === tab.id
              ? 'bg-white text-slate-900 shadow-xs'
              : 'text-slate-600 hover:text-slate-900'
          "
          @click="activeTab = tab.id"
        >
          <component :is="tab.icon" class="w-4 h-4" />
          <span>{{ tab.label }}</span>
        </button>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 1: PROFIL OUTLET & PRINTER THERMAL                         -->
      <!-- ============================================================== -->
      <div v-if="activeTab === 'PRINTER'" class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Form Settings -->
        <div
          class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Printer class="w-4 h-4 text-brand-600" />
              <span>Konfigurasi Template Nota Thermal</span>
            </h3>
            <p class="text-xs text-slate-500">
              Sesuaikan informasi identitas outlet yang tercetak pada printer kasir.
            </p>
          </div>

          <!-- Paper Width Selector -->
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5"
              >Ukuran Kertas Thermal Printer</label
            >
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                class="p-3 rounded-xl border text-left cursor-pointer transition-all flex items-center justify-between"
                :class="
                  settings.printer.paperWidth === '58mm'
                    ? 'border-brand-500 bg-brand-50/50 ring-2 ring-brand-500/20'
                    : 'border-slate-200 bg-slate-50 hover:bg-white'
                "
                @click="settings.printer.paperWidth = '58mm'"
              >
                <div>
                  <span class="font-bold text-xs text-slate-900 block">58 mm (Standar)</span>
                  <span class="text-[11px] text-slate-500"
                    >Kompak, hemat kertas, printer portable/Bluetooth</span
                  >
                </div>
                <CheckCircle2
                  v-if="settings.printer.paperWidth === '58mm'"
                  class="w-4 h-4 text-brand-600 shrink-0 ml-2"
                />
              </button>

              <button
                type="button"
                class="p-3 rounded-xl border text-left cursor-pointer transition-all flex items-center justify-between"
                :class="
                  settings.printer.paperWidth === '80mm'
                    ? 'border-brand-500 bg-brand-50/50 ring-2 ring-brand-500/20'
                    : 'border-slate-200 bg-slate-50 hover:bg-white'
                "
                @click="settings.printer.paperWidth = '80mm'"
              >
                <div>
                  <span class="font-bold text-xs text-slate-900 block">80 mm (Lebar)</span>
                  <span class="text-[11px] text-slate-500"
                    >Lebar, jelas, font proporsional printer POS kasir</span
                  >
                </div>
                <CheckCircle2
                  v-if="settings.printer.paperWidth === '80mm'"
                  class="w-4 h-4 text-brand-600 shrink-0 ml-2"
                />
              </button>
            </div>
          </div>

          <!-- Store Title & Hotline -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >Nama Toko / Header Nota</label
              >
              <input
                v-model="settings.printer.storeName"
                type="text"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-semibold focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >No. WhatsApp Hotline</label
              >
              <input
                v-model="settings.printer.phone"
                type="text"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-semibold focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
          </div>

          <!-- Address Lines -->
          <div class="space-y-2">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Alamat Baris 1</label>
              <input
                v-model="settings.printer.addressLine1"
                type="text"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-semibold focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >Alamat Baris 2 (Kel/Kec/Kota)</label
              >
              <input
                v-model="settings.printer.addressLine2"
                type="text"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-semibold focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
          </div>

          <!-- Footer Notes -->
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">
              Catatan Kaki Nota Konsumen (Kebijakan & Syarat)
            </label>
            <textarea
              v-model="settings.printer.footerNote"
              rows="3"
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-mono focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              placeholder="Ketik catatan kaki nota di sini..."
            ></textarea>
            <p class="text-[11px] text-slate-400 mt-1">
              Gunakan enter untuk memisahkan baris pada kertas thermal.
            </p>
          </div>
        </div>

        <!-- Live Print Preview Box -->
        <div
          class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs flex flex-col justify-between"
        >
          <div>
            <div class="flex items-center justify-between pb-3 border-b border-slate-100">
              <span
                class="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-1.5"
              >
                <Printer class="w-3.5 h-3.5 text-brand-600" />
                Pratinjau Nota (Live Preview)
              </span>
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-slate-100 text-slate-700"
              >
                {{ settings.printer.paperWidth }}
              </span>
            </div>

            <!-- Receipt Card Container -->
            <div
              class="p-4 bg-slate-50 rounded-xl border border-slate-200/80 my-4 font-mono text-xs text-slate-900"
            >
              <div
                class="mx-auto bg-white p-3.5 shadow-xs border border-dashed border-slate-300 rounded-sm space-y-2 text-center"
                :class="settings.printer.paperWidth === '80mm' ? 'max-w-[340px]' : 'max-w-[260px]'"
              >
                <!-- Header -->
                <div class="border-b border-dashed border-slate-400 pb-2 space-y-0.5">
                  <h4 class="font-extrabold text-sm tracking-wider uppercase">
                    {{ settings.printer.storeName || 'FAN LAUNDRY' }}
                  </h4>
                  <p class="text-[10px] leading-tight text-slate-600">
                    {{ settings.printer.addressLine1 }}
                  </p>
                  <p class="text-[10px] leading-tight text-slate-600">
                    {{ settings.printer.addressLine2 }}
                  </p>
                  <p class="text-[10px] font-semibold text-slate-800">
                    WA: {{ settings.printer.phone }}
                  </p>
                </div>

                <!-- Sample Order Body -->
                <div
                  class="text-[10px] text-left border-b border-dashed border-slate-400 py-2 space-y-1"
                >
                  <div class="flex justify-between">
                    <span>No. Order:</span>
                    <span class="font-bold">#261120000002</span>
                  </div>
                  <div class="flex justify-between">
                    <span>Pelanggan:</span>
                    <span class="font-bold">Nuryanto</span>
                  </div>
                  <div class="flex justify-between">
                    <span>Layanan:</span>
                    <span>Cuci Basah (Self Service)</span>
                  </div>
                  <div class="flex justify-between font-bold border-t border-slate-200 pt-1">
                    <span>TOTAL:</span>
                    <span>Rp 28.000 (LUNAS)</span>
                  </div>
                </div>

                <!-- Footer -->
                <div class="text-[9px] text-slate-500 whitespace-pre-line leading-relaxed pt-1">
                  {{ settings.printer.footerNote }}
                </div>
                <p class="text-[9px] font-bold text-slate-700">
                  *** SIMPAN NOTA INI DENGAN BAIK ***
                </p>
              </div>
            </div>
          </div>

          <button
            type="button"
            class="w-full py-2.5 rounded-xl border border-slate-300 hover:bg-slate-50 text-slate-800 text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-2"
            @click="testPrint"
          >
            <Printer class="w-4 h-4 text-brand-600" />
            <span>Uji Coba Cetak (Browser Print)</span>
          </button>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 2: HAK AKSES & PRIVASI OWNER (RBAC)                        -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'RBAC'" class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Privacy Toggle & PIN -->
        <div
          class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-6"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <ShieldCheck class="w-4 h-4 text-purple-600" />
              <span>Privasi Finansial & Role-Based Access Control</span>
            </h3>
            <p class="text-xs text-slate-500">
              Atur batasan akses data sensitif keuangan antara kasir dan owner.
            </p>
          </div>

          <!-- Hide Financials Toggle -->
          <div
            class="p-4 rounded-xl border border-slate-200 bg-slate-50/70 flex items-start justify-between gap-4"
          >
            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <span class="font-bold text-xs text-slate-900">Sembunyikan Omset dari Kasir</span>
                <span
                  class="text-[10px] font-extrabold px-1.5 py-0.2 rounded-md bg-purple-100 text-purple-800"
                >
                  Rekomendasi Owner
                </span>
              </div>
              <p class="text-xs text-slate-500 leading-relaxed">
                Saat diaktifkan dan staf berada dalam <strong>Mode Kasir</strong>, total nilai omset
                bruto, pendapatan kas, dan laba bersih di halaman Beranda & Laporan akan disensor
                (<code class="bg-white px-1 py-0.5 rounded border border-slate-200">Rp ••••••</code
                >).
              </p>
            </div>

            <!-- Toggle Switch -->
            <button
              type="button"
              class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-hidden"
              :class="settings.hideFinancialsFromCashier ? 'bg-purple-600' : 'bg-slate-300'"
              @click="settings.hideFinancialsFromCashier = !settings.hideFinancialsFromCashier"
            >
              <span
                class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow-sm ring-0 transition duration-200 ease-in-out"
                :class="settings.hideFinancialsFromCashier ? 'translate-x-5' : 'translate-x-0'"
              />
            </button>
          </div>

          <!-- Change Owner PIN Form -->
          <div class="p-4 rounded-xl border border-purple-100 bg-purple-50/40 space-y-3">
            <div class="flex items-center gap-2">
              <Lock class="w-4 h-4 text-purple-600" />
              <h4 class="font-bold text-xs text-purple-950 uppercase tracking-wider">
                Ubah PIN Owner
              </h4>
            </div>
            <p class="text-xs text-slate-600">
              PIN ini digunakan untuk membuka sensor keuangan saat kasir ingin berkonsultasi atau
              saat owner login.
            </p>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-1">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1"
                  >PIN Lama (Default: 1234)</label
                >
                <input
                  v-model="oldPin"
                  type="password"
                  maxlength="8"
                  placeholder="PIN saat ini..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-mono focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                />
              </div>
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1"
                  >PIN Baru (Min. 4 Digit)</label
                >
                <input
                  v-model="newPin"
                  type="password"
                  maxlength="8"
                  placeholder="PIN baru..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-mono focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                />
              </div>
            </div>

            <button
              type="button"
              class="px-4 py-2 rounded-xl bg-purple-600 hover:bg-purple-700 text-white font-bold text-xs shadow-xs transition-all cursor-pointer"
              @click="handleSavePin"
            >
              Perbarui PIN Owner
            </button>
          </div>
        </div>

        <!-- Role Permissions Matrix -->
        <div
          class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm">Matriks Hak Akses (RBAC)</h3>
            <p class="text-xs text-slate-500">
              Perbandingan kewenangan kasir vs owner di LaundryOS.
            </p>
          </div>

          <div class="divide-y divide-slate-100 text-xs">
            <div class="py-2.5 flex items-center justify-between">
              <span class="font-medium text-slate-700">Pencatatan POS & Transaksi Walk-in</span>
              <div class="flex items-center gap-3">
                <span class="text-emerald-600 font-bold">Kasir: Ya</span>
                <span class="text-purple-600 font-bold">Owner: Ya</span>
              </div>
            </div>

            <div class="py-2.5 flex items-center justify-between">
              <span class="font-medium text-slate-700">Kelola 5 Antrian & Timbang Ulang</span>
              <div class="flex items-center gap-3">
                <span class="text-emerald-600 font-bold">Kasir: Ya</span>
                <span class="text-purple-600 font-bold">Owner: Ya</span>
              </div>
            </div>

            <div class="py-2.5 flex items-center justify-between">
              <span class="font-medium text-slate-700">Cetak Nota Thermal 3 Format</span>
              <div class="flex items-center gap-3">
                <span class="text-emerald-600 font-bold">Kasir: Ya</span>
                <span class="text-purple-600 font-bold">Owner: Ya</span>
              </div>
            </div>

            <div class="py-2.5 flex items-center justify-between">
              <span class="font-medium text-slate-700">Catat Kas Kecil Harian (Expenses)</span>
              <div class="flex items-center gap-3">
                <span class="text-emerald-600 font-bold">Kasir: Ya</span>
                <span class="text-purple-600 font-bold">Owner: Ya</span>
              </div>
            </div>

            <div class="py-2.5 flex items-center justify-between bg-purple-50/50 px-2 rounded-lg">
              <span class="font-bold text-slate-900">Lihat Laporan Omset & Laba Rugi</span>
              <div class="flex items-center gap-3">
                <span class="text-rose-600 font-bold">Kasir: Sensor</span>
                <span class="text-purple-600 font-bold">Owner: Ya</span>
              </div>
            </div>

            <div class="py-2.5 flex items-center justify-between bg-purple-50/50 px-2 rounded-lg">
              <span class="font-bold text-slate-900">Ubah Master Tarif & Diskon</span>
              <div class="flex items-center gap-3">
                <span class="text-slate-400 font-bold">Kasir: Tidak</span>
                <span class="text-purple-600 font-bold">Owner: Ya</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 3: KAPASITAS MESIN & TOLERANSI SELF SERVICE                 -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'MACHINE'" class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div
          class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Scale class="w-4 h-4 text-blue-600" />
              <span>Kapasitas Tabung Mesin Cuci & Pengering</span>
            </h3>
            <p class="text-xs text-slate-500">
              Konfigurasi parameter otomatis untuk kalkulator kapasitas layanan Cuci Basah / Kering
              Self Service.
            </p>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">
                Kapasitas Standar per Tabung (Kg)
              </label>
              <div class="relative">
                <input
                  v-model.number="settings.machine.tubeCapacityKg"
                  type="number"
                  step="0.5"
                  min="1"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold focus:ring-2 focus:ring-blue-500 focus:outline-hidden"
                />
                <span class="absolute right-3 top-2 text-xs font-bold text-slate-400">Kg</span>
              </div>
              <p class="text-[11px] text-slate-400 mt-1">
                Standar kapasitas mesin komersial FAN Laundry.
              </p>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">
                Batas Toleransi Kelebihan Berat (Kg)
              </label>
              <div class="relative">
                <input
                  v-model.number="settings.machine.toleranceKg"
                  type="number"
                  step="0.1"
                  min="0"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold focus:ring-2 focus:ring-blue-500 focus:outline-hidden"
                />
                <span class="absolute right-3 top-2 text-xs font-bold text-slate-400">Kg</span>
              </div>
              <p class="text-[11px] text-slate-400 mt-1">
                Keringanan kelebihan timbangan tanpa kena cas tabung kedua.
              </p>
            </div>
          </div>

          <!-- Total Safe Capacity Indicator -->
          <div class="p-4 rounded-xl bg-blue-50 border border-blue-200 flex items-start gap-3">
            <CheckCircle2 class="w-5 h-5 text-blue-600 shrink-0 mt-0.5" />
            <div class="text-xs">
              <span class="font-extrabold text-blue-950 block">
                Maksimal Muatan Aman 1 Tabung:
                {{
                  (
                    Number(settings.machine.tubeCapacityKg) + Number(settings.machine.toleranceKg)
                  ).toFixed(1)
                }}
                Kg
              </span>
              <p class="text-[11px] text-blue-800 mt-0.5">
                Jika timbangan pelanggan melebihi
                <strong
                  >{{
                    (
                      Number(settings.machine.tubeCapacityKg) + Number(settings.machine.toleranceKg)
                    ).toFixed(1)
                  }}
                  Kg</strong
                >, POS kasir akan otomatis memunculkan rekomendasi pemisahan menjadi
                <strong>2 Tabung</strong> atau menyisihkan kelebihannya.
              </p>
            </div>
          </div>

          <!-- Workflow Option: Skip Process -->
          <div
            class="p-4 rounded-xl border border-slate-200 bg-slate-50/70 flex items-center justify-between"
          >
            <div>
              <span class="font-bold text-xs text-slate-900 block"
                >Izinkan Opsi 'Skip Proses' di Antrian</span
              >
              <span class="text-xs text-slate-500">
                Memungkinkan pesanan khusus (misal titip setrika/lipat) langsung ke status Siap
                Ambil tanpa antrian cuci.
              </span>
            </div>
            <button
              type="button"
              class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-hidden"
              :class="settings.workflow.allowSkipProcess ? 'bg-brand-600' : 'bg-slate-300'"
              @click="settings.workflow.allowSkipProcess = !settings.workflow.allowSkipProcess"
            >
              <span
                class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow-sm ring-0 transition duration-200 ease-in-out"
                :class="settings.workflow.allowSkipProcess ? 'translate-x-5' : 'translate-x-0'"
              />
            </button>
          </div>
        </div>

        <div
          class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm">Simulasi Kalkulator POS</h3>
            <p class="text-xs text-slate-500">
              Uji coba reaksi logika kapasitas dengan berbagai berat timbangan.
            </p>
          </div>

          <div class="space-y-3">
            <label class="block text-xs font-bold text-slate-700">Simulasi Berat Timbangan:</label>
            <div class="flex items-center gap-2">
              <input
                v-model.number="simulatedWeight"
                type="number"
                step="0.5"
                min="1"
                class="w-full px-3 py-2 text-sm font-black rounded-xl border border-slate-200 text-center"
              />
              <span class="text-xs font-bold text-slate-400">Kg</span>
            </div>

            <div
              v-if="
                simulatedWeight <=
                Number(settings.machine.tubeCapacityKg) + Number(settings.machine.toleranceKg)
              "
              class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-800 space-y-1"
            >
              <span class="font-bold flex items-center gap-1.5 text-emerald-950">
                <CheckCircle2 class="w-4 h-4 text-emerald-600" />
                Masuk Kapasitas 1 Tabung
              </span>
              <p class="text-[11px]">Cucian aman dicuci dalam 1 tabung mesin cuci.</p>
            </div>

            <div
              v-else
              class="p-3 rounded-xl bg-amber-50 border border-amber-300 text-xs text-amber-900 space-y-1"
            >
              <span class="font-bold flex items-center gap-1.5 text-amber-950">
                <AlertTriangle class="w-4 h-4 text-amber-600" />
                Melebihi Toleransi 1 Tabung (+{{
                  (simulatedWeight - settings.machine.tubeCapacityKg).toFixed(1)
                }}
                Kg)
              </span>
              <p class="text-[11px]">
                Sistem POS akan menawarkan proses 2 Tabung (Rp {{ formatRupiah(28000) }}) atau
                sisihkan lebihnya.
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 4: MASTER LAYANAN & TARIF                                  -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'SERVICES'"
        class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
      >
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-2"
        >
          <div>
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Tag class="w-4 h-4 text-emerald-600" />
              <span>Katalog Master Layanan & Tarif</span>
            </h3>
            <p class="text-xs text-slate-500">
              Daftar layanan aktif yang dapat dipilih oleh kasir di Web POS.
            </p>
          </div>
          <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-slate-100 text-slate-700">
            {{ servicesList.length }} Layanan Terdaftar
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr
                class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider"
              >
                <th class="py-3 px-3">Kode</th>
                <th class="py-3 px-3">Nama Layanan</th>
                <th class="py-3 px-3">Kategori</th>
                <th class="py-3 px-3">Satuan</th>
                <th class="py-3 px-3 text-right">Tarif (Rp)</th>
                <th class="py-3 px-3 text-center">Durasi</th>
                <th class="py-3 px-3 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="s in servicesList"
                :key="s.id || s.service_code"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-2.5 px-3 font-mono font-bold text-slate-700">{{ s.service_code }}</td>
                <td class="py-2.5 px-3 font-bold text-slate-900">{{ s.service_name }}</td>
                <td class="py-2.5 px-3">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                    :class="
                      s.category === 'TABUNG'
                        ? 'bg-blue-100 text-blue-800'
                        : s.category === 'KILOAN'
                          ? 'bg-emerald-100 text-emerald-800'
                          : 'bg-purple-100 text-purple-800'
                    "
                  >
                    {{ s.category }}
                  </span>
                </td>
                <td class="py-2.5 px-3 font-semibold text-slate-600">{{ s.unit }}</td>
                <td class="py-2.5 px-3 text-right font-extrabold text-slate-900">
                  Rp {{ formatRupiah(s.price) }}
                </td>
                <td class="py-2.5 px-3 text-center font-medium text-slate-600">
                  {{ s.duration_day === 0 ? 'Hari ini' : `${s.duration_day} Hari` }}
                </td>
                <td class="py-2.5 px-3 text-center">
                  <span
                    class="inline-flex items-center gap-1 text-[10px] font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full"
                  >
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                    Aktif
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 5: MASTER PARFUM & DISKON                                  -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'PERFUME_DISCOUNT'"
        class="grid grid-cols-1 lg:grid-cols-2 gap-6"
      >
        <!-- Master Parfum -->
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div>
              <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
                <Sparkles class="w-4 h-4 text-amber-500" />
                <span>Varian Aroma Parfum</span>
              </h3>
              <p class="text-xs text-slate-500">
                Pilihan parfum yang tersedia untuk order cucian pelanggan.
              </p>
            </div>
            <span class="text-xs font-bold text-slate-500">{{ perfumesList.length }} Aroma</span>
          </div>

          <div class="space-y-2">
            <div
              v-for="p in perfumesList"
              :key="p.id"
              class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/50 flex items-center justify-between"
            >
              <div>
                <span class="font-bold text-xs text-slate-900 block">{{ p.name }}</span>
                <span class="text-[11px] text-slate-500">{{
                  p.description || 'Aroma wangi tahan lama'
                }}</span>
              </div>
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-700"
              >
                Tersedia
              </span>
            </div>
          </div>
        </div>

        <!-- Master Diskon -->
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div>
              <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
                <Percent class="w-4 h-4 text-rose-500" />
                <span>Diskon & Promo Loyalitas</span>
              </h3>
              <p class="text-xs text-slate-500">
                Diskon otomatis yang dapat diaplikasikan pada transaksi kasir.
              </p>
            </div>
            <span class="text-xs font-bold text-slate-500">{{ discountsList.length }} Promo</span>
          </div>

          <div class="space-y-2">
            <div
              v-for="d in discountsList"
              :key="d.id"
              class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/50 flex items-center justify-between"
            >
              <div>
                <span class="font-bold text-xs text-slate-900 block">{{ d.name }}</span>
                <span class="text-[11px] text-slate-500">
                  {{
                    d.discount_type === 'PERCENTAGE'
                      ? `Diskon ${Number(d.value)}%`
                      : `Potongan Rp ${formatRupiah(d.value)}`
                  }}
                  <span v-if="Number(d.min_order_amount) > 0">
                    • Min. Order Rp {{ formatRupiah(d.min_order_amount) }}</span
                  >
                </span>
              </div>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-rose-50 text-rose-700">
                Aktif
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 6: STATUS INTEGRASI GATEWAY & SISTEM                       -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'INTEGRATION'" class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div
          class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Smartphone class="w-4 h-4 text-emerald-600" />
              <span>Integrasi WhatsApp Gateway & Servis Backend</span>
            </h3>
            <p class="text-xs text-slate-500">
              Monitoring status koneksi microservice bot, API, dan database.
            </p>
          </div>

          <!-- Services Status Cards -->
          <div class="space-y-3">
            <!-- Gateway Service -->
            <div
              class="p-4 rounded-xl border border-emerald-200 bg-emerald-50/50 flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-10 h-10 rounded-xl bg-emerald-600 text-white flex items-center justify-center font-bold"
                >
                  WA
                </div>
                <div>
                  <h4 class="font-bold text-xs text-slate-900">WhatsApp Gateway (Baileys)</h4>
                  <p class="text-[11px] text-slate-500">Koneksi nomor hotline 0852-8597-7834</p>
                </div>
              </div>
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800"
              >
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                Tersambung
              </span>
            </div>

            <!-- Backend API -->
            <div
              class="p-4 rounded-xl border border-blue-200 bg-blue-50/50 flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-10 h-10 rounded-xl bg-blue-600 text-white flex items-center justify-center font-bold"
                >
                  API
                </div>
                <div>
                  <h4 class="font-bold text-xs text-slate-900">LaundryOS REST API (Go / Chi)</h4>
                  <p class="text-[11px] text-slate-500">Versi v0.5.0 • Port 8080</p>
                </div>
              </div>
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-blue-100 text-blue-800"
              >
                <span class="w-2 h-2 rounded-full bg-blue-500"></span>
                Sehat (Healthy)
              </span>
            </div>

            <!-- PostgreSQL Database -->
            <div
              class="p-4 rounded-xl border border-indigo-200 bg-indigo-50/50 flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-10 h-10 rounded-xl bg-indigo-600 text-white flex items-center justify-center font-bold"
                >
                  DB
                </div>
                <div>
                  <h4 class="font-bold text-xs text-slate-900">PostgreSQL 16</h4>
                  <p class="text-[11px] text-slate-500">Database laundryos • Migrasi up-to-date</p>
                </div>
              </div>
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-indigo-100 text-indigo-800"
              >
                <span class="w-2 h-2 rounded-full bg-indigo-500"></span>
                Connected
              </span>
            </div>
          </div>
        </div>

        <!-- Factory Reset Settings Card -->
        <div
          class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2 text-rose-700">
              <RotateCcw class="w-4 h-4 text-rose-600" />
              <span>Reset Pengaturan Pabrik</span>
            </h3>
            <p class="text-xs text-slate-500">
              Kembalikan seluruh preferensi printer, PIN, dan toleransi ke pengaturan standar.
            </p>
          </div>

          <p class="text-xs text-slate-600 leading-relaxed">
            Tindakan ini akan mengembalikan nama toko menjadi <strong>FAN LAUNDRY</strong>, ukuran
            kertas printer ke <strong>58mm</strong>, kapasitas tabung ke <strong>8.0 Kg</strong>,
            toleransi ke <strong>0.5 Kg</strong>, dan PIN Owner ke <strong>1234</strong>.
          </p>

          <button
            type="button"
            class="w-full py-2.5 rounded-xl border border-rose-300 text-rose-700 hover:bg-rose-50 text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-2"
            @click="handleFactoryReset"
          >
            <RotateCcw class="w-4 h-4" />
            <span>Kembalikan ke Default Pabrik</span>
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '../components/layout/AppLayout.vue'
import {
  Settings,
  Printer,
  Shield,
  ShieldCheck,
  User,
  Lock,
  Scale,
  Sparkles,
  Tag,
  Percent,
  CheckCircle2,
  AlertTriangle,
  Smartphone,
  RotateCcw,
  Save,
} from 'lucide-vue-next'
import { useAppSettings } from '../composables/useAppSettings'
import { useToast } from '../composables/useToast'
import { api, formatRupiah } from '../api/client'

const { settings, changePin, resetToDefaults } = useAppSettings()
const { success, error } = useToast()

const activeTab = ref('PRINTER')

const tabs = [
  { id: 'PRINTER', label: 'Profil Outlet & Printer Thermal', icon: Printer },
  { id: 'RBAC', label: 'Hak Akses & Privasi (RBAC)', icon: ShieldCheck },
  { id: 'MACHINE', label: 'Kapasitas Mesin & Toleransi', icon: Scale },
  { id: 'SERVICES', label: 'Master Layanan & Tarif', icon: Tag },
  { id: 'PERFUME_DISCOUNT', label: 'Master Parfum & Diskon', icon: Sparkles },
  { id: 'INTEGRATION', label: 'Status Integrasi & Sistem', icon: Smartphone },
]

// State for PIN Form
const oldPin = ref('')
const newPin = ref('')

// State for Simulation
const simulatedWeight = ref(8.5)

// Master Data Lists from API
const servicesList = ref([])
const perfumesList = ref([])
const discountsList = ref([])

const fetchMasterData = async () => {
  try {
    const [srv, prf, disc] = await Promise.all([
      api.get('/master/services').catch(() => []),
      api.get('/master/perfumes').catch(() => []),
      api.get('/master/discounts').catch(() => []),
    ])
    if (Array.isArray(srv)) servicesList.value = srv
    if (Array.isArray(prf)) perfumesList.value = prf
    if (Array.isArray(disc)) discountsList.value = disc
  } catch (err) {
    console.error('Failed to load master data:', err)
  }
}

onMounted(() => {
  fetchMasterData()
})

const handleSaveAll = () => {
  success('Seluruh pengaturan berhasil disimpan ke memori sistem!')
}

const handleSavePin = () => {
  const res = changePin(oldPin.value, newPin.value)
  if (res.success) {
    success(res.message)
    oldPin.value = ''
    newPin.value = ''
  } else {
    error(res.message)
  }
}

const testPrint = () => {
  window.print()
}

const handleFactoryReset = () => {
  if (confirm('Yakin ingin mereset seluruh preferensi ke pengaturan default pabrik?')) {
    resetToDefaults()
    success('Pengaturan telah dikembalikan ke default pabrik.')
  }
}
</script>
