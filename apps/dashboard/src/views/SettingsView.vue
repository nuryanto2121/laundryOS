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
              Kelola privasi kasir, template printer thermal, master layanan, pegawai & hak akses,
              pengeluaran kas kecil, dan pelanggan.
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
            class="px-3.5 py-1.5 rounded-xl bg-slate-900 text-white text-xs font-bold hover:bg-slate-800 transition-all cursor-pointer flex items-center gap-1.5 shadow-xs"
            @click="handleSaveAll"
          >
            <Save class="w-3.5 h-3.5" />
            <span>Simpan</span>
          </button>
        </div>
      </div>

      <!-- 2. Mobile Quick-Picker Dropdown (< md) -->
      <div class="block md:hidden bg-white p-3.5 rounded-2xl border border-slate-200/90 shadow-xs">
        <label
          class="block text-[10px] font-extrabold text-slate-400 uppercase tracking-wider mb-1.5"
        >
          Pilih Menu Pengaturan & Master Data
        </label>
        <div class="relative">
          <select
            v-model="activeTab"
            class="w-full pl-10 pr-10 py-3 bg-slate-50 border border-slate-200 rounded-xl font-bold text-xs text-slate-800 appearance-none focus:outline-hidden focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all cursor-pointer"
          >
            <option v-for="tab in tabs" :key="tab.id" :value="tab.id">
              {{ tab.label }}
            </option>
          </select>
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-brand-600"
          >
            <component :is="activeTabObj.icon" class="w-4 h-4" />
          </div>
          <div
            class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-slate-400"
          >
            <ChevronDown class="w-4 h-4" />
          </div>
        </div>
      </div>

      <!-- 3. Desktop Navigation Tabs (>= md) -->
      <div
        class="hidden md:flex items-center gap-1.5 overflow-x-auto pb-1 scrollbar-none bg-slate-100/80 p-1.5 rounded-2xl text-xs font-bold"
      >
        <button
          v-for="tab in tabs"
          :key="tab.id"
          type="button"
          class="flex items-center gap-2 px-3 py-2 rounded-xl transition-all whitespace-nowrap cursor-pointer shrink-0"
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
      <!-- TAB 2: PEGAWAI & HAK AKSES (RBAC)                              -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'EMPLOYEES'" class="space-y-6">
        <!-- Privacy Toggle & PIN Section -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <div
            class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5"
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
                  Saat diaktifkan dan staf berada dalam <strong>Mode Kasir</strong>, total nilai
                  omset bruto, pendapatan kas, dan laba bersih di halaman Beranda & Laporan akan
                  disensor (<code class="bg-white px-1 py-0.5 rounded border border-slate-200"
                    >Rp ••••••</code
                  >).
                </p>
              </div>

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
          </div>

          <!-- Change Owner PIN Form -->
          <div
            class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
          >
            <div class="border-b border-slate-100 pb-3">
              <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
                <Lock class="w-4 h-4 text-purple-600" />
                <span>Ubah PIN Akses Owner</span>
              </h3>
              <p class="text-xs text-slate-500">PIN cepat 4-digit untuk membuka sensor keuangan.</p>
            </div>

            <div class="space-y-3">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1"
                  >PIN Lama (Default: 1234)</label
                >
                <input
                  v-model="oldPin"
                  type="password"
                  maxlength="8"
                  placeholder="PIN saat ini..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono focus:bg-white focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                />
              </div>
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1"
                  >PIN Baru (4-8 Digit)</label
                >
                <input
                  v-model="newPin"
                  type="password"
                  maxlength="8"
                  placeholder="Ketik PIN baru..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono focus:bg-white focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                />
              </div>

              <button
                type="button"
                :disabled="!oldPin || !newPin"
                class="w-full py-2.5 rounded-xl bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white font-bold text-xs transition-all shadow-xs cursor-pointer"
                @click="handleSavePin"
              >
                Perbarui PIN Keamanan
              </button>
            </div>
          </div>
        </div>

        <!-- Master Pegawai CRUD Table -->
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
          <div
            class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3"
          >
            <div>
              <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
                <Users class="w-4 h-4 text-purple-600" />
                <span>Daftar Pegawai & Hak Akses</span>
              </h3>
              <p class="text-xs text-slate-500">
                Kelola akun staf kasir, kurir, operator, dan hak akses otorisasi sistem.
              </p>
            </div>
            <button
              type="button"
              class="px-3.5 py-2 rounded-xl bg-purple-600 hover:bg-purple-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1.5 cursor-pointer self-start sm:self-auto"
              @click="openAddEmployeeModal"
            >
              <Plus class="w-4 h-4" />
              <span>Tambah Pegawai Baru</span>
            </button>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr
                  class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider"
                >
                  <th class="py-3 px-3">Username</th>
                  <th class="py-3 px-3">Nama Lengkap</th>
                  <th class="py-3 px-3">No. HP</th>
                  <th class="py-3 px-3 text-center">Hak Akses (Role)</th>
                  <th class="py-3 px-3 text-center">Status</th>
                  <th class="py-3 px-3 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-if="employeesList.length === 0">
                  <td colspan="6" class="py-6 text-center text-slate-400">
                    Belum ada pegawai terdaftar.
                  </td>
                </tr>
                <tr
                  v-for="emp in employeesList"
                  :key="emp.id"
                  class="hover:bg-slate-50/70 transition-colors"
                >
                  <td class="py-2.5 px-3 font-mono font-bold text-slate-800">{{ emp.username }}</td>
                  <td class="py-2.5 px-3 font-bold text-slate-900">{{ emp.full_name }}</td>
                  <td class="py-2.5 px-3 font-medium text-slate-600">{{ emp.phone || '-' }}</td>
                  <td class="py-2.5 px-3 text-center">
                    <span
                      class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                      :class="getRoleBadgeClass(emp.role)"
                    >
                      {{ emp.role }}
                    </span>
                  </td>
                  <td class="py-2.5 px-3 text-center">
                    <span
                      class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full"
                      :class="
                        emp.is_active
                          ? 'bg-emerald-50 text-emerald-700'
                          : 'bg-slate-100 text-slate-500'
                      "
                    >
                      <span
                        class="w-1.5 h-1.5 rounded-full"
                        :class="emp.is_active ? 'bg-emerald-500' : 'bg-slate-400'"
                      ></span>
                      {{ emp.is_active ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </td>
                  <td class="py-2.5 px-3 text-center">
                    <div class="flex items-center justify-center gap-1.5">
                      <button
                        type="button"
                        class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-slate-100 transition-all cursor-pointer"
                        title="Edit Pegawai"
                        @click="openEditEmployeeModal(emp)"
                      >
                        <Edit2 class="w-3.5 h-3.5" />
                      </button>
                      <button
                        type="button"
                        class="p-1.5 rounded-lg text-slate-500 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                        title="Hapus Pegawai"
                        @click="confirmDeleteEmployee(emp)"
                      >
                        <Trash2 class="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 3: MASTER LAYANAN & TARIF                                  -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'SERVICES'"
        class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
      >
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3"
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
          <button
            type="button"
            class="px-3.5 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1.5 cursor-pointer self-start sm:self-auto"
            @click="openAddServiceModal"
          >
            <Plus class="w-4 h-4" />
            <span>Tambah Layanan Baru</span>
          </button>
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
                <th class="py-3 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="servicesList.length === 0">
                <td colspan="8" class="py-6 text-center text-slate-400">
                  Belum ada layanan terdaftar.
                </td>
              </tr>
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
                    class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full"
                    :class="
                      s.is_active !== false
                        ? 'bg-emerald-50 text-emerald-700'
                        : 'bg-slate-100 text-slate-500'
                    "
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full"
                      :class="s.is_active !== false ? 'bg-emerald-500' : 'bg-slate-400'"
                    ></span>
                    {{ s.is_active !== false ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </td>
                <td class="py-2.5 px-3 text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-slate-100 transition-all cursor-pointer"
                      title="Edit Layanan"
                      @click="openEditServiceModal(s)"
                    >
                      <Edit2 class="w-3.5 h-3.5" />
                    </button>
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-500 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                      title="Hapus Layanan"
                      @click="confirmDeleteService(s)"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 4: MASTER PARFUM & DISKON                                  -->
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
              <p class="text-xs text-slate-500">Pilihan parfum yang tersedia untuk order cucian.</p>
            </div>
            <button
              type="button"
              class="px-3 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1 cursor-pointer"
              @click="openAddPerfumeModal"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Aroma Baru</span>
            </button>
          </div>

          <div class="space-y-2">
            <div v-if="perfumesList.length === 0" class="py-6 text-center text-xs text-slate-400">
              Belum ada parfum terdaftar.
            </div>
            <div
              v-for="p in perfumesList"
              :key="p.id"
              class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/50 flex items-center justify-between hover:bg-slate-50 transition-colors"
            >
              <div>
                <span class="font-bold text-xs text-slate-900 block">{{ p.name }}</span>
                <span class="text-[11px] text-slate-500">{{
                  p.description || 'Aroma wangi tahan lama'
                }}</span>
              </div>
              <div class="flex items-center gap-1.5">
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="
                    p.is_active !== false
                      ? 'bg-emerald-50 text-emerald-700'
                      : 'bg-slate-100 text-slate-500'
                  "
                >
                  {{ p.is_active !== false ? 'Aktif' : 'Nonaktif' }}
                </span>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-brand-600 hover:bg-slate-200/60"
                  @click="openEditPerfumeModal(p)"
                >
                  <Edit2 class="w-3.5 h-3.5" />
                </button>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                  @click="confirmDeletePerfume(p)"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </div>
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
            <button
              type="button"
              class="px-3 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1 cursor-pointer"
              @click="openAddDiscountModal"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Promo Baru</span>
            </button>
          </div>

          <div class="space-y-2">
            <div v-if="discountsList.length === 0" class="py-6 text-center text-xs text-slate-400">
              Belum ada diskon terdaftar.
            </div>
            <div
              v-for="d in discountsList"
              :key="d.id"
              class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/50 flex items-center justify-between hover:bg-slate-50 transition-colors"
            >
              <div>
                <div class="flex items-center gap-2">
                  <span class="font-bold text-xs text-slate-900">{{ d.name }}</span>
                  <span
                    class="font-mono text-[10px] font-extrabold px-1.5 py-0.2 rounded bg-slate-200 text-slate-700"
                  >
                    {{ d.code }}
                  </span>
                </div>
                <span class="text-[11px] text-slate-500">
                  {{
                    d.discount_type === 'PERCENTAGE' || d.discount_percent > 0
                      ? `Diskon ${Number(d.discount_percent || d.value)}%`
                      : `Potongan Rp ${formatRupiah(d.fixed_amount || d.value)}`
                  }}
                  <span v-if="Number(d.min_order_amount) > 0">
                    • Min. Order Rp {{ formatRupiah(d.min_order_amount) }}
                  </span>
                </span>
              </div>
              <div class="flex items-center gap-1.5">
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="
                    d.is_active !== false
                      ? 'bg-rose-50 text-rose-700'
                      : 'bg-slate-100 text-slate-500'
                  "
                >
                  {{ d.is_active !== false ? 'Aktif' : 'Nonaktif' }}
                </span>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-brand-600 hover:bg-slate-200/60"
                  @click="openEditDiscountModal(d)"
                >
                  <Edit2 class="w-3.5 h-3.5" />
                </button>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                  @click="confirmDeleteDiscount(d)"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 5: KATEGORI PENGELUARAN                                    -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'EXPENSE_CATEGORIES'"
        class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
      >
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3"
        >
          <div>
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Receipt class="w-4 h-4 text-amber-600" />
              <span>Master Kategori Pengeluaran Kas Kecil</span>
            </h3>
            <p class="text-xs text-slate-500">
              Daftar pos anggaran pengeluaran untuk klasifikasi petty cash dan laporan laba rugi.
            </p>
          </div>
          <button
            type="button"
            class="px-3.5 py-2 rounded-xl bg-amber-600 hover:bg-amber-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1.5 cursor-pointer self-start sm:self-auto"
            @click="openAddCategoryModal"
          >
            <Plus class="w-4 h-4" />
            <span>Tambah Kategori Baru</span>
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr
                class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider"
              >
                <th class="py-3 px-3">Kode Kategori</th>
                <th class="py-3 px-3">Nama Kategori Pengeluaran</th>
                <th class="py-3 px-3 text-center">Status</th>
                <th class="py-3 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="expenseCategoriesList.length === 0">
                <td colspan="4" class="py-6 text-center text-slate-400">
                  Belum ada kategori pengeluaran.
                </td>
              </tr>
              <tr
                v-for="cat in expenseCategoriesList"
                :key="cat.id"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-2.5 px-3 font-mono font-bold text-slate-700">
                  {{ cat.category_code }}
                </td>
                <td class="py-2.5 px-3 font-bold text-slate-900">{{ cat.category_name }}</td>
                <td class="py-2.5 px-3 text-center">
                  <span
                    class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full"
                    :class="
                      cat.is_active
                        ? 'bg-emerald-50 text-emerald-700'
                        : 'bg-slate-100 text-slate-500'
                    "
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full"
                      :class="cat.is_active ? 'bg-emerald-500' : 'bg-slate-400'"
                    ></span>
                    {{ cat.is_active ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </td>
                <td class="py-2.5 px-3 text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-slate-100 transition-all cursor-pointer"
                      title="Edit Kategori"
                      @click="openEditCategoryModal(cat)"
                    >
                      <Edit2 class="w-3.5 h-3.5" />
                    </button>
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-500 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                      title="Hapus Kategori"
                      @click="confirmDeleteCategory(cat)"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 6: PENGOREKSI PENGELUARAN (AUDIT KAS KECIL)                 -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'EXPENSE_CORRECTION'"
        class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5"
      >
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3"
        >
          <div>
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <DollarSign class="w-4 h-4 text-rose-600" />
              <span>Pengoreksi & Audit Pengeluaran Kas Kecil</span>
            </h3>
            <p class="text-xs text-slate-500">
              Koreksi kesalahan input tanggal, pos biaya, atau nominal pengeluaran kas kecil kasir.
            </p>
          </div>
        </div>

        <!-- Filter Bar -->
        <div
          class="p-4 bg-slate-50 rounded-xl border border-slate-200/80 grid grid-cols-1 sm:grid-cols-3 gap-3"
        >
          <div>
            <label class="block text-[11px] font-bold text-slate-600 mb-1">Mulai Tanggal</label>
            <input
              v-model="expenseAuditFilter.startDate"
              type="date"
              class="w-full px-3 py-2 bg-white border border-slate-200 rounded-xl text-xs font-semibold text-slate-800"
            />
          </div>
          <div>
            <label class="block text-[11px] font-bold text-slate-600 mb-1">Sampai Tanggal</label>
            <input
              v-model="expenseAuditFilter.endDate"
              type="date"
              class="w-full px-3 py-2 bg-white border border-slate-200 rounded-xl text-xs font-semibold text-slate-800"
            />
          </div>
          <div>
            <label class="block text-[11px] font-bold text-slate-600 mb-1">Filter Kategori</label>
            <select
              v-model="expenseAuditFilter.category"
              class="w-full px-3 py-2 bg-white border border-slate-200 rounded-xl text-xs font-semibold text-slate-800"
            >
              <option value="">Semua Kategori</option>
              <option v-for="cat in expenseCategoriesList" :key="cat.id" :value="cat.category_code">
                {{ cat.category_name }}
              </option>
            </select>
          </div>
        </div>

        <!-- Expense Audit Table -->
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr
                class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider"
              >
                <th class="py-3 px-3">Tanggal</th>
                <th class="py-3 px-3">Kategori</th>
                <th class="py-3 px-3 text-right">Nominal (Rp)</th>
                <th class="py-3 px-3">Keterangan Biaya</th>
                <th class="py-3 px-3">Dicatat Oleh</th>
                <th class="py-3 px-3 text-center">Aksi Koreksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="expensesAuditList.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400">
                  Tidak ada catatan pengeluaran pada rentang tanggal ini.
                </td>
              </tr>
              <tr
                v-for="exp in expensesAuditList"
                :key="exp.id"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-2.5 px-3 font-semibold text-slate-800 whitespace-nowrap">
                  {{ formatDate(exp.expense_date) }}
                </td>
                <td class="py-2.5 px-3">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-50 text-amber-800 border border-amber-200"
                  >
                    {{ exp.category }}
                  </span>
                </td>
                <td class="py-2.5 px-3 text-right font-extrabold text-rose-600 whitespace-nowrap">
                  Rp {{ formatRupiah(exp.amount) }}
                </td>
                <td class="py-2.5 px-3 font-medium text-slate-700 max-w-xs truncate">
                  {{ exp.description || exp.notes || '-' }}
                </td>
                <td class="py-2.5 px-3 text-slate-500 font-medium">
                  {{ exp.created_by || 'Kasir' }}
                </td>
                <td class="py-2.5 px-3 text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <button
                      type="button"
                      class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-brand-50 hover:text-brand-600 text-slate-700 font-bold text-xs transition-all flex items-center gap-1 cursor-pointer"
                      title="Koreksi Pengeluaran"
                      @click="openEditExpenseModal(exp)"
                    >
                      <Edit2 class="w-3 h-3" />
                      <span>Koreksi</span>
                    </button>
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                      title="Hapus Pengeluaran"
                      @click="confirmDeleteExpense(exp)"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Bar -->
        <PaginationBar
          v-if="totalExpensesAudit > 0"
          v-model:current-page="pageExpensesAudit"
          v-model:items-per-page="limitExpensesAudit"
          :total-items="totalExpensesAudit"
        />
      </div>

      <!-- ============================================================== -->
      <!-- TAB 7: KELOLA PELANGGAN (CRM MASTER)                           -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'CUSTOMERS'"
        class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
      >
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3"
        >
          <div>
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Users class="w-4 h-4 text-brand-600" />
              <span>Kelola Master Pelanggan (CRM)</span>
            </h3>
            <p class="text-xs text-slate-500">
              Daftar seluruh database pelanggan outlet dengan riwayat kontak & alamat.
            </p>
          </div>
          <button
            type="button"
            class="px-3.5 py-2 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1.5 cursor-pointer self-start sm:self-auto"
            @click="openAddCustomerModal"
          >
            <Plus class="w-4 h-4" />
            <span>Tambah Pelanggan Baru</span>
          </button>
        </div>

        <!-- Search Bar -->
        <div class="relative">
          <input
            v-model="customerSearchQuery"
            type="text"
            placeholder="Cari nama pelanggan, nomor WhatsApp, atau alamat..."
            class="w-full pl-9 pr-8 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs font-medium text-slate-800 focus:bg-white focus:outline-hidden focus:ring-2 focus:ring-brand-500 transition-all"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-3 pointer-events-none" />
          <button
            v-if="customerSearchQuery"
            class="absolute right-2.5 top-2.5 text-slate-400 hover:text-slate-600"
            @click="customerSearchQuery = ''"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr
                class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider"
              >
                <th class="py-3 px-3">Nama Pelanggan</th>
                <th class="py-3 px-3">No. WhatsApp</th>
                <th class="py-3 px-3">Alamat</th>
                <th class="py-3 px-3">Terdaftar</th>
                <th class="py-3 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="customersList.length === 0">
                <td colspan="5" class="py-8 text-center text-slate-400">
                  Tidak ada data pelanggan yang cocok.
                </td>
              </tr>
              <tr
                v-for="c in customersList"
                :key="c.id"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-2.5 px-3 font-bold text-slate-900">{{ c.full_name }}</td>
                <td class="py-2.5 px-3 font-mono font-medium text-slate-700">{{ c.phone }}</td>
                <td class="py-2.5 px-3 text-slate-600 max-w-xs truncate">{{ c.address || '-' }}</td>
                <td class="py-2.5 px-3 text-slate-500 text-[11px] whitespace-nowrap">
                  {{ formatDate(c.created_at) }}
                </td>
                <td class="py-2.5 px-3 text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-slate-100 transition-all cursor-pointer"
                      title="Edit Pelanggan"
                      @click="openEditCustomerModal(c)"
                    >
                      <Edit2 class="w-3.5 h-3.5" />
                    </button>
                    <button
                      type="button"
                      class="p-1.5 rounded-lg text-slate-500 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                      title="Hapus Pelanggan"
                      @click="confirmDeleteCustomer(c)"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Bar -->
        <PaginationBar
          v-if="totalCustomers > 0"
          v-model:current-page="pageCustomers"
          v-model:items-per-page="limitCustomers"
          :total-items="totalCustomers"
        />
      </div>

      <!-- ============================================================== -->
      <!-- TAB 8: KAPASITAS MESIN & TOLERANSI                             -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'MACHINE'" class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div
          class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Scale class="w-4 h-4 text-blue-600" />
              <span>Parameter Kapasitas Tabung Mesin Cuci</span>
            </h3>
            <p class="text-xs text-slate-500">
              Konfigurasi batasan muatan tabung dan ambang batas toleransi lebihan timbangan.
            </p>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >Kapasitas Standar 1 Tabung (Kg)</label
              >
              <input
                v-model.number="settings.machine.tubeCapacityKg"
                type="number"
                step="0.5"
                min="1"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold focus:ring-2 focus:ring-blue-500 focus:outline-hidden"
              />
              <p class="text-[11px] text-slate-400 mt-1">Default mesin cuci FAN Laundry: 8.0 Kg</p>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >Batas Toleransi Berat (Kg)</label
              >
              <input
                v-model.number="settings.machine.toleranceKg"
                type="number"
                step="0.1"
                min="0"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold focus:ring-2 focus:ring-blue-500 focus:outline-hidden"
              />
              <p class="text-[11px] text-slate-400 mt-1">
                Lebihan gratis sebelum disarankan 2 tabung (contoh: 0.5 Kg)
              </p>
            </div>
          </div>
        </div>

        <!-- Simulation Box -->
        <div
          class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4"
        >
          <div class="border-b border-slate-100 pb-3">
            <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
              <Scale class="w-4 h-4 text-blue-600" />
              <span>Simulasi Logika Timbangan Kasir</span>
            </h3>
            <p class="text-xs text-slate-500">
              Tes hasil rekomendasi kasir saat pelanggan membawa cucian kiloan.
            </p>
          </div>

          <div class="space-y-3">
            <label class="block text-xs font-bold text-slate-700"
              >Berat Cucian Uji Coba:
              <strong class="text-blue-600">{{ simulatedWeight }} Kg</strong></label
            >
            <div class="flex items-center gap-3">
              <input
                v-model.number="simulatedWeight"
                type="range"
                min="1"
                max="20"
                step="0.1"
                class="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-blue-600"
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
      <!-- TAB 9: STATUS INTEGRASI GATEWAY & SISTEM                       -->
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

          <div class="space-y-3">
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

      <!-- ============================================================== -->
      <!-- MODALS SECTION                                                 -->
      <!-- ============================================================== -->

      <!-- 1. Modal Tambah/Edit Pegawai -->
      <div
        v-if="showEmployeeModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Users class="w-4 h-4 text-purple-600" />
              <span>{{ isEditingEmployee ? 'Edit Data Pegawai' : 'Tambah Pegawai Baru' }}</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showEmployeeModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitEmployeeForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Username Login *</label>
              <input
                v-model="employeeForm.username"
                type="text"
                required
                :disabled="isEditingEmployee"
                placeholder="Contoh: kasir2"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Lengkap *</label>
              <input
                v-model="employeeForm.full_name"
                type="text"
                required
                placeholder="Contoh: Siti Aminah"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-purple-500"
              />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1"
                  >Peran / Hak Akses *</label
                >
                <select
                  v-model="employeeForm.role"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500"
                >
                  <option value="OWNER">👑 Owner</option>
                  <option value="KASIR">👤 Kasir</option>
                  <option value="KURIR">🛵 Kurir</option>
                  <option value="OPERATOR">⚙️ Operator</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">
                  {{
                    isEditingEmployee
                      ? 'PIN Baru (Kosongkan jika tdk ubah)'
                      : 'PIN Akses (4 Digit) *'
                  }}
                </label>
                <input
                  v-model="employeeForm.pin"
                  type="password"
                  maxlength="6"
                  :required="!isEditingEmployee"
                  placeholder="1234"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono focus:bg-white focus:ring-2 focus:ring-purple-500"
                />
              </div>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >No. Handphone / WhatsApp</label
              >
              <input
                v-model="employeeForm.phone"
                type="text"
                placeholder="08123456789"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-purple-500"
              />
            </div>
            <div class="flex items-center gap-2 pt-1">
              <input
                id="emp-active"
                v-model="employeeForm.is_active"
                type="checkbox"
                class="rounded text-purple-600 focus:ring-purple-500 w-4 h-4 cursor-pointer"
              />
              <label for="emp-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Status Akun Aktif (Dapat login)
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showEmployeeModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-purple-600 hover:bg-purple-700 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Pegawai' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 2. Modal Tambah/Edit Layanan -->
      <div
        v-if="showServiceModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Tag class="w-4 h-4 text-emerald-600" />
              <span>{{ isEditingService ? 'Edit Layanan' : 'Tambah Layanan Baru' }}</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showServiceModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitServiceForm">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Kode Layanan *</label>
                <input
                  v-model="serviceForm.service_code"
                  type="text"
                  required
                  placeholder="CK-01"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono font-bold uppercase focus:bg-white focus:ring-2 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Kategori *</label>
                <select
                  v-model="serviceForm.category"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-emerald-500"
                >
                  <option value="KILOAN">KILOAN</option>
                  <option value="TABUNG">TABUNG</option>
                  <option value="SATUAN">SATUAN</option>
                  <option value="METERAN">METERAN</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Layanan *</label>
              <input
                v-model="serviceForm.service_name"
                type="text"
                required
                placeholder="Contoh: Cuci Kering Seterika"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-emerald-500"
              />
            </div>
            <div class="grid grid-cols-3 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Satuan *</label>
                <input
                  v-model="serviceForm.unit"
                  type="text"
                  required
                  placeholder="Kg / Tabung"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Tarif (Rp) *</label>
                <input
                  v-model.number="serviceForm.price"
                  type="number"
                  min="0"
                  step="500"
                  required
                  placeholder="10000"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Durasi (Hari)</label>
                <input
                  v-model.number="serviceForm.duration_day"
                  type="number"
                  min="0"
                  placeholder="1"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>
            <div class="flex items-center gap-2 pt-1">
              <input
                id="srv-active"
                v-model="serviceForm.is_active"
                type="checkbox"
                class="rounded text-emerald-600 focus:ring-emerald-500 w-4 h-4 cursor-pointer"
              />
              <label for="srv-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Layanan Aktif (Tampil di POS)
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showServiceModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Layanan' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 3. Modal Tambah/Edit Parfum -->
      <div
        v-if="showPerfumeModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Sparkles class="w-4 h-4 text-amber-500" />
              <span>{{ isEditingPerfume ? 'Edit Aroma Parfum' : 'Tambah Aroma Baru' }}</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showPerfumeModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitPerfumeForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Aroma *</label>
              <input
                v-model="perfumeForm.name"
                type="text"
                required
                placeholder="Contoh: Sakura Blossom"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Deskripsi Aroma</label>
              <input
                v-model="perfumeForm.description"
                type="text"
                placeholder="Aroma segar & lembut tahan hingga 7 hari"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-medium focus:bg-white focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div class="flex items-center gap-2 pt-1">
              <input
                id="prf-active"
                v-model="perfumeForm.is_active"
                type="checkbox"
                class="rounded text-amber-600 focus:ring-amber-500 w-4 h-4 cursor-pointer"
              />
              <label for="prf-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Tersedia untuk Cucian
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showPerfumeModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Parfum' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 4. Modal Tambah/Edit Diskon -->
      <div
        v-if="showDiscountModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Percent class="w-4 h-4 text-rose-500" />
              <span>{{ isEditingDiscount ? 'Edit Diskon / Promo' : 'Tambah Promo Baru' }}</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showDiscountModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitDiscountForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Kode Promo *</label>
              <input
                v-model="discountForm.code"
                type="text"
                required
                placeholder="PROMO10"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono font-bold uppercase focus:bg-white focus:ring-2 focus:ring-rose-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Diskon *</label>
              <input
                v-model="discountForm.name"
                type="text"
                required
                placeholder="Diskon Pelanggan Baru 10%"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-rose-500"
              />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Diskon Persen (%)</label>
                <input
                  v-model.number="discountForm.discount_percent"
                  type="number"
                  min="0"
                  max="100"
                  placeholder="10"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-rose-500"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1"
                  >Potongan Tetap (Rp)</label
                >
                <input
                  v-model.number="discountForm.fixed_amount"
                  type="number"
                  min="0"
                  step="500"
                  placeholder="5000"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-rose-500"
                />
              </div>
            </div>
            <div class="flex items-center gap-2 pt-1">
              <input
                id="disc-active"
                v-model="discountForm.is_active"
                type="checkbox"
                class="rounded text-rose-600 focus:ring-rose-500 w-4 h-4 cursor-pointer"
              />
              <label for="disc-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Promo Aktif
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showDiscountModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Promo' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 5. Modal Tambah/Edit Kategori Pengeluaran -->
      <div
        v-if="showCategoryModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Receipt class="w-4 h-4 text-amber-600" />
              <span>{{
                isEditingCategory ? 'Edit Kategori Pengeluaran' : 'Tambah Kategori Baru'
              }}</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showCategoryModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitCategoryForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Kode Kategori *</label>
              <input
                v-model="categoryForm.category_code"
                type="text"
                required
                :disabled="isEditingCategory"
                placeholder="Contoh: PLASTIK"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono font-bold uppercase focus:bg-white focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >Nama Kategori Pengeluaran *</label
              >
              <input
                v-model="categoryForm.category_name"
                type="text"
                required
                placeholder="Plastik & Roll Packaging"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div class="flex items-center gap-2 pt-1">
              <input
                id="cat-active"
                v-model="categoryForm.is_active"
                type="checkbox"
                class="rounded text-amber-600 focus:ring-amber-500 w-4 h-4 cursor-pointer"
              />
              <label for="cat-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Kategori Aktif (Bisa Dipilih Kasir)
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showCategoryModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-amber-600 hover:bg-amber-700 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Kategori' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 6. Modal Koreksi Pengeluaran (Edit Expense) -->
      <div
        v-if="showEditExpenseModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <DollarSign class="w-4 h-4 text-rose-600" />
              <span>Koreksi Data Pengeluaran Kas Kecil</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showEditExpenseModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitEditExpenseForm">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Tanggal *</label>
                <input
                  v-model="editExpenseForm.expense_date"
                  type="date"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-rose-500"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Kategori *</label>
                <select
                  v-model="editExpenseForm.category"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-rose-500"
                >
                  <option
                    v-for="cat in expenseCategoriesList"
                    :key="cat.id"
                    :value="cat.category_code"
                  >
                    {{ cat.category_name }}
                  </option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nominal (Rp) *</label>
              <input
                v-model.number="editExpenseForm.amount"
                type="number"
                min="500"
                step="500"
                required
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-black text-rose-600 focus:bg-white focus:ring-2 focus:ring-rose-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Keterangan Biaya *</label>
              <textarea
                v-model="editExpenseForm.description"
                rows="2"
                required
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-medium focus:bg-white focus:ring-2 focus:ring-rose-500"
              ></textarea>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Dicatat Oleh</label>
              <input
                v-model="editExpenseForm.created_by"
                type="text"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-medium focus:bg-white focus:ring-2 focus:ring-rose-500"
              />
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showEditExpenseModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Perbarui Pengeluaran' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 7. Modal Tambah/Edit Pelanggan -->
      <div
        v-if="showCustomerModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Users class="w-4 h-4 text-brand-600" />
              <span>{{ isEditingCustomer ? 'Edit Data Pelanggan' : 'Tambah Pelanggan Baru' }}</span>
            </h4>
            <button
              class="p-1 rounded-lg text-slate-400 hover:text-slate-600"
              @click="showCustomerModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitCustomerForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Lengkap *</label>
              <input
                v-model="customerForm.full_name"
                type="text"
                required
                placeholder="Contoh: Budi Santoso"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-brand-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >Nomor Handphone / WhatsApp *</label
              >
              <input
                v-model="customerForm.phone"
                type="text"
                required
                placeholder="Contoh: 081234567890"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono font-medium focus:bg-white focus:ring-2 focus:ring-brand-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Alamat Domisili</label>
              <textarea
                v-model="customerForm.address"
                rows="2"
                placeholder="Contoh: Jl. Kalibaru Timur No. 12, RT 02/05"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-medium focus:bg-white focus:ring-2 focus:ring-brand-500"
              ></textarea>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100"
                @click="showCustomerModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all shadow-xs"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Pelanggan' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import AppLayout from '../components/layout/AppLayout.vue'
import PaginationBar from '../components/ui/PaginationBar.vue'
import {
  Settings,
  Printer,
  Shield,
  ShieldCheck,
  User,
  Users,
  Lock,
  Scale,
  Sparkles,
  Tag,
  Percent,
  Receipt,
  DollarSign,
  CheckCircle2,
  AlertTriangle,
  Smartphone,
  RotateCcw,
  Save,
  ChevronDown,
  Plus,
  Edit2,
  Trash2,
  X,
  Search,
} from 'lucide-vue-next'
import { useAppSettings } from '../composables/useAppSettings'
import { useToast } from '../composables/useToast'
import { api, formatRupiah, formatDate } from '../api/client'

const { settings, changePin, resetToDefaults } = useAppSettings()
const { success, error } = useToast()

const activeTab = ref('PRINTER')

const tabs = [
  { id: 'PRINTER', label: 'Profil Outlet & Printer Thermal', icon: Printer },
  { id: 'EMPLOYEES', label: 'Pegawai & Hak Akses (RBAC)', icon: ShieldCheck },
  { id: 'SERVICES', label: 'Master Layanan & Tarif', icon: Tag },
  { id: 'PERFUME_DISCOUNT', label: 'Master Parfum & Diskon', icon: Sparkles },
  { id: 'EXPENSE_CATEGORIES', label: 'Kategori Pengeluaran', icon: Receipt },
  { id: 'EXPENSE_CORRECTION', label: 'Pengoreksi Kas Kecil', icon: DollarSign },
  { id: 'CUSTOMERS', label: 'Kelola Pelanggan (CRM)', icon: Users },
  { id: 'MACHINE', label: 'Kapasitas Mesin & Toleransi', icon: Scale },
  { id: 'INTEGRATION', label: 'Status Integrasi & Sistem', icon: Smartphone },
]

const activeTabObj = computed(() => {
  return tabs.find((t) => t.id === activeTab.value) || tabs[0]
})

// State for PIN Form
const oldPin = ref('')
const newPin = ref('')

// State for Simulation
const simulatedWeight = ref(8.5)

// Master Data Lists from API
const servicesList = ref([])
const perfumesList = ref([])
const discountsList = ref([])
const expenseCategoriesList = ref([])
const employeesList = ref([])

// Expense Audit & Correction State
const getThirtyDaysAgoStr = () => {
  const d = new Date()
  d.setDate(d.getDate() - 30)
  return d.toISOString().slice(0, 10)
}
const getTodayStr = () => new Date().toISOString().slice(0, 10)

const expenseAuditFilter = reactive({
  startDate: getThirtyDaysAgoStr(),
  endDate: getTodayStr(),
  category: '',
})
const expensesAuditList = ref([])
const pageExpensesAudit = ref(1)
const limitExpensesAudit = ref(10)
const totalExpensesAudit = ref(0)

// Customer Management State
const customerSearchQuery = ref('')
const customersList = ref([])
const pageCustomers = ref(1)
const limitCustomers = ref(10)
const totalCustomers = ref(0)

// Modal states & Loading
const modalSaving = ref(false)

// 1. Employee Modal State
const showEmployeeModal = ref(false)
const isEditingEmployee = ref(false)
const currentEmployeeId = ref(null)
const employeeForm = reactive({
  username: '',
  full_name: '',
  phone: '',
  role: 'KASIR',
  pin: '',
  is_active: true,
})

// 2. Service Modal State
const showServiceModal = ref(false)
const isEditingService = ref(false)
const currentServiceId = ref(null)
const serviceForm = reactive({
  service_code: '',
  service_name: '',
  category: 'KILOAN',
  unit: 'Kg',
  price: 10000,
  duration_day: 1,
  is_active: true,
})

// 3. Perfume Modal State
const showPerfumeModal = ref(false)
const isEditingPerfume = ref(false)
const currentPerfumeId = ref(null)
const perfumeForm = reactive({
  name: '',
  description: '',
  is_active: true,
})

// 4. Discount Modal State
const showDiscountModal = ref(false)
const isEditingDiscount = ref(false)
const currentDiscountId = ref(null)
const discountForm = reactive({
  code: '',
  name: '',
  discount_percent: 10,
  fixed_amount: 0,
  is_active: true,
})

// 5. Category Modal State
const showCategoryModal = ref(false)
const isEditingCategory = ref(false)
const currentCategoryId = ref(null)
const categoryForm = reactive({
  category_code: '',
  category_name: '',
  is_active: true,
})

// 6. Edit Expense Modal State
const showEditExpenseModal = ref(false)
const currentExpenseId = ref(null)
const editExpenseForm = reactive({
  expense_date: '',
  category: '',
  amount: 0,
  description: '',
  created_by: '',
})

// 7. Customer Modal State
const showCustomerModal = ref(false)
const isEditingCustomer = ref(false)
const currentCustomerId = ref(null)
const customerForm = reactive({
  full_name: '',
  phone: '',
  address: '',
})

// ==========================================
// API Loaders
// ==========================================
const fetchMasterData = async () => {
  try {
    const [srv, prf, disc, expCat, emp] = await Promise.all([
      api.get('/master/services').catch(() => []),
      api.get('/master/perfumes').catch(() => []),
      api.get('/master/discounts').catch(() => []),
      api.get('/master/expense-categories').catch(() => []),
      api.get('/master/employees').catch(() => []),
    ])
    if (Array.isArray(srv)) servicesList.value = srv
    if (Array.isArray(prf)) perfumesList.value = prf
    if (Array.isArray(disc)) discountsList.value = disc
    if (Array.isArray(expCat)) expenseCategoriesList.value = expCat
    if (Array.isArray(emp)) employeesList.value = emp
  } catch (err) {
    console.error('Failed to load master data:', err)
  }
}

const loadExpensesAudit = async () => {
  try {
    const params = {
      page: pageExpensesAudit.value,
      limit: limitExpensesAudit.value,
      start: expenseAuditFilter.startDate,
      end: expenseAuditFilter.endDate,
      category: expenseAuditFilter.category || undefined,
    }
    const res = await api.get('/expenses', params)
    if (res && res.data && res.meta) {
      expensesAuditList.value = res.data
      totalExpensesAudit.value = res.meta.total_items
    } else if (Array.isArray(res)) {
      expensesAuditList.value = res
      totalExpensesAudit.value = res.length
    }
  } catch (err) {
    console.error('Failed to load expenses audit:', err)
  }
}

const loadCustomers = async () => {
  try {
    const params = {
      page: pageCustomers.value,
      limit: limitCustomers.value,
      search: customerSearchQuery.value.trim() || undefined,
    }
    const res = await api.get('/customers', params)
    if (res && res.data && res.meta) {
      customersList.value = res.data
      totalCustomers.value = res.meta.total_items
    } else if (Array.isArray(res)) {
      customersList.value = res
      totalCustomers.value = res.length
    }
  } catch (err) {
    console.error('Failed to load customers:', err)
  }
}

// Watchers for Tabs & Pagination
watch(activeTab, (tab) => {
  if (tab === 'EXPENSE_CORRECTION') {
    loadExpensesAudit()
  } else if (tab === 'CUSTOMERS') {
    loadCustomers()
  }
})

watch(
  [
    () => expenseAuditFilter.startDate,
    () => expenseAuditFilter.endDate,
    () => expenseAuditFilter.category,
  ],
  () => {
    pageExpensesAudit.value = 1
    loadExpensesAudit()
  }
)

watch([pageExpensesAudit, limitExpensesAudit], () => {
  loadExpensesAudit()
})

watch(customerSearchQuery, () => {
  pageCustomers.value = 1
  loadCustomers()
})

watch([pageCustomers, limitCustomers], () => {
  loadCustomers()
})

onMounted(() => {
  fetchMasterData()
  loadExpensesAudit()
  loadCustomers()
})

// ==========================================
// Handlers: Save All & PIN
// ==========================================
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

const getRoleBadgeClass = (role) => {
  switch (role) {
    case 'OWNER':
      return 'bg-purple-100 text-purple-900 border border-purple-200'
    case 'KASIR':
      return 'bg-blue-100 text-blue-900 border border-blue-200'
    case 'KURIR':
      return 'bg-emerald-100 text-emerald-900 border border-emerald-200'
    case 'OPERATOR':
      return 'bg-amber-100 text-amber-900 border border-amber-200'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

// ==========================================
// CRUD: Pegawai (Employees)
// ==========================================
const openAddEmployeeModal = () => {
  isEditingEmployee.value = false
  currentEmployeeId.value = null
  employeeForm.username = ''
  employeeForm.full_name = ''
  employeeForm.phone = ''
  employeeForm.role = 'KASIR'
  employeeForm.pin = '1234'
  employeeForm.is_active = true
  showEmployeeModal.value = true
}

const openEditEmployeeModal = (emp) => {
  isEditingEmployee.value = true
  currentEmployeeId.value = emp.id
  employeeForm.username = emp.username
  employeeForm.full_name = emp.full_name
  employeeForm.phone = emp.phone || ''
  employeeForm.role = emp.role
  employeeForm.pin = ''
  employeeForm.is_active = emp.is_active
  showEmployeeModal.value = true
}

const submitEmployeeForm = async () => {
  modalSaving.value = true
  try {
    if (isEditingEmployee.value) {
      await api.put(`/master/employees/${currentEmployeeId.value}`, {
        username: employeeForm.username,
        full_name: employeeForm.full_name,
        phone: employeeForm.phone,
        role: employeeForm.role,
        pin: employeeForm.pin,
        is_active: employeeForm.is_active,
      })
      success('Data pegawai berhasil diperbarui!')
    } else {
      await api.post('/master/employees', {
        username: employeeForm.username,
        full_name: employeeForm.full_name,
        phone: employeeForm.phone,
        role: employeeForm.role,
        pin: employeeForm.pin,
        is_active: employeeForm.is_active,
      })
      success('Pegawai baru berhasil ditambahkan!')
    }
    showEmployeeModal.value = false
    const empList = await api.get('/master/employees')
    if (Array.isArray(empList)) employeesList.value = empList
  } catch (err) {
    error(err.message || 'Gagal menyimpan pegawai')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeleteEmployee = async (emp) => {
  if (confirm(`Yakin ingin menonaktifkan/menghapus pegawai "${emp.full_name}"?`)) {
    try {
      await api.delete(`/master/employees/${emp.id}`)
      success('Pegawai berhasil dihapus!')
      employeesList.value = employeesList.value.filter((e) => e.id !== emp.id)
    } catch (err) {
      error(err.message || 'Gagal menghapus pegawai')
    }
  }
}

// ==========================================
// CRUD: Layanan (Services)
// ==========================================
const openAddServiceModal = () => {
  isEditingService.value = false
  currentServiceId.value = null
  serviceForm.service_code = ''
  serviceForm.service_name = ''
  serviceForm.category = 'KILOAN'
  serviceForm.unit = 'Kg'
  serviceForm.price = 10000
  serviceForm.duration_day = 1
  serviceForm.is_active = true
  showServiceModal.value = true
}

const openEditServiceModal = (s) => {
  isEditingService.value = true
  currentServiceId.value = s.id
  serviceForm.service_code = s.service_code
  serviceForm.service_name = s.service_name
  serviceForm.category = s.category
  serviceForm.unit = s.unit
  serviceForm.price = Number(s.price) || 0
  serviceForm.duration_day = Number(s.duration_day) || 1
  serviceForm.is_active = s.is_active !== false
  showServiceModal.value = true
}

const submitServiceForm = async () => {
  modalSaving.value = true
  try {
    if (isEditingService.value) {
      await api.put(`/master/services/${currentServiceId.value}`, {
        service_code: serviceForm.service_code,
        service_name: serviceForm.service_name,
        category: serviceForm.category,
        unit: serviceForm.unit,
        price: serviceForm.price,
        duration_day: serviceForm.duration_day,
        is_active: serviceForm.is_active,
      })
      success('Layanan berhasil diperbarui!')
    } else {
      await api.post('/master/services', {
        service_code: serviceForm.service_code,
        service_name: serviceForm.service_name,
        category: serviceForm.category,
        unit: serviceForm.unit,
        price: serviceForm.price,
        duration_day: serviceForm.duration_day,
        is_active: serviceForm.is_active,
      })
      success('Layanan baru berhasil ditambahkan!')
    }
    showServiceModal.value = false
    const srv = await api.get('/master/services')
    if (Array.isArray(srv)) servicesList.value = srv
  } catch (err) {
    error(err.message || 'Gagal menyimpan layanan')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeleteService = async (s) => {
  if (confirm(`Yakin ingin menghapus layanan "${s.service_name}"?`)) {
    try {
      await api.delete(`/master/services/${s.id}`)
      success('Layanan berhasil dihapus!')
      servicesList.value = servicesList.value.filter((item) => item.id !== s.id)
    } catch (err) {
      error(err.message || 'Gagal menghapus layanan')
    }
  }
}

// ==========================================
// CRUD: Parfum (Perfumes)
// ==========================================
const openAddPerfumeModal = () => {
  isEditingPerfume.value = false
  currentPerfumeId.value = null
  perfumeForm.name = ''
  perfumeForm.description = ''
  perfumeForm.is_active = true
  showPerfumeModal.value = true
}

const openEditPerfumeModal = (p) => {
  isEditingPerfume.value = true
  currentPerfumeId.value = p.id
  perfumeForm.name = p.name
  perfumeForm.description = p.description || ''
  perfumeForm.is_active = p.is_active !== false
  showPerfumeModal.value = true
}

const submitPerfumeForm = async () => {
  modalSaving.value = true
  try {
    if (isEditingPerfume.value) {
      await api.put(`/master/perfumes/${currentPerfumeId.value}`, {
        name: perfumeForm.name,
        description: perfumeForm.description,
        is_active: perfumeForm.is_active,
      })
      success('Varian aroma berhasil diperbarui!')
    } else {
      await api.post('/master/perfumes', {
        name: perfumeForm.name,
        description: perfumeForm.description,
        is_active: perfumeForm.is_active,
      })
      success('Aroma parfum baru berhasil ditambahkan!')
    }
    showPerfumeModal.value = false
    const prf = await api.get('/master/perfumes')
    if (Array.isArray(prf)) perfumesList.value = prf
  } catch (err) {
    error(err.message || 'Gagal menyimpan parfum')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeletePerfume = async (p) => {
  if (confirm(`Yakin ingin menghapus parfum "${p.name}"?`)) {
    try {
      await api.delete(`/master/perfumes/${p.id}`)
      success('Parfum berhasil dihapus!')
      perfumesList.value = perfumesList.value.filter((item) => item.id !== p.id)
    } catch (err) {
      error(err.message || 'Gagal menghapus parfum')
    }
  }
}

// ==========================================
// CRUD: Diskon (Discounts)
// ==========================================
const openAddDiscountModal = () => {
  isEditingDiscount.value = false
  currentDiscountId.value = null
  discountForm.code = ''
  discountForm.name = ''
  discountForm.discount_percent = 10
  discountForm.fixed_amount = 0
  discountForm.is_active = true
  showDiscountModal.value = true
}

const openEditDiscountModal = (d) => {
  isEditingDiscount.value = true
  currentDiscountId.value = d.id
  discountForm.code = d.code
  discountForm.name = d.name
  discountForm.discount_percent = Number(d.discount_percent || d.value) || 0
  discountForm.fixed_amount = Number(d.fixed_amount || 0)
  discountForm.is_active = d.is_active !== false
  showDiscountModal.value = true
}

const submitDiscountForm = async () => {
  modalSaving.value = true
  try {
    if (isEditingDiscount.value) {
      await api.put(`/master/discounts/${currentDiscountId.value}`, {
        code: discountForm.code,
        name: discountForm.name,
        discount_percent: discountForm.discount_percent,
        fixed_amount: discountForm.fixed_amount,
        is_active: discountForm.is_active,
      })
      success('Promo diskon berhasil diperbarui!')
    } else {
      await api.post('/master/discounts', {
        code: discountForm.code,
        name: discountForm.name,
        discount_percent: discountForm.discount_percent,
        fixed_amount: discountForm.fixed_amount,
        is_active: discountForm.is_active,
      })
      success('Promo diskon baru berhasil ditambahkan!')
    }
    showDiscountModal.value = false
    const disc = await api.get('/master/discounts')
    if (Array.isArray(disc)) discountsList.value = disc
  } catch (err) {
    error(err.message || 'Gagal menyimpan diskon')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeleteDiscount = async (d) => {
  if (confirm(`Yakin ingin menghapus promo "${d.name}"?`)) {
    try {
      await api.delete(`/master/discounts/${d.id}`)
      success('Diskon berhasil dihapus!')
      discountsList.value = discountsList.value.filter((item) => item.id !== d.id)
    } catch (err) {
      error(err.message || 'Gagal menghapus diskon')
    }
  }
}

// ==========================================
// CRUD: Kategori Pengeluaran
// ==========================================
const openAddCategoryModal = () => {
  isEditingCategory.value = false
  currentCategoryId.value = null
  categoryForm.category_code = ''
  categoryForm.category_name = ''
  categoryForm.is_active = true
  showCategoryModal.value = true
}

const openEditCategoryModal = (cat) => {
  isEditingCategory.value = true
  currentCategoryId.value = cat.id
  categoryForm.category_code = cat.category_code
  categoryForm.category_name = cat.category_name
  categoryForm.is_active = cat.is_active
  showCategoryModal.value = true
}

const submitCategoryForm = async () => {
  modalSaving.value = true
  try {
    if (isEditingCategory.value) {
      await api.put(`/master/expense-categories/${currentCategoryId.value}`, {
        category_code: categoryForm.category_code,
        category_name: categoryForm.category_name,
        is_active: categoryForm.is_active,
      })
      success('Kategori pengeluaran berhasil diperbarui!')
    } else {
      await api.post('/master/expense-categories', {
        category_code: categoryForm.category_code,
        category_name: categoryForm.category_name,
        is_active: categoryForm.is_active,
      })
      success('Kategori baru berhasil ditambahkan!')
    }
    showCategoryModal.value = false
    const expCat = await api.get('/master/expense-categories')
    if (Array.isArray(expCat)) expenseCategoriesList.value = expCat
  } catch (err) {
    error(err.message || 'Gagal menyimpan kategori')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeleteCategory = async (cat) => {
  if (confirm(`Yakin ingin menghapus kategori "${cat.category_name}"?`)) {
    try {
      await api.delete(`/master/expense-categories/${cat.id}`)
      success('Kategori pengeluaran berhasil dihapus!')
      expenseCategoriesList.value = expenseCategoriesList.value.filter((item) => item.id !== cat.id)
    } catch (err) {
      error(err.message || 'Gagal menghapus kategori')
    }
  }
}

// ==========================================
// Audit & Koreksi Pengeluaran (Expense Audit)
// ==========================================
const openEditExpenseModal = (exp) => {
  currentExpenseId.value = exp.id
  editExpenseForm.expense_date = exp.expense_date ? exp.expense_date.slice(0, 10) : getTodayStr()
  editExpenseForm.category = exp.category
  editExpenseForm.amount = Number(exp.amount) || 0
  editExpenseForm.description = exp.description || exp.notes || ''
  editExpenseForm.created_by = exp.created_by || 'Kasir'
  showEditExpenseModal.value = true
}

const submitEditExpenseForm = async () => {
  modalSaving.value = true
  try {
    await api.put(`/expenses/${currentExpenseId.value}`, {
      expense_date: editExpenseForm.expense_date,
      category: editExpenseForm.category,
      amount: editExpenseForm.amount,
      description: editExpenseForm.description,
      created_by: editExpenseForm.created_by,
    })
    success('Data pengeluaran berhasil dikoreksi!')
    showEditExpenseModal.value = false
    loadExpensesAudit()
  } catch (err) {
    error(err.message || 'Gagal mengoreksi pengeluaran')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeleteExpense = async (exp) => {
  if (
    confirm(`Yakin ingin menghapus transaksi pengeluaran sebesar Rp ${formatRupiah(exp.amount)}?`)
  ) {
    try {
      await api.delete(`/expenses/${exp.id}`)
      success('Pengeluaran berhasil dihapus!')
      loadExpensesAudit()
    } catch (err) {
      error(err.message || 'Gagal menghapus pengeluaran')
    }
  }
}

// ==========================================
// CRUD: Pelanggan (Customers)
// ==========================================
const openAddCustomerModal = () => {
  isEditingCustomer.value = false
  currentCustomerId.value = null
  customerForm.full_name = ''
  customerForm.phone = ''
  customerForm.address = ''
  showCustomerModal.value = true
}

const openEditCustomerModal = (c) => {
  isEditingCustomer.value = true
  currentCustomerId.value = c.id
  customerForm.full_name = c.full_name
  customerForm.phone = c.phone
  customerForm.address = c.address || ''
  showCustomerModal.value = true
}

const submitCustomerForm = async () => {
  modalSaving.value = true
  try {
    if (isEditingCustomer.value) {
      await api.put(`/customers/${currentCustomerId.value}`, {
        full_name: customerForm.full_name,
        phone: customerForm.phone,
        address: customerForm.address,
      })
      success('Data pelanggan berhasil diperbarui!')
    } else {
      await api.post('/customers', {
        full_name: customerForm.full_name,
        phone: customerForm.phone,
        address: customerForm.address,
      })
      success('Pelanggan baru berhasil ditambahkan!')
    }
    showCustomerModal.value = false
    loadCustomers()
  } catch (err) {
    error(err.message || 'Gagal menyimpan pelanggan')
  } finally {
    modalSaving.value = false
  }
}

const confirmDeleteCustomer = async (c) => {
  if (confirm(`Yakin ingin menghapus data pelanggan "${c.full_name}"?`)) {
    try {
      await api.delete(`/customers/${c.id}`)
      success('Pelanggan berhasil dihapus!')
      loadCustomers()
    } catch (err) {
      error(err.message || 'Gagal menghapus pelanggan')
    }
  }
}
</script>
