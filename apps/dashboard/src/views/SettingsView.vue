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
              Kelola privasi kasir, template printer thermal, katalog layanan, pegawai, kas kecil, dan pelanggan.
            </p>
          </div>
        </div>

        <!-- Role Badge & Action -->
        <div class="flex items-center gap-2 shrink-0">
          <div
            class="px-3 py-1.5 rounded-xl border text-xs font-bold flex items-center gap-2 shadow-2xs"
            :class="
              settings.role === 'SUPER_ADMIN'
                ? 'bg-amber-50 text-amber-900 border-amber-200'
                : settings.role === 'OWNER'
                  ? 'bg-purple-50 text-purple-900 border-purple-200'
                  : 'bg-slate-100 text-slate-700 border-slate-200'
            "
          >
            <Crown v-if="settings.role === 'SUPER_ADMIN'" class="w-4 h-4 text-amber-500" />
            <Shield v-else-if="settings.role === 'OWNER'" class="w-4 h-4 text-purple-600" />
            <User v-else class="w-4 h-4 text-slate-500" />
            <span>Peran: {{ settings.role === 'SUPER_ADMIN' ? '⚡ Super Admin (Multi-Tenant)' : settings.role === 'OWNER' ? '👑 Owner (Akses Penuh)' : '👤 Kasir' }}</span>
          </div>
        </div>
      </div>

      <!-- Static Main Settings Layout: Left Sidebar Navigation (col-span-3) + Right Active Content (col-span-9) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        <!-- Static Sidebar Navigation for Desktop (No horizontal sliding) -->
        <div class="hidden lg:block lg:col-span-3 space-y-4 sticky top-20">
          <div class="bg-white rounded-2xl border border-slate-200/80 p-3 shadow-xs space-y-3.5">
            <div v-for="group in tabGroups" :key="group.title" class="space-y-1">
              <div class="text-[10px] font-black uppercase tracking-wider text-slate-400 px-3 py-1">
                {{ group.title }}
              </div>
              <button
                v-for="tab in group.tabs"
                :key="tab.id"
                type="button"
                class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer text-left"
                :class="
                  activeTab === tab.id
                    ? 'bg-brand-50 text-brand-700 shadow-2xs font-extrabold border border-brand-200/60'
                    : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                "
                @click="activeTab = tab.id"
              >
                <div class="flex items-center gap-2.5 truncate">
                  <component :is="tab.icon" class="w-4 h-4 shrink-0" :class="activeTab === tab.id ? 'text-brand-600' : 'text-slate-400'" />
                  <span class="truncate">{{ tab.label }}</span>
                </div>
                <span
                  v-if="tab.id === 'LEAVE_APPROVAL' && pendingLeavesCount > 0"
                  class="px-1.5 py-0.2 rounded-full text-[10px] font-black bg-amber-500 text-white animate-pulse"
                >
                  {{ pendingLeavesCount }}
                </span>
              </button>
            </div>
          </div>
        </div>

        <!-- Right Content Area (lg:col-span-9) -->
        <div class="col-span-1 lg:col-span-9 space-y-6">
          <!-- Mobile Dropdown Selector (Static, visible on mobile only, < lg) -->
          <div class="block lg:hidden bg-white p-3.5 rounded-2xl border border-slate-200/90 shadow-xs">
            <label class="block text-[10px] font-extrabold text-slate-400 uppercase tracking-wider mb-1.5">
              Pilih Menu Pengaturan
            </label>
            <div class="relative">
              <select
                v-model="activeTab"
                class="w-full pl-10 pr-10 py-3 bg-slate-50 border border-slate-200 rounded-xl font-bold text-xs text-slate-800 appearance-none focus:outline-hidden focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all cursor-pointer"
              >
                <optgroup v-for="group in tabGroups" :key="group.title" :label="group.title">
                  <option v-for="tab in group.tabs" :key="tab.id" :value="tab.id">
                    {{ tab.label }}
                  </option>
                </optgroup>
              </select>
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-brand-600">
                <component :is="activeTabObj.icon" class="w-4 h-4" />
              </div>
              <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-slate-400">
                <ChevronDown class="w-4 h-4" />
              </div>
            </div>
          </div>

          <!-- Outlet Context Selector Banner (Filter Per-Outlet Settings) -->
          <div class="bg-gradient-to-r from-brand-600 via-sky-600 to-indigo-600 rounded-3xl p-4 sm:p-5 text-white shadow-md flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-white/20 backdrop-blur-xs flex items-center justify-center shrink-0 shadow-xs">
                <Store class="w-5 h-5 text-white" />
              </div>
              <div>
                <span class="text-[10px] font-extrabold uppercase tracking-widest text-sky-200 block">Cabang Aktif Yang Dikonfigurasi</span>
                <h2 class="text-base font-black text-white flex items-center gap-2">
                  <span>{{ activeOutlet?.name || 'Cabang Belum Dipilih' }}</span>
                  <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-white/20 text-white">
                    {{ activeOutlet?.code || 'Pusat' }}
                  </span>
                </h2>
              </div>
            </div>
            <!-- Outlet Selector Dropdown -->
            <div class="flex items-center gap-2">
              <span class="text-xs text-sky-100 hidden sm:inline font-bold">Ganti Cabang:</span>
              <select
                :value="activeOutletId"
                class="px-3 py-2 bg-white text-slate-800 font-bold text-xs rounded-xl shadow-xs border border-white/40 focus:outline-hidden cursor-pointer"
                @change="setActiveOutlet($event.target.value)"
              >
                <option v-for="o in outletsList" :key="o.id" :value="o.id">
                  {{ o.name }} ({{ o.code }})
                </option>
              </select>
            </div>
          </div>

          <!-- ============================================================== -->
          <!-- TAB 0: MASTER CABANG & TITIK GPS OUTLET                         -->
          <!-- ============================================================== -->
          <div v-if="activeTab === 'OUTLETS'" class="space-y-6">
            <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div>
                <h3 class="font-black text-slate-900 text-sm flex items-center gap-2">
                  <Store class="w-4 h-4 text-brand-600" />
                  <span>Master Cabang & Titik Koordinat GPS Outlet</span>
                </h3>
                <p class="text-xs text-slate-500 mt-0.5">
                  Kelola outlet/cabang laundry di bawah akun Owner Anda, tentukan alamat dan titik GPS untuk validasi absensi selfie kasir.
                </p>
              </div>
              <button
                type="button"
                class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-2 cursor-pointer w-fit"
                @click="openAddOutletModal"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Tambah Cabang Baru</span>
              </button>
            </div>

            <!-- Outlet Cards Grid -->
            <div v-if="outletsList.length === 0" class="py-12 bg-white rounded-2xl border border-slate-100 text-center text-xs text-slate-400 space-y-2">
              <Store class="w-8 h-8 mx-auto text-slate-300" />
              <p>Belum ada data cabang outlet. Klik "Tambah Cabang Baru" untuk menambahkan.</p>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div
                v-for="outlet in outletsList"
                :key="outlet.id"
                class="bg-white rounded-2xl border border-slate-100 p-5 shadow-xs space-y-3.5 flex flex-col justify-between hover:border-slate-200 transition-all"
              >
                <div class="space-y-3">
                  <div class="flex items-center justify-between pb-2.5 border-b border-slate-100">
                    <div class="flex items-center gap-2.5">
                      <div class="w-9 h-9 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center font-black text-xs">
                        {{ outlet.code }}
                      </div>
                      <div>
                        <h4 class="font-black text-sm text-slate-900">{{ outlet.name }}</h4>
                        <p class="text-[10px] text-slate-400 font-mono">{{ outlet.phone || 'Tanpa No. WhatsApp' }}</p>
                      </div>
                    </div>
                    <span
                      class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                      :class="outlet.is_active ? 'bg-emerald-100 text-emerald-800' : 'bg-slate-100 text-slate-500'"
                    >
                      {{ outlet.is_active ? 'Aktif' : 'Non-Aktif' }}
                    </span>
                  </div>

                  <p class="text-xs text-slate-600 line-clamp-2">
                    {{ outlet.address || 'Alamat outlet belum diset.' }}
                  </p>

                  <!-- GPS Radar Info -->
                  <div class="p-3 rounded-xl bg-slate-50 border border-slate-200/80 text-[11px] space-y-1 font-mono">
                    <div class="flex items-center justify-between text-slate-700">
                      <span class="text-slate-400 flex items-center gap-1 font-sans font-bold">
                        <MapPin class="w-3 h-3 text-brand-600" /> Koordinat GPS:
                      </span>
                      <span class="font-bold">{{ outlet.latitude }}, {{ outlet.longitude }}</span>
                    </div>
                    <div class="flex items-center justify-between text-slate-700">
                      <span class="text-slate-400 font-sans font-bold">Radius Absen:</span>
                      <span class="font-bold text-emerald-700 bg-emerald-50 px-1.5 py-0.2 rounded">{{ outlet.geofence_radius_meters || 25 }} Meter</span>
                    </div>
                  </div>
                </div>

                <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-100">
                  <button
                    type="button"
                    class="px-3 py-1.5 rounded-xl border border-slate-200 text-slate-700 text-xs font-bold hover:bg-slate-50 transition-colors cursor-pointer flex items-center gap-1.5"
                    @click="openEditOutletModal(outlet)"
                  >
                    <Edit2 class="w-3.5 h-3.5 text-slate-500" />
                    <span>Edit Cabang & GPS</span>
                  </button>
                  <button
                    type="button"
                    class="p-1.5 rounded-xl border border-rose-200 text-rose-600 hover:bg-rose-50 transition-colors cursor-pointer"
                    title="Nonaktifkan Cabang"
                    @click="deleteOutlet(outlet.id, outlet.name)"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- ============================================================== -->
          <!-- TAB WA_GATEWAY: WHATSAPP BOT & GRUP NOTIFIKASI CABANG          -->
          <!-- ============================================================== -->
          <div v-if="activeTab === 'WA_GATEWAY'" class="space-y-6">
            <!-- Header Banner -->
            <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div>
                <h3 class="font-black text-slate-900 text-sm flex items-center gap-2">
                  <MessageSquare class="w-4 h-4 text-emerald-600" />
                  <span>WhatsApp Gateway & Manajemen Grup Bot Cabang</span>
                </h3>
                <p class="text-xs text-slate-500 mt-0.5">
                  1 Akun WhatsApp Bot resmi untuk melayani seluruh pelanggan. Setiap cabang dapat dihubungkan ke Grup WhatsApp kurir/staf masing-masing.
                </p>
              </div>
              <button
                type="button"
                class="px-4 py-2 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-2 cursor-pointer w-fit"
                :disabled="loadingGateway"
                @click="loadGatewayStatusAndGroups"
              >
                <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingGateway }" />
                <span>Segarkan Status</span>
              </button>
            </div>

            <!-- Connection Status Cards -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-2">
                <span class="text-[10px] uppercase font-black text-slate-400 tracking-wider">Status Koneksi Gateway</span>
                <div class="flex items-center gap-2">
                  <span
                    class="w-2.5 h-2.5 rounded-full"
                    :class="gatewayStatus.state === 'CONNECTED' ? 'bg-emerald-500 animate-pulse' : gatewayStatus.state === 'CONNECTING' ? 'bg-amber-500 animate-ping' : 'bg-rose-500'"
                  ></span>
                  <span class="font-black text-sm text-slate-800">
                    {{ gatewayStatus.state === 'CONNECTED' ? 'Terhubung (Online)' : gatewayStatus.state === 'CONNECTING' ? 'Menunggu Scan QR...' : 'Terputus (Offline)' }}
                  </span>
                </div>
                <p class="text-xs text-slate-500">
                  {{ gatewayStatus.state === 'CONNECTED' ? 'Bot siap menerima pesan & mengirim notifikasi.' : 'Scan QR code untuk menautkan nomor WhatsApp.' }}
                </p>
              </div>

              <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-2">
                <span class="text-[10px] uppercase font-black text-slate-400 tracking-wider">Nomor Akun Bot</span>
                <p class="font-black text-sm text-slate-800 font-mono">
                  {{ gatewayStatus.bot?.id ? '+' + gatewayStatus.bot.id.split(':')[0].split('@')[0] : 'Belum Tertaut' }}
                </p>
                <p class="text-xs text-slate-500">
                  Nama Bot: {{ gatewayStatus.bot?.name || 'LaundryOS Bot' }}
                </p>
              </div>

              <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-2">
                <span class="text-[10px] uppercase font-black text-slate-400 tracking-wider">Grup Cabang Aktif Ini</span>
                <p class="font-black text-sm text-brand-600 truncate">
                  {{ activeOutlet?.whatsapp_group_name || 'Belum Dihubungkan' }}
                </p>
                <p class="text-[11px] text-slate-400 font-mono truncate">
                  {{ activeOutlet?.whatsapp_group_id || 'Belum ada Group JID' }}
                </p>
              </div>
            </div>

            <!-- QR Code Section (if disconnected / connecting) -->
            <div
              v-if="gatewayStatus.state !== 'CONNECTED'"
              class="bg-white rounded-3xl border border-slate-100 p-6 sm:p-8 shadow-xs text-center space-y-6 max-w-xl mx-auto"
            >
              <div>
                <h4 class="font-black text-base text-slate-900">Tautkan Akun WhatsApp Bot</h4>
                <p class="text-xs text-slate-500 mt-1 max-w-md mx-auto">
                  Scan QR code di bawah menggunakan aplikasi WhatsApp pada smartphone resmi laundry Anda.
                </p>
              </div>

              <!-- Visual QR Code Container -->
              <div class="p-6 bg-slate-50 rounded-2xl border border-slate-200 inline-block">
                <img
                  v-if="gatewayStatus.qr_data_url"
                  :src="gatewayStatus.qr_data_url"
                  alt="WhatsApp QR Code"
                  class="w-64 h-64 mx-auto rounded-xl shadow-xs border bg-white p-2"
                />
                <div v-else class="w-64 h-64 flex flex-col items-center justify-center text-slate-400 space-y-2 bg-white rounded-xl border">
                  <div class="w-8 h-8 border-3 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
                  <span class="text-xs font-bold">Menyiapkan QR Code...</span>
                </div>
              </div>

              <!-- Steps Guide -->
              <div class="text-left bg-slate-50/80 rounded-2xl p-4 border border-slate-100 text-xs space-y-2 text-slate-700">
                <p class="font-bold text-slate-900">Cara Menautkan Perangkat:</p>
                <ol class="list-decimal list-inside space-y-1 text-slate-600">
                  <li>Buka aplikasi WhatsApp di HP Anda.</li>
                  <li>Ketuk ikon <strong>Titik Tiga (⋮)</strong> di Android atau menu <strong>Pengaturan</strong> di iPhone.</li>
                  <li>Pilih menu <strong>Perangkat Tertaut</strong> > <strong>Tautkan Perangkat</strong>.</li>
                  <li>Arahkan kamera HP ke QR Code di atas sampai terhubung.</li>
                </ol>
              </div>
            </div>

            <!-- Connected: List of Participating WhatsApp Groups -->
            <div v-else class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
              <div class="flex items-center justify-between pb-3 border-b border-slate-100">
                <div>
                  <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
                    <Users class="w-4 h-4 text-brand-600" />
                    <span>Daftar Grup WhatsApp yang Diikuti Bot</span>
                  </h4>
                  <p class="text-xs text-slate-500 mt-0.5">
                    Pilih grup kurir/staf yang akan menerima notifikasi order jemput (pickup) untuk cabang <strong>{{ activeOutlet?.name }}</strong>.
                  </p>
                </div>
                <span class="px-2.5 py-1 rounded-full bg-slate-100 text-slate-700 text-xs font-bold">
                  {{ gatewayGroups.length }} Grup Ditemukan
                </span>
              </div>

              <div v-if="gatewayGroups.length === 0" class="py-12 text-center text-xs text-slate-400 space-y-2">
                <Users class="w-8 h-8 mx-auto text-slate-300" />
                <p>Bot belum dimasukkan ke grup WhatsApp manapun.</p>
                <p class="text-slate-500">Silakan buat Grup WhatsApp di HP, lalu tambahkan nomor bot ke grup tersebut.</p>
              </div>

              <div v-else class="space-y-2.5">
                <div
                  v-for="group in gatewayGroups"
                  :key="group.id"
                  class="p-4 rounded-xl border flex flex-col sm:flex-row sm:items-center justify-between gap-3 transition-all"
                  :class="activeOutlet?.whatsapp_group_id === group.id ? 'border-brand-500 bg-brand-50/40 ring-1 ring-brand-500' : 'border-slate-200/80 hover:border-slate-300 bg-slate-50/50'"
                >
                  <div class="space-y-1">
                    <div class="flex items-center gap-2">
                      <h5 class="font-bold text-sm text-slate-900">{{ group.name }}</h5>
                      <span
                        v-if="activeOutlet?.whatsapp_group_id === group.id"
                        class="px-2 py-0.5 rounded-full text-[10px] font-black bg-brand-100 text-brand-700"
                      >
                        Grup Cabang Ini ✅
                      </span>
                    </div>
                    <div class="flex items-center gap-3 text-xs text-slate-500">
                      <span>👥 {{ group.participants_count }} Anggota</span>
                      <span class="text-slate-300">•</span>
                      <span class="font-mono text-[11px] text-slate-400">{{ group.id }}</span>
                    </div>
                  </div>

                  <div class="flex items-center gap-2 shrink-0">
                    <button
                      type="button"
                      class="px-3 py-1.5 rounded-xl border border-slate-200 text-slate-700 text-xs font-bold hover:bg-slate-100 transition-colors cursor-pointer flex items-center gap-1.5"
                      title="Salin ID Grup ke Clipboard"
                      @click="copyToClipboard(group.id, 'ID Grup WhatsApp disalin!')"
                    >
                      <Copy class="w-3.5 h-3.5 text-slate-500" />
                      <span>Salin ID</span>
                    </button>
                    <button
                      type="button"
                      class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all shadow-xs cursor-pointer flex items-center gap-1.5"
                      :class="activeOutlet?.whatsapp_group_id === group.id ? 'bg-emerald-600 text-white' : 'bg-brand-600 hover:bg-brand-700 text-white'"
                      @click="assignGroupToActiveOutlet(group)"
                    >
                      <CheckCircle2 class="w-3.5 h-3.5" />
                      <span>{{ activeOutlet?.whatsapp_group_id === group.id ? 'Terpasang' : 'Pilih untuk Cabang Ini' }}</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
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

          <!-- Save Button for Profile & Printer -->
          <div class="pt-2">
            <button
              type="button"
              class="px-4 py-2.5 rounded-xl bg-slate-900 text-white text-xs font-bold hover:bg-slate-800 transition-all cursor-pointer flex items-center gap-1.5 shadow-xs"
              @click="handleSavePrinterSettings"
            >
              <Save class="w-3.5 h-3.5" />
              <span>Simpan Konfigurasi Profil & Nota</span>
            </button>
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

            <!-- Receipt Card Container (Only this element will print) -->
            <div
              class="p-4 bg-slate-50 rounded-xl border border-slate-200/80 my-4 font-mono text-xs text-slate-900"
            >
              <div
                id="thermal-preview-receipt"
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
            class="w-full py-2.5 rounded-xl border border-slate-300 hover:bg-slate-50 text-slate-800 text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-2 shadow-xs"
            @click="testPrint"
          >
            <Printer class="w-4 h-4 text-brand-600" />
            <span>Uji Coba Cetak Struk Thermal</span>
          </button>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 2: PEGAWAI & HAK AKSES (RBAC)                              -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'EMPLOYEES'" class="space-y-6">
        <!-- Privacy Toggle & PIN Section -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <div class="lg:col-span-7 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-5">
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
                  (<code class="bg-white px-1 py-0.5 rounded border border-slate-200">Rp ••••••</code>).
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
          <div class="lg:col-span-5 bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
            <div class="border-b border-slate-100 pb-3">
              <h3 class="font-extrabold text-slate-900 text-sm flex items-center gap-2">
                <Lock class="w-4 h-4 text-purple-600" />
                <span>Ubah PIN Akses Owner</span>
              </h3>
              <p class="text-xs text-slate-500">
                PIN cepat 4-digit untuk membuka sensor keuangan.
              </p>
            </div>

            <div class="space-y-3">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">PIN Lama (Default: 1234)</label>
                <input
                  v-model="oldPin"
                  type="password"
                  maxlength="8"
                  placeholder="PIN saat ini..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono focus:bg-white focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                />
              </div>
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">PIN Baru (4-8 Digit)</label>
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
          <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3">
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
                <tr class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider">
                  <th class="py-3 px-3 cursor-pointer select-none" @click="sortEmployee('username')">
                    <div class="flex items-center gap-1">
                      <span>Username</span>
                      <component :is="getSortIcon('employee', 'username')" class="w-3 h-3 text-slate-400" />
                    </div>
                  </th>
                  <th class="py-3 px-3 cursor-pointer select-none" @click="sortEmployee('full_name')">
                    <div class="flex items-center gap-1">
                      <span>Nama Pegawai</span>
                      <component :is="getSortIcon('employee', 'full_name')" class="w-3 h-3 text-slate-400" />
                    </div>
                  </th>
                  <th class="py-3 px-3">Penugasan Outlet</th>
                  <th class="py-3 px-3">No. HP</th>
                  <th class="py-3 px-3 text-center cursor-pointer select-none" @click="sortEmployee('role')">
                    <div class="flex items-center justify-center gap-1">
                      <span>Hak Akses (Role & ACL)</span>
                      <component :is="getSortIcon('employee', 'role')" class="w-3 h-3 text-slate-400" />
                    </div>
                  </th>
                  <th class="py-3 px-3 text-center">Status</th>
                  <th class="py-3 px-3 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-if="sortedEmployees.length === 0">
                  <td colspan="7" class="py-6 text-center text-slate-400">Belum ada pegawai terdaftar.</td>
                </tr>
                <tr
                  v-for="emp in sortedEmployees"
                  :key="emp.id"
                  class="hover:bg-slate-50/70 transition-colors"
                >
                  <td class="py-2.5 px-3 font-mono font-bold text-slate-800">{{ emp.username }}</td>
                  <td class="py-2.5 px-3">
                    <p class="font-bold text-slate-900 leading-tight">{{ emp.full_name }}</p>
                    <span v-if="emp.email" class="text-[10px] text-slate-400 block">{{ emp.email }}</span>
                  </td>
                  <td class="py-2.5 px-3">
                    <span class="inline-flex items-center gap-1 text-[11px] font-bold text-slate-700 bg-slate-100 px-2 py-0.5 rounded-lg">
                      <Store class="w-3 h-3 text-brand-600" />
                      {{ emp.outlet_name || 'Semua Cabang' }}
                    </span>
                  </td>
                  <td class="py-2.5 px-3 font-medium text-slate-600">{{ emp.phone || '-' }}</td>
                  <td class="py-2.5 px-3 text-center">
                    <div class="flex flex-col items-center gap-0.5">
                      <span
                        class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                        :class="getRoleBadgeClass(emp.role)"
                      >
                        {{ emp.role }}
                      </span>
                      <span
                        v-if="emp.role === 'KASIR'"
                        class="text-[9px] font-bold text-purple-700"
                      >
                        {{ Object.values(emp.permissions || {}).filter(Boolean).length }} Hak Akses Aktif
                      </span>
                    </div>
                  </td>
                  <td class="py-2.5 px-3 text-center">
                    <span
                      class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full"
                      :class="emp.is_active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-500'"
                    >
                      <span class="w-1.5 h-1.5 rounded-full" :class="emp.is_active ? 'bg-emerald-500' : 'bg-slate-400'"></span>
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
                        @click="triggerDeleteModal('Pegawai', emp.full_name, () => deleteEmployee(emp.id))"
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
      <!-- TAB 3: MASTER LAYANAN & TARIF (GROUPED PER KATEGORI)          -->
      <!-- ============================================================== -->
      <div
        v-else-if="activeTab === 'SERVICES'"
        class="space-y-6"
      >
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h3 class="font-extrabold text-slate-900 text-base flex items-center gap-2">
              <Tag class="w-5 h-5 text-emerald-600" />
              <span>Katalog Master Layanan & Tarif</span>
            </h3>
            <p class="text-xs text-slate-500 mt-0.5">
              Kelola katalog layanan yang dikelompokkan berdasarkan kategori produk laundry.
            </p>
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="px-3 py-2 rounded-xl border border-slate-200 hover:bg-slate-50 text-slate-700 text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
              @click="openAddCategoryModal('SERVICE')"
            >
              <FolderTree class="w-3.5 h-3.5 text-slate-500" />
              <span>Kategori Baru</span>
            </button>
            <button
              type="button"
              class="px-3.5 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
              @click="openAddServiceModal"
            >
              <Plus class="w-4 h-4" />
              <span>Tambah Layanan Baru</span>
            </button>
          </div>
        </div>

        <!-- Grouped Services by Category -->
        <div v-if="serviceGroups.length === 0" class="bg-white rounded-2xl p-8 text-center text-slate-400 text-xs border border-slate-100">
          Belum ada layanan terdaftar.
        </div>

        <div
          v-for="grp in serviceGroups"
          :key="grp.categoryName"
          class="bg-white rounded-2xl border border-slate-100 overflow-hidden shadow-xs space-y-2"
        >
          <!-- Category Header -->
          <div class="bg-slate-50/80 px-5 py-3 border-b border-slate-100 flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
              <h4 class="font-black text-xs sm:text-sm text-slate-900 uppercase tracking-wider">
                {{ grp.categoryName }}
              </h4>
              <span class="px-2 py-0.5 rounded-md bg-emerald-50 text-emerald-800 text-[10px] font-bold border border-emerald-200/60">
                Satuan: {{ grp.defaultUnit }}
              </span>
            </div>
            <span class="text-xs font-extrabold text-slate-500">
              {{ grp.services.length }} Layanan
            </span>
          </div>

          <!-- Services Table inside category -->
          <div class="overflow-x-auto p-2">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr class="border-b border-slate-100 text-[11px] font-extrabold text-slate-400 uppercase tracking-wider">
                  <th class="py-2.5 px-3">Kode</th>
                  <th class="py-2.5 px-3 cursor-pointer select-none" @click="sortService(grp.categoryName, 'service_name')">
                    <div class="flex items-center gap-1">
                      <span>Nama Layanan</span>
                      <component :is="getSortIcon('service_' + grp.categoryName, 'service_name')" class="w-3 h-3 text-slate-400" />
                    </div>
                  </th>
                  <th class="py-2.5 px-3">Satuan</th>
                  <th class="py-2.5 px-3 text-right cursor-pointer select-none" @click="sortService(grp.categoryName, 'price')">
                    <div class="flex items-center justify-end gap-1">
                      <span>Tarif (Rp)</span>
                      <component :is="getSortIcon('service_' + grp.categoryName, 'price')" class="w-3 h-3 text-slate-400" />
                    </div>
                  </th>
                  <th class="py-2.5 px-3 text-center">Durasi</th>
                  <th class="py-2.5 px-3 text-center">Status</th>
                  <th class="py-2.5 px-3 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="s in grp.services"
                  :key="s.id || s.service_code"
                  class="hover:bg-slate-50/70 transition-colors"
                >
                  <td class="py-2.5 px-3 font-mono font-bold text-slate-700">{{ s.service_code }}</td>
                  <td class="py-2.5 px-3 font-bold text-slate-900">{{ s.service_name }}</td>
                  <td class="py-2.5 px-3 font-semibold text-slate-600">{{ s.unit || grp.defaultUnit }}</td>
                  <td class="py-2.5 px-3 text-right font-black text-slate-900">
                    Rp {{ formatRupiah(s.price || s.price_per_kg) }}
                  </td>
                  <td class="py-2.5 px-3 text-center font-medium text-slate-600">
                    {{ s.duration_day === 0 ? 'Hari ini' : `${s.duration_day} Hari` }}
                  </td>
                  <td class="py-2.5 px-3 text-center">
                    <span
                      class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full"
                      :class="s.is_active !== false && s.active !== false ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-500'"
                    >
                      <span class="w-1.5 h-1.5 rounded-full" :class="s.is_active !== false && s.active !== false ? 'bg-emerald-500' : 'bg-slate-400'"></span>
                      {{ s.is_active !== false && s.active !== false ? 'Aktif' : 'Nonaktif' }}
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
                        @click="triggerDeleteModal('Layanan', s.service_name, () => deleteService(s.id))"
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
              <p class="text-xs text-slate-500">
                Pilihan parfum yang tersedia untuk order cucian.
              </p>
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
                <span class="text-[11px] text-slate-500">{{ p.description || 'Aroma wangi tahan lama' }}</span>
              </div>
              <div class="flex items-center gap-1.5">
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="p.is_active !== false ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-500'"
                >
                  {{ p.is_active !== false ? 'Aktif' : 'Nonaktif' }}
                </span>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-brand-600 hover:bg-slate-200/60 cursor-pointer"
                  @click="openEditPerfumeModal(p)"
                >
                  <Edit2 class="w-3.5 h-3.5" />
                </button>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-rose-600 hover:bg-rose-50 cursor-pointer"
                  @click="triggerDeleteModal('Parfum', p.name, () => deletePerfume(p.id))"
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
                  <span class="font-mono text-[10px] font-extrabold px-1.5 py-0.2 rounded bg-slate-200 text-slate-700">
                    {{ d.code || 'PROMO' }}
                  </span>
                </div>
                <span class="text-[11px] text-slate-500">
                  {{
                    d.discount_type === 'PERCENTAGE'
                      ? `Diskon ${Number(d.value || d.discount_percent)}%`
                      : `Potongan Rp ${formatRupiah(d.value || d.fixed_amount)}`
                  }}
                  <span v-if="Number(d.min_order_amount) > 0">
                    • Min. Order Rp {{ formatRupiah(d.min_order_amount) }}
                  </span>
                </span>
              </div>
              <div class="flex items-center gap-1.5">
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="d.is_active !== false ? 'bg-rose-50 text-rose-700' : 'bg-slate-100 text-slate-500'"
                >
                  {{ d.is_active !== false ? 'Aktif' : 'Nonaktif' }}
                </span>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-brand-600 hover:bg-slate-200/60 cursor-pointer"
                  @click="openEditDiscountModal(d)"
                >
                  <Edit2 class="w-3.5 h-3.5" />
                </button>
                <button
                  type="button"
                  class="p-1 rounded-md text-slate-400 hover:text-rose-600 hover:bg-rose-50 cursor-pointer"
                  @click="triggerDeleteModal('Diskon/Promo', d.name, () => deleteDiscount(d.id))"
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
        <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3">
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
            @click="openAddCategoryModal('EXPENSE')"
          >
            <Plus class="w-4 h-4" />
            <span>Tambah Kategori Baru</span>
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider">
                <th class="py-3 px-3 cursor-pointer select-none" @click="sortCategory('code')">
                  <div class="flex items-center gap-1">
                    <span>Kode Kategori</span>
                    <component :is="getSortIcon('category', 'code')" class="w-3 h-3 text-slate-400" />
                  </div>
                </th>
                <th class="py-3 px-3 cursor-pointer select-none" @click="sortCategory('name')">
                  <div class="flex items-center gap-1">
                    <span>Nama Kategori Pengeluaran</span>
                    <component :is="getSortIcon('category', 'name')" class="w-3 h-3 text-slate-400" />
                  </div>
                </th>
                <th class="py-3 px-3 text-center">Status</th>
                <th class="py-3 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="sortedExpenseCategories.length === 0">
                <td colspan="4" class="py-6 text-center text-slate-400">Belum ada kategori pengeluaran.</td>
              </tr>
              <tr
                v-for="cat in sortedExpenseCategories"
                :key="cat.id"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-2.5 px-3 font-mono font-bold text-slate-700">
                  {{ cat.code || cat.category_code }}
                </td>
                <td class="py-2.5 px-3 font-bold text-slate-900">
                  {{ cat.name || cat.category_name }}
                </td>
                <td class="py-2.5 px-3 text-center">
                  <span
                    class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full"
                    :class="cat.is_active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-500'"
                  >
                    <span class="w-1.5 h-1.5 rounded-full" :class="cat.is_active ? 'bg-emerald-500' : 'bg-slate-400'"></span>
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
                      @click="triggerDeleteModal('Kategori', cat.name || cat.category_name, () => deleteCategory(cat.id))"
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
        <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3">
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
        <div class="p-4 bg-slate-50 rounded-xl border border-slate-200/80 grid grid-cols-1 sm:grid-cols-3 gap-3">
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
              <option v-for="cat in expenseCategoriesList" :key="cat.id" :value="cat.code || cat.category_code">
                {{ cat.name || cat.category_name }}
              </option>
            </select>
          </div>
        </div>

        <!-- Expense Audit Table -->
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider">
                <th class="py-3 px-3 cursor-pointer select-none" @click="sortExpenseAudit('expense_date')">
                  <div class="flex items-center gap-1">
                    <span>Tanggal</span>
                    <component :is="getSortIcon('expense_audit', 'expense_date')" class="w-3 h-3 text-slate-400" />
                  </div>
                </th>
                <th class="py-3 px-3">Kategori</th>
                <th class="py-3 px-3 text-right cursor-pointer select-none" @click="sortExpenseAudit('amount')">
                  <div class="flex items-center justify-end gap-1">
                    <span>Nominal (Rp)</span>
                    <component :is="getSortIcon('expense_audit', 'amount')" class="w-3 h-3 text-slate-400" />
                  </div>
                </th>
                <th class="py-3 px-3">Keterangan Biaya</th>
                <th class="py-3 px-3">Dicatat Oleh</th>
                <th class="py-3 px-3 text-center">Aksi Koreksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="sortedExpensesAudit.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400">
                  Tidak ada catatan pengeluaran pada rentang tanggal ini.
                </td>
              </tr>
              <tr
                v-for="exp in sortedExpensesAudit"
                :key="exp.id"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-2.5 px-3 font-semibold text-slate-800 whitespace-nowrap">
                  {{ formatDate(exp.expense_date) }}
                </td>
                <td class="py-2.5 px-3">
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-50 text-amber-800 border border-amber-200">
                    {{ exp.category }}
                  </span>
                </td>
                <td class="py-2.5 px-3 text-right font-extrabold text-rose-600 whitespace-nowrap">
                  Rp {{ formatRupiah(exp.amount) }}
                </td>
                <td class="py-2.5 px-3 font-medium text-slate-700 max-w-xs truncate">
                  {{ exp.description || exp.notes || '-' }}
                </td>
                <td class="py-2.5 px-3 text-slate-500 font-medium whitespace-nowrap">
                  <span class="px-2 py-0.5 rounded-md bg-slate-100 text-slate-700 font-semibold text-[10px]">
                    {{ exp.created_by || 'Kasir' }}
                  </span>
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
                      @click="triggerDeleteModal('Pengeluaran', `Rp ${formatRupiah(exp.amount)}`, () => deleteExpense(exp.id))"
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
        <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3">
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
            class="absolute right-2.5 top-2.5 text-slate-400 hover:text-slate-600 cursor-pointer"
            @click="customerSearchQuery = ''"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-y border-slate-200/80 text-[11px] font-extrabold text-slate-500 uppercase tracking-wider">
                <th class="py-3 px-3 cursor-pointer select-none" @click="sortCustomer('full_name')">
                  <div class="flex items-center gap-1">
                    <span>Nama Pelanggan</span>
                    <component :is="getSortIcon('customer', 'full_name')" class="w-3 h-3 text-slate-400" />
                  </div>
                </th>
                <th class="py-3 px-3">No. WhatsApp</th>
                <th class="py-3 px-3">Alamat</th>
                <th class="py-3 px-3 cursor-pointer select-none" @click="sortCustomer('created_at')">
                  <div class="flex items-center gap-1">
                    <span>Terdaftar</span>
                    <component :is="getSortIcon('customer', 'created_at')" class="w-3 h-3 text-slate-400" />
                  </div>
                </th>
                <th class="py-3 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="sortedCustomers.length === 0">
                <td colspan="5" class="py-8 text-center text-slate-400">
                  Tidak ada data pelanggan yang cocok.
                </td>
              </tr>
              <tr
                v-for="c in sortedCustomers"
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
                      @click="triggerDeleteModal('Pelanggan', c.full_name, () => deleteCustomer(c.id))"
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
            class="w-full py-2.5 rounded-xl border border-rose-300 text-rose-700 hover:bg-rose-50 text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-2 shadow-xs"
            @click="handleFactoryReset"
          >
            <RotateCcw class="w-4 h-4" />
            <span>Kembalikan ke Default Pabrik</span>
          </button>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB: MASTER SHIFT & JAM KERJA                                  -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'SHIFTS_ATTENDANCE'" class="space-y-6">
        <!-- Header Info -->
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h3 class="font-black text-slate-900 text-sm flex items-center gap-2">
              <Clock class="w-4 h-4 text-brand-600" />
              <span>Master Shift Operasional & Toleransi Keterlambatan</span>
            </h3>
            <p class="text-xs text-slate-500 mt-0.5">
              Konfigurasi jam masuk, jam pulang, istirahat, dan toleransi keterlambatan shift FAN Laundry.
            </p>
          </div>
          <button
            type="button"
            class="px-4 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold transition-all flex items-center gap-2 cursor-pointer w-fit"
            @click="loadShifts"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            <span>Muat Ulang Shift</span>
          </button>
        </div>

        <!-- Shift Cards Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div
            v-for="shift in shiftsList"
            :key="shift.id"
            class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4 flex flex-col justify-between"
          >
            <div class="space-y-4">
              <div class="flex items-center justify-between pb-3 border-b border-slate-100">
                <div class="flex items-center gap-2.5">
                  <div class="w-9 h-9 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center font-black text-sm">
                    {{ shift.shift_number }}
                  </div>
                  <div>
                    <h4 class="font-extrabold text-sm text-slate-900">{{ shift.shift_name }}</h4>
                    <p class="text-[11px] text-slate-400">Jadwal Standar: {{ shift.start_time }} - {{ shift.end_time }} WIB</p>
                  </div>
                </div>
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                  :class="shift.is_active ? 'bg-emerald-100 text-emerald-800' : 'bg-slate-100 text-slate-500'"
                >
                  {{ shift.is_active ? 'Aktif' : 'Non-Aktif' }}
                </span>
              </div>

              <!-- Shift Form Inputs -->
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-[11px] font-bold text-slate-700 mb-1">Jam Masuk (Start)</label>
                  <input
                    v-model="shift.start_time"
                    type="time"
                    step="1"
                    class="w-full px-3 py-2 text-xs font-bold rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                  />
                </div>
                <div>
                  <label class="block text-[11px] font-bold text-slate-700 mb-1">Jam Pulang (End)</label>
                  <input
                    v-model="shift.end_time"
                    type="time"
                    step="1"
                    class="w-full px-3 py-2 text-xs font-bold rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                  />
                </div>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-[11px] font-bold text-slate-700 mb-1">Mulai Istirahat</label>
                  <input
                    v-model="shift.break_start"
                    type="time"
                    step="1"
                    class="w-full px-3 py-2 text-xs font-medium rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                  />
                </div>
                <div>
                  <label class="block text-[11px] font-bold text-slate-700 mb-1">Selesai Istirahat</label>
                  <input
                    v-model="shift.break_end"
                    type="time"
                    step="1"
                    class="w-full px-3 py-2 text-xs font-medium rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                  />
                </div>
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">
                  Toleransi Keterlambatan (Menit)
                </label>
                <div class="relative">
                  <input
                    v-model.number="shift.late_tolerance_minutes"
                    type="number"
                    min="0"
                    max="60"
                    class="w-full px-3 py-2 text-xs font-bold rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500 pr-16"
                  />
                  <span class="absolute right-3 top-2 text-xs text-slate-400 font-medium pointer-events-none">
                    Menit
                  </span>
                </div>
                <p class="text-[10px] text-slate-400 mt-1">
                  Check-in lewat dari jam masuk + toleransi akan tercatat "Terlambat".
                </p>
              </div>
            </div>

            <button
              type="button"
              class="w-full mt-4 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all shadow-xs flex items-center justify-center gap-2 cursor-pointer"
              @click="saveShiftSettings(shift)"
            >
              <Save class="w-3.5 h-3.5" />
              <span>Simpan Perubahan {{ shift.shift_name }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB: PERSETUJUAN CUTI (LEAVE APPROVAL)                         -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'LEAVE_APPROVAL'" class="space-y-6">
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3">
            <div>
              <div class="flex items-center gap-2">
                <CalendarCheck class="w-5 h-5 text-indigo-600" />
                <h3 class="font-black text-sm text-slate-900">Persetujuan Cuti Karyawan (Leave Requests)</h3>
                <span
                  v-if="pendingLeavesCount > 0"
                  class="px-2 py-0.5 rounded-full text-[10px] font-black bg-amber-500 text-white animate-pulse"
                >
                  {{ pendingLeavesCount }} Menunggu
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-0.5">
                Jatah libur karyawan: maksimal 2x / bulan kalender. Review pengajuan libur staf Anda.
              </p>
            </div>

            <!-- Filter Status -->
            <div class="flex items-center gap-2">
              <select
                v-model="leaveFilterStatus"
                class="px-3 py-1.5 rounded-xl border border-slate-200 text-xs font-bold text-slate-700 bg-slate-50 focus:bg-white cursor-pointer"
                @change="loadLeaves"
              >
                <option value="">Semua Status</option>
                <option value="PENDING">Menunggu (Pending)</option>
                <option value="APPROVED">Disetujui (Approved)</option>
                <option value="REJECTED">Ditolak (Rejected)</option>
              </select>
              <button
                type="button"
                class="p-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-600 transition-colors cursor-pointer"
                title="Refresh Pengajuan Cuti"
                @click="loadLeaves"
              >
                <RotateCcw class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- Leaves List -->
          <div v-if="loadingLeaves" class="py-12 text-center text-xs text-slate-400">
            Memuat pengajuan libur...
          </div>
          <div
            v-else-if="leavesList.length === 0"
            class="py-8 text-center text-xs text-slate-400"
          >
            Tidak ada data pengajuan libur dengan filter saat ini.
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3.5">
            <div
              v-for="leave in leavesList"
              :key="leave.id"
              class="p-4 rounded-2xl border transition-all space-y-3"
              :class="
                leave.status === 'PENDING'
                  ? 'bg-amber-50/50 border-amber-200 shadow-xs'
                  : leave.status === 'APPROVED'
                    ? 'bg-emerald-50/30 border-emerald-200'
                    : 'bg-slate-50 border-slate-200'
              "
            >
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="font-extrabold text-xs text-slate-900">{{ leave.employee_name }}</h4>
                  <p class="text-[10px] text-slate-400">{{ formatDate(leave.created_at) }}</p>
                </div>
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                  :class="
                    leave.status === 'APPROVED'
                      ? 'bg-emerald-100 text-emerald-800'
                      : leave.status === 'REJECTED'
                        ? 'bg-rose-100 text-rose-800'
                        : 'bg-amber-100 text-amber-800'
                  "
                >
                  {{ leave.status === 'APPROVED' ? 'Disetujui ✅' : leave.status === 'REJECTED' ? 'Ditolak ❌' : 'Pending ⏳' }}
                </span>
              </div>

              <div class="p-2.5 rounded-xl bg-white/80 border border-slate-100 text-xs space-y-1">
                <div class="flex items-center justify-between text-[11px] text-slate-500">
                  <span>Tanggal Libur:</span>
                  <span class="font-black text-indigo-700">{{ formatDate(leave.date) }}</span>
                </div>
                <p class="text-[11px] text-slate-700 italic pt-1">
                  "{{ leave.reason }}"
                </p>
              </div>

              <div v-if="leave.reviewer_notes" class="text-[10px] text-slate-500">
                <span class="font-bold">Catatan Review:</span> {{ leave.reviewer_notes }}
              </div>

              <!-- Owner Action Buttons for PENDING -->
              <div v-if="leave.status === 'PENDING'" class="flex items-center gap-2 pt-1 border-t border-slate-200/50">
                <button
                  type="button"
                  class="flex-1 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs transition-colors cursor-pointer flex items-center justify-center gap-1 shadow-xs"
                  @click="handleApproveLeave(leave)"
                >
                  <CheckCircle2 class="w-3.5 h-3.5" />
                  <span>Setujui</span>
                </button>
                <button
                  type="button"
                  class="flex-1 py-1.5 rounded-xl bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 font-bold text-xs transition-colors cursor-pointer flex items-center justify-center gap-1"
                  @click="openRejectLeaveModal(leave)"
                >
                  <X class="w-3.5 h-3.5" />
                  <span>Tolak</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB: REKAP ABSENSI KARYAWAN                                    -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'ATTENDANCE_LOG'" class="space-y-6">
        <div class="bg-white rounded-2xl border border-slate-100 p-5 sm:p-6 shadow-xs space-y-4">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-100 gap-3">
            <div>
              <h3 class="font-black text-sm text-slate-900 flex items-center gap-2">
                <Camera class="w-4 h-4 text-emerald-600" />
                <span>Log Kehadiran, Foto Selfie & Serah Terima Shift</span>
              </h3>
              <p class="text-xs text-slate-500 mt-0.5">
                Riwayat absen masuk, foto di lokasi outlet, GPS radius, dan saldo kas awal/akhir kasir.
              </p>
            </div>

            <!-- Filter Controls -->
            <div class="flex items-center gap-2 flex-wrap">
              <input
                v-model="attendanceFilter.employee"
                type="text"
                placeholder="Cari nama..."
                class="px-3 py-1.5 rounded-xl border border-slate-200 text-xs font-semibold text-slate-700 bg-slate-50 focus:bg-white w-28 sm:w-36"
                @keyup.enter="loadAttendances"
              />
              <button
                type="button"
                class="px-3 py-1.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all cursor-pointer flex items-center gap-1 shadow-xs"
                @click="loadAttendances"
              >
                <Search class="w-3.5 h-3.5" />
                <span>Filter</span>
              </button>
            </div>
          </div>

          <!-- Attendance Table -->
          <div v-if="loadingAttendances" class="py-12 text-center text-xs text-slate-400">
            Memuat riwayat absensi...
          </div>
          <div
            v-else-if="attendanceList.length === 0"
            class="py-12 text-center text-xs text-slate-400 space-y-2"
          >
            <Camera class="w-8 h-8 mx-auto text-slate-300" />
            <p>Belum ada rekaman absensi untuk periode ini.</p>
          </div>
          <div v-else class="overflow-x-auto rounded-xl border border-slate-200">
            <table class="w-full text-left text-xs">
              <thead class="bg-slate-50 border-b border-slate-200 text-slate-600 font-extrabold uppercase text-[10px] tracking-wider">
                <tr>
                  <th class="px-4 py-3">Foto Selfie</th>
                  <th class="px-4 py-3">Karyawan & Shift</th>
                  <th class="px-4 py-3">Jam Masuk (Check-In)</th>
                  <th class="px-4 py-3">Jam Pulang (Check-Out)</th>
                  <th class="px-4 py-3">Status Ketepatan</th>
                  <th class="px-4 py-3">Kas Laci Awal / Akhir</th>
                  <th class="px-4 py-3">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-slate-700">
                <tr v-for="att in attendanceList" :key="att.id" class="hover:bg-slate-50/70 transition-colors">
                  <!-- Foto Thumbnail -->
                  <td class="px-4 py-3">
                    <div
                      v-if="att.check_in_photo"
                      class="w-12 h-12 rounded-xl overflow-hidden border border-slate-200 cursor-pointer relative group shadow-2xs"
                      @click="openPhotoZoom(att)"
                    >
                      <img :src="att.check_in_photo" class="w-full h-full object-cover group-hover:scale-110 transition-transform" />
                      <div class="absolute inset-0 bg-slate-900/30 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity text-white">
                        <Eye class="w-3.5 h-3.5" />
                      </div>
                    </div>
                    <div v-else class="w-12 h-12 rounded-xl bg-slate-100 border border-slate-200 flex items-center justify-center text-slate-400">
                      <Camera class="w-4 h-4" />
                    </div>
                  </td>

                  <!-- Karyawan & Shift -->
                  <td class="px-4 py-3">
                    <p class="font-extrabold text-slate-900 text-sm">{{ att.employee_name }}</p>
                    <span class="text-[11px] font-bold text-brand-600">{{ att.shift_name }}</span>
                  </td>

                  <!-- Check In Time & Geofence -->
                  <td class="px-4 py-3">
                    <p class="font-bold text-slate-800">{{ att.clock_in_time ? formatDate(att.clock_in_time) : '-' }}</p>
                    <span
                      v-if="att.distance_meters != null"
                      class="inline-flex items-center gap-1 text-[10px] font-bold text-emerald-700 bg-emerald-50 px-1.5 py-0.5 rounded-md mt-0.5"
                    >
                      <MapPin class="w-3 h-3 text-emerald-600" />
                      {{ att.distance_meters.toFixed(1) }}m dari outlet
                    </span>
                  </td>

                  <!-- Check Out Time -->
                  <td class="px-4 py-3">
                    <p class="font-bold text-slate-800">{{ att.clock_out_time ? formatDate(att.clock_out_time) : 'Masih Bertugas' }}</p>
                    <span v-if="!att.clock_out_time" class="inline-flex items-center gap-1 text-[10px] font-bold text-amber-700 bg-amber-50 px-1.5 py-0.5 rounded-md mt-0.5 animate-pulse">
                      Aktif Sekarang
                    </span>
                  </td>

                  <!-- Punctuality Status -->
                  <td class="px-4 py-3">
                    <span
                      class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider inline-block"
                      :class="
                        att.punctuality_status === 'ON_TIME'
                          ? 'bg-emerald-100 text-emerald-800'
                          : att.punctuality_status === 'LATE'
                            ? 'bg-rose-100 text-rose-800'
                            : 'bg-amber-100 text-amber-800'
                      "
                    >
                      {{
                        att.punctuality_status === 'ON_TIME'
                          ? 'Tepat Waktu'
                          : att.punctuality_status === 'LATE'
                            ? `Terlambat ${att.late_minutes} Mnt`
                            : att.punctuality_status
                      }}
                    </span>
                  </td>

                  <!-- Cash Drawer -->
                  <td class="px-4 py-3 font-mono">
                    <div class="text-[11px]">
                      <span class="text-slate-400">Awal:</span> {{ formatRupiah(att.starting_cash || 0) }}
                    </div>
                    <div v-if="att.ending_cash != null" class="text-[11px]">
                      <span class="text-slate-400">Akhir:</span> {{ formatRupiah(att.ending_cash) }}
                    </div>
                  </td>

                  <!-- Action -->
                  <td class="px-4 py-3">
                    <button
                      v-if="att.check_in_photo"
                      type="button"
                      class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold transition-colors cursor-pointer flex items-center gap-1"
                      @click="openPhotoZoom(att)"
                    >
                      <Eye class="w-3.5 h-3.5" />
                      <span>Foto</span>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>

      <!-- ============================================================== -->
      <!-- MODALS SECTION                                                 -->
      <!-- ============================================================== -->

      <!-- 0. Modal Tambah/Edit Outlet Cabang & GPS -->
      <Teleport to="body">
        <div v-if="showOutletModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
          <div class="bg-white rounded-3xl max-w-xl w-full p-5 shadow-2xl border border-slate-100 space-y-4 max-h-[92vh] overflow-y-auto">
            <div class="flex items-center justify-between pb-2 border-b border-slate-100 sticky top-0 bg-white z-10">
              <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
                <Store class="w-4 h-4 text-brand-600" />
                <span>{{ isEditingOutlet ? 'Edit Profil & Koordinat Cabang' : 'Tambah Outlet Cabang Baru' }}</span>
              </h4>
              <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showOutletModal = false">
                <X class="w-4 h-4" />
              </button>
            </div>
            <form class="space-y-3" @submit.prevent="submitOutletForm">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Nama Outlet *</label>
                <input
                  v-model="outletForm.name"
                  type="text"
                  required
                  placeholder="Contoh: FAN Laundry Kalibaru"
                  class="w-full px-3 py-2 text-xs font-bold rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                />
              </div>
              <div class="grid grid-cols-2 gap-2">
                <div>
                  <label class="block text-xs font-bold text-slate-700 mb-1">Kode Cabang *</label>
                  <input
                    v-model="outletForm.code"
                    type="text"
                    required
                    placeholder="Contoh: OUT-01"
                    class="w-full px-3 py-2 text-xs font-mono font-bold rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500 uppercase"
                  />
                </div>
                <div>
                  <label class="block text-xs font-bold text-slate-700 mb-1">No. WhatsApp</label>
                  <input
                    v-model="outletForm.phone"
                    type="text"
                    placeholder="Contoh: 0852-8597-7834"
                    class="w-full px-3 py-2 text-xs font-medium rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                  />
                </div>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Alamat Outlet</label>
                <textarea
                  v-model="outletForm.address"
                  rows="2"
                  placeholder="Alamat lengkap outlet cabang..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-brand-500"
                ></textarea>
              </div>

              <!-- Interactive GPS Map Picker -->
              <div class="p-3 rounded-2xl bg-slate-50 border border-slate-200 space-y-3">
                <LocationPickerMap
                  v-model:latitude="outletForm.latitude"
                  v-model:longitude="outletForm.longitude"
                  :radius-meters="outletForm.geofence_radius_meters"
                />

                <div class="grid grid-cols-2 gap-2 pt-1 border-t border-slate-200/70">
                  <div>
                    <label class="block text-[10px] font-bold text-slate-500 mb-0.5">Latitude (Manual/Otomatis)</label>
                    <input
                      v-model.number="outletForm.latitude"
                      type="number"
                      step="0.000001"
                      required
                      placeholder="-6.113425"
                      class="w-full px-2.5 py-1.5 text-xs font-mono font-bold rounded-lg border border-slate-200 bg-white"
                    />
                  </div>
                  <div>
                    <label class="block text-[10px] font-bold text-slate-500 mb-0.5">Longitude (Manual/Otomatis)</label>
                    <input
                      v-model.number="outletForm.longitude"
                      type="number"
                      step="0.000001"
                      required
                      placeholder="106.914230"
                      class="w-full px-2.5 py-1.5 text-xs font-mono font-bold rounded-lg border border-slate-200 bg-white"
                    />
                  </div>
                </div>
                <div>
                  <label class="block text-[10px] font-bold text-slate-500 mb-0.5">Radius Geofence Absensi (Meter)</label>
                  <input
                    v-model.number="outletForm.geofence_radius_meters"
                    type="number"
                    min="5"
                    max="500"
                    required
                    placeholder="25"
                    class="w-full px-2.5 py-1.5 text-xs font-bold rounded-lg border border-slate-200 bg-white"
                  />
                  <p class="text-[10px] text-slate-400 mt-1">Disarankan 20-30 meter untuk toleransi sinyal GPS saat karyawan berada di dalam ruangan toko.</p>
                </div>
              </div>

              <!-- WhatsApp Gateway Group Integration -->
              <div class="p-3 rounded-2xl bg-emerald-50/60 border border-emerald-200/80 space-y-2">
                <span class="block text-xs font-bold text-emerald-900 flex items-center gap-1.5">
                  <MessageSquare class="w-3.5 h-3.5 text-emerald-600" />
                  <span>Integrasi Grup WhatsApp Cabang (Bot Notifikasi)</span>
                </span>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <div>
                    <label class="block text-[10px] font-bold text-slate-500 mb-0.5">Nama WhatsApp Group</label>
                    <input
                      v-model="outletForm.whatsapp_group_name"
                      type="text"
                      placeholder="Contoh: Tim FAN Cabang 2"
                      class="w-full px-2.5 py-1.5 text-xs font-medium rounded-lg border border-slate-200 bg-white"
                    />
                  </div>
                  <div>
                    <label class="block text-[10px] font-bold text-slate-500 mb-0.5">WhatsApp Group ID (JID)</label>
                    <input
                      v-model="outletForm.whatsapp_group_id"
                      type="text"
                      placeholder="Contoh: 12036302...@g.us"
                      class="w-full px-2.5 py-1.5 text-xs font-mono font-bold rounded-lg border border-slate-200 bg-white"
                    />
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-2 pt-1">
                <input
                  id="outlet-active"
                  v-model="outletForm.is_active"
                  type="checkbox"
                  class="rounded border-slate-300 text-brand-600 focus:ring-brand-500 w-4 h-4 cursor-pointer"
                />
                <label for="outlet-active" class="text-xs font-bold text-slate-700 cursor-pointer">
                  Outlet Aktif Beroperasi
                </label>
              </div>
              <div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
                <button
                  type="button"
                  class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                  @click="showOutletModal = false"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  :disabled="modalSaving"
                  class="px-4 py-1.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
                >
                  {{ modalSaving ? 'Menyimpan...' : 'Simpan Cabang' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Teleport>

      <!-- 1. Modal Tambah/Edit Pegawai & Granular ACL -->
      <Teleport to="body">
        <div v-if="showEmployeeModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
          <div class="bg-white rounded-3xl max-w-xl w-full p-5 sm:p-6 shadow-2xl border border-slate-100 space-y-4 max-h-[92vh] overflow-y-auto">
            <div class="flex items-center justify-between pb-3 border-b border-slate-100 sticky top-0 bg-white z-10">
              <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
                <Users class="w-4 h-4 text-purple-600" />
                <span>{{ isEditingEmployee ? 'Edit Data Pegawai & Hak Akses' : 'Tambah Pegawai Baru' }}</span>
              </h4>
              <button class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer" @click="showEmployeeModal = false">
                <X class="w-4 h-4" />
              </button>
            </div>
            <form class="space-y-4" @submit.prevent="submitEmployeeForm">
            <!-- Data Diri Pegawai -->
            <div class="space-y-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Nama Pegawai *</label>
                <input
                  v-model="employeeForm.full_name"
                  type="text"
                  required
                  placeholder="Contoh: Staff 1"
                  class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500 transition-all"
                />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs font-bold text-slate-700 mb-1">Username Login *</label>
                  <input
                    v-model="employeeForm.username"
                    type="text"
                    required
                    :disabled="isEditingEmployee"
                    placeholder="Contoh: kasir1"
                    class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500 disabled:opacity-60"
                  />
                </div>
                <div>
                  <label class="block text-xs font-bold text-slate-700 mb-1">No Handphone</label>
                  <input
                    v-model="employeeForm.phone"
                    type="text"
                    placeholder="Contoh: 085781258821"
                    class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Email</label>
                <input
                  v-model="employeeForm.email"
                  type="email"
                  placeholder="Contoh: staff@fanlaundry.com"
                  class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-medium focus:bg-white focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">
                  {{ isEditingEmployee ? 'Password (harap kosongkan jika tidak merubah password)' : 'Password Akun Login *' }}
                </label>
                <input
                  v-model="employeeForm.password"
                  type="password"
                  :required="!isEditingEmployee"
                  placeholder="Masukkan Password untuk pegawai login"
                  class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-mono focus:bg-white focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs font-bold text-slate-700 mb-1">Penugasan Outlet / Cabang</label>
                  <select
                    v-model="employeeForm.outlet_id"
                    class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500"
                  >
                    <option value="">-- Semua Cabang (Pusat) --</option>
                    <option v-for="o in outletsList" :key="o.id" :value="o.id">
                      {{ o.name }} ({{ o.code }})
                    </option>
                  </select>
                </div>
                <div>
                  <label class="block text-xs font-bold text-slate-700 mb-1">Peran Jabatan *</label>
                  <select
                    v-model="employeeForm.role"
                    required
                    class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-purple-500"
                  >
                    <option value="KASIR">👤 Kasir (POS & Operasional)</option>
                    <option value="OWNER">👑 Owner (Akses Penuh Cabang)</option>
                    <option value="KURIR">🛵 Kurir (Antar Jemput)</option>
                    <option value="OPERATOR">⚙️ Operator (Produksi Cucian)</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Granular ACL Permissions (Checklist Sesuai Contoh Edobi) -->
            <div class="p-4 bg-slate-50 border border-slate-200/90 rounded-2xl space-y-3.5">
              <div class="flex items-center justify-between border-b border-slate-200 pb-2">
                <div>
                  <h5 class="text-xs font-black text-slate-900 flex items-center gap-1.5">
                    <ShieldCheck class="w-4 h-4 text-purple-600" />
                    <span>Hak Akses Pegawai</span>
                  </h5>
                  <p class="text-[10px] text-slate-500">Centang fitur yang diizinkan untuk dikerjakan pegawai ini</p>
                </div>
                <button
                  type="button"
                  class="text-[10px] font-bold text-purple-600 hover:text-purple-800 cursor-pointer"
                  @click="employeeForm.permissions = { ...defaultAcl }"
                >
                  Reset Default Kasir
                </button>
              </div>

              <div v-for="grp in aclPermissionsList" :key="grp.group" class="space-y-2">
                <span class="text-[10px] font-black uppercase tracking-wider text-purple-700/80 block mt-2">
                  {{ grp.group }}
                </span>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pl-1">
                  <label
                    v-for="p in grp.permissions"
                    :key="p.key"
                    class="flex items-center gap-2 p-2 rounded-xl bg-white border border-slate-200/70 hover:border-purple-300 transition-all cursor-pointer select-none"
                  >
                    <input
                      v-model="employeeForm.permissions[p.key]"
                      type="checkbox"
                      class="rounded border-slate-300 text-purple-600 focus:ring-purple-500 w-4 h-4 cursor-pointer"
                    />
                    <span class="text-[11px] font-bold text-slate-700 leading-tight">{{ p.label }}</span>
                  </label>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-2 pt-1">
              <input
                id="emp-active"
                v-model="employeeForm.is_active"
                type="checkbox"
                class="rounded border-slate-300 text-purple-600 focus:ring-purple-500 w-4 h-4 cursor-pointer"
              />
              <label for="emp-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Status Akun Aktif (Dapat Login ke Aplikasi)
              </label>
            </div>

            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100 sticky bottom-0 bg-white py-2">
              <button
                type="button"
                class="px-4 py-2 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showEmployeeModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-5 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-xs font-black transition-all shadow-md shadow-purple-500/20 cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'SIMPAN DATA PEGAWAI' }}
              </button>
            </div>
          </form>
        </div>
      </div>
      </Teleport>

      <!-- 2. Modal Tambah/Edit Layanan (Flow: Pick Category First -> Auto Code & Unit) -->
      <div v-if="showServiceModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Tag class="w-4 h-4 text-emerald-600" />
              <span>{{ isEditingService ? 'Edit Layanan' : 'Tambah Layanan Baru' }}</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showServiceModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3.5" @submit.prevent="submitServiceForm">
            <!-- Step 1: Kategori Layanan (Wajib Pilih Dulu) -->
            <div>
              <label class="block text-xs font-black text-slate-800 mb-1 flex items-center justify-between">
                <span>1. Pilih Kategori Layanan *</span>
                <span class="text-[10px] text-emerald-700 font-bold">Langkah Pertama</span>
              </label>
              <select
                v-model="serviceForm.category_id"
                required
                class="w-full px-3 py-2.5 text-xs rounded-xl border border-emerald-300 bg-emerald-50/40 font-bold text-slate-900 focus:bg-white focus:ring-2 focus:ring-emerald-500 cursor-pointer"
                @change="handleCategorySelectionChange"
              >
                <option value="" disabled>-- Pilih Kategori Layanan --</option>
                <option v-for="cat in serviceCategoriesList" :key="cat.id" :value="cat.id">
                  {{ cat.name }} (Satuan: {{ cat.default_unit }})
                </option>
              </select>
            </div>

            <!-- Step 2: Auto-Generated Code & Satuan -->
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">
                  Kode Layanan <span class="text-[10px] text-slate-400 font-normal">(Otomatis)</span>
                </label>
                <input
                  v-model="serviceForm.service_code"
                  type="text"
                  required
                  readonly
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-mono font-bold uppercase text-slate-600 cursor-not-allowed select-none"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Satuan Tarif *</label>
                <input
                  v-model="serviceForm.unit"
                  type="text"
                  required
                  placeholder="Kg / Tabung / Pcs"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>

            <!-- Step 3: Nama Layanan -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Layanan *</label>
              <input
                v-model="serviceForm.service_name"
                type="text"
                required
                placeholder="Contoh: Cuci Kering Seterika Kilat"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <!-- Step 4: Tarif & Durasi -->
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Tarif Layanan (Rp) *</label>
                <input
                  v-model.number="serviceForm.price"
                  type="number"
                  min="500"
                  step="500"
                  required
                  placeholder="10000"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-black text-slate-900 focus:bg-white focus:ring-2 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Durasi Kerja (Hari)</label>
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
                v-model="serviceForm.is_active"
                type="checkbox"
                id="srv-active"
                class="rounded text-emerald-600 focus:ring-emerald-500 w-4 h-4 cursor-pointer"
              />
              <label for="srv-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Layanan Aktif (Dapat dipilih kasir di Web POS)
              </label>
            </div>

            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showServiceModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Layanan' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 3. Modal Tambah/Edit Parfum -->
      <div v-if="showPerfumeModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Sparkles class="w-4 h-4 text-amber-500" />
              <span>{{ isEditingPerfume ? 'Edit Aroma Parfum' : 'Tambah Aroma Baru' }}</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showPerfumeModal = false">
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
                v-model="perfumeForm.is_active"
                type="checkbox"
                id="prf-active"
                class="rounded text-amber-600 focus:ring-amber-500 w-4 h-4 cursor-pointer"
              />
              <label for="prf-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Tersedia untuk Cucian
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showPerfumeModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Parfum' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 4. Modal Tambah/Edit Diskon (Mutually Exclusive: Persentase vs Potongan Tetap) -->
      <div v-if="showDiscountModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Percent class="w-4 h-4 text-rose-500" />
              <span>{{ isEditingDiscount ? 'Edit Diskon / Promo' : 'Tambah Promo Baru' }}</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showDiscountModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3.5" @submit.prevent="submitDiscountForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">
                Kode Promo <span class="text-[10px] text-slate-400 font-normal">(Otomatis oleh Sistem)</span>
              </label>
              <input
                v-model="discountForm.code"
                type="text"
                required
                readonly
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-mono font-bold uppercase text-slate-600 cursor-not-allowed select-none"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Diskon / Promo *</label>
              <input
                v-model="discountForm.name"
                type="text"
                required
                placeholder="Contoh: Diskon Pelanggan Baru"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-rose-500"
                @input="handleDiscountNameInput"
              />
            </div>

            <!-- Tipe Diskon (Mutually Exclusive Toggle) -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1.5">
                Bentuk Potongan Diskon * (Pilih Salah Satu)
              </label>
              <div class="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  class="py-2 px-3 rounded-xl border text-xs font-bold transition-all text-center cursor-pointer"
                  :class="
                    discountForm.type === 'PERCENTAGE'
                      ? 'bg-rose-50 border-rose-400 text-rose-700 ring-2 ring-rose-400/20'
                      : 'bg-slate-50 border-slate-200 text-slate-600 hover:bg-slate-100'
                  "
                  @click="setDiscountType('PERCENTAGE')"
                >
                  Persentase (%)
                </button>
                <button
                  type="button"
                  class="py-2 px-3 rounded-xl border text-xs font-bold transition-all text-center cursor-pointer"
                  :class="
                    discountForm.type === 'NOMINAL'
                      ? 'bg-rose-50 border-rose-400 text-rose-700 ring-2 ring-rose-400/20'
                      : 'bg-slate-50 border-slate-200 text-slate-600 hover:bg-slate-100'
                  "
                  @click="setDiscountType('NOMINAL')"
                >
                  Potongan Tetap (Rp)
                </button>
              </div>
            </div>

            <!-- Input Nilai Sesuai Tipe -->
            <div v-if="discountForm.type === 'PERCENTAGE'">
              <label class="block text-xs font-bold text-rose-700 mb-1">Nilai Diskon Persen (%) *</label>
              <div class="relative">
                <input
                  v-model.number="discountForm.percentValue"
                  type="number"
                  min="1"
                  max="100"
                  required
                  placeholder="Contoh: 10"
                  class="w-full pr-8 px-3 py-2 text-xs rounded-xl border border-rose-300 bg-rose-50/30 font-black text-rose-700 focus:bg-white focus:ring-2 focus:ring-rose-500"
                />
                <span class="absolute right-3 top-2.5 font-bold text-xs text-rose-500">%</span>
              </div>
            </div>

            <div v-else>
              <label class="block text-xs font-bold text-rose-700 mb-1">Nilai Potongan Tetap (Rp) *</label>
              <div class="relative">
                <span class="absolute left-3 top-2 font-bold text-xs text-rose-500">Rp</span>
                <input
                  v-model.number="discountForm.nominalValue"
                  type="number"
                  min="500"
                  step="500"
                  required
                  placeholder="Contoh: 10000"
                  class="w-full pl-9 pr-3 py-2 text-xs rounded-xl border border-rose-300 bg-rose-50/30 font-black text-rose-700 focus:bg-white focus:ring-2 focus:ring-rose-500"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Minimal Order (Rp)</label>
              <input
                v-model.number="discountForm.min_order_amount"
                type="number"
                min="0"
                step="1000"
                placeholder="0 (Tanpa minimum)"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-rose-500"
              />
            </div>

            <div class="flex items-center gap-2 pt-1">
              <input
                v-model="discountForm.is_active"
                type="checkbox"
                id="disc-active"
                class="rounded text-rose-600 focus:ring-rose-500 w-4 h-4 cursor-pointer"
              />
              <label for="disc-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Promo Aktif
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showDiscountModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Promo' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 5. Modal Tambah/Edit Kategori (Pengeluaran / Layanan) -->
      <div v-if="showCategoryModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <FolderTree class="w-4 h-4 text-amber-600" />
              <span>{{ categoryModalTarget === 'SERVICE' ? 'Tambah Kategori Layanan' : (isEditingCategory ? 'Edit Kategori Pengeluaran' : 'Tambah Kategori Pengeluaran') }}</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showCategoryModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>
          <form class="space-y-3" @submit.prevent="submitCategoryForm">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Kategori *</label>
              <input
                v-model="categoryForm.name"
                type="text"
                required
                placeholder="Contoh: Perlengkapan Bayi / Plastik"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-amber-500"
                @input="handleCategoryNameInput"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">
                Kode Kategori <span class="text-[10px] text-slate-400 font-normal">(Otomatis oleh Sistem)</span>
              </label>
              <input
                v-model="categoryForm.code"
                type="text"
                required
                readonly
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-mono font-bold uppercase text-slate-600 cursor-not-allowed select-none"
              />
            </div>
            <div v-if="categoryModalTarget === 'SERVICE'">
              <label class="block text-xs font-bold text-slate-700 mb-1">Satuan Default Layanan</label>
              <input
                v-model="categoryForm.default_unit"
                type="text"
                placeholder="Kg / Tabung / Pcs / Meter"
                class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-semibold focus:bg-white focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div class="flex items-center gap-2 pt-1">
              <input
                v-model="categoryForm.is_active"
                type="checkbox"
                id="cat-active"
                class="rounded text-amber-600 focus:ring-amber-500 w-4 h-4 cursor-pointer"
              />
              <label for="cat-active" class="text-xs font-semibold text-slate-700 cursor-pointer">
                Kategori Aktif
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showCategoryModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-amber-600 hover:bg-amber-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Kategori' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 6. Modal Koreksi Pengeluaran (created_by Strictly Read-Only) -->
      <div v-if="showEditExpenseModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <DollarSign class="w-4 h-4 text-rose-600" />
              <span>Koreksi Data Pengeluaran Kas Kecil</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showEditExpenseModal = false">
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
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 font-bold focus:bg-white focus:ring-2 focus:ring-rose-500 cursor-pointer"
                >
                  <option v-for="cat in expenseCategoriesList" :key="cat.id" :value="cat.code || cat.category_code">
                    {{ cat.name || cat.category_name }}
                  </option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nominal Pengeluaran (Rp) *</label>
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
            <!-- created_by: STRICTLY READ-ONLY BADGE (Cannot be edited) -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1 flex items-center justify-between">
                <span>Dicatat Oleh</span>
                <span class="text-[10px] text-slate-400 font-normal">Terkunci (Audit Keamanan)</span>
              </label>
              <div class="px-3 py-2 rounded-xl bg-slate-100 border border-slate-200 text-xs font-bold text-slate-600 select-none flex items-center gap-1.5 cursor-not-allowed">
                <Shield class="w-3.5 h-3.5 text-slate-400" />
                <span>{{ editExpenseForm.created_by || 'Kasir Utama' }}</span>
              </div>
            </div>
            <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
              <button
                type="button"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showEditExpenseModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Perbarui Pengeluaran' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 7. Modal Tambah/Edit Pelanggan -->
      <div v-if="showCustomerModal" class="fixed inset-0 z-50 flex items-center justify-center p-3 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-black text-sm text-slate-900 flex items-center gap-2">
              <Users class="w-4 h-4 text-brand-600" />
              <span>{{ isEditingCustomer ? 'Edit Data Pelanggan' : 'Tambah Pelanggan Baru' }}</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showCustomerModal = false">
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
              <label class="block text-xs font-bold text-slate-700 mb-1">Nomor Handphone / WhatsApp *</label>
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
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
                @click="showCustomerModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="modalSaving"
                class="px-4 py-1.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer"
              >
                {{ modalSaving ? 'Menyimpan...' : 'Simpan Pelanggan' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 8. UNIVERSAL DELETE CONFIRMATION MODAL (Clean Vue Modal Replacing Alert JS) -->
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in">
        <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-bold text-sm text-rose-600 flex items-center gap-2">
              <AlertTriangle class="w-4 h-4 text-rose-600" />
              <span>Konfirmasi Hapus Data</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showDeleteModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>
          <p class="text-xs text-slate-600 leading-relaxed">
            Apakah Anda yakin ingin menghapus {{ deleteModalPayload.type }}
            <strong class="text-slate-900">"{{ deleteModalPayload.title }}"</strong>?
          </p>
          <div class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-[11px] text-rose-800 space-y-0.5">
            <span class="font-bold block">Peringatan:</span>
            <p>Data yang dihapus akan dinonaktifkan dari sistem LaundryOS.</p>
          </div>
          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-100">
            <button
              type="button"
              class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
              @click="showDeleteModal = false"
            >
              Batal
            </button>
            <button
              type="button"
              :disabled="deleteLoading"
              class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 disabled:opacity-50 text-white text-xs font-bold transition-all shadow-xs cursor-pointer flex items-center gap-1.5"
              @click="executeConfirmedDelete"
            >
              <Trash2 class="w-3.5 h-3.5" />
              <span>{{ deleteLoading ? 'Menghapus...' : 'Ya, Hapus' }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 9. MODAL ZOOM FOTO ABSENSI -->
      <div
        v-if="showPhotoZoomModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/70 backdrop-blur-xs animate-in fade-in"
        @click.self="showPhotoZoomModal = false"
      >
        <div class="bg-white rounded-3xl max-w-sm w-full p-4 shadow-2xl border border-slate-100 space-y-3">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h4 class="font-black text-sm text-slate-900">{{ photoZoomData.title }}</h4>
              <p class="text-[10px] text-slate-500">{{ photoZoomData.time }} • {{ photoZoomData.distance }}</p>
            </div>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showPhotoZoomModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>
          <div class="rounded-2xl overflow-hidden border border-slate-200 aspect-square bg-slate-950">
            <img :src="photoZoomData.url" class="w-full h-full object-cover" />
          </div>
          <button
            type="button"
            class="w-full py-2 rounded-xl bg-slate-100 text-slate-700 text-xs font-bold hover:bg-slate-200 cursor-pointer transition-colors"
            @click="showPhotoZoomModal = false"
          >
            Tutup Preview Foto
          </button>
        </div>
      </div>

      <!-- 10. MODAL TOLAK CUTI DENGAN ALASAN -->
      <div
        v-if="showRejectModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
        @click.self="showRejectModal = false"
      >
        <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <h4 class="font-bold text-sm text-rose-600 flex items-center gap-2">
              <CalendarX class="w-4 h-4 text-rose-600" />
              <span>Tolak Pengajuan Libur</span>
            </h4>
            <button class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer" @click="showRejectModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>
          <p class="text-xs text-slate-600">
            Tolak pengajuan libur <strong>{{ selectedLeaveForReject?.employee_name }}</strong> pada tanggal <strong>{{ formatDate(selectedLeaveForReject?.date) }}</strong>?
          </p>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Catatan / Alasan Penolakan</label>
            <textarea
              v-model="rejectReasonInput"
              rows="2"
              placeholder="Contoh: Jadwal shift bentrok dengan rekan kerja..."
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-rose-500 resize-none font-medium"
            ></textarea>
          </div>
          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-100">
            <button
              type="button"
              class="px-3.5 py-1.5 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 cursor-pointer"
              @click="showRejectModal = false"
            >
              Batal
            </button>
            <button
              type="button"
              class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-xs cursor-pointer flex items-center gap-1.5"
              @click="handleConfirmRejectLeave"
            >
              <CalendarX class="w-3.5 h-3.5" />
              <span>Konfirmasi Tolak</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import AppLayout from '../components/layout/AppLayout.vue'
import PaginationBar from '../components/ui/PaginationBar.vue'
import LocationPickerMap from '../components/common/LocationPickerMap.vue'
import { useAuth } from '../composables/useAuth'
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
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  FolderTree,
  Clock,
  CalendarCheck,
  CalendarX,
  MapPin,
  Camera,
  Eye,
  Store,
  Crown,
  MessageSquare,
  QrCode,
  Copy,
  RefreshCw,
} from 'lucide-vue-next'
import { useAppSettings } from '../composables/useAppSettings'
import { useToast } from '../composables/useToast'
import { useActiveOutlet } from '../composables/useActiveOutlet'
import { api, formatRupiah, formatDate } from '../api/client'
import { printThermalReceipt } from '../utils/printReceipt'

const { settings, changePin, resetToDefaults } = useAppSettings()
const { success, error, info } = useToast()
const auth = useAuth()
const { activeOutlet, activeOutletId, setActiveOutlet } = useActiveOutlet()

const activeTab = ref('OUTLETS')

const tabGroups = [
  {
    title: 'Cabang & Outlet',
    tabs: [
      { id: 'OUTLETS', label: 'Master Cabang (GPS)', icon: Store },
      { id: 'WA_GATEWAY', label: 'WhatsApp Bot & Cabang', icon: MessageSquare },
      { id: 'PRINTER', label: 'Profil Outlet & Printer', icon: Printer },
    ],
  },
  {
    title: 'Operasional & Shift',
    tabs: [
      { id: 'SHIFTS_ATTENDANCE', label: 'Master Shift & Jam Kerja', icon: Clock },
      { id: 'MACHINE', label: 'Kapasitas Mesin & Toleransi', icon: Scale },
    ],
  },
  {
    title: 'Kepegawaian & Absensi',
    tabs: [
      { id: 'ATTENDANCE_LOG', label: 'Rekap Absensi Karyawan', icon: Camera },
      { id: 'LEAVE_APPROVAL', label: 'Persetujuan Cuti (Libur)', icon: CalendarCheck },
      { id: 'EMPLOYEES', label: 'Pegawai & Hak Akses', icon: ShieldCheck },
    ],
  },
  {
    title: 'Katalog & Keuangan',
    tabs: [
      { id: 'SERVICES', label: 'Master Layanan & Tarif', icon: Tag },
      { id: 'PERFUME_DISCOUNT', label: 'Master Parfum & Diskon', icon: Sparkles },
      { id: 'EXPENSE_CATEGORIES', label: 'Kategori Pengeluaran', icon: Receipt },
      { id: 'EXPENSE_CORRECTION', label: 'Pengoreksi Kas Kecil', icon: DollarSign },
      { id: 'CUSTOMERS', label: 'Kelola Pelanggan (CRM)', icon: Users },
    ],
  },
  {
    title: 'Sistem',
    tabs: [
      { id: 'INTEGRATION', label: 'Status Integrasi & Sistem', icon: Smartphone },
    ],
  },
]

const tabs = computed(() => tabGroups.flatMap((g) => g.tabs))

const activeTabObj = computed(() => {
  return tabs.value.find((t) => t.id === activeTab.value) || tabs.value[0]
})

// State for PIN Form
const oldPin = ref('')
const newPin = ref('')

// State for Simulation
const simulatedWeight = ref(8.5)

// Master Data Lists from API
const servicesList = ref([])
const serviceCategoriesList = ref([])
const perfumesList = ref([])
const discountsList = ref([])
const expenseCategoriesList = ref([])
const employeesList = ref([])

// Sorting state maps
const sortStates = reactive({})

const getSortIcon = (table, col) => {
  const current = sortStates[table]
  if (!current || current.col !== col) return ArrowUpDown
  return current.dir === 'asc' ? ArrowUp : ArrowDown
}

const toggleSort = async (table, col) => {
  if (!sortStates[table] || sortStates[table].col !== col) {
    sortStates[table] = { col, dir: 'asc' }
  } else if (sortStates[table].dir === 'asc') {
    sortStates[table] = { col, dir: 'desc' }
  } else {
    delete sortStates[table]
  }

  // Server-side backend sorting logic
  if (table === 'customer') {
    await loadCustomers()
  } else if (table === 'expense_audit') {
    await loadExpensesAudit()
  } else if (table === 'category') {
    await fetchExpenseCategoriesMaster()
  } else if (table.startsWith('service')) {
    await fetchServicesMaster()
  }
}

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

// Universal Delete Confirmation Modal State
const showDeleteModal = ref(false)
const deleteLoading = ref(false)
const deleteModalPayload = reactive({
  type: '',
  title: '',
  callback: null,
})

const triggerDeleteModal = (type, title, callback) => {
  deleteModalPayload.type = type
  deleteModalPayload.title = title
  deleteModalPayload.callback = callback
  showDeleteModal.value = true
}

const executeConfirmedDelete = async () => {
  if (!deleteModalPayload.callback) return
  deleteLoading.value = true
  try {
    await deleteModalPayload.callback()
    showDeleteModal.value = false
  } catch (err) {
    error(err.message || 'Gagal menghapus data')
  } finally {
    deleteLoading.value = false
  }
}

// 1. Employee Modal State & Granular ACL
const defaultAcl = {
  can_create_order: true,
  can_backdate_order: false,
  can_edit_order: false,
  can_cancel_order: false,
  can_create_expense: true,
  can_process_wash: true,
  can_process_dry: true,
  can_process_iron: true,
  can_process_pack: true,
  can_manage_services: false,
  can_manage_customers: true,
  can_manage_employees: false,
  can_view_financials: false,
  can_view_order_reports: false,
  can_view_financial_reports: false,
  can_view_performance_reports: false,
  can_view_customer_reports: false,
}

const aclPermissionsList = [
  {
    group: 'Transaksi & Order POS',
    permissions: [
      { key: 'can_create_order', label: 'Membuat Order / Transaksi' },
      { key: 'can_backdate_order', label: 'Membuat Order Backdate (Tanggal Mundur)' },
      { key: 'can_edit_order', label: 'Edit Order / Transaksi' },
      { key: 'can_cancel_order', label: 'Membatalkan Order / Transaksi' },
    ],
  },
  {
    group: 'Operasional & Produksi',
    permissions: [
      { key: 'can_create_expense', label: 'Membuat Pengeluaran' },
      { key: 'can_process_wash', label: 'Cuci (Khusus Paket Cuci)' },
      { key: 'can_process_dry', label: 'Kering (Khusus Paket Kering)' },
      { key: 'can_process_iron', label: 'Setrika (Khusus Paket Setrika)' },
      { key: 'can_process_pack', label: 'Packing (Khusus Paket Packing)' },
    ],
  },
  {
    group: 'Manajemen & Data',
    permissions: [
      { key: 'can_manage_services', label: 'Mengelola Layanan / Produk' },
      { key: 'can_manage_customers', label: 'Mengelola Data Pelanggan' },
      { key: 'can_manage_employees', label: 'Mengelola Data Pegawai' },
    ],
  },
  {
    group: 'Keuangan & Laporan',
    permissions: [
      { key: 'can_view_financials', label: 'Menampilkan Nilai Omset' },
      { key: 'can_view_order_reports', label: 'Akses Laporan Transaksi' },
      { key: 'can_view_financial_reports', label: 'Akses Laporan Keuangan' },
      { key: 'can_view_performance_reports', label: 'Akses Laporan Kinerja' },
      { key: 'can_view_customer_reports', label: 'Akses Laporan Pelanggan' },
    ],
  },
]

const showEmployeeModal = ref(false)
const isEditingEmployee = ref(false)
const currentEmployeeId = ref(null)
const employeeForm = reactive({
  username: '',
  full_name: '',
  phone: '',
  email: '',
  role: 'KASIR',
  password: '',
  outlet_id: '',
  permissions: { ...defaultAcl },
  is_active: true,
})

const sortEmployee = (col) => toggleSort('employee', col)
const sortedEmployees = computed(() => {
  const list = [...employeesList.value]
  const s = sortStates['employee']
  if (!s) return list
  return list.sort((a, b) => {
    const valA = (a[s.col] || '').toString().toLowerCase()
    const valB = (b[s.col] || '').toString().toLowerCase()
    return s.dir === 'asc' ? valA.localeCompare(valB) : valB.localeCompare(valA)
  })
})

// 2. Service Modal State & Category Flow
const showServiceModal = ref(false)
const isEditingService = ref(false)
const currentServiceId = ref(null)
const serviceForm = reactive({
  category_id: '',
  category_name: '',
  service_code: '',
  service_name: '',
  unit: 'Kg',
  price: 10000,
  duration_day: 1,
  is_active: true,
})

const generateServiceCode = (catName) => {
  if (!catName) return 'SRV-' + Math.floor(100 + Math.random() * 900)
  const clean = catName.toUpperCase()
  let prefix = 'SRV'
  if (clean.includes('SELF SERVICE') || clean.includes('TABUNG')) prefix = 'SLF'
  else if (clean.includes('SETRIKA') && clean.includes('CUCI')) prefix = 'CS'
  else if (clean.includes('SETRIKA')) prefix = 'ST'
  else if (clean.includes('KARPET')) prefix = 'KRP'
  else if (clean.includes('BONEKA')) prefix = 'BNK'
  else if (clean.includes('SEPATU') || clean.includes('TAS')) prefix = 'SPT'
  else if (clean.includes('BED COVER') || clean.includes('SELIMUT')) prefix = 'BC'
  else if (clean.includes('SATUAN')) prefix = 'SAT'

  const count = servicesList.value.filter((s) => (s.service_code || '').startsWith(prefix)).length + 1
  return `${prefix}-${String(count).padStart(2, '0')}`
}

const handleCategorySelectionChange = () => {
  const selectedCat = serviceCategoriesList.value.find((c) => c.id === serviceForm.category_id)
  if (selectedCat) {
    serviceForm.category_name = selectedCat.name
    serviceForm.unit = selectedCat.default_unit || 'Kg'
    if (!isEditingService.value) {
      serviceForm.service_code = generateServiceCode(selectedCat.name)
    }
  }
}

// Grouped Services computation for categorized catalog
const sortService = (catName, col) => toggleSort('service_' + catName, col)

const serviceGroups = computed(() => {
  const groupsMap = {}

  // Initialize from serviceCategoriesList to keep order
  serviceCategoriesList.value.forEach((cat) => {
    groupsMap[cat.name] = {
      categoryName: cat.name,
      defaultUnit: cat.default_unit || 'Kg',
      services: [],
    }
  })

  // Distribute services
  servicesList.value.forEach((s) => {
    const grpName = s.category_name || s.category || 'Layanan Lain-lain'
    if (!groupsMap[grpName]) {
      groupsMap[grpName] = {
        categoryName: grpName,
        defaultUnit: s.unit || 'Kg',
        services: [],
      }
    }
    groupsMap[grpName].services.push(s)
  })

  // Apply sorting within each group
  return Object.values(groupsMap)
    .filter((g) => g.services.length > 0 || serviceCategoriesList.value.some((c) => c.name === g.categoryName))
    .map((g) => {
      const s = sortStates['service_' + g.categoryName]
      if (!s) return g
      const sorted = [...g.services].sort((a, b) => {
        if (s.col === 'price') {
          const pA = Number(a.price || a.price_per_kg) || 0
          const pB = Number(b.price || b.price_per_kg) || 0
          return s.dir === 'asc' ? pA - pB : pB - pA
        }
        const valA = (a[s.col] || '').toString().toLowerCase()
        const valB = (b[s.col] || '').toString().toLowerCase()
        return s.dir === 'asc' ? valA.localeCompare(valB) : valB.localeCompare(valA)
      })
      return { ...g, services: sorted }
    })
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

// 4. Discount Modal State (Mutually Exclusive: Percentage vs Nominal)
const showDiscountModal = ref(false)
const isEditingDiscount = ref(false)
const currentDiscountId = ref(null)
const discountForm = reactive({
  code: '',
  name: '',
  type: 'PERCENTAGE', // 'PERCENTAGE' | 'NOMINAL'
  percentValue: 10,
  nominalValue: 5000,
  min_order_amount: 0,
  is_active: true,
})

const setDiscountType = (t) => {
  discountForm.type = t
}

const handleDiscountNameInput = () => {
  if (!isEditingDiscount.value && !discountForm.code) {
    const clean = discountForm.name.replace(/[^a-zA-Z0-9]/g, '').toUpperCase().slice(0, 6)
    discountForm.code = clean ? `PRM-${clean}` : 'PRM-' + Math.floor(100 + Math.random() * 900)
  }
}

// 5. Category Modal State (Supports either Expense or Service Category)
const showCategoryModal = ref(false)
const categoryModalTarget = ref('EXPENSE') // 'EXPENSE' | 'SERVICE'
const isEditingCategory = ref(false)
const currentCategoryId = ref(null)
const categoryForm = reactive({
  code: '',
  name: '',
  default_unit: 'Kg',
  is_active: true,
})

const handleCategoryNameInput = () => {
  if (!isEditingCategory.value) {
    categoryForm.code = categoryForm.name
      .trim()
      .toUpperCase()
      .replace(/[^A-Z0-9]+/g, '_')
      .replace(/^_+|_+$/g, '')
  }
}

const sortCategory = (col) => toggleSort('category', col)
const sortedExpenseCategories = computed(() => expenseCategoriesList.value)

// 6. Edit Expense Modal State (created_by strictly read-only)
const showEditExpenseModal = ref(false)
const currentExpenseId = ref(null)
const editExpenseForm = reactive({
  expense_date: '',
  category: '',
  amount: 0,
  description: '',
  created_by: '',
})

const sortExpenseAudit = (col) => toggleSort('expense_audit', col)
const sortedExpensesAudit = computed(() => expensesAuditList.value)

// 7. Customer Modal State
const showCustomerModal = ref(false)
const isEditingCustomer = ref(false)
const currentCustomerId = ref(null)
const customerForm = reactive({
  full_name: '',
  phone: '',
  address: '',
})

const sortCustomer = (col) => toggleSort('customer', col)
const sortedCustomers = computed(() => customersList.value)

// ==========================================
// API Loaders
// ==========================================
const fetchServicesMaster = async () => {
  try {
    const params = {}
    const srvSort = Object.keys(sortStates).find((k) => k.startsWith('service'))
    if (srvSort && sortStates[srvSort]) {
      params.sort_by = sortStates[srvSort].col
      params.sort_dir = sortStates[srvSort].dir
    }
    const srv = await api.get('/master/services', params).catch(() => [])
    if (Array.isArray(srv)) servicesList.value = srv
  } catch (err) {
    console.error('Failed to load services master:', err)
  }
}

const fetchExpenseCategoriesMaster = async () => {
  try {
    const params = {}
    if (sortStates['category']) {
      params.sort_by = sortStates['category'].col
      params.sort_dir = sortStates['category'].dir
    }
    const expCat = await api.get('/master/expense-categories', params).catch(() => [])
    if (Array.isArray(expCat)) expenseCategoriesList.value = expCat
  } catch (err) {
    console.error('Failed to load expense categories master:', err)
  }
}

const fetchMasterData = async () => {
  try {
    const srvParams = {}
    const srvSort = Object.keys(sortStates).find((k) => k.startsWith('service'))
    if (srvSort && sortStates[srvSort]) {
      srvParams.sort_by = sortStates[srvSort].col
      srvParams.sort_dir = sortStates[srvSort].dir
    }
    const expParams = {}
    if (sortStates['category']) {
      expParams.sort_by = sortStates['category'].col
      expParams.sort_dir = sortStates['category'].dir
    }

    const [srv, srvCat, prf, disc, expCat, emp] = await Promise.all([
      api.get('/master/services', srvParams).catch(() => []),
      api.get('/master/service-categories').catch(() => []),
      api.get('/master/perfumes').catch(() => []),
      api.get('/master/discounts').catch(() => []),
      api.get('/master/expense-categories', expParams).catch(() => []),
      api.get('/master/employees').catch(() => []),
    ])
    if (Array.isArray(srv)) servicesList.value = srv
    if (Array.isArray(srvCat)) serviceCategoriesList.value = srvCat
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
    if (sortStates['expense_audit']) {
      params.sort_by = sortStates['expense_audit'].col
      params.sort_dir = sortStates['expense_audit'].dir
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
    }
    const q = customerSearchQuery.value.trim()
    if (q.length >= 3) {
      params.search = q
    }
    if (sortStates['customer']) {
      params.sort_by = sortStates['customer'].col
      params.sort_dir = sortStates['customer'].dir
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
  if (tab === 'OUTLETS') {
    loadOutlets()
  } else if (tab === 'EXPENSE_CORRECTION') {
    loadExpensesAudit()
  } else if (tab === 'CUSTOMERS') {
    loadCustomers()
  } else if (tab === 'SHIFTS_ATTENDANCE') {
    loadShifts()
  } else if (tab === 'ATTENDANCE_LOG') {
    loadAttendances()
  } else if (tab === 'LEAVE_APPROVAL') {
    loadLeaves()
  } else if (tab === 'WA_GATEWAY') {
    loadGatewayStatusAndGroups()
  }
})

watch(
  [() => expenseAuditFilter.startDate, () => expenseAuditFilter.endDate, () => expenseAuditFilter.category],
  () => {
    pageExpensesAudit.value = 1
    loadExpensesAudit()
  }
)

watch([pageExpensesAudit, limitExpensesAudit], () => {
  loadExpensesAudit()
})

// 3-character threshold filter watcher for customer search
let customerSearchTimer = null
watch(customerSearchQuery, (val) => {
  const trimmed = (val || '').trim()
  if (customerSearchTimer) clearTimeout(customerSearchTimer)

  if (trimmed.length === 0) {
    pageCustomers.value = 1
    loadCustomers()
  } else if (trimmed.length >= 3) {
    customerSearchTimer = setTimeout(() => {
      pageCustomers.value = 1
      loadCustomers()
    }, 350)
  }
  // If 1 or 2 characters: DO NOT hit API
})

watch([pageCustomers, limitCustomers], () => {
  loadCustomers()
})

onMounted(() => {
  fetchMasterData()
  loadExpensesAudit()
  loadCustomers()
  loadShifts()
  loadLeaves()
  loadAttendances()
  loadOutlets()
  loadGatewayStatusAndGroups()
})

// ==========================================
// Handlers: Save All & PIN
// ==========================================
const handleSavePrinterSettings = () => {
  success('Konfigurasi profil outlet & template nota thermal berhasil disimpan!')
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
  printThermalReceipt('thermal-preview-receipt', {
    paperWidth: settings.printer?.paperWidth || '58mm',
    title: 'Uji-Coba-Cetak-Struk',
  })
}

const handleFactoryReset = () => {
  triggerDeleteModal(
    'Pengaturan',
    'Reset ke Default Pabrik',
    async () => {
      resetToDefaults()
      success('Pengaturan telah dikembalikan ke default pabrik.')
    }
  )
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
// CRUD: Pegawai (Employees & Granular ACL)
// ==========================================
const openAddEmployeeModal = () => {
  isEditingEmployee.value = false
  currentEmployeeId.value = null
  employeeForm.username = ''
  employeeForm.full_name = ''
  employeeForm.phone = ''
  employeeForm.email = ''
  employeeForm.role = 'KASIR'
  employeeForm.password = '123456'
  employeeForm.outlet_id = outletsList.value[0]?.id || ''
  employeeForm.permissions = { ...defaultAcl }
  employeeForm.is_active = true
  showEmployeeModal.value = true
}

const openEditEmployeeModal = (emp) => {
  isEditingEmployee.value = true
  currentEmployeeId.value = emp.id
  employeeForm.username = emp.username
  employeeForm.full_name = emp.full_name
  employeeForm.phone = emp.phone || ''
  employeeForm.email = emp.email || ''
  employeeForm.role = emp.role
  employeeForm.password = ''
  employeeForm.outlet_id = emp.outlet_id || (outletsList.value[0]?.id || '')
  employeeForm.permissions = { ...defaultAcl, ...(emp.permissions || {}) }
  employeeForm.is_active = emp.is_active
  showEmployeeModal.value = true
}

const submitEmployeeForm = async () => {
  modalSaving.value = true
  try {
    const payload = {
      username: employeeForm.username,
      full_name: employeeForm.full_name,
      phone: employeeForm.phone,
      email: employeeForm.email,
      role: employeeForm.role,
      password: employeeForm.password,
      outlet_id: employeeForm.outlet_id,
      permissions: employeeForm.permissions,
      is_active: employeeForm.is_active,
    }
    if (isEditingEmployee.value) {
      await api.put(`/master/employees/${currentEmployeeId.value}`, payload)
      success('Data pegawai & hak akses berhasil diperbarui!')
    } else {
      await api.post('/master/employees', payload)
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

const deleteEmployee = async (id) => {
  await api.delete(`/master/employees/${id}`)
  success('Pegawai berhasil dihapus!')
  employeesList.value = employeesList.value.filter((e) => e.id !== id)
}

// ==========================================
// CRUD: Layanan (Services)
// ==========================================
const openAddServiceModal = () => {
  isEditingService.value = false
  currentServiceId.value = null

  // Default to first service category
  const firstCat = serviceCategoriesList.value[0]
  serviceForm.category_id = firstCat ? firstCat.id : ''
  serviceForm.category_name = firstCat ? firstCat.name : 'Cuci Setrika Reguler'
  serviceForm.unit = firstCat ? (firstCat.default_unit || 'Kg') : 'Kg'
  serviceForm.service_code = generateServiceCode(serviceForm.category_name)
  serviceForm.service_name = ''
  serviceForm.price = 10000
  serviceForm.duration_day = 1
  serviceForm.is_active = true
  showServiceModal.value = true
}

const openEditServiceModal = (s) => {
  isEditingService.value = true
  currentServiceId.value = s.id
  serviceForm.category_id = s.category_id || ''
  serviceForm.category_name = s.category_name || s.category || ''
  serviceForm.service_code = s.service_code
  serviceForm.service_name = s.service_name
  serviceForm.unit = s.unit || 'Kg'
  serviceForm.price = Number(s.price || s.price_per_kg) || 0
  serviceForm.duration_day = Number(s.duration_day) || 1
  serviceForm.is_active = s.is_active !== false && s.active !== false
  showServiceModal.value = true
}

const submitServiceForm = async () => {
  modalSaving.value = true
  try {
    const payload = {
      service_code: serviceForm.service_code,
      service_name: serviceForm.service_name,
      category_id: serviceForm.category_id,
      category_name: serviceForm.category_name,
      category: serviceForm.category_name,
      unit: serviceForm.unit,
      price: serviceForm.price,
      duration_day: serviceForm.duration_day,
      active: serviceForm.is_active,
      is_active: serviceForm.is_active,
    }
    if (isEditingService.value) {
      await api.put(`/master/services/${currentServiceId.value}`, payload)
      success('Layanan berhasil diperbarui!')
    } else {
      await api.post('/master/services', payload)
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

const deleteService = async (id) => {
  await api.delete(`/master/services/${id}`)
  success('Layanan berhasil dihapus!')
  servicesList.value = servicesList.value.filter((item) => item.id !== id)
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

const deletePerfume = async (id) => {
  await api.delete(`/master/perfumes/${id}`)
  success('Parfum berhasil dihapus!')
  perfumesList.value = perfumesList.value.filter((item) => item.id !== id)
}

// ==========================================
// CRUD: Diskon & Promo (Mutually Exclusive)
// ==========================================
const openAddDiscountModal = () => {
  isEditingDiscount.value = false
  currentDiscountId.value = null
  discountForm.code = 'PRM-' + Math.floor(100 + Math.random() * 900)
  discountForm.name = ''
  discountForm.type = 'PERCENTAGE'
  discountForm.percentValue = 10
  discountForm.nominalValue = 5000
  discountForm.min_order_amount = 0
  discountForm.is_active = true
  showDiscountModal.value = true
}

const openEditDiscountModal = (d) => {
  isEditingDiscount.value = true
  currentDiscountId.value = d.id
  discountForm.code = d.code || 'PRM-' + Math.floor(100 + Math.random() * 900)
  discountForm.name = d.name

  const isPercent = d.discount_type === 'PERCENTAGE' || Number(d.discount_percent) > 0
  discountForm.type = isPercent ? 'PERCENTAGE' : 'NOMINAL'
  discountForm.percentValue = isPercent ? Number(d.value || d.discount_percent || 10) : 10
  discountForm.nominalValue = !isPercent ? Number(d.value || d.fixed_amount || 5000) : 5000
  discountForm.min_order_amount = Number(d.min_order_amount) || 0
  discountForm.is_active = d.is_active !== false
  showDiscountModal.value = true
}

const submitDiscountForm = async () => {
  modalSaving.value = true
  try {
    const isPercent = discountForm.type === 'PERCENTAGE'
    const value = isPercent ? discountForm.percentValue : discountForm.nominalValue

    const payload = {
      code: discountForm.code,
      name: discountForm.name,
      discount_type: discountForm.type,
      value: value,
      discount_percent: isPercent ? value : 0,
      fixed_amount: !isPercent ? value : 0,
      min_order_amount: discountForm.min_order_amount || 0,
      is_active: discountForm.is_active,
    }

    if (isEditingDiscount.value) {
      await api.put(`/master/discounts/${currentDiscountId.value}`, payload)
      success('Promo diskon berhasil diperbarui!')
    } else {
      await api.post('/master/discounts', payload)
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

const deleteDiscount = async (id) => {
  await api.delete(`/master/discounts/${id}`)
  success('Diskon berhasil dihapus!')
  discountsList.value = discountsList.value.filter((item) => item.id !== id)
}

// ==========================================
// CRUD: Kategori (Expense & Service)
// ==========================================
const openAddCategoryModal = (target = 'EXPENSE') => {
  categoryModalTarget.value = target
  isEditingCategory.value = false
  currentCategoryId.value = null
  categoryForm.name = ''
  categoryForm.code = ''
  categoryForm.default_unit = 'Kg'
  categoryForm.is_active = true
  showCategoryModal.value = true
}

const openEditCategoryModal = (cat) => {
  categoryModalTarget.value = 'EXPENSE'
  isEditingCategory.value = true
  currentCategoryId.value = cat.id
  categoryForm.code = cat.code || cat.category_code
  categoryForm.name = cat.name || cat.category_name
  categoryForm.is_active = cat.is_active !== false
  showCategoryModal.value = true
}

const submitCategoryForm = async () => {
  modalSaving.value = true
  try {
    if (categoryModalTarget.value === 'SERVICE') {
      await api.post('/master/service-categories', {
        code: categoryForm.code,
        name: categoryForm.name,
        default_unit: categoryForm.default_unit || 'Kg',
        display_order: serviceCategoriesList.value.length + 1,
        is_active: categoryForm.is_active,
      })
      success('Kategori layanan baru berhasil ditambahkan!')
      const srvCat = await api.get('/master/service-categories')
      if (Array.isArray(srvCat)) serviceCategoriesList.value = srvCat
    } else {
      if (isEditingCategory.value) {
        await api.put(`/master/expense-categories/${currentCategoryId.value}`, {
          code: categoryForm.code,
          name: categoryForm.name,
          is_active: categoryForm.is_active,
        })
        success('Kategori pengeluaran berhasil diperbarui!')
      } else {
        await api.post('/master/expense-categories', {
          code: categoryForm.code,
          name: categoryForm.name,
          is_active: categoryForm.is_active,
        })
        success('Kategori pengeluaran baru berhasil ditambahkan!')
      }
      const expCat = await api.get('/master/expense-categories')
      if (Array.isArray(expCat)) expenseCategoriesList.value = expCat
    }
    showCategoryModal.value = false
  } catch (err) {
    error(err.message || 'Gagal menyimpan kategori')
  } finally {
    modalSaving.value = false
  }
}

const deleteCategory = async (id) => {
  await api.delete(`/master/expense-categories/${id}`)
  success('Kategori berhasil dihapus!')
  expenseCategoriesList.value = expenseCategoriesList.value.filter((item) => item.id !== id)
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

const deleteExpense = async (id) => {
  await api.delete(`/expenses/${id}`)
  success('Pengeluaran berhasil dihapus!')
  loadExpensesAudit()
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

const deleteCustomer = async (id) => {
  await api.delete(`/customers/${id}`)
  success('Pelanggan berhasil dihapus!')
  loadCustomers()
}

// ==========================================
// Milestone: Outlets, Shifts, Attendance & Leaves
// ==========================================
const outletsList = ref([])
const showOutletModal = ref(false)
const isEditingOutlet = ref(false)
const outletForm = reactive({
  id: '',
  name: '',
  code: '',
  address: '',
  phone: '',
  latitude: -6.113425,
  longitude: 106.914230,
  geofence_radius_meters: 25,
  whatsapp_group_id: '',
  whatsapp_group_name: '',
  is_active: true,
})

const loadOutlets = async () => {
  try {
    const data = await api.get('/master/outlets')
    outletsList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Failed to load outlets:', err)
  }
}

const openAddOutletModal = () => {
  isEditingOutlet.value = false
  Object.assign(outletForm, {
    id: '',
    name: '',
    code: 'OUT-' + Math.floor(100 + Math.random() * 900),
    address: '',
    phone: '',
    latitude: -6.113425,
    longitude: 106.914230,
    geofence_radius_meters: 25,
    whatsapp_group_id: '',
    whatsapp_group_name: '',
    is_active: true,
  })
  showOutletModal.value = true
}

const openEditOutletModal = (outlet) => {
  isEditingOutlet.value = true
  Object.assign(outletForm, {
    id: outlet.id,
    name: outlet.name,
    code: outlet.code,
    address: outlet.address,
    phone: outlet.phone,
    latitude: outlet.latitude || -6.113425,
    longitude: outlet.longitude || 106.914230,
    geofence_radius_meters: outlet.geofence_radius_meters || 25,
    whatsapp_group_id: outlet.whatsapp_group_id || '',
    whatsapp_group_name: outlet.whatsapp_group_name || '',
    is_active: outlet.is_active,
  })
  showOutletModal.value = true
}

const submitOutletForm = async () => {
  modalSaving.value = true
  try {
    const payload = {
      name: outletForm.name,
      code: outletForm.code,
      address: outletForm.address,
      phone: outletForm.phone,
      latitude: parseFloat(outletForm.latitude) || -6.113425,
      longitude: parseFloat(outletForm.longitude) || 106.914230,
      geofence_radius_meters: parseInt(outletForm.geofence_radius_meters) || 25,
      whatsapp_group_id: outletForm.whatsapp_group_id || '',
      whatsapp_group_name: outletForm.whatsapp_group_name || '',
      is_active: outletForm.is_active,
    }
    if (isEditingOutlet.value) {
      await api.put(`/master/outlets/${outletForm.id}`, payload)
      success('Profil & koordinat GPS outlet berhasil diperbarui!')
    } else {
      await api.post('/master/outlets', payload)
      success('Outlet cabang baru berhasil ditambahkan!')
    }
    showOutletModal.value = false
    await loadOutlets()
  } catch (err) {
    error(err.message || 'Gagal menyimpan outlet')
  } finally {
    modalSaving.value = false
  }
}

const deleteOutlet = async (id, name) => {
  triggerDeleteModal(
    'Outlet Cabang',
    name,
    async () => {
      await api.delete(`/master/outlets/${id}`)
      success(`Outlet ${name} berhasil dinonaktifkan!`)
      await loadOutlets()
    }
  )
}

const shiftsList = ref([])
const attendanceList = ref([])
const attendanceFilter = reactive({
  employee: '',
  startDate: '',
  endDate: '',
})
const loadingAttendances = ref(false)

const leavesList = ref([])
const leaveFilterStatus = ref('')
const loadingLeaves = ref(false)
const showRejectModal = ref(false)
const selectedLeaveForReject = ref(null)
const rejectReasonInput = ref('')

const showPhotoZoomModal = ref(false)
const photoZoomData = ref({ url: '', title: '', time: '', distance: '' })

const pendingLeavesCount = computed(() => {
  return leavesList.value.filter((l) => l.status === 'PENDING').length
})

const loadShifts = async () => {
  try {
    const data = await api.get('/master/shifts')
    shiftsList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Failed to load shifts:', err)
  }
}

const saveShiftSettings = async (shift) => {
  try {
    await api.put(`/master/shifts/${shift.id}`, {
      shift_name: shift.shift_name || shift.name || ('Shift ' + (shift.code || 'Operasional')),
      start_time: shift.start_time,
      end_time: shift.end_time,
      break_start: shift.break_start,
      break_end: shift.break_end,
      late_tolerance_minutes: parseInt(shift.late_tolerance_minutes) || 15,
      is_active: shift.is_active,
    })
    success(`Pengaturan jam ${shift.shift_name || 'shift'} berhasil disimpan!`)
    await loadShifts()
  } catch (err) {
    error(err.message || 'Gagal menyimpan shift')
  }
}

const loadAttendances = async () => {
  loadingAttendances.value = true
  try {
    const params = {}
    if (attendanceFilter.employee) params.employee = attendanceFilter.employee
    if (attendanceFilter.startDate) params.start_date = attendanceFilter.startDate
    if (attendanceFilter.endDate) params.end_date = attendanceFilter.endDate
    const data = await api.get('/attendances/history', params)
    attendanceList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Failed to load attendances:', err)
  } finally {
    loadingAttendances.value = false
  }
}

const loadLeaves = async () => {
  loadingLeaves.value = true
  try {
    const params = {}
    if (leaveFilterStatus.value) params.status = leaveFilterStatus.value
    const data = await api.get('/leaves', params)
    leavesList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Failed to load leaves:', err)
  } finally {
    loadingLeaves.value = false
  }
}

const handleApproveLeave = async (leave) => {
  try {
    await api.put(`/leaves/${leave.id}/review`, {
      status: 'APPROVED',
      reviewer_notes: 'Disetujui oleh Owner.',
    })
    success(`Pengajuan libur ${leave.employee_name} disetujui! ✅`)
    await loadLeaves()
  } catch (err) {
    error(err.message || 'Gagal menyetujui libur')
  }
}

const openRejectLeaveModal = (leave) => {
  selectedLeaveForReject.value = leave
  rejectReasonInput.value = ''
  showRejectModal.value = true
}

const handleConfirmRejectLeave = async () => {
  if (!selectedLeaveForReject.value) return
  try {
    await api.put(`/leaves/${selectedLeaveForReject.value.id}/review`, {
      status: 'REJECTED',
      reviewer_notes: rejectReasonInput.value || 'Ditolak karena pertimbangan jadwal shift operasional.',
    })
    success(`Pengajuan libur ${selectedLeaveForReject.value.employee_name} ditolak. ❌`)
    showRejectModal.value = false
    selectedLeaveForReject.value = null
    await loadLeaves()
  } catch (err) {
    error(err.message || 'Gagal menolak libur')
  }
}

const openPhotoZoom = (att) => {
  photoZoomData.value = {
    url: att.check_in_photo || att.check_out_photo,
    title: `${att.employee_name} • ${att.shift_name}`,
    time: att.clock_in_time ? formatDate(att.clock_in_time) : '',
    distance: att.distance_meters != null ? `${att.distance_meters.toFixed(1)} meter dari outlet` : '',
  }
  showPhotoZoomModal.value = true
}

// ==========================================
// WhatsApp Gateway & Multi-Outlet Routing
// ==========================================
const loadingGateway = ref(false)
const loadingGroups = ref(false)
const searchGroupQuery = ref('')
const gatewayStatus = ref({
  state: 'DISCONNECTED',
  bot: null,
  qr: null,
  qrDataURL: null,
})
const gatewayGroups = ref([])

const filteredGatewayGroups = computed(() => {
  if (!searchGroupQuery.value.trim()) return gatewayGroups.value
  const q = searchGroupQuery.value.toLowerCase()
  return gatewayGroups.value.filter(
    (g) => (g.name || '').toLowerCase().includes(q) || (g.id || '').toLowerCase().includes(q)
  )
})

const loadGatewayStatusAndGroups = async () => {
  loadingGateway.value = true
  try {
    const statusData = await api.get('/gateway/status')
    if (statusData) {
      gatewayStatus.value = statusData
    }
    if (statusData?.state === 'CONNECTED') {
      loadingGroups.value = true
      try {
        const groupsData = await api.get('/gateway/groups')
        gatewayGroups.value = Array.isArray(groupsData) ? groupsData : []
      } catch (gErr) {
        console.error('Failed to load gateway groups:', gErr)
      } finally {
        loadingGroups.value = false
      }
    }
  } catch (err) {
    console.error('Failed to load gateway status:', err)
  } finally {
    loadingGateway.value = false
  }
}

const assignGroupToActiveOutlet = async (group) => {
  if (!activeOutlet.value) {
    error('Pilih cabang aktif terlebih dahulu di header pengaturan.')
    return
  }
  try {
    await api.put(`/master/outlets/${activeOutlet.value.id}`, {
      name: activeOutlet.value.name,
      code: activeOutlet.value.code,
      address: activeOutlet.value.address,
      phone: activeOutlet.value.phone,
      latitude: activeOutlet.value.latitude,
      longitude: activeOutlet.value.longitude,
      geofence_radius_meters: activeOutlet.value.geofence_radius_meters,
      whatsapp_group_id: group.id,
      whatsapp_group_name: group.name,
      is_active: activeOutlet.value.is_active,
    })
    success(`Cabang "${activeOutlet.value.name}" berhasil dihubungkan ke grup "${group.name}"!`)
    await loadOutlets()
  } catch (err) {
    error(err.message || 'Gagal menghubungkan grup ke cabang')
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    info('ID Grup WhatsApp berhasil disalin ke clipboard!')
  } catch (e) {
    info(`ID: ${text}`)
  }
}
</script>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }
  #thermal-preview-receipt,
  #thermal-preview-receipt * {
    visibility: visible;
  }
  #thermal-preview-receipt {
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    margin: 0;
    padding: 0;
    box-shadow: none !important;
    border: none !important;
  }
}
</style>
